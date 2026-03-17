# Phase 1 Breakdown

## Scope

Bu dosya, `Phase 1: Infrastructure Foundation` icin uygulanabilir alt is kirilimini tanimlar.

## 01-01 Structure and Registration

Hedef: Calisan modul kayit modelini ve fiziksel proje iskeletini netlestirmek.

- `cmd/api/main.go` icin ince entrypoint akisini tanimla.
- `internal/app/bootstrap.go`, `internal/app/modules.go`, `internal/app/routes.go` dosyalarinin sorumluluk sinirlarini yaz.
- `internal/platform` altinda config, db, http, middleware, logger, metrics, validation, events klasorlerini netlestir.
- `internal/modules/identity/auth` modulunu ilk concrete ornek olarak fiziksel agaca yerlestir.
- Modul constructor, dependency gecisi ve `RegisterRoutes(r chi.Router)` modelini kilitle.
- Ortak module kayit kontratini ve app seviyesinde modul toplama mantigini tanimla.
- Bu plan sonunda hangi dosyalarin gercekten olusacagi ve hangilerinin sonraya kalacagi net olsun.

Teslimler:
- Nihai fiziksel klasor agaci
- `main.go`, `bootstrap.go`, `modules.go`, `routes.go` sorumluluk dagilimi
- `identity/auth` modul iskeleti kurali

Dogrulama:
- Fiziksel agac baska moduller gelince tekrar adlandirma gerektirmemeli.
- `identity/auth` ilk concrete modul olarak sisteme eklenebilir durumda olmali.

## 01-02 Runtime and Environment

Hedef: Uygulamanin kalkis, config ve veritabani baglanti omurgasini sabitlemek.

- `.env.example` dosyasini `NOVASCANS_` standardi ile doldur.
- Typed config struct'larini ve fail-fast validation kurallarini yaz.
- `docker-compose.yml` icinde `api + postgres` gelistirme akisini netlestir.
- Ayni PostgreSQL icinde ana DB ve test DB isimlendirme kurallarini yaz.
- PostgreSQL pool ayarlari, SSL modu ve timeout alanlarini standardize et.
- Goose migration'larin ne zaman ve hangi komutla calistirilacagini netlestir.
- Baslangic icin gerekli `Makefile` veya komut seti beklentisini yaz.

Teslimler:
- Env semasi
- Runtime config sozlesmesi
- Docker gelistirme akisi
- Migration calistirma kurali

Dogrulama:
- Eksik zorunlu env ile uygulama fail etmeli.
- Docker gelistirme akisi tekrar edilebilir olmali.

## 01-03 HTTP and API Surface

Hedef: HTTP giris yuzeyi, middleware zinciri ve ortak API sozlesmesini sabitlemek.

- Chi router omurgasini ve global route mount yapisini kur.
- `/healthz`, `/readyz`, `/metrics` endpoint sahipligini netlestir.
- `/api/v1` prefix ve modul bazli route standardini yerlestir.
- Middleware zincirini sira ve sorumluluk bazinda sabitle: request id, real ip, recover, timeout, logger, metrics.
- Merkezi error writer, response helper ve validation akisini tanimla.
- `identity/auth` icin `ping`, `user create/read`, `session create/revoke` endpoint kontratlarini yaz.
- Status code ve response zarfi kurallarini endpoint bazinda eslestir.

Teslimler:
- Route/versioning standardi
- Middleware sirasi
- Error/response sozlesmesi
- Ilk auth endpoint kontratlari

Dogrulama:
- Sistem endpointleri API versiyonundan bagimsiz olmali.
- Tum auth endpointleri ayni response/error standardini kullanmali.

## 01-04 Data and Transactions

Hedef: Auth veri zemini, repository modeli ve transaction sahipligini netlestirmek.

- `users`, `auth_password_credentials`, `auth_sessions` migrationlarini ve isimlendirme kurallarini yaz.
- Temel index, unique constraint ve iliski kurallarini acikla.
- `db/queries/identity/auth` altinda `users.sql`, `password_credentials.sql`, `sessions.sql` dosyalarini tanimla.
- `sqlc.yaml` cikti yolunu ve package isimlendirme standardini sabitle.
- `store/repo.go` ile generated koda ince adaptor katmanini planla.
- Service seviyesinde transaction sahipligini ve `WithTx(...)` kullanim kuralini yaz.
- Repository ve transaction davranisi icin minimum integration test senaryolarini listele.

Teslimler:
- Ilk auth veri modeli
- SQL kaynak duzeni
- Repository + transaction kurali
- Integration test kapsam listesi

Dogrulama:
- Birden fazla yazma islemi ayni service akisi icinde atomic kalmali.
- Auth veri modeli profile gibi sonraki modulleri yanlis yere baglamamali.

## 01-05 Observability and Delivery Discipline

Hedef: Gozlemlenebilirlik, test disiplini ve teslim/version akisini standartlastirmak.

- `EventBus` interface ve in-memory implementasyon kapsamlarini tanimla.
- `slog` tabanli yapilandirilmis loglama alanlarini standardize et.
- Metrics endpoint ve request metriklerinin hangi seviyede toplanacagini netlestir.
- Hizli testler ile integration testleri ayiran komut stratejisini yaz.
- `CHANGELOG.md` giris kurallarini ve semver bump mantigini faz teslim akisina bagla.
- Git branch, release tag ve docs snapshot adlandirma standardini dokumante et.
- Faz sonu `git status` review ve kapatma adimini teslim kontrol listesine sabitle.

Teslimler:
- Observability kurallari
- Test komut politikasi
- Changelog + semver + git naming standardi

Dogrulama:
- Faz kapatilirken hem teknik degisiklik hem de surum kaydi tutarli kalmali.
- Eski faz ve dokuman snapshot'lari tag seviyesinde kolayca karsilastirilabilmeli.

## Done Criteria

- Faz 1 route, data ve delivery omurgasi calisiyor.
- `identity/auth` modulu minimum gercek davranisla ayakta.
- Testler ve changelog/version disiplini tanimli.
- Sonraki auth veya domain fazlari bu temel uzerine kurulabilir.
