# StreamPulse 🎧

Real-time audio streaming platform — **Projet Semestriel 5A Tech Lead, Bloc 3 RNCP 38822**
(résilience, observabilité, industrialisation).

A broadcaster publishes a live audio flow that is multiplexed in real time to N
listeners using **Go goroutines + channels**, with context-based cancellation to
avoid leaks. Everything else (auth, roles, playlists, admin, observability)
orbits this core feature.

## Stack

| Layer | Tech |
|-------|------|
| Backend | Go (Gin, GORM, PostgreSQL), Clean Architecture / DDD |
| Frontend | Flutter mobile (BLoC, go_router, just_audio) |
| Observability | OpenTelemetry traces, Prometheus metrics, Grafana dashboards, structured JSON logs (slog), Loki |
| Infra | Docker multi-stage, docker-compose, 12-Factor config |

## Quickstart (full stack)

```bash
docker compose up -d --build        # API + Postgres + OTel + Prometheus + Grafana + Loki
curl localhost:8080/health          # {"status":"ok"}
```

- API: http://localhost:8080
- Metrics: http://localhost:8080/metrics
- Grafana: http://localhost:3000 (admin / admin) — dashboard *StreamPulse — Bloc 3 Observability* auto-provisioned
- Prometheus: http://localhost:9090

An **admin** account is seeded from `ADMIN_EMAIL` / `ADMIN_PASSWORD`
(default `admin@streampulse.local` / `admin1234`).

## Run the backend locally (without Docker)

```bash
docker compose up -d postgres
cd backend
cp .env.example .env                # then set JWT_SECRET, ADMIN_EMAIL, ADMIN_PASSWORD
go run ./cmd/api
```

## Tests

```bash
cd backend
go test ./... -race -cover          # unit tests (core: auth ~85%, streaming ~88%)
go vet ./...
```

## End-to-end demo (real-time fan-out)

```bash
API=http://localhost:8080/api/v1

# 1. login as admin (a broadcaster-capable role)
TOKEN=$(curl -s -X POST $API/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@streampulse.local","password":"admin1234"}' | jq -r .token)

# 2. start a live stream
SID=$(curl -s -X POST $API/streams -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Radio Jazz","description":"live"}' | jq -r .id)

# 3. listen (in another terminal) — streams audio bytes as they arrive
curl -N $API/streams/$SID/listen -H "Authorization: Bearer $TOKEN"

# 4. publish audio (any binary source / file)
cat song.mp3 | curl -X POST --data-binary @- \
  $API/streams/$SID/publish -H "Authorization: Bearer $TOKEN"
```

## API overview

| Method | Path | Role | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/register` | public | sign up, returns JWT |
| POST | `/api/v1/auth/login` | public | login, returns JWT |
| GET | `/api/v1/users/me` | user | current profile |
| PUT | `/api/v1/users/me` | user | update profile |
| GET | `/api/v1/streams` | user | list live streams |
| GET | `/api/v1/streams/:id` | user | stream detail |
| GET | `/api/v1/streams/:id/listen` | user | **listen (audio fan-out)** |
| POST | `/api/v1/streams` | broadcaster | **start a live stream** |
| POST | `/api/v1/streams/:id/publish` | broadcaster | **push audio** |
| POST | `/api/v1/streams/:id/stop` | broadcaster | stop stream |
| DELETE | `/api/v1/streams/:id` | broadcaster/admin | delete stream |
| CRUD | `/api/v1/playlists` (+ `/tracks`) | user | playlists & queue |
| CRUD | `/api/v1/tracks` | broadcaster (write) | audio sources |
| GET | `/api/v1/admin/users` | admin | list users |
| PUT | `/api/v1/admin/users/:id/role` | admin | change role |
| GET | `/api/v1/admin/stats` | admin | global stats |

## Roles

`anonymous` → `user` (listen, favorites, playlists) → `broadcaster` (create live
streams, upload sources) → `admin` (manage users, global metrics). Promote a user
to broadcaster via `PUT /api/v1/admin/users/:id/role`.

## Architecture

See [`CLAUDE.md`](./CLAUDE.md) for the full Clean Architecture / DDD layout and
the [`docs/`](./docs) folder for the cahier des charges (FR/EN) and tickets.
