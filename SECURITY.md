# Politique de sécurité — StreamPulse

## Versions supportées

| Version | Statut |
|---|---|
| 1.x | ✅ Activement maintenue |
| < 1.0 | ❌ Non supportée |

## Signaler une vulnérabilité

Merci de **ne pas ouvrir d'issue GitHub publique** pour un problème de sécurité.

Envoyer un email à **`contact@ecole-decode.fr`** avec :

- une description du problème,
- les étapes de reproduction,
- l'impact estimé,
- la version concernée.

Une réponse initiale est envoyée sous **72 heures**. Un correctif est publié dans les meilleurs délais selon la criticité (jusqu'à 30 jours pour les failles non critiques).

## Engagements

- Pas de poursuite contre les chercheurs en sécurité agissant de bonne foi.
- Reconnaissance publique (avec accord) dans le `CHANGELOG.md` lors du correctif.
- Suivi des bonnes pratiques OWASP et de la procédure RGPD ([`docs/RGPD.md`](./docs/RGPD.md) § 8).

## Mesures préventives intégrées

- Scan de dépendances : `govulncheck`, `trivy`, Dependabot (Ce3.3.4).
- Détection de secrets : `gitleaks` à chaque PR.
- Authentification : JWT signé, bcrypt (coût ≥ 12), rate limiting.
- TLS 1.2+ obligatoire en production.
- Headers de sécurité : `X-Content-Type-Options`, `X-Frame-Options`, `Content-Security-Policy`.
- Audit hebdomadaire automatisé.

## Hall of Fame

*(Aucune divulgation pour le moment.)*
