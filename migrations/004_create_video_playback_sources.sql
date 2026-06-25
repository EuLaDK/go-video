CREATE TABLE IF NOT EXISTS video_playback_sources (
  video_id TEXT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
  quality TEXT NOT NULL,
  label TEXT NOT NULL,
  source_url TEXT NOT NULL,
  mime_type TEXT NOT NULL,
  display_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (video_id, quality)
);

CREATE INDEX IF NOT EXISTS idx_video_playback_sources_order
ON video_playback_sources(video_id, display_order);
