package render

import (
	"os"
	"text/template"
)

// Données pour le template Tiltfile
type TiltfileData struct {
	ImageName    string
	LiveSyncPath string
	Framework    string
}

// GenerateTiltfile génère un Tiltfile personnalisé pour le projet
func GenerateTiltfile(framework string) error {
	f, err := os.Create("Tiltfile")
	if err != nil {
		return err
	}
	defer f.Close()

	// Chemin de sync par défaut
	syncPath := "./src"

	// Ajuster le chemin de sync selon le framework
	switch framework {
	case "spring":
		syncPath = "./src/main/java"
	case "quarkus":
		syncPath = "./src/main/java"
	}

	tiltData := TiltfileData{
		ImageName:    "app",
		LiveSyncPath: syncPath,
		Framework:    framework,
	}

	// Contenu du Tiltfile selon le framework
	tiltContent := `# Turbotilt - Tiltfile généré automatiquement

docker_build('{{ .ImageName }}', '.', 
  dockerfile='Dockerfile',
  live_update=[
    sync('{{ .LiveSyncPath }}', '/app/{{ .LiveSyncPath }}'),
    run('echo "Files synced to container"', trigger=['{{ .LiveSyncPath }}'])
  ]
)

# Configuration spécifique pour {{ .Framework }}
{{ if eq .Framework "spring" }}
# Configuration Spring Boot
{{ else if eq .Framework "quarkus" }}
# Configuration Quarkus
{{ else }}
# Configuration générique
{{ end }}

k8s_yaml('docker-compose.yml')
`

	tmpl, err := template.New("tiltfile").Parse(tiltContent)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, tiltData)
}
