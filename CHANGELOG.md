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

- Phase 1 planlama kurallari ve versiyonlama disiplini tanimlandi.
- Git branch, tag ve dokuman snapshot isimlendirme standardi sikilastirildi.
