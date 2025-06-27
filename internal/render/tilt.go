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
	Framework  string
	AppName    string
	Port       string
	Date       string
	DevMode    bool
	Services   []interface{}
	IsMultiService bool // Indique s'il s'agit d'un projet multi-services
}

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
		Framework:     opts.Framework,
		AppName:       appName,
		Port:          port,
		Date:          time.Now().Format("2006-01-02 15:04:05"),
		DevMode:       opts.DevMode,
		Services:      services,
		IsMultiService: false,
	}
	
	// Charger le template depuis le fichier
	tmplPath := "templates/Tiltfile.tmpl"
	tmpl, err := template.New("Tiltfile").Delims("[[", "]]").Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		return err
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
			"Name":     svc.ServiceName,
			"Path":     svc.Path,
			"Framework": svc.Framework,
			"Port":     svc.Port,
			"DevMode":  svc.DevMode,
			"JDKVersion": svc.JDKVersion,
		}
	}
	
	data := TiltfileTemplateData{
		Date:          time.Now().Format("2006-01-02 15:04:05"),
		Services:      []interface{}{services},
		IsMultiService: true,
	}
	
	// Utiliser un template spécifique pour les projets multi-services
	tmplPath := "templates/Tiltfile.multi.tmpl"
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		// Utiliser le template standard s'il n'y a pas de template multi-service
		tmplPath = "templates/Tiltfile.tmpl"
	}
	
	tmpl, err := template.New("Tiltfile").Delims("[[", "]]").Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("erreur parsing template: %w", err)
	}
	
	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("erreur exécution template: %w", err)
	}
	
	return nil
}

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
