# PowerShell script d'installation pour Turbotilt
# Usage: iwr -useb https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.ps1 | iex

$ErrorActionPreference = 'Stop'

# Variables de configuration
$GithubRepo = "Fascinax/turbotilt"
$BinaryName = "turbotilt.exe"
$InstallDir = "$env:LOCALAPPDATA\Turbotilt"

# Couleurs pour les messages
function Write-Color {
    param (
        [Parameter(Position=0, Mandatory=$true)]
        [string]$Message,
        [Parameter(Position=1, Mandatory=$false)]
        [string]$Color = "White"
    )
    
    $prevForeground = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $Color
    Write-Output $Message
    $host.UI.RawUI.ForegroundColor = $prevForeground
}

Write-Color "Détection du système..." "Cyan"

# Détecter l'architecture
$Arch = "amd64"
if ([Environment]::Is64BitOperatingSystem -eq $false) {
    Write-Color "Architecture 32-bit détectée." "Yellow"
    $Arch = "386"
}
else {
    if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq [System.Runtime.InteropServices.Architecture]::Arm64) {
        Write-Color "Architecture ARM64 détectée." "Yellow"
        $Arch = "arm64"
    }
    else {
        Write-Color "Architecture AMD64 détectée." "Green"
    }
}

# Récupérer la dernière version depuis GitHub
Write-Color "Récupération de la dernière version..." "Cyan"
$apiUrl = "https://api.github.com/repos/$GithubRepo/releases/latest"

try {
    $releaseInfo = Invoke-RestMethod -Uri $apiUrl -UseBasicParsing
    $version = $releaseInfo.tag_name
    Write-Color "Dernière version: $version" "Green"
}
catch {
    Write-Color "Erreur lors de la récupération de la dernière version: $_" "Red"
    exit 1
}

# URL de téléchargement
$downloadUrl = "https://github.com/$GithubRepo/releases/download/$version/turbotilt-$($version.TrimStart('v'))-windows-$Arch.zip"
Write-Color "URL de téléchargement: $downloadUrl" "Cyan"

# Créer dossier temporaire
$tempDir = [System.IO.Path]::GetTempPath() + [System.IO.Path]::GetRandomFileName()
New-Item -ItemType Directory -Path $tempDir | Out-Null
$zipPath = Join-Path -Path $tempDir -ChildPath "turbotilt.zip"

# Télécharger le fichier
Write-Color "Téléchargement de Turbotilt..." "Cyan"
try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
}
catch {
    Write-Color "Erreur lors du téléchargement: $_" "Red"
    Remove-Item -Path $tempDir -Recurse -Force
    exit 1
}

# Extraire le contenu
Write-Color "Extraction du fichier ZIP..." "Cyan"
try {
    Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force
}
catch {
    Write-Color "Erreur lors de l'extraction: $_" "Red"
    Remove-Item -Path $tempDir -Recurse -Force
    exit 1
}

# Créer le dossier d'installation si nécessaire
if (-not (Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

# Copier le binaire dans le dossier d'installation
$binaryPath = Join-Path -Path $tempDir -ChildPath $BinaryName
$installPath = Join-Path -Path $InstallDir -ChildPath $BinaryName

Write-Color "Installation de Turbotilt dans $InstallDir..." "Cyan"
try {
    Copy-Item -Path $binaryPath -Destination $installPath -Force
}
catch {
    Write-Color "Erreur lors de l'installation: $_" "Red"
    Remove-Item -Path $tempDir -Recurse -Force
    exit 1
}

# Nettoyer les fichiers temporaires
Remove-Item -Path $tempDir -Recurse -Force

# Ajouter au PATH si ce n'est pas déjà fait
$userPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
if (-not $userPath.Contains($InstallDir)) {
    Write-Color "Ajout de Turbotilt au PATH utilisateur..." "Cyan"
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", [EnvironmentVariableTarget]::User)
    $env:Path += ";$InstallDir"
}

# Vérifier l'installation
Write-Color "Vérification de l'installation..." "Cyan"
try {
    $output = & "$installPath" --version
    Write-Color "✅ Turbotilt a été installé avec succès!" "Green"
    Write-Color "$output" "White"
}
catch {
    Write-Color "⚠️ Turbotilt a été installé dans $installPath mais il y a eu une erreur lors de la vérification." "Yellow"
}

Write-Color "`n🚀 Pour démarrer, exécutez: turbotilt init" "Green"
