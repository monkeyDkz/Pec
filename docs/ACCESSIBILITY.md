# Politique d'accessibilité — StreamPulse

| | |
|---|---|
| **Version** | 1.0.0 |
| **Couverture RNCP** | A3.6 — Ce3.6.4 |
| **Référentiels** | WCAG 2.1 niveau AA · RGAA 4.1 |
| **Statut** | Adopté |

> Ce document décrit la stratégie d'accessibilité de l'application mobile **et** de la documentation technique. L'objectif est que StreamPulse soit utilisable par les personnes en situation de handicap (visuel, moteur, cognitif).

---

## 1. Principes directeurs

Inspirés des quatre piliers WCAG (POUR) :

| Pilier | Engagement |
|---|---|
| **Perceptible** | Tout contenu utile a une alternative textuelle ou audio. Contraste suffisant. |
| **Utilisable** | Navigable au clavier, au lecteur d'écran, par des cibles tactiles ≥ 48×48 dp. |
| **Compréhensible** | Langage clair, ordre logique de lecture, messages d'erreur explicites. |
| **Robuste** | Compatibilité avec TalkBack (Android), VoiceOver (iOS), lecteurs d'écran de bureau. |

---

## 2. Application mobile

### 2.1 Exigences techniques

| Critère | Implémentation Flutter |
|---|---|
| Contraste texte / fond | Ratio ≥ 4.5:1 (texte normal), ≥ 3:1 (texte large). Vérifié via Material 3 ColorScheme. |
| Cibles tactiles | ≥ 48×48 dp (recommandation Material). |
| Labels d'accessibilité | Widget `Semantics(label: ...)` pour toute action sans texte visible. |
| Ordre de focus | `FocusTraversalGroup`, ordre logique de lecture. |
| Texte évolutif | Respect du paramètre système (`MediaQuery.textScalerOf(context)`). |
| Couleur seule | Jamais comme unique vecteur d'information (toujours doublé d'une icône ou texte). |
| Animations | Respect de `MediaQuery.disableAnimationsOf(context)` (réduction d'animation). |
| Lecture audio | Contrôles accessibles depuis le lockscreen via `audio_service`. |

### 2.2 Patterns d'implémentation

```dart
// Icône-bouton avec label sémantique
IconButton(
  icon: const Icon(Icons.play_arrow),
  tooltip: 'Lancer la lecture',
  onPressed: _play,
);

// Image décorative ignorée par le lecteur d'écran
ExcludeSemantics(child: Image.asset('assets/wave.png'));

// Image porteuse de sens
Semantics(
  label: 'Pochette de l\'album Atmospheric Sessions',
  image: true,
  child: Image.network(track.cover),
);

// Annonce dynamique
SemanticsService.announce(
  'Stream en direct démarré',
  Directionality.of(context),
);
```

### 2.3 Tests obligatoires

| Test | Outil | Fréquence |
|---|---|---|
| Navigation complète au lecteur d'écran | TalkBack (Android) | Avant chaque release |
| Idem | VoiceOver (iOS) | Avant chaque release |
| Contraste visuel | Outil Figma + DevTools | Sur chaque écran |
| Tailles de police × 2 | Paramètre système | Avant chaque release |
| Daltonisme (deuteranopie, protanopie) | Filtres macOS / Android | Avant chaque release |

### 2.4 Roadmap fonctionnelle

| US | Statut |
|---|---|
| US-051 — Navigation au lecteur d'écran | ✅ Cible v1.0.0 |
| Sous-titrage / transcription du flux audio | 🟡 Bonus v1.1 |
| Mode contraste élevé | 🟡 Bonus v1.1 |

---

## 3. Documentation technique

### 3.1 Format

- **Markdown sémantique** (titres `<h1>` à `<h4>` hiérarchisés, listes, tableaux).
- Rendu HTML accessible via GitHub Pages ou `pandoc` (`pandoc -s --toc -o docs.html`).
- PDF généré uniquement en complément (jamais comme seul format).

### 3.2 Critères

| Critère | Action |
|---|---|
| Structure | Titres hiérarchisés, pas de saut de niveau. |
| Alternatives | Tous les diagrammes ont une description textuelle (titre + paragraphe explicatif). |
| Images | Texte alternatif `![alt](url)` obligatoire. |
| Contraste | Thème GitHub clair et sombre vérifiés. |
| Lien | Libellé explicite (jamais « cliquer ici »). |
| Tableaux | En-têtes (`|---|`) et navigation logique. |
| Mermaid | Description textuelle accompagnant chaque diagramme. |

### 3.3 Audit

Lighthouse cible : **score accessibilité ≥ 80**.

```bash
# Audit local du site doc
npx http-server docs &
npx lighthouse http://localhost:8080 \
  --only-categories=accessibility \
  --output=html \
  --output-path=./docs/lighthouse-report.html
```

---

## 4. Communication accessible

| Canal | Adaptations |
|---|---|
| Messages d'erreur | Texte clair, pas uniquement un code (« 401 » → « Identifiants incorrects »). |
| Notifications | Doublées d'un signal visuel **et** audio. |
| Vidéos de formation *(plan formation)* | Sous-titrées + transcription textuelle. |
| Tutoriels | Description textuelle pas-à-pas en plus des captures. |

---

## 5. Procédure de signalement

Un utilisateur rencontrant un obstacle peut :

1. Utiliser l'écran *Paramètres → Signaler un problème d'accessibilité*.
2. Écrire à `contact@ecole-decode.fr`.

Réponse sous 7 jours ouvrés. Tickets ouverts dans GitHub avec label `a11y`.

---

## 6. Plan d'amélioration continue

| # | Action | Sprint |
|---|---|---|
| 1 | Tests automatisés a11y dans la CI Flutter | S4 |
| 2 | Audit externe par une personne en situation de handicap | Pré-release |
| 3 | Documentation sous-titrée (transcription audio) | Bonus v1.1 |
| 4 | Mode contraste élevé | Bonus v1.1 |
| 5 | Génération automatique de HTML accessible des docs (pandoc en CI) | S5 |

---

## 7. Annexes

- [WCAG 2.1 — Comprendre les directives](https://www.w3.org/WAI/WCAG21/Understanding/)
- [RGAA 4.1](https://accessibilite.numerique.gouv.fr/)
- [Material Design — Accessibility](https://m3.material.io/foundations/accessible-design)
- [Flutter — Accessibility](https://docs.flutter.dev/accessibility-and-internationalization/accessibility)
- [Lighthouse](https://developer.chrome.com/docs/lighthouse/accessibility/)
