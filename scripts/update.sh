#!/bin/bash
# Script de mise à jour pour Turbotilt

set -e

echo "🔄 Mise à jour de Turbotilt..."

# Déterminer la plateforme et l'architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Architecture non supportée: $ARCH"
        exit 1
        ;;
esac

# Obtenir la dernière version depuis l'API GitHub
echo "🔍 Recherche de la dernière version..."
API_URL="https://api.github.com/repos/Fascinax/turbotilt/releases/latest"
LATEST_VERSION=$(curl -s $API_URL | grep -o '"tag_name": "v[^"]*"' | sed 's/"tag_name": "v//' | sed 's/"//')

if [ -z "$LATEST_VERSION" ]; then
    echo "❌ Impossible de déterminer la dernière version"
    exit 1
fi

echo "✅ Dernière version trouvée: $LATEST_VERSION"

# Déterminer le chemin d'installation actuel
CURRENT_PATH=$(which turbotilt 2>/dev/null || echo "$HOME/.local/bin/turbotilt")
INSTALL_DIR=$(dirname "$CURRENT_PATH")

if [ ! -d "$INSTALL_DIR" ]; then
    echo "📁 Création du répertoire d'installation $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR"
fi

# Télécharger la dernière version
DOWNLOAD_URL="https://github.com/Fascinax/turbotilt/releases/download/v${LATEST_VERSION}/turbotilt-${LATEST_VERSION}-${OS}-${ARCH}.zip"
TEMP_DIR=$(mktemp -d)
TEMP_FILE="$TEMP_DIR/turbotilt.zip"

echo "📥 Téléchargement depuis $DOWNLOAD_URL..."
curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"

# Extraire l'archive
echo "📦 Extraction de l'archive..."
unzip -o "$TEMP_FILE" -d "$TEMP_DIR"

# Installer le binaire
echo "📋 Installation de Turbotilt dans $INSTALL_DIR..."
cp "$TEMP_DIR/turbotilt" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/turbotilt"

# Nettoyer
rm -rf "$TEMP_DIR"

echo "✨ Turbotilt a été mis à jour vers la version $LATEST_VERSION"
echo "📍 Emplacement: $INSTALL_DIR/turbotilt"

# Vérifier l'installation
if command -v turbotilt >/dev/null 2>&1; then
    echo "🚀 Vous pouvez maintenant utiliser Turbotilt!"
    echo "💡 Exécutez 'turbotilt --version' pour vérifier l'installation"
else
    echo "⚠️ Le répertoire $INSTALL_DIR n'est pas dans votre PATH"
    echo "💡 Ajoutez-le à votre PATH ou utilisez le chemin complet: $INSTALL_DIR/turbotilt"
fi
