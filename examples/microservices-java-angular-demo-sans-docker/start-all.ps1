# Script pour démarrer tous les microservices avec Turbotilt
Write-Host "Démarrage de l'environnement de développement microservices avec Turbotilt..." -ForegroundColor Cyan

# Vérifier si Turbotilt est installé
try {
    $turbotiltVersion = turbotilt version
    Write-Host "Turbotilt détecté: $turbotiltVersion" -ForegroundColor Green
} catch {
    Write-Host "Turbotilt n'est pas installé ou n'est pas dans le PATH." -ForegroundColor Red
    Write-Host "Veuillez installer Turbotilt: https://github.com/Fascinax/Turbotilt" -ForegroundColor Yellow
    exit 1
}

# Répertoire du projet
$projectDir = $PSScriptRoot

# Se positionner dans le répertoire du projet
Set-Location $projectDir

# Lancer Turbotilt
Write-Host "Lancement de tous les services..." -ForegroundColor Cyan
turbotilt up

# Capturer Ctrl+C pour un arrêt propre
try {
    while ($true) {
        Start-Sleep -Seconds 1
    }
} finally {
    # Arrêter les services lors de l'interruption
    Write-Host "`nArrêt des services..." -ForegroundColor Yellow
    turbotilt stop
}
