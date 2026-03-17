# Changelog

Bu dosya proje degisikliklerini kaydeder.

## Versioning Policy

- Proje `semver` kullanir.
- Erken gelistirme doneminde surumler `0.x.y` cizgisinde ilerler.
- Geriye donuk uyumlu yeni yetenekler `minor` artisi ile kaydedilir.
- Davranis kirilimi yaratmayan duzeltmeler `patch` artisi ile kaydedilir.
- Proje kararlilik kazanana kadar buyuk kirilimlar `0.x` icinde dikkatli sekilde yonetilir.

## Entry Rules

- Her tamamlanan uygulama fazi bu dosyaya islenir.
- Kayitlar degisiklik kapsamlarini, teknik etkilerini ve gerekiyorsa migration/test notlarini icermelidir.
- Planning tartismalari tek basina changelog girdisi olusturmaz; uygulanmis sonuc kaydedilir.

## Naming Convention

- Repo reference: `https://github.com/Tokuchi61/Novascans`
- Ana branch `main` olarak sabit kalir; surumleme branch adi uzerinden degil tag ve release uzerinden izlenir.
- Faz branch'leri okunabilir format kullanir:
  - `phase/01-infrastructure-foundation`
  - `phase/02-auth-core`
- Gerektiginde feature ve docs branch'leri de ayni mantigi izler:
  - `feature/identity-auth-sessions`
  - `docs/phase-01-foundation-rules`
- Kod release etiketleri su formati izler:
  - `release/v0.1.0-phase-01-infrastructure-foundation`
  - `release/v0.2.0-phase-02-auth-core`
- Dokuman snapshot etiketleri su formati izler:
  - `docs/v0.1.0-phase-01-infrastructure-foundation`
  - `docs/v0.2.0-phase-02-auth-core`
- Plan veya milestone snapshot'lari gerekiyorsa ayri etiketlenir:
  - `plan/v0.1.0-phase-01-infrastructure-foundation`
- Surum, faz numarasi ve kisa kapsami ayni etikette bulunur.
- Bosluk yerine `-` kullanilir ve tum branch/tag adlari kucuk harfle yazilir.

## Unreleased

- No unreleased changes yet.

## 0.2.0 - 2026-03-17

- Phase 2 tamamlandi; `identity/auth` modulu `http / app / domain / store` sinirlarina refactor edildi ve eski root-level handler/service/types yapisi kaldirildi.
- Auth veri modeli ve SQLC uretilen kodu `uuid` standardina tasindi; generated kod `internal/gen/sqlc/identity/...` altina ayrildi ve eski `store/sqlc` auth paketi kaldirildi.
- Auth core akislar eklendi: `register`, `login`, `refresh`, `logout`, `logout-all`, `me`.
- Access token + refresh session modeli kuruldu; refresh rotation ve session revoke davranisi eklendi.
- `email verify request`, `email verify`, `forgot password`, `reset password` backend akislarina tek kullanimlik token tablolari ve endpointleri eklendi.
- `identity/access` modulu eklendi; sabit base role modeli (`guest`, `user`, `moderator`, `admin`), permission katalogu, coklu sub-role atamasi ve principal cozumleme omurgasi kuruldu.
- Authorization guard zinciri eklendi; `guest` runtime principal destegi ve admin/base-role tabanli koruma davranisi standartlastirildi.
- Access yonetim endpointleri eklendi: permission listesi, sub-role liste/olusturma, user base role guncelleme ve user sub-role atama/kaldirma.
- `cmd/seed` komutu eklendi; permission katalogu, `manga_moderator`, `comment_moderator`, `chapter_moderator` sub-role'leri ile 1 user, 3 moderator ve 1 admin seed kullanicisi idempotent olarak yukleniyor.
- Docker image artik `seed` binary'sini de uretiyor; `Makefile` icine `seed` hedefi eklendi.
- Auth ve access icin yeni route testleri ve PostgreSQL integration testleri eklendi.
- Dogrulama tekrarlandi: `go test ./...`, `go test -tags=integration ./...`, `docker compose up -d --build`, container ici `migrate up`, `seed`, admin login, `access/me`, refresh, verify/reset ve coklu sub-role smoke akislari canli serviste dogrulandi.

## 0.1.1 - 2026-03-17

- Phase 1 sonrasi uyumlandirma duzeltmeleri yapildi; plan dokumanlari ile calisan kod arasindaki sapmalar kapatildi.
- Resmi migration akisi hosta bagli `localhost` DSN yerine `cmd/migrate` binary'si ve container ici `/app/migrate` komutu uzerinden standartlastirildi.
- `Dockerfile` guncellendi; `api` yanina `migrate` binary'si eklendi, migration dosyalari runtime image'a kopyalandi ve build toolchain `golang:1.26.1-alpine` olarak pinlendi.
- `.dockerignore` eklendi; `.git`, `.planning`, `.codex`, `.env` ve diger gelistirme artifaktlari artik build context'e girmiyor.
- Docker gelistirme akisi host ve container icin ayrildi: `.env.example` host tarafinda `localhost` baz aliyor, `docker-compose.yml` ise `api` servisi icin `NOVASCANS_DB_HOST=postgres` override ediyor.
- Ayni PostgreSQL container'i icinde `novascans_test` veritabani icin init script eklendi ve `test-db-ensure` akisi ile tekrar cagirilabilir hale getirildi.
- Integration testlerin resmi yolu compose agi icindeki gecici Go 1.26.1 container'i olacak sekilde netlestirildi; host tarafinda da `.env` autoload ve test DB olusturma fallback'i eklendi.
- `internal/platform/validation` paketi `go-playground/validator/v10` tabanli wrapper'a donusturuldu; mevcut error mesaji sozlesmesi korundu.
- `config.LoadFromEnv()` artik `.env` dosyasini varsa otomatik yukluyor; boylece lokal `go test`, `go run` ve yardimci komutlar ayni env sozlesmesini paylasiyor.
- Dogrulama tekrarlandi: `go test ./...`, `go test -tags=integration ./...`, Docker agi icinden integration test calistirmasi, `docker compose up -d --build`, container ici `migrate status/up`, `/readyz`, `/metrics` ve auth user/session akisi yeniden dogrulandi.

## 0.1.0 - 2026-03-17

- Phase 1 planlama kurallari ve versiyonlama disiplini tanimlandi.
- Git branch, tag ve dokuman snapshot isimlendirme standardi sikilastirildi.
- `01-01` kapsaminda Go modul iskeleti, `cmd/api` entrypoint'i, `internal/app` bootstrap akisi, ortak modul kontrati ve `identity/auth` ornek modul zinciri eklendi.
- `01-02` kapsaminda `NOVASCANS_` env semasi, fail-fast config yukleyici, ilk PostgreSQL baglanti paketi, `.env.example`, `Dockerfile`, `docker-compose.yml`, `.env`, `Makefile` ve `VERSION` dosyasi eklendi.
- `01-03` kapsaminda `chi` router zinciri, `healthz`, `readyz`, `metrics`, request-id, recover, timeout, logger, metrics middleware'leri ile merkezi error/response yapisi eklendi.
- `01-03` kapsaminda `identity/auth` modulu `ping`, `user create/read`, `session create/revoke` endpointleriyle calisir hale getirildi.
- `01-04` kapsaminda `users`, `auth_password_credentials`, `auth_sessions` migrationlari, schema SQL'i, `sqlc.yaml`, query dosyalari, generated `sqlc` kodu ve PostgreSQL repository/transaction zemini eklendi.
- `01-05` kapsaminda Prometheus uyumlu metrik cikisi, test komut ayrimi, integration test iskeleti ve semver/git naming teslim disiplini tamamlandi.
- Smoke, unit ve integration-tag test yollari calistirildi; Docker compose dosyasi `docker compose config` ile dogrulandi.
- Canli dogrulama tamamlandi: `docker compose up -d --build` calistirildi, goose migration'i compose agi icinden uygulandi, `readyz`, `metrics` ve auth user/session akisi gercek container'lar uzerinde dogrulandi.
