# Supported Frameworks & Services

This document lists all the Java frameworks and dependent services that are supported by Turbotilt.

## Table of Contents

- [Java Frameworks](#java-frameworks)
- [Dependent Services](#dependent-services)
- [Framework-specific Features](#framework-specific-features)
- [Service-specific Configuration](#service-specific-configuration)

## Java Frameworks

Turbotilt supports the following Java frameworks:

| Framework | Auto-detection | Live Reload | Build Systems | Versions |
|-----------|----------------|-------------|---------------|----------|
| Spring Boot | ✅ | ✅ | Maven, Gradle | 2.x, 3.x |
| Quarkus | ✅ | ✅ | Maven, Gradle | 2.x, 3.x |
| Micronaut | ✅ | ✅ | Maven, Gradle | 3.x, 4.x |

### Framework Detection

Frameworks are detected through:
- Analysis of `pom.xml` or `build.gradle` dependencies
- Presence of framework-specific configuration files
- Project structure analysis

## Dependent Services

Turbotilt supports the following dependent services:

### Databases

| Database | Versions | Auto-detection |
|----------|----------|----------------|
| MySQL | 5.7, 8.0 | ✅ |
| PostgreSQL | 12, 13, 14, 15 | ✅ |
| MongoDB | 4, 5, 6 | ✅ |

### Caching

| Service | Versions | Auto-detection |
|---------|----------|----------------|
| Redis | 6, 7 | ✅ |

### Messaging

| Service | Versions | Auto-detection |
|---------|----------|----------------|
| Kafka | 2.8, 3.0, 3.5 | ✅ |
| RabbitMQ | 3.8, 3.9, 3.10 | ✅ |

### Search

| Service | Versions | Auto-detection |
|---------|----------|----------------|
| Elasticsearch | 7, 8 | ✅ |

## Framework-specific Features

### Spring Boot

- Optimized multi-stage Dockerfile
- Dev mode with spring-boot-devtools
- Environment variables for configuration
- Hot reload with Tilt

```yaml
services:
  - name: spring-app
    runtime: spring
    port: "8080"
    env:
      SPRING_PROFILES_ACTIVE: dev
```

### Quarkus

- Optimized multi-stage Dockerfile
- Dev mode with Quarkus development mode
- Environment variables for configuration
- Hot reload with Tilt

```yaml
services:
  - name: quarkus-app
    runtime: quarkus
    port: "8080"
    env:
      QUARKUS_PROFILE: dev
```

### Micronaut

- Optimized multi-stage Dockerfile
- Dev mode with Micronaut development mode
- Environment variables for configuration
- Hot reload with Tilt

```yaml
services:
  - name: micronaut-app
    runtime: micronaut
    port: "8080"
    env:
      MICRONAUT_ENVIRONMENTS: dev
```

## Service-specific Configuration

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

For more detailed configuration options, refer to the [Configuration Guide](./configuration.md).
