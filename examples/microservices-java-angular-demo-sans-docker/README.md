# Exemple de Microservices avec Java et Angular (Sans Docker)

Ce projet est un exemple d'architecture microservices utilisant différents frameworks Java pour le backend, Angular pour le frontend, et RabbitMQ pour la communication asynchrone. Cette version est conçue pour être utilisée directement avec Turbotilt sans configuration préalable.

## Architecture

L'architecture du projet est composée des éléments suivants :

- **User Service** : Service de gestion des utilisateurs développé avec Spring Boot
- **Product Service** : Service de gestion des produits développé avec Quarkus
- **Order Service** : Service de gestion des commandes développé avec Micronaut
- **Frontend** : Interface utilisateur développée avec Angular
- **RabbitMQ** : Broker de messages pour la communication asynchrone entre les services

## Structure des Services

### User Service (Spring Boot)
- Port : 8081
- Fonctionnalités : CRUD des utilisateurs, authentification
- Endpoint principal : `/api/users`

### Product Service (Quarkus)
- Port : 8082
- Fonctionnalités : CRUD des produits, gestion des stocks
- Endpoint principal : `/api/products`

### Order Service (Micronaut)
- Port : 8083
- Fonctionnalités : Création et gestion des commandes
- Endpoint principal : `/api/orders`

### Frontend (Angular)
- Port : 4200
- Pages principales : Produits, Utilisateurs, Commandes

## Communication entre Services

Les services communiquent entre eux de deux façons :

1. **Communication synchrone** : Appels HTTP REST entre les services
2. **Communication asynchrone** : Messages échangés via RabbitMQ

Exemples de flux de communication :
- Quand un utilisateur est créé/modifié, un message est envoyé sur RabbitMQ
- Quand une commande est créée, un message est envoyé pour mettre à jour les stocks de produits

## Démarrage avec Turbotilt

Pour démarrer l'ensemble des services avec Turbotilt (sans configuration préalable) :

```bash
cd microservices-java-angular-demo-sans-docker
turbotilt init
turbotilt up
```

Turbotilt s'occupera de :
1. Analyser la structure du projet
2. Détecter les frameworks utilisés (Spring Boot, Quarkus, Micronaut)
3. Générer les fichiers nécessaires à l'exécution
4. Démarrer tous les services

## Développement Manuel

Pour le développement manuel, chaque service peut être lancé individuellement :

```bash
# Service utilisateur
cd user-service
./mvnw spring-boot:run

# Service produit
cd product-service
./mvnw quarkus:dev

# Service commande
cd order-service
./mvnw mn:run

# Frontend
cd frontend
npm start
```

## Points d'intérêt de l'exemple

- **Différents frameworks Java** : Démonstration de Spring Boot, Quarkus et Micronaut
- **Communication asynchrone** : Utilisation de RabbitMQ entre les services
- **Frontend moderne** : Interface Angular avec routage et communication avec les APIs
- **Utilisation de Turbotilt sans configuration** : Démonstration de la détection automatique et de la génération de configuration
