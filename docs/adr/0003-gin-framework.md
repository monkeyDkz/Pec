# ADR 0003 — Gin comme framework HTTP

| Statut | Date | Auteurs |
|---|---|---|
| Adopté | 2026-01-20 | Équipe StreamPulse |

## Contexte

Le backend Go doit exposer une **API REST** et gérer du **streaming HTTP chunked** vers les listeners. Il faut un framework qui propose :

- un routeur performant,
- un système de **middlewares** (auth, métriques, logs, CORS),
- un **écosystème mature** d'extensions (Swagger, OTel),
- une intégration aisée avec le streaming `http.ResponseWriter`.

## Décision

L'API utilise **Gin (v1.10+)** comme framework HTTP.

## Alternatives considérées

| Framework | Forces | Faiblesses | Verdict |
|---|---|---|---|
| **Gin** | Routeur très rapide (httprouter), middlewares simples, intégration OTel (`otelgin`), Swag, popularité massive | Validation moins ergonomique qu'Echo, certains paterns "magiques" | ✅ Retenu |
| Echo | Très propre, validation intégrée, performance comparable | Communauté plus petite, moins de docs/middlewares | ❌ |
| `net/http` + chi | Idiomatique, minimal | Plus de code à écrire (middlewares, JSON, validation) | ❌ |
| Fiber | Très rapide (fasthttp) | Incompatibilité partielle avec `net/http` standard → casse OTel/Prometheus tooling | ❌ |
| go-kit | Excellent pour microservices | Trop lourd pour un monolithe pédagogique | ❌ |

## Conséquences

### Positives

- Démarrage rapide : `r := gin.New()` puis `r.GET(...)`.
- Middleware d'authentification JWT trivial à brancher.
- `c.Stream(...)` adapté au streaming audio chunked.
- Intégration native avec `otelgin` (traces) et un middleware Prometheus maison.

### Négatives

- Pas de validation intégrée → utilisation de `go-playground/validator` ou validation manuelle dans les usecases.
- Le contexte Gin (`*gin.Context`) doit être converti en `context.Context` pour les opérations longues (`c.Request.Context()`).

### Neutres

- Risque de "fat handlers" → atténué par la Clean Architecture (handlers ⇒ usecases).
- Stack trace par défaut peu lisible → middleware `gin.Recovery()` configuré pour logger via slog.

## Liens

- [Gin Web Framework](https://gin-gonic.com/)
- [`otelgin`](https://github.com/open-telemetry/opentelemetry-go-contrib)
- [`swaggo/swag`](https://github.com/swaggo/swag)
