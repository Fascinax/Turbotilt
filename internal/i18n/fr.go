package i18n

// frTranslations contient les traductions françaises
var frTranslations = map[string]string{
	// Messages généraux
	"app.name":                "Turbotilt",
	"app.description":         "CLI pour environnements dev cloud-native",
	
	// Commandes
	"cmd.init":                "Initialiser un environnement de développement",
	"cmd.up":                  "Démarrer l'environnement",
	"cmd.stop":                "Arrêter l'environnement",
	"cmd.doctor":              "Vérifier l'environnement et la configuration",
	"cmd.version":             "Afficher les informations de version",
	
	// Messages de succès
	"success.init":            "✅ Configuration Turbotilt terminée!",
	"success.up":              "✅ Environnement démarré avec succès!",
	"success.stop":            "✅ Environnement arrêté avec succès!",
	
	// Messages d'erreur
	"error.not_found":         "❌ %s introuvable",
	"error.docker":            "❌ Erreur Docker: %s",
	"error.init":              "❌ Erreur lors de l'initialisation: %s",
	"error.up":                "❌ Erreur lors du démarrage: %s",
	"error.stop":              "❌ Erreur lors de l'arrêt: %s",
	"error.doctor":            "❌ Problème détecté: %s",
	"error.config":            "❌ Erreur de configuration: %s",
	
	// Détection
	"detect.framework":        "🔍 Détection du framework...",
	"detect.services":         "🔍 Détection des services dépendants...",
	"detect.success":          "✅ Détecté: %s",
	"detect.fail":             "❌ Aucun framework connu détecté",
	
	// Génération
	"generate.dockerfile":     "📄 Génération du Dockerfile...",
	"generate.compose":        "📄 Génération du docker-compose.yml...",
	"generate.tiltfile":       "📄 Génération du Tiltfile...",
	"generate.manifest":       "📄 Génération du manifeste turbotilt.yaml...",
	
	// Doctor
	"doctor.checking":         "🔍 Vérification de %s...",
	"doctor.success":          "✅ %s est correctement configuré",
	"doctor.warning":          "⚠️ %s présente un problème mineur: %s",
	"doctor.error":            "❌ %s présente un problème: %s",
	"doctor.score":            "📊 Score de santé: %d/100",
	
	// Mise à jour
	"update.available":        "📦 Une nouvelle version est disponible: %s (actuelle: %s)",
	"update.download":         "💾 Téléchargez-la sur: %s",
	"update.current":          "✅ Vous utilisez la dernière version",
}
