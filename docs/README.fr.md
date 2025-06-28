# 🛠️ Turbotilt

> **Initialisation, lancement et diagnostic d'environnements de développement** pour projets Java (Spring Boot, Quarkus, Micronaut...), avec support de Tilt pour le live reload.

![status-badge](https://img.shields.io/badge/status-beta-orange)
![version](https://img.shields.io/github/v/release/Fascinax/Turbotilt?include_prereleases)
![license](https://img.shields.io/github/license/Fascinax/Turbotilt)
![go-version](https://img.shields.io/github/go-mod/go-version/Fascinax/Turbotilt)
[![ci](https://github.com/Fascinax/Turbotilt/actions/workflows/ci.yml/badge.svg)](https://github.com/Fascinax/Turbotilt/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Fascinax/Turbotilt/branch/main/graph/badge.svg)](https://codecov.io/gh/Fascinax/Turbotilt)

*[English documentation](../README.md)*

---

## ✨ Fonctionnalités

| Fonctionnalité                                                            | Description                                                             |
| ------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| 🔍 **Détection automatique** des frameworks Java (Maven/Gradle)           | Analyse des fichiers `pom.xml` & `build.gradle` et de la structure      |
| �️ **Sélection multi-microservices** avec la commande `select`           | Détecte et sélectionne quels services exécuter dans des environnements complexes |
| 🧹 **Environnements temporaires** avec la commande `tup`                  | Génère les configurations, démarre les services et nettoie à la fin     |
| �🐳 **Génération dynamique** de Dockerfile adapté au framework détecté     | Crée un Dockerfile optimisé pour Spring, Quarkus ou Micronaut           |
| 🧩 **Docker Compose intégré** avec détection des services dépendants      | Détecte et configure automatiquement MySQL, PostgreSQL, Redis, etc.     |
| ⚡ **Support de Tilt** pour le live-reload                                | Génère un Tiltfile avec règles live-update adaptées au framework        |
| 🏥 **Diagnostic avancé** (doctor)                                         | Vérifie l'installation, l'environnement et génère un score de santé     |
| 🔧 **Configuration flexible**                                             | Configuration par fichier YAML et flags en ligne de commande             |
| 📝 **Manifeste déclaratif**                                              | Support de configuration multi-service avec validation de schéma           |

---

## 📦 Installation

Plusieurs méthodes d'installation sont disponibles :

### Homebrew (macOS et Linux)

```bash
brew tap VOTRE_UTILISATEUR/turbotilt
brew install turbotilt
```

### Script d'installation (macOS, Linux, Windows)

```bash
# macOS et Linux
curl -fsSL https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.sh | bash

# Windows PowerShell
iwr -useb https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.ps1 | iex
```

### Téléchargement direct

Téléchargez la dernière version depuis la [page des releases](https://github.com/Fascinax/turbotilt/releases).

---

## 🚀 Démarrage rapide

```bash
# Initialiser un projet (auto-détection du framework)
turbotilt init

# Démarrer l'environnement avec Tilt
turbotilt up

# Vérifier l'environnement et la configuration
turbotilt doctor

# Arrêter l'environnement et nettoyer
turbotilt stop
```

Pour des exemples d'utilisation plus détaillés, consultez le [Guide d'utilisation](./usage.fr.md).

---

## 📖 Documentation

- [Guide d'utilisation](./usage.fr.md) - Instructions d'utilisation détaillées et exemples
- [Configuration](./configuration.fr.md) - Options de configuration et format du manifeste
- [Intégration](./integration.fr.md) - Comment intégrer Turbotilt à votre projet
- [Frameworks et Services supportés](./supported.fr.md) - Liste des frameworks Java et services dépendants supportés
- [Guide Multi-Microservices](./guide-pratique-multiservices.fr.md) - Guide pour projets complexes multi-services
- [Avantages pour les Équipes](./avantages-equipes.fr.md) - Bénéfices pour les équipes de développement
- [Commande Select](./select-command.fr.md) - Documentation détaillée de la commande `select`
- [Commande Tup](./tup-command.fr.md) - Documentation de la commande `tup`
- [Contribution](../CONTRIBUTING.md) - Comment contribuer au projet
- [Tests](../TESTING.md) - Directives et procédures de test
- [Notes de version](../CHANGELOG-IMPROVEMENTS.md) - Derniers changements et améliorations

---

## 🤝 Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir des issues ou des pull requests.

## 📄 Licence

MIT
