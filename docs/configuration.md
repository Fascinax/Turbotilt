# Configuration Guide

This guide provides detailed information about Turbotilt's configuration options and manifest format.

## Table of Contents

- [Configuration Methods](#configuration-methods)
- [Manifest File Format](#manifest-file-format)
- [Service Configuration](#service-configuration)
- [Dependent Services](#dependent-services)
- [Environment Variables](#environment-variables)
- [Volume Configuration](#volume-configuration)
- [Examples](#examples)

## Configuration Methods

Turbotilt can be configured in multiple ways:

1. **Command-line flags**: For quick settings
2. **Manifest file (`turbotilt.yaml`)**: For persistent project configuration
3. **Auto-detection**: For framework and service detection

The priority order is:
1. Command-line flags (highest priority)
2. Manifest file configuration
3. Auto-detected values (lowest priority)

## Manifest File Format

Turbotilt supports a declarative approach with a YAML manifest file that defines all services in your project.

### Basic Structure

```yaml
services:
  - name: my-service
    path: .
    java: "17"
    build: maven
    runtime: spring
    port: "8080"
    devMode: true
```

### Multi-service Example

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

## Service Configuration

Here are all available options for configuring your Java services:

| Option | Description | Possible Values | Default |
|--------|-------------|-----------------|---------|
| `name` | Service name | String | Directory name |
| `path` | Relative path to the service | Path | `.` |
| `java` | Java version | `"8"`, `"11"`, `"17"`, `"21"` | `"17"` |
| `build` | Build system | `maven`, `gradle` | Auto-detected |
| `runtime` | Java framework | `spring`, `quarkus`, `micronaut` | Auto-detected |
| `port` | Exposed port | Numeric string | `"8080"` |
| `devMode` | Enable live reload | `true`, `false` | `true` |
| `env` | Environment variables | Key-value map | `{}` |
| `watchPaths` | Paths to watch | List of paths | Auto-detected |

## Dependent Services

You can declare dependent services such as databases or caches:

```yaml
services:
  # Your application
  - name: app
    path: .
    java: "17"
    runtime: spring
    port: "8080"
    
  # MySQL database
  - name: db
    type: mysql
    version: "8.0"
    port: "3306"
    env:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: mydb
```

### Supported Service Types

| Type | Versions | Specific Options |
|------|----------|------------------|
| `mysql` | `5.7`, `8.0` | `MYSQL_ROOT_PASSWORD`, `MYSQL_DATABASE` |
| `postgres` | `12`, `13`, `14`, `15` | `POSTGRES_PASSWORD`, `POSTGRES_DB` |
| `mongodb` | `4`, `5`, `6` | `MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD` |
| `redis` | `6`, `7` | - |
| `kafka` | `2.8`, `3.0`, `3.5` | - |
| `rabbitmq` | `3.8`, `3.9`, `3.10` | - |
| `elasticsearch` | `7`, `8` | - |

## Environment Variables

You can set environment variables for each service:

```yaml
services:
  - name: my-service
    # ...other settings...
    env:
      SPRING_PROFILES_ACTIVE: dev
      LOGGING_LEVEL_ROOT: INFO
      DATABASE_URL: jdbc:mysql://db:3306/mydb
```

## Volume Configuration

You can configure volumes for data persistence:

```yaml
services:
  - name: db
    type: mysql
    version: "8.0"
    volumes:
      - mysql-data:/var/lib/mysql
```

## Examples

### Minimal Spring Boot Project

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

### Complete Multi-service Project

```yaml
services:
  # API Service
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

  # Auth Service
  - name: auth-service
    path: ./services/auth
    java: "17"
    build: maven
    runtime: quarkus
    port: "8081"
    devMode: true
    env:
      QUARKUS_PROFILE: dev

  # Database
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

For more detailed examples, check the `examples` directory in the repository.
