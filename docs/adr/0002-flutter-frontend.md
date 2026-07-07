# ADR 0002 — Flutter pour le mobile multi-plateforme

| Statut | Date | Auteurs |
|---|---|---|
| Adopté | 2026-01-15 | Équipe StreamPulse |

## Contexte

L'application cliente doit fonctionner sur **iOS et Android** avec une seule équipe et une seule base de code. Les exigences du sujet incluent :

- **Lecture audio en arrière-plan** (background services, gestion d'interruptions).
- **State management robuste** (Bloc, Riverpod ou équivalent).
- **60 FPS** maintenu pendant le streaming.
- Cours dédié au framework (Thomas Coichot).

## Décision

Le frontend mobile est développé en **Flutter (≥ 3.22) avec Dart (≥ 3.4)**.

## Alternatives considérées

| Stack | Forces | Faiblesses | Verdict |
|---|---|---|---|
| **Flutter** | Une seule base, performance proche du natif, hot reload, écosystème mature pour audio (`just_audio`, `audio_service`) | Apps plus volumineuses (~ 15 Mo de base), nécessite Dart | ✅ Retenu |
| React Native | Écosystème JS, partage de logique avec web | Bridge avec le natif moins performant, audio en arrière-plan plus complexe | ❌ |
| Natif iOS + Android (Swift + Kotlin) | Performance maximale | Double équipe, double code, double timeline | ❌ |
| Kotlin Multiplatform | Code partagé, natif | UI non partagée (à reconstruire), maturité audio moindre | ❌ |
| PWA | Pas d'app store | Limitations audio background sur iOS, intégration système faible | ❌ |

## Conséquences

### Positives

- **Une seule base** Dart → vitesse de développement.
- `just_audio` + `audio_service` couvrent nativement le streaming HTTP et la lecture en arrière-plan.
- **Hot reload** permettant des itérations UI rapides.
- Cours dédié → alignement pédagogique.
- Excellent support pour BLoC (voir ADR 0004) et `go_router`.

### Négatives

- Apprentissage de **Dart** spécifique (mais proche de TS/Kotlin).
- Taille de l'APK plus élevée qu'une app native pure.
- Quelques composants natifs (mic capture broadcaster) nécessitent du code Kotlin/Swift via `platform channels`.

### Neutres

- Nécessite un Mac pour builder iOS.
- Build Flutter Web possible mais hors périmètre v1.

## Liens

- [Documentation Flutter](https://docs.flutter.dev/)
- [ADR 0004 — BLoC](./0004-bloc-state-management.md)
- [`just_audio`](https://pub.dev/packages/just_audio)
- [`audio_service`](https://pub.dev/packages/audio_service)
