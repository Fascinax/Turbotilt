# Sélection de Microservices avec Turbotilt

La commande `select` vous permet de scanner un répertoire pour détecter les microservices et de sélectionner ceux que vous souhaitez lancer. C'est particulièrement utile pour les projets comportant plusieurs microservices où vous ne souhaitez travailler qu'avec un sous-ensemble d'entre eux.

## Utilisation de base

```bash
turbotilt select [répertoire]
```

Si aucun répertoire n'est spécifié, le répertoire actuel sera scanné.

## Comment ça marche

1. La commande analyse le répertoire spécifié pour détecter les microservices en recherchant des fichiers de projet comme `pom.xml`, `build.gradle`, `package.json`, etc.
2. Elle affiche une liste des microservices détectés, chacun avec un numéro d'index.
3. Vous sélectionnez les services à inclure en entrant les numéros correspondants (séparés par des virgules).
4. Selon les options que vous avez utilisées, Turbotilt peut :
   - Créer un fichier de configuration permanent (`turbotilt.yaml`)
   - Lancer immédiatement les services sélectionnés
   - Les deux options ci-dessus

## Options

- `-o, --output <filename>` : Spécifie un nom pour le fichier de configuration généré (par défaut : `turbotilt.yaml`)
- `-c, --create-config` : Crée un fichier `turbotilt.yaml` avec les services sélectionnés
- `-l, --launch` : Lance les services sélectionnés après la sélection

## Exemples

### Scanner et sélectionner uniquement

```bash
# Scanner le répertoire actuel et sélectionner les services
turbotilt select

# Scanner un répertoire spécifique et sélectionner les services
turbotilt select ./mon-projet
```

### Créer un fichier de configuration

```bash
# Sélectionner les services et créer un fichier turbotilt.yaml
turbotilt select -c

# Sélectionner les services et créer un fichier de configuration personnalisé
turbotilt select --output ma-config-personnalisee.yaml
```

### Sélectionner et lancer

```bash
# Sélectionner les services et les lancer immédiatement
turbotilt select -l

# Sélectionner les services, créer un fichier de configuration et les lancer
turbotilt select -c -l
```

## Cas d'utilisation

### Environnement de développement temporaire

Lorsque vous souhaitez démarrer rapidement un sous-ensemble de services sans créer de configuration permanente :

```bash
turbotilt select -l
```

Cette commande scannera les services, vous permettra de sélectionner ceux à inclure, et les lancera sans créer de fichier de configuration permanent.

### Création de plusieurs fichiers de configuration

Pour un projet avec de nombreux microservices, vous pourriez vouloir créer différentes configurations pour différents scénarios de développement :

```bash
# Créer une configuration pour les services backend
turbotilt select --output services-backend.yaml

# Créer une configuration pour les services frontend
turbotilt select --output services-frontend.yaml
```

Plus tard, vous pouvez utiliser ces configurations avec la commande `up` :

```bash
turbotilt up -f services-backend.yaml
```

### Intégration de nouveaux membres d'équipe

La commande `select` offre un moyen simple pour les nouveaux membres d'équipe de comprendre l'architecture d'un projet complexe :

1. Exécutez `turbotilt select` pour voir tous les microservices disponibles
2. Sélectionnez des services spécifiques pour comprendre comment ils fonctionnent ensemble
3. Sauvegardez la configuration pour une utilisation future avec `--create-config`
