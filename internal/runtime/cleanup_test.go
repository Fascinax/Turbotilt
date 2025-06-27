package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupTempFiles(t *testing.T) {
	// Créer des fichiers temporaires pour le test
	tempDir := t.TempDir()
	tempFiles := []string{
		filepath.Join(tempDir, "file1.tmp"),
		filepath.Join(tempDir, "file2.tmp"),
	}
	
	// Créer les fichiers
	for _, file := range tempFiles {
		if err := os.WriteFile(file, []byte("test content"), 0644); err != nil {
			t.Fatalf("Impossible de créer le fichier temporaire: %v", err)
		}
	}
	
	// Vérifier que les fichiers existent
	for _, file := range tempFiles {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("Le fichier %s devrait exister: %v", file, err)
		}
	}
	
	// Exécuter la fonction de nettoyage
	CleanupTempFiles(tempFiles)
	
	// Vérifier que les fichiers ont été supprimés
	for _, file := range tempFiles {
		if _, err := os.Stat(file); err == nil {
			t.Errorf("Le fichier %s devrait avoir été supprimé", file)
		} else if !os.IsNotExist(err) {
			t.Errorf("Erreur inattendue lors de la vérification du fichier %s: %v", file, err)
		}
	}
}
