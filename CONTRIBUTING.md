# Contribuer à StreamPulse

Merci de l'intérêt porté au projet ! Ce guide décrit comment proposer une contribution, du clone au merge.

## Sommaire

- [Code de conduite](#code-de-conduite)
- [Workflow général](#workflow-général)
- [Préparer son environnement](#préparer-son-environnement)
- [Conventions de commit](#conventions-de-commit)
- [Stratégie de branches](#stratégie-de-branches)
- [Pull Requests](#pull-requests)
- [Code Review](#code-review)
- [Tests](#tests)
- [Documentation](#documentation)
- [Signaler un bug](#signaler-un-bug)

---

## Code de conduite

Le projet adhère à un principe simple : **bienveillance, rigueur, ouverture**. Tout comportement irrespectueux est inacceptable.

## Workflow général

1. Créer un *issue* GitHub décrivant le besoin (sauf typo / fix trivial).
2. Forker (ou créer une branche `feature/<slug>` sur le dépôt).
3. Coder + tester + documenter.
4. Ouvrir une *Pull Request* vers `develop`.
5. Réviser et merger après validation CI + 1 review minimum.

## Préparer son environnement

```bash
git clone https://github.com/<org>/streampulse.git
cd streampulse
cp .env.example .env
docker compose up -d
```

Configurer la signature GPG des commits :

```bash
gpg --full-generate-key
git config user.signingkey <KEY_ID>
git config commit.gpgsign true
```

Détails : [`docs/GIT_STRATEGY.md`](./docs/GIT_STRATEGY.md).

## Conventions de commit

Format **[Conventional Commits](https://www.conventionalcommits.org/fr/v1.0.0/)** :

```
<type>(<scope>): <description courte impérative>

<corps optionnel — pourquoi, contexte>

<footer optionnel — references issues, breaking changes>
```

| Type | Quand l'utiliser |
|---|---|
| `feat` | Nouvelle fonctionnalité |
| `fix` | Correction de bug |
| `docs` | Documentation uniquement |
| `style` | Formatage, point-virgules manquants… (pas de logique) |
| `refactor` | Refactoring sans nouvelle fonctionnalité ni fix |
| `perf` | Amélioration de performance |
| `test` | Ajout / modification de tests |
| `build` | Système de build, dépendances |
| `ci` | Pipeline CI/CD |
| `chore` | Maintenance diverse |

Exemples :

```
feat(streaming): implement pub/sub hub with backpressure

Adds StreamHub with goroutine-per-listener model. Buffers of 256 chunks
prevent slow listeners from stalling the broadcaster.

Closes #42
```

```
fix(auth): reject JWTs signed with weak HMAC key
```

## Stratégie de branches

| Branche | Rôle |
|---|---|
| `main` | Code en production. Protégée. |
| `develop` | Intégration des features avant release. |
| `feature/<slug>` | Nouvelle fonctionnalité. |
| `fix/<slug>` | Bug non bloquant. |
| `hotfix/<slug>` | Correctif urgent direct sur `main`. |
| `release/x.y.z` | Préparation d'une release. |

Voir [`docs/GIT_STRATEGY.md`](./docs/GIT_STRATEGY.md).

## Pull Requests

### Template

Chaque PR doit contenir :

- **Contexte** : pourquoi ce changement ?
- **Changements** : liste des modifications principales.
- **Tests** : comment vérifier ? (commandes, scénarios)
- **Liens** : issues, ADR, tickets RNCP couverts.
- **Screenshots** : pour tout changement UI.

### Checklist

- [ ] Le code compile (`go build ./...` / `flutter analyze`).
- [ ] Les tests passent (`go test ./... -race -cover` / `flutter test`).
- [ ] La couverture reste **≥ 80 %** côté backend.
- [ ] La doc associée est à jour (README, ADR, CHANGELOG).
- [ ] Aucun secret dans le diff (`gitleaks` doit passer).
- [ ] Commits signés GPG.
- [ ] CI verte.

## Code Review

- Au moins **1 review approuvée** requise avant merge sur `develop`.
- Au moins **2 reviews** pour un merge sur `main`.
- Reviewer : challenger la conception, vérifier les tests, suggérer concrètement.
- Auteur : répondre à chaque commentaire (résoudre ou expliquer).

## Tests

Toute nouvelle fonctionnalité doit s'accompagner de :

- Tests unitaires sur les usecases / blocs.
- Tests d'intégration pour les chemins critiques.
- Mise à jour du [cahier de recette](./docs/PLAN_DE_TESTS.md) si nouvelle US.

## Documentation

Une PR qui change le comportement public doit mettre à jour :

- Le `README.md` si nécessaire.
- Le `CHANGELOG.md` (section `[Unreleased]`).
- Le cahier des charges FR + EN si une spec change.
- La doc Swagger (annotations sur les handlers).
- Un nouvel ADR si décision structurante.

## Signaler un bug

Ouvrir une issue avec :

1. **Version** concernée (`git rev-parse HEAD`).
2. **Étapes** pour reproduire.
3. **Comportement attendu** vs **observé**.
4. **Logs / captures**.
5. **Environnement** (OS, navigateur, device).

Pour un problème de sécurité, **ne pas ouvrir d'issue publique** : contacter directement l'équipe via `contact@ecole-decode.fr` ou la procédure documentée dans [`docs/RGPD.md`](./docs/RGPD.md).

---

Merci pour ta contribution !
