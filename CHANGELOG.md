# Changelog

Toutes les évolutions notables du projet sont documentées dans ce fichier.

Le format suit [Keep a Changelog 1.1.0](https://keepachangelog.com/fr/1.1.0/) et le versionnage suit [Semantic Versioning 2.0.0](https://semver.org/lang/fr/).

> **Convention** : `Added` / `Changed` / `Deprecated` / `Removed` / `Fixed` / `Security`.
> Chaque release est taguée Git (`vMAJOR.MINOR.PATCH`) et associée à une note de release GitHub.

---

## [Unreleased]

### Added
- Cahier des charges complet FR/EN (B2) avec diagrammes UML et BPMN.
- Set complet d'ADR couvrant les décisions structurantes (Go, Flutter, PostgreSQL, BLoC, observabilité, streaming, déploiement).
- Plan de tests itératif + cahier de recette mappé sur les user stories.
- Plan de formation utilisateurs FR/EN (auditeur, broadcaster, admin) incluant adaptations accessibilité.
- Registre RGPD et procédures d'export / suppression.
- Documentation de la chaîne d'outils CI/CD (toolchain).
- Stratégie Git documentée (GitFlow simplifié + commits signés).
- Politique d'accessibilité (WCAG 2.1 AA cible).

### Changed
- *(à compléter au fil du sprint)*

### Security
- *(à compléter au fil du sprint)*

---

## [1.0.0] — *Prévu*

Première release stable, présentée au jury RNCP.

### Added
- Authentification JWT (register / login).
- 4 rôles : `anonymous`, `user`, `broadcaster`, `admin`.
- API CRUD utilisateurs, streams, playlists, tracks.
- Moteur de streaming multiplexé (goroutines + channels).
- Application Flutter : login, listing streams, lecteur audio, broadcaster UI.
- Observabilité native : OTEL traces, Prometheus, Loki, dashboard Grafana.
- Pipeline CI/CD GitHub Actions (lint, test, build, scan, deploy).
- Image Docker multi-stage Alpine (< 30 Mo).
- Déploiement public HTTPS (Let's Encrypt).
- Endpoints RGPD (`/users/me/data`, `DELETE /users/me`).

### Security
- bcrypt (coût 12) pour les mots de passe.
- Rate limiting sur `/auth/*`.
- Scan dépendances (`govulncheck`, `trivy`, `gitleaks`) dans la CI.

---

## Politique de versionnage

| Type | Quand | Exemple |
|---|---|---|
| **MAJOR** | Changement incompatible de l'API publique | `1.x` → `2.0.0` |
| **MINOR** | Ajout rétrocompatible de fonctionnalité | `1.0` → `1.1.0` |
| **PATCH** | Correction de bug rétrocompatible | `1.0.0` → `1.0.1` |
| **Pre-release** | Suffixe `-rc.N` ou `-beta.N` | `1.0.0-rc.1` |

Chaque release :

1. Crée un tag Git signé : `git tag -s v1.0.0 -m "Release 1.0.0"`.
2. Pousse les tags : `git push origin --tags`.
3. Génère une note de release GitHub à partir des entrées de ce fichier.
4. Met à jour la documentation API (`/swagger/`) avec la nouvelle version.

---

[Unreleased]: https://github.com/<org>/streampulse/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/<org>/streampulse/releases/tag/v1.0.0
