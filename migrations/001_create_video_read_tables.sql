CREATE TABLE IF NOT EXISTS channels (
  slug TEXT PRIMARY KEY,
  label TEXT NOT NULL,
  description TEXT NOT NULL,
  keywords TEXT[] NOT NULL DEFAULT '{}',
  accent TEXT NOT NULL,
  display_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS videos (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  subtitle TEXT NOT NULL,
  description TEXT NOT NULL,
  score TEXT NOT NULL,
  heat TEXT NOT NULL,
  update_text TEXT NOT NULL,
  category TEXT NOT NULL,
  year_text TEXT NOT NULL,
  region TEXT NOT NULL,
  total_episodes INTEGER NOT NULL,
  quality TEXT NOT NULL,
  badge TEXT NOT NULL,
  progress TEXT NOT NULL,
  duration TEXT NOT NULL,
  source_url TEXT NOT NULL,
  cover_gradient TEXT NOT NULL,
  tags TEXT[] NOT NULL DEFAULT '{}',
  cast_names TEXT[] NOT NULL DEFAULT '{}',
  search_text TEXT NOT NULL DEFAULT '',
  display_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS video_episodes (
  video_id TEXT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
  episode INTEGER NOT NULL,
  title TEXT NOT NULL,
  duration TEXT NOT NULL,
  status TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (video_id, episode)
);

CREATE TABLE IF NOT EXISTS video_release_calendar (
  video_id TEXT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
  item_order INTEGER NOT NULL,
  time_text TEXT NOT NULL,
  detail TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (video_id, item_order)
);

CREATE TABLE IF NOT EXISTS video_related (
  video_id TEXT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
  related_video_id TEXT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
  display_order INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (video_id, related_video_id)
);

CREATE INDEX IF NOT EXISTS idx_channels_display_order ON channels(display_order);
CREATE INDEX IF NOT EXISTS idx_videos_display_order ON videos(display_order);
CREATE INDEX IF NOT EXISTS idx_videos_year_text ON videos(year_text);
CREATE INDEX IF NOT EXISTS idx_videos_quality ON videos(quality);
CREATE INDEX IF NOT EXISTS idx_video_related_order ON video_related(video_id, display_order);
