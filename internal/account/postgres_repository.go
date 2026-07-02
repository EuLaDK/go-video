package account

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository 创建账号 PostgreSQL 仓库；pool 为数据库连接池。
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// GetUser 读取用户资料；ctx 为请求上下文，userID 为用户标识。
func (repository *PostgresRepository) GetUser(ctx context.Context, userID string) (UserProfile, error) {
	var profile UserProfile
	err := repository.pool.QueryRow(ctx, `
		SELECT id, avatar_url, email, is_logged_in, is_vip, nickname, password_hash, phone, vip_until
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&profile.ID,
		&profile.AvatarURL,
		&profile.Email,
		&profile.IsLoggedIn,
		&profile.IsVip,
		&profile.Nickname,
		&profile.PasswordHash,
		&profile.Phone,
		&profile.VipUntil,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return UserProfile{}, ErrUserNotFound
	}
	if err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}

// GetUserByEmail 按邮箱读取用户资料；ctx 为请求上下文，email 为规范化邮箱。
func (repository *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (UserProfile, error) {
	var profile UserProfile
	err := repository.pool.QueryRow(ctx, `
		SELECT id, avatar_url, email, is_logged_in, is_vip, nickname, password_hash, phone, vip_until
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`, email).Scan(
		&profile.ID,
		&profile.AvatarURL,
		&profile.Email,
		&profile.IsLoggedIn,
		&profile.IsVip,
		&profile.Nickname,
		&profile.PasswordHash,
		&profile.Phone,
		&profile.VipUntil,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return UserProfile{}, ErrUserNotFound
	}
	if err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}

// UpsertUser 写入用户资料；ctx 为请求上下文，profile 为用户资料。
func (repository *PostgresRepository) UpsertUser(ctx context.Context, profile UserProfile) error {
	_, err := repository.pool.Exec(ctx, `
		INSERT INTO users (
			id,
			avatar_url,
			email,
			is_logged_in,
			is_vip,
			nickname,
			password_hash,
			phone,
			vip_until,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (id) DO UPDATE SET
			avatar_url = EXCLUDED.avatar_url,
			email = EXCLUDED.email,
			is_logged_in = EXCLUDED.is_logged_in,
			is_vip = EXCLUDED.is_vip,
			nickname = EXCLUDED.nickname,
			password_hash = EXCLUDED.password_hash,
			phone = EXCLUDED.phone,
			vip_until = EXCLUDED.vip_until,
			updated_at = NOW()
	`, profile.ID, profile.AvatarURL, profile.Email, profile.IsLoggedIn, profile.IsVip, profile.Nickname, profile.PasswordHash, profile.Phone, profile.VipUntil)

	return err
}

// ListFavorites 读取收藏列表；ctx 为请求上下文，userID 为用户标识。
func (repository *PostgresRepository) ListFavorites(ctx context.Context, userID string) ([]FavoriteItem, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT video_id, title, category, progress, cover_gradient, description, added_at
		FROM user_favorites
		WHERE user_id = $1
		ORDER BY added_at DESC, video_id ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []FavoriteItem{}
	for rows.Next() {
		var item FavoriteItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Category, &item.Progress, &item.CoverGradient, &item.Description, &item.AddedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// UpsertFavorite 写入收藏；ctx 为请求上下文，userID 为用户标识，item 为收藏项。
func (repository *PostgresRepository) UpsertFavorite(ctx context.Context, userID string, item FavoriteItem) error {
	_, err := repository.pool.Exec(ctx, `
		INSERT INTO user_favorites (
			user_id,
			video_id,
			title,
			category,
			progress,
			cover_gradient,
			description,
			added_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (user_id, video_id) DO UPDATE SET
			title = EXCLUDED.title,
			category = EXCLUDED.category,
			progress = EXCLUDED.progress,
			cover_gradient = EXCLUDED.cover_gradient,
			description = EXCLUDED.description,
			added_at = EXCLUDED.added_at,
			updated_at = NOW()
	`, userID, item.ID, item.Title, item.Category, item.Progress, item.CoverGradient, item.Description, item.AddedAt)

	return err
}

// DeleteFavorite 删除收藏；ctx 为请求上下文，userID 和 videoID 定位收藏。
func (repository *PostgresRepository) DeleteFavorite(ctx context.Context, userID string, videoID string) error {
	_, err := repository.pool.Exec(ctx, `
		DELETE FROM user_favorites
		WHERE user_id = $1 AND video_id = $2
	`, userID, videoID)

	return err
}

// ListWatchHistory 读取观看历史；ctx 为请求上下文，userID 为用户标识。
func (repository *PostgresRepository) ListWatchHistory(ctx context.Context, userID string) ([]WatchHistoryItem, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT
			video_id,
			title,
			category,
			progress,
			cover_gradient,
			episode,
			watch_seconds,
			duration_seconds,
			watched_at
		FROM user_watch_history
		WHERE user_id = $1
		ORDER BY watched_at DESC, video_id ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []WatchHistoryItem{}
	for rows.Next() {
		var item WatchHistoryItem
		var episode int
		var watchSeconds pgtype.Int4
		var durationSeconds pgtype.Int4
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Category,
			&item.Progress,
			&item.CoverGradient,
			&episode,
			&watchSeconds,
			&durationSeconds,
			&item.WatchedAt,
		); err != nil {
			return nil, err
		}
		if episode > 0 {
			item.Episode = intPointer(episode)
		}
		if watchSeconds.Valid {
			item.WatchSeconds = intPointer(int(watchSeconds.Int32))
		}
		if durationSeconds.Valid {
			item.DurationSeconds = intPointer(int(durationSeconds.Int32))
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// UpsertWatchHistory 写入观看历史；ctx 为请求上下文，userID 为用户标识，item 为历史项。
func (repository *PostgresRepository) UpsertWatchHistory(ctx context.Context, userID string, item WatchHistoryItem) error {
	_, err := repository.pool.Exec(ctx, `
		INSERT INTO user_watch_history (
			user_id,
			video_id,
			episode,
			title,
			category,
			progress,
			cover_gradient,
			watch_seconds,
			duration_seconds,
			watched_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
		ON CONFLICT (user_id, video_id, episode) DO UPDATE SET
			title = EXCLUDED.title,
			category = EXCLUDED.category,
			progress = EXCLUDED.progress,
			cover_gradient = EXCLUDED.cover_gradient,
			watch_seconds = EXCLUDED.watch_seconds,
			duration_seconds = EXCLUDED.duration_seconds,
			watched_at = EXCLUDED.watched_at,
			updated_at = NOW()
	`, userID, item.ID, episodeValue(item.Episode), item.Title, item.Category, item.Progress, item.CoverGradient, item.WatchSeconds, item.DurationSeconds, item.WatchedAt)

	return err
}

// DeleteWatchHistory 删除观看历史；ctx 为请求上下文，userID、videoID 和 episode 定位历史。
func (repository *PostgresRepository) DeleteWatchHistory(ctx context.Context, userID string, videoID string, episode *int) error {
	if episode == nil {
		_, err := repository.pool.Exec(ctx, `
			DELETE FROM user_watch_history
			WHERE user_id = $1 AND video_id = $2
		`, userID, videoID)
		return err
	}

	_, err := repository.pool.Exec(ctx, `
		DELETE FROM user_watch_history
		WHERE user_id = $1 AND video_id = $2 AND episode = $3
	`, userID, videoID, episodeValue(episode))

	return err
}

// ClearWatchHistory 清空观看历史；ctx 为请求上下文，userID 为用户标识。
func (repository *PostgresRepository) ClearWatchHistory(ctx context.Context, userID string) error {
	_, err := repository.pool.Exec(ctx, `
		DELETE FROM user_watch_history
		WHERE user_id = $1
	`, userID)

	return err
}

// intPointer 返回 int 指针；value 为待包装数值。
func intPointer(value int) *int {
	return &value
}

// episodeValue 返回数据库集数值；episode 为空时返回 0。
func episodeValue(episode *int) int {
	if episode == nil {
		return 0
	}

	return *episode
}
