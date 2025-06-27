package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Langue par défaut
const (
	DefaultLang = "fr"
)

// Translator gère les traductions
type Translator struct {
	currentLang  string
	translations map[string]map[string]string
}

// NewTranslator crée un nouveau traducteur
func NewTranslator() *Translator {
	t := &Translator{
		currentLang:  DefaultLang,
		translations: make(map[string]map[string]string),
	}

	// Charger les traductions intégrées
	t.translations["fr"] = frTranslations
	t.translations["en"] = enTranslations

	// Essayer de détecter la langue du système
	lang := os.Getenv("LANG")
	if lang != "" {
		lang = strings.Split(lang, "_")[0]
		if _, ok := t.translations[lang]; ok {
			t.currentLang = lang
		}
	}

	return t
}

// LoadTranslations charge les traductions depuis un fichier JSON
func (t *Translator) LoadTranslations(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			lang := strings.TrimSuffix(file.Name(), ".json")

			data, err := os.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				continue
			}

			var translations map[string]string
			if err := json.Unmarshal(data, &translations); err != nil {
				continue
			}

			t.translations[lang] = translations
		}
	}

	return nil
}

// SetLang définit la langue courante
func (t *Translator) SetLang(lang string) bool {
	if _, ok := t.translations[lang]; ok {
		t.currentLang = lang
		return true
	}
	return false
}

// GetLang renvoie la langue courante
func (t *Translator) GetLang() string {
	return t.currentLang
}

// T traduit une clé
func (t *Translator) T(key string, args ...interface{}) string {
	// Essayer la langue courante
	if trans, ok := t.translations[t.currentLang]; ok {
		if str, ok := trans[key]; ok {
			if len(args) > 0 {
				return fmt.Sprintf(str, args...)
			}
			return str
		}
	}

	// Essayer la langue par défaut si différente
	if t.currentLang != DefaultLang {
		if trans, ok := t.translations[DefaultLang]; ok {
			if str, ok := trans[key]; ok {
				if len(args) > 0 {
					return fmt.Sprintf(str, args...)
				}
				return str
			}
		}
	}

	// Fallback à la clé
	return key
}

// Global instance
var global *Translator

// Init initialise le traducteur global
func Init() {
	global = NewTranslator()
}

// SetLanguage définit la langue globale
func SetLanguage(lang string) bool {
	if global == nil {
		Init()
	}
	return global.SetLang(lang)
}

// T traduit une clé avec le traducteur global
func T(key string, args ...interface{}) string {
	if global == nil {
		Init()
	}
	return global.T(key, args...)
}
