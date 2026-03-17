# Roadmap: Novascans

## Overview

Roadmap dar kapsamli fazlarla ilerliyor. Phase 1, Phase 2 ve Phase 3 tamamlandi. Projenin omurgasi artik `infrastructure + identity/auth + identity/access + user/account` seviyesinde stabil. Siradaki mantikli adim `content/manga` tarafini bu omurga ustune kurmak.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

- [x] **Phase 1: Infrastructure Foundation** - Proje iskeleti, router, veri erisim katmani, Docker, PostgreSQL ve temel backend konvansiyonlari netlesir.
- [x] **Phase 2: Identity Core and Access Control** - Auth akislari tamamlandi, access/RBAC modulu eklendi, UUID standardi yerlestirildi ve auth modulu ic yapisi buyumeye uygun sekilde toparlandi.
- [x] **Phase 3: Account Core** - Auth kimligine bagli profil, ayarlar ve gizlilik zemini kurulur; register akisi default account bootstrap ile genisletilir.
- [ ] **Phase 4: Manga Content Core** - Manga, chapter ve sayfa veri modeli; public okuma endpointleri ve auth/access ile korunan icerik yonetim API'leri kurulur.

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
Phases execute in numeric order: 1 -> 2 -> 3 -> 4

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Infrastructure Foundation | 5/5 | Completed | 2026-03-17 (aligned with 0.1.1 patch) |
| 2. Identity Core and Access Control | 6/6 | Completed | 2026-03-17 |
| 3. Account Core | 5/5 | Completed | 2026-03-17 |
| 4. Manga Content Core | 0/5 | Planned | - |

## Phase Details

### Phase 2: Identity Core and Access Control
**Goal**: `identity/auth` modulunu production'a yakin cekirdek seviyeye tasimak, `identity/access` modulunu eklemek ve authorization modelini sabit roller + coklu alt roller mantigi ile kurmak.
**Depends on**: Phase 1
**Requirements**: [AUTH-01, ACCESS-01, ACCESS-02, ACCESS-03, ACCESS-04, ACCESS-05, ACCESS-06, AUTH-02, AUTH-03, AUTH-04, AUTH-05, ARCH-10, DATA-08, DEV-10]
**Success Criteria** (what must be TRUE):
  1. Auth modulu `register`, `login`, `refresh`, `logout current session`, `logout all sessions`, `me`, `email verify request`, `email verify`, `forgot password` ve `reset password` akislarini destekler.
  2. `OAuth`, sosyal giris ve `2FA` bu faza dahil edilmeden disarida tutulur.
  3. Sistemde sabit `guest`, `user`, `moderator`, `admin` rollerinin anlami nettir; `guest` runtime principal olarak calisir.
  4. `moderator` ve alt rol modeli moderasyon/backoffice alanlari icin tekrar kullanilabilir authorization zemini sunar; birden fazla alt rol tek kullaniciya atanabilir.
  5. Seed verisi en az `user`, `moderator`, `admin` sistem rollerini ve `manga_moderator`, `comment_moderator`, `chapter_moderator` alt rollerini yukler.
  6. Auth ve access veri modeli ile ilgili kimlik alanlari `uuid` standardina tasinir; PostgreSQL tarafinda native `uuid` tipi kullanilir.
  7. Auth modulu ic refactor gecirir; DTO, app/domain model ve persistence sinirlari netlesir, generated `sqlc` kodu el yazisi koddan daha temiz ayrilir.
  8. Authorization middleware ve access karar katmani sonraki manga, comment ve admin fazlarina tekrar kullanilabilir bir temel saglar.
**Plans**: 6 plans
**Post-completion maintenance**: `0.2.0` ile auth core, verification/reset, access module, seed komutu, UUID standardizasyonu ve canli smoke dogrulamasi tamamlandi.

Plans:
- [x] 02-01: `identity/auth` modulunu `http`, `app`, `domain`, `store` sinirlarina gore refactor et; HTTP hata bagimliligini service katmanindan ayir ve generated `sqlc` kod konumunu netlestir
- [x] 02-02: Tum auth kimlik alanlarini `uuid` standardina tasi; PostgreSQL `uuid` kolonlari, Go `uuid` tipi ve migration/query/model donusumunu tamamla
- [x] 02-03: Auth core akislarini tamamla: register, login, refresh rotation, me, current/all session logout ve session yonetim davranislarini sabitle
- [x] 02-04: Email verification ve password reset token zeminini kur; dis mail provider olmadan request/consume akislarini backend seviyesinde tamamla
- [x] 02-05: `identity/access` modulunu kur; sabit sistem rollerini, permission katalogunu, coklu alt rol atamasini, principal/cozumleme ve authorization middleware yapisini uygula
- [x] 02-06: Seed verisini, rol/sub-role baslangic kayitlarini, auth-access entegrasyon testlerini, changelog ve teslim dogrulamalarini tamamla

### Phase 3: Account Core
**Goal**: `user/account` modulunu auth kimliginden ayri ama ona bagli sekilde kurmak; profil, ayarlar ve gizlilik zeminini saglamak ve register akisina senkron account bootstrap eklemek.
**Depends on**: Phase 2
**Requirements**: [ACCOUNT-01, ACCOUNT-02, ACCOUNT-03, ACCOUNT-04, ACCOUNT-05, ACCOUNT-06, ACCOUNT-07, DATA-09, API-08, DEV-11]
**Success Criteria** (what must be TRUE):
  1. `auth` kullanici kimligini olusturmaya devam eder; `account` yeni kullanici olusturmaz.
  2. Register akisi tek transaction icinde default `profile`, `settings` ve `privacy` kayitlarini olusturur.
  3. `account` veri modeli `users.id` uzerine bagli tablolarla kurulur; auth kimlik verisi account alanlarina sizmaz.
  4. Kullanici kendi account verisini goruntuleyebilir ve guncelleyebilir.
  5. Public profile okuma davranisi username uzerinden calisir ve privacy kurallariyla uyumludur.
  6. `wall`, `friends`, `follow`, `dm`, `library` ve `history` bu faza dahil edilmez.
**Plans**: 5 plans
**Post-completion maintenance**: `0.3.0` ile account tablolari, register bootstrap, own-account API'leri, public profile ve privacy davranisi tamamlandi; seed ve canli smoke akislari account verisiyle hizalandi. `0.3.1` patch'i ile account modulu ortak access guard uzerine tasindi, transaction standardi ortak `TxManager` kullanimiyla hizalandi ve Faz 3 planning dokuman drift'i kapatildi.

Plans:
- [x] 03-01: `user/account` domain sinirini netlestir; `profile`, `settings`, `privacy` veri modelini ve auth ile entegrasyon kurallarini sabitle
- [x] 03-02: Account schema, migration, query ve module iskeletini kur; register provisioning portunu ve transaction akisini tamamla
- [x] 03-03: Account API'lerini uygula: `me`, profile read/update, settings read/update, privacy read/update
- [x] 03-04: Public profile ve auth-account bootstrap akislarini uygula; username uniqueness ve privacy davranisini dogrula
- [x] 03-05: Seed, integration/smoke test, changelog ve teslim kayitlarini tamamla

### Phase 4: Manga Content Core
**Goal**: `content/manga` alanini auth/access temeli ve account cekirdegi ustune kurmak; public okuma, chapter/page modeli ve permission-korumali icerik yonetim API'lerini saglamak.
**Depends on**: Phase 3
**Requirements**: [MANGA-02, MANGA-03, MANGA-04, MANGA-05, MANGA-06, MANGA-07, ACCESS-06, DATA-10, API-09, DEV-12]
**Success Criteria** (what must be TRUE):
  1. Sistemde manga, chapter ve ordered page veri modeli bulunur.
  2. Public API tarafinda manga listeleme, manga detay ve chapter okuma endpointleri calisir.
  3. Manga/chapter olusturma ve guncelleme akislari access permission kontrolu ile korunur.
  4. Icerik gorunurluk ve yayin durumu alanlari sonraki VIP/erken erisim fazlarina genisleyebilecek sekilde tanimlanir.
  5. Seed ve test verisi en az bir manga, birden fazla chapter ve okunabilir page setiyle public/read smoke akisini dogrular.
**Plans**: 5 plans

Plans:
- [ ] 04-01: `content/manga` domain sinirini netlestir; manga, chapter, page, publication status ve visibility veri modelini sabitle
- [ ] 04-02: Manga schema, migration, query ve module iskeletini kur; `sqlc` ve repository akisini tamamla
- [ ] 04-03: Public read API'lerini uygula: manga list, manga detail, chapter read
- [ ] 04-04: Permission-korumali management API'lerini uygula: manga/chapter create, update, publish-state degisimi
- [ ] 04-05: Seed, integration/smoke test, changelog ve teslim kayitlarini tamamla
