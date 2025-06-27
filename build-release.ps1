# Script PowerShell pour générer des releases sur Windows
# Usage: .\build-release.ps1

# Configuration
$BinaryName = "turbotilt"
$Version = "0.1.0" # Ou utilisez git describe --tags --always --dirty
$ReleaseDir = "release"
$DistDir = "dist"
$Platforms = @("windows", "linux", "darwin")
$Architectures = @("amd64", "arm64")

# Nettoyage préalable
Write-Host "🧹 Nettoyage des fichiers générés..." -ForegroundColor Cyan
if (Test-Path $BinaryName) { Remove-Item $BinaryName -Force }
if (Test-Path "$BinaryName.exe") { Remove-Item "$BinaryName.exe" -Force }
if (Test-Path $DistDir) { Remove-Item $DistDir -Recurse -Force }
if (Test-Path $ReleaseDir) { Remove-Item $ReleaseDir -Recurse -Force }

# Création des répertoires
New-Item -ItemType Directory -Path $DistDir -Force | Out-Null
New-Item -ItemType Directory -Path $ReleaseDir -Force | Out-Null

# Cross-compilation pour toutes les plateformes
Write-Host "🔨 Compilation pour toutes les plateformes..." -ForegroundColor Green

foreach ($Platform in $Platforms) {
    foreach ($Arch in $Architectures) {
        $OutputDir = "$DistDir\$Platform-$Arch"
        $Extension = if ($Platform -eq "windows") { ".exe" } else { "" }
        $OutputFile = "$OutputDir\$BinaryName$Extension"
        
        Write-Host "📦 Compilation pour $Platform-$Arch..." -ForegroundColor Yellow
        New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
        
        $env:GOOS = $Platform
        $env:GOARCH = $Arch
        
        # Construction avec ldflags
        go build -ldflags "-X turbotilt/cmd.Version=$Version -X turbotilt/cmd.BuildTime=$(Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ') -X turbotilt/cmd.GitCommit=$(git rev-parse --short HEAD)" -o $OutputFile
        
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ Échec de la compilation pour $Platform-$Arch" -ForegroundColor Red
            continue
        }
        
        # Création du ZIP
        $ZipFile = "$ReleaseDir\$BinaryName-$Version-$Platform-$Arch.zip"
        Write-Host "📦 Création de l'archive $ZipFile..." -ForegroundColor Blue
        
        Compress-Archive -Path $OutputFile -DestinationPath $ZipFile -Force
        
        # Calcul et enregistrement du SHA256
        $Hash = Get-FileHash -Path $ZipFile -Algorithm SHA256
        $Hash.Hash | Out-File "$ZipFile.sha256"
    }
}

Write-Host "✅ Release packages et checksums créés dans le répertoire $ReleaseDir" -ForegroundColor Green

# Afficher les checksums pour faciliter la mise à jour de la formule Homebrew
Write-Host "`n📋 SHA256 pour la formule Homebrew:" -ForegroundColor Magenta
foreach ($Platform in @("darwin", "linux")) {
    foreach ($Arch in @("arm64", "amd64")) {
        $ZipFile = "$ReleaseDir\$BinaryName-$Version-$Platform-$Arch.zip"
        $Hash = Get-Content "$ZipFile.sha256"
        Write-Host "$Platform-$Arch : $Hash" -ForegroundColor Cyan
    }
}
