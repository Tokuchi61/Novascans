# Novascans

## What This Is

Novascans icin tum urunu bir defada planlamak yerine adim adim ilerleyen bir backend-first kurulum yapiliyor. Bu ilk milestone yalnizca temel altyapiya odaklanir: proje iskeleti, router, veri erisim katmani, PostgreSQL, Docker, migration akisi, gozlemlenebilirlik ve sonraki modullerin oturacagi ortak backend konvansiyonlari.

## Core Value

Sonraki auth, manga, topluluk ve admin modullerini tekrar temel degistirmeden tasiyabilecek temiz, hizli ve genisletilebilir bir backend omurgasi kurmak.

## Requirements

### Validated

(None yet - ship to validate)

### Active

- [ ] Go, Docker ve PostgreSQL tabanli temel backend altyapisi calisir.
- [ ] Router, veri erisim katmani, migration ve config konvansiyonlari netlesir.
- [ ] Sonraki moduller icin tekrar kullanilabilir proje iskeleti, hata modeli ve gelistirme akisi tanimlanir.

### Out of Scope

- Auth, manga, topluluk ve admin is mantigi - Bu basliklar sonraki fazlarda tek tek ele alinacak.
- Kullanici veya admin web arayuzu - Once backend temeli sabitlenecek.
- Oyunlastirma, VIP ve rollout is kurallari - Temel altyapi netlestikten sonra planlanacak.

## Context

Planlama yaklasimi degisti: artik tum ozellikler tek roadmap icinde zorlanmayacak. Bunun yerine her seferinde dar bir problem secilecek, planlanacak, uygulanacak ve ondan sonra siradaki alan secilecek. Bu ilk alan "temel backend altyapisi" olarak secildi. Ancak bu, gecici veya oyuncak seviyesinde bir kurulum anlamina gelmeyecek; Faz 1 kararlarinin amaci ileride kritik mimari veya veri butunlugu sorunlari yaratmayacak net bir temel kurmaktir.

Teknik sabitler su an icin nettir: Go kullanilacak, servis Docker icinde calisacak ve PostgreSQL ana veritabani olacak. Ilk altyapi fazi icin cekirdek secimler de netlesti: router olarak `chi`, veri erisim katmani olarak `sqlc`, migration araci olarak `goose` kullanilacak. Surum politikasi guncel stabil surumleri hedefleyecek; bu planlama turunda referans Go surumu `1.26.1`, PostgreSQL surumu `18.3` olarak alindi.

Modul organizasyonu icin de temel karar alindi: fiziksel yapi `kategori -> gercek modul` seklinde kurulacak. Ust seviye klasorler yalnizca gruplama amaci tasiyacak; is kurallari alttaki gercek modullerde yasayacak. Ornek hedef aileler su sekildedir:
- `identity/auth`, `identity/access`
- `user/profile`, `user/settings`, `user/privacy`
- `community/wall`, `community/social`, `community/messaging`, `community/comment`
- `content/manga`, `content/chapter`, `content/reading`, `content/library`
- `progression/entitlement`, `progression/gamification`
- `backoffice/admin`, `backoffice/moderation`

Bu karar, `user` gibi tek bir buyuk modul altina wall, friends ve dm gibi davranislari yigmayi engeller. Altyapi tercihleri icin ek kararlar da alindi:
- Config yalnizca `env` tabanli olacak; `.env.example` bulunacak, uygulama acilisinda typed config parse edilip eksik zorunlu alanlarda fail-fast davranacak.
- SQL kaynaklari `db/queries/<kategori>/<modul>` altinda tutulacak; `sqlc` ile uretilen Go kodu ilgili modullerin `store/sqlc` dizinlerinde yer alacak.
- Event altyapisi Faz 1'de `EventBus` interface'i ve basit bir `in-memory` implementasyonuyla kurulacak; dis broker entegrasyonu sonra eklenecek.
- Docker gelistirme ortami varsayilan olarak `api + postgres` seklinde olacak; testler ayni PostgreSQL container'inda ayri bir test veritabani kullanacak.

Bootstrap yapisi icin secilen yon su olacak:
- `cmd/api/main.go` ince bir entrypoint olarak kalacak.
- `internal/app/bootstrap.go` ortak bilesenlerin kurulumunu yapacak.
- `internal/app/modules.go` modullerin kaydini ve baglanmasini yonetecek.
- `internal/app/routes.go` global route mount islemlerini yuruterek `/api/v1`, `/healthz`, `/readyz`, `/metrics` rotalarini toplayacak.

Transaction siniri icin secilen yon su olacak:
- Transaction yonetimi `service/use-case` seviyesinde yer alacak.
- Repository katmani kendi basina transaction baslatmayacak.
- `sqlc` sorgulari transaction icinde `WithTx(...)` yaklasimiyla kullanilacak.
- Birden fazla repository'yi kapsayan yazma akislarinda atomic davranis service katmaninda koordine edilecek.

Hata ve response standardi icin secilen yon su olacak:
- Tum API hatalari tek bir error zarfi ile donulecek: `code`, `message`, opsiyonel `details`.
- Ham panic, stack trace veya veritabani hatasi istemciye sizmayacak.
- Basari cevaplari varsayilan olarak `data` zarfi ile donecek; ekstra `meta` alani ihtiyaca gore sonradan eklenecek.
- Validation hatalari alan bazli `details.fields` yapisinda donecek.
- Handler katmani hata formatini elle kurmayacak; merkezi bir error writer veya response helper kullanacak.

Test tabani icin secilen yon su olacak:
- Test katmanlari `unit`, `handler/http`, `integration` ve `smoke/startup` olarak ayrilacak.
- Faz 1'de zorunlu guvence alanlari config parse, migration, db baglantisi, repository akislari, transaction davranisi, error/response writer, `healthz`, `readyz` ve en az bir route/middleware akisi olacak.
- Unit ve handler testleri hizli calisacak; integration testleri ayri komut veya tag ile kosturulecek.
- Integration testler gercek PostgreSQL test veritabani kullanacak.
- Yuksek coverage hedefi yerine kritik altyapi davranislarinin guvencesi hedeflenecek.

Modul kayit modeli icin secilen yon su olacak:
- Her gercek modul constructor ile dependency alacak.
- `app` katmani modul olusumunu merkezden yonetecek.
- Moduller Faz 1'de minimum arayuz ile calisacak: `Key()` ve `RegisterRoutes(r chi.Router)`.
- Route mount islemi merkezi olarak `internal/app/routes.go` icinde gerceklesecek.
- Job registration, event subscriber registration veya health contributor gibi gelismis lifecycle alanlari ilk fazda zorunlu tutulmayacak.
- Ilk fiziksel ornek modul `system` yerine `identity/auth` olacak; boylece modul deseni gercek bir domain alani uzerinde dogrulanacak.

Config semasi ve env isimlendirme standardi icin secilen yon su olacak:
- Tum env degiskenleri `NOVASCANS_` prefix'i ile baslayacak.
- Config yapisi Go tarafinda gruplu typed struct'lar ile temsil edilecek: `App`, `HTTP`, `DB`, `Log`, `Metrics`.
- Env isimleri duz ama gruplu olacak: `NOVASCANS_HTTP_PORT`, `NOVASCANS_DB_HOST` gibi.
- Config sadece altyapi ve runtime davranisi tasiyacak; VIP, XP, rollout gibi is kurallari env uzerinden yonetilmeyecek.
- `.env.example` tum zorunlu alanlari gosterecek; uygulama eksik zorunlu alanlarda startup'ta fail edecek.

Middleware seti ve siralamasi icin secilen yon su olacak:
- Faz 1 temel middleware zinciri `request_id`, `real_ip`, `recover`, `timeout`, `logger`, `metrics` seklinde kurulacak.
- Panic korumasi ve timeout global seviyede saglanacak.
- Request id her istekte uretilecek ve log/response baglamina tasinacak.
- Auth, rate limit, CORS, CSRF ve compression middleware'leri Faz 1 temel zincirine dahil edilmeyecek; ama sonra kolayca eklenebilir olacak.

API route yapisi ve versioning standardi icin secilen yon su olacak:
- Is kurali endpoint'leri sabit olarak `/api/v1/<module>` yapisini kullanacak.
- Sistem endpoint'leri API disinda kok seviyede kalacak: `/healthz`, `/readyz`, `/metrics`.
- Fiziksel kategori klasorleri public URL yapisina yansimayacak; URL'lerde yalnizca gercek modul isimleri gorunecek.
- Versioning route tabanli olacak; header veya query param tabanli versioning ile baslanmayacak.
- Her modul kendi base path sahipligini `RegisterRoutes(...)` icinde tasiyacak.

Migration naming ve `sqlc` duzen standardi icin secilen yon su olacak:
- Migration dosyalari `goose` ile sirali numara + kisa aciklama formatinda adlandirilacak: `000001_init_extensions.sql` gibi.
- SQL kaynaklari `db/queries/<kategori>/<modul>` yapisinda tutulacak.
- `sqlc` ile uretilen kod ilgili modullerin `store/sqlc` dizinlerinde yer alacak.
- Elle yazilan repository katmani generated koda ince bir adaptor olarak ayni modulde `store/repo.go` altinda bulunacak.
- Timestamp tabanli migration adlandirmasi ile baslanmayacak.

Hala netlestirilmesi gereken alan ilk fiziksel klasor agaci ve placeholder dosyalardir.
`identity/auth` modulunun Faz 1 kapsami icin secilen yon su olacak:
- Faz 1'de `identity/auth` tam auth feature setini bitirmeyecek; ama gercek bir domain modul olarak kurulacak.
- Modul; route registration, handler-service-repo zinciri, `sqlc` akisi ve migration sahipligini gercekten gosterecek.
- Auth alani icin temel veri modeli ve tablo zemini bu fazda baslayabilecek; ancak register/login/refresh/logout is kurallarinin tam uygulanmasi sonraki fazda genisletilecek.
- Faz 1'de modulu dogrulamak icin en az bir basit calisan endpoint ve repository akisi bulunacak.

Ilk fiziksel klasor agaci icin secilen yon su olacak:
- `cmd/api` uygulama giris noktasi olacak.
- `internal/app` bootstrap, module registration ve route mount alanini tasiyacak.
- `internal/platform` ortak altyapiyi tasiyacak: config, db, http, middleware, logger, metrics, validation, events.
- `internal/modules/identity/auth` Faz 1'deki ilk gercek modul olacak; `store/repo.go` ve `store/sqlc/` yapisini kullanacak.
- SQL kaynaklari `db/queries/identity/auth` altinda, migration dosyalari `db/migrations` altinda duracak.
- Yalnizca gercekten kullanilacak klasor ve dosyalar olusturulacak; gelecekteki moduller icin bos placeholder kod acilmayacak.

`identity/auth` icin Faz 1 tablo ve query kapsami icin secilen yon su olacak:
- `users` tablosu auth alaninin sahip oldugu cekirdek identity tablosu olacak.
- Parola verisi `users` tablosundan ayrilacak ve `auth_password_credentials` altinda tutulacak.
- Session/refresh temeli `auth_sessions` tablosu ile kurulacak.
- `profile` benzeri alanlar auth tablosuna konmayacak; sonraki moduller `users.id` uzerine insa edilecek.
- Faz 1'de `email_verification_tokens`, `password_reset_tokens`, `oauth_identities`, `mfa` ve benzeri genislemeler acilmayacak.
- Ilk SQL kaynaklari `users.sql`, `password_credentials.sql` ve `sessions.sql` dosyalari olarak acilacak.

`identity/auth` icin Faz 1 endpoint ve minimal davranis kapsami icin secilen yon su olacak:
- `GET /api/v1/auth/ping` modulin route registration ve response standardini dogrulayacak.
- `POST /api/v1/auth/users` identity + credential olusturma omurgasini gosterecek; ancak tam kayit urun akisi seviyesine zorlanmayacak.
- `GET /api/v1/auth/users/{id}` repository, sqlc ve response standardini gercek bir okuma akisi ile dogrulayacak.
- `POST /api/v1/auth/sessions` session olusturma omurgasini dogrulayacak; ancak tam login hardening seti bu fazda zorunlu olmayacak.
- `DELETE /api/v1/auth/sessions/{id}` session revoke temelini gosterecek.
- Bu endpointler production-ready auth urununun tamami olarak degil, auth modulunun gercek veri ve servis omurgasini ispatlayan minimum davranis seti olarak ele alinacak.

Surumleme ve faz tamamlama protokolu icin secilen yon su olacak:
- Tum asama ve uygulama degisiklikleri kok dizindeki `CHANGELOG.md` dosyasina detayli sekilde islenecek.
- Surumleme `semver` ile yonetilecek; proje erken donemde `0.x.y` cizgisinde ilerleyecek.
- Geriye donuk uyumlu yeni yetenekler `minor`, duzeltmeler ve davranis kirilimi yaratmayan iyilestirmeler `patch` olarak islenecek.
- Faz tamamlama akisi standart olacak: dokumanlari incele, mevcut yapilari incele, fazi uygula, test et, changelog ve surumu guncelle, git durumunu gozden gecir, fazi kapat.
- Bu kural planning notlari icin degil, tamamlanan uygulama fazlari icin zorunlu kabul edilecek.
- Git tarafinda karsilastirma branch adlariyla degil, okunabilir tag ve release etiketleriyle yapilacak.
- Kalici ana branch `main` olarak korunacak; faz, feature ve dokuman calismalari ayri branch adlariyla tasinacak: `phase/01-infrastructure-foundation`, `feature/identity-auth-sessions`, `docs/phase-01-foundation-rules` gibi.
- Release etiketleri `release/v0.1.0-phase-01-infrastructure-foundation`, dokuman snapshot'lari `docs/v0.1.0-phase-01-infrastructure-foundation`, plan snapshot'lari gerekirse `plan/v0.1.0-phase-01-infrastructure-foundation` formatini izleyecek.
- Repo referansi `https://github.com/Tokuchi61/Novascans` olarak kabul edilecek.

Bu asama icin su an kritik acik soru kalmadi; bir sonraki mantikli adim Phase 1 planlarini alt islere bolmek.

## Constraints

- **Tech stack**: Go, Docker, PostgreSQL - Temel platform degismeyecek.
- **Planning style**: Dar kapsamli fazlar - Tum urunu bastan detaylandirmak yerine parca parca ilerleniyor.
- **Maintainability**: Sonraki moduller temeli bozmayacak sekilde eklenmeli - Faz 1'in degeri tekrar kullanilabilirlik.
- **Operational simplicity**: Lokal kurulum ve container akisi sade olmali - Gelistirme hizi erken donemde kritik.
- **Future readiness**: Auth ve domain modulleri sonraki fazlarda kolayca takilabilmeli - Erken kararlar kilitlenme yaratmamali.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Planlama tek seferlik buyuk roadmap yerine faz faz ilerleyecek | Kapsam sisip karar kalitesi dusuyordu | - Pending |
| Ilk milestone yalnizca temel altyapi olacak | En once tekrar kullanilabilir backend zemini kurulmak isteniyor | - Pending |
| Faz 1 sade olacak ama gecici veya oyuncak seviyesinde olmayacak | Erken kolaycilik ileride kritik hata ve yeniden yazim riski yaratmamali | - Pending |
| Go + Docker + PostgreSQL cekirdek teknoloji olarak korunacak | Performans ve isletim tercihi sabit | - Pending |
| Router olarak chi secilecek | Hafif, composable ve `net/http` uyumlu bir omurga isteniyor | - Pending |
| Veri erisim katmani olarak sqlc secilecek | SQL kontrolunu kaybetmeden type-safe kod uretimi isteniyor | - Pending |
| Migration araci olarak goose secilecek | SQL-first akisla uyumlu, sade migration yonetimi gerekiyor | - Pending |
| Faz 1 referans surumleri Go 1.26.1 ve PostgreSQL 18.3 olacak | Planlanan altyapi guncel stabil surumler uzerine kurulacak | - Pending |
| Modul organizasyonu kategori -> gercek modul seklinde kurulacak | Ust seviye gruplama saglarken mega-modul olusmasini engellemek gerekiyor | - Pending |
| Config yalnizca env tabanli olacak | Container ve deploy ortamlarinda en sade ve dogal yaklasim bu | - Pending |
| SQL kaynaklari merkezi, generated kod modul icinde tutulacak | SQL sahipligi ve modul sinirlari birlikte korunmali | - Pending |
| EventBus interface'i ve in-memory implementasyon Faz 1'e dahil olacak | Moduller arasi gevsek baglanti icin broker'siz bir baslangic yeterli | - Pending |
| Docker varsayilaninda api + postgres olacak, test veritabani ayni container icinde ayrilacak | Ilk asamada compose karmasasini buyutmadan test izolasyonu saglamak yeterli | - Pending |
| Bootstrap yapisi `cmd/api` + `internal/app` ayrimi ile kurulacak | Main ince kalmali, wiring tek merkezde toplanmali | - Pending |
| Transaction yonetimi service/use-case seviyesinde olacak | Coklu repository akislarinda atomic davranis merkezi olarak koordine edilmeli | - Pending |
| API hata ve response formati merkezi ve tek tip olacak | Sonraki tum moduller ayni sozlesmeye oturmali | - Pending |
| Test tabani unit, http, integration ve smoke olarak ayrilacak | Hizli geri bildirim ile gercek altyapi guvencesi birlikte korunmali | - Pending |
| Modul kayit modeli constructor + merkezi route registration ile kurulacak | Modul wiring'i dagilmadan buyume noktasi saglanmali | - Pending |
| Config semasi `NOVASCANS_` prefix'i ve typed gruplu struct'lar ile kurulacak | Cevre degiskenleri net, tasinabilir ve parse edilebilir olmali | - Pending |
| Temel middleware zinciri request-id ile baslayip metrics ile bitecek | Gozlemlenebilirlik ve hata guvencesi tum isteklerde tutarli olmali | - Pending |
| API route yapisi `/api/v1/<module>` standardi ile kurulacak | Public API sozlesmesi fiziksel klasor yapisindan bagimsiz ve sade olmali | - Pending |
| Migration naming sirali numara + kisa aciklama formatinda olacak | Okunabilir ve dallarda yonetilebilir bir migration akisi gerekli | - Pending |
| SQL kaynaklari merkezi, generated sqlc kodu modul icinde olacak | SQL sahipligi ile modul sinirlari birlikte korunmali | - Pending |
| Ilk concrete modul `identity/auth` olacak ama tam auth feature seti bu fazda bitmeyecek | Altyapiyi gercek bir domain uzerinde dogrularken kapsam kaymasi olmamali | - Pending |
| Ilk fiziksel agac `cmd/api`, `internal/app`, `internal/platform` ve `internal/modules/identity/auth` ekseninde kurulacak | Yalnizca kullanilacak iskelet acilmali, gelecege donuk bos kod uretilmemeli | - Pending |
| Auth alaninin ilk tablo zemini `users`, `auth_password_credentials` ve `auth_sessions` ile kurulacak | Identity cekirdegi ile auth verileri birbirine karismadan buyumeli | - Pending |
| Auth modulunun Faz 1 endpoint seti ping, user create/read ve session create/revoke ile sinirlanacak | Gercek modul davranisi ispatlanirken kapsam kontrollu kalmali | - Pending |
| Tum tamamlanan uygulama fazlari `CHANGELOG.md` ve semver ile kaydedilecek | Surum gecmisi ve teslim disiplini proje boyunca izlenebilir olmali | - Pending |
| Faz tamamlama protokolu dokuman inceleme -> uygulama -> test -> version/changelog -> git review sirasi ile ilerleyecek | Her fazda ayni teslim kalitesi korunmali | - Pending |
| Git branch, tag ve release adlari tip + surum + faz + kapsam formatinda okunabilir tutulacak | Eski surumler ve fazlar arasi karsilastirma kolay olmali | - Pending |

---
*Last updated: 2026-03-17 after infrastructure scoping*
