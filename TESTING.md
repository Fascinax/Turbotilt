# Tests pour Turbotilt

## Structure des tests

Les tests sont organisés selon la structure du projet pour faciliter la maintenance :

```
turbotilt/
├── internal/
│   ├── scan/
│   │   ├── scan_test.go         # Tests de détection de framework
│   │   └── services_test.go     # Tests de détection de services
│   ├── config/
│   │   └── config_test.go       # Tests de gestion de configuration
│   ├── render/
│   │   └── render_test.go       # Tests de génération de fichiers
```

## Exécution des tests

### Tous les tests

Pour exécuter tous les tests, utilisez la commande :

```powershell
go test ./...
```

Ou avec le Makefile :

```powershell
make test
```

### Tests d'un package spécifique

Pour tester uniquement un package spécifique :

```powershell
go test ./internal/scan
go test ./internal/config
go test ./internal/render
```

### Tests avec couverture

Pour générer un rapport de couverture de tests :

```powershell
make coverage
```

Cela générera un fichier `coverage.html` que vous pourrez ouvrir dans un navigateur.

## Ajout de nouveaux tests

Pour ajouter de nouveaux tests, suivez les conventions Go standard :

1. Créez un fichier `*_test.go` dans le même package que le code testé
2. Utilisez des noms de fonctions commençant par `Test`
3. Utilisez le package `testing` et la fonction `t.Run()` pour les sous-tests
4. Préférez des tests unitaires indépendants qui ne dépendent pas d'état externe

## Bonnes pratiques

- Utilisez des répertoires temporaires pour les tests nécessitant des fichiers
- Nettoyez tous les fichiers créés pendant les tests
- Gardez les tests simples et concentrés sur une seule fonctionnalité
- Documentez les tests complexes avec des commentaires explicatifs
