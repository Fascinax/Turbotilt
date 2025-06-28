# Frameworks et Services supportés

Ce document répertorie tous les frameworks Java et services dépendants qui sont supportés par Turbotilt.

## Table des matières

- [Frameworks Java](#frameworks-java)
- [Services dépendants](#services-dépendants)
- [Fonctionnalités spécifiques aux frameworks](#fonctionnalités-spécifiques-aux-frameworks)
- [Configuration spécifique aux services](#configuration-spécifique-aux-services)

## Frameworks Java

Turbotilt supporte les frameworks Java suivants :

| Framework | Auto-détection | Live Reload | Systèmes de build | Versions |
|-----------|----------------|-------------|-------------------|----------|
| Spring Boot | ✅ | ✅ | Maven, Gradle | 2.x, 3.x |
| Quarkus | ✅ | ✅ | Maven, Gradle | 2.x, 3.x |
| Micronaut | ✅ | ✅ | Maven, Gradle | 3.x, 4.x |

### Détection des frameworks

Les frameworks sont détectés par :
- Analyse des dépendances `pom.xml` ou `build.gradle`
- Présence de fichiers de configuration spécifiques au framework
- Analyse de la structure du projet

## Services dépendants

Turbotilt supporte les services dépendants suivants :

### Bases de données

| Base de données | Versions | Auto-détection |
|----------------|----------|----------------|
| MySQL | 5.7, 8.0 | ✅ |
| PostgreSQL | 12, 13, 14, 15 | ✅ |
| MongoDB | 4, 5, 6 | ✅ |

### Mise en cache

| Service | Versions | Auto-détection |
|---------|----------|----------------|
| Redis | 6, 7 | ✅ |

### Messagerie

| Service | Versions | Auto-détection |
|---------|----------|----------------|
| Kafka | 2.8, 3.0, 3.5 | ✅ |
| RabbitMQ | 3.8, 3.9, 3.10 | ✅ |

### Recherche

| Service | Versions | Auto-détection |
|---------|----------|----------------|
| Elasticsearch | 7, 8 | ✅ |

## Fonctionnalités spécifiques aux frameworks

### Spring Boot

- Dockerfile multi-étapes optimisé
- Mode développement avec spring-boot-devtools
- Variables d'environnement pour la configuration
- Rechargement à chaud avec Tilt

```yaml
services:
  - name: spring-app
    runtime: spring
    port: "8080"
    env:
      SPRING_PROFILES_ACTIVE: dev
```

### Quarkus

- Dockerfile multi-étapes optimisé
- Mode développement avec le mode développement de Quarkus
- Variables d'environnement pour la configuration
- Rechargement à chaud avec Tilt

```yaml
services:
  - name: quarkus-app
    runtime: quarkus
    port: "8080"
    env:
      QUARKUS_PROFILE: dev
```

### Micronaut

- Dockerfile multi-étapes optimisé
- Mode développement avec le mode développement de Micronaut
- Variables d'environnement pour la configuration
- Rechargement à chaud avec Tilt

```yaml
services:
  - name: micronaut-app
    runtime: micronaut
    port: "8080"
    env:
      MICRONAUT_ENVIRONMENTS: dev
```

## Configuration spécifique aux services

### MySQL

```yaml
services:
  - name: db
    type: mysql
    version: "8.0"
    port: "3306"
    env:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: mydb
      MYSQL_USER: user
      MYSQL_PASSWORD: password
```

### PostgreSQL

```yaml
services:
  - name: db
    type: postgres
    version: "14"
    port: "5432"
    env:
      POSTGRES_PASSWORD: password
      POSTGRES_DB: mydb
      POSTGRES_USER: user
```

### MongoDB

```yaml
services:
  - name: db
    type: mongodb
    version: "6"
    port: "27017"
    env:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: password
```

### Redis

```yaml
services:
  - name: cache
    type: redis
    version: "7"
    port: "6379"
```

### Kafka

```yaml
services:
  - name: kafka
    type: kafka
    version: "3.5"
    port: "9092"
```

### RabbitMQ

```yaml
services:
  - name: rabbitmq
    type: rabbitmq
    version: "3.10"
    port: "5672"
    env:
      RABBITMQ_DEFAULT_USER: user
      RABBITMQ_DEFAULT_PASS: password
```

### Elasticsearch

```yaml
services:
  - name: elasticsearch
    type: elasticsearch
    version: "8"
    port: "9200"
    env:
      discovery.type: single-node
```

Pour des options de configuration plus détaillées, référez-vous au [Guide de configuration](./configuration.fr.md).
