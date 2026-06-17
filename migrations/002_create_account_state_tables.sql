CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  avatar_url TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL DEFAULT '',
  is_logged_in BOOLEAN NOT NULL DEFAULT FALSE,
  is_vip BOOLEAN NOT NULL DEFAULT FALSE,
  nickname TEXT NOT NULL DEFAULT 'Next Video 用户',
  phone TEXT NOT NULL DEFAULT '',
  vip_until TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_favorites (
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  video_id TEXT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  category TEXT NOT NULL,
  progress TEXT NOT NULL,
  cover_gradient TEXT NOT NULL,
  description TEXT NOT NULL,
  added_at BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, video_id)
);

CREATE TABLE IF NOT EXISTS user_watch_history (
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  video_id TEXT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
  episode INTEGER NOT NULL DEFAULT 0,
  title TEXT NOT NULL,
  category TEXT NOT NULL,
  progress TEXT NOT NULL,
  cover_gradient TEXT NOT NULL,
  watch_seconds INTEGER,
  duration_seconds INTEGER,
  watched_at BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, video_id, episode)
);

CREATE INDEX IF NOT EXISTS idx_user_favorites_added_at ON user_favorites(user_id, added_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_watch_history_watched_at ON user_watch_history(user_id, watched_at DESC);
