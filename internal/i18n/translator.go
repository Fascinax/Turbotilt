package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Default language
const (
	DefaultLang = "fr"
)

// Translator manages translations
type Translator struct {
	currentLang  string
	translations map[string]map[string]string
}

// NewTranslator creates a new translator
func NewTranslator() *Translator {
	t := &Translator{
		currentLang:  DefaultLang,
		translations: make(map[string]map[string]string),
	}

	// Load built-in translations
	t.translations["fr"] = frTranslations
	t.translations["en"] = enTranslations

	// Try to detect system language
	t.detectSystemLanguage()

	return t
}

// GetTranslator returns a singleton instance of the translator
func GetTranslator() *Translator {
	return NewTranslator()
}

// detectSystemLanguage tries to detect the system language
func (t *Translator) detectSystemLanguage() {
	lang := os.Getenv("LANG")
	if lang != "" {
		lang = strings.Split(lang, "_")[0]
		if _, ok := t.translations[lang]; ok {
			t.currentLang = lang
		}
	}
}

// LoadTranslations loads translations from a JSON file
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

// SetLang sets the current language
func (t *Translator) SetLang(lang string) bool {
	if _, ok := t.translations[lang]; ok {
		t.currentLang = lang
		return true
	}
	return false
}

// GetLang returns the current language
func (t *Translator) GetLang() string {
	return t.currentLang
}

// T translates a key
func (t *Translator) T(key string, args ...interface{}) string {
	// Try current language
	if trans, ok := t.translations[t.currentLang]; ok {
		if str, ok := trans[key]; ok {
			if len(args) > 0 {
				return fmt.Sprintf(str, args...)
			}
			return str
		}
	}

	// Try default language if different
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

	// Fallback to key
	return key
}

// Tr translates a key (alias for T)
func (t *Translator) Tr(key string, args ...interface{}) string {
	return t.T(key, args...)
}

// Global instance
var global *Translator

// Init initializes the global translator
func Init() {
	global = NewTranslator()
}

// SetLanguage sets the global language
func SetLanguage(lang string) bool {
	if global == nil {
		Init()
	}
	return global.SetLang(lang)
}

// T translates a key with the global translator
func T(key string, args ...interface{}) string {
	if global == nil {
		Init()
	}
	return global.T(key, args...)
}
