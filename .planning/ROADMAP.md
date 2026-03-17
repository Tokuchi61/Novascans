# Roadmap: Novascans

## Overview

Roadmap sifirlandi ve daraltildi. Su an yalnizca ilk faz planlaniyor: temel backend altyapisi. Bu faz tamamlandiginda proje, sonraki auth ve domain fazlarini kirilmadan tasiyabilecek net bir teknik zemine sahip olacak.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

- [x] **Phase 1: Infrastructure Foundation** - Proje iskeleti, router, veri erisim katmani, Docker, PostgreSQL ve temel backend konvansiyonlari netlesir.

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
Phases execute in numeric order: 1

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Infrastructure Foundation | 5/5 | Completed | 2026-03-17 (aligned with 0.1.1 patch) |
