# Usage Guide

This guide provides detailed instructions on how to use Turbotilt effectively for your Java development projects.

## Table of Contents

- [Commands Overview](#commands-overview)
- [Initializing a Project](#initializing-a-project)
- [Starting Your Environment](#starting-your-environment)
- [Checking Your Environment](#checking-your-environment)
- [Stopping Your Environment](#stopping-your-environment)
- [Advanced Usage](#advanced-usage)

## Commands Overview

Turbotilt offers the following main commands:

| Command | Description |
|---------|-------------|
| `init`  | Initialize a project by creating Dockerfile, docker-compose.yml and Tiltfile |
| `up`    | Start the development environment using Tilt or Docker Compose |
| `doctor`| Check the environment and configuration, providing diagnostics |
| `stop`  | Stop the environment and clean up resources |
| `version`| Display the current version of Turbotilt |

## Initializing a Project

The `init` command analyzes your project, detects the Java framework, and creates the necessary files for your development environment.

### Basic Initialization

```bash
turbotilt init
```

This will:
1. Scan your project to detect the Java framework (Spring Boot, Quarkus, Micronaut)
2. Generate an appropriate Dockerfile
3. Create a docker-compose.yml file, including dependent services (if detected)
4. Generate a Tiltfile for live reload

### Options

```bash
# Specify the framework explicitly
turbotilt init --framework spring

# Set the application port
turbotilt init --port 8080

# Specify JDK version
turbotilt init --jdk 17

# Enable development mode (default)
turbotilt init --dev

# Generate a manifest file
turbotilt init --generate-manifest

# Initialize from an existing manifest
turbotilt init --from-manifest
```

## Starting Your Environment

The `up` command starts your development environment using Tilt (default) or Docker Compose.

### Basic Start

```bash
turbotilt up
```

This will:
1. Build your application using the generated Dockerfile
2. Start all services defined in docker-compose.yml
3. Configure live reload with Tilt
4. Show logs from all services

### Options

```bash
# Start without Tilt (Docker Compose only)
turbotilt up --tilt=false

# Start a specific service from the manifest
turbotilt up --service payment-service

# Enable debug mode with detailed logs
turbotilt up --debug
```

## Checking Your Environment

The `doctor` command checks your environment and configuration, helping you troubleshoot issues.

```bash
# Basic health check
turbotilt doctor

# Validate the manifest file
turbotilt doctor --validate-manifest

# Detailed check with verbose output
turbotilt doctor --verbose --log
```

The doctor command checks:
- Docker and Docker Compose installation and configuration
- Tilt installation for live reload
- JDK and Java environment
- Network configuration and permissions
- Manifest syntax and validity

## Stopping Your Environment

The `stop` command stops your environment and cleans up resources.

```bash
turbotilt stop
```

This will:
1. Stop all running containers
2. Remove temporary resources
3. Keep your configuration files intact

## Advanced Usage

### Global Flags

All commands accept these options:

- `--dry-run`: Simulate execution without making actual changes
- `--debug`: Enable debug mode with detailed output
- `--config-file`: Specify a custom configuration file path

### Auto-update of Tiltfiles

In developer mode, Turbotilt automatically monitors changes to your source files and updates Tiltfiles accordingly, ensuring that your changes are always taken into account.

### Multi-service Projects

For projects with multiple services, you can define all services in a manifest file (`turbotilt.yaml`). See the [Configuration Guide](./configuration.md) for more details.

```bash
# Start all services defined in the manifest
turbotilt up

# Start a specific service
turbotilt up --service user-service
```

### Working with Existing Projects

To integrate Turbotilt with an existing project:

1. Navigate to your project directory
2. Run `turbotilt init` to generate configuration files
3. Customize the generated files as needed
4. Start your environment with `turbotilt up`

For more detailed integration information, see the [Integration Guide](./integration.md).
