# Guide d'intégration

Ce guide décrit comment intégrer votre projet avec Turbotilt pour optimiser votre environnement de développement.

## Table des matières

- [Premiers pas](#premiers-pas)
- [Intégration avec des projets existants](#intégration-avec-des-projets-existants)
- [Dockerfiles personnalisés](#dockerfiles-personnalisés)
- [Configuration Docker Compose personnalisée](#configuration-docker-compose-personnalisée)
- [Configuration Tiltfile personnalisée](#configuration-tiltfile-personnalisée)
- [Intégration CI/CD](#intégration-cicd)

## Premiers pas

### Pour les nouveaux projets

1. Créez un nouveau répertoire pour votre projet
2. Initialisez votre projet Java (Spring Boot, Quarkus ou Micronaut)
3. Exécutez `turbotilt init` pour générer les fichiers de configuration nécessaires
4. Démarrez votre environnement de développement avec `turbotilt up`

### Pour les projets existants

1. Naviguez vers le répertoire de votre projet
2. Exécutez `turbotilt init` pour générer les fichiers de configuration
3. Révisez et personnalisez les fichiers générés selon vos besoins
4. Démarrez votre environnement de développement avec `turbotilt up`

## Intégration avec des projets existants

### Auto-détection

Turbotilt détecte automatiquement :

- Les frameworks Java (Spring Boot, Quarkus, Micronaut) à partir de `pom.xml` ou `build.gradle`
- Les dépendances de base de données (MySQL, PostgreSQL, MongoDB, etc.)
- Les dépendances de cache (Redis)
- Les dépendances de messagerie (Kafka, RabbitMQ)

### Fichier manifeste

Pour plus de contrôle, créez un fichier `turbotilt.yaml` à la racine de votre projet :

```yaml
services:
  - name: mon-service
    path: .
    java: "17"
    build: maven
    runtime: spring
    port: "8080"
    devMode: true
```

Consultez le [Guide de configuration](./configuration.fr.md) pour toutes les options disponibles.

## Dockerfiles personnalisés

Si vous avez un Dockerfile personnalisé, Turbotilt peut l'utiliser au lieu d'en générer un :

1. Créez un `Dockerfile` à la racine de votre projet
2. Exécutez `turbotilt init --use-existing-dockerfile`

Turbotilt détectera votre Dockerfile personnalisé et l'utilisera pour construire votre application.

Alternativement, vous pouvez spécifier le chemin du Dockerfile dans votre `turbotilt.yaml` :

```yaml
services:
  - name: mon-service
    path: .
    dockerfilePath: ./custom/Dockerfile
```

## Configuration Docker Compose personnalisée

Pour une configuration Docker Compose avancée :

1. Créez un fichier `docker-compose.yml` à la racine de votre projet
2. Exécutez `turbotilt init --use-existing-compose`

Turbotilt détectera votre fichier Docker Compose personnalisé et l'utilisera.

Vous pouvez également spécifier le chemin du fichier Docker Compose dans votre `turbotilt.yaml` :

```yaml
services:
  - name: mon-service
    path: .
    composeFilePath: ./custom/docker-compose.yml
```

## Configuration Tiltfile personnalisée

Pour une configuration Tilt avancée :

1. Créez un `Tiltfile` à la racine de votre projet
2. Exécutez `turbotilt init --use-existing-tiltfile`

Turbotilt détectera votre Tiltfile personnalisé et l'utilisera.

Vous pouvez également spécifier le chemin du Tiltfile dans votre `turbotilt.yaml` :

```yaml
services:
  - name: mon-service
    path: .
    tiltfilePath: ./custom/Tiltfile
```

## Intégration CI/CD

Turbotilt peut être intégré dans votre pipeline CI/CD :

### Exemple GitHub Actions

```yaml
name: Turbotilt CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Installation de Turbotilt
      run: curl -fsSL https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.sh | bash
    
    - name: Validation de la configuration
      run: turbotilt doctor --validate-manifest
    
    - name: Initialisation du projet
      run: turbotilt init --dry-run
    
    - name: Construction de l'image Docker
      run: turbotilt up --tilt=false --build-only
```

### Exemple Jenkins Pipeline

```groovy
pipeline {
    agent any
    
    stages {
        stage('Installation de Turbotilt') {
            steps {
                sh 'curl -fsSL https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.sh | bash'
            }
        }
        
        stage('Validation de la configuration') {
            steps {
                sh 'turbotilt doctor --validate-manifest'
            }
        }
        
        stage('Initialisation du projet') {
            steps {
                sh 'turbotilt init --dry-run'
            }
        }
        
        stage('Construction de l\'image Docker') {
            steps {
                sh 'turbotilt up --tilt=false --build-only'
            }
        }
    }
}
```

Pour plus d'informations sur l'intégration de Turbotilt avec votre pipeline CI/CD, veuillez consulter la documentation de votre plateforme CI/CD.
