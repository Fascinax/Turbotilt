# Exemple de projet complexe

Ce répertoire démontre un projet complexe de microservices avec plusieurs services. C'est un exemple parfait pour illustrer la détection automatique de services et la fonctionnalité de la commande `select` de Turbotilt.

## Structure du projet

Ce projet se compose des microservices suivants **sans aucun fichier de configuration Docker** :

- **Service Utilisateur** (Java Spring Boot)
- **Service d'Authentification** (Java Quarkus)
- **Service de Paiement** (Java Micronaut)
- **Frontend** (Angular)
- **Service d'Analytique** (Python)
- **Service de Notification** (Node.js)
- **Passerelle API** (Go)

## Détection automatique des services

L'une des fonctionnalités clés de Turbotilt est sa capacité à détecter automatiquement différents types de microservices sans nécessiter de fichiers de configuration Docker. Dans cet exemple :

- Les services Java sont détectés par leurs fichiers `pom.xml` et les marqueurs spécifiques aux frameworks
- Le frontend Angular est détecté par son fichier `angular.json`
- Le service Python est détecté par son fichier `setup.py`
- Le service Node.js est détecté par son fichier `package.json`
- Le service Go est détecté par son fichier `go.mod`

## Utilisation de la commande `select`

La commande `select` est particulièrement utile pour des projets complexes comme celui-ci, où vous n'avez pas forcément besoin d'exécuter tous les services en même temps.

### Sélection basique

Pour scanner ce répertoire et sélectionner quels services lancer :

```bash
cd complex-project
turbotilt select
```

Vous verrez une liste des microservices détectés et pourrez choisir lesquels inclure.

### Cas d'utilisation typiques

#### Développement Frontend

Si vous travaillez sur le frontend et avez seulement besoin des services d'authentification et d'utilisateur :

```bash
turbotilt select --launch
```

Puis sélectionnez les numéros pour le Frontend, le Service Utilisateur et le Service d'Authentification lorsque vous y êtes invité.

#### Développement Backend

Si vous travaillez sur les services backend et n'avez pas besoin du frontend :

```bash
turbotilt select --create-config --output services-backend.yaml
```

Puis sélectionnez tous les services backend. Plus tard, vous pourrez démarrer cet environnement avec :

```bash
turbotilt up -f services-backend.yaml
```

#### Environnement minimal

Pour créer un environnement minimal pour des tests rapides :

```bash
turbotilt select --launch
```

Puis sélectionnez uniquement les services essentiels pour votre tâche actuelle.

## Utilisation avancée

Vous pouvez combiner `select` avec d'autres commandes Turbotilt pour des flux de travail plus avancés :

```bash
# Créer une sélection de services puis exécuter en mode temporaire
turbotilt select --output ma-selection.yaml
turbotilt tup -f ma-selection.yaml
```

Cette approche vous permet de :

1. Faire une sélection précise des services
2. Les exécuter en mode temporaire (configurations nettoyées à la fin)
3. Itérer rapidement sur votre développement
