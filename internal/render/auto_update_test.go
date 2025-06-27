package render

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Implémentation mock de ConfigInterface pour les tests
type mockConfig struct {
	serviceName string
	framework   string
}

func (m *mockConfig) GetServiceName() string {
	return m.serviceName
}

func (m *mockConfig) GetFramework() string {
	return m.framework
}

func TestNewAutoUpdateWatcher(t *testing.T) {
	// Créer un mock de configuration
	conf := &mockConfig{
		serviceName: "test-service",
		framework:   "spring",
	}

	// Créer un watcher
	osPath := filepath.Join("test", "path", "Tiltfile")
	watcher := NewAutoUpdateWatcher(osPath, conf)

	// Vérifier les propriétés du watcher
	if watcher.tiltfilePath != osPath {
		t.Errorf("tiltfilePath devrait être '%s', mais est '%s'", osPath, watcher.tiltfilePath)
	}

	expected := filepath.Join("test", "path")
	if watcher.projectRoot != expected {
		t.Errorf("projectRoot devrait être '%s', mais est '%s'", expected, watcher.projectRoot)
	}

	if watcher.conf != conf {
		t.Error("conf n'a pas été correctement assigné")
	}

	if watcher.isRunning {
		t.Error("isRunning devrait être false initialement")
	}

	if watcher.checkInterval != 5*time.Second {
		t.Errorf("checkInterval devrait être 5s, mais est %v", watcher.checkInterval)
	}
}

func TestAutoUpdateWatcherStart(t *testing.T) {
	// Créer un répertoire temporaire pour le test
	tempDir := t.TempDir()
	tiltfilePath := filepath.Join(tempDir, "Tiltfile")

	// Créer un Tiltfile vide
	if err := os.WriteFile(tiltfilePath, []byte("# Test Tiltfile"), 0644); err != nil {
		t.Fatalf("Impossible de créer le Tiltfile: %v", err)
	}

	// Créer un mock de configuration
	conf := &mockConfig{
		serviceName: "test-service",
		framework:   "spring",
	}

	// Créer un watcher
	watcher := NewAutoUpdateWatcher(tiltfilePath, conf)

	// Démarrer le watcher
	err := watcher.Start()
	if err != nil {
		t.Fatalf("Start() a retourné une erreur: %v", err)
	}

	// Vérifier que le watcher est en cours d'exécution
	if !watcher.isRunning {
		t.Error("isRunning devrait être true après Start()")
	}

	// Test de démarrage multiple
	err = watcher.Start()
	if err == nil {
		t.Error("Start() devrait retourner une erreur si le watcher est déjà en cours d'exécution")
	}

	// Arrêter le watcher
	watcher.Stop()

	// Vérifier que le watcher est arrêté
	if watcher.isRunning {
		t.Error("isRunning devrait être false après Stop()")
	}
}

func TestTriggerUpdate(t *testing.T) {
	// Créer un répertoire temporaire pour le test
	tempDir := t.TempDir()
	tiltfilePath := filepath.Join(tempDir, "Tiltfile")

	// Créer un Tiltfile vide
	if err := os.WriteFile(tiltfilePath, []byte("# Test Tiltfile"), 0644); err != nil {
		t.Fatalf("Impossible de créer le Tiltfile: %v", err)
	}

	// Créer un mock de configuration
	conf := &mockConfig{
		serviceName: "test-service",
		framework:   "spring",
	}

	// Créer un watcher
	watcher := NewAutoUpdateWatcher(tiltfilePath, conf)

	// Déclencher une mise à jour
	watcher.TriggerUpdate()

	// Vérifier que la mise à jour a été déclenchée
	select {
	case triggered := <-watcher.updateTriggered:
		if !triggered {
			t.Error("updateTriggered devrait être true")
		}
	default:
		t.Error("TriggerUpdate() n'a pas envoyé de signal sur le canal updateTriggered")
	}
}
