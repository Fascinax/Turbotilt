package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewTranslator(t *testing.T) {
	// Reset the LANG environment variable
	origLang := os.Getenv("LANG")
	defer os.Setenv("LANG", origLang)

	// Test with default language
	os.Setenv("LANG", "")
	translator := NewTranslator()
	if translator.currentLang != DefaultLang {
		t.Errorf("Default language should be %s, but got %s", DefaultLang, translator.currentLang)
	}

	// Test with a supported language
	os.Setenv("LANG", "en_US.UTF-8")
	translator = NewTranslator()
	if translator.currentLang != "en" {
		t.Errorf("Language should be 'en', but got %s", translator.currentLang)
	}

	// Test with an unsupported language
	os.Setenv("LANG", "de_DE.UTF-8")
	translator = NewTranslator()
	if translator.currentLang != DefaultLang {
		t.Errorf("Language should fall back to %s, but got %s", DefaultLang, translator.currentLang)
	}
}

func TestSetLang(t *testing.T) {
	translator := NewTranslator()

	// Test with a supported language
	if !translator.SetLang("en") {
		t.Error("SetLang should return true for a supported language")
	}
	if translator.currentLang != "en" {
		t.Errorf("Language should be 'en', but got %s", translator.currentLang)
	}

	// Test with an unsupported language
	if translator.SetLang("de") {
		t.Error("SetLang should return false for an unsupported language")
	}
	if translator.currentLang != "en" {
		t.Errorf("Language should not change, but changed to %s", translator.currentLang)
	}
}

func TestGetLang(t *testing.T) {
	translator := NewTranslator()
	translator.currentLang = "en"

	if lang := translator.GetLang(); lang != "en" {
		t.Errorf("GetLang devrait retourner 'en', mais a retourné %s", lang)
	}
}

func TestT(t *testing.T) {
	translator := NewTranslator()

	// Add translations for testing
	if translator.translations["fr"] == nil {
		translator.translations["fr"] = make(map[string]string)
	}
	if translator.translations["en"] == nil {
		translator.translations["en"] = make(map[string]string)
	}

	translator.translations["fr"]["test_key"] = "Valeur de test"
	translator.translations["en"]["test_key"] = "Test value"
	translator.translations["fr"]["test_format"] = "Valeur %s avec %d paramètres"
	translator.translations["en"]["test_format"] = "Value %s with %d parameters"

	// Test with default language (French)
	translator.currentLang = "fr" // Force French language
	if msg := translator.T("test_key"); msg != "Valeur de test" {
		t.Errorf("T should return 'Valeur de test', but got %s", msg)
	}

	// Test with another language
	translator.SetLang("en")
	if msg := translator.T("test_key"); msg != "Test value" {
		t.Errorf("T should return 'Test value', but got %s", msg)
	}

	// Test with formatting parameters
	if msg := translator.T("test_format", "test", 2); msg != "Value test with 2 parameters" {
		t.Errorf("T should return 'Value test with 2 parameters', but got %s", msg)
	}

	// Test with a non-existent key
	if msg := translator.T("nonexistent_key"); msg != "nonexistent_key" {
		t.Errorf("T should return the key itself, but got %s", msg)
	}

	// Test with a key that exists in the default language but not in the current language
	translator.translations["fr"]["only_fr"] = "Seulement en français"
	translator.SetLang("en")
	if msg := translator.T("only_fr"); msg != "Seulement en français" {
		t.Errorf("T should fall back to the default language, but got %s", msg)
	}
}

func TestLoadTranslations(t *testing.T) {
	// Create a temporary directory for translation files
	tempDir := t.TempDir()

	// Create a test.json translation file
	testContent := `{
		"test_key": "Test translation",
		"format_key": "Format with %s"
	}`
	if err := os.WriteFile(filepath.Join(tempDir, "test.json"), []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create translation file: %v", err)
	}

	// Create a non-JSON file
	if err := os.WriteFile(filepath.Join(tempDir, "not-json.txt"), []byte("Not a JSON file"), 0644); err != nil {
		t.Fatalf("Failed to create non-JSON file: %v", err)
	}

	translator := NewTranslator()
	if err := translator.LoadTranslations(tempDir); err != nil {
		t.Fatalf("LoadTranslations returned an error: %v", err)
	}

	// Verify that the translation was loaded
	if _, ok := translator.translations["test"]; !ok {
		t.Error("Language 'test' was not loaded")
	}

	// Change language and test translation
	translator.SetLang("test")
	if msg := translator.T("test_key"); msg != "Test translation" {
		t.Errorf("T should return 'Test translation', but got %s", msg)
	}

	// Test with a non-existent directory
	if err := translator.LoadTranslations("/nonexistent"); err == nil {
		t.Error("LoadTranslations should return an error for a non-existent directory")
	}
}

func TestGlobalTranslator(t *testing.T) {
	// Reset the global translator
	global = nil

	// Initialize the global translator
	Init()

	// Force French language
	global.currentLang = "fr"

	// Add a test translation
	if global.translations["fr"] == nil {
		global.translations["fr"] = make(map[string]string)
	}
	global.translations["fr"]["test_key"] = "Clé de test"

	// Test the global T function
	if msg := T("test_key"); msg != "Clé de test" {
		t.Errorf("Global T should return 'Clé de test', but got %s", msg)
	}

	// Test SetLanguage
	if !SetLanguage("en") {
		t.Error("SetLanguage should return true for a supported language")
	}

	// Verify that the language was changed
	if global.currentLang != "en" {
		t.Errorf("Language should be 'en', but got %s", global.currentLang)
	}
}
