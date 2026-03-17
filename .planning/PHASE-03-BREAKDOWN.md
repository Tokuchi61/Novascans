# Phase 3 Breakdown

## Scope

Bu dosya, `Phase 3: Account Core` icin uygulanabilir alt is kirilimini tanimlar.

## Locked Decisions

- `auth` kullanici kimligini olusturmaya devam eder; `account` yeni kullanici olusturmaz
- `users` tablosu auth sahipliginde kalir
- `account` verisi `users.id` uzerine bagli tablolarda tutulur
- register akisi default account kayitlarini tek transaction icinde ve senkron olarak olusturur
- bu entegrasyon event tabanli degil, port/servis uzerinden kurulur
- ilk faz account kapsaminda `profile`, `settings`, `privacy` bulunur
- public profile username uzerinden okunur
- `wall`, `friends`, `follow`, `dm`, `library`, `history` bu fazin disindadir

## 03-01 Account Domain Modeling

Hedef: `user/account` alaninin veri ve davranis sinirlarini Phase 3 icin netlestirmek.

- `profile`, `settings` ve `privacy` veri modellerini kesinlestir.
- Auth ile account sahiplik sinirini sabitle: `users` auth'da kalir, account bagli tablolar uzerinde calisir.
- Public profile icin `username` alanini ve benzersizlik kuralini sabitle.
- Privacy davranisini en az `public`, `authenticated`, `private` seviyelerinde netlestir.
- Bu faza dahil olmayan alanlari net ayir: wall, friends, follow, dm, library, history.

Teslimler:
- Nihai domain modeli
- Enum/field karar listesi
- Faz disi alanlar listesi

Dogrulama:
- Veri modeli own-account write ve public profile read akislarini birlikte desteklemeli.
- Faz disi alanlar sonraki moduller icin acik ama simdilik bagimsiz kalmali.

## 03-02 Persistence and Module Skeleton

Hedef: Account alani icin schema, migration ve module iskeletini auth/access standartlariyla uyumlu kurmak.

- `user/account` module iskeletini olustur.
- Gerekli tablolar, index'ler ve foreign key'leri migration olarak ekle.
- `db/queries/user/account` sorgularini yaz.
- `sqlc` generated kodu yeni access/auth yerlesim standardiyla ayni sekilde uret.
- Repository ve service sinirlarini netlestir.
- Register akisina account bootstrap eklemek icin `AccountProvisioner` portunu ve transaction entegrasyonunu kur.

Teslimler:
- Module iskeleti
- Migration dosyalari
- Query dosyalari
- Generated SQLC kodu

Dogrulama:
- `sqlc generate` temiz gecmeli.
- Migration ile temiz veritabani ayaga kalkmali.

## 03-03 Account APIs

Hedef: Kullanicinin kendi account verisini yonetebilecegi ve public profile gorebilecegi API yuzeyini acmak.

- `GET /api/v1/account/me`
- `GET/PATCH /api/v1/account/profile`
- `GET/PATCH /api/v1/account/settings`
- `GET/PATCH /api/v1/account/privacy`
- `GET /api/v1/account/profile/{username}`
- Privacy kurallarini public profile akisina uygula.

Teslimler:
- Account endpoint seti
- DTO ve response maplemeleri
- Read-path testleri
- Write-path testleri

Dogrulama:
- Kullanici kendi account verisini gorebilmeli ve guncelleyebilmeli.
- Public profile akisi username ile dogrulanmali.

## 03-04 Auth-Account Bootstrap and Privacy

Hedef: Auth ile account arasindaki baglantiyi tek transaction icinde dogrulamak ve privacy davranisini netlestirmek.

- Register akisinda default profile/settings/privacy kayitlarini olustur.
- Username uniqueness ve default deger politikalarini uygula.
- Public profile gorunurlugunu `public`, `authenticated`, `private` seviyelerinde dogrula.
- Account bootstrap basarisiz olursa auth kaydinin da rollback oldugunu test et.

Teslimler:
- Auth-account entegrasyonu
- Privacy ve bootstrap testleri
- Username davranis kurallari

Dogrulama:
- Register sonrasi account kayitlari otomatik olusmali.
- Bootstrap hatasi transaction rollback ile sonuclanmali.
- Privacy seviyesi public profile sonucunu dogru etkilemeli.

## 03-05 Seed, Verification and Delivery

Hedef: Faz 3 teslimini tekrar uretilebilir seed ve test disiplini ile kapatmak.

- Account tablolari icin geriye donuk seed uyumunu koru.
- Register bootstrap smoke senaryosunu yaz.
- Public profile ve own-account update smoke veya integration senaryolarini yaz.
- Changelog, roadmap, state ve version guncellemelerini yap.

Teslimler:
- Seed veri paketi
- Smoke/integration dogrulamalari
- Faz kapanis dokumani

Dogrulama:
- Temiz ortamda register sonrasi account bootstrap akisi hemen calismali.
- Own-account update ve public profile read birlikte dogrulanmali.

## Done Criteria

- Account cekirdek veri modeli ayakta.
- Register -> account bootstrap zinciri calisiyor.
- Kendi account verisini yonetme ve public profile okuma akislar calisiyor.
- Sonraki manga, wall, friends, library ve history fazlari icin uygun genisleme zemini ayrilmis oluyor.
