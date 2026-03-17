# Phase 3 Context

**Phase:** 3
**Name:** Account Core
**Captured:** 2026-03-17
**Status:** Ready for planning

## Phase Boundary

Bu fazin amaci `user/account` alanini auth kimliginden ayri ama ona bagli sekilde kurmaktir. Faz 3 yeni kullanici olusturmaz; auth tarafinda olusan kullanicilar icin profil, ayarlar ve gizlilik kayitlarini saglar.

## Locked Decisions

- `auth` kullanici kimligini olusturmaya devam eder
- `account` yeni kullanici olusturmaz
- `users` tablosu auth sahipliginde kalir
- `account` verisi `users.id` uzerine bagli tablolarda tutulur
- register akisi sirasinda default account kayitlari senkron ve tek transaction icinde olusturulur
- bu entegrasyon event tabanli degil, port/servis uzerinden saglanir
- ilk account kapsami `profile`, `settings`, `privacy` ile sinirlidir
- public profile okumasi `username` uzerinden yapilir
- `wall`, `friends`, `follow`, `dm`, `library`, `history` bu fazin disindadir

## Suggested Data Shape

- `account_profiles`
  - `user_id`
  - `username`
  - `display_name`
  - `bio`
  - `avatar_path`
  - `banner_path`
  - `created_at`
  - `updated_at`
- `account_settings`
  - `user_id`
  - `locale`
  - `timezone`
  - `created_at`
  - `updated_at`
- `account_privacy_settings`
  - `user_id`
  - `profile_visibility`
  - `created_at`
  - `updated_at`

## Suggested API Surface

- `GET /api/v1/account/me`
- `GET /api/v1/account/profile`
- `PATCH /api/v1/account/profile`
- `GET /api/v1/account/settings`
- `PATCH /api/v1/account/settings`
- `GET /api/v1/account/privacy`
- `PATCH /api/v1/account/privacy`
- `GET /api/v1/account/profile/{username}`

## Integration Notes

- register akisi `auth` icinde kalir
- auth service, account bootstrap icin bir port kullanir
- account implementasyonu `internal/app` icinde auth'a baglanir
- bootstrap basarisiz olursa kullanici, credential ve session kaydi rollback olur

## Deferred Ideas

- wall
- friends
- follow
- dm
- library
- history
- manga

---
*Last updated: 2026-03-17*
