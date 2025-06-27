# Guide de contribution à Turbotilt

Nous sommes ravis que vous envisagiez de contribuer à Turbotilt ! Ce document vous guidera à travers le processus de contribution.

## Code de conduite

Nous attendons de tous les contributeurs qu'ils respectent notre code de conduite. Veuillez être respectueux et professionnel dans toutes vos interactions.

## Comment contribuer

### Soumettre des problèmes (issues)

Si vous avez trouvé un bug ou avez une idée d'amélioration :

1. Vérifiez d'abord si le problème a déjà été signalé
2. Utilisez le modèle d'issue fourni
3. Incluez autant de détails que possible:
   - Version de Turbotilt
   - Système d'exploitation
   - Étapes pour reproduire le problème
   - Comportement attendu vs comportement observé

### Soumettre des modifications (pull requests)

1. Forkez le dépôt
2. Créez une branche pour votre fonctionnalité (`git checkout -b feature/amazing-feature`)
3. Commitez vos changements (`git commit -m 'feat: ajout d'une fonctionnalité incroyable'`)
4. Poussez vers votre branche (`git push origin feature/amazing-feature`)
5. Ouvrez une Pull Request

### Conventions de commit

Nous utilisons [Conventional Commits](https://www.conventionalcommits.org/) pour nos messages de commit:

- `feat:` pour les nouvelles fonctionnalités
- `fix:` pour les corrections de bugs
- `docs:` pour les changements de documentation
- `style:` pour les changements de formatage
- `refactor:` pour les changements de code qui n'ajoutent pas de fonctionnalités ou ne corrigent pas de bugs
- `test:` pour ajouter des tests manquants ou corriger des tests existants
- `chore:` pour les tâches de maintenance

### Processus de développement

1. Assurez-vous que tous les tests passent avant de soumettre votre PR
2. Ajoutez des tests pour les nouvelles fonctionnalités
3. Mettez à jour la documentation si nécessaire
4. Votre PR doit être revue par au moins un mainteneur avant d'être fusionnée

## Configuration de l'environnement de développement

### Prérequis

- Go 1.23 ou supérieur
- Docker et Docker Compose
- Tilt (pour tester les fonctionnalités de live reload)
- Make

### Configuration initiale

```bash
# Cloner le dépôt
git clone https://github.com/Fascinax/turbotilt.git
cd turbotilt

# Installer les dépendances
go mod download

# Exécuter les tests
make test

# Construire l'application
make build
```

### Structure du projet

- `cmd/` - Point d'entrée CLI (utilisant Cobra)
- `internal/` - Code interne de l'application
- `templates/` - Templates pour la génération de fichiers
- `scripts/` - Scripts utilitaires
- `test-projet/` - Projet de test pour valider les fonctionnalités

## Ressources supplémentaires

- [Documentation Go](https://golang.org/doc/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Documentation Tilt](https://docs.tilt.dev/)

Merci pour votre contribution à Turbotilt !
