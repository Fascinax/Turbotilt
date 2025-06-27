# Makefile pour Turbotilt

# Variables
BINARY_NAME=turbotilt.exe
GO=$(shell where go)
GOTEST=$(GO) test
GOVET=$(GO) vet
GOBUILD=$(GO) build

# Cibles
.PHONY: all build clean test coverage vet lint

all: test build

build:
	@echo "Compilation de $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) 

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
	del /f $(BINARY_NAME)
	del /f coverage.out coverage.html

run: build
	@echo "Exécution de $(BINARY_NAME)..."
	./$(BINARY_NAME)

## Règles d'installation
install-dev-deps:
	@echo "Installation des dépendances de développement..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

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
