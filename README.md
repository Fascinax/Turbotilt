# 🛠️ Turbotilt : CLI for cloud-native dev environments

> **Initializat## � Documentation

- [Usage Guide](./docs/usage.md) - Detailed usage instructions and examples
- [Configuration](./docs/configuration.md) - Configuration options and manifest format
- [Integration](./docs/integration.md) - How to integrate Turbotilt with your project
- [Supported Frameworks & Services](./docs/supported.md) - List of supported Java frameworks and dependent services
- [Multi-Microservice Guide](./docs/practical-guide-multiservices.md) - Guide for complex multi-service projects
- [Team Benefits](./docs/team-benefits.md) - Benefits for development teams
- [Select Command](./docs/select-command.md) - Detailed documentation for the `select` command
- [Tup Command](./docs/tup-command.md) - Documentation for the `tup` command
- [Contributing](./CONTRIBUTING.md) - How to contribute to the project
- [Testing](./TESTING.md) - Testing guidelines and procedures
- [Release Notes](./CHANGELOG-IMPROVEMENTS.md) - Latest changes and improvementsch and diagnosis of development environments** for Java projects (Spring Boot, Quarkus, Micronaut...), with Tilt support for live reload.

![status-badge](https://img.shields.io/badge/status-beta-orange)
![version](https://img.shields.io/github/v/release/Fascinax/Turbotilt?include_prereleases)
![license](https://img.shields.io/github/license/Fascinax/Turbotilt)
![go-version](https://img.shields.io/github/go-mod/go-version/Fascinax/Turbotilt)
[![ci](https://github.com/Fascinax/Turbotilt/actions/workflows/ci.yml/badge.svg)](https://github.com/Fascinax/Turbotilt/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Fascinax/Turbotilt/branch/main/graph/badge.svg)](https://codecov.io/gh/Fascinax/Turbotilt)

*[Documentation en français](./docs/README.fr.md)*

---

## ✨ Features

| Feature                                                                | Description                                                        |
| ---------------------------------------------------------------------- | ------------------------------------------------------------------ |
| 🔍 **Automatic detection** of Java frameworks (Maven/Gradle)           | Analysis of `pom.xml` & `build.gradle` files and project structure |
| �️ **Multi-microservice selection** with `select` command             | Detect and select which services to run in complex environments    |
| 🧹 **Temporary environments** with `tup` command                       | Generate configs, start services, and clean up when done           |
| �🐳 **Dynamic generation** of Dockerfile adapted to detected framework  | Creates optimized Dockerfile for Spring, Quarkus or Micronaut      |
| 🧩 **Integrated Docker Compose** with dependent services detection     | Automatically detects and configures MySQL, PostgreSQL, Redis, etc.|
| ⚡ **Tilt support** for live-reload                                    | Generates Tiltfile with live-update rules adapted to the framework |
| 🏥 **Advanced diagnostics** (doctor)                                   | Checks installation, environment and generates a health score      |
| 🔧 **Flexible configuration**                                          | Configuration via YAML file and command line flags                 |
| 📝 **Declarative manifest**                                            | Multi-service configuration support with schema validation         |

---

## 📦 Installation

Multiple installation methods are available:

### Homebrew (macOS and Linux)

```bash
brew tap YOUR_USERNAME/turbotilt
brew install turbotilt
```

### Installation script (macOS, Linux, Windows)

```bash
# macOS and Linux
curl -fsSL https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.sh | bash

# Windows PowerShell
iwr -useb https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.ps1 | iex
```

### Direct download

Download the latest version from the [releases page](https://github.com/Fascinax/turbotilt/releases).

## 🚀 Quick Start

```bash
# Initialize a project (auto-detection of framework)
turbotilt init

# Start the environment with Tilt
turbotilt up

# Check the environment and configuration
turbotilt doctor

# Stop the environment and clean up
turbotilt stop
```

For more detailed usage examples, see the [Usage Guide](./docs/usage.md).

## � Documentation

- [Usage Guide](./docs/usage.md) - Detailed usage instructions and examples
- [Configuration](./docs/configuration.md) - Configuration options and manifest format
- [Integration](./docs/integration.md) - How to integrate Turbotilt with your project
- [Supported Frameworks & Services](./docs/supported.md) - List of supported Java frameworks and dependent services
- [Contributing](./CONTRIBUTING.md) - How to contribute to the project
- [Testing](./TESTING.md) - Testing guidelines and procedures
- [Release Notes](./CHANGELOG-IMPROVEMENTS.md) - Latest changes and improvements

## 🤝 Contributing

Contributions are welcome! Please feel free to open issues or pull requests.

## 📄 License

MIT
