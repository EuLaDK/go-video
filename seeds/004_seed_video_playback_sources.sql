BEGIN;

INSERT INTO video_playback_sources (
  video_id,
  quality,
  label,
  source_url,
  mime_type,
  display_order,
  updated_at
)
SELECT
  id,
  quality,
  quality,
  source_url,
  'video/mp4',
  1,
  NOW()
FROM videos
ON CONFLICT (video_id, quality) DO UPDATE SET
  label = EXCLUDED.label,
  source_url = EXCLUDED.source_url,
  mime_type = EXCLUDED.mime_type,
  display_order = EXCLUDED.display_order,
  updated_at = NOW();

INSERT INTO video_playback_sources (
  video_id,
  quality,
  label,
  source_url,
  mime_type,
  display_order,
  updated_at
) VALUES
  ('xinghe', '1080P', '1080P 高清', '/assets/video/staticTest.mp4', 'video/mp4', 2, NOW()),
  ('xinghe', '720P', '720P 流畅', '/assets/video/staticTest.mp4', 'video/mp4', 3, NOW()),
  ('lingdian', '720P', '720P 流畅', '/assets/video/staticTest.mp4', 'video/mp4', 2, NOW())
ON CONFLICT (video_id, quality) DO UPDATE SET
  label = EXCLUDED.label,
  source_url = EXCLUDED.source_url,
  mime_type = EXCLUDED.mime_type,
  display_order = EXCLUDED.display_order,
  updated_at = NOW();

COMMIT;
