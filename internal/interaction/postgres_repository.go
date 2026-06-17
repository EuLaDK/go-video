package interaction

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

type commentQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// NewPostgresRepository 创建视频互动 PostgreSQL 仓库；pool 为数据库连接池。
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// ListComments 读取评论列表；ctx 为请求上下文，videoID/userID/sort 用于定位和排序。
func (repository *PostgresRepository) ListComments(ctx context.Context, videoID string, userID string, sortMode CommentSort) ([]CommentItem, error) {
	orderBy := "c.created_at_ms DESC, c.id ASC"
	if sortMode == CommentSortHot {
		orderBy = "COUNT(l.user_id) DESC, c.created_at_ms DESC, c.id ASC"
	}

	rows, err := repository.pool.Query(ctx, `
		SELECT
			c.id,
			c.video_id,
			c.content,
			c.author,
			EXISTS (
				SELECT 1
				FROM video_comment_likes current_like
				WHERE current_like.comment_id = c.id AND current_like.user_id = $2
			) AS liked_by_me,
			COUNT(l.user_id) AS likes,
			c.created_at_ms
		FROM video_comments c
		LEFT JOIN video_comment_likes l ON l.comment_id = c.id
		WHERE c.video_id = $1
		GROUP BY c.id, c.video_id, c.content, c.author, c.created_at_ms
		ORDER BY `+orderBy, videoID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []CommentItem{}
	for rows.Next() {
		item, err := scanCommentRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// InsertComment 写入评论；ctx 为请求上下文，userID 为作者，item 为评论数据。
func (repository *PostgresRepository) InsertComment(ctx context.Context, userID string, item CommentItem) error {
	_, err := repository.pool.Exec(ctx, `
		INSERT INTO video_comments (
			id,
			video_id,
			user_id,
			author,
			content,
			created_at_ms,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`, item.ID, item.VideoID, userID, item.Author, item.Content, item.CreatedAt)

	return err
}

// ToggleCommentLike 切换评论点赞；ctx 为请求上下文，videoID/commentID/userID 定位点赞。
func (repository *PostgresRepository) ToggleCommentLike(ctx context.Context, videoID string, commentID string, userID string) (CommentItem, error) {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return CommentItem{}, err
	}
	defer tx.Rollback(ctx)

	var exists bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM video_comments
			WHERE id = $1 AND video_id = $2
		)
	`, commentID, videoID).Scan(&exists); err != nil {
		return CommentItem{}, err
	}
	if !exists {
		return CommentItem{}, ErrCommentNotFound
	}

	var liked bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM video_comment_likes
			WHERE comment_id = $1 AND user_id = $2
		)
	`, commentID, userID).Scan(&liked); err != nil {
		return CommentItem{}, err
	}

	if liked {
		if _, err := tx.Exec(ctx, `
			DELETE FROM video_comment_likes
			WHERE comment_id = $1 AND user_id = $2
		`, commentID, userID); err != nil {
			return CommentItem{}, err
		}
	} else {
		if _, err := tx.Exec(ctx, `
			INSERT INTO video_comment_likes (comment_id, user_id, created_at)
			VALUES ($1, $2, NOW())
			ON CONFLICT (comment_id, user_id) DO NOTHING
		`, commentID, userID); err != nil {
			return CommentItem{}, err
		}
	}

	item, err := queryComment(ctx, tx, videoID, commentID, userID)
	if err != nil {
		return CommentItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return CommentItem{}, err
	}

	return item, nil
}

// DeleteComment 删除当前用户自己的评论；ctx 为请求上下文，videoID/commentID/userID 定位评论。
func (repository *PostgresRepository) DeleteComment(ctx context.Context, videoID string, commentID string, userID string) error {
	tag, err := repository.pool.Exec(ctx, `
		DELETE FROM video_comments
		WHERE id = $1 AND video_id = $2 AND user_id = $3
	`, commentID, videoID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrCommentNotFound
	}

	return nil
}

// ListDanmaku 读取弹幕列表；ctx 为请求上下文，videoID 定位视频。
func (repository *PostgresRepository) ListDanmaku(ctx context.Context, videoID string) ([]DanmakuItem, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT id, video_id, content, color, created_at_ms
		FROM video_danmaku
		WHERE video_id = $1
		ORDER BY created_at_ms DESC, id ASC
	`, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []DanmakuItem{}
	for rows.Next() {
		var item DanmakuItem
		if err := rows.Scan(&item.ID, &item.VideoID, &item.Content, &item.Color, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// InsertDanmaku 写入弹幕；ctx 为请求上下文，userID 为发送用户，item 为弹幕数据。
func (repository *PostgresRepository) InsertDanmaku(ctx context.Context, userID string, item DanmakuItem) error {
	_, err := repository.pool.Exec(ctx, `
		INSERT INTO video_danmaku (
			id,
			video_id,
			user_id,
			content,
			color,
			created_at_ms,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`, item.ID, item.VideoID, userID, item.Content, item.Color, item.CreatedAt)

	return err
}

// queryComment 读取单条评论及当前用户点赞状态；ctx 为请求上下文，querier 为查询器。
func queryComment(ctx context.Context, querier commentQuerier, videoID string, commentID string, userID string) (CommentItem, error) {
	row := querier.QueryRow(ctx, `
		SELECT
			c.id,
			c.video_id,
			c.content,
			c.author,
			EXISTS (
				SELECT 1
				FROM video_comment_likes current_like
				WHERE current_like.comment_id = c.id AND current_like.user_id = $3
			) AS liked_by_me,
			COUNT(l.user_id) AS likes,
			c.created_at_ms
		FROM video_comments c
		LEFT JOIN video_comment_likes l ON l.comment_id = c.id
		WHERE c.video_id = $1 AND c.id = $2
		GROUP BY c.id, c.video_id, c.content, c.author, c.created_at_ms
	`, videoID, commentID, userID)

	item, err := scanCommentRow(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return CommentItem{}, ErrCommentNotFound
	}
	if err != nil {
		return CommentItem{}, err
	}

	return item, nil
}

type commentScanner interface {
	Scan(dest ...any) error
}

// scanCommentRow 扫描单条评论；row 为数据库单行结果。
func scanCommentRow(row commentScanner) (CommentItem, error) {
	var item CommentItem
	var likes int64
	if err := row.Scan(&item.ID, &item.VideoID, &item.Content, &item.Author, &item.LikedByMe, &likes, &item.CreatedAt); err != nil {
		return CommentItem{}, err
	}
	item.Likes = int(likes)

	return item, nil
}

type commentRows interface {
	Scan(dest ...any) error
}

// scanCommentRows 扫描评论列表中的一行；rows 为数据库多行结果。
func scanCommentRows(rows commentRows) (CommentItem, error) {
	return scanCommentRow(rows)
}
