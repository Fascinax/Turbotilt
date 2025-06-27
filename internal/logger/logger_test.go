package logger

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureOutput capture la sortie standard pour les tests
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestSetLevel(t *testing.T) {
	// Sauvegarde et restauration du niveau de log
	oldLevel := currentLevel
	defer func() { currentLevel = oldLevel }()

	// Test de modification du niveau
	SetLevel(DEBUG)
	if currentLevel != DEBUG {
		t.Errorf("Le niveau devrait être DEBUG, mais il est %v", currentLevel)
	}

	SetLevel(ERROR)
	if currentLevel != ERROR {
		t.Errorf("Le niveau devrait être ERROR, mais il est %v", currentLevel)
	}
}

func TestFileLogging(t *testing.T) {
	// Désactiver le logging de fichier pour commencer
	DisableFileLogging()
	if logFile != nil {
		t.Error("Le fichier de log devrait être nil après DisableFileLogging")
	}

	// Créer un fichier temporaire pour les logs
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	// Activer le logging de fichier
	err := EnableFileLogging(logPath)
	if err != nil {
		t.Fatalf("EnableFileLogging a retourné une erreur: %v", err)
	}
	if logFile == nil {
		t.Error("Le fichier de log ne devrait pas être nil après EnableFileLogging")
	}

	// Écrire un message de log
	Info("Test message")

	// Fermer le fichier de log
	DisableFileLogging()
	if logFile != nil {
		t.Error("Le fichier de log devrait être nil après DisableFileLogging")
	}

	// Vérifier que le message a été écrit dans le fichier
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Impossible de lire le fichier de log: %v", err)
	}
	if !strings.Contains(string(content), "Test message") {
		t.Error("Le message de log n'a pas été écrit dans le fichier")
	}
}

func TestSetUseColors(t *testing.T) {
	// Sauvegarde et restauration du paramètre
	oldUseColors := useColors
	defer func() { useColors = oldUseColors }()

	// Test de modification
	SetUseColors(false)
	if useColors != false {
		t.Error("useColors devrait être false")
	}

	SetUseColors(true)
	if useColors != true {
		t.Error("useColors devrait être true")
	}
}

func TestSetUseEmojis(t *testing.T) {
	// Sauvegarde et restauration du paramètre
	oldUseEmojis := useEmojis
	defer func() { useEmojis = oldUseEmojis }()

	// Test de modification
	SetUseEmojis(false)
	if useEmojis != false {
		t.Error("useEmojis devrait être false")
	}

	SetUseEmojis(true)
	if useEmojis != true {
		t.Error("useEmojis devrait être true")
	}
}

func TestLogLevels(t *testing.T) {
	// Sauvegarde et restauration des paramètres
	oldLevel := currentLevel
	oldUseColors := useColors
	oldUseEmojis := useEmojis
	defer func() {
		currentLevel = oldLevel
		useColors = oldUseColors
		useEmojis = oldUseEmojis
	}()

	// Désactiver les couleurs et emojis pour simplifier les tests
	SetUseColors(false)
	SetUseEmojis(false)

	tests := []struct {
		level       LogLevel
		logFunc     func(string, ...interface{})
		shouldPrint bool
		contains    string
	}{
		{DEBUG, Debug, true, "[DEBUG]"},
		{INFO, Info, true, "[INFO]"},
		{WARNING, Warning, true, "[WARN]"},
		{ERROR, Error, true, "[ERROR]"},
		{FATAL, func(format string, args ...interface{}) {
			// Remplacer Fatal par une fonction qui ne quitte pas le programme
			log(FATAL, format, args...)
		}, true, "[FATAL]"},
	}

	for _, test := range tests {
		// Régler le niveau minimum pour que le message soit affiché
		SetLevel(test.level)

		output := captureOutput(func() {
			test.logFunc("Test message")
		})

		if test.shouldPrint && !strings.Contains(output, test.contains) {
			t.Errorf("Le log de niveau %v devrait contenir '%s', mais a produit: %s",
				test.level, test.contains, output)
		}

		// Tester le filtrage en réglant le niveau minimum au-dessus
		if test.level < FATAL {
			SetLevel(test.level + 1)
			output = captureOutput(func() {
				test.logFunc("Test message")
			})
			if output != "" {
				t.Errorf("Le log de niveau %v ne devrait pas être affiché quand le niveau minimum est %v",
					test.level, test.level+1)
			}
		}
	}
}

func TestFormatting(t *testing.T) {
	// Sauvegarde et restauration des paramètres
	oldLevel := currentLevel
	oldUseColors := useColors
	oldUseEmojis := useEmojis
	defer func() {
		currentLevel = oldLevel
		useColors = oldUseColors
		useEmojis = oldUseEmojis
	}()

	// Désactiver les couleurs et emojis
	SetUseColors(false)
	SetUseEmojis(false)
	SetLevel(INFO)

	output := captureOutput(func() {
		Info("Test %s with %d parameters", "message", 2)
	})

	if !strings.Contains(output, "Test message with 2 parameters") {
		t.Errorf("Le formatage du message n'est pas correct: %s", output)
	}
}
