# Turbotilt Temporary Up (tup) Command

The `tup` command (Temporary Up) provides a way to generate configuration files, start services, and automatically clean up when done. This document explains how to use this command and its various use cases.

## Table of Contents

- [Command Overview](#command-overview)
- [Basic Usage](#basic-usage)
- [Command Options](#command-options)
- [Use Cases](#use-cases)
  - [Temporary Development Environments](#temporary-development-environments)
  - [Evaluating Turbotilt on Existing Projects](#evaluating-turbotilt-on-existing-projects)
  - [Multi-Framework Microservices](#multi-framework-microservices)
  - [CI/CD Pipelines](#cicd-pipelines)
  - [Team Collaboration](#team-collaboration)
- [Example Architectures](#example-architectures)
- [Best Practices](#best-practices)

## Command Overview

The `tup` command combines the functionality of `init`, `up`, and `stop` with automatic cleanup of generated configuration files. It's designed for scenarios where you want to quickly start a development environment without permanently adding configuration files to your project.

## Basic Usage

```bash
# Start a temporary development environment
turbotilt tup

# Start a specific service in a temporary environment
turbotilt tup --service my-service

# Start in detached mode (background)
turbotilt tup --detached
```

When you run `turbotilt tup`:

1. Configuration files are generated (Dockerfile, docker-compose.yml, Tiltfile)
2. Services are started (with live reload by default)
3. When you press Ctrl+C, services are stopped and configuration files are removed

## Command Options

| Option | Description |
|--------|-------------|
| `--tilt` | Use Tilt for live reload (default: true) |
| `--detached` | Run in detached mode (background) |
| `--service [name]` | Run only a specific service |
| `--debug` | Enable debug mode with verbose output |
| `--dry-run` | Simulate execution without making changes |

## Use Cases

### Temporary Development Environments

**Scenario**: You need to quickly set up a development environment without adding configuration files to your project.

**Solution**: Use `turbotilt tup` to generate temporary files, start the environment, and clean up when done.

```bash
cd my-java-project
turbotilt tup
# Work with hot reload enabled
# Press Ctrl+C when done
```

**Benefits**:
- No configuration files cluttering your workspace
- No need to add files to `.gitignore`
- Clean working directory for version control

### Evaluating Turbotilt on Existing Projects

**Scenario**: You want to see how Turbotilt would configure your project without making permanent changes.

```bash
cd existing-project
turbotilt tup --detached
# Check how the services are running
turbotilt stop
```

**Benefits**:
- Risk-free evaluation
- No conflicts with existing configuration
- Test before committing to Turbotilt

### Multi-Framework Microservices

**Scenario**: You have a complex microservices architecture using different Java frameworks (Spring Boot, Quarkus, Micronaut, etc.).

```bash
cd microservices-project
turbotilt tup
```

**Benefits**:
- Auto-detection of different frameworks
- Unified development environment
- Standardized developer experience across teams

### CI/CD Pipelines

**Scenario**: You need to start services temporarily for running tests in a CI/CD pipeline.

```bash
# In CI script
turbotilt tup --detached
sleep 10  # Wait for services to start
./run-integration-tests.sh
turbotilt stop
```

**Benefits**:
- Clean test environment for each CI run
- No configuration management required
- Isolated test environments

### Team Collaboration

**Scenario**: Team members have different configuration preferences or setup requirements.

```bash
git pull  # Get latest code
turbotilt tup  # Start with auto-detected configuration
```

**Benefits**:
- No configuration file conflicts
- Consistent developer experience
- Each developer gets a fresh environment

## Example Architectures

The `tup` command is particularly useful for these architectural patterns:

### Monorepo with Multiple Projects

```
monorepo/
├── service-a/   # Spring Boot service
├── service-b/   # Quarkus service
└── frontend/    # Angular frontend
```

Each service can be detected and configured appropriately, even with different frameworks.

### Polyglot Projects

```
polyglot-app/
├── api/         # Java API
├── workers/     # Python processing
└── dashboard/   # Node.js frontend
```

Turbotilt can generate appropriate configurations for the Java components while working alongside other language services.

### Microservices with Shared Dependencies

```
project/
├── user-service/     # Spring Boot
├── order-service/    # Micronaut
├── product-service/  # Quarkus
└── shared-lib/       # Common code
```

The `tup` command can detect relationships between services and configure them to work together.

## Best Practices

1. **Initial Testing**: Use `tup` for initial testing before committing to permanent configuration files.

2. **CI Integration**: Incorporate `tup` into CI pipelines for clean test environments.

3. **Demo Environments**: Use for creating quick demo environments without configuration overhead.

4. **Onboarding**: Simplify onboarding by letting new team members start with `tup` instead of manual configuration.

5. **Legacy Projects**: Test how Turbotilt works with legacy projects without modifying them.
