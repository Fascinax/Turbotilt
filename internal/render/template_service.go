package render

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateService fournit des méthodes pour gérer les templates
type TemplateService struct {
	FuncMap       template.FuncMap
	TemplatesDirs []string
	Delimiters    [2]string
}

// NewTemplateService crée une nouvelle instance de TemplateService
func NewTemplateService() *TemplateService {
	return &TemplateService{
		FuncMap: template.FuncMap{
			"eq": func(a, b interface{}) bool {
				return a == b
			},
		},
		TemplatesDirs: []string{
			"templates",
			"../templates",
			".", // Pour les tests qui créent le template directement
		},
		Delimiters: [2]string{"[[", "]]"},
	}
}

// FindTemplateFile recherche un fichier template dans les répertoires configurés
func (ts *TemplateService) FindTemplateFile(baseNames ...string) (string, error) {
	for _, dir := range ts.TemplatesDirs {
		for _, baseName := range baseNames {
			tmplPath := filepath.Join(dir, baseName)
			if _, statErr := os.Stat(tmplPath); statErr == nil {
				return tmplPath, nil
			}
		}
	}
	return "", fmt.Errorf("aucun fichier template trouvé parmi: %v", baseNames)
}

// LoadTemplate charge un template à partir d'un chemin ou utilise un template par défaut si non trouvé
func (ts *TemplateService) LoadTemplate(name string, templatePaths []string, defaultTemplate string) (*template.Template, error) {
	// Essayer de charger depuis les fichiers
	tmplPath, err := ts.FindTemplateFile(templatePaths...)
	if err == nil {
		return template.New(name).
			Delims(ts.Delimiters[0], ts.Delimiters[1]).
			Funcs(ts.FuncMap).
			ParseFiles(tmplPath)
	}

	// Si non trouvé, utiliser le template par défaut
	return template.New(name).
		Delims(ts.Delimiters[0], ts.Delimiters[1]).
		Funcs(ts.FuncMap).
		Parse(defaultTemplate)
}

// RenderToFile génère un fichier à partir d'un template et de données
func (ts *TemplateService) RenderToFile(filePath string, tmpl *template.Template, data interface{}) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier %s: %w", filePath, err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("erreur lors de l'exécution du template: %w", err)
	}

	return nil
}
