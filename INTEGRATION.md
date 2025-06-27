# Integration avec Turbotilt

Ce document décrit comment intégrer votre projet avec Turbotilt pour optimiser votre environnement de développement.

## Fichier turbotilt.yaml

Turbotilt supporte une configuration déclarative via un fichier `turbotilt.yaml` à la racine de votre projet.

### Structure de base

```yaml
services:
  - name: mon-service
    path: .
    java: "17"
    build: maven
    runtime: spring
    port: "8080"
    devMode: true
```

### Options disponibles

| Option | Description | Valeurs possibles | Défaut |
|--------|-------------|-------------------|--------|
| `name` | Nom du service | Chaîne de caractères | Nom du répertoire |
| `path` | Chemin relatif vers le service | Chemin | `.` |
| `java` | Version de Java | `"8"`, `"11"`, `"17"`, `"21"` | `"17"` |
| `build` | Système de build | `maven`, `gradle` | Auto-détecté |
| `runtime` | Framework Java | `spring`, `quarkus`, `micronaut` | Auto-détecté |
| `port` | Port exposé | Chaîne numérique | `"8080"` |
| `devMode` | Activer le live reload | `true`, `false` | `true` |
| `env` | Variables d'environnement | Map clé-valeur | `{}` |
| `watchPaths` | Chemins à surveiller | Liste de chemins | Auto-détecté |

## Services dépendants

Vous pouvez déclarer des services dépendants comme des bases de données ou des caches:

```yaml
services:
  # Votre application
  - name: app
    path: .
    java: "17"
    runtime: spring
    port: "8080"
    
  # Base de données MySQL
  - name: db
    type: mysql
    version: "8.0"
    port: "3306"
    env:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: mydb
```

### Types de services supportés

| Type | Versions | Options spécifiques |
|------|----------|---------------------|
| `mysql` | `5.7`, `8.0` | `MYSQL_ROOT_PASSWORD`, `MYSQL_DATABASE` |
| `postgres` | `12`, `13`, `14`, `15` | `POSTGRES_PASSWORD`, `POSTGRES_DB` |
| `mongodb` | `4`, `5`, `6` | `MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD` |
| `redis` | `6`, `7` | - |
| `kafka` | `2.8`, `3.0`, `3.5` | - |
| `rabbitmq` | `3.8`, `3.9`, `3.10` | - |
| `elasticsearch` | `7`, `8` | - |

## Exemples

Consultez les exemples dans le répertoire `examples/` pour des configurations typiques:

- `turbotilt.yaml` - Configuration multi-services complète
- `turbotilt-simple.yaml` - Configuration minimale pour un projet Spring Boot

## Commandes utiles

```bash
# Initialiser un projet avec auto-détection
turbotilt init

# Générer un manifeste à partir de l'auto-détection
turbotilt init --generate-manifest

# Démarrer l'environnement
turbotilt up

# Vérifier la configuration
turbotilt doctor --validate-manifest
```
