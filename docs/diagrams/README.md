# Diagrammes

Ce dossier rassemble les diagrammes **standardisés** du projet (UML, BPMN, C4) au-delà de leurs versions intégrées en Mermaid dans les documents Markdown.

## Sources Mermaid (versionnables)

Les diagrammes sont **d'abord** maintenus en Mermaid directement dans :

- [`../cahier-des-charges-fr.md`](../cahier-des-charges-fr.md) (sections 6, 7, 8, 9, 10)
- [`../cahier-des-charges-en.md`](../cahier-des-charges-en.md) (sections 6, 7, 8, 9, 10)
- Les ADR (`../adr/0008-streaming-pubsub.md`, etc.)

Mermaid est rendu nativement par GitHub : versionnable, accessible, lisible en revue.

## Exports BPMN 2.0 stricts

Pour les diagrammes BPMN à exporter en notation BPMN 2.0 standardisée (`.bpmn`) afin d'être ouverts dans des éditeurs spécialisés (bpmn.io, Camunda Modeler, Signavio) :

| Fichier prévu | Contenu | Source |
|---|---|---|
| `bpmn-registration.bpmn` | Processus d'inscription | Cahier des charges § 10.1 |
| `bpmn-broadcasting.bpmn` | Processus de diffusion d'un stream | Cahier des charges § 10.2 |
| `bpmn-cicd.bpmn` | Processus CI/CD | Cahier des charges § 10.3 |
| `bpmn-incident.bpmn` | Processus de gestion des incidents | Cahier des charges § 10.4 |

## Exports UML

| Fichier prévu | Contenu | Source |
|---|---|---|
| `uml-class.puml` | Diagramme de classes domaine | Cahier des charges § 9.1 |
| `uml-usecases.puml` | Cas d'utilisation | Cahier des charges § 9.2 |
| `uml-seq-auth.puml` | Séquence inscription / login | Cahier des charges § 9.3 |
| `uml-seq-streaming.puml` | Séquence diffusion | Cahier des charges § 9.4 |
| `uml-deployment.puml` | Déploiement | Cahier des charges § 9.5 |

## Outils recommandés

- **Mermaid Live Editor** — https://mermaid.live (copier-coller le bloc depuis le Markdown).
- **bpmn.io** — https://bpmn.io (BPMN 2.0).
- **PlantUML** — https://plantuml.com (UML standardisé exportable en PNG/SVG).
- **draw.io / diagrams.net** — pour les schémas C4 plus libres.

## Conventions

- Tout diagramme exporté en `.bpmn` ou `.puml` reste **lisible en source** dans Git.
- Les exports `.png` / `.svg` sont régénérés à chaque release et accompagnés d'une **description textuelle alternative** dans le Markdown qui les référence (Ce3.6.4).
- Pas de fichiers binaires propriétaires non versionnables (Visio, Lucidchart…).

## Génération automatique (CI prévue)

Un job CI optionnel exécutera :

```bash
mmdc -i ../cahier-des-charges-fr.md -o ./mermaid-fr/  # rendu PNG des Mermaid
plantuml -tsvg uml-*.puml                              # rendu SVG des UML
```

afin de pré-générer les images pour la version HTML/PDF accessible de la documentation.
