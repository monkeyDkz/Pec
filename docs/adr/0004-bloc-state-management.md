# ADR 0004 — BLoC pour la gestion d'état Flutter

| Statut | Date | Auteurs |
|---|---|---|
| Adopté | 2026-01-20 | Équipe StreamPulse |

## Contexte

L'application Flutter gère plusieurs flux d'état complexes :

- **Authentification** (loading → authenticated / error).
- **Lecture audio** (idle / loading / playing / paused / buffering / error).
- **Listing de streams** (chargement, pull-to-refresh, mise à jour temps réel).
- **Broadcaster** (offline / live / publishing chunks).

Le sujet impose un pattern « robuste » (Bloc ou Riverpod). Le cours Thomas Coichot privilégie **BLoC**.

## Décision

La gestion d'état utilise **flutter_bloc (v8)** avec une séparation stricte :

- **Events** : actions utilisateur ou systèmes.
- **States** : représentations immuables de l'UI.
- **Bloc** : transition Event × State → nouvel State.

## Alternatives considérées

| Approche | Forces | Faiblesses | Verdict |
|---|---|---|---|
| **flutter_bloc** | Séparation explicite Event/State, prédictible, testable via `bloc_test`, ouvert | Plus de boilerplate qu'un simple `setState` ou Riverpod | ✅ Retenu |
| Riverpod | Plus ergonomique, composable, gestion fine de dépendances | Pattern moins enseigné dans le cours, magie implicite | ❌ |
| Provider seul | Suffisant pour cas simples | Insuffisant pour stream temps réel + side effects | ❌ |
| GetX | Concis, productif | Magie, anti-pattern selon plusieurs guidelines Flutter, communauté divisée | ❌ |
| `setState` | Trivial | Ingérable au-delà de 2 écrans | ❌ |

## Conséquences

### Positives

- Tests unitaires des Blocs via `bloc_test` : trivial à 80 % de couverture.
- Logique métier **séparée des widgets** → meilleurs widgets, plus testables.
- Patron alignée avec le cours (Coichot).
- Excellent debugging via `BlocObserver`.

### Négatives

- Boilerplate (Event class, State class, transitions).
- Apprentissage initial pour des contributeurs venant de Redux ou Riverpod.

### Neutres

- Une feature = un dossier `feature/<x>/{bloc, repositories, screens, models}/`.
- Repositories injectés dans les Blocs (DI manuelle via constructeur), pas de framework DI.

## Liens

- [flutter_bloc](https://bloclibrary.dev/)
- [`bloc_test`](https://pub.dev/packages/bloc_test)
- [ADR 0002 — Flutter](./0002-flutter-frontend.md)
