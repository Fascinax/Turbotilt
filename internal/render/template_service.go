package render

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateService provides methods for managing templates
type TemplateService struct {
	FuncMap       template.FuncMap
	TemplatesDirs []string
	Delimiters    [2]string
}

// NewTemplateService creates a new instance of TemplateService
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
			".", // For tests that create the template directly
		},
		Delimiters: [2]string{"[[", "]]"},
	}
}

// FindTemplateFile searches for a template file in the configured directories
func (ts *TemplateService) FindTemplateFile(baseNames ...string) (string, error) {
	for _, dir := range ts.TemplatesDirs {
		for _, baseName := range baseNames {
			tmplPath := filepath.Join(dir, baseName)
			if _, statErr := os.Stat(tmplPath); statErr == nil {
				return tmplPath, nil
			}
		}
	}
	return "", fmt.Errorf("no template file found among: %v", baseNames)
}

// LoadTemplate loads a template from a path or uses a default template if not found
func (ts *TemplateService) LoadTemplate(name string, templatePaths []string, defaultTemplate string) (*template.Template, error) {
	// Try to load from files
	tmplPath, err := ts.FindTemplateFile(templatePaths...)
	if err == nil {
		return template.New(name).
			Delims(ts.Delimiters[0], ts.Delimiters[1]).
			Funcs(ts.FuncMap).
			ParseFiles(tmplPath)
	}

	// If not found, use the default template
	return template.New(name).
		Delims(ts.Delimiters[0], ts.Delimiters[1]).
		Funcs(ts.FuncMap).
		Parse(defaultTemplate)
}

// RenderToFile generates a file from a template and data
func (ts *TemplateService) RenderToFile(filePath string, tmpl *template.Template, data interface{}) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filePath, err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
