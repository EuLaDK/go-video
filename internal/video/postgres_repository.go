package video

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository 创建 PostgreSQL 视频仓库；pool 为已建立的数据库连接池。
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// ListChannels 从 PostgreSQL 读取频道列表；ctx 为请求上下文。
func (repository *PostgresRepository) ListChannels(ctx context.Context) ([]Channel, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT slug, label, description, keywords, accent
		FROM channels
		ORDER BY display_order ASC, slug ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	channels := []Channel{}
	for rows.Next() {
		var channel Channel
		if err := rows.Scan(
			&channel.Slug,
			&channel.Label,
			&channel.Description,
			&channel.Keywords,
			&channel.Accent,
		); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return channels, nil
}

// ListVideos 从 PostgreSQL 读取完整视频列表；ctx 为请求上下文。
func (repository *PostgresRepository) ListVideos(ctx context.Context) ([]Video, error) {
	videos, err := repository.listVideoRows(ctx)
	if err != nil {
		return nil, err
	}

	episodes, err := repository.listEpisodes(ctx)
	if err != nil {
		return nil, err
	}

	calendars, err := repository.listReleaseCalendar(ctx)
	if err != nil {
		return nil, err
	}

	relatedIDs, err := repository.listRelatedVideoIDs(ctx)
	if err != nil {
		return nil, err
	}

	for index := range videos {
		videoID := videos[index].ID
		videos[index].Episodes = episodes[videoID]
		videos[index].ReleaseCalendar = calendars[videoID]
		videos[index].RelatedVideoIDs = relatedIDs[videoID]
	}

	return videos, nil
}

// ListPlaybackSources 从 PostgreSQL 读取单个视频的播放源；ctx 为请求上下文，videoID 为视频 id。
func (repository *PostgresRepository) ListPlaybackSources(ctx context.Context, videoID string) ([]PlaybackSource, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT quality, label, source_url, mime_type
		FROM video_playback_sources
		WHERE video_id = $1
		ORDER BY display_order ASC, quality ASC
	`, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := []PlaybackSource{}
	for rows.Next() {
		var source PlaybackSource
		if err := rows.Scan(
			&source.Quality,
			&source.Label,
			&source.SourceURL,
			&source.MimeType,
		); err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}

// listVideoRows 读取视频主体字段；ctx 为请求上下文。
func (repository *PostgresRepository) listVideoRows(ctx context.Context) ([]Video, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT
			id,
			title,
			subtitle,
			description,
			score,
			heat,
			update_text,
			category,
			year_text,
			region,
			total_episodes,
			quality,
			badge,
			progress,
			duration,
			source_url,
			cover_gradient,
			tags,
			cast_names
		FROM videos
		ORDER BY display_order ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := []Video{}
	for rows.Next() {
		var item Video
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Subtitle,
			&item.Description,
			&item.Score,
			&item.Heat,
			&item.Update,
			&item.Category,
			&item.Year,
			&item.Region,
			&item.TotalEpisodes,
			&item.Quality,
			&item.Badge,
			&item.Progress,
			&item.Duration,
			&item.SourceURL,
			&item.CoverGradient,
			&item.Tags,
			&item.CastNames,
		); err != nil {
			return nil, err
		}
		videos = append(videos, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}

// listEpisodes 读取视频选集并按 video_id 分组；ctx 为请求上下文。
func (repository *PostgresRepository) listEpisodes(ctx context.Context) (map[string][]Episode, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT video_id, episode, title, duration, status
		FROM video_episodes
		ORDER BY video_id ASC, episode ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	episodes := map[string][]Episode{}
	for rows.Next() {
		var videoID string
		var item Episode
		var status pgtype.Text
		if err := rows.Scan(&videoID, &item.Episode, &item.Title, &item.Duration, &status); err != nil {
			return nil, err
		}
		if status.Valid {
			item.Status = status.String
		}
		episodes[videoID] = append(episodes[videoID], item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return episodes, nil
}

// listReleaseCalendar 读取更新时间线并按 video_id 分组；ctx 为请求上下文。
func (repository *PostgresRepository) listReleaseCalendar(ctx context.Context) (map[string][]ReleaseCalendarItem, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT video_id, time_text, detail, active
		FROM video_release_calendar
		ORDER BY video_id ASC, item_order ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	calendars := map[string][]ReleaseCalendarItem{}
	for rows.Next() {
		var videoID string
		var item ReleaseCalendarItem
		if err := rows.Scan(&videoID, &item.Time, &item.Detail, &item.Active); err != nil {
			return nil, err
		}
		calendars[videoID] = append(calendars[videoID], item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return calendars, nil
}

// listRelatedVideoIDs 读取相关推荐关系并按 video_id 分组；ctx 为请求上下文。
func (repository *PostgresRepository) listRelatedVideoIDs(ctx context.Context) (map[string][]string, error) {
	rows, err := repository.pool.Query(ctx, `
		SELECT video_id, related_video_id
		FROM video_related
		ORDER BY video_id ASC, display_order ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	relatedIDs := map[string][]string{}
	for rows.Next() {
		var videoID string
		var relatedVideoID string
		if err := rows.Scan(&videoID, &relatedVideoID); err != nil {
			return nil, err
		}
		relatedIDs[videoID] = append(relatedIDs[videoID], relatedVideoID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return relatedIDs, nil
}
