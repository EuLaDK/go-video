BEGIN;

CREATE TEMP TABLE seed_channels (
  display_order INTEGER,
  slug TEXT,
  label TEXT,
  description TEXT,
  keywords TEXT[],
  accent TEXT
) ON COMMIT DROP;

INSERT INTO seed_channels (display_order, slug, label, description, keywords, accent) VALUES
  (1, 'featured', '精选', '聚合全站热播内容，覆盖剧集、电影、综艺、纪录片与少儿科普。', ARRAY[]::TEXT[], 'from-emerald-400/28 via-sky-400/18 to-rose-400/18'),
  (2, 'tv', '电视剧', '追新剧、看完结，高热度剧集和口碑长剧都集中在这里。', ARRAY['剧集','都市','悬疑','青春'], 'from-emerald-400/26 via-teal-400/16 to-slate-900/20'),
  (3, 'movie', '电影', '本周新片、冒险大片和高分电影，适合快速进入沉浸观影。', ARRAY['电影'], 'from-rose-500/24 via-indigo-500/18 to-slate-900/20'),
  (4, 'variety', '综艺', '真人秀、挑战企划和轻松陪伴内容，适合碎片时间连续观看。', ARRAY['综艺','真人秀','挑战'], 'from-amber-400/26 via-orange-500/16 to-slate-900/20'),
  (5, 'anime', '动漫', '动画、科幻和亲子向内容，先用现有片库撑起频道形态。', ARRAY['动画','少儿','科幻'], 'from-sky-400/24 via-emerald-400/16 to-slate-900/20'),
  (6, 'documentary', '纪录片', '自然、人文与探索类内容，突出 4K、慢节奏和真实质感。', ARRAY['纪录片','自然','人文'], 'from-cyan-400/22 via-blue-500/16 to-slate-900/20'),
  (7, 'kids', '少儿', '适合家庭和亲子场景的科普、动画与轻松成长内容。', ARRAY['少儿','亲子','科普','动画'], 'from-blue-400/24 via-teal-400/16 to-slate-900/20'),
  (8, 'vip', 'VIP', '会员抢先看、高清片源和独播内容，突出付费权益入口。', ARRAY['会员','独播','会员抢先看'], 'from-emerald-300/26 via-yellow-300/14 to-slate-900/20'),
  (9, 'sports', '体育', '先用竞技、挑战类内容占位，后续可接赛事直播和赛程数据。', ARRAY['竞技','挑战','热血'], 'from-lime-400/22 via-emerald-400/16 to-slate-900/20'),
  (10, 'game', '游戏', '先展示竞技、科幻和热血内容，后续可扩展游戏赛事与直播。', ARRAY['竞技','科幻','动作'], 'from-violet-400/22 via-blue-500/16 to-slate-900/20');

INSERT INTO channels (slug, label, description, keywords, accent, display_order, updated_at)
SELECT slug, label, description, keywords, accent, display_order, NOW()
FROM seed_channels
ON CONFLICT (slug) DO UPDATE SET
  label = EXCLUDED.label,
  description = EXCLUDED.description,
  keywords = EXCLUDED.keywords,
  accent = EXCLUDED.accent,
  display_order = EXCLUDED.display_order,
  updated_at = NOW();

CREATE TEMP TABLE seed_videos (
  display_order INTEGER,
  id TEXT,
  title TEXT,
  subtitle TEXT,
  description TEXT,
  score TEXT,
  heat TEXT,
  update_text TEXT,
  category TEXT,
  year_text TEXT,
  region TEXT,
  total_episodes INTEGER,
  quality TEXT,
  badge TEXT,
  progress TEXT,
  duration TEXT,
  source_url TEXT,
  cover_gradient TEXT,
  tags TEXT[],
  cast_names TEXT[],
  episode_count INTEGER
) ON COMMIT DROP;

INSERT INTO seed_videos (
  display_order,
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
  cast_names,
  episode_count
) VALUES
  (1, 'xinghe', '星河回响', '全网热播 · 科幻悬疑 · 会员抢先看', '近未来的深空勘探计划突然收到来自失联星舰的回声信号，一支临时调查小队被迫进入未知航线，在记忆、谎言与时间错位之间寻找真相。', '9.3', '热度 10026', '更新至 18 集 / 每周五 20:00 更新', '科幻 / 悬疑', '2026', '中国大陆', 24, '4K HDR', '独播', '会员抢先看', '45:00', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#0f766e,#111827 50%,#be123c),radial-gradient(circle_at_72%_28%,rgba(125,211,252,0.55),transparent_25%)', ARRAY['科幻','悬疑','冒险','会员抢先看'], ARRAY['林舟','许念','周砚','陈白'], 12),
  (2, 'anye', '暗夜追光', '犯罪悬疑 · 迷雾追凶 · 高能反转', '一桩旧案在城市停电夜重新浮出水面，刑警与记者沿着光源消失的方向追查，逐步逼近被刻意掩埋的真相。', '8.9', '热度 9821', '全 16 集已完结', '犯罪 / 悬疑', '2026', '中国大陆', 16, '1080P', '高分完结', '全季可看', '42:18', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#1d4ed8,#111827 58%,#0f172a)', ARRAY['犯罪','悬疑','反转','完结'], ARRAY['秦越','唐棠','陆衡','宋知'], 16),
  (3, 'yunduan', '云端餐厅', '都市治愈 · 深夜食堂 · 温暖群像', '一间开在高楼天台的餐厅，只在雨夜营业。不同来客把没说出口的遗憾交给一道菜，也在烟火气里重新出发。', '8.5', '热度 8742', '更新至 10 集 / 每周三更新', '都市 / 治愈', '2026', '中国大陆', 20, '4K', '温暖新剧', '更新至 10 集', '39:46', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#c2410c,#27272a 55%,#111827)', ARRAY['都市','治愈','美食','群像'], ARRAY['沈安','叶晴','乔木','罗一'], 10),
  (4, 'shaonian', '少年棋局', '青春竞技 · 热血成长 · 棋逢对手', '少年棋手从街边棋摊走向全国赛场，在一次次胜负之间学会面对天赋、压力与友情。', '8.3', '热度 8016', '更新至 8 集 / 周末连更', '青春 / 竞技', '2026', '中国大陆', 18, '1080P', '热血成长', '更新至 8 集', '41:20', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#7c3aed,#1f2937 58%,#0f172a)', ARRAY['青春','竞技','成长','热血'], ARRAY['祁远','江夏','孟宁','白辰'], 8),
  (5, 'hai-an', '海岸来信', '纪录片 · 海岛人文 · 慢节奏治愈', '镜头沿着漫长海岸线记录渔村、灯塔与迁徙的人们，在潮汐更替中寻找普通生活里的辽阔。', '8.8', '热度 7388', '全 6 集已上线', '纪录片', '2026', '中国大陆', 6, '4K HDR', '口碑纪录', '4K HDR', '48:08', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#0f766e,#172554 56%,#111827)', ARRAY['纪录片','自然','人文','4K'], ARRAY['旁白：周闻'], 6),
  (6, 'guitu', '归途列车', '电影 · 冒险 · 归途重逢', '暴雪封路后，一列夜行列车被迫停在无人山谷，陌生乘客在共同求生中拼出一段迟到多年的回家路。', '8.1', '热度 6920', '本周新片', '电影 / 冒险', '2026', '中国大陆', 1, '4K', '本周新片', '本周新片', '01:48:20', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#be123c,#312e81 55%,#111827)', ARRAY['电影','冒险','公路','亲情'], ARRAY['梁舟','何曼','苏临'], 1),
  (7, 'chunri', '春日事务所', '剧集 · 都市 · 轻喜治愈', '三位年轻人在旧街区开了一间万能事务所，从修理小物件开始，也慢慢修补人与人之间的关系。', '8.2', '热度 6638', '更新至 12 集', '剧集 / 都市', '2026', '中国大陆', 24, '1080P', '轻喜热播', '更新至 12 集', '40:10', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#15803d,#0f172a 58%,#1f2937)', ARRAY['都市','轻喜','治愈','友情'], ARRAY['唐青','顾南','齐愿','米兰'], 12),
  (8, 'jixian', '极限搭档', '综艺 · 真人秀 · 高能挑战', '六位嘉宾组成临时搭档，在城市与荒野之间完成连续挑战，默契、体力和临场判断都被推到极限。', '8.0', '热度 6501', '第 6 期上线', '综艺 / 真人秀', '2026', '中国大陆', 12, '1080P', '高能综艺', '第 6 期上线', '01:12:06', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#ca8a04,#1e293b 55%,#111827)', ARRAY['综艺','真人秀','挑战','搞笑'], ARRAY['常驻嘉宾团'], 6),
  (9, 'xingqiu', '星球课堂', '少儿 · 科普 · 趣味探索', '用轻松动画和真实实验讲解宇宙、海洋与日常科学，让孩子在故事里理解世界如何运转。', '8.6', '热度 6122', '适合 7+', '少儿 / 科普', '2026', '中国大陆', 30, '1080P', '适合 7+', '适合 7+', '24:30', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#2563eb,#0f766e 56%,#0f172a)', ARRAY['少儿','科普','动画','亲子'], ARRAY['小宇宙讲解团'], 10),
  (10, 'xuexian', '雪线之上', '纪录片 · 自然 · 雪山生态', '摄制组穿越高海拔雪线，记录极端气候下的动物迁徙、冰川变化与守护者的日常。', '8.9', '热度 5988', '4K HDR', '纪录片 / 自然', '2026', '中国大陆', 5, '4K HDR', '自然纪录', '4K HDR', '50:18', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#0369a1,#334155 56%,#020617)', ARRAY['纪录片','自然','雪山','4K'], ARRAY['旁白：林默'], 5),
  (11, 'shen-kong', '深空来客', '科幻 · 悬疑 · 未知信号', '一颗来自太阳系边缘的探测器突然回传陌生影像，科学团队在解析数据时发现它并不孤单。', '8.4', '热度 5772', '同类型热播', '科幻 / 悬疑', '2026', '中国大陆', 12, '4K', '同类型热播', '同类型热播', '42:18', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#1d4ed8,#111827 58%,#0f172a)', ARRAY['科幻','悬疑','太空','探索'], ARRAY['郑原','叶知','罗森'], 6),
  (12, 'lingdian', '零点航线', '冒险 · 灾难 · 极限救援', '跨海航班在零点穿越风暴区，机组和乘客必须在有限时间内完成一场不可能的迫降。', '8.0', '热度 5569', '会员抢先看', '冒险 / 灾难', '2026', '中国大陆', 1, '1080P', '会员抢先看', '会员抢先看', '39:46', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#be123c,#312e81 55%,#111827)', ARRAY['冒险','灾难','救援','电影'], ARRAY['江川','宁栀','贺言'], 1),
  (13, 'jiyi', '记忆穹顶', '悬疑 · 剧情 · 记忆迷局', '一座能储存记忆的城市突然出现集体失忆事件，修复师必须进入别人的过去寻找缺失的一小时。', '8.7', '热度 5420', '高分剧集', '悬疑 / 剧情', '2026', '中国大陆', 18, '4K', '高分剧集', '高分剧集', '45:02', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#0f766e,#172554 56%,#111827)', ARRAY['悬疑','剧情','记忆','反转'], ARRAY['程望','许岚','顾醒'], 8),
  (14, 'bianjing', '边境星门', '科幻 · 动作 · 星际防线', '边境星门意外开启，守备队在未知文明与人类命令之间做出选择，战斗由此改变两个世界的命运。', '8.2', '热度 5304', '正在热播', '科幻 / 动作', '2026', '中国大陆', 20, '4K HDR', '正在热播', '正在热播', '47:31', '/assets/video/staticTest.mp4', 'linear-gradient(135deg,#ca8a04,#1e293b 55%,#111827)', ARRAY['科幻','动作','星际','热血'], ARRAY['陆行','白鹿','韩野'], 6);

INSERT INTO videos (
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
  cast_names,
  search_text,
  display_order,
  updated_at
)
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
  cast_names,
  CONCAT_WS(' ', title, subtitle, description, category, badge, progress, year_text, region, array_to_string(tags, ' '), array_to_string(cast_names, ' ')),
  display_order,
  NOW()
FROM seed_videos
ON CONFLICT (id) DO UPDATE SET
  title = EXCLUDED.title,
  subtitle = EXCLUDED.subtitle,
  description = EXCLUDED.description,
  score = EXCLUDED.score,
  heat = EXCLUDED.heat,
  update_text = EXCLUDED.update_text,
  category = EXCLUDED.category,
  year_text = EXCLUDED.year_text,
  region = EXCLUDED.region,
  total_episodes = EXCLUDED.total_episodes,
  quality = EXCLUDED.quality,
  badge = EXCLUDED.badge,
  progress = EXCLUDED.progress,
  duration = EXCLUDED.duration,
  source_url = EXCLUDED.source_url,
  cover_gradient = EXCLUDED.cover_gradient,
  tags = EXCLUDED.tags,
  cast_names = EXCLUDED.cast_names,
  search_text = EXCLUDED.search_text,
  display_order = EXCLUDED.display_order,
  updated_at = NOW();

DELETE FROM video_related WHERE video_id IN (SELECT id FROM seed_videos);
DELETE FROM video_release_calendar WHERE video_id IN (SELECT id FROM seed_videos);
DELETE FROM video_episodes WHERE video_id IN (SELECT id FROM seed_videos);

INSERT INTO video_episodes (video_id, episode, title, duration, status, updated_at)
SELECT
  seed_videos.id,
  series.episode,
  '第 ' || series.episode || ' 集',
  CASE WHEN series.episode = 1 THEN '正在播放' ELSE '45 分钟' END,
  CASE WHEN series.episode = 1 THEN 'active' ELSE NULL END,
  NOW()
FROM seed_videos
CROSS JOIN LATERAL generate_series(1, seed_videos.episode_count) AS series(episode)
ON CONFLICT (video_id, episode) DO UPDATE SET
  title = EXCLUDED.title,
  duration = EXCLUDED.duration,
  status = EXCLUDED.status,
  updated_at = NOW();

INSERT INTO video_release_calendar (video_id, item_order, time_text, detail, active, updated_at)
SELECT id, 1, '当前', update_text, TRUE, NOW()
FROM seed_videos
ON CONFLICT (video_id, item_order) DO UPDATE SET
  time_text = EXCLUDED.time_text,
  detail = EXCLUDED.detail,
  active = EXCLUDED.active,
  updated_at = NOW();

CREATE TEMP TABLE seed_related (
  video_id TEXT,
  related_video_id TEXT,
  display_order INTEGER
) ON COMMIT DROP;

INSERT INTO seed_related (video_id, related_video_id, display_order) VALUES
  ('xinghe','shen-kong',1),('xinghe','lingdian',2),('xinghe','jiyi',3),('xinghe','bianjing',4),
  ('anye','jiyi',1),('anye','xinghe',2),('anye','lingdian',3),('anye','hai-an',4),
  ('yunduan','chunri',1),('yunduan','guitu',2),('yunduan','xingqiu',3),('yunduan','hai-an',4),
  ('shaonian','jixian',1),('shaonian','xingqiu',2),('shaonian','chunri',3),('shaonian','yunduan',4),
  ('hai-an','xuexian',1),('hai-an','guitu',2),('hai-an','xingqiu',3),('hai-an','shen-kong',4),
  ('guitu','lingdian',1),('guitu','hai-an',2),('guitu','yunduan',3),('guitu','jixian',4),
  ('chunri','yunduan',1),('chunri','shaonian',2),('chunri','xingqiu',3),('chunri','guitu',4),
  ('jixian','shaonian',1),('jixian','guitu',2),('jixian','chunri',3),('jixian','xingqiu',4),
  ('xingqiu','shen-kong',1),('xingqiu','hai-an',2),('xingqiu','shaonian',3),('xingqiu','chunri',4),
  ('xuexian','hai-an',1),('xuexian','xingqiu',2),('xuexian','shen-kong',3),('xuexian','guitu',4),
  ('shen-kong','xinghe',1),('shen-kong','lingdian',2),('shen-kong','jiyi',3),('shen-kong','xingqiu',4),
  ('lingdian','guitu',1),('lingdian','xinghe',2),('lingdian','shen-kong',3),('lingdian','jixian',4),
  ('jiyi','anye',1),('jiyi','xinghe',2),('jiyi','shen-kong',3),('jiyi','lingdian',4),
  ('bianjing','xinghe',1),('bianjing','shen-kong',2),('bianjing','lingdian',3),('bianjing','shaonian',4);

INSERT INTO video_related (video_id, related_video_id, display_order)
SELECT video_id, related_video_id, display_order
FROM seed_related
ON CONFLICT (video_id, related_video_id) DO UPDATE SET
  display_order = EXCLUDED.display_order;

COMMIT;
