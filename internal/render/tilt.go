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
const tiltfileTemplateMulti = `# Tiltfile généré par Turbotilt
# Date: {{.Date}}

k8s_yaml('docker-compose.yml')
`

// GenerateTiltfile génère un Tiltfile personnalisé pour le projet mono-service
func GenerateTiltfile(opts Options) error {
	f, err := os.Create("Tiltfile")
	if err != nil {
		return err
	}
	defer f.Close()

	// Définir des délimiteurs personnalisés pour éviter les conflits avec la syntaxe Tilt
	funcMap := template.FuncMap{
		"eq": func(a, b interface{}) bool {
			return a == b
		},
	}

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

	port := "8080" // Port par défaut
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

	// Charger le template depuis le fichier
	// Chercher les templates dans cet ordre : répertoire courant, répertoire du projet, répertoire parent
	tmplPaths := []string{"templates/Tiltfile.tmpl", "../templates/Tiltfile.tmpl"}
	var tmpl *template.Template
	var templateErr error

	for _, tmplPath := range tmplPaths {
		tmpl, err = template.New("Tiltfile").Delims("[[", "]]").Funcs(funcMap).ParseFiles(tmplPath)
		if err == nil {
			break
		}
		templateErr = err
	}

	if err != nil {
		// Si le template n'a pas été trouvé, utiliser un template par défaut
		const defaultTemplate = `# Tiltfile généré par Turbotilt (default template)
# Date: [[.Date]]
# Framework: [[.Framework]]
docker_compose('docker-compose.yml')
`
		tmpl, err = template.New("Tiltfile").Delims("[[", "]]").Funcs(funcMap).Parse(defaultTemplate)
		if err != nil {
			return fmt.Errorf("erreur avec le template par défaut: %w (après %v)", err, templateErr)
		}
	}

	// Exécuter le template
	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}

// GenerateMultiServiceTiltfile génère un Tiltfile pour un projet multi-services
func GenerateMultiServiceTiltfile(serviceList ServiceList) error {
	f, err := os.Create("Tiltfile")
	if err != nil {
		return fmt.Errorf("erreur création Tiltfile: %w", err)
	}
	defer f.Close()

	// Définir des délimiteurs personnalisés pour éviter les conflits avec la syntaxe Tilt
	funcMap := template.FuncMap{
		"eq": func(a, b interface{}) bool {
			return a == b
		},
	}

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

	// Chercher un template pour les projets multi-services
	tmplPaths := []string{
		"templates/Tiltfile.multi.tmpl",
		"../templates/Tiltfile.multi.tmpl",
		"templates/Tiltfile.tmpl",
		"../templates/Tiltfile.tmpl",
	}

	var tmpl *template.Template
	var templateErr error

	for _, tmplPath := range tmplPaths {
		tmpl, err = template.New("Tiltfile").Delims("[[", "]]").Funcs(funcMap).ParseFiles(tmplPath)
		if err == nil {
			break
		}
		templateErr = err
	}

	if err != nil {
		// Si le template n'a pas été trouvé, utiliser un template par défaut
		const defaultTemplate = `# Tiltfile multi-service généré par Turbotilt (default template)
# Date: [[.Date]]
docker_compose('docker-compose.yml')
`
		tmpl, err = template.New("Tiltfile").Delims("[[", "]]").Funcs(funcMap).Parse(defaultTemplate)
		if err != nil {
			return fmt.Errorf("erreur parsing template: %w (après %v)", err, templateErr)
		}
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("erreur exécution template: %w", err)
	}

	return nil
}

// Note: La constante tiltfileTemplateMulti et la fonction generateMultiServiceTiltfile
// ont été déplacées en haut du fichier pour éviter la duplication.
