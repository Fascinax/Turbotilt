# 🛠️ Turbotilt : CLI pour environnements dev cloud-native

> **Initialisation, lancement et diagnostic d'environnements de développement** pour projets Java (Spring Boot, Quarkus, Micronaut…), avec support de Tilt pour le live reload.

![status-badge](https://img.shields.io/badge/status-beta-orange)

---

## ✨ Fonctionnalités

| Fonctionnalité                                                            | Description                                                             |
| ------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| 🔍 **Détection automatique** des frameworks Java (Maven/Gradle)           | Analyse des fichiers `pom.xml` & `build.gradle` et de la structure      |
| 🐳 **Génération dynamique** de Dockerfile adapté au framework détecté     | Crée un Dockerfile optimisé pour Spring, Quarkus ou Micronaut           |
| 🧩 **Docker Compose intégré** avec détection des services dépendants      | Détecte et configure automatiquement MySQL, PostgreSQL, Redis, etc.     |
| ⚡ **Support de Tilt** pour le live-reload                                | Génère un Tiltfile avec règles live-update adaptées au framework        |
| 🏥 **Diagnostic avancé** (doctor)                                         | Vérifie l'installation, l'environnement et génère un score de santé     |
| 🔧 **Configuration flexible**                                             | Configuration par fichier YAML et flags en ligne de commande             |
| 📝 **Manifeste déclaratif**                                              | Support de configuration multi-service avec validation de schéma           |

---

## 📦 Structure du projet

```
turbotilt/
├── cmd/                  # point d'entrée CLI (cobra)
│   ├── root.go
│   ├── init.go          # turbotilt init
│   ├── up.go            # turbotilt up
│   ├── doctor.go        # turbotilt doctor
│   └── ...
├── internal/
│   ├── config/          # gestion de la configuration et validation du manifeste
│   ├── scan/            # détection frameworks et services
│   ├── render/          # génération des fichiers Docker, etc.
│   └── runtime/         # exécution avec Tilt ou Docker Compose
├── templates/           # templates pour Tiltfile et autres
│   └── Tiltfile.tmpl    # template de base pour Tiltfile
├── go.mod               # dépendances Go
└── README.md            # documentation
```

---

## 🚀 Démarrage rapide

### Installation

```bash
# Cloner le repository
git clone https://github.com/votre-nom/turbotilt.git
cd turbotilt

# Compiler l'application
go build

# Créer un lien symbolique (optionnel)
# sudo ln -s $(pwd)/turbotilt /usr/local/bin/turbotilt
```

### Utilisation

```bash
# Initialiser un projet (auto-détection du framework)
turbotilt init

# Options disponibles
turbotilt init --framework spring --port 8080 --jdk 17 --dev

# Générer un manifeste turbotilt.yaml (nouveau)
turbotilt init --generate-manifest

# Initialiser un projet à partir d'un manifeste existant
turbotilt init --from-manifest

# Démarrer l'environnement avec Tilt
turbotilt up

# Démarrer avec Docker Compose uniquement (sans Tilt)
turbotilt up --tilt=false

# Démarrer un service spécifique du manifeste
turbotilt up --service payment-service

# Vérifier l'environnement et la configuration
turbotilt doctor

# Valider la syntaxe du manifeste
turbotilt doctor --validate-manifest

# Vérification détaillée
turbotilt doctor --verbose --log

# Arrêter l'environnement et nettoyer
turbotilt stop
```

---

## 🔍 Comment ça fonctionne

1. **Phase de configuration** - Lecture du manifeste `turbotilt.yaml` si présent
2. **Phase de scan** - Détection du framework et des services dépendants (si non spécifiés dans le manifeste)
3. **Génération des fichiers** - Création de Dockerfile, docker-compose.yml et Tiltfile basés sur le manifeste ou l'auto-détection
4. **Démarrage de l'environnement** - Exécution via Tilt ou Docker Compose
5. **Surveillance du code** - Live reload avec Tilt (pour un développement rapide)

> La configuration déclarative du manifeste a toujours priorité sur les valeurs auto-détectées.

---

## ⚙️ Configuration

Turbotilt peut être configuré via:

1. **Flags en ligne de commande** - Pour les réglages rapides
2. **Fichier turbotilt.yml** - Pour une configuration persistante du projet

### Configuration déclarative (nouveauté)

Depuis la version 2.0, Turbotilt supporte une approche entièrement déclarative permettant de définir tous les services de votre projet dans un seul fichier manifeste.

Exemple de `turbotilt.yaml` multi-services :

```yaml
services:
  - name: user-service
    path: services/user
    java: 17
    build: maven
    runtime: spring
    port: 8081
  - name: payment-service
    path: services/payment
    java: 21
    runtime: quarkus
```

**Avantages du manifeste déclaratif :**
- Définition complète de l'environnement dans un seul fichier
- Prise en charge de multiples services avec configurations indépendantes
- Possibilité de surcharger les paramètres auto-détectés
- Validation stricte du schéma de configuration

> **Note :** La configuration du manifeste a toujours priorité sur l'auto-détection.

### Configuration simple (héritage)

Pour les projets simples, vous pouvez toujours utiliser un format plus concis :

```yaml
framework: spring
port: 8080
jdk: "17"
dev_mode: true
services:
  - type: mysql
    version: "8.0"
    credentials:
      username: root
      password: example
```

---

## 🧩 Frameworks supportés

- **Spring Boot** - Détection automatique via pom.xml ou build.gradle
- **Quarkus** - Support complet avec live-reload
- **Micronaut** - Support basique

## 🛢️ Services dépendants supportés

- **Bases de données** - MySQL, PostgreSQL, MongoDB
- **Cache** - Redis
- **Messagerie** - Kafka, RabbitMQ
- **Autres** - Elasticsearch

---

## 🤝 Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir des issues ou des pull requests.

## 📄 Licence

MIT
