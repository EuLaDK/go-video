INSERT INTO video_comments (id, video_id, user_id, author, content, created_at_ms)
VALUES
  ('seed-comment-xinghe-1', 'xinghe', 'demo-user', '我', '第一集的悬念铺得很稳，适合继续追。', 1766000001000),
  ('seed-comment-xinghe-2', 'xinghe', 'demo-user', '我', '太空站广播那段氛围感很好。', 1766000002000)
ON CONFLICT (id) DO UPDATE SET
  author = EXCLUDED.author,
  content = EXCLUDED.content,
  created_at_ms = EXCLUDED.created_at_ms;

INSERT INTO video_comment_likes (comment_id, user_id)
VALUES
  ('seed-comment-xinghe-1', 'demo-user')
ON CONFLICT (comment_id, user_id) DO NOTHING;

INSERT INTO video_danmaku (id, video_id, user_id, content, color, created_at_ms)
VALUES
  ('seed-danmaku-xinghe-1', 'xinghe', 'demo-user', '前方信号出现了', 'green', 1766000003000),
  ('seed-danmaku-xinghe-2', 'xinghe', 'demo-user', '这个镜头好漂亮', 'yellow', 1766000004000)
ON CONFLICT (id) DO UPDATE SET
  content = EXCLUDED.content,
  color = EXCLUDED.color,
  created_at_ms = EXCLUDED.created_at_ms;
