# Phase 2 Breakdown

## Scope

Bu dosya, `Phase 2: Identity Core and Access Control` icin uygulanabilir alt is kirilimini tanimlar.

## 02-01 Auth Module Refactor

Hedef: `identity/auth` modulunu Faz 2 buyuklugunu tasiyabilecek ic sinirlara ayirmak.

- HTTP request/response DTO'larini modulun tasima modellerinden ayir.
- Service katmaninin `platform/http` hata tiplerine dogrudan bagliligini kes.
- Repository interface sahipligini `store` katmanindan cikarip `app` veya `domain` tarafina tasi.
- `store` katmaninda sadece persistence implementasyonlarini birak.
- Generated `sqlc` kodunun modulle iliskisini korurken el yazisi koddan gorunur sekilde ayrildigi nihai dizini sabitle.
- `identity/access` modulunun ayni dizin desenini kullanabilmesi icin ortak naming kararlarini kilitle.

Teslimler:
- Nihai auth modul ici klasor yapisi
- DTO / app / domain / store sinirlari
- Repository interface sahipligi kurali
- Generated kod yerlesim karari

Dogrulama:
- Service katmani HTTP veya SQLC generated modellere dogrudan bagli olmamali.
- Faz 2 sonunda `identity/access` ayni deseni tekrar kullanabilmeli.

## 02-02 UUID Standardization

Hedef: Kimlik alanlarini rastgele hex string yerine tutarli `uuid` standardina tasimak.

- Go tarafinda kullanilacak UUID paketini sabitle: `github.com/google/uuid`.
- PostgreSQL tarafinda `uuid` kolon tipi kullanma kararini kilitle.
- Mevcut auth tablolarindaki `id`, `user_id` ve token iliskili kimlik kolonlarinin donusum stratejisini yaz.
- SQL sorgulari, generated kod, domain modelleri ve response maplemelerinde UUID tip gecisini planla.
- Seed ve test verisinin UUID standardiyla uretilmesini sagla.

Teslimler:
- UUID paket karari
- Migration/donusum stratejisi
- Go model ve persistence tip standardi

Dogrulama:
- Yeni kayitlar tutarli UUID ile uretilmeli.
- DB, query ve uygulama katmanlari ayni kimlik semasini kullanmali.

## 02-03 Auth Core Completion

Hedef: Faz 1 iskeletini cekirdek auth davranislarina tamamlamak.

- `register`
- `login`
- `refresh`
- `logout current session`
- `logout all sessions`
- `me`
- session lifecycle ve refresh rotation kurallarini netlestir
- access token ve refresh/session iliskisini sabitle

Teslimler:
- Nihai auth endpoint seti
- Session/token lifecycle kurallari
- Auth middleware giris omurgasi

Dogrulama:
- Login olan kullanici `me` endpoint'ine erisebilmeli.
- Refresh edilen session kurali deterministic olmali.
- Revoke edilen session tekrar kullanilamamali.

## 02-04 Verification and Reset Flows

Hedef: Mail provider olmadan da auth tamamlama akislarini backend seviyesinde kurmak.

- `email verify request`
- `email verify`
- `forgot password`
- `reset password`
- verification/reset token veri modelini ayir
- token olusturma, tuketme, gecerlilik suresi ve yeniden kullanimi engelleme kurallarini yaz
- dis mail provider bu fazda olmadan generic request response politikasini sabitle

Teslimler:
- Verification token zemini
- Password reset token zemini
- Request/consume endpointleri

Dogrulama:
- Token olmadan verify/reset tamamlanmamali.
- Token ikinci kez kullanilamamali.
- Mail delivery olmadan backend akislar test edilebilir kalmali.

## 02-05 Access and Authorization

Hedef: Sabit roller + coklu alt roller mantigiyla `identity/access` modulunu kurmak.

- Sabit sistem rollerini tanimla: `guest`, `user`, `moderator`, `admin`
- `guest` icin runtime principal davranisini netlestir.
- `moderator` rolunun moderasyon alani icin kilit rol oldugunu sabitle.
- Permission katalog yapisini tanimla.
- Admin tarafindan olusturulan alt rollerin mevcut permission katalogundan secim yapmasi kuralini yaz.
- Kullaniciya birden fazla alt rol atanabilmesi kuralini uygula.
- Authorization middleware/policy kontrol akisini tanimla.

Teslimler:
- `identity/access` veri modeli
- Principal ve permission karar modeli
- Authorization middleware kurali

Dogrulama:
- `moderator` tek basina panel girisini saglamali.
- Alt roller panel ici yetki kapsamlarini belirlemeli.
- `admin` tam yetki bypass veya tam izin cozumuyle tutarli davranmali.

## 02-06 Seed, Tests and Delivery

Hedef: Faz 2 teslimini tekrar uretilebilir seed, test ve dokuman disiplini ile kapatmak.

- Sistem rol seed'lerini ekle: `user`, `moderator`, `admin`
- Alt rol seed'lerini ekle: `manga_moderator`, `comment_moderator`, `chapter_moderator`
- Uygun moderator seed kullanicilarina birer alt rol ata
- Auth + access entegrasyon testlerini yaz
- Authorization ve auth core smoke senaryolarini belirle
- Changelog, roadmap, state ve version kapanis kurallarini uygula

Teslimler:
- Seed veri paketi
- Entegrasyon ve smoke test kapsam listesi
- Faz kapanis dokumantasyonu

Dogrulama:
- Temiz ortamda seed sonrasi roller ve alt roller beklenen sekilde gorulmeli.
- Auth ve access birlikte calisan temel UAT akislari testten gecmeli.

## Done Criteria

- Auth core akislar production'a yakin cekirdek seviyede tamamlanmis.
- Access/RBAC omurgasi sabit roller ve coklu alt roller ile ayaga kalkmis.
- UUID standardi kimlik katmaninda yerlestirilmis.
- Sonraki manga, moderation ve admin fazlari bu auth/access temelini tekrar kullanabilir.
