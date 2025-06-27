package render

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// TiltfileTemplateData contient les données pour le template du Tiltfile
type TiltfileTemplateData struct {
	Framework      string
	AppName        string
	Port           string
	Date           string
	DevMode        bool
	Services       []interface{}
	IsMultiService bool // Indique s'il s'agit d'un projet multi-services
}

// Template de base pour le Tiltfile multi-service au cas où le fichier de template n'existe pas
// Commented out as it's currently not used
// const tiltfileTemplateMulti = `# Tiltfile généré par Turbotilt
// # Date: {{.Date}}
//
// k8s_yaml('docker-compose.yml')
// `

const (
	// DefaultTiltfileTemplate est le template par défaut pour le Tiltfile
	DefaultTiltfileTemplate = `# Tiltfile généré par Turbotilt (default template)
# Date: [[.Date]]
# Framework: [[.Framework]]
docker_compose('docker-compose.yml')
`

	// DefaultTiltfileMultiTemplate est le template par défaut pour le Tiltfile multi-service
	DefaultTiltfileMultiTemplate = `# Tiltfile multi-service généré par Turbotilt (default template)
# Date: [[.Date]]
docker_compose('docker-compose.yml')
`

	// TemplatePathTiltfile est le chemin vers le fichier de template Tiltfile
	TemplatePathTiltfile = "Tiltfile.tmpl"

	// TemplatePathTiltfileMulti est le chemin vers le fichier de template Tiltfile multi-service
	TemplatePathTiltfileMulti = "Tiltfile.multi.tmpl"
)

// GenerateTiltfile génère un Tiltfile personnalisé pour le projet mono-service
func GenerateTiltfile(opts Options) error {
	// Initialiser le service de templates
	ts := NewTemplateService()

	// Préparer les données pour le template
	appName := "app"
	if opts.AppName == "" {
		// Utiliser le nom du répertoire courant comme nom par défaut
		cwd, err := os.Getwd()
		if err == nil {
			appName = filepath.Base(cwd)
		}
	} else {
		appName = opts.AppName
	}

	port := DefaultPort // Port par défaut
	if opts.Port != "" {
		port = opts.Port
	}

	// Préparer les données à injecter dans le template
	// Convertir les services en []interface{} pour le template
	services := make([]interface{}, len(opts.Services))
	for i, svc := range opts.Services {
		services[i] = svc
	}

	data := TiltfileTemplateData{
		Framework:      opts.Framework,
		AppName:        appName,
		Port:           port,
		Date:           time.Now().Format("2006-01-02 15:04:05"),
		DevMode:        opts.DevMode,
		Services:       services,
		IsMultiService: false,
	}

	// Template par défaut à utiliser si aucun fichier template n'est trouvé
	const defaultTemplate = `# Tiltfile généré par Turbotilt (default template)
# Date: [[.Date]]
# Framework: [[.Framework]]
docker_compose('docker-compose.yml')
`

	// Essayer de charger le template
	tmpl, err := ts.LoadTemplate("Tiltfile",
		[]string{TemplatePathTiltfile}, defaultTemplate)
	if err != nil {
		return fmt.Errorf("erreur lors du chargement du template: %w", err)
	}

	// Générer le fichier Tiltfile
	return ts.RenderToFile("Tiltfile", tmpl, data)
}

// GenerateMultiServiceTiltfile génère un Tiltfile pour un projet multi-services
func GenerateMultiServiceTiltfile(serviceList ServiceList) error {
	// Initialiser le service de templates
	ts := NewTemplateService()

	// Créer des données combinées pour tous les services
	services := make([]map[string]interface{}, len(serviceList.Services))
	for i, svc := range serviceList.Services {
		services[i] = map[string]interface{}{
			"Name":       svc.ServiceName,
			"Path":       svc.Path,
			"Framework":  svc.Framework,
			"Port":       svc.Port,
			"DevMode":    svc.DevMode,
			"JDKVersion": svc.JDKVersion,
		}
	}

	data := TiltfileTemplateData{
		Date:           time.Now().Format("2006-01-02 15:04:05"),
		Services:       []interface{}{services},
		IsMultiService: true,
	}

	// Template par défaut à utiliser si aucun fichier template n'est trouvé
	const defaultTemplate = `# Tiltfile multi-service généré par Turbotilt (default template)
# Date: [[.Date]]
docker_compose('docker-compose.yml')
`

	// Essayer de charger le template
	tmpl, err := ts.LoadTemplate("Tiltfile",
		[]string{TemplatePathTiltfileMulti, TemplatePathTiltfile},
		DefaultTiltfileMultiTemplate)
	if err != nil {
		return fmt.Errorf("erreur lors du chargement du template: %w", err)
	}

	// Générer le fichier Tiltfile
	return ts.RenderToFile("Tiltfile", tmpl, data)
}

// GenerateTiltfileFromTemplate génère un Tiltfile personnalisé pour le test
func GenerateTiltfileFromTemplate(opts Options, templatePath string, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Préparer les données pour le template
	data := TiltfileTemplateData{
		Framework:      opts.Framework,
		AppName:        opts.AppName,
		Port:           opts.Port,
		Date:           time.Now().Format("2006-01-02 15:04:05"),
		DevMode:        opts.DevMode,
		IsMultiService: false,
	}

	// Charger le template depuis le fichier
	tmpl, err := template.New(filepath.Base(templatePath)).Delims("[[", "]]").Parse(string(mustReadFile(templatePath)))
	if err != nil {
		return err
	}

	// Exécuter le template
	return tmpl.Execute(f, data)
}

// mustReadFile lit un fichier et panique en cas d'erreur
func mustReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}

// Commented out as it's currently not used
// createTestTemplate creates a template file for testing and returns its path
// func createTestTemplate(content, filename string) (string, error) {
// 	// Create in current directory for tests
// 	err := os.WriteFile(filename, []byte(content), 0644)
// 	if err != nil {
// 		return "", err
// 	}
// 	return filename, nil
// }

// Note: La constante tiltfileTemplateMulti et la fonction generateMultiServiceTiltfile
// ont été déplacées en haut du fichier pour éviter la duplication.
