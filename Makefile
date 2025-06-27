# Makefile pour Turbotilt

# Variables
BINARY_NAME=turbotilt
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILDTIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X turbotilt/cmd.Version=$(VERSION) -X turbotilt/cmd.BuildTime=$(BUILDTIME) -X turbotilt/cmd.GitCommit=$(COMMIT)"

GO=$(shell where go 2>nul || which go)
GOTEST=$(GO) test
GOVET=$(GO) vet
GOBUILD=$(GO) build
GOFMT=$(GO) fmt
GOLINT=golangci-lint

# Build targets per platform
PLATFORMS=linux windows darwin
ARCHITECTURES=amd64 arm64

# Output directories
RELEASE_DIR=release
DIST_DIR=dist

# Cibles
.PHONY: all build clean test coverage vet lint release dist install docs examples $(PLATFORMS) $(ARCHITECTURES)

all: test build

build:
	@echo "Compilation de $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)$(if $(filter $(OS),Windows_NT),.exe,)

# Cross-platform builds
$(PLATFORMS):
	@echo "Building for $@..."
	-@mkdir -p $(DIST_DIR)/$@-amd64
	GOOS=$@ GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$@-amd64/$(BINARY_NAME)$(if $(filter windows,$@),.exe,)
	-@mkdir -p $(DIST_DIR)/$@-arm64
	GOOS=$@ GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$@-arm64/$(BINARY_NAME)$(if $(filter windows,$@),.exe,)

# Build for all platforms
dist: clean
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	@$(MAKE) $(PLATFORMS)

# Create release packages
release: dist
	@echo "Creating release packages..."
	@mkdir -p $(RELEASE_DIR)
	@for platform in $(PLATFORMS); do \
		for arch in amd64 arm64; do \
			echo "Packaging for $$platform-$$arch..."; \
			cd $(DIST_DIR)/$$platform-$$arch && \
			zip -q ../../$(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-$$platform-$$arch.zip $(BINARY_NAME)$(if $(filter windows,$$platform),.exe,) && \
			cd ../.. && \
			shasum -a 256 $(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-$$platform-$$arch.zip > $(RELEASE_DIR)/$(BINARY_NAME)-$(VERSION)-$$platform-$$arch.zip.sha256; \
		done \
	done
	@echo "Release packages and checksums created in $(RELEASE_DIR) directory"

test:
	@echo "Exécution des tests..."
	$(GOTEST) -v ./...

coverage:
	@echo "Génération de la couverture de tests..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

vet:
	@echo "Vérification du code avec go vet..."
	$(GOVET) ./...

lint:
	@echo "Linting du code (nécessite golangci-lint)..."
	golangci-lint run -c .golangci.yml

clean:
	@echo "Nettoyage des fichiers générés..."
	-@if exist $(BINARY_NAME) del /f $(BINARY_NAME)
	-@if exist $(BINARY_NAME).exe del /f $(BINARY_NAME).exe
	-@if exist coverage.out del /f coverage.out
	-@if exist coverage.html del /f coverage.html
	-@if exist $(DIST_DIR) rmdir /s /q $(DIST_DIR)
	-@if exist $(RELEASE_DIR) rmdir /s /q $(RELEASE_DIR)

run: build
	@echo "Exécution de $(BINARY_NAME)..."
	./$(BINARY_NAME)

## Règles d'installation
install-dev-deps:
	@echo "Installation des dépendances de développement..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install:
	@echo "Installation de $(BINARY_NAME)..."
	@mkdir -p $(HOME)/.local/bin
	@cp $(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
	@chmod +x $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "Installé dans $(HOME)/.local/bin/$(BINARY_NAME)"
	@echo "Assurez-vous que ce chemin est dans votre PATH"

docs:
	@echo "Génération de la documentation Go..."
	@mkdir -p docs
	godoc -http=:8080 -index

examples:
	@echo "Vérification des exemples..."
	@cd examples && $(GO) build -o /dev/null ./...

# Support pour les tests par package
test/%:
	$(GOTEST) -v ./$*

help:
	@echo "Cibles disponibles:"
	@echo "  make all        - Exécute les tests et compile le binaire"
	@echo "  make build      - Compile le binaire"
	@echo "  make test       - Exécute les tests"
	@echo "  make coverage   - Génère un rapport de couverture de tests"
	@echo "  make vet        - Exécute go vet pour vérifier le code"
	@echo "  make lint       - Exécute golangci-lint"
	@echo "  make clean      - Supprime les fichiers générés"
	@echo "  make run        - Compile et exécute le binaire"
	@echo "  make install-dev-deps - Installe les outils de développement"
	@echo "  make install    - Installe le binaire dans ~/.local/bin"
	@echo "  make docs       - Génère la documentation Go"
	@echo "  make examples    - Vérifie les exemples du projet"
