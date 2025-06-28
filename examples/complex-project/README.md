# Complex Project Example

This directory demonstrates a complex microservices project with multiple services. It's a perfect example to showcase Turbotilt's automatic service detection and the `select` command functionality.

## Project Structure

This project consists of the following microservices **without any Docker configuration files**:

- **User Service** (Java Spring Boot)
- **Authentication Service** (Java Quarkus)
- **Payment Service** (Java Micronaut)
- **Frontend** (Angular)
- **Analytics Service** (Python)
- **Notification Service** (Node.js)
- **API Gateway** (Go)

## Automatic Service Detection

One of the key features of Turbotilt is its ability to automatically detect different types of microservices without requiring any Docker configuration files. In this example:

- Java services are detected by their `pom.xml` files and framework-specific markers
- Angular frontend is detected by its `angular.json` file
- Python service is detected by its `setup.py` file
- Node.js service is detected by its `package.json` file
- Go service is detected by its `go.mod` file

## Using the `select` Command

The `select` command is particularly useful for complex projects like this one, where you might not need to run all services at once.

### Basic Selection

To scan this directory and select which services to launch:

```bash
cd complex-project
turbotilt select
```

You'll see a list of detected microservices and can choose which ones to include.

### Typical Use Cases

#### Frontend Development

If you're working on the frontend and only need the authentication and user services:

```bash
turbotilt select --launch
```

Then select the numbers for the Frontend, User Service, and Authentication Service when prompted.

#### Backend Development

If you're working on the backend services and don't need the frontend:

```bash
turbotilt select --create-config --output backend-services.yaml
```

Then select all the backend services. Later, you can start this environment with:

```bash
turbotilt up -f backend-services.yaml
```

#### Minimal Environment

To create a minimal environment for quick testing:

```bash
turbotilt select --launch
```

Then select only the essential services for your current task.

## Advanced Usage

You can combine `select` with other Turbotilt commands for more advanced workflows:

```bash
# Create a service selection and then run in temporary mode
turbotilt select --output my-selection.yaml
turbotilt tup -f my-selection.yaml
```

This approach lets you:

1. Make a careful selection of services
2. Run them in temporary mode (configs cleaned up when done)
3. Iterate quickly on your development
