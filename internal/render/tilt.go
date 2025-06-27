package render

import (
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
}

// GenerateTiltfile génère un Tiltfile personnalisé pour le projet
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
		Framework:  opts.Framework,
		AppName:    appName,
		Port:       port,
		Date:       time.Now().Format("2006-01-02 15:04:05"),
		DevMode:    opts.DevMode,
		Services:   services,
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
