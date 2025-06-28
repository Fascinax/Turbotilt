# Selecting Microservices with Turbotilt

The `select` command allows you to scan a directory for microservices and select which ones to launch. This is particularly useful for projects with multiple microservices where you may only want to work with a subset of them.

## Basic Usage

```bash
turbotilt select [directory]
```

If no directory is specified, the current directory will be scanned.

## How It Works

1. The command scans the specified directory for microservices by looking for project files like `pom.xml`, `build.gradle`, `package.json`, etc.
2. It displays a list of detected microservices, each with an index number.
3. You select which services to include by entering the corresponding numbers (comma-separated).
4. Depending on the flags you used, Turbotilt can:
   - Create a permanent configuration file (`turbotilt.yaml`)
   - Launch the selected services immediately
   - Both of the above

## Options

- `-o, --output <filename>`: Specify a name for the generated configuration file (default: `turbotilt.yaml`)
- `-c, --create-config`: Create a `turbotilt.yaml` file with the selected services
- `-l, --launch`: Launch the selected services after selection

## Examples

### Scan and Select Only

```bash
# Scan current directory and select services
turbotilt select

# Scan a specific directory and select services
turbotilt select ./my-project
```

### Create a Configuration File

```bash
# Select services and create a turbotilt.yaml file
turbotilt select -c

# Select services and create a custom named configuration file
turbotilt select --output my-custom-config.yaml
```

### Select and Launch

```bash
# Select services and launch them immediately
turbotilt select -l

# Select services, create a config file, and launch them
turbotilt select -c -l
```

## Use Cases

### Temporary Development Environment

When you want to quickly start a subset of services without creating permanent configuration:

```bash
turbotilt select -l
```

This will scan for services, let you select which ones to include, and launch them without creating a permanent configuration file.

### Creating Multiple Configuration Files

For a project with many microservices, you might want to create different configurations for different development scenarios:

```bash
# Create a configuration for the backend services
turbotilt select --output backend-services.yaml

# Create a configuration for the frontend services
turbotilt select --output frontend-services.yaml
```

Later, you can use these configurations with the `up` command:

```bash
turbotilt up -f backend-services.yaml
```

### Onboarding New Team Members

The `select` command provides a simple way for new team members to understand the architecture of a complex project:

1. Run `turbotilt select` to see all available microservices
2. Select specific services to understand how they work together
3. Save the configuration for future use with `--create-config`
