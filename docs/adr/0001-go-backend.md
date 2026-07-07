# ADR 0001 — Go pour le backend

| Statut | Date | Auteurs |
|---|---|---|
| Adopté | 2026-01-15 | Équipe StreamPulse |

## Contexte

StreamPulse doit pouvoir gérer **N auditeurs simultanés** sur un même flux audio en temps réel. Le sujet RNCP Bloc 3 exige notamment :

- une utilisation experte de la **concurrence** et du **parallélisme** ;
- une consommation mémoire minimale pour `N` listeners ;
- une architecture testable à ≥ 80 % de couverture ;
- une gestion fine du cycle de vie via `context.Context` (timeouts, annulation).

Le langage du backend doit donc être taillé pour :

1. la **concurrence native** légère (modèle CSP, goroutines) ;
2. la **performance brute** sur de la diffusion I/O ;
3. la **simplicité opérationnelle** d'un binaire statique à déployer.

## Décision

Le backend de StreamPulse est implémenté en **Go (≥ 1.22)**.

## Alternatives considérées

| Langage | Forces | Faiblesses pour ce projet | Verdict |
|---|---|---|---|
| **Go** | Goroutines/channels, binaire statique, écosystème SRE (Prometheus, OTel), GC bas-latence | Génériques moins matures, ergonomie ORM perfectible | ✅ Retenu |
| Node.js (TypeScript) | Écosystème, partage de types | Single-thread event-loop, mémoire élevée pour N connexions, GC moins prévisible | ❌ |
| Rust | Performance maximale, sécurité mémoire | Courbe d'apprentissage, vitesse de développement plus lente sur cadre académique | ❌ |
| Java / Spring | Maturité, écosystème | Empreinte mémoire, démarrage lent, complexité pour la concurrence audio | ❌ |
| Elixir / Erlang | Concurrence excellente, modèle acteur | Écosystème observabilité moins riche, équipe non formée | ❌ |

## Conséquences

### Positives

- Modèle **CSP** idéal pour le streaming pub/sub (voir ADR 0008).
- Binaire statique → image Docker **< 30 Mo** (Alpine), surface d'attaque réduite.
- Écosystème **SRE-native** : OpenTelemetry, Prometheus, gRPC.
- Démarrage en quelques millisecondes → adapté au scale horizontal.
- Cours dédié (Thomas Guillier) → alignement pédagogique.

### Négatives

- Ergonomie ORM (GORM) plus verbeuse qu'en JS/TS ou Python.
- Pas de runtime hot-reload natif (nécessite `air` ou équivalent).
- Génériques arrivés tardivement (1.18) → certaines libs récentes seulement.

### Neutres

- Gestion d'erreurs explicite par `if err != nil` — discipline d'équipe à maintenir.
- Modules standardisés via `go mod`, pas de package manager tiers.

## Liens

- [Sujet du projet — Bloc 3](../../Projet%20Semestriel%20Bloc%203%20%281%29.pdf)
- [ADR 0008 — Streaming pub/sub](./0008-streaming-pubsub.md)
- [The Twelve-Factor App](https://12factor.net)
