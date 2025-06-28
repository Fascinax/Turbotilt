# Guide d'utilisation

Ce guide fournit des instructions détaillées sur l'utilisation efficace de Turbotilt pour vos projets de développement Java.

## Table des matières

- [Aperçu des commandes](#aperçu-des-commandes)
- [Initialisation d'un projet](#initialisation-dun-projet)
- [Démarrage de votre environnement](#démarrage-de-votre-environnement)
- [Vérification de votre environnement](#vérification-de-votre-environnement)
- [Arrêt de votre environnement](#arrêt-de-votre-environnement)
- [Utilisation avancée](#utilisation-avancée)

## Aperçu des commandes

Turbotilt propose les commandes principales suivantes :

| Commande | Description |
|----------|-------------|
| `init`   | Initialise un projet en créant Dockerfile, docker-compose.yml et Tiltfile |
| `up`     | Démarre l'environnement de développement avec Tilt ou Docker Compose |
| `tup`    | Temporary Up - génère les configurations, démarre les services et nettoie à la fin |
| `select` | Détecte les microservices dans un répertoire et permet de sélectionner lesquels lancer |
| `doctor` | Vérifie l'environnement et la configuration, fournissant des diagnostics |
| `stop`   | Arrête l'environnement et nettoie les ressources |
| `version`| Affiche la version actuelle de Turbotilt |

## Initialisation d'un projet

La commande `init` analyse votre projet, détecte le framework Java et crée les fichiers nécessaires pour votre environnement de développement.

### Initialisation de base

```bash
turbotilt init
```

Cela va :
1. Scanner votre projet pour détecter le framework Java (Spring Boot, Quarkus, Micronaut)
2. Générer un Dockerfile approprié
3. Créer un fichier docker-compose.yml, incluant les services dépendants (si détectés)
4. Générer un Tiltfile pour le live reload

### Options

```bash
# Spécifier explicitement le framework
turbotilt init --framework spring

# Définir le port de l'application
turbotilt init --port 8080

# Spécifier la version JDK
turbotilt init --jdk 17

# Activer le mode développement (par défaut)
turbotilt init --dev

# Générer un fichier manifeste
turbotilt init --generate-manifest

# Initialiser à partir d'un manifeste existant
turbotilt init --from-manifest
```

## Démarrage de votre environnement

La commande `up` démarre votre environnement de développement en utilisant Tilt (par défaut) ou Docker Compose.

### Démarrage de base

```bash
turbotilt up
```

Cela va :
1. Construire votre application en utilisant le Dockerfile généré
2. Démarrer tous les services définis dans docker-compose.yml
3. Configurer le live reload avec Tilt
4. Afficher les logs de tous les services

### Options

```bash
# Démarrer sans Tilt (Docker Compose uniquement)
turbotilt up --tilt=false

# Démarrer un service spécifique du manifeste
turbotilt up --service payment-service

# Activer le mode debug avec logs détaillés
turbotilt up --debug
```

## Vérification de votre environnement

La commande `doctor` vérifie votre environnement et votre configuration, vous aidant à résoudre les problèmes.

```bash
# Vérification de santé basique
turbotilt doctor

# Valider le fichier manifeste
turbotilt doctor --validate-manifest

# Vérification détaillée avec sortie verbeuse
turbotilt doctor --verbose --log
```

La commande doctor vérifie :
- L'installation et la configuration de Docker et Docker Compose
- L'installation de Tilt pour le live reload
- L'environnement JDK et Java
- La configuration réseau et les permissions
- La syntaxe et la validité du manifeste

## Arrêt de votre environnement

La commande `stop` arrête votre environnement et nettoie les ressources.

```bash
turbotilt stop
```

Cela va :
1. Arrêter tous les conteneurs en cours d'exécution
2. Supprimer les ressources temporaires
3. Conserver vos fichiers de configuration intacts

## Utilisation avancée

### Flags globaux

Toutes les commandes acceptent ces options :

- `--dry-run` : Simule l'exécution sans effectuer de modifications réelles
- `--debug` : Active le mode debug avec sortie détaillée
- `--config-file` : Spécifie un chemin de fichier de configuration personnalisé

### Auto-update des Tiltfiles

En mode développeur, Turbotilt surveille automatiquement les modifications de vos fichiers sources et met à jour les Tiltfiles en conséquence, garantissant que vos changements sont toujours pris en compte.

### Projets multi-services

Pour les projets avec plusieurs services, vous pouvez définir tous les services dans un fichier manifeste (`turbotilt.yaml`). Voir le [Guide de configuration](./configuration.fr.md) pour plus de détails.

```bash
# Démarrer tous les services définis dans le manifeste
turbotilt up

# Démarrer un service spécifique
turbotilt up --service user-service
```

### Travailler avec des projets existants

Pour intégrer Turbotilt à un projet existant :

1. Naviguez vers le répertoire de votre projet
2. Exécutez `turbotilt init` pour générer les fichiers de configuration
3. Personnalisez les fichiers générés selon vos besoins
4. Démarrez votre environnement avec `turbotilt up`

Pour des informations d'intégration plus détaillées, consultez le [Guide d'intégration](./integration.fr.md).

## Utilisation d'environnements temporaires avec `tup`

La commande `tup` (Temporary Up) permet de générer des fichiers de configuration, démarrer des services et nettoyer automatiquement à la fin.

### Utilisation de base

```bash
# Démarrer un environnement temporaire (configurations générées et nettoyées à la sortie)
turbotilt tup

# Démarrer un service spécifique
turbotilt tup --service auth-service

# Démarrer en mode détaché (arrière-plan)
turbotilt tup --detached
```

### Principaux avantages

- Pas de fichiers de configuration laissés dans votre espace de travail
- Idéal pour tester rapidement des projets
- Parfait pour les pipelines CI/CD
- Collaboration d'équipe simplifiée avec détection automatique

Pour une documentation complète sur la commande `tup` et ses cas d'utilisation, consultez le [Guide de la commande Tup](./tup-command.fr.md).

## Sélection de microservices avec `select`

La commande `select` vous permet de scanner un répertoire pour détecter les microservices et choisir lesquels lancer, offrant une flexibilité pour les grands projets multi-services.

### Utilisation de base

```bash
# Scanner le répertoire actuel et sélectionner interactivement les services
turbotilt select

# Scanner un répertoire spécifique
turbotilt select ./mon-repertoire-de-projet

# Sélectionner des services, créer un fichier de configuration et les lancer
turbotilt select --create-config --launch
```

### Principaux avantages

- Sélection interactive des services dans les grands projets
- Possibilité de créer des groupements de services personnalisés
- Support pour des configurations temporaires ou permanentes
- Idéal pour l'intégration de nouveaux membres d'équipe dans des architectures complexes

Pour des scénarios d'utilisation détaillés et des options, consultez le [Guide de la commande Select](./select-command.fr.md).
