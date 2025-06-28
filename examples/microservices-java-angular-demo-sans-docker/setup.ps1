# Script PowerShell pour préparer l'environnement de développement

Write-Host "Configuration de l'environnement pour le projet Microservices Java-Angular Demo..." -ForegroundColor Cyan

# Vérifier Java
try {
    $javaVersion = java -version 2>&1
    if ($javaVersion -match "version") {
        Write-Host "Java est installé." -ForegroundColor Green
    } else {
        Write-Host "Java n'est pas correctement installé. Veuillez installer Java 17." -ForegroundColor Red
        exit
    }
} catch {
    Write-Host "Java n'est pas installé. Veuillez installer Java 17." -ForegroundColor Red
    exit
}

# Vérifier Maven
try {
    $mvnVersion = mvn --version 2>&1
    if ($mvnVersion -match "Apache Maven") {
        Write-Host "Maven est installé." -ForegroundColor Green
    } else {
        Write-Host "Maven n'est pas correctement installé. Veuillez installer Maven." -ForegroundColor Red
        exit
    }
} catch {
    Write-Host "Maven n'est pas installé. Veuillez installer Maven." -ForegroundColor Red
    exit
}

# Vérifier Node.js
try {
    $nodeVersion = node --version 2>&1
    if ($nodeVersion -match "v") {
        Write-Host "Node.js est installé." -ForegroundColor Green
    } else {
        Write-Host "Node.js n'est pas correctement installé. Veuillez installer Node.js." -ForegroundColor Red
        exit
    }
} catch {
    Write-Host "Node.js n'est pas installé. Veuillez installer Node.js." -ForegroundColor Red
    exit
}

# Vérifier npm
try {
    $npmVersion = npm --version 2>&1
    if ($npmVersion -match "\d+\.\d+\.\d+") {
        Write-Host "npm est installé." -ForegroundColor Green
    } else {
        Write-Host "npm n'est pas correctement installé. Veuillez installer npm." -ForegroundColor Red
        exit
    }
} catch {
    Write-Host "npm n'est pas installé. Veuillez installer npm." -ForegroundColor Red
    exit
}

Write-Host "Installation des dépendances pour le frontend..." -ForegroundColor Cyan
cd $PSScriptRoot\frontend
npm install

Write-Host "Compilation du frontend..." -ForegroundColor Cyan
npm run build

Write-Host "Préparation terminée! Utilisez le script start-services.ps1 pour démarrer tous les services." -ForegroundColor Green
