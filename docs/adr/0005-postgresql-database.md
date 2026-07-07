# ADR 0005 — PostgreSQL comme base relationnelle

| Statut | Date | Auteurs |
|---|---|---|
| Adopté | 2026-01-22 | Équipe StreamPulse |

## Contexte

Les données métier (utilisateurs, streams, playlists, tracks, association many-to-many playlist↔track, feedbacks) sont :

- **structurées** et fortement relationnelles ;
- soumises à des **contraintes d'intégrité** (cascade de suppression RGPD, unicité email) ;
- consultées par des **agrégations** simples (compter les streams live, lister par owner, etc.).

Le projet doit aussi prouver la **maîtrise SQL** dans un contexte Tech Lead.

## Décision

La base de données est **PostgreSQL 16**, accédée via **GORM** côté Go.

## Alternatives considérées

| Base | Forces | Faiblesses | Verdict |
|---|---|---|---|
| **PostgreSQL** | ACID, JSONB, écosystème mature, support natif `gen_random_uuid()`, performance, observabilité | Hébergement plus lourd que SQLite en dev | ✅ Retenu |
| MySQL / MariaDB | Très répandu, ACID | Moins de fonctionnalités modernes (JSONB plus pauvre, CTE longtemps absentes) | ❌ |
| SQLite | Embarqué, simple en dev | Non adapté à plusieurs replicas / concurrence write élevée en prod | ❌ |
| MongoDB | Flexible, JSON natif | Mauvais fit pour relations (`playlist_tracks`, `broadcaster_id`), pas d'ACID multi-document avant 4.x | ❌ |
| CockroachDB / Spanner | Distribué, ACID | Surdimensionné pour un MVP académique, coût | ❌ |

## Conséquences

### Positives

- **Contraintes** d'intégrité fortes : `FOREIGN KEY`, `UNIQUE`, `ON DELETE CASCADE`.
- **JSONB** disponible pour les exports RGPD à la volée.
- **Healthcheck** trivial dans docker-compose.
- Support natif des extensions (`pgcrypto` pour les UUID, `pg_trgm` pour recherche, etc.).
- GORM offre une couche ergonomique tout en restant en SQL.

### Négatives

- Setup plus lourd qu'un SQLite (mais résolu par docker-compose).
- GORM ajoute une abstraction supplémentaire à comprendre.
- Quelques requêtes complexes (statistiques d'écoute) nécessiteront du SQL brut.

### Neutres

- Migrations versionnées avec `golang-migrate` à venir (pour aller au-delà d'`AutoMigrate`).
- Pour le développement local, une seule instance partagée via docker-compose.

## Liens

- [PostgreSQL 16](https://www.postgresql.org/docs/16/)
- [GORM](https://gorm.io)
- [`golang-migrate`](https://github.com/golang-migrate/migrate)
