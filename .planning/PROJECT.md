# Novascans

## What This Is

Novascans icin backend-first, faz bazli ilerleyen bir urun gelistirme reposu. Temel altyapi Faz 1'de kuruldu; Faz 2 ile `identity/auth` ve `identity/access` omurgasi production'a yakin cekirdek seviyeye tasindi.

## Core Value

Auth, access, manga, topluluk ve admin modullerini yeniden altyapi kurmadan tasiyabilecek temiz, hizli ve genisletilebilir bir Go backend zemini kurmak.

## Current State

- [x] Faz 1: altyapi temeli tamamlandi
- [x] Faz 2: identity core ve access control tamamlandi
- [ ] Siradaki faz: domain odakli modullerden biri (`manga`, `moderation`, `admin` veya `user/account`) secilecek

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

Bir sonraki faz secilirken auth/access altyapisini tekrar kullanacak bir domain secilmeli. En mantikli adaylar:

- `content/manga`
- `backoffice/moderation`
- `backoffice/admin`
- `user/account`

---
*Last updated: 2026-03-17 after Phase 2 completion*
