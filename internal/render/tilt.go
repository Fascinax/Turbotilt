package render

import (
	"os"
	"text/template"
)

// GenerateTiltfile génère un Tiltfile personnalisé pour le projet
func GenerateTiltfile(opts Options) error {
	f, err := os.Create("Tiltfile")
	if err != nil {
		return err
	}
	defer f.Close()

	// Chemin de sync par défaut
	syncPath := "./src"
	
	// Ajuster le chemin de sync selon le framework
	switch opts.Framework {
	case "spring", "quarkus", "micronaut":
		syncPath = "./src/main/java"
	}

	// Contenu du Tiltfile selon le framework
	tiltContent := `# Turbotilt - Tiltfile généré automatiquement
# Framework: {{ .Framework }}

docker_build('app', '.', 
  dockerfile='Dockerfile',
  live_update=[
    sync('{{ .SyncPath }}', '/app/{{ .SyncPath }}'),
{{if eq .Framework "spring"}}    # Rechargement à chaud pour Spring
    run('touch /app/src/main/resources/application.properties', trigger=['{{ .SyncPath }}/**/*.java']),
{{else if eq .Framework "quarkus"}}    # Rechargement à chaud pour Quarkus
    run('touch /app/src/main/resources/application.properties', trigger=['{{ .SyncPath }}/**/*.java']),
{{else if eq .Framework "micronaut"}}    # Rechargement à chaud pour Micronaut
    run('touch /app/src/main/resources/application.yml', trigger=['{{ .SyncPath }}/**/*.java']),
{{else}}    # Rechargement simple
{{end}}
    run('echo "Files synced to container at $(date)"', trigger=['{{ .SyncPath }}'])
  ]
)

# Configuration spécifique
{{if eq .Framework "spring"}}
# Spring Boot configuration
# Hot reload using Spring DevTools
{{else if eq .Framework "quarkus"}}
# Quarkus configuration
# Hot reload using Quarkus Dev Services
{{else if eq .Framework "micronaut"}}
# Micronaut configuration
# Hot reload using Micronaut DevTools
{{end}}

# Mode: {{if .DevMode}}Development{{else}}Production{{end}}

# Port: {{ .Port }}

k8s_yaml('docker-compose.yml')
`
	
	data := struct {
		Options
		SyncPath string
	}{
		Options: opts,
		SyncPath: syncPath,
	}
	
	tmpl, err := template.New("tiltfile").Parse(tiltContent)
	if err != nil {
		return err
	}
	
	return tmpl.Execute(f, data)
}
