# ADR 0008 — Architecture streaming en pub/sub goroutines

| Statut | Date | Auteurs |
|---|---|---|
| Adopté | 2026-02-10 | Équipe StreamPulse |

## Contexte

Le **cœur fonctionnel** de StreamPulse est la diffusion d'un flux audio d'un broadcaster vers **N listeners**. Les contraintes :

- temps réel (latence < 2 s end-to-end) ;
- protection contre les **memory leaks** (cleanup `context`) ;
- **race-safe** (validation par `go test -race`) ;
- gestion de la **backpressure** : un listener lent ne doit pas bloquer le broadcaster ;
- métriques business par stream (listeners, déconnexions).

## Décision

Implémentation d'un **StreamHub** centralisé, un par stream actif, basé sur le **pattern pub/sub** :

- une **goroutine producteur** lit les chunks audio publiés par le broadcaster ;
- chaque **listener** dispose d'une **goroutine consommateur** + un **channel bufferisé** (256 chunks) ;
- la **backpressure** est gérée par `select { case ch <- chunk: default: drop & metric }`. Un listener qui ne suit pas voit ses chunks droppés sans bloquer les autres ;
- le **context.Context** du broadcaster propage l'annulation à tous les listeners → cleanup automatique ;
- un `sync.RWMutex` protège la map des subscribers.

```go
type StreamHub struct {
    mu          sync.RWMutex
    subscribers map[string]chan []byte
    ctx         context.Context
    cancel      context.CancelFunc
}

func (h *StreamHub) Subscribe(id string) <-chan []byte { ... }
func (h *StreamHub) Publish(chunk []byte) { ... }
func (h *StreamHub) Unsubscribe(id string) { ... }
func (h *StreamHub) Close() { ... } // cancel + close all channels
```

## Alternatives considérées

| Approche | Forces | Faiblesses | Verdict |
|---|---|---|---|
| **Hub goroutines + channels bufferisés** | Idiomatique Go, race-safe, simple à tester | Bufferisation à dimensionner | ✅ Retenu |
| Redis Pub/Sub | Permet plusieurs instances API | Latence supplémentaire, dépendance externe, surdimensionné pour MVP | ❌ (envisageable post-MVP) |
| WebSocket par listener | Bi-directionnel | Plus complexe, pas besoin de bi-directionnel pour de l'écoute | ❌ |
| WebRTC SFU | Latence ultra basse | Complexité énorme, hors périmètre académique | ❌ |
| HLS / DASH (segments) | Standard streaming web | Latence ≥ 6-10 s, pas du temps réel | ❌ |
| HTTP/2 server push | Standard, multiplexé | Support inégal, peu d'écosystème streaming audio | ❌ |

## Conséquences

### Positives

- **Memory leak prévenu** : `context.Done()` ferme channels et termine goroutines.
- **Linear scaling** sur une instance jusqu'à ~ 1 000 listeners/stream sur un VPS 2 vCPU.
- **Tests `-race`** garantissent l'absence de data races.
- **Métriques business** émises directement depuis le hub :
  - `streampulse_active_listeners{stream_id=...}`
  - `streampulse_stream_disconnections_total{reason="slow_consumer"}` quand un listener voit ses chunks droppés.
- Backpressure **mesurable** plutôt que cachée.

### Négatives

- **Limité à une seule instance API** pour un même stream. Pour plusieurs instances → externalisation pub/sub (Redis Streams, NATS) ; à envisager en v2.
- Si tous les listeners sont lents en même temps, l'effet "snowball" est possible → mitigation par un *drain* asynchrone par listener (déjà prévu).

### Neutres

- Le choix d'un buffer de **256 chunks** est arbitraire. À tuner lors des tests de charge (TICK-074).
- Le **transcodage** à la volée (bonus TICK-104) s'insérerait entre la goroutine producteur et la map des subscribers.

## Tests associés

- **Unitaires** : table-driven `TestStreamHub_*` avec `-race`.
- **Benchmark** : `BenchmarkStreamHub_NSubscribers` pour 10 / 100 / 1000 listeners.
- **Charge** : `k6` avec 100 listeners simultanés (cible MVP).

## Estimation de coût CPU/RAM (à confirmer après TICK-074)

| Flux × listeners | CPU (vCPU) | RAM | Bande passante | Coût VPS estimé |
|---|---|---|---|---|
| 1 × 10 | ~ 1 % | ~ 30 Mo | ~ 1.3 Mbps | ~ 5 €/mois |
| 1 × 100 | ~ 5 % | ~ 50 Mo | ~ 13 Mbps | ~ 5 €/mois |
| 10 × 100 | ~ 40 % | ~ 250 Mo | ~ 130 Mbps | ~ 15 €/mois |

> Ces chiffres seront actualisés à partir des résultats de **TICK-074 — Tests de performance / charge**.

## Liens

- [ADR 0001 — Go backend](./0001-go-backend.md)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- TICK-020 — Architecture du streaming multiplexé
- TICK-074 — Tests de performance / charge
