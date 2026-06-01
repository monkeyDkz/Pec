# StreamPulse - Tickets Projet Complet

> Chaque ticket est taggue avec les criteres RNCP Bloc 3 qu'il couvre.
> Priorite: P0 = bloquant, P1 = essentiel, P2 = important, P3 = bonus

---

## EPIC 1 : Infrastructure & DevOps

### TICK-001 : Setup Docker Compose complet

**Priorite**: P0 | **RNCP**: A3.4, A3.5, C3.4.1, C3.5.1
**Assignable a**: Backend

**Description**: Configurer l'environnement de developpement local complet avec Docker Compose.

**Taches**:

- [x] Service PostgreSQL 16 avec healthcheck
- [x] Service API Go (build multi-stage)
- [x] Service OpenTelemetry Collector
- [x] Service Prometheus
- [x] Service Grafana
- [x] Service Loki
- [ ] Verifier que `docker compose up` demarre tout correctement
- [ ] Documenter les ports et credentials dans le README

**Criteres d'acceptation**:

- `docker compose up -d` lance tous les services sans erreur
- L'API repond sur `localhost:8080/health`
- Grafana accessible sur `localhost:3000`

---

### TICK-002 : Pipeline CI/CD (GitHub Actions)

**Priorite**: P0 | **RNCP**: A3.1, A3.3, A3.4, Ce3.1.2, Ce3.1.3, Ce3.3.1, Ce3.3.2, Ce3.4.1, Ce3.4.2, Ce3.4.4
**Assignable a**: Backend/DevOps

**Description**: Mettre en place un pipeline d'integration continue ET de deploiement continu automatise (Ce3.1.2 : build + tests sans intervention humaine, Ce3.1.3 : resolution fiable des erreurs rapidement).

**Taches**:

- [ ] Workflow CI backend : `go vet`, `go test -race -cover`, `golangci-lint`, `govulncheck`
- [ ] Workflow CI frontend : `flutter analyze`, `flutter test`
- [ ] Workflow CD : build Docker image et push sur registry (GHCR)
- [ ] Workflow de deploiement automatique sur merge dans main (Ce3.4.1, Ce3.4.4)
- [ ] Deploiement staging sur push dans develop
- [ ] Badge de statut CI dans le README
- [ ] Verification de la couverture de tests (seuil 80%)
- [ ] **Notifications d'echec** : alerter l'equipe en cas de build/test casse (Ce3.1.3)
- [ ] **Temps de resolution** : le pipeline doit fournir des logs clairs pour debug rapide (Ce3.1.3)
- [ ] Pipeline multi-plateforme : build Flutter iOS + Android (Ce3.4.2)

**Criteres d'acceptation**:

- Chaque push/PR declenche les tests automatiquement sans intervention humaine (Ce3.1.2)
- Le pipeline echoue si la couverture < 80% (backend)
- Les images Docker sont buildees automatiquement sur merge dans main
- Le deploiement en production est automatise (Ce3.4.1)
- Les builds Flutter iOS et Android sont generes (Ce3.4.2)
- Les echecs CI notifient l'equipe pour resolution rapide (Ce3.1.3)

---

### TICK-003 : Configuration 12-Factor App

**Priorite**: P0 | **RNCP**: A3.4, C3.4.1
**Assignable a**: Backend

**Description**: Garantir que toute la configuration est externalisee via variables d'environnement.

**Taches**:

- [x] Utiliser Viper pour charger la config depuis env vars
- [x] Creer `.env.example` avec toutes les variables documentees
- [ ] Valider que zero valeur n'est hardcodee dans le code
- [ ] Supporter les profils dev/staging/prod via `ENVIRONMENT`
- [ ] Documenter chaque variable dans le README

**Criteres d'acceptation**:

- Aucune valeur de configuration hardcodee
- L'app demarre correctement avec seulement des env vars (sans fichier .env)

---

### TICK-004 : Dockerfile multi-stage optimise

**Priorite**: P1 | **RNCP**: A3.4, A3.5, C3.4.1
**Assignable a**: Backend

**Description**: Image Docker optimisee (alpine/distroless) pour reduire la surface d'attaque.

**Taches**:

- [x] Stage builder avec golang:alpine
- [x] Stage final avec alpine minimal
- [ ] Verifier que l'image finale < 30MB
- [ ] Scanner l'image avec `trivy` ou `grype` pour vulnerabilites
- [ ] Ajouter un user non-root dans l'image

**Criteres d'acceptation**:

- Image finale < 30MB
- Aucune vulnerabilite critique dans le scan
- Le process tourne en non-root

---

### TICK-005 : Mise en production

**Priorite**: P1 | **RNCP**: A3.3, A3.5, C3.3.1, C3.5.1
**Assignable a**: DevOps

**Description**: Deployer l'application sur un environnement cloud (VPS, AWS, GCP ou Fly.io).

**Taches**:

- [ ] Choisir la plateforme de deploiement et documenter le choix (ADR)
- [ ] Configurer le deploiement automatique depuis la CI/CD
- [ ] Configurer un nom de domaine + certificat TLS (Let's Encrypt)
- [ ] Verifier que l'API est accessible en HTTPS
- [ ] Configurer les variables d'environnement de production

**Criteres d'acceptation**:

- L'API est accessible publiquement en HTTPS
- Le deploiement est automatise depuis main
- ADR documente le choix d'infra

---

## EPIC 2 : Authentification & Utilisateurs

### TICK-010 : Inscription utilisateur (Register)

**Priorite**: P0 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Backend

**Description**: Endpoint POST /api/v1/auth/register avec validation et hash bcrypt.

**Taches**:

- [x] Entity User avec roles (anonymous, user, broadcaster, admin)
- [x] Interface UserRepository
- [ ] Implementation GORM du UserRepository
- [x] UseCase Register avec hash bcrypt
- [ ] Handler HTTP avec validation des inputs (email format, password min length)
- [ ] Tests unitaires du usecase (table-driven, testify)
- [ ] Tests d'integration avec base de donnees reelle

**Criteres d'acceptation**:

- L'email et le username sont uniques (erreur 409 si doublon)
- Le mot de passe est hashe en bcrypt (jamais stocke en clair)
- Retourne un JWT a l'inscription
- Tests couvrent les cas nominaux et d'erreur

---

### TICK-011 : Connexion utilisateur (Login)

**Priorite**: P0 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Backend

**Description**: Endpoint POST /api/v1/auth/login retournant un JWT.

**Taches**:

- [x] UseCase Login (verify bcrypt + generate JWT)
- [x] JWT Manager (generate + validate)
- [ ] Handler HTTP avec validation
- [ ] Tests unitaires avec mocks (mockgen)
- [ ] Tests de securite : brute force, token expiration

**Criteres d'acceptation**:

- Login avec email + password valides retourne un JWT
- JWT contient user_id et role
- JWT expire apres une duree configurable
- Credentials invalides retournent 401

---

### TICK-012 : Middleware d'authentification JWT

**Priorite**: P0 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Backend

**Description**: Middleware Gin qui valide le token JWT sur les routes protegees.

**Taches**:

- [x] Middleware AuthMiddleware (extraction Bearer token, validation)
- [x] Middleware RequireRole (verification des roles)
- [ ] Tests du middleware avec tokens valides, expires, malformes
- [ ] Gestion du role "anonymous" pour les routes partiellement protegees

**Criteres d'acceptation**:

- Routes protegees retournent 401 sans token
- Routes admin retournent 403 si role != admin
- Le user_id est disponible dans le context Gin

---

### TICK-013 : CRUD Utilisateurs (Admin)

**Priorite**: P1 | **RNCP**: A3.1, C3.1.2
**Assignable a**: Backend

**Description**: Endpoints d'administration pour gerer les utilisateurs.

**Taches**:

- [ ] GET /api/v1/admin/users (liste paginee)
- [ ] PUT /api/v1/admin/users/:id/role (changer le role)
- [ ] DELETE /api/v1/admin/users/:id (desactiver un compte)
- [ ] Tests unitaires et d'integration

**Criteres d'acceptation**:

- Seuls les admins peuvent acceder a ces endpoints
- La pagination fonctionne (offset/limit)
- Le changement de role est immediatement effectif

---

### TICK-014 : Profil utilisateur (GET/PUT /me)

**Priorite**: P1 | **RNCP**: A3.1
**Assignable a**: Backend

**Description**: Permettre a l'utilisateur de consulter et modifier son profil.

**Taches**:

- [ ] GET /api/v1/users/me (retourne le profil)
- [ ] PUT /api/v1/users/me (modifier username, email)
- [ ] Validation des inputs
- [ ] Tests

**Criteres d'acceptation**:

- Le mot de passe n'est jamais retourne dans la reponse
- Le changement d'email verifie l'unicite

---

## EPIC 3 : Moteur de Streaming Audio

### TICK-020 : Architecture du streaming multiplexe

**Priorite**: P0 | **RNCP**: A3.1, A3.4, C3.1.1, C3.4.1
**Assignable a**: Backend

**Description**: Concevoir et implementer le moteur de streaming capable de diffuser un flux audio vers N clients simultanement via goroutines et channels.

**Taches**:

- [ ] Definir l'architecture du hub de streaming (pattern pub/sub)
- [ ] Implementer un StreamHub qui gere les connexions/deconnexions
- [ ] Utiliser des goroutines par listener avec des channels pour le broadcast
- [ ] Implementer le cleanup avec context.Context (eviter les memory leaks)
- [ ] Gerer la backpressure (buffer overflow quand un client est lent)
- [ ] Metriques : compter les listeners actifs, les deconnexions brutales
- [ ] Tests unitaires du hub avec race detector (`-race`)
- [ ] Benchmark : mesurer la consommation memoire pour N listeners

**Criteres d'acceptation**:

- Un broadcaster peut envoyer un flux audio
- N clients recoivent le flux en temps reel simultanement
- La deconnexion d'un client ne crash pas le serveur
- Les goroutines sont proprement nettoyees (pas de leak)
- Le test `-race` passe sans erreur

---

### TICK-021 : Endpoint de publication audio (Broadcaster)

**Priorite**: P0 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Backend

**Description**: POST /api/v1/streams/:id/publish - le broadcaster envoie le flux audio.

**Taches**:

- [ ] Handler qui lit le body en streaming (chunked transfer)
- [ ] Connecter au StreamHub pour broadcaster
- [ ] Valider que seul le role "broadcaster" ou "admin" peut publier
- [ ] Gerer le context (timeout, annulation)
- [ ] Mettre a jour le statut du stream (live/offline)
- [ ] Incrementer/decrementer la metrique ActiveStreams
- [ ] Tests avec un flux audio simule

**Criteres d'acceptation**:

- Le flux audio est recu et redistribue aux listeners
- Le statut du stream passe a "live" au debut et "offline" a la fin
- Un non-broadcaster recoit 403

---

### TICK-022 : Endpoint d'ecoute audio (Listener)

**Priorite**: P0 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Backend

**Description**: GET /api/v1/streams/:id/listen - un utilisateur recoit le flux audio en streaming.

**Taches**:

- [ ] Handler qui s'abonne au StreamHub et streame la reponse
- [ ] Headers corrects (Content-Type: audio/mpeg, Transfer-Encoding: chunked)
- [ ] Gestion du context (deconnexion propre du client)
- [ ] Incrementer/decrementer ActiveListeners
- [ ] Logger les deconnexions brutales (metrique business)
- [ ] Tests avec un client HTTP simulant un listener

**Criteres d'acceptation**:

- Le flux audio est recu en continu par le client
- La deconnexion du client est detectee et les ressources liberees
- Les metriques de listeners sont mises a jour en temps reel

---

### TICK-023 : CRUD Streams

**Priorite**: P1 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Backend

**Description**: Endpoints REST pour creer, lister, consulter et supprimer des streams.

**Taches**:

- [ ] Implementation GORM du StreamRepository
- [ ] POST /api/v1/streams (creer un stream, role broadcaster)
- [ ] GET /api/v1/streams (lister les streams live)
- [ ] GET /api/v1/streams/:id (detail d'un stream)
- [ ] DELETE /api/v1/streams/:id (supprimer, role owner ou admin)
- [ ] Tests unitaires et d'integration

**Criteres d'acceptation**:

- Un broadcaster peut creer un stream
- La liste montre les streams live en priorite
- Seul le proprietaire ou un admin peut supprimer

---

## EPIC 4 : Playlists & Tracks

### TICK-030 : CRUD Playlists

**Priorite**: P1 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Backend

**Description**: Gestion complete des playlists avec logique de file d'attente (Queue).

**Taches**:

- [ ] Implementation GORM du PlaylistRepository
- [ ] POST /api/v1/playlists (creer)
- [ ] GET /api/v1/playlists (lister les playlists de l'utilisateur)
- [ ] GET /api/v1/playlists/:id (detail avec tracks)
- [ ] PUT /api/v1/playlists/:id (modifier nom/description)
- [ ] DELETE /api/v1/playlists/:id
- [ ] POST /api/v1/playlists/:id/tracks (ajouter un track)
- [ ] DELETE /api/v1/playlists/:id/tracks/:trackId (retirer un track)
- [ ] Logique de queue (ordre des tracks, next/previous)
- [ ] Tests unitaires et d'integration

**Criteres d'acceptation**:

- CRUD complet fonctionnel
- L'ordre des tracks dans la playlist est maintenu
- Seul le proprietaire peut modifier sa playlist

---

### TICK-031 : Upload et gestion des Tracks

**Priorite**: P1 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Backend

**Description**: Upload de fichiers audio et gestion des tracks.

**Taches**:

- [ ] Implementation GORM du TrackRepository
- [ ] POST /api/v1/tracks (upload multipart/form-data)
- [ ] Stockage des fichiers (local en dev, S3/MinIO en prod)
- [ ] GET /api/v1/tracks (lister)
- [ ] GET /api/v1/tracks/:id (detail)
- [ ] DELETE /api/v1/tracks/:id
- [ ] Validation du format audio (mp3, aac, ogg)
- [ ] Extraction des metadonnees (duree, codec)
- [ ] Tests

**Criteres d'acceptation**:

- Upload de fichiers audio fonctionne
- Les formats non supportes sont rejetes
- Les fichiers sont stockes de maniere securisee

---

## EPIC 5 : Base de donnees & Persistance

### TICK-040 : Setup PostgreSQL + GORM migrations

**Priorite**: P0 | **RNCP**: A3.1, A3.4, C3.1.1, C3.4.1
**Assignable a**: Backend

**Description**: Configurer la connexion PostgreSQL et les migrations automatiques.

**Taches**:

- [ ] Initialiser la connexion GORM avec le DSN depuis config
- [ ] AutoMigrate pour les entites User, Stream, Playlist, Track
- [ ] Creer un script de migration versionne (pour les changements futurs)
- [ ] Creer un seed pour les donnees de test (admin user, streams de demo)
- [ ] Connection pooling configure (max open/idle connections)
- [ ] Tests d'integration avec une vraie base PostgreSQL

**Criteres d'acceptation**:

- Les tables sont creees automatiquement au demarrage
- Le seed cree un admin par defaut en dev
- La connexion gere les timeouts et reconnexions

---

### TICK-041 : Implementation des repositories GORM

**Priorite**: P0 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Backend

**Description**: Implementer toutes les interfaces repository avec GORM.

**Taches**:

- [ ] UserRepository GORM
- [ ] StreamRepository GORM
- [ ] PlaylistRepository GORM (avec relations many2many)
- [ ] TrackRepository GORM
- [ ] Utiliser context.Context dans toutes les requetes
- [ ] Tests d'integration pour chaque repository

**Criteres d'acceptation**:

- Toutes les interfaces du domain sont implementees
- Les requetes utilisent le context pour les timeouts
- Les tests passent contre une vraie base PostgreSQL

---

## EPIC 6 : Observabilite (Coeur du Bloc 3)

### TICK-050 : Instrumentation OpenTelemetry

**Priorite**: P0 | **RNCP**: A3.5, C3.5.1, C3.5.2, C3.5.4
**Assignable a**: Backend

**Description**: Instrumenter le code Go pour generer des traces distribuees end-to-end.

**Taches**:

- [x] Initialiser le TracerProvider avec export OTLP gRPC
- [ ] Instrumenter automatiquement Gin (otelgin middleware)
- [ ] Instrumenter GORM (otelgorm plugin)
- [ ] Ajouter des spans manuels sur les operations critiques (streaming, auth)
- [ ] Propager le trace context dans les headers HTTP
- [ ] Verifier qu'une trace complete (mobile -> API -> DB) est visible dans Grafana

**Criteres d'acceptation**:

- Chaque requete HTTP genere une trace
- Les requetes DB apparaissent comme child spans
- Le flux complet est visible dans Grafana Tempo/Jaeger

---

### TICK-051 : Logs structures JSON

**Priorite**: P0 | **RNCP**: A3.5, C3.5.1
**Assignable a**: Backend

**Description**: Tous les logs en format JSON structure via slog.

**Taches**:

- [x] Configurer slog avec JSONHandler
- [x] Middleware de logging des requetes HTTP
- [ ] Ajouter le trace_id dans chaque log (correlation traces/logs)
- [ ] Logger les evenements business (login, stream start/stop, disconnection)
- [ ] Configurer l'envoi des logs vers Loki
- [ ] Verifier les logs dans Grafana/Loki

**Criteres d'acceptation**:

- Tous les logs sont en JSON
- Chaque log contient : timestamp, level, message, trace_id, champs contextuels
- Les logs sont consultables dans Grafana via Loki

---

### TICK-052 : Metriques Prometheus (techniques + business)

**Priorite**: P0 | **RNCP**: A3.5, C3.5.1, C3.5.2, C3.5.4
**Assignable a**: Backend

**Description**: Exposer des metriques techniques ET business differenciees.

**Taches**:

- [x] Metriques techniques : http_requests_total, http_request_duration, error rates
- [x] Metriques business : active_streams, active_listeners, stream_disconnections
- [x] Middleware Gin pour collecter les metriques HTTP
- [ ] Endpoint /metrics expose pour Prometheus
- [ ] Ajouter des metriques de streaming : bytes envoyes, latence du flux
- [ ] Tests de l'endpoint /metrics

**Criteres d'acceptation**:

- /metrics retourne les metriques au format Prometheus
- On peut distinguer les erreurs techniques (5xx) des metriques business
- Prometheus scrape correctement les metriques

---

### TICK-053 : Dashboard Grafana

**Priorite**: P0 | **RNCP**: A3.5, C3.5.1, C3.5.4
**Assignable a**: Backend/DevOps

**Description**: Creer un dashboard Grafana complet avec metriques techniques et business.

**Taches**:

- [ ] Configurer les datasources (Prometheus, Loki, Tempo)
- [ ] Panel : Nombre d'utilisateurs en ligne
- [ ] Panel : Nombre de streams actifs
- [ ] Panel : Listeners par stream
- [ ] Panel : Debit de streaming (bytes/sec)
- [ ] Panel : Taux d'erreurs HTTP (5xx)
- [ ] Panel : Temps de reponse API (p50, p95, p99)
- [ ] Panel : Deconnexions brutales (metrique business)
- [ ] Panel : Logs en temps reel (Loki)
- [ ] Exporter le dashboard en JSON provisionne
- [ ] Configurer des alertes (ex: error rate > 5%)

**Criteres d'acceptation**:

- Le dashboard montre clairement la distinction entre metriques techniques et business
- Les panels se rafraichissent en temps reel
- Le dashboard est reproductible (JSON provisionne dans le repo)

---

### TICK-054 : Alerting et supervision

**Priorite**: P1 | **RNCP**: A3.5, C3.5.2
**Assignable a**: DevOps

**Description**: Configurer des alertes automatiques sur les metriques critiques.

**Taches**:

- [ ] Alertes Grafana : taux d'erreur > seuil
- [ ] Alertes Grafana : API response time > seuil
- [ ] Alertes Grafana : nombre de streams actifs = 0 (si attendu > 0)
- [ ] Canal de notification (email, Slack, webhook)
- [ ] Documenter la strategie d'alerting

**Criteres d'acceptation**:

- Les alertes se declenchent correctement
- Les notifications arrivent dans le canal configure

---

## EPIC 7 : Application Mobile Flutter

### TICK-060 : Ecran de Login/Register

**Priorite**: P0 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Frontend

**Description**: Ecrans d'authentification avec validation et gestion d'erreurs.

**Taches**:

- [x] AuthBloc (events: login, register, check, logout)
- [x] AuthRepository (appels API via Dio)
- [x] LoginScreen avec formulaire valide
- [ ] RegisterScreen avec champs email, username, password, confirm password
- [ ] Stockage securise du token JWT (flutter_secure_storage)
- [ ] Redirection automatique si token valide
- [ ] Gestion des erreurs (affichage snackbar)
- [ ] Tests du AuthBloc (bloc_test)

**Criteres d'acceptation**:

- L'utilisateur peut s'inscrire et se connecter
- Le token est stocke de maniere securisee
- Les erreurs sont affichees clairement
- Le bloc est teste

---

### TICK-061 : Lecteur audio avance

**Priorite**: P0 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Frontend

**Description**: Lecteur audio complet avec streaming, progression, volume et lecture en arriere-plan.

**Taches**:

- [ ] PlayerBloc (play, pause, stop, seek, volume)
- [ ] Integration just_audio pour le streaming HTTP
- [ ] Barre de progression interactive
- [ ] Controle du volume
- [ ] Integration audio_service pour lecture en arriere-plan
- [ ] Gestion des interruptions (appels entrants, notifications)
- [ ] Notification media controls (lock screen)
- [ ] Gestion du changement d'etat reseau (reconnexion auto)
- [ ] Tests du PlayerBloc

**Criteres d'acceptation**:

- Le streaming audio fonctionne en temps reel
- La lecture continue en arriere-plan
- Les controles apparaissent sur le lockscreen
- L'interface reste fluide (60 FPS)

---

### TICK-062 : Liste des streams live

**Priorite**: P0 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Frontend

**Description**: Ecran listant les streams disponibles avec possibilite d'ecouter.

**Taches**:

- [ ] StreamsBloc (load, refresh)
- [ ] StreamsRepository (API calls)
- [ ] StreamsScreen avec liste scrollable
- [ ] Pull-to-refresh
- [ ] Indicateur de stream live (badge, animation)
- [ ] Nombre de listeners affiche
- [ ] Tap sur un stream -> demarrer l'ecoute
- [ ] Shimmer loading effect
- [ ] Tests du bloc

**Criteres d'acceptation**:

- Les streams live sont affiches en temps reel
- Le nombre de listeners est mis a jour
- L'UX est fluide avec loading states

---

### TICK-063 : Interface Diffuseur

**Priorite**: P1 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Frontend

**Description**: Dashboard pour les broadcasters: lancer/arreter un flux en direct.

**Taches**:

- [ ] BroadcasterBloc (start stream, stop stream, status)
- [ ] Ecran de creation de stream (titre, description)
- [ ] Bouton Start/Stop avec feedback visuel
- [ ] Indicateur de statut (live/offline)
- [ ] Affichage du nombre de listeners en temps reel
- [ ] Capture audio depuis le micro du device
- [ ] Envoi du flux audio vers l'API
- [ ] Tests du bloc

**Criteres d'acceptation**:

- Le broadcaster peut demarrer et arreter un stream
- L'audio du micro est capture et envoye au serveur
- Le nombre de listeners est affiche en temps reel

---

### TICK-064 : Gestion des playlists (mobile)

**Priorite**: P1 | **RNCP**: A3.1, C3.1.1
**Assignable a**: Frontend

**Description**: CRUD playlists et gestion des tracks dans le mobile.

**Taches**:

- [ ] PlaylistsBloc (CRUD, add/remove track)
- [ ] PlaylistsRepository
- [ ] Ecran liste des playlists
- [ ] Ecran detail playlist avec liste de tracks
- [ ] Ajouter/retirer des tracks
- [ ] Reordonner les tracks (drag & drop)
- [ ] Lecture sequentielle des tracks d'une playlist
- [ ] Tests du bloc

**Criteres d'acceptation**:

- CRUD complet des playlists
- Les tracks peuvent etre reordonnees
- La lecture enchainee fonctionne

---

### TICK-065 : Ecran Admin (mobile)

**Priorite**: P2 | **RNCP**: A3.1
**Assignable a**: Frontend

**Description**: Interface admin pour gerer les utilisateurs et voir les stats globales.

**Taches**:

- [ ] AdminBloc (list users, change role, stats)
- [ ] Ecran liste des utilisateurs
- [ ] Modifier le role d'un utilisateur
- [ ] Ecran de statistiques globales (graphiques)
- [ ] Proteger l'acces (visible uniquement pour role admin)
- [ ] Tests du bloc

**Criteres d'acceptation**:

- Seuls les admins voient l'onglet admin
- Les stats globales sont affichees

---

### TICK-066a : Configuration complete du ThemeData

**Priorite**: P0 | **RNCP**: —
**Assignable a**: Frontend

**Description**: Configurer le ThemeData complet AVANT de commencer les views. Toutes les views doivent utiliser `Theme.of(context)` et jamais de valeurs hardcodees.

**Taches**:

- [ ] ColorScheme complet (primary, secondary, error, surface, background) light + dark
- [ ] TextTheme : definir les styles (headlineLarge, titleMedium, bodyLarge, labelSmall, etc.)
- [ ] Styles de boutons (FilledButton, OutlinedButton, TextButton) via ThemeData
- [ ] Style des InputDecoration (champs de formulaire) global
- [ ] Style des Cards (elevation, shape, padding)
- [ ] Style de l'AppBar global
- [ ] Style de la BottomNavigationBar / NavigationBar
- [ ] Style des SnackBar
- [ ] Dark mode complet et coherent
- [ ] Verifier qu'aucune view n'utilise de couleurs/polices en dur

**Criteres d'acceptation**:

- Tout le style passe par `Theme.of(context)` dans les widgets
- Le switch light/dark mode fonctionne sans casser l'UI
- Zero valeur de couleur ou de taille hardcodee dans les views

---

### TICK-066b : UI/UX Design moderne + Accessibilite

**Priorite**: P1 | **RNCP**: A3.6, C3.6.4
**Assignable a**: Frontend

**Description**: Design moderne, responsive iOS/Android, accessibilite.

**Taches**:

- [ ] Navigation bottom bar (Streams, Playlists, Broadcaster, Profile)
- [ ] Animations et transitions fluides
- [ ] Labels d'accessibilite (Semantics)
- [ ] Tester avec TalkBack/VoiceOver
- [ ] Responsive pour differentes tailles d'ecran (MediaQuery.sizeOf)

**Criteres d'acceptation**:

- L'app fonctionne sur iOS et Android
- Le dark mode fonctionne
- Les elements interactifs ont des labels d'accessibilite
- 60 FPS mesure avec Flutter DevTools

---

## EPIC 8 : Tests (Coeur du Bloc 3)

### TICK-070 : Plan de tests iteratif detaille + Cahier de recette

**Priorite**: P0 | **RNCP**: A3.2, Ce3.2.1, Ce3.2.2, Ce3.2.4
**Assignable a**: Commun

**Description**: Rediger un plan de tests iteratif complet couvrant TOUS les cas d'utilisation (Ce3.2.1), organise en parallele du developpement (Ce3.2.2), verifiant le bon fonctionnement selon les attentes documentees (Ce3.2.4).

**Taches**:

- [ ] Document `docs/PLAN_DE_TESTS.md`
- [ ] **Cahier de recette** : lister TOUS les cas d'utilisation et scenarii a valider (Ce3.2.1)
- [ ] Tests fonctionnels : par feature (auth, streaming, playlists, admin) mappes aux user stories
- [ ] Tests non-fonctionnels : performance, securite, charge, accessibilite
- [ ] Tests unitaires, d'integration, de securite detailles (Ce3.2.1)
- [ ] Definir les criteres de reussite pour chaque test
- [ ] **Planification iterative** : phases de tests en parallele du dev, pas en fin de projet (Ce3.2.2)
- [ ] Documenter les problemes techniques rencontres et les scenarii de regression
- [ ] Documenter les outils utilises (testify, mockgen, bloc_test, flutter_test, k6)
- [ ] **Verification des attentes** : chaque test est lie a une attente documentee des utilisateurs (Ce3.2.4)

**Criteres d'acceptation**:

- Le plan couvre tests unitaires, integration, securite, E2E
- Les tests fonctionnels et non-fonctionnels sont differencies
- Le cahier de recette est complet (tous les cas d'utilisation)
- La planification montre que les tests se font en parallele du dev (Ce3.2.2)
- Chaque test reference une attente documentee (Ce3.2.4)

---

### TICK-071 : Tests unitaires backend (80% couverture)

**Priorite**: P0 | **RNCP**: A3.1, A3.2, C3.1.1, C3.1.4, C3.2.3
**Assignable a**: Backend

**Description**: Atteindre 80% de couverture de tests unitaires sur le backend Go.

**Taches**:

- [ ] Tests des usecases (auth, streaming, playlists) avec mocks
- [ ] Tests des handlers HTTP (httptest)
- [ ] Tests des middlewares
- [ ] Tests du JWT Manager
- [ ] Tests du password hasher
- [ ] Generer les mocks avec mockgen pour les interfaces
- [ ] Table-driven tests systematiques
- [ ] Coverage report dans la CI (`go test -coverprofile`)

**Criteres d'acceptation**:

- Couverture >= 80%
- Tous les tests passent avec `-race`
- Les mocks sont generes automatiquement

---

### TICK-072 : Tests d'integration backend

**Priorite**: P0 | **RNCP**: A3.2, C3.2.1, C3.2.2, C3.2.3, C3.2.4
**Assignable a**: Backend

**Description**: Tests d'integration contre une vraie base PostgreSQL.

**Taches**:

- [ ] Setup/teardown de base de test (Docker ou testcontainers)
- [ ] Tests d'integration des repositories
- [ ] Tests d'integration des endpoints API (httptest + vraie DB)
- [ ] Tests du flux complet : register -> login -> create stream -> listen
- [ ] Corriger les problemes detectes lors des tests

**Criteres d'acceptation**:

- Les tests tournent contre une vraie base PostgreSQL
- Le flux E2E complet est teste
- Les corrections sont tracees dans les commits

---

### TICK-073 : Tests frontend (Flutter)

**Priorite**: P1 | **RNCP**: A3.1, A3.2, C3.1.4, C3.2.3
**Assignable a**: Frontend

**Description**: Tests unitaires des BLoCs et tests widget.

**Taches**:

- [ ] Tests du AuthBloc (bloc_test)
- [ ] Tests du PlayerBloc
- [ ] Tests du StreamsBloc
- [ ] Tests du PlaylistsBloc
- [ ] Tests widget du LoginScreen
- [ ] Tests widget du StreamsScreen
- [ ] Mocks avec mocktail

**Criteres d'acceptation**:

- Tous les BLoCs sont testes
- Les principaux widgets ont des tests
- `flutter test` passe sans erreur

---

### TICK-074 : Tests de performance / charge

**Priorite**: P1 | **RNCP**: A3.2, A3.5, C3.2.2, C3.5.4
**Assignable a**: Backend

**Description**: Tester la performance du serveur sous charge (N listeners simultanes).

**Taches**:

- [ ] Script de benchmark Go (`go test -bench`)
- [ ] Test de charge avec k6 ou vegeta
- [ ] Mesurer : latence API, throughput streaming, consommation memoire
- [ ] Tester avec 10, 50, 100 listeners simultanes
- [ ] Documenter les resultats et les limites
- [ ] Estimer le cout CPU/memoire par flux (justification des couts RNCP)

**Criteres d'acceptation**:

- Les benchmarks sont reproductibles
- Les resultats sont documentes
- L'estimation de cout est chiffree

---

## EPIC 8b : Conformite RGPD & Securite reglementaire

### TICK-075 : Conformite RGPD

**Priorite**: P0 | **RNCP**: A3.1, Ce3.1.4
**Assignable a**: Commun

**Description**: S'assurer que la solution respecte les contraintes reglementaires RGPD. Ce critere est OBLIGATOIRE (Ce3.1.4).

**Taches**:

- [ ] Politique de confidentialite affichee dans l'app (ecran settings)
- [ ] Consentement explicite a la collecte de donnees a l'inscription
- [ ] Endpoint DELETE /api/v1/users/me pour le droit a la suppression (droit a l'oubli)
- [ ] Endpoint GET /api/v1/users/me/data pour l'export des donnees personnelles (droit d'acces)
- [ ] Anonymisation des logs (ne pas logger d'email, IP hashee)
- [ ] Duree de retention des donnees definie et documentee
- [ ] Chiffrement des donnees sensibles au repos (mots de passe bcrypt = OK, verifier le reste)
- [ ] Registre de traitements simplifie dans `docs/RGPD.md`
- [ ] Mentions legales dans l'app et la doc

**Criteres d'acceptation**:

- L'utilisateur peut supprimer son compte et toutes ses donnees
- L'utilisateur peut exporter ses donnees personnelles
- Aucune donnee personnelle identifiable dans les logs
- Le registre des traitements est documente

---

### TICK-076 : Tests de securite automatises

**Priorite**: P1 | **RNCP**: A3.2, Ce3.2.1, Ce3.3.4
**Assignable a**: Backend

**Description**: Integrer des tests de securite dans le pipeline CI pour minimiser les vulnerabilites (Ce3.3.4).

**Taches**:

- [ ] Scanner de vulnerabilites des dependances Go (`govulncheck`)
- [ ] Scanner de vulnerabilites des dependances Flutter (`flutter pub outdated`)
- [ ] Scanner d'image Docker (`trivy` dans la CI)
- [ ] Tests OWASP basiques (injection SQL, XSS, CSRF)
- [ ] Audit des secrets (pas de secrets dans le code : `gitleaks` dans la CI)
- [ ] Rapport de securite genere a chaque build CI

**Criteres d'acceptation**:

- La CI echoue si une vulnerabilite critique est detectee
- Aucun secret n'est commite dans le repo
- Le rapport de securite est archive

---

## EPIC 9 : Securite

### TICK-080 : Securisation des endpoints

**Priorite**: P0 | **RNCP**: A3.1, A3.4, C3.1.1, C3.4.2
**Assignable a**: Backend

**Description**: Securiser tous les endpoints selon les bonnes pratiques OWASP.

**Taches**:

- [ ] Rate limiting sur les endpoints d'auth (anti brute-force)
- [ ] Validation des inputs (taille, format, injection)
- [ ] Headers de securite (X-Content-Type-Options, X-Frame-Options, etc.)
- [ ] CORS configure strictement en production
- [ ] Protection endpoint /metrics (pas expose publiquement en prod)
- [ ] Requetes SQL parametrees (GORM le fait par defaut)
- [ ] Tests de securite automatises

**Criteres d'acceptation**:

- Aucune injection SQL possible
- Rate limiting actif sur /auth/\*
- Les headers de securite sont presents

---

### TICK-081 : Scanning de vulnerabilites en continu (surveillance)

**Priorite**: P0 | **RNCP**: A3.3, Ce3.3.1, Ce3.3.4
**Assignable a**: DevOps

**Description**: Mettre en place une surveillance continue des vulnerabilites et de l'integrite des systemes (Ce3.3.4).

**Taches**:

- [ ] Dependabot ou Renovate sur le repo GitHub (alertes dependances)
- [ ] Workflow CI hebdomadaire de scan de vulnerabilites (`govulncheck`, `trivy`)
- [ ] Alertes GitHub Security activees
- [ ] Documenter la procedure de mise a jour en cas de vulnerabilite
- [ ] Verifier que la surveillance guide l'iteration (Ce3.3.3) : documenter les decisions prises suite aux alertes

**Criteres d'acceptation**:

- Dependabot/Renovate est actif et cree des PR automatiques
- Les alertes de securite sont traitees dans un delai defini
- Un log des decisions prises suite aux alertes est maintenu

---

### TICK-082 : TLS en production

**Priorite**: P1 | **RNCP**: A3.4, C3.4.2
**Assignable a**: DevOps

**Description**: Configurer HTTPS/TLS pour toutes les communications.

**Taches**:

- [ ] Certificat Let's Encrypt (ou reverse proxy avec TLS termination)
- [ ] Forcer HTTPS en production (redirect HTTP -> HTTPS)
- [ ] Configurer TLS 1.2 minimum
- [ ] Tester avec SSL Labs

**Criteres d'acceptation**:

- Toutes les communications sont chiffrees en prod
- Score SSL Labs >= A

---

## EPIC 9b : Feedback utilisateurs & Boucle d'amelioration

### TICK-083 : Mecanisme de retour utilisateurs

**Priorite**: P1 | **RNCP**: A3.4, Ce3.4.3, Ce3.3.2, Ce3.3.3
**Assignable a**: Full-stack

**Description**: Permettre la collecte de feedback utilisateurs pour evaluer l'adequation de la solution (Ce3.4.3) et orienter la feuille de route (Ce3.3.2).

**Taches**:

- [ ] Endpoint POST /api/v1/feedback (envoyer un feedback)
- [ ] Ecran de feedback dans l'app mobile (note + commentaire)
- [ ] Stockage en base (table Feedback : user_id, rating, comment, created_at)
- [ ] Endpoint GET /api/v1/admin/feedback (consulter les retours, role admin)
- [ ] Documenter comment les feedbacks orientent les priorites de dev (Ce3.3.2)
- [ ] Changelog visible dans l'app montrant les evolutions liees aux retours

**Criteres d'acceptation**:

- Les utilisateurs peuvent soumettre un feedback depuis l'app
- L'admin peut consulter tous les feedbacks
- Un document ou changelog montre que les retours influencent le developpement

---

### TICK-084 : Chaine d'outils integree et travail collaboratif

**Priorite**: P1 | **RNCP**: A3.5, Ce3.5.3
**Assignable a**: Commun

**Description**: Documenter et configurer la chaine d'outils integree permettant l'automatisation et le travail collaboratif (Ce3.5.3).

**Taches**:

- [ ] Documenter la toolchain complete dans `docs/TOOLCHAIN.md` : Git -> GitHub -> CI/CD -> Docker -> Registry -> Deploy -> Monitoring
- [ ] Schema visuel de la chaine d'outils (diagramme)
- [ ] Configurer les notifications CI/CD (Slack ou email aux membres)
- [ ] Board de gestion de projet (GitHub Projects ou equivalent)
- [ ] Documenter le workflow de collaboration (PR review, merge, deploy)

**Criteres d'acceptation**:

- La chaine d'outils est documentee avec un schema
- Le board de projet est actif et utilise
- Les notifications CI sont configurees

---

## EPIC 10 : Documentation (Critique pour A3.6)

### TICK-090 : README complet

**Priorite**: P0 | **RNCP**: A3.6, C3.6.1, C3.6.2
**Assignable a**: Commun

**Description**: README detaille du projet.

**Taches**:

- [ ] Presentation du projet et contexte
- [ ] Architecture technique (schema)
- [ ] Prerequisites (Go, Flutter, Docker)
- [ ] Installation et demarrage rapide
- [ ] Configuration (variables d'environnement)
- [ ] Structure du projet
- [ ] API endpoints documentation
- [ ] Comment lancer les tests
- [ ] Comment deployer
- [ ] Screenshots de l'app mobile
- [ ] Repartition des taches par membre

**Criteres d'acceptation**:

- Un nouveau developpeur peut setup le projet en suivant le README
- La repartition des taches est clairement documentee

---

### TICK-091 : Architecture Decision Records (ADR)

**Priorite**: P0 | **RNCP**: A3.6, C3.6.1, C3.6.3
**Assignable a**: Commun

**Description**: Documenter les choix techniques et leurs justifications.

**Taches**:

- [ ] ADR-001 : Choix de Go pour le backend (vs Node.js, Rust)
- [ ] ADR-002 : Choix de Flutter pour le mobile (vs React Native)
- [ ] ADR-003 : Choix de Gin comme framework HTTP (vs Echo, net/http)
- [ ] ADR-004 : Choix de BLoC pour le state management (vs Riverpod)
- [ ] ADR-005 : Choix de PostgreSQL (vs MongoDB, SQLite)
- [ ] ADR-006 : Strategie d'observabilite (OTEL + Prometheus + Grafana)
- [ ] ADR-007 : Choix de la plateforme de deploiement
- [ ] ADR-008 : Architecture streaming (pub/sub avec goroutines)

**Criteres d'acceptation**:

- Chaque ADR suit le format : Contexte, Decision, Consequences
- Les alternatives considerees sont mentionnees
- Les ADR sont dans `docs/adr/`

---

### TICK-092 : Documentation API (Swagger/OpenAPI)

**Priorite**: P1 | **RNCP**: A3.6, C3.6.1, C3.6.2
**Assignable a**: Backend

**Description**: Documentation API auto-generee avec Swag.

**Taches**:

- [ ] Installer swag (`go install github.com/swaggo/swag/cmd/swag`)
- [ ] Annoter les handlers avec les commentaires Swag
- [ ] Endpoint /swagger/ pour l'UI interactive
- [ ] Generer la spec OpenAPI
- [ ] Mettre a jour la doc a chaque changement d'API

**Criteres d'acceptation**:

- La doc Swagger est accessible et a jour
- Tous les endpoints sont documentes avec request/response schemas

---

### TICK-093 : Cahier des charges FR/EN

**Priorite**: P0 | **RNCP**: A3.6, Ce3.6.1, Ce3.6.2, Ce3.6.3
**Assignable a**: Commun

**Description**: Document de specifications techniques en francais et anglais (Ce3.6.2 exige niveau B2 minimum).

**Taches**:

- [ ] Contexte et objectifs du projet
- [ ] Specifications fonctionnelles avec **User Stories** detaillees (Ce3.6.1)
- [ ] Specifications techniques
- [ ] Architecture globale (schemas)
- [ ] **Schema de la base de donnees** (Ce3.6.1)
- [ ] **Schema general de la securite** (Ce3.6.1)
- [ ] **Diagrammes UML** : classes, sequences, cas d'utilisation (Ce3.6.1)
- [ ] **Diagrammes BPMN** : processus de streaming, processus d'inscription (Ce3.6.1)
- [ ] Exemples pertinents integres dans la doc (Ce3.6.1)
- [ ] Contraintes et limites
- [ ] **Version francaise** complete
- [ ] **Version anglaise** complete (equivalent B2, Ce3.6.2)
- [ ] Accessible (PDF ou site web)

**Criteres d'acceptation**:

- Le document est complet et professionnel
- Les deux versions (FR/EN) sont coherentes et de qualite B2+
- User stories, schema BDD, schema securite, UML et BPMN presents
- Les diagrammes utilisent un langage standardise (UML/BPMN)

---

### TICK-094 : Documentation technique integree + Diagrammes standardises

**Priorite**: P0 | **RNCP**: A3.6, Ce3.6.1, Ce3.6.2
**Assignable a**: Backend

**Description**: Documentation dans le code (godoc) et diagrammes standardises UML/BPMN (Ce3.6.1).

**Taches**:

- [ ] Commentaires godoc sur toutes les fonctions exportees
- [ ] **Diagramme de classe UML** : entites du domaine et leurs relations
- [ ] **Diagramme de sequence UML** : flux d'authentification (register -> login -> token)
- [ ] **Diagramme de sequence UML** : flux de streaming (publish -> broadcast -> listen)
- [ ] **Diagramme de deploiement UML** : conteneurs Docker, reseau, services
- [ ] **Diagramme de cas d'utilisation UML** : par role (anonymous, user, broadcaster, admin)
- [ ] **Diagramme BPMN** : processus CI/CD (commit -> test -> build -> deploy)
- [ ] **Diagramme BPMN** : processus de gestion des incidents (alerte -> diagnostic -> fix)
- [ ] Diagramme d'architecture (C4 model)
- [ ] Schema ERD de la base de donnees
- [ ] Guide de contribution (CONTRIBUTING.md)
- [ ] Tous les diagrammes dans `docs/diagrams/`

**Criteres d'acceptation**:

- `godoc` genere une doc lisible
- Au minimum 3 diagrammes UML et 2 diagrammes BPMN
- Les diagrammes utilisent la notation standardisee

---

### TICK-095 : Documentation versionee (Ce3.6.3)

**Priorite**: P0 | **RNCP**: A3.6, Ce3.6.3
**Assignable a**: Commun

**Description**: La documentation doit etre adaptee aux differentes versions de la solution (Ce3.6.3). Chaque release doit avoir sa doc associee.

**Taches**:

- [ ] Versionner la documentation avec les tags Git (v1.0.0, v1.1.0, etc.)
- [ ] Changelog automatique ou manuel (CHANGELOG.md) lie aux versions
- [ ] La doc API (Swagger) reflète la version deployee
- [ ] Les ADR sont dates et lies a la version concernee
- [ ] Le cahier des charges mentionne la version de la solution qu'il decrit
- [ ] Taguer les releases GitHub avec notes de release

**Criteres d'acceptation**:

- Chaque release a un tag Git et un changelog associe
- La doc Swagger correspond a la version deployee
- On peut retrouver la doc d'une version anterieure via Git

---

### TICK-096 : Documentation accessible handicap (Ce3.6.4)

**Priorite**: P0 | **RNCP**: A3.6, Ce3.6.4
**Assignable a**: Commun

**Description**: La documentation technique doit inclure des solutions pour les utilisateurs en situation de handicap : lisible, accessible, audible (Ce3.6.4).

**Taches**:

- [ ] Documentation en format HTML accessible (pas uniquement PDF)
- [ ] Structure semantique (headings, listes, alt-text sur les images/diagrammes)
- [ ] Contraste de couleurs suffisant (WCAG AA minimum)
- [ ] Documentation lisible par les lecteurs d'ecran (screen readers)
- [ ] Alternative textuelle pour tous les diagrammes et schemas
- [ ] Taille de police ajustable ou suffisamment grande
- [ ] Tester l'accessibilite avec un outil (axe, Lighthouse, WAVE)

**Criteres d'acceptation**:

- La doc passe un audit d'accessibilite basique (Lighthouse > 80)
- Tous les diagrammes ont une description textuelle alternative
- La doc est navigable au clavier et lisible par un screen reader

---

### TICK-097 : Plan de formation utilisateurs (Ce3.6.5)

**Priorite**: P0 | **RNCP**: A3.6, Ce3.6.5
**Assignable a**: Commun

**Description**: Rediger un plan de formation des utilisateurs adapte a la diversite du public, y compris les personnes en situation de handicap (Ce3.6.5).

**Taches**:

- [ ] Document `docs/PLAN_FORMATION.md`
- [ ] Identifier les differents profils d'utilisateurs (auditeur, diffuseur, admin)
- [ ] Pour chaque profil : guide de prise en main pas-a-pas avec captures d'ecran
- [ ] Tutoriels video ou screencast (ou lien vers)
- [ ] FAQ des problemes courants
- [ ] Adaptations pour les personnes en situation de handicap :
    - [ ] Instructions en texte clair (pas uniquement visuel)
    - [ ] Strategies d'enseignement alternatives (ex: description audio des ecrans)
    - [ ] Navigation clavier documentee
- [ ] Version du plan de formation en francais ET en anglais
- [ ] Parcours de formation progressif (debutant -> avance)

**Criteres d'acceptation**:

- Chaque profil utilisateur a un guide dedie
- Les adaptations handicap sont explicites et concretes
- Le plan est disponible en FR et EN
- Un nouvel utilisateur peut se former en autonomie

---

## EPIC 11 : Fonctionnalites Bonus (/5 pts)

### TICK-100 : Mode Offline (cache playlists)

**Priorite**: P3 | **RNCP**: A3.1
**Assignable a**: Frontend

**Description**: Mise en cache des tracks de playlists pour ecoute sans reseau.

**Taches**:

- [ ] Telecharger les tracks d'une playlist pour ecoute offline
- [ ] Detecter l'etat reseau (connectivity_plus)
- [ ] Jouer depuis le cache si pas de reseau
- [ ] Gestion de l'espace disque (limite, purge)
- [ ] Indicateur visuel du mode offline

---

### TICK-101 : WebSocket Chat en direct

**Priorite**: P3 | **RNCP**: A3.1
**Assignable a**: Full-stack

**Description**: Chat en temps reel entre les auditeurs d'un meme flux.

**Taches**:

- [ ] Endpoint WebSocket backend (gorilla/websocket)
- [ ] Hub de chat par stream (goroutines + channels)
- [ ] Widget de chat dans le mobile (web_socket_channel)
- [ ] Messages avec username et timestamp
- [ ] Moderation basique (admin peut muter)

---

### TICK-102 : Deploiement Kubernetes

**Priorite**: P3 | **RNCP**: A3.4, A3.5
**Assignable a**: DevOps

**Description**: Deployer sur un cluster K8s avec gestion des ressources.

**Taches**:

- [ ] Manifestes Kubernetes (Deployment, Service, Ingress)
- [ ] ConfigMap et Secrets pour la config
- [ ] Resources limits/requests
- [ ] HPA (Horizontal Pod Autoscaler)
- [ ] Helm chart (optionnel)

---

### TICK-103 : Algorithme de recommandation

**Priorite**: P3 | **RNCP**: A3.1
**Assignable a**: Backend

**Description**: Recommandations simples basees sur l'historique d'ecoute.

**Taches**:

- [ ] Table d'historique d'ecoute (user_id, track_id, listened_at)
- [ ] Endpoint GET /api/v1/recommendations
- [ ] Algorithme : collaborative filtering simple ou content-based
- [ ] Tests

---

### TICK-104 : Transcodage audio a la volee

**Priorite**: P3 | **RNCP**: A3.1
**Assignable a**: Backend

**Description**: Adapter la qualite du flux audio selon la bande passante du client.

**Taches**:

- [ ] Detection de la bande passante client
- [ ] Integration FFmpeg pour le transcodage
- [ ] Qualites : low (64kbps), medium (128kbps), high (320kbps)
- [ ] Selection automatique ou manuelle de la qualite

---

## EPIC 12 : Rendus & Soutenance

### TICK-110 : Depot de code, controle de version et strategie de branches

**Priorite**: P0 | **RNCP**: A3.1, Ce3.1.1, A3.3
**Assignable a**: Commun

**Description**: Configurer le depot Git, la strategie de branches/forks et le suivi de l'historique (Ce3.1.1). Ce critere est OBLIGATOIRE.

**Taches**:

- [ ] Creer le repository GitHub avec README, LICENSE et .gitignore
- [ ] Configurer GPG signing pour chaque membre (commits signes)
- [ ] Convention de commit documentee (Conventional Commits : feat/fix/docs/ci/...)
- [ ] **Strategie de branches documentee** (GitFlow ou trunk-based, avec schema)
- [ ] **Branches** : main (production), develop (integration), feature/_, hotfix/_
- [ ] **Protection de la branche main** : PR obligatoires, reviews requises, CI doit passer
- [ ] **Protection de la branche develop** : au moins 1 review
- [ ] Code review obligatoire avant merge (template de PR)
- [ ] **Suivi de l'historique** : pas de force push, historique lineaire ou merge commits documentes
- [ ] Documenter la strategie dans `docs/GIT_STRATEGY.md`

---

### TICK-111 : Presentation soutenance (10 min)

**Priorite**: P1
**Assignable a**: Commun

**Description**: Preparer la presentation collective de 10 minutes.

**Taches**:

- [ ] Slides : contexte et problematique
- [ ] Demo live du produit
- [ ] Montrer la coherence globale du projet
- [ ] Repeter la presentation

---

### TICK-112 : Preparation soutenance individuelle (20 min)

**Priorite**: P1 | **RNCP**: Tous les criteres
**Assignable a**: Chaque membre

**Description**: Preparer la soutenance individuelle devant le jury externe.

**Taches**:

- [ ] Documenter sa contribution personnelle
- [ ] Preparer les justifications de chaque choix technique
- [ ] Preparer des alternatives a proposer
- [ ] Identifier les limites et perspectives du projet
- [ ] S'entrainer a repondre aux questions techniques

---

### TICK-113 : Archive finale

**Priorite**: P0
**Assignable a**: Commun

**Description**: Preparer l'archive finale avec tout le code et les livrables.

**Taches**:

- [ ] Code source complet (backend + frontend)
- [ ] Documentation complete (README, ADR, cahier des charges, API doc)
- [ ] Docker Compose fonctionnel
- [ ] Instructions de deploiement
- [ ] Captures d'ecran de l'app
- [ ] Dashboard Grafana exporte
- [ ] Plan de tests et resultats

---

## Recapitulatif RNCP -> Tickets (mapping complet)

### A3.1 — Integration des changements de code (CI)

| Critere     | Libelle                                                   | Tickets couvrants  |
| ----------- | --------------------------------------------------------- | ------------------ |
| **Ce3.1.1** | Depot de code, controle de version, branches, historique  | TICK-110           |
| **Ce3.1.2** | CI : build + tests automatiques sans intervention humaine | TICK-002           |
| **Ce3.1.3** | CI : resolution fiable des erreurs rapidement             | TICK-002, 071, 072 |
| **Ce3.1.4** | Contraintes securite et RGPD respectees                   | TICK-075, 076, 080 |

### A3.2 — Tests automatises

| Critere     | Libelle                                                                  | Tickets couvrants            |
| ----------- | ------------------------------------------------------------------------ | ---------------------------- |
| **Ce3.2.1** | Plan de tests couvre tous les cas (unitaires, integration, securite)     | TICK-070, 071, 072, 073, 076 |
| **Ce3.2.2** | Planification en parallele du dev, cahier de recette                     | TICK-070, 072, 074           |
| **Ce3.2.3** | Tests automatises : identification/correction bogues + suivi performance | TICK-071, 072, 073, 074      |
| **Ce3.2.4** | Tests verifient le bon fonctionnement selon attentes documentees         | TICK-070, 072                |

### A3.3 — Surveillance continue des mises a jour

| Critere     | Libelle                                                       | Tickets couvrants       |
| ----------- | ------------------------------------------------------------- | ----------------------- |
| **Ce3.3.1** | Outils/methodologies avancees de surveillance                 | TICK-002, 053, 054, 081 |
| **Ce3.3.2** | Surveillance oriente la feuille de route dev                  | TICK-083, 053           |
| **Ce3.3.3** | Surveillance guide l'iteration, adaptation aux utilisateurs   | TICK-083, 081           |
| **Ce3.3.4** | Minimiser vulnerabilites, proteger integrite systemes/donnees | TICK-076, 080, 081, 082 |

### A3.4 — Distribution automatique (CD)

| Critere     | Libelle                                                     | Tickets couvrants |
| ----------- | ----------------------------------------------------------- | ----------------- |
| **Ce3.4.1** | Mise en production = distribution automatique               | TICK-002, 005     |
| **Ce3.4.2** | Distribution sur toutes les plateformes (iOS, Android, API) | TICK-002, 005     |
| **Ce3.4.3** | Deploiement permet retour utilisateurs                      | TICK-083          |
| **Ce3.4.4** | Deploiement continu = mises a jour frequentes               | TICK-002, 005     |

### A3.5 — Operations continues DevOps

| Critere     | Libelle                                                        | Tickets couvrants       |
| ----------- | -------------------------------------------------------------- | ----------------------- |
| **Ce3.5.1** | Operations continues, ajustements en cycles courts             | TICK-050, 051, 052, 053 |
| **Ce3.5.2** | Alertes et notifications detectent anomalies                   | TICK-054, 081           |
| **Ce3.5.3** | Chaine d'outils integree, automatisation, travail collaboratif | TICK-084                |
| **Ce3.5.4** | Garantir stabilite et performances                             | TICK-053, 074           |

### A3.6 — Documentation technique

| Critere     | Libelle                                                       | Tickets couvrants  |
| ----------- | ------------------------------------------------------------- | ------------------ |
| **Ce3.6.1** | Doc claire : user stories, BDD, securite, UML, BPMN, exemples | TICK-093, 094      |
| **Ce3.6.2** | Documentation FR et EN (B2)                                   | TICK-090, 092, 093 |
| **Ce3.6.3** | Doc adaptee aux differentes versions                          | TICK-091, 095      |
| **Ce3.6.4** | Doc accessible handicap (lisible, audible)                    | TICK-066, 096      |
| **Ce3.6.5** | Plan de formation adapte diversite + handicap                 | TICK-097           |

---

> **ATTENTION CRITIQUE** : Un seul critere "non acquis" invalide TOUT le Bloc 3.
> Les tickets P0 sont TOUS obligatoires. Verifier systematiquement que chaque Ce est couvert avant la soutenance.
>
> **Total** : 6 competences, 22 criteres d'evaluation, ~50 tickets dont ~30 P0.
