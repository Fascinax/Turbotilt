# Commande Turbotilt Temporary Up (tup)

La commande `tup` (Temporary Up) permet de générer des fichiers de configuration, démarrer des services, et nettoyer automatiquement à la fin. Ce document explique comment utiliser cette commande et ses différents cas d'utilisation.

## Table des matières

- [Aperçu de la commande](#aperçu-de-la-commande)
- [Utilisation de base](#utilisation-de-base)
- [Options de la commande](#options-de-la-commande)
- [Cas d'utilisation](#cas-dutilisation)
  - [Environnements de développement temporaires](#environnements-de-développement-temporaires)
  - [Évaluation de Turbotilt sur des projets existants](#évaluation-de-turbotilt-sur-des-projets-existants)
  - [Microservices multi-frameworks](#microservices-multi-frameworks)
  - [Pipelines CI/CD](#pipelines-cicd)
  - [Collaboration d'équipe](#collaboration-déquipe)
- [Architectures d'exemple](#architectures-dexemple)
- [Bonnes pratiques](#bonnes-pratiques)

## Aperçu de la commande

La commande `tup` combine les fonctionnalités de `init`, `up`, et `stop` avec le nettoyage automatique des fichiers de configuration générés. Elle est conçue pour les scénarios où vous souhaitez rapidement démarrer un environnement de développement sans ajouter de façon permanente des fichiers de configuration à votre projet.

## Utilisation de base

```bash
# Démarrer un environnement de développement temporaire
turbotilt tup

# Démarrer un service spécifique dans un environnement temporaire
turbotilt tup --service mon-service

# Démarrer en mode détaché (arrière-plan)
turbotilt tup --detached
```

Lorsque vous exécutez `turbotilt tup` :

1. Les fichiers de configuration sont générés (Dockerfile, docker-compose.yml, Tiltfile)
2. Les services sont démarrés (avec rechargement en direct par défaut)
3. Lorsque vous appuyez sur Ctrl+C, les services sont arrêtés et les fichiers de configuration sont supprimés

## Options de la commande

| Option | Description |
|--------|-------------|
| `--tilt` | Utiliser Tilt pour le rechargement en direct (par défaut : true) |
| `--detached` | Exécuter en mode détaché (arrière-plan) |
| `--service [nom]` | Exécuter uniquement un service spécifique |
| `--debug` | Activer le mode débogage avec sortie détaillée |
| `--dry-run` | Simuler l'exécution sans apporter de modifications |

## Cas d'utilisation

### Environnements de développement temporaires

**Scénario** : Vous avez besoin de configurer rapidement un environnement de développement sans ajouter de fichiers de configuration à votre projet.

**Solution** : Utilisez `turbotilt tup` pour générer des fichiers temporaires, démarrer l'environnement et nettoyer à la fin.

```bash
cd mon-projet-java
turbotilt tup
# Travailler avec le rechargement à chaud activé
# Appuyer sur Ctrl+C lorsque vous avez terminé
```

**Avantages** :
- Pas de fichiers de configuration encombrant votre espace de travail
- Pas besoin d'ajouter des fichiers à `.gitignore`
- Répertoire de travail propre pour le contrôle de version

### Évaluation de Turbotilt sur des projets existants

**Scénario** : Vous souhaitez voir comment Turbotilt configurerait votre projet sans apporter de modifications permanentes.

```bash
cd projet-existant
turbotilt tup --detached
# Vérifier comment les services fonctionnent
turbotilt stop
```

**Avantages** :
- Évaluation sans risque
- Pas de conflits avec la configuration existante
- Test avant de s'engager avec Turbotilt

### Microservices multi-frameworks

**Scénario** : Vous avez une architecture de microservices complexe utilisant différents frameworks Java (Spring Boot, Quarkus, Micronaut, etc.).

```bash
cd projet-microservices
turbotilt tup
```

**Avantages** :
- Détection automatique des différents frameworks
- Environnement de développement unifié
- Expérience développeur standardisée entre les équipes

### Pipelines CI/CD

**Scénario** : Vous devez démarrer temporairement des services pour exécuter des tests dans un pipeline CI/CD.

```bash
# Dans un script CI
turbotilt tup --detached
sleep 10  # Attendre que les services démarrent
./run-integration-tests.sh
turbotilt stop
```

**Avantages** :
- Environnement de test propre pour chaque exécution CI
- Pas de gestion de configuration requise
- Environnements de test isolés

### Collaboration d'équipe

**Scénario** : Les membres de l'équipe ont des préférences de configuration ou des besoins de configuration différents.

```bash
git pull  # Obtenir le code le plus récent
turbotilt tup  # Démarrer avec la configuration détectée automatiquement
```

**Avantages** :
- Pas de conflits de fichiers de configuration
- Expérience développeur cohérente
- Chaque développeur obtient un environnement frais

## Architectures d'exemple

La commande `tup` est particulièrement utile pour ces modèles d'architecture :

### Monorepo avec plusieurs projets

```
monorepo/
├── service-a/   # Service Spring Boot
├── service-b/   # Service Quarkus
└── frontend/    # Frontend Angular
```

Chaque service peut être détecté et configuré de manière appropriée, même avec différents frameworks.

### Projets polyglotte

```
polyglot-app/
├── api/         # API Java
├── workers/     # Traitement Python
└── dashboard/   # Frontend Node.js
```

Turbotilt peut générer des configurations appropriées pour les composants Java tout en travaillant aux côtés d'autres services en d'autres langages.

### Microservices avec dépendances partagées

```
projet/
├── user-service/     # Spring Boot
├── order-service/    # Micronaut
├── product-service/  # Quarkus
└── shared-lib/       # Code commun
```

La commande `tup` peut détecter les relations entre les services et les configurer pour qu'ils fonctionnent ensemble.

## Bonnes pratiques

1. **Test initial** : Utilisez `tup` pour les tests initiaux avant de vous engager avec des fichiers de configuration permanents.

2. **Intégration CI** : Incorporez `tup` dans les pipelines CI pour des environnements de test propres.

3. **Environnements de démonstration** : Utilisez pour créer rapidement des environnements de démonstration sans surcharge de configuration.

4. **Intégration** : Simplifiez l'intégration en laissant les nouveaux membres de l'équipe commencer avec `tup` au lieu d'une configuration manuelle.

5. **Projets existants** : Testez comment Turbotilt fonctionne avec des projets existants sans les modifier.
