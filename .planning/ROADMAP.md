# Roadmap: Novascans

## Overview

Roadmap dar kapsamli fazlarla ilerliyor. Phase 1 ve Phase 2 tamamlandi. Projenin omurgasi artik `infrastructure + identity/auth + identity/access` seviyesinde stabil.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

- [x] **Phase 1: Infrastructure Foundation** - Proje iskeleti, router, veri erisim katmani, Docker, PostgreSQL ve temel backend konvansiyonlari netlesir.
- [x] **Phase 2: Identity Core and Access Control** - Auth akislari tamamlandi, access/RBAC modulu eklendi, UUID standardi yerlestirildi ve auth modulu ic yapisi buyumeye uygun sekilde toparlandi.

## Phase Details

### Phase 1: Infrastructure Foundation
**Goal**: Sonraki moduller icin tekrar kullanilabilir bir backend zemini kurmak ve `chi + sqlc + goose` tabanli cekirdek teknik secimleri netlestirmek.
**Depends on**: Nothing (first phase)
**Requirements**: [ARCH-01, ARCH-04, ARCH-02, ARCH-03, ARCH-05, ARCH-06, ARCH-07, ARCH-08, ARCH-09, DATA-01, DATA-02, DATA-03, DATA-04, DATA-05, DATA-06, DATA-07, API-01, API-02, API-03, API-04, API-05, API-06, API-07, DEV-01, DEV-02, DEV-03, DEV-04, DEV-05, DEV-06, DEV-07, DEV-08, DEV-09]
**Success Criteria** (what must be TRUE):
  1. Gelistirici servis ve PostgreSQL'i Docker ile tekrar edilebilir sekilde ayaga kaldirabilir.
  2. Proje iskeleti yeni moduller eklenirken route, config ve veri erisim katmanini tekrar yazmayi gerektirmez.
  3. Router, middleware, migration, veri erisim ve event stratejisi belgelenmis ve uygulanabilir hale gelir.
  4. Hata modeli, validation, loglama, `healthz`, `readyz`, `metrics`, env tabanli config ve temel test/gelistirme akisi standartlastirilir.
  5. `identity/auth` modulu gercek bir ornek olarak ayaga kalkar, `users`, `auth_password_credentials` ve `auth_sessions` zeminiyle, `ping`, `user create/read` ve `session create/revoke` endpointleriyle dogrulanir, ancak tum auth feature setini bu fazda zorunlu kilmaz.
  6. Faz tamamlama ve versiyonlama disiplini `CHANGELOG.md` ve standart teslim protokolu ile kurulmus olur.
**Plans**: 5 plans
**Post-completion maintenance**: `0.1.1` patch guncellemesi ile migration akisi container ici komuta tasindi, test veritabani otomasyonu netlestirildi, `.dockerignore` eklendi ve validation karari kodla hizalandi.

Plans:
- [x] 01-01: Proje klasor yapisi, kategori -> modul hiyerarsisi, `cmd/api` + `internal/app` bootstrap akisi, fiziksel agac ve `identity/auth` ornek modulunun kapsam taniminin yapilmasi
- [x] 01-02: Docker, PostgreSQL 18.3, Go 1.26.1, `NOVASCANS_` env semasi, config yukleme ve migration stratejisinin kurulmasi
- [x] 01-03: Chi router, middleware zinciri, route/versioning standardi, auth modulu ilk endpoint seti, validation, hata modeli ve response standardinin netlestirilmesi
- [x] 01-04: Migration naming, auth tablo zemini, sqlc veri erisim akisi, repository kurallari ve service-seviyesi transaction sinirlari
- [x] 01-05: Event altyapisi, loglama, health/readiness/metrics, test tabani, changelog/semver disiplini, git naming kurallari ve gelistirici deneyimi komutlari

## Progress

**Execution Order:**
Phases execute in numeric order: 1 -> 2

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Infrastructure Foundation | 5/5 | Completed | 2026-03-17 (aligned with 0.1.1 patch) |
| 2. Identity Core and Access Control | 6/6 | Completed | 2026-03-17 |

## Phase Details

### Phase 2: Identity Core and Access Control
**Goal**: `identity/auth` modulunu production'a yakin cekirdek seviyeye tasimak, `identity/access` modulunu eklemek ve authorization modelini sabit roller + coklu alt roller mantigi ile kurmak.
**Depends on**: Phase 1
**Requirements**: [AUTH-01, ACCESS-01, ACCESS-02, ACCESS-03, ACCESS-04, ACCESS-05, ACCESS-06, AUTH-02, AUTH-03, AUTH-04, AUTH-05, ARCH-10, DATA-08, DEV-10]
**Success Criteria** (what must be TRUE):
  1. Auth modulu `register`, `login`, `refresh`, `logout current session`, `logout all sessions`, `me`, `email verify request`, `email verify`, `forgot password` ve `reset password` akislarini destekler.
  2. `OAuth`, sosyal giris ve `2FA` bu faza dahil edilmeden disarida tutulur.
  3. Sistemde sabit `guest`, `user`, `moderator`, `admin` rollerinin anlami nettir; `guest` runtime principal olarak calisir.
  4. `moderator` rolune sahip kullanici moderasyon alanina girebilir, ama gorecegi/isletecegi kisimlar admin tarafindan atanabilen birden fazla alt rol ile belirlenir.
  5. Seed verisi en az `user`, `moderator`, `admin` sistem rollerini ve `manga_moderator`, `comment_moderator`, `chapter_moderator` alt rollerini yukler.
  6. Auth ve access veri modeli ile ilgili kimlik alanlari `uuid` standardina tasinir; PostgreSQL tarafinda native `uuid` tipi kullanilir.
  7. Auth modulu ic refactor gecirir; DTO, app/domain model ve persistence sinirlari netlesir, generated `sqlc` kodu el yazisi koddan daha temiz ayrilir.
  8. Authorization middleware ve access karar katmani sonraki manga, comment ve admin fazlarina tekrar kullanilabilir bir temel saglar.
**Plans**: 6 plans
**Post-completion maintenance**: `0.2.0` ile auth core, verification/reset, access module, seed komutu, UUID standardizasyonu ve canlı smoke dogrulamasi tamamlandi.

Plans:
- [x] 02-01: `identity/auth` modulunu `http`, `app`, `domain`, `store` sinirlarina gore refactor et; HTTP hata bagimliligini service katmanindan ayir ve generated `sqlc` kod konumunu netlestir
- [x] 02-02: Tum auth kimlik alanlarini `uuid` standardina tasi; PostgreSQL `uuid` kolonlari, Go `uuid` tipi ve migration/query/model donusumunu tamamla
- [x] 02-03: Auth core akislarini tamamla: register, login, refresh rotation, me, current/all session logout ve session yonetim davranislarini sabitle
- [x] 02-04: Email verification ve password reset token zeminini kur; dis mail provider olmadan request/consume akislarini backend seviyesinde tamamla
- [x] 02-05: `identity/access` modulunu kur; sabit sistem rollerini, permission katalogunu, coklu alt rol atamasini, principal/cozumleme ve authorization middleware yapisini uygula
- [x] 02-06: Seed verisini, rol/sub-role baslangic kayitlarini, auth-access entegrasyon testlerini, changelog ve teslim dogrulamalarini tamamla
