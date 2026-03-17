# Phase Workflow

## Purpose

Bu dosya, tamamlanan uygulama fazlarinda izlenecek standart teslim protokolunu tanimlar.

## Standard Flow

1. Ilgili planlama dokumanlarini incele.
2. Mevcut kod yapisini ve etkilenmis alanlari incele.
3. Faz kapsamini ve uygulanacak degisiklikleri netlestir.
4. Uygulamayi gelistir.
5. Gerekli testleri calistir.
6. `CHANGELOG.md` dosyasini ve surum bilgisini guncelle.
7. `git status` ile degisiklikleri gozden gecir.
8. Fazi kapat.

## Git Naming Rules

- Repo reference: `https://github.com/Tokuchi61/Novascans`
- Kalici ana branch yalnizca `main` olarak korunur; surum veya kapsam bilgisi branch adina yazilmaz.
- Faz calisma branch'leri okunabilir ve tek amacli olur:
  - `phase/01-infrastructure-foundation`
  - `phase/02-auth-core`
- Faz icindeki daha kucuk uygulama branch'leri gerekiyorsa kapsam bazli ayrilir:
  - `feature/identity-auth-sessions`
  - `feature/platform-metrics`
  - `docs/phase-01-foundation-rules`
- Faz veya release karsilastirmalari branch adlariyla degil, tag ve release etiketleriyle yapilir.
- Kod release etiketleri su formati izler:
  - `release/v0.1.0-phase-01-infrastructure-foundation`
  - `release/v0.2.0-phase-02-auth-core`
- Dokuman snapshot etiketleri su formati izler:
  - `docs/v0.1.0-phase-01-infrastructure-foundation`
  - `docs/v0.2.0-phase-02-auth-core`
- Gerekirse plan snapshot'lari da ayri etiketlenir:
  - `plan/v0.1.0-phase-01-infrastructure-foundation`
- Commit mesajlari da okunabilir ve karsilastirma dostu olur:
  - `phase(01-03): add chi router and middleware chain`
  - `docs(phase-01): refine infrastructure workflow rules`
- Dosya, branch, tag ve release adlarinda bosluk kullanilmaz; kucuk harf ve `-` kullanilir.
- Surum, faz numarasi ve kisa kapsam bilgisi tag veya release etiketinde mutlaka yer alir.

## Notes

- Bu protokol planning tartismalari icin degil, tamamlanan uygulama fazlari icin zorunludur.
- Uzak repoya push islemi remote yapisi mevcutsa ve o asama icin is akisi gerektiriyorsa yapilir.
- Semver ve changelog ayrintilari kok dizindeki `CHANGELOG.md` dosyasinda tanimlidir.
