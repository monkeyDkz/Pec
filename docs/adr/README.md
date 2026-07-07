# Architecture Decision Records (ADR)

Cette section regroupe les **décisions structurantes** prises sur le projet StreamPulse. Chaque décision est documentée selon le format proposé par Michael Nygard, en français.

## Pourquoi des ADRs ?

> *« On documente non seulement ce qu'on a fait, mais surtout pourquoi on a choisi une technologie plutôt qu'une autre. »* — Sujet RNCP Bloc 3.

Les ADRs permettent :

- de rendre les choix **traçables** (qui, quand, pourquoi) ;
- d'éviter de **rejouer** des débats déjà tranchés ;
- de faciliter l'**onboarding** de nouveaux contributeurs ;
- de **justifier** les choix techniques en soutenance.

## Format d'un ADR

```
# ADR XXXX — Titre court

| Statut | Date | Auteurs |

## Contexte
Pourquoi a-t-on dû trancher ?

## Décision
Ce que l'on a choisi (clair, affirmatif).

## Alternatives considérées
Ce que l'on a évalué et écarté, avec raisons.

## Conséquences
Positives, négatives, neutres.

## Liens
Tickets, issues, autres ADRs.
```

## Cycle de vie

Un ADR peut avoir l'un des statuts suivants :

- **Proposé** — discussion en cours.
- **Adopté** — la décision est en vigueur.
- **Déprécié** — toujours actif mais à éviter pour de nouveaux développements.
- **Remplacé** — annulé par un ADR plus récent (indiquer le numéro).

## Index

| # | Titre | Statut | Date |
|---|---|---|---|
| [0001](./0001-go-backend.md) | Go pour le backend | Adopté | 2026-01-15 |
| [0002](./0002-flutter-frontend.md) | Flutter pour le mobile multi-plateforme | Adopté | 2026-01-15 |
| [0003](./0003-gin-framework.md) | Gin comme framework HTTP | Adopté | 2026-01-20 |
| [0004](./0004-bloc-state-management.md) | BLoC pour la gestion d'état Flutter | Adopté | 2026-01-20 |
| [0005](./0005-postgresql-database.md) | PostgreSQL comme base relationnelle | Adopté | 2026-01-22 |
| [0006](./0006-observability-stack.md) | OpenTelemetry + Prometheus + Grafana + Loki | Adopté | 2026-01-25 |
| [0007](./0007-deployment-platform.md) | Plateforme de déploiement (Fly.io) | Adopté | 2026-02-05 |
| [0008](./0008-streaming-pubsub.md) | Architecture streaming en pub/sub goroutines | Adopté | 2026-02-10 |

## Outils

Pour générer un nouvel ADR :

```bash
# Numéro suivant
NEXT=$(printf "%04d" $(( $(ls docs/adr/[0-9]*.md | wc -l) + 1 )))
cp docs/adr/template.md docs/adr/${NEXT}-mon-titre.md
```
