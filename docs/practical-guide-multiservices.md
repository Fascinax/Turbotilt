# Turbotilt Usage Guide for Multi-microservice Projects

This guide explains how to use Turbotilt to efficiently manage a complex development environment with multiple microservices using different technologies.

## Example Project Structure

The example project in the `complex-project` folder contains the following services:

| Service | Technology | Description |
|---------|------------|-------------|
| user-service | Java Spring Boot | User management |
| auth-service | Java Quarkus | Authentication and authorization |
| payment-service | Java Micronaut | Payment processing |
| frontend | Angular | User interface |
| analytics-service | Python | Data analysis |
| notification-service | Node.js | Notification management |
| api-gateway | Go | API Gateway |

## Main Turbotilt Commands

### Detecting and Selecting Services

The `select` command is perfect for complex projects like this one:

```bash
cd complex-project
turbotilt select
```

This command will:
1. Scan the directory for microservices
2. Display a list of detected services
3. Allow you to choose which ones you want to launch

### Useful Options

- To create a permanent configuration file:
  ```bash
  turbotilt select --create-config
  ```

- To directly launch the selected services:
  ```bash
  turbotilt select --launch
  ```

- To combine both options:
  ```bash
  turbotilt select --create-config --launch
  ```

- To specify a custom configuration filename:
  ```bash
  turbotilt select --output my-environment.yaml
  ```

## Common Usage Scenarios

### 1. Frontend Development

If you're only working on the frontend:

```bash
turbotilt select --launch
```

Then select `frontend`, `user-service`, and `auth-service`.

### 2. Working on the Payment System

```bash
turbotilt select --launch
```

Then select `payment-service`, `user-service`, and optionally `frontend`.

### 3. Creating Specialized Environments

You can create different configuration files for different use cases:

```bash
# Configuration for frontend development
turbotilt select --output frontend-dev.yaml
# Select frontend, user-service, auth-service

# Configuration for backend development
turbotilt select --output backend-dev.yaml
# Select user-service, auth-service, payment-service, etc.
```

Then, use these configurations with:

```bash
turbotilt up -f frontend-dev.yaml
# or
turbotilt up -f backend-dev.yaml
```

## Temporary Usage (Without Docker Files)

Combining the `select` and `tup` commands is particularly powerful:

```bash
# Select services and create a configuration
turbotilt select --output temp-config.yaml

# Launch temporarily (Docker files are generated then deleted)
turbotilt tup -f temp-config.yaml
```

This approach is ideal for teams that don't want to keep Docker files in their code repository.

## Integration with CI/CD Tools

For continuous integration, you can use:

```bash
# Generate a configuration for all services
echo "all" | turbotilt select --output ci-config.yaml --create-config

# Launch all services for testing
turbotilt up -f ci-config.yaml
```

## Practical Tips

- Use `select` to explore a new project and understand its structure
- Create specific configurations for different aspects of development
- Combine with `tup` for a clean environment without generated files
- Use meaningful names for your configuration files
