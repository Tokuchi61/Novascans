# Requirements: Novascans

**Defined:** 2026-03-17
**Core Value:** Sonraki auth, manga, topluluk ve admin modullerini tekrar temel degistirmeden tasiyabilecek temiz, hizli ve genisletilebilir bir backend omurgasi kurmak.

## v1 Requirements

### Architecture

- [x] **ARCH-01**: Backend, sonradan auth ve domain modulleri eklenebilecek moduler bir proje iskeletiyle baslar.
- [x] **ARCH-04**: Modul dizinleri `kategori -> gercek modul` seklinde duzenlenir; ust klasorler yalnizca gruplama amaci tasir.
- [x] **ARCH-02**: HTTP server baslatma, middleware zinciri, graceful shutdown ve merkezi bootstrap akisi standartlastirilir.
- [x] **ARCH-03**: Asenkron isler veya event tabanli genisleme icin temel extension point'ler ayrilir.
- [x] **ARCH-05**: Event yayini icin `EventBus` interface'i ve en az bir `in-memory` implementasyon saglanir.
- [x] **ARCH-06**: Modul kayit modeli constructor tabanli dependency injection ve merkezi route registration ile standartlastirilir.
- [x] **ARCH-07**: Ilk fiziksel modul ornegi `identity/auth` olarak kurulur ve modul deseni gercek bir alan uzerinde dogrulanir.
- [x] **ARCH-08**: `identity/auth` modulu Faz 1'de tam auth is kurali yerine gercek route, repo, sqlc ve migration omurgasi ile dogrulanir.
- [x] **ARCH-09**: Fiziksel klasor agaci yalnizca aktif kullanim ihtiyacina gore kurulur; `cmd/api`, `internal/app`, `internal/platform` ve `internal/modules/identity/auth` temel omurgayi olusturur.

### Data Layer

- [x] **DATA-01**: PostgreSQL baglantisi, migration akisi ve lokal veritabani yasam dongusu Docker icinde standart hale gelir.
- [x] **DATA-02**: Secilen query katmani sonraki auth, user ve manga modullerine uygun bir veri erisim modeli sunar.
- [x] **DATA-03**: Repository/service sinirlari, transaction kurallari ve veri erisim konvansiyonlari yazili hale getirilir.
- [x] **DATA-04**: SQL kaynaklari `db/queries/<kategori>/<modul>` yapisinda saklanir ve `sqlc` generated kodu ilgili modullerde konumlanir.
- [x] **DATA-05**: Transaction yonetimi service/use-case seviyesinde standartlastirilir ve repository katmani gizli transaction acmaz.
- [x] **DATA-06**: Migration dosyalari sirali numara + kisa aciklama standardi ile yonetilir.
- [x] **DATA-07**: `identity/auth` icin ilk veri zemini `users`, `auth_password_credentials` ve `auth_sessions` tablolarini ayri sorumluluklarla kurar.

### API Foundation

- [x] **API-01**: Secilen router moduler route registration, versioning ve middleware kompozisyonunu destekler.
- [x] **API-02**: Standart hata cevabi, request validation ve yapilandirilmis loglama tum API katmaninda tutarli uygulanir.
- [x] **API-03**: Health ve readiness endpoint'leri container ve lokal calisma senaryolari icin hazir olur.
- [x] **API-04**: Basari ve hata response formatlari merkezi helper veya writer ile standartlastirilir.
- [x] **API-05**: Global middleware zinciri request-id, recover, timeout, logging ve metrics davranislarini standartlastirir.
- [x] **API-06**: Public route yapisi `/api/v1/<module>` standardini izler; sistem endpoint'leri kok seviyede kalir.
- [x] **API-07**: `identity/auth` modulu Faz 1'de `ping`, `user create/read` ve `session create/revoke` endpointleriyle dogrulanir.

### Developer Experience

- [x] **DEV-01**: Lokal bootstrap, env yonetimi ve temel calistirma komutlari tekrar edilebilir hale getirilir.
- [x] **DEV-02**: Gelecek fazlar icin minimum test ve kalite kapilari tanimlanir.
- [x] **DEV-03**: Config yalnizca `env` uzerinden yuklenir, typed parse edilir ve eksik zorunlu alanlarda uygulama fail-fast davranir.
- [x] **DEV-04**: Docker gelistirme ortami varsayilan olarak `api + postgres` ile gelir ve testler ayri test veritabani kullanir.
- [x] **DEV-05**: Test stratejisi unit, handler/http, integration ve smoke katmanlarini ayirir; integration testleri ayri calistirma yoluna sahiptir.
- [x] **DEV-06**: Env isimlendirme standardi `NOVASCANS_` prefix'i ve gruplu alan yapisi ile standartlastirilir.
- [x] **DEV-07**: Tum tamamlanan uygulama fazlari `CHANGELOG.md` uzerinden semver mantigi ile kaydedilir.
- [x] **DEV-08**: Faz tamamlama protokolu dokuman inceleme, kod inceleme, uygulama, test, changelog/version guncelleme ve git review adimlarini standartlastirir.
- [x] **DEV-09**: Git tag veya release isimleri surum, faz numarasi ve kapsam bilgisi ile okunabilir formatta tutulur.

## v2 Requirements

### Domain Modules

- **AUTH-01**: Auth ve authorization modulu bu temel uzerine ayrica planlanir.
- **USER-01**: User ve profile modulu ayrica planlanir.
- **MANGA-01**: Manga ve chapter modulu ayrica planlanir.
- **COMM-01**: Wall, yorum ve DM modulu ayrica planlanir.

### Phase 2 Candidate Requirements

- [x] **AUTH-02**: Auth modulu `register`, `login`, `refresh`, `logout current session`, `logout all sessions` ve `me` akislarini tamamlar.
- [x] **AUTH-03**: Email verification ve password reset backend akislarina token bazli zemin eklenir; mail provider bu fazda zorunlu degildir.
- [x] **AUTH-04**: `OAuth`, sosyal girisler ve `2FA` Phase 2 kapsami disinda tutulur.
- [x] **AUTH-05**: Auth modulunun ic yapisi DTO, app/domain model ve persistence sinirlarini daha net ayiracak sekilde refactor edilir.
- [x] **ACCESS-01**: `identity/access` modulu sabit sistem rolleri `guest`, `user`, `moderator`, `admin` ile kurulur.
- [x] **ACCESS-02**: `guest` runtime principal olarak ele alinir ve kullanici kaydi olmadan authorization kararina girebilir.
- [x] **ACCESS-03**: Kullaniciya birden fazla alt rol atanabilir; alt roller permission katalogundan yetki alir.
- [x] **ACCESS-04**: `moderator` moderasyon alani icin kilit roldur; panel ici gercek yetki alt rollerle belirlenir.
- [x] **ACCESS-05**: Sistem baslangicinda `manga_moderator`, `comment_moderator`, `chapter_moderator` seed alt rolleri yuklenir.
- [x] **ACCESS-06**: Authorization middleware ve principal cozumleme mantigi sonraki moduller tarafindan tekrar kullanilabilir hale gelir.
- [x] **ARCH-10**: Generated `sqlc` kodu ile el yazisi persistence kodu arasindaki sinir daha belirgin hale getirilir.
- [x] **DATA-08**: Auth ve access veri modelindeki kimlik alanlari `uuid` standardina tasinir.
- [x] **DEV-10**: Seed ve entegrasyon test stratejisi auth + access modullerini birlikte dogrular.

## Out of Scope

| Feature | Reason |
|---------|--------|
| OAuth / sosyal giris / 2FA | Faz 2'de bilincli olarak disarida tutuldu |
| Admin control plane | Altyapi kararlarindan sonra planlanacak |
| Business rules such as VIP, XP, rollout | Temel veri ve API omurgasi netlesmeden karar verilmemeli |
| Frontend and admin UI | Bu faz yalnizca backend altyapisina odakli |

### Phase 3 Candidate Requirements

- [x] **ACCOUNT-01**: `user/account` modulu auth tarafinda olusan kullanicilar icin profil, ayarlar ve gizlilik zeminini kurar.
- [x] **ACCOUNT-02**: `users` tablosu auth sahipliginde kalir; account verisi `users.id` uzerine bagli tablolarda tutulur.
- [x] **ACCOUNT-03**: Register akisi tek transaction icinde default `profile`, `settings` ve `privacy` kayitlarini olusturur.
- [x] **ACCOUNT-04**: Account modulu kendi kullanicisi icin `me`, profile, settings ve privacy read/update endpointlerini saglar.
- [x] **ACCOUNT-05**: Public profile okuma davranisi username uzerinden calisir ve username benzersizligi korunur.
- [x] **ACCOUNT-06**: Privacy davranisi en az profile gorunurlugunu `public`, `authenticated`, `private` seviyelerinde destekler.
- [x] **ACCOUNT-07**: `wall`, `friends`, `follow`, `dm`, `library` ve `history` bu faza dahil edilmeden disarida tutulur.
- [x] **DATA-09**: Account tablolari ve sorgulari `sqlc` akisiyla auth/access standartlariyla uyumlu kurulur.
- [x] **API-08**: Account endpointleri ayni `/api/v1/<module>` sozlesmesi icinde tutarli davranir.
- [x] **DEV-11**: Account smoke ve integration testleri register bootstrap, own-account write ve public-profile read akislarini birlikte dogrular.

### Phase 4 Candidate Requirements

- **MANGA-02**: `content/manga` veri modeli manga, chapter ve ordered page iliskilerini kapsar.
- **MANGA-03**: Public API manga listesi, manga detay ve chapter okuma endpointlerini saglar.
- **MANGA-04**: Manga/chapter yonetim endpointleri auth/access permission zinciri ile korunur.
- **MANGA-05**: Icerik yayin durumu ve gorunurluk alani gelecekteki VIP/erken erisim kurallarina acik olacak sekilde modellenir.
- **MANGA-06**: Seed verisi okunabilir bir manga + chapter akisini tekrar uretilebilir hale getirir.
- **MANGA-07**: Reading history, library ve comment gibi bagli davranislar bu faza dahil edilmeden disarida tutulur.
- **DATA-10**: Manga/chapter/page tablolari ve sorgulari `sqlc` akisiyla auth/access standartlariyla uyumlu kurulur.
- **API-09**: Public ve protected manga endpointleri ayni `/api/v1/<module>` sozlesmesi icinde ayrik ama tutarli davranir.
- **DEV-12**: Manga smoke ve integration testleri public read ve permission-korumali write akislarini birlikte dogrular.

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| ARCH-01 | Phase 1 | Complete |
| ARCH-04 | Phase 1 | Complete |
| ARCH-02 | Phase 1 | Complete |
| ARCH-03 | Phase 1 | Complete |
| ARCH-05 | Phase 1 | Complete |
| ARCH-06 | Phase 1 | Complete |
| ARCH-07 | Phase 1 | Complete |
| ARCH-08 | Phase 1 | Complete |
| ARCH-09 | Phase 1 | Complete |
| DATA-01 | Phase 1 | Complete |
| DATA-02 | Phase 1 | Complete |
| DATA-03 | Phase 1 | Complete |
| DATA-04 | Phase 1 | Complete |
| DATA-05 | Phase 1 | Complete |
| DATA-06 | Phase 1 | Complete |
| DATA-07 | Phase 1 | Complete |
| API-01 | Phase 1 | Complete |
| API-02 | Phase 1 | Complete |
| API-03 | Phase 1 | Complete |
| API-04 | Phase 1 | Complete |
| API-05 | Phase 1 | Complete |
| API-06 | Phase 1 | Complete |
| API-07 | Phase 1 | Complete |
| DEV-01 | Phase 1 | Complete |
| DEV-02 | Phase 1 | Complete |
| DEV-03 | Phase 1 | Complete |
| DEV-04 | Phase 1 | Complete |
| DEV-05 | Phase 1 | Complete |
| DEV-06 | Phase 1 | Complete |
| DEV-07 | Phase 1 | Complete |
| DEV-08 | Phase 1 | Complete |
| DEV-09 | Phase 1 | Complete |
| AUTH-02 | Phase 2 | Complete |
| AUTH-03 | Phase 2 | Complete |
| AUTH-04 | Phase 2 | Complete |
| AUTH-05 | Phase 2 | Complete |
| ACCESS-01 | Phase 2 | Complete |
| ACCESS-02 | Phase 2 | Complete |
| ACCESS-03 | Phase 2 | Complete |
| ACCESS-04 | Phase 2 | Complete |
| ACCESS-05 | Phase 2 | Complete |
| ACCESS-06 | Phase 2 | Complete |
| ARCH-10 | Phase 2 | Complete |
| DATA-08 | Phase 2 | Complete |
| DEV-10 | Phase 2 | Complete |
| ACCOUNT-01 | Phase 3 | Complete |
| ACCOUNT-02 | Phase 3 | Complete |
| ACCOUNT-03 | Phase 3 | Complete |
| ACCOUNT-04 | Phase 3 | Complete |
| ACCOUNT-05 | Phase 3 | Complete |
| ACCOUNT-06 | Phase 3 | Complete |
| ACCOUNT-07 | Phase 3 | Complete |
| DATA-09 | Phase 3 | Complete |
| API-08 | Phase 3 | Complete |
| DEV-11 | Phase 3 | Complete |

**Coverage:**
- v1 requirements: 32 total
- Mapped to phases: 32
- Unmapped: 0

## Verification Notes

- `go test ./...` calistirildi.
- `go test -tags=integration ./...` lokal fallback yoldan calistirildi; `.env` autoload ve test DB hazirlama davranisi dogrulandi.
- Docker agi icinden gecici `golang:1.26.1-alpine` container'i ile integration testler tekrar calistirildi.
- `docker compose config` dogrulandi.
- `docker compose up -d --build` ile `api` ve `postgres` container'lari guncel image'larla yeniden ayaga kaldirildi.
- Container ici `/app/migrate status` ve `/app/migrate up` komutlari dogrulandi.
- PostgreSQL container'i icinde `novascans_test` veritabani tekrar cagrilabilir script ile dogrulandi.
- `readyz`, `metrics`, `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh` ve `GET /api/v1/access/me` canli servis uzerinde dogrulandi.
- Phase 2 icin `go test ./...` ve `go test -tags=integration ./...` auth+access degisiklikleriyle yeniden calistirildi.
- Canli Docker dogrulamasi tekrarlandi: `docker compose up -d --build`, container ici `migrate up`, `seed`, admin login, `access/me`, `refresh`, verify/reset akis smoke senaryolari ve coklu sub-role atama akislari dogrulandi.
- Phase 3 icin `sqlc generate`, `go test ./...` ve `go test -tags=integration ./...` account modulu ve auth-account bootstrap degisiklikleriyle calistirildi.
- Canli Docker dogrulamasi tekrarlandi: `docker compose up -d --build`, container ici `migrate up`, `seed`, yeni register kullanicisi ile `account/me`, profile/settings/privacy update ve public profile smoke akislari dogrulandi.
- Seeded admin kullanicisi ile `account/me` smoke akisi tekrar calistirildi; seed komutunun account defaults uretebildigi dogrulandi.

---
*Requirements defined: 2026-03-17*
*Last updated: 2026-03-17 after Phase 3 completion*
