#!/bin/bash

# Script pour démarrer tous les microservices avec Turbotilt
echo -e "\033[36mDémarrage de l'environnement de développement microservices avec Turbotilt...\033[0m"

# Vérifier si Turbotilt est installé
if ! command -v turbotilt &> /dev/null; then
    echo -e "\033[31mTurbotilt n'est pas installé ou n'est pas dans le PATH.\033[0m"
    echo -e "\033[33mVeuillez installer Turbotilt: https://github.com/Fascinax/Turbotilt\033[0m"
    exit 1
fi

# Afficher la version de Turbotilt
turbotiltVersion=$(turbotilt version)
echo -e "\033[32mTurbotilt détecté: $turbotiltVersion\033[0m"

# Répertoire du projet
projectDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Se positionner dans le répertoire du projet
cd "$projectDir"

# Lancer Turbotilt
echo -e "\033[36mLancement de tous les services...\033[0m"
turbotilt up

# Gérer l'interruption proprement
function cleanup {
    echo -e "\n\033[33mArrêt des services...\033[0m"
    turbotilt stop
    exit 0
}

# Capturer SIGINT (Ctrl+C) et SIGTERM
trap cleanup SIGINT SIGTERM

# Garder le script en vie
while true; do
    sleep 1
done
