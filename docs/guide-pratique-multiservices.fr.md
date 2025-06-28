# Guide d'utilisation de Turbotilt avec un projet multi-microservices

Ce guide vous explique comment utiliser Turbotilt pour gérer efficacement un environnement de développement complexe comportant plusieurs microservices de différentes technologies.

## Structure du projet d'exemple

Le projet d'exemple dans le dossier `complex-project` contient les services suivants :

| Service | Technologie | Description |
|---------|------------|-------------|
| user-service | Java Spring Boot | Gestion des utilisateurs |
| auth-service | Java Quarkus | Authentification et autorisation |
| payment-service | Java Micronaut | Traitement des paiements |
| frontend | Angular | Interface utilisateur |
| analytics-service | Python | Analyse de données |
| notification-service | Node.js | Gestion des notifications |
| api-gateway | Go | Passerelle API |

## Commandes Turbotilt principales

### Détecter et sélectionner des services

La commande `select` est parfaite pour les projets complexes comme celui-ci :

```bash
cd complex-project
turbotilt select
```

Cette commande va :
1. Scanner le répertoire à la recherche de microservices
2. Afficher une liste des services détectés
3. Vous permettre de choisir lesquels vous souhaitez lancer

### Options utiles

- Pour créer un fichier de configuration permanent :
  ```bash
  turbotilt select --create-config
  ```

- Pour lancer directement les services sélectionnés :
  ```bash
  turbotilt select --launch
  ```

- Pour combiner les deux options :
  ```bash
  turbotilt select --create-config --launch
  ```

- Pour spécifier un nom de fichier de configuration personnalisé :
  ```bash
  turbotilt select --output mon-environnement.yaml
  ```

## Scénarios d'utilisation courants

### 1. Développement frontend

Si vous travaillez uniquement sur le frontend :

```bash
turbotilt select --launch
```

Puis sélectionnez `frontend`, `user-service` et `auth-service`.

### 2. Travail sur le système de paiement

```bash
turbotilt select --launch
```

Puis sélectionnez `payment-service`, `user-service` et éventuellement `frontend`.

### 3. Création d'environnements spécialisés

Vous pouvez créer différents fichiers de configuration pour différents cas d'utilisation :

```bash
# Configuration pour le développement frontend
turbotilt select --output frontend-dev.yaml
# Sélectionnez frontend, user-service, auth-service

# Configuration pour le développement backend
turbotilt select --output backend-dev.yaml
# Sélectionnez user-service, auth-service, payment-service, etc.
```

Puis, utilisez ces configurations avec :

```bash
turbotilt up -f frontend-dev.yaml
# ou
turbotilt up -f backend-dev.yaml
```

## Utilisation temporaire (sans fichiers Docker)

La combinaison des commandes `select` et `tup` est particulièrement puissante :

```bash
# Sélectionner des services et créer une configuration
turbotilt select --output temp-config.yaml

# Lancer temporairement (les fichiers Docker sont générés puis supprimés)
turbotilt tup -f temp-config.yaml
```

Cette approche est idéale pour les équipes qui ne souhaitent pas conserver de fichiers Docker dans leur dépôt de code.

## Intégration avec des outils CI/CD

Pour l'intégration continue, vous pouvez utiliser :

```bash
# Générer une configuration pour tous les services
echo "all" | turbotilt select --output ci-config.yaml --create-config

# Lancer tous les services pour les tests
turbotilt up -f ci-config.yaml
```

## Conseils pratiques

- Utilisez `select` pour explorer un nouveau projet et comprendre sa structure
- Créez des configurations spécifiques pour différents aspects du développement
- Combinez avec `tup` pour un environnement propre sans fichiers générés
- Utilisez des noms significatifs pour vos fichiers de configuration
