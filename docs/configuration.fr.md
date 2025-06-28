# Guide de configuration

Ce guide fournit des informations détaillées sur les options de configuration de Turbotilt et le format du manifeste.

## Table des matières

- [Méthodes de configuration](#méthodes-de-configuration)
- [Format du fichier manifeste](#format-du-fichier-manifeste)
- [Configuration des services](#configuration-des-services)
- [Services dépendants](#services-dépendants)
- [Variables d'environnement](#variables-denvironnement)
- [Configuration des volumes](#configuration-des-volumes)
- [Exemples](#exemples)

## Méthodes de configuration

Turbotilt peut être configuré de plusieurs façons :

1. **Flags en ligne de commande** : Pour des réglages rapides
2. **Fichier manifeste (`turbotilt.yaml`)** : Pour une configuration persistante du projet
3. **Auto-détection** : Pour la détection du framework et des services

L'ordre de priorité est :
1. Flags en ligne de commande (priorité la plus élevée)
2. Configuration du fichier manifeste
3. Valeurs auto-détectées (priorité la plus basse)

## Format du fichier manifeste

Turbotilt supporte une approche déclarative avec un fichier manifeste YAML qui définit tous les services de votre projet.

### Structure de base

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

### Exemple multi-services

```yaml
services:
  - name: user-service
    path: services/user
    java: "17"
    build: maven
    runtime: spring
    port: "8081"
  
  - name: payment-service
    path: services/payment
    java: "21"
    runtime: quarkus
    port: "8082"
```

## Configuration des services

Voici toutes les options disponibles pour configurer vos services Java :

| Option | Description | Valeurs possibles | Défaut |
|--------|-------------|-------------------|--------|
| `name` | Nom du service | Chaîne de caractères | Nom du répertoire |
| `path` | Chemin relatif vers le service | Chemin | `.` |
| `java` | Version de Java | `"8"`, `"11"`, `"17"`, `"21"` | `"17"` |
| `build` | Système de build | `maven`, `gradle` | Auto-détecté |
| `runtime` | Framework Java | `spring`, `quarkus`, `micronaut` | Auto-détecté |
| `port` | Port exposé | Chaîne numérique | `"8080"` |
| `devMode` | Activer le live reload | `true`, `false` | `true` |
| `env` | Variables d'environnement | Map clé-valeur | `{}` |
| `watchPaths` | Chemins à surveiller | Liste de chemins | Auto-détecté |

## Services dépendants

Vous pouvez déclarer des services dépendants comme des bases de données ou des caches :

```yaml
services:
  # Votre application
  - name: app
    path: .
    java: "17"
    runtime: spring
    port: "8080"
    
  # Base de données MySQL
  - name: db
    type: mysql
    version: "8.0"
    port: "3306"
    env:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: mydb
```

### Types de services supportés

| Type | Versions | Options spécifiques |
|------|----------|---------------------|
| `mysql` | `5.7`, `8.0` | `MYSQL_ROOT_PASSWORD`, `MYSQL_DATABASE` |
| `postgres` | `12`, `13`, `14`, `15` | `POSTGRES_PASSWORD`, `POSTGRES_DB` |
| `mongodb` | `4`, `5`, `6` | `MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD` |
| `redis` | `6`, `7` | - |
| `kafka` | `2.8`, `3.0`, `3.5` | - |
| `rabbitmq` | `3.8`, `3.9`, `3.10` | - |
| `elasticsearch` | `7`, `8` | - |

## Variables d'environnement

Vous pouvez définir des variables d'environnement pour chaque service :

```yaml
services:
  - name: mon-service
    # ...autres paramètres...
    env:
      SPRING_PROFILES_ACTIVE: dev
      LOGGING_LEVEL_ROOT: INFO
      DATABASE_URL: jdbc:mysql://db:3306/mydb
```

## Configuration des volumes

Vous pouvez configurer des volumes pour la persistance des données :

```yaml
services:
  - name: db
    type: mysql
    version: "8.0"
    volumes:
      - mysql-data:/var/lib/mysql
```

## Exemples

### Projet Spring Boot minimal

```yaml
services:
  - name: spring-app
    path: .
    java: "17"
    build: maven
    runtime: spring
    port: "8080"
    devMode: true
```

### Projet multi-services complet

```yaml
services:
  # Service API
  - name: api-service
    path: ./services/api
    java: "17"
    build: maven
    runtime: spring
    port: "8080"
    devMode: true
    watchPaths:
      - src/main/java
      - src/main/resources
    env:
      SPRING_PROFILES_ACTIVE: dev
      LOGGING_LEVEL_ROOT: INFO

  # Service d'authentification
  - name: auth-service
    path: ./services/auth
    java: "17"
    build: maven
    runtime: quarkus
    port: "8081"
    devMode: true
    env:
      QUARKUS_PROFILE: dev

  # Base de données
  - name: db
    type: mysql
    version: "8.0"
    port: "3306"
    env:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: appdb
    volumes:
      - mysql-data:/var/lib/mysql

  # Cache
  - name: cache
    type: redis
    version: "6.2"
    port: "6379"
    volumes:
      - redis-data:/data
```

Pour des exemples plus détaillés, consultez le répertoire `examples` dans le dépôt.
