# StreamPulse

> Plateforme de streaming audio temps réel — Projet semestriel 5A Tech Lead, S2 Bloc 3 (RNCP 38822) — École .decode.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](./backend/go.mod)
[![Flutter](https://img.shields.io/badge/Flutter-3.22+-02569B.svg)](./frontend/pubspec.yaml)
[![Docs](https://img.shields.io/badge/docs-FR%20%7C%20EN-success.svg)](./docs/cahier-des-charges-fr.md)

StreamPulse permet à un *broadcaster* de diffuser un flux audio en direct vers N *listeners* simultanés, via une API Go performante et une application mobile Flutter. Le projet illustre les compétences attendues d'un Tech Lead : résilience, observabilité, industrialisation et conformité.

---

## Sommaire

- [Aperçu](#aperçu)
- [Architecture](#architecture)
- [Prérequis](#prérequis)
- [Démarrage rapide](#démarrage-rapide)
- [Structure du projet](#structure-du-projet)
- [Configuration](#configuration)
- [Endpoints API](#endpoints-api)
- [Tests](#tests)
- [Observabilité](#observabilité)
- [Déploiement](#déploiement)
- [Documentation](#documentation)
- [Équipe et répartition](#équipe-et-répartition)
- [Licence](#licence)

---

## Aperçu

| Brique | Stack | Description |
|---|---|---|
| **Backend** | Go 1.22 · Gin · GORM · PostgreSQL 16 | API REST + moteur de streaming audio basé sur goroutines / channels |
| **Frontend** | Flutter 3.22 · BLoC · go_router · just_audio | Application mobile iOS / Android |
| **Observabilité** | OpenTelemetry · Prometheus · Grafana · Loki | Traces, métriques techniques et métier, logs JSON |
| **Infra** | Docker multi-stage · docker-compose · GitHub Actions | Conteneurisation, CI/CD, déploiement |

### Fonctionnalités principales

- 🔐 Authentification JWT (inscription, connexion, multi-rôles)
- 🎙️ Diffusion audio temps réel multi-clients
- 🎧 Lecteur audio mobile avec lecture en arrière-plan
- 📋 Gestion de playlists avec file d'attente
- 👥 4 rôles : `anonymous`, `user`, `broadcaster`, `admin`
- 📊 Dashboard Grafana avec métriques techniques **et** métier
- 🔍 Traces distribuées OTEL de l'app mobile jusqu'à la base
- 🛡️ Conformité RGPD (export, suppression, registre des traitements)

---

## Architecture

```mermaid
graph LR
    M[Mobile Flutter] -->|HTTPS / JWT| A[API Go - Gin]
    A -->|SQL| DB[(PostgreSQL)]
    A -->|OTLP gRPC| C[OTel Collector]
    A -->|/metrics| P[Prometheus]
    A -->|stdout JSON| L[Loki]
    C --> P
    G[Grafana] --> P
    G --> L
```

Détails complets : [`docs/cahier-des-charges-fr.md`](./docs/cahier-des-charges-fr.md).

---

## Prérequis

| Outil | Version minimale | Vérification |
|---|---|---|
| Go | 1.22 | `go version` |
| Flutter | 3.22 | `flutter --version` |
| Docker | 24+ | `docker --version` |
| Docker Compose | v2 | `docker compose version` |
| Make *(optionnel)* | — | `make --version` |
| Git | 2.40+ | `git --version` |

---

## Démarrage rapide

### 1. Cloner le dépôt

```bash
git clone https://github.com/<org>/streampulse.git
cd streampulse
```

### 2. Configurer les variables d'environnement

```bash
cp .env.example .env
# Éditer .env et renseigner au minimum :
# - DATABASE_URL
# - JWT_SECRET (générer avec : openssl rand -base64 48)
```

### 3. Lancer la stack complète

```bash
docker compose up -d
```

Les services exposés en local :

| Service | URL | Credentials |
|---|---|---|
| API Go | http://localhost:8080 | — |
| Healthcheck | http://localhost:8080/health | — |
| Métriques Prometheus | http://localhost:8080/metrics | — |
| Prometheus | http://localhost:9090 | — |
| Grafana | http://localhost:3000 | `admin / admin` |
| PostgreSQL | `localhost:5432` | voir `.env` |

### 4. Lancer l'application mobile

```bash
cd frontend
flutter pub get
flutter run --dart-define=API_URL=http://localhost:8080/api/v1
```

---

## Structure du projet

```
.
├── backend/                     # API Go (Clean Architecture / DDD)
│   ├── cmd/api/                 # Entry point + graceful shutdown
│   ├── internal/
│   │   ├── domain/              # Entities, repository interfaces
│   │   ├── application/         # UseCases + DTOs
│   │   ├── infrastructure/      # GORM, JWT, observability, streaming
│   │   └── transport/http/      # Gin handlers, middlewares, router
│   ├── migrations/              # Migrations SQL versionnées
│   └── Dockerfile               # Multi-stage Alpine
├── frontend/                    # App Flutter (Feature-based + BLoC)
│   └── lib/
│       ├── core/                # API client, theme, router, storage
│       └── features/            # auth, player, streams, playlists, broadcaster, admin
├── deployments/                 # Configs OTel Collector, Prometheus
├── docs/                        # Documentation projet
│   ├── cahier-des-charges-fr.md # Spécifications complètes FR
│   ├── cahier-des-charges-en.md # Spécifications complètes EN (B2)
│   ├── TICKETS.md               # Tickets RNCP-mappés
│   ├── adr/                     # Architecture Decision Records
│   ├── PLAN_DE_TESTS.md         # Plan de tests + cahier de recette
│   ├── PLAN_FORMATION-fr.md     # Plan de formation utilisateurs
│   ├── RGPD.md                  # Registre des traitements
│   ├── TOOLCHAIN.md             # Chaîne d'outils CI/CD
│   ├── GIT_STRATEGY.md          # Stratégie de branches
│   └── ACCESSIBILITY.md         # Politique d'accessibilité
├── docker-compose.yml
├── CHANGELOG.md
├── CONTRIBUTING.md
├── LICENSE
└── README.md
```

---

## Configuration

Toutes les variables sont injectées via environnement (12-Factor App). Voir [`.env.example`](./.env.example).

| Variable | Obligatoire | Défaut | Description |
|---|:---:|---|---|
| `PORT` | non | `8080` | Port d'écoute HTTP |
| `DATABASE_URL` | **oui** | — | DSN PostgreSQL (`postgres://user:pass@host:5432/db`) |
| `JWT_SECRET` | **oui** | — | Secret HMAC-SHA256 (≥ 32 caractères) |
| `OTEL_ENDPOINT` | non | `localhost:4317` | Endpoint OTLP gRPC |
| `LOG_LEVEL` | non | `info` | `debug`, `info`, `warn`, `error` |
| `ENVIRONMENT` | non | `development` | `development`, `staging`, `production` |
| `CORS_ORIGINS` | non | `*` | Origines CORS (CSV) |

---

## Endpoints API

Documentation interactive complète : [`/swagger/`](http://localhost:8080/swagger/) (générée via `swag`).

Résumé :

| Méthode | Chemin | Rôle | Description |
|---|---|---|---|
| POST | `/api/v1/auth/register` | Anonyme | Inscription |
| POST | `/api/v1/auth/login` | Anonyme | Connexion (JWT) |
| GET | `/api/v1/users/me` | User+ | Profil |
| DELETE | `/api/v1/users/me` | User+ | Suppression (RGPD) |
| GET | `/api/v1/users/me/data` | User+ | Export données (RGPD) |
| GET | `/api/v1/streams` | Anonyme | Liste streams |
| POST | `/api/v1/streams` | Broadcaster+ | Créer stream |
| GET | `/api/v1/streams/:id/listen` | Anonyme | Écouter (chunked HTTP) |
| POST | `/api/v1/streams/:id/publish` | Broadcaster+ | Diffuser audio |
| GET / POST | `/api/v1/playlists` | User+ | CRUD playlists |
| GET | `/api/v1/admin/users` | Admin | Liste utilisateurs |
| GET | `/health` | — | Healthcheck |
| GET | `/metrics` | Interne | Métriques Prometheus |

---

## Tests

### Backend

```bash
cd backend
go test ./... -race -cover            # Tests + race detector + coverage
go test ./... -coverprofile=cover.out # Génère le profil
go tool cover -html=cover.out         # Visualise la couverture
go test -bench=. ./internal/...       # Benchmarks
```

Cible : **≥ 80 % de couverture**.

### Frontend

```bash
cd frontend
flutter test                          # Unit + widget tests
flutter test integration_test/        # Tests d'intégration
flutter analyze                       # Lint statique
```

### Tests de charge

```bash
k6 run scripts/load-test.js           # 100 listeners simultanés
```

Détails : [`docs/PLAN_DE_TESTS.md`](./docs/PLAN_DE_TESTS.md).

---

## Observabilité

| Composant | URL | Usage |
|---|---|---|
| **Grafana** | http://localhost:3000 | Dashboards, alertes |
| **Prometheus** | http://localhost:9090 | Métriques techniques + métier |
| **Loki** | http://localhost:3100 | Logs JSON corrélés (`trace_id`) |
| **Tempo** | http://localhost:3200 | Traces distribuées |

Dashboards livrés (JSON provisionné dans `deployments/grafana/`) :

- **Technique** : taux d'erreur HTTP, latence p50/p95/p99, throughput.
- **Métier** : utilisateurs en ligne, streams actifs, listeners par stream, déconnexions brutales.

---

## Déploiement

Pipeline CI/CD : voir [`docs/TOOLCHAIN.md`](./docs/TOOLCHAIN.md).

Stratégie de branches : voir [`docs/GIT_STRATEGY.md`](./docs/GIT_STRATEGY.md).

Production : URL publique en HTTPS via Let's Encrypt (`https://streampulse.<domaine>`).

---

## Documentation

| Document | Description |
|---|---|
| [Cahier des charges FR](./docs/cahier-des-charges-fr.md) | Spécifications fonctionnelles + techniques |
| [Cahier des charges EN](./docs/cahier-des-charges-en.md) | English specifications (B2) |
| [Tickets RNCP](./docs/TICKETS.md) | Mapping critères RNCP → tickets |
| [ADRs](./docs/adr/) | Architecture Decision Records |
| [Plan de tests](./docs/PLAN_DE_TESTS.md) | Cahier de recette + plan iteratif |
| [Plan de formation FR](./docs/PLAN_FORMATION-fr.md) / [EN](./docs/PLAN_FORMATION-en.md) | Guide utilisateur par profil |
| [RGPD](./docs/RGPD.md) | Registre des traitements |
| [Toolchain](./docs/TOOLCHAIN.md) | Chaîne d'outils CI/CD |
| [Stratégie Git](./docs/GIT_STRATEGY.md) | Branches, commits, reviews |
| [Accessibilité](./docs/ACCESSIBILITY.md) | Politique WCAG / a11y |
| [Contribuer](./CONTRIBUTING.md) | Workflow contributeur |
| [Changelog](./CHANGELOG.md) | Historique des versions |

---

## Équipe et répartition

| Membre | Rôle principal | Périmètre |
|---|---|---|
| *(à compléter)* | Backend Lead | API Go, streaming hub, observabilité, CI |
| *(à compléter)* | Frontend Lead | Application Flutter, BLoC, lecteur audio |
| *(à compléter)* | DevOps / Infra | Docker, déploiement, monitoring |
| *(à compléter)* | QA / Documentation | Plan de tests, ADR, formation |

Détail des contributions : `git shortlog -sne` et page « Insights » du dépôt.

---

## Licence

Distribué sous licence MIT. Voir [`LICENSE`](./LICENSE).

---

## Contact

École .decode — `contact@ecole-decode.fr`
Encadrants : Thomas Guillier (Backend Go), Thomas Coichot (Flutter), Karl Marques Bernardo (Projet).
