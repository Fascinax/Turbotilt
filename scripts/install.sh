#!/bin/bash

# Script d'installation pour Turbotilt
# Usage: curl -fsSL https://raw.githubusercontent.com/Fascinax/turbotilt/main/scripts/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architecture
case "${ARCH}" in
    x86_64) ARCH="amd64" ;;
    aarch64 | arm64) ARCH="arm64" ;;
    *)
        echo -e "${RED}Architecture non supportée: ${ARCH}${NC}"
        exit 1
        ;;
esac

# Map OS
case "${OS}" in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    msys*|mingw*|cygwin*|windows*) 
        OS="windows"
        if [ "${ARCH}" = "arm64" ]; then
            echo -e "${RED}Windows ARM64 n'est pas pris en charge actuellement.${NC}"
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}Système d'exploitation non supporté: ${OS}${NC}"
        exit 1
        ;;
esac

echo -e "${BLUE}Détection: ${OS}-${ARCH}${NC}"

# Get latest release info from GitHub
echo -e "${YELLOW}Récupération des informations de la dernière version...${NC}"
GITHUB_REPO="Fascinax/turbotilt"
LATEST_RELEASE_URL="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
LATEST_VERSION=$(curl -sL ${LATEST_RELEASE_URL} | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Impossible de récupérer la dernière version. Veuillez vérifier votre connexion internet.${NC}"
    exit 1
fi

echo -e "${GREEN}Version: ${LATEST_VERSION}${NC}"

# Determine file extension
if [ "$OS" = "windows" ]; then
    FILE_EXT=".exe"
else
    FILE_EXT=""
fi

# Determine installation directory
if [ "$OS" = "windows" ]; then
    # Default to Program Files for Windows
    INSTALL_DIR="$HOME/turbotilt"
    mkdir -p "$INSTALL_DIR"
elif [ "$OS" = "darwin" ]; then
    INSTALL_DIR="/usr/local/bin"
else
    INSTALL_DIR="/usr/local/bin"
fi

# Download URL
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/turbotilt-${LATEST_VERSION#v}-${OS}-${ARCH}.zip"
TEMP_DIR=$(mktemp -d)
TEMP_FILE="${TEMP_DIR}/turbotilt.zip"

echo -e "${YELLOW}Téléchargement de Turbotilt depuis ${DOWNLOAD_URL}...${NC}"
if ! curl -L -o "${TEMP_FILE}" "${DOWNLOAD_URL}"; then
    echo -e "${RED}Échec du téléchargement.${NC}"
    rm -rf "${TEMP_DIR}"
    exit 1
fi

echo -e "${YELLOW}Extraction...${NC}"
unzip -q -o "${TEMP_FILE}" -d "${TEMP_DIR}"

BINARY="${TEMP_DIR}/turbotilt${FILE_EXT}"
if [ ! -f "${BINARY}" ]; then
    echo -e "${RED}Binary not found in the downloaded archive.${NC}"
    rm -rf "${TEMP_DIR}"
    exit 1
fi

echo -e "${YELLOW}Installation dans ${INSTALL_DIR}...${NC}"

# Ensure install directory is writable, use sudo if necessary
if [ -w "${INSTALL_DIR}" ]; then
    cp "${BINARY}" "${INSTALL_DIR}/turbotilt${FILE_EXT}"
    chmod +x "${INSTALL_DIR}/turbotilt${FILE_EXT}"
else
    echo -e "${YELLOW}Privilèges root nécessaires pour installer dans ${INSTALL_DIR}...${NC}"
    sudo cp "${BINARY}" "${INSTALL_DIR}/turbotilt${FILE_EXT}"
    sudo chmod +x "${INSTALL_DIR}/turbotilt${FILE_EXT}"
fi

# Clean up
rm -rf "${TEMP_DIR}"

# Check installation
if command -v turbotilt >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Turbotilt a été installé avec succès!${NC}"
    echo -e "${BLUE}Version installée:${NC}"
    turbotilt --version
else
    echo -e "${YELLOW}⚠️ Turbotilt installé dans ${INSTALL_DIR}/turbotilt${FILE_EXT} mais n'est pas dans votre PATH.${NC}"
    echo -e "${YELLOW}Ajoutez ${INSTALL_DIR} à votre PATH ou exécutez directement via ${INSTALL_DIR}/turbotilt${FILE_EXT}${NC}"
fi

echo -e "${GREEN}🚀 Pour démarrer, exécutez: turbotilt init${NC}"
