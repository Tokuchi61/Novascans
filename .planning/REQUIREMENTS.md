# Requirements: Novascans Infrastructure Foundation

**Defined:** 2026-03-17
**Core Value:** Sonraki auth, manga, topluluk ve admin modullerini tekrar temel degistirmeden tasiyabilecek temiz, hizli ve genisletilebilir bir backend omurgasi kurmak.

## v1 Requirements

### Architecture

- [ ] **ARCH-01**: Backend, sonradan auth ve domain modulleri eklenebilecek moduler bir proje iskeletiyle baslar.
- [ ] **ARCH-04**: Modul dizinleri `kategori -> gercek modul` seklinde duzenlenir; ust klasorler yalnizca gruplama amaci tasir.
- [ ] **ARCH-02**: HTTP server baslatma, middleware zinciri, graceful shutdown ve merkezi bootstrap akisi standartlastirilir.
- [ ] **ARCH-03**: Asenkron isler veya event tabanli genisleme icin temel extension point'ler ayrilir.
- [ ] **ARCH-05**: Event yayini icin `EventBus` interface'i ve en az bir `in-memory` implementasyon saglanir.
- [ ] **ARCH-06**: Modul kayit modeli constructor tabanli dependency injection ve merkezi route registration ile standartlastirilir.
- [ ] **ARCH-07**: Ilk fiziksel modul ornegi `identity/auth` olarak kurulur ve modul deseni gercek bir alan uzerinde dogrulanir.
- [ ] **ARCH-08**: `identity/auth` modulu Faz 1'de tam auth is kurali yerine gercek route, repo, sqlc ve migration omurgasi ile dogrulanir.
- [ ] **ARCH-09**: Fiziksel klasor agaci yalnizca aktif kullanim ihtiyacina gore kurulur; `cmd/api`, `internal/app`, `internal/platform` ve `internal/modules/identity/auth` temel omurgayi olusturur.

### Data Layer

- [ ] **DATA-01**: PostgreSQL baglantisi, migration akisi ve lokal veritabani yasam dongusu Docker icinde standart hale gelir.
- [ ] **DATA-02**: Secilen query katmani sonraki auth, user ve manga modullerine uygun bir veri erisim modeli sunar.
- [ ] **DATA-03**: Repository/service sinirlari, transaction kurallari ve veri erisim konvansiyonlari yazili hale getirilir.
- [ ] **DATA-04**: SQL kaynaklari `db/queries/<kategori>/<modul>` yapisinda saklanir ve `sqlc` generated kodu ilgili modullerde konumlanir.
- [ ] **DATA-05**: Transaction yonetimi service/use-case seviyesinde standartlastirilir ve repository katmani gizli transaction acmaz.
- [ ] **DATA-06**: Migration dosyalari sirali numara + kisa aciklama standardi ile yonetilir.
- [ ] **DATA-07**: `identity/auth` icin ilk veri zemini `users`, `auth_password_credentials` ve `auth_sessions` tablolarini ayri sorumluluklarla kurar.

### API Foundation

- [ ] **API-01**: Secilen router moduler route registration, versioning ve middleware kompozisyonunu destekler.
- [ ] **API-02**: Standart hata cevabi, request validation ve yapilandirilmis loglama tum API katmaninda tutarli uygulanir.
- [ ] **API-03**: Health ve readiness endpoint'leri container ve lokal calisma senaryolari icin hazir olur.
- [ ] **API-04**: Basari ve hata response formatlari merkezi helper veya writer ile standartlastirilir.
- [ ] **API-05**: Global middleware zinciri request-id, recover, timeout, logging ve metrics davranislarini standartlastirir.
- [ ] **API-06**: Public route yapisi `/api/v1/<module>` standardini izler; sistem endpoint'leri kok seviyede kalir.
- [ ] **API-07**: `identity/auth` modulu Faz 1'de `ping`, `user create/read` ve `session create/revoke` endpointleriyle dogrulanir.

### Developer Experience

- [ ] **DEV-01**: Lokal bootstrap, env yonetimi ve temel calistirma komutlari tekrar edilebilir hale getirilir.
- [ ] **DEV-02**: Gelecek fazlar icin minimum test ve kalite kapilari tanimlanir.
- [ ] **DEV-03**: Config yalnizca `env` uzerinden yuklenir, typed parse edilir ve eksik zorunlu alanlarda uygulama fail-fast davranir.
- [ ] **DEV-04**: Docker gelistirme ortami varsayilan olarak `api + postgres` ile gelir ve testler ayri test veritabani kullanir.
- [ ] **DEV-05**: Test stratejisi unit, handler/http, integration ve smoke katmanlarini ayirir; integration testleri ayri calistirma yoluna sahiptir.
- [ ] **DEV-06**: Env isimlendirme standardi `NOVASCANS_` prefix'i ve gruplu alan yapisi ile standartlastirilir.
- [ ] **DEV-07**: Tum tamamlanan uygulama fazlari `CHANGELOG.md` uzerinden semver mantigi ile kaydedilir.
- [ ] **DEV-08**: Faz tamamlama protokolu dokuman inceleme, kod inceleme, uygulama, test, changelog/version guncelleme ve git review adimlarini standartlastirir.
- [ ] **DEV-09**: Git tag veya release isimleri surum, faz numarasi ve kapsam bilgisi ile okunabilir formatta tutulur.

## v2 Requirements

### Domain Modules

- **AUTH-01**: Auth ve authorization modulu bu temel uzerine ayrica planlanir.
- **USER-01**: User ve profile modulu ayrica planlanir.
- **MANGA-01**: Manga ve chapter modulu ayrica planlanir.
- **COMM-01**: Wall, yorum ve DM modulu ayrica planlanir.

## Out of Scope

| Feature | Reason |
|---------|--------|
| Full auth implementation | Bu fazda sadece auth'a uygun temel altyapi kurulacak |
| Admin control plane | Altyapi kararlarindan sonra planlanacak |
| Business rules such as VIP, XP, rollout | Temel veri ve API omurgasi netlesmeden karar verilmemeli |
| Frontend and admin UI | Bu faz yalnizca backend altyapisina odakli |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| ARCH-01 | Phase 1 | Pending |
| ARCH-04 | Phase 1 | Pending |
| ARCH-02 | Phase 1 | Pending |
| ARCH-03 | Phase 1 | Pending |
| ARCH-05 | Phase 1 | Pending |
| ARCH-06 | Phase 1 | Pending |
| ARCH-07 | Phase 1 | Pending |
| ARCH-08 | Phase 1 | Pending |
| ARCH-09 | Phase 1 | Pending |
| DATA-01 | Phase 1 | Pending |
| DATA-02 | Phase 1 | Pending |
| DATA-03 | Phase 1 | Pending |
| DATA-04 | Phase 1 | Pending |
| DATA-05 | Phase 1 | Pending |
| DATA-06 | Phase 1 | Pending |
| DATA-07 | Phase 1 | Pending |
| API-01 | Phase 1 | Pending |
| API-02 | Phase 1 | Pending |
| API-03 | Phase 1 | Pending |
| API-04 | Phase 1 | Pending |
| API-05 | Phase 1 | Pending |
| API-06 | Phase 1 | Pending |
| API-07 | Phase 1 | Pending |
| DEV-01 | Phase 1 | Pending |
| DEV-02 | Phase 1 | Pending |
| DEV-03 | Phase 1 | Pending |
| DEV-04 | Phase 1 | Pending |
| DEV-05 | Phase 1 | Pending |
| DEV-06 | Phase 1 | Pending |
| DEV-07 | Phase 1 | Pending |
| DEV-08 | Phase 1 | Pending |
| DEV-09 | Phase 1 | Pending |

**Coverage:**
- v1 requirements: 32 total
- Mapped to phases: 32
- Unmapped: 0

---
*Requirements defined: 2026-03-17*
*Last updated: 2026-03-17 after infrastructure scoping*
