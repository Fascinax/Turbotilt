#!/usr/bin/env bash

# Script pour mettre à jour la formule Homebrew avec les bons checksums
# Utilisation: ./update-homebrew.sh <version>

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
  echo "Usage: ./update-homebrew.sh <version>"
  echo "Exemple: ./update-homebrew.sh 0.1.0"
  exit 1
fi

# Télécharger toutes les versions
TEMP_DIR=$(mktemp -d)
cd $TEMP_DIR

echo "Téléchargement des archives..."
curl -sL "https://github.com/VOTRE_UTILISATEUR/turbotilt/releases/download/v${VERSION}/turbotilt-${VERSION}-darwin-amd64.zip" -o darwin-amd64.zip
curl -sL "https://github.com/VOTRE_UTILISATEUR/turbotilt/releases/download/v${VERSION}/turbotilt-${VERSION}-darwin-arm64.zip" -o darwin-arm64.zip
curl -sL "https://github.com/VOTRE_UTILISATEUR/turbotilt/releases/download/v${VERSION}/turbotilt-${VERSION}-linux-amd64.zip" -o linux-amd64.zip
curl -sL "https://github.com/VOTRE_UTILISATEUR/turbotilt/releases/download/v${VERSION}/turbotilt-${VERSION}-linux-arm64.zip" -o linux-arm64.zip

# Calculer les checksums
DARWIN_AMD64_SHA=$(shasum -a 256 darwin-amd64.zip | awk '{print $1}')
DARWIN_ARM64_SHA=$(shasum -a 256 darwin-arm64.zip | awk '{print $1}')
LINUX_AMD64_SHA=$(shasum -a 256 linux-amd64.zip | awk '{print $1}')
LINUX_ARM64_SHA=$(shasum -a 256 linux-arm64.zip | awk '{print $1}')

cd -
rm -rf $TEMP_DIR

FORMULA_PATH="scripts/homebrew/turbotilt.rb"

# Mettre à jour la formule
sed -i -e "s/version \".*\"/version \"${VERSION}\"/" $FORMULA_PATH
sed -i -e "s/sha256 \"REPLACE_WITH_ACTUAL_SHA256\".*# darwin-arm64/sha256 \"${DARWIN_ARM64_SHA}\"  # darwin-arm64/" $FORMULA_PATH
sed -i -e "s/sha256 \"REPLACE_WITH_ACTUAL_SHA256\".*# darwin-amd64/sha256 \"${DARWIN_AMD64_SHA}\"  # darwin-amd64/" $FORMULA_PATH
sed -i -e "s/sha256 \"REPLACE_WITH_ACTUAL_SHA256\".*# linux-arm64/sha256 \"${LINUX_ARM64_SHA}\"  # linux-arm64/" $FORMULA_PATH
sed -i -e "s/sha256 \"REPLACE_WITH_ACTUAL_SHA256\".*# linux-amd64/sha256 \"${LINUX_AMD64_SHA}\"  # linux-amd64/" $FORMULA_PATH

echo "Formule Homebrew mise à jour avec la version ${VERSION} et checksums"
echo "N'oubliez pas de commit et push les changements dans le tap Homebrew"
