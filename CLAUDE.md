# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**StreamPulse** — Real-time audio streaming platform. Semester project (5A Tech Lead, Bloc 3 RNCP 38822) focused on resilience, observability, and industrialization.

**Feature principale : le streaming audio en live** (diffusion temps reel d'un broadcaster vers N listeners via goroutines + channels). Tout le reste (auth, playlists, admin) gravite autour de cette feature.

- **Backend**: Go (Gin, GORM, PostgreSQL)
- **Frontend**: Flutter mobile (BLoC pattern, go_router, just_audio)
- **Observability**: OpenTelemetry traces, Prometheus metrics, Grafana dashboards, structured JSON logs (slog)
- **Infra**: Docker multi-stage, docker-compose (Postgres, OTel Collector, Prometheus, Grafana, Loki)

## Architecture

### Backend (Clean Architecture / DDD)

```
backend/
├── cmd/api/              # Entry point, graceful shutdown
├── internal/
│   ├── domain/
│   │   ├── entity/       # User, Stream, Playlist, Track
│   │   ├── repository/   # Interfaces only (no implementations)
│   │   └── service/      # Domain service interfaces
│   ├── application/
│   │   ├── usecase/      # Business logic (depends on domain interfaces)
│   │   └── dto/          # Request/response DTOs
│   ├── infrastructure/
│   │   ├── persistence/  # GORM repository implementations
│   │   ├── streaming/    # Audio stream multiplexing (goroutines + channels)
│   │   ├── auth/         # JWT + bcrypt
│   │   ├── config/       # Viper-based, 12-Factor (env vars only)
│   │   └── observability/# OTEL tracer, slog JSON logger, Prometheus metrics
│   └── transport/
│       └── http/
│           ├── handler/  # Gin handlers
│           ├── middleware/# Auth, CORS, logging, metrics
│           └── router/   # Route definitions
├── migrations/
├── deployments/
└── Dockerfile            # Multi-stage (golang:alpine -> alpine distroless)
```

**Dependency rule**: domain -> 0 deps | application -> domain | infrastructure -> domain | transport -> application + infrastructure

### Frontend (Feature-based + BLoC)

```
frontend/lib/
├── core/
│   ├── api/          # Dio client with JWT interceptor
│   ├── router/       # go_router config
│   ├── theme/        # Material 3 theme
│   └── storage/      # flutter_secure_storage wrapper
└── features/
    ├── auth/         # bloc/ repositories/ screens/ models/
    ├── player/       # Audio player (just_audio + audio_service)
    ├── streams/      # Live streams listing/listening
    ├── playlists/    # CRUD playlists
    ├── broadcaster/  # Stream creation/publishing
    └── admin/        # User management, global stats
```

### Roles

| Role | Capabilities |
|------|-------------|
| Anonymous | Browse limited content |
| User | Listen, favorites, playlists |
| Broadcaster | Create live streams, upload audio |
| Admin | Manage users, access global metrics |

## Common Commands

### Backend
```bash
cd backend
go build ./...                           # Build all
go test ./... -cover                     # Tests with coverage
go test ./... -race                      # Race condition detection
go test -run TestFunctionName ./path/... # Single test
go vet ./...                             # Static analysis
```

### Frontend
```bash
cd frontend
flutter pub get                    # Install deps
flutter run                        # Run on connected device
flutter test                       # All tests
flutter test test/path_test.dart   # Single test
flutter analyze                    # Lint
dart run build_runner build        # Code generation (json_serializable)
flutter build apk                  # Build Android APK
flutter build appbundle            # Build Android App Bundle (recommended for PlayStore)
flutter build ios                  # Build iOS (Mac required)
flutter build web --release        # Build web release
flutter test integration_test/     # Integration tests
```

### Infrastructure
```bash
docker compose up -d               # Start all services
docker compose up -d postgres      # Postgres only
docker compose logs -f api         # Follow API logs
docker compose down -v             # Tear down with volumes
```

### Endpoints
- API: http://localhost:8080
- Health: GET /health
- Metrics: GET /metrics (Prometheus)
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090

## Go Conventions (per course)

- **Error handling**: Always check `err != nil`. Wrap errors with `fmt.Errorf("context: %w", err)`. Never use panic for control flow.
- **Config**: Zero hardcoding. All config via env vars (Viper). `.env` for local only.
- **Logging**: `log/slog` with JSON handler. Structured key-value pairs. Never log secrets.
- **Testing**: Table-driven tests with `t.Run()`. Use testify assertions. Target 80% coverage. Use `mockgen` for interface mocks.
- **Concurrency**: goroutines + channels for streaming multiplexing. Always use `context.Context` for timeouts/cancellation. `sync.Mutex` with `defer Unlock()`. Run `-race` in CI.
- **Naming**: Exported = UpperCase, private = lowerCase. Short package names. Interfaces for all repository/service contracts.
- **Security**: JWT auth, bcrypt passwords, parameterized SQL queries, TLS in prod, CORS headers.
- **Formatting**: `go fmt` always. `go vet` in CI. Consider `golangci-lint`.

## Flutter Conventions (per course — Thomas Coichot)

- **State management**: BLoC pattern (flutter_bloc). Events = user actions, States = UI state. The course also covers Provider (`context.watch<T>()` for reactive rebuild, `context.read<T>()` for one-shot access).
- **Architecture**: Feature-based folders. Each feature has bloc/, models/, repositories/, screens/. `lib/` is the main code directory, `main.dart` is the entry point.
- **Theme FIRST**: Configurer le ThemeData complet (couleurs, typographies, styles de boutons, cards, inputs) AVANT de commencer les views. Utiliser `Theme.of(context)` partout dans les widgets.
- **Widgets**: Use `StatelessWidget` for static content, `StatefulWidget` when state changes. Always `.dispose()` controllers in `dispose()` to avoid memory leaks.
- **Lifecycle**: Know `initState()` -> `build()` -> `didUpdateWidget()` -> `deactivate()` -> `dispose()`.
- **Null safety**: Dart sound null safety — variables cannot be null unless explicitly declared with `?`. Avoid `!` operator unless certain.
- **Navigation**: GoRouter (`context.go()` replaces current route, `context.push()` adds to stack, `context.pop()` goes back). Use named routes and path parameters (`:id`).
- **Forms**: Use `GlobalKey<FormState>` + `TextEditingController`. Validate with `_formKey.currentState!.validate()`. Always dispose controllers.
- **Responsive**: Use `MediaQuery.sizeOf(context)` (not the old `.of(context).size`) for performance.
- **Networking**: Dio with interceptors for JWT. API URL via `--dart-define`. The course mentions the `http` package for simple calls but Dio is preferred for interceptors.
- **Storage**: `flutter_secure_storage` for tokens/secrets. `SharedPreferences` for non-sensitive key-value data (persists across restarts).
- **Testing**: 3 levels per course:
  - Unit tests (`test/` directory, mirrors `lib/` structure, matchers: `equals`, `isNull`, `throwsA`)
  - Widget tests (`testWidgets`, `tester.pumpWidget()`, matchers: `findsOneWidget`, `findsNothing`)
  - Integration tests (`integration_test/` directory, `flutter test integration_test/`)
- **Build/Deploy**: `flutter build apk`, `flutter build appbundle` (recommended for PlayStore), `flutter build ios` (Mac required), `flutter build web --release`.
- **Background audio**: audio_service for background playback continuity.
- **Useful packages** (per course): `intl` (i18n/dates), `file_picker`, `sentry_flutter` (crash monitoring), `url_launcher`, `toastification`, `collection`.

## Observability Stack (critical for Bloc 3)

- **Traces**: OpenTelemetry SDK in Go -> OTLP gRPC -> OTel Collector -> backend
- **Metrics**: Prometheus client in Go -> /metrics endpoint -> Prometheus scrapes -> Grafana
- **Logs**: slog JSON -> stdout -> Loki -> Grafana
- **Business metrics**: active streams, active listeners per stream, disconnection count (separate from HTTP 5xx)
- **Dashboard must show**: users online, streaming throughput, error rates, API response times
