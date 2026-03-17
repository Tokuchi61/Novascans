# Novascans

## What This Is

Novascans icin backend-first, faz bazli ilerleyen bir urun gelistirme reposu. Temel altyapi Faz 1'de kuruldu; Faz 2 ile `identity/auth` ve `identity/access` omurgasi production'a yakin cekirdek seviyeye tasindi.

## Core Value

Auth, access, manga, topluluk ve admin modullerini yeniden altyapi kurmadan tasiyabilecek temiz, hizli ve genisletilebilir bir Go backend zemini kurmak.

## Current State

- [x] Faz 1: altyapi temeli tamamlandi
- [x] Faz 2: identity core ve access control tamamlandi
- [x] Faz 3: `user/account` tamamlandi
- [ ] Faz 4 planlaniyor: `content/manga`

## Delivered So Far

### Phase 1

- `chi + sqlc + goose` tabanli modul omurgasi
- `cmd/api`, `cmd/migrate`, typed env config, Docker, PostgreSQL
- merkezi error/response standardi
- `healthz`, `readyz`, `metrics`
- `identity/auth` icin ilk CRUD/session iskeleti

### Phase 2

- `identity/auth` modulunun `http / app / domain / store` sinirlarina refactor edilmesi
- auth kimlik alanlarinin `uuid` standardina tasinmasi
- access token + refresh session modeli
- `register`, `login`, `refresh`, `logout`, `logout-all`, `me`
- `email verify request`, `email verify`, `forgot password`, `reset password`
- `identity/access` modulunun eklenmesi
- sabit base role modeli: `guest`, `user`, `moderator`, `admin`
- coklu sub-role modeli ve permission katalogu
- admin tarafindan sub-role olusturma ve kullaniciya atama endpointleri
- tekrar uretilebilir seed komutu

### Phase 3

- `user/account` modulunun `profile / settings / privacy` sinirlariyla eklenmesi
- account tablolari, migrationlari, query dosyalari ve `sqlc` generated kodunun eklenmesi
- register akisinda default account kayitlarinin tek transaction icinde olusturulmasi
- `account/me`, own profile/settings/privacy read-update ve public profile endpointlerinin eklenmesi
- `username` tabanli public profile ve `public / authenticated / private` privacy davranisinin eklenmesi
- seed kullanicilar icin default account kayitlarinin da idempotent olarak uretilmesi

## Active Technical Shape

### Stack

- Go
- Docker
- PostgreSQL
- Router: `chi`
- Query layer: `sqlc`
- Migrations: `goose`

### Module Layout

- `identity/auth`
- `identity/access`
- `user/account`

Fiziksel organizasyon `kategori -> gercek modul` seklini izler. Ust klasorler yalnizca gruplama amaclidir.

### Auth Decisions

- access token imzali ve kisa omurludur
- refresh/session kaydi veritabaninda tutulur
- refresh rotation uygulanir
- verify/reset akislarinda tek kullanimlik token modeli vardir
- mail provider bu asamada bilincli olarak disaridadir
- development ortaminda verify/reset akislarini test etmek icin debug token donulebilir

### Access Decisions

- `guest` veritabaninda kullanici kaydi degil, runtime principal'dir
- kalici kullanicilar icin base role `users.base_role` alaninda tutulur
- base role seti sabittir: `user`, `moderator`, `admin`
- kullaniciya `0..n` adet sub-role atanabilir
- sub-role yetkileri permission katalogundan secilir
- `admin` authorization tarafinda tam yetki bypass olarak davranir

### Account Decisions

- `auth` kullanici kimligini olusturmaya devam eder; `account` yeni kullanici olusturmaz
- `users` tablosu auth sahipliginde kalir ve kimlik koku olarak kullanilir
- `account` fazi `users.id` uzerine bagli kayitlar uretir
- register akisi sirasinda default account kayitlari senkron ve tek transaction icinde olusturulur
- bu entegrasyon event tabanli degil, `AccountProvisioner` benzeri bir port/servis uzerinden saglanir
- ilk account kapsaminda `profile`, `settings` ve `privacy` davranislari ayni modul siniri icinde ele alinir
- public profile icin `username` benzersiz alan olarak account tarafinda sahiplenilir
- `wall`, `friends`, `follow`, `dm`, `library` ve `history` bu fazin disindadir

## Seed Baseline

Seed komutu asagidaki veriyi idempotent olarak yukler:

- permission katalogu:
  - `manga.create`
  - `manga.update`
  - `manga.delete`
  - `comment.moderate`
  - `chapter.create`
  - `chapter.update`
  - `chapter.delete`
- sub-role'ler:
  - `manga_moderator`
  - `comment_moderator`
  - `chapter_moderator`
- ornek kullanicilar:
  - 1 `user`
  - 3 `moderator`
  - 1 `admin`

## Constraints

- frontend ve admin UI bu repoda sonraki fazlara aittir
- `OAuth`, sosyal giris ve `2FA` halen kapsam disidir
- VIP, rollout ve oyunlastirma sistemleri sonraki domain fazlarina aittir

## Next Planning Direction

Siradaki aktif plan yonu `content/manga`. Bu fazda hedef:

- manga, chapter ve page veri modelini kurmak
- public okuma endpointlerini acmak
- management API'lerini mevcut auth/access temeliyle korumak
- reading history, library, comment ve VIP erisim gibi bagli davranislari sonraki fazlara birakmak

---
*Last updated: 2026-03-17 after Phase 3 completion*
