# ADR 0006 — Stack d'observabilité OTEL + Prometheus + Grafana + Loki

| Statut | Date | Auteurs |
|---|---|---|
| Adopté | 2026-01-25 | Équipe StreamPulse |

## Contexte

Le **cœur du Bloc 3 RNCP** est l'observabilité. Les exigences :

- **Traces distribuées** : suivre une requête mobile → API → DB.
- **Logs JSON structurés** indexables (Loki, Elasticsearch).
- **Métriques** différenciant **technique** (HTTP 5xx) et **métier** (déconnexions brutales, listeners actifs).
- **Dashboard Grafana** : utilisateurs en ligne, débit streaming, taux d'erreurs, temps de réponse.
- **Alertes** automatiques (Ce3.5.2).

## Décision

Adoption d'une stack **CNCF-native** :

| Signal | Outil |
|---|---|
| **Traces** | OpenTelemetry Go SDK → OTel Collector → Tempo |
| **Métriques** | `prometheus/client_golang` (scrape direct) + métriques OTEL via Collector |
| **Logs** | `log/slog` JSONHandler → Loki (via Promtail) |
| **Visualisation** | Grafana |

## Alternatives considérées

| Stack | Verdict |
|---|---|
| **OTEL + Prometheus + Grafana + Loki** | ✅ Standard CNCF, gratuit, parfaite démonstration RNCP |
| Datadog | ❌ Payant, dépendance SaaS, pas pédagogique |
| New Relic | ❌ Idem |
| ELK (Elasticsearch + Kibana) | ❌ Plus gourmand, configuration plus lourde, pas de tracing natif |
| Stack maison | ❌ Réinventer la roue, contraire au critère RNCP « standard OpenTelemetry » |

## Conséquences

### Positives

- **OTEL** est le standard CNCF → instrumentation portable.
- **Corrélation logs ↔ traces** via `trace_id` injecté dans chaque log slog.
- **Grafana** unifie la visualisation des trois signaux.
- **Métriques métier** clairement séparées :
  - `streampulse_active_streams` (gauge)
  - `streampulse_active_listeners` (gauge)
  - `streampulse_stream_disconnections_total` (counter)
  - `http_requests_total` (counter, technique)
  - `http_request_duration_seconds` (histogram, technique)

### Négatives

- Stack complète = 5 conteneurs supplémentaires en dev → coût RAM local.
- Apprentissage OTel parfois confus (SDK / API / Collector).
- Configuration du Collector (`deployments/otel-collector.yaml`) à maintenir.

### Neutres

- Possibilité future de remplacer Tempo par Jaeger sans changer le code applicatif.
- Possibilité de remplacer Loki par Elasticsearch si volumétrie élevée (improbable en académique).

## Métriques métier exposées (extrait)

| Métrique | Type | Labels | Usage |
|---|---|---|---|
| `streampulse_active_streams` | gauge | — | Dashboard temps réel |
| `streampulse_active_listeners` | gauge | `stream_id` | Audience par flux |
| `streampulse_stream_disconnections_total` | counter | `reason` (`client_close`/`timeout`/`error`) | Identifier les déconnexions brutales (Bloc 3 exige cette distinction) |
| `streampulse_stream_bytes_total` | counter | `stream_id` | Volume diffusé |
| `http_requests_total` | counter | `method`, `path`, `status` | Erreurs techniques |
| `http_request_duration_seconds` | histogram | `method`, `path` | Latence p50/p95/p99 |

## Liens

- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Prometheus client_golang](https://github.com/prometheus/client_golang)
- [Grafana Loki](https://grafana.com/oss/loki/)
- [ADR 0001 — Go backend](./0001-go-backend.md)
