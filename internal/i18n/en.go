package i18n

// enTranslations contient les traductions anglaises
var enTranslations = map[string]string{
	// General messages
	"app.name":                "Turbotilt",
	"app.description":         "CLI for cloud-native dev environments",
	
	// Commands
	"cmd.init":                "Initialize a development environment",
	"cmd.up":                  "Start the environment",
	"cmd.stop":                "Stop the environment",
	"cmd.doctor":              "Check the environment and configuration",
	"cmd.version":             "Display version information",
	
	// Success messages
	"success.init":            "✅ Turbotilt configuration completed!",
	"success.up":              "✅ Environment started successfully!",
	"success.stop":            "✅ Environment stopped successfully!",
	
	// Error messages
	"error.not_found":         "❌ %s not found",
	"error.docker":            "❌ Docker error: %s",
	"error.init":              "❌ Initialization error: %s",
	"error.up":                "❌ Startup error: %s",
	"error.stop":              "❌ Shutdown error: %s",
	"error.doctor":            "❌ Issue detected: %s",
	"error.config":            "❌ Configuration error: %s",
	
	// Detection
	"detect.framework":        "🔍 Detecting framework...",
	"detect.services":         "🔍 Detecting dependent services...",
	"detect.success":          "✅ Detected: %s",
	"detect.fail":             "❌ No known framework detected",
	
	// Generation
	"generate.dockerfile":     "📄 Generating Dockerfile...",
	"generate.compose":        "📄 Generating docker-compose.yml...",
	"generate.tiltfile":       "📄 Generating Tiltfile...",
	"generate.manifest":       "📄 Generating turbotilt.yaml manifest...",
	
	// Doctor
	"doctor.checking":         "🔍 Checking %s...",
	"doctor.success":          "✅ %s is properly configured",
	"doctor.warning":          "⚠️ %s has a minor issue: %s",
	"doctor.error":            "❌ %s has an issue: %s",
	"doctor.score":            "📊 Health score: %d/100",
	
	// Update
	"update.available":        "📦 A new version is available: %s (current: %s)",
	"update.download":         "💾 Download it at: %s",
	"update.current":          "✅ You are using the latest version",
}
