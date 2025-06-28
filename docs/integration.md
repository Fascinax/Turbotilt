# Integration Guide

This guide describes how to integrate your project with Turbotilt to optimize your development environment.

## Table of Contents

- [Getting Started](#getting-started)
- [Integration with Existing Projects](#integration-with-existing-projects)
- [Custom Dockerfiles](#custom-dockerfiles)
- [Custom Docker Compose Configuration](#custom-docker-compose-configuration)
- [Custom Tiltfile Configuration](#custom-tiltfile-configuration)
- [CI/CD Integration](#cicd-integration)

## Getting Started

### For New Projects

1. Create a new directory for your project
2. Initialize your Java project (Spring Boot, Quarkus, or Micronaut)
3. Run `turbotilt init` to generate the necessary configuration files
4. Start your development environment with `turbotilt up`

### For Existing Projects

1. Navigate to your project directory
2. Run `turbotilt init` to generate configuration files
3. Review and customize the generated files as needed
4. Start your development environment with `turbotilt up`

## Integration with Existing Projects

### Auto-detection

Turbotilt automatically detects:

- Java frameworks (Spring Boot, Quarkus, Micronaut) from `pom.xml` or `build.gradle`
- Database dependencies (MySQL, PostgreSQL, MongoDB, etc.)
- Cache dependencies (Redis)
- Messaging dependencies (Kafka, RabbitMQ)

### Manifest File

For more control, create a `turbotilt.yaml` file in your project root:

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

See the [Configuration Guide](./configuration.md) for all available options.

## Custom Dockerfiles

If you have a custom Dockerfile, Turbotilt can use it instead of generating one:

1. Create a `Dockerfile` in your project root
2. Run `turbotilt init --use-existing-dockerfile`

Turbotilt will detect your custom Dockerfile and use it for building your application.

Alternatively, you can specify the Dockerfile path in your `turbotilt.yaml`:

```yaml
services:
  - name: my-service
    path: .
    dockerfilePath: ./custom/Dockerfile
```

## Custom Docker Compose Configuration

For advanced Docker Compose configuration:

1. Create a `docker-compose.yml` file in your project root
2. Run `turbotilt init --use-existing-compose`

Turbotilt will detect your custom Docker Compose file and use it.

You can also specify the Docker Compose file path in your `turbotilt.yaml`:

```yaml
services:
  - name: my-service
    path: .
    composeFilePath: ./custom/docker-compose.yml
```

## Custom Tiltfile Configuration

For advanced Tilt configuration:

1. Create a `Tiltfile` in your project root
2. Run `turbotilt init --use-existing-tiltfile`

Turbotilt will detect your custom Tiltfile and use it.

You can also specify the Tiltfile path in your `turbotilt.yaml`:

```yaml
services:
  - name: my-service
    path: .
    tiltfilePath: ./custom/Tiltfile
```

## CI/CD Integration

Turbotilt can be integrated into your CI/CD pipeline:

### GitHub Actions Example

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
    
    - name: Install Turbotilt
      run: curl -fsSL https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.sh | bash
    
    - name: Validate configuration
      run: turbotilt doctor --validate-manifest
    
    - name: Initialize project
      run: turbotilt init --dry-run
    
    - name: Build Docker image
      run: turbotilt up --tilt=false --build-only
```

### Jenkins Pipeline Example

```groovy
pipeline {
    agent any
    
    stages {
        stage('Install Turbotilt') {
            steps {
                sh 'curl -fsSL https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.sh | bash'
            }
        }
        
        stage('Validate configuration') {
            steps {
                sh 'turbotilt doctor --validate-manifest'
            }
        }
        
        stage('Initialize project') {
            steps {
                sh 'turbotilt init --dry-run'
            }
        }
        
        stage('Build Docker image') {
            steps {
                sh 'turbotilt up --tilt=false --build-only'
            }
        }
    }
}
```

For more information about integrating Turbotilt with your CI/CD pipeline, please refer to your CI/CD platform's documentation.
