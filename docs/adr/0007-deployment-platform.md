# ADR 0007 — Plateforme de déploiement

| Statut | Date | Auteurs |
|---|---|---|
| Adopté | 2026-02-05 | Équipe StreamPulse |

## Contexte

L'application doit être **publiquement accessible** au jury en **HTTPS** :

- mise en production automatisée depuis la CI ;
- certificat TLS valide (Let's Encrypt) ;
- coût compatible avec un budget étudiant (≤ ~ 10 €/mois) ;
- infrastructure capable d'exposer plusieurs conteneurs (API + Postgres + observabilité minimum).

## Décision

Le déploiement production cible **Fly.io** (région `cdg`), avec PostgreSQL géré via Fly Postgres ou Neon.

L'environnement **staging** utilise la même plateforme pour parité.

Une alternative VPS chez OVH/Scaleway est documentée en plan B (voir conséquences).

## Alternatives considérées

| Plateforme | Forces | Faiblesses | Verdict |
|---|---|---|---|
| **Fly.io** | Déploiement Docker simple (`fly launch`), TLS auto, multi-régions, Free tier généreux, IPv6, healthchecks | Coûts au-delà du free-tier, persistance volumes spécifique | ✅ Retenu |
| VPS classique (OVH/Hetzner) | Contrôle total, prix fixe | Configuration manuelle (TLS, monitoring, sécu), pas de scale auto | 🟡 Plan B |
| AWS / GCP | Maturité maximale | Complexité IAM, factuation surprenante en académique | ❌ |
| Render / Railway | Très simple | Pricing moins prévisible, dépendance forte au PaaS | ❌ |
| Kubernetes managed (GKE/EKS) | Industriel | Surdimensionné, coûts > budget | ❌ (bonus TICK-102) |

## Conséquences

### Positives

- `fly deploy` en CI couvre Ce3.4.1 (CD automatisé) et Ce3.4.4 (mises à jour fréquentes).
- TLS Let's Encrypt automatique → score SSL Labs ≥ A.
- Healthchecks Fly.io connectés à `/health`.
- Conteneur unique pour l'API → image Docker multi-stage Alpine (< 30 Mo).

### Négatives

- L'observabilité (Grafana, Loki, Tempo) reste **locale** ou hébergée sur un second conteneur Fly → coût supplémentaire.
- Volumes persistants Fly.io plus chers que stockage S3 → uploads tracks redirigés vers MinIO/S3 si volume.

### Neutres

- Configuration : `fly.toml` versionné dans le repo.
- Secrets : `fly secrets set DATABASE_URL=...` (jamais commité).
- Le passage à Kubernetes est possible plus tard (TICK-102 bonus).

## Procédure de déploiement

1. CI build l'image Docker, push sur GHCR.
2. CI exécute `flyctl deploy --image ghcr.io/<org>/streampulse-api:<tag>`.
3. Healthcheck `/health` valide la mise en service.
4. Si KO → rollback automatique vers l'image précédente.

## Liens

- [Fly.io Docs](https://fly.io/docs/)
- [`flyctl`](https://fly.io/docs/flyctl/)
- [ADR 0006 — Observabilité](./0006-observability-stack.md)
- TICK-005 — Mise en production
