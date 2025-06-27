# Script PowerShell pour exécuter les tests
# Ce script exécute les tests package par package pour faciliter le débogage

Write-Host "Exécution des tests unitaires pour Turbotilt" -ForegroundColor Green

# Définir la fonction pour exécuter un test avec une présentation améliorée
function Run-Test {
    param (
        [string]$PackagePath
    )
    
    Write-Host ""
    Write-Host "=====================================================" -ForegroundColor Cyan
    Write-Host "Exécution des tests pour le package: $PackagePath" -ForegroundColor Cyan
    Write-Host "=====================================================" -ForegroundColor Cyan
    
    & go test $PackagePath
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Tests réussis pour $PackagePath" -ForegroundColor Green
    } else {
        Write-Host "❌ Échec des tests pour $PackagePath" -ForegroundColor Red
    }
}

# Exécuter les tests package par package
Run-Test "./internal/config"
Run-Test "./internal/scan"
Run-Test "./internal/logger" 
Run-Test "./internal/render"
Run-Test "./internal/runtime"

Write-Host ""
Write-Host "Exécution terminée" -ForegroundColor Green
