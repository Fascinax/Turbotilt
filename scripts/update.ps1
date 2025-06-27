# Script PowerShell de mise à jour pour Turbotilt

Write-Host "🔄 Mise à jour de Turbotilt..." -ForegroundColor Cyan

# Déterminer l'architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ([System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture -eq [System.Runtime.InteropServices.Architecture]::Arm64) {
        "arm64"
    } else {
        "amd64"
    }
} else {
    Write-Host "❌ Architecture 32 bits non supportée" -ForegroundColor Red
    exit 1
}

# Obtenir la dernière version depuis l'API GitHub
Write-Host "🔍 Recherche de la dernière version..." -ForegroundColor Yellow
$ApiUrl = "https://api.github.com/repos/Fascinax/turbotilt/releases/latest"
$LatestVersionInfo = Invoke-RestMethod -Uri $ApiUrl -UseBasicParsing
$LatestVersion = $LatestVersionInfo.tag_name -replace "v", ""

if (-not $LatestVersion) {
    Write-Host "❌ Impossible de déterminer la dernière version" -ForegroundColor Red
    exit 1
}

Write-Host "✅ Dernière version trouvée: $LatestVersion" -ForegroundColor Green

# Déterminer le chemin d'installation actuel
$CurrentPath = $null
try {
    $CurrentPath = (Get-Command turbotilt -ErrorAction Stop).Source
} catch {
    $CurrentPath = Join-Path $env:LOCALAPPDATA "turbotilt\turbotilt.exe"
}

$InstallDir = Split-Path -Parent $CurrentPath

if (-not (Test-Path $InstallDir)) {
    Write-Host "📁 Création du répertoire d'installation $InstallDir..." -ForegroundColor Yellow
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# Télécharger la dernière version
$DownloadUrl = "https://github.com/Fascinax/turbotilt/releases/download/v${LatestVersion}/turbotilt-${LatestVersion}-windows-${Arch}.zip"
$TempDir = Join-Path $env:TEMP "turbotilt-update"
$TempFile = Join-Path $TempDir "turbotilt.zip"

if (Test-Path $TempDir) {
    Remove-Item $TempDir -Recurse -Force
}
New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

Write-Host "📥 Téléchargement depuis $DownloadUrl..." -ForegroundColor Blue
Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile -UseBasicParsing

# Extraire l'archive
Write-Host "📦 Extraction de l'archive..." -ForegroundColor Magenta
Expand-Archive -Path $TempFile -DestinationPath $TempDir -Force

# Installer le binaire
Write-Host "📋 Installation de Turbotilt dans $InstallDir..." -ForegroundColor Cyan
Copy-Item -Path "$TempDir\turbotilt.exe" -Destination $InstallDir -Force

# Nettoyer
Remove-Item $TempDir -Recurse -Force

Write-Host "✨ Turbotilt a été mis à jour vers la version $LatestVersion" -ForegroundColor Green
Write-Host "📍 Emplacement: $InstallDir\turbotilt.exe" -ForegroundColor White

# Vérifier l'installation
try {
    $Version = & "$InstallDir\turbotilt.exe" --version
    Write-Host "🚀 Installation réussie!" -ForegroundColor Green
    Write-Host "💡 Version installée: $Version" -ForegroundColor Cyan
} catch {
    Write-Host "⚠️ L'installation semble avoir réussi, mais il y a eu un problème lors de l'exécution du binaire" -ForegroundColor Yellow
    Write-Host "💡 Vérifiez manuellement avec: $InstallDir\turbotilt.exe --version" -ForegroundColor White
}

# Vérifier si le répertoire est dans le PATH
$EnvPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if (-not $EnvPath.Contains($InstallDir)) {
    Write-Host "⚠️ Le répertoire $InstallDir n'est pas dans votre PATH" -ForegroundColor Yellow
    Write-Host "💡 Vous pouvez l'ajouter en exécutant:" -ForegroundColor White
    Write-Host "[Environment]::SetEnvironmentVariable('PATH', `"$InstallDir;`" + [Environment]::GetEnvironmentVariable('PATH', 'User'), 'User')" -ForegroundColor Gray
}
