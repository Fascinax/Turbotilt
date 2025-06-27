package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewTranslator(t *testing.T) {
	// Réinitialiser la variable d'environnement LANG
	origLang := os.Getenv("LANG")
	defer os.Setenv("LANG", origLang)

	// Test avec la langue par défaut
	os.Setenv("LANG", "")
	translator := NewTranslator()
	if translator.currentLang != DefaultLang {
		t.Errorf("La langue par défaut devrait être %s, mais elle est %s", DefaultLang, translator.currentLang)
	}

	// Test avec une langue supportée
	os.Setenv("LANG", "en_US.UTF-8")
	translator = NewTranslator()
	if translator.currentLang != "en" {
		t.Errorf("La langue devrait être 'en', mais elle est %s", translator.currentLang)
	}

	// Test avec une langue non supportée
	os.Setenv("LANG", "de_DE.UTF-8")
	translator = NewTranslator()
	if translator.currentLang != DefaultLang {
		t.Errorf("La langue devrait revenir à %s, mais elle est %s", DefaultLang, translator.currentLang)
	}
}

func TestSetLang(t *testing.T) {
	translator := NewTranslator()

	// Test avec une langue supportée
	if !translator.SetLang("en") {
		t.Error("SetLang devrait retourner true pour une langue supportée")
	}
	if translator.currentLang != "en" {
		t.Errorf("La langue devrait être 'en', mais elle est %s", translator.currentLang)
	}

	// Test avec une langue non supportée
	if translator.SetLang("de") {
		t.Error("SetLang devrait retourner false pour une langue non supportée")
	}
	if translator.currentLang != "en" {
		t.Errorf("La langue ne devrait pas changer, mais elle est devenue %s", translator.currentLang)
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
	
	// Ajout de traductions pour le test
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
	
	// Test avec la langue par défaut (français)
	translator.currentLang = "fr" // Forcer la langue française
	if msg := translator.T("test_key"); msg != "Valeur de test" {
		t.Errorf("T devrait retourner 'Valeur de test', mais a retourné %s", msg)
	}
	
	// Test avec une autre langue
	translator.SetLang("en")
	if msg := translator.T("test_key"); msg != "Test value" {
		t.Errorf("T devrait retourner 'Test value', mais a retourné %s", msg)
	}
	
	// Test avec des paramètres de formatage
	if msg := translator.T("test_format", "test", 2); msg != "Value test with 2 parameters" {
		t.Errorf("T devrait retourner 'Value test with 2 parameters', mais a retourné %s", msg)
	}
	
	// Test avec une clé non existante
	if msg := translator.T("nonexistent_key"); msg != "nonexistent_key" {
		t.Errorf("T devrait retourner la clé elle-même, mais a retourné %s", msg)
	}
	
	// Test avec une clé existante dans la langue par défaut mais pas dans la langue courante
	translator.translations["fr"]["only_fr"] = "Seulement en français"
	translator.SetLang("en")
	if msg := translator.T("only_fr"); msg != "Seulement en français" {
		t.Errorf("T devrait revenir à la langue par défaut, mais a retourné %s", msg)
	}
}

func TestLoadTranslations(t *testing.T) {
	// Créer un répertoire temporaire pour les fichiers de traduction
	tempDir := t.TempDir()
	
	// Créer un fichier de traduction test.json
	testContent := `{
		"test_key": "Test translation",
		"format_key": "Format with %s"
	}`
	if err := os.WriteFile(filepath.Join(tempDir, "test.json"), []byte(testContent), 0644); err != nil {
		t.Fatalf("Impossible de créer le fichier de traduction: %v", err)
	}
	
	// Créer un fichier non-JSON
	if err := os.WriteFile(filepath.Join(tempDir, "not-json.txt"), []byte("Not a JSON file"), 0644); err != nil {
		t.Fatalf("Impossible de créer le fichier non-JSON: %v", err)
	}
	
	translator := NewTranslator()
	if err := translator.LoadTranslations(tempDir); err != nil {
		t.Fatalf("LoadTranslations a retourné une erreur: %v", err)
	}
	
	// Vérifier que la traduction a été chargée
	if _, ok := translator.translations["test"]; !ok {
		t.Error("La langue 'test' n'a pas été chargée")
	}
	
	// Changer la langue et tester la traduction
	translator.SetLang("test")
	if msg := translator.T("test_key"); msg != "Test translation" {
		t.Errorf("T devrait retourner 'Test translation', mais a retourné %s", msg)
	}
	
	// Test avec un dossier inexistant
	if err := translator.LoadTranslations("/nonexistent"); err == nil {
		t.Error("LoadTranslations devrait retourner une erreur pour un dossier inexistant")
	}
}

func TestGlobalTranslator(t *testing.T) {
	// Réinitialiser le traducteur global
	global = nil
	
	// Initialiser le traducteur global
	Init()
	
	// Forcer la langue française
	global.currentLang = "fr"
	
	// Ajouter une traduction de test
	if global.translations["fr"] == nil {
		global.translations["fr"] = make(map[string]string)
	}
	global.translations["fr"]["test_key"] = "Clé de test"
	
	// Test de la fonction T globale
	if msg := T("test_key"); msg != "Clé de test" {
		t.Errorf("T global devrait retourner 'Clé de test', mais a retourné %s", msg)
	}
	
	// Test de SetLanguage
	if !SetLanguage("en") {
		t.Error("SetLanguage devrait retourner true pour une langue supportée")
	}
	
	// Vérifier que la langue a été changée
	if global.currentLang != "en" {
		t.Errorf("La langue devrait être 'en', mais elle est %s", global.currentLang)
	}
}
