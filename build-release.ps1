# Script PowerShell pour générer des releases sur Windows
# Usage: .\build-release.ps1 [version]

# Configuration
$BinaryName = "turbotilt"
$Version = if ($args[0]) { $args[0] } else { "0.1.0" } # Utiliser l'argument ou la valeur par défaut
$ReleaseDir = "release"
$DistDir = "dist"
$Platforms = @("windows", "linux", "darwin")
$Architectures = @("amd64", "arm64")

# Afficher la version
Write-Host "🔍 Génération de la release pour $BinaryName version $Version" -ForegroundColor Cyan

# Nettoyage préalable
Write-Host "🧹 Nettoyage des fichiers générés..." -ForegroundColor Cyan
if (Test-Path $BinaryName) { Remove-Item $BinaryName -Force }
if (Test-Path "$BinaryName.exe") { Remove-Item "$BinaryName.exe" -Force }
if (Test-Path $DistDir) { Remove-Item $DistDir -Recurse -Force }
if (Test-Path $ReleaseDir) { Remove-Item $ReleaseDir -Recurse -Force }

# Création des répertoires
New-Item -ItemType Directory -Path $DistDir -Force | Out-Null
New-Item -ItemType Directory -Path $ReleaseDir -Force | Out-Null

$HashTable = @{}

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
        
        # Stocker le hash pour la mise à jour de la formule Homebrew
        $HashTable["$Platform-$Arch"] = $Hash.Hash
    }
}

Write-Host "✅ Release packages et checksums créés dans le répertoire $ReleaseDir" -ForegroundColor Green

# Afficher les checksums pour faciliter la mise à jour de la formule Homebrew
Write-Host "`n📋 SHA256 pour la formule Homebrew:" -ForegroundColor Magenta
foreach ($Platform in @("darwin", "linux")) {
    foreach ($Arch in @("arm64", "amd64")) {
        $Hash = $HashTable["$Platform-$Arch"]
        Write-Host "$Platform-$Arch : $Hash" -ForegroundColor Cyan
    }
}

# Mise à jour automatique de la formule Homebrew
$HomebrewPath = "scripts\homebrew\turbotilt.rb"
if (Test-Path $HomebrewPath) {
    Write-Host "`n📝 Mise à jour de la formule Homebrew..." -ForegroundColor Magenta
    
    $HomebrewContent = Get-Content $HomebrewPath -Raw
    
    # Mettre à jour la version
    $HomebrewContent = $HomebrewContent -replace 'version "[0-9\.]+', "version `"$Version"
    
    # Mettre à jour les checksums
    $HomebrewContent = $HomebrewContent -replace 'sha256 "[A-F0-9]+"  # darwin-arm64', "sha256 `"$($HashTable['darwin-arm64'])`"  # darwin-arm64"
    $HomebrewContent = $HomebrewContent -replace 'sha256 "[A-F0-9]+"  # darwin-amd64', "sha256 `"$($HashTable['darwin-amd64'])`"  # darwin-amd64"
    $HomebrewContent = $HomebrewContent -replace 'sha256 "[A-F0-9]+"  # linux-arm64', "sha256 `"$($HashTable['linux-arm64'])`"  # linux-arm64"
    $HomebrewContent = $HomebrewContent -replace 'sha256 "[A-F0-9]+"  # linux-amd64', "sha256 `"$($HashTable['linux-amd64'])`"  # linux-amd64"
    
    # Écrire le fichier mis à jour
    $HomebrewContent | Set-Content $HomebrewPath
    
    Write-Host "✅ Formule Homebrew mise à jour avec succès" -ForegroundColor Green
}

Write-Host "`n🎉 Build terminé avec succès!" -ForegroundColor Green
Write-Host "💡 Pour publier cette version, créez un tag git et utilisez GitHub Actions:" -ForegroundColor White
Write-Host "git tag v$Version" -ForegroundColor Gray
Write-Host "git push origin v$Version" -ForegroundColor Gray
