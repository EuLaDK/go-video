INSERT INTO users (
  id,
  avatar_url,
  email,
  is_logged_in,
  is_vip,
  nickname,
  phone,
  vip_until,
  updated_at
) VALUES (
  'demo-user',
  '',
  '',
  FALSE,
  FALSE,
  'Next Video 用户',
  '',
  '',
  NOW()
),
(
  'demo-vip',
  '',
  '',
  TRUE,
  TRUE,
  'Next Video VIP 用户',
  '',
  '2099-12-31',
  NOW()
)
ON CONFLICT (id) DO NOTHING;
