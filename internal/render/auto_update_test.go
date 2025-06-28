package render

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Mock implementation of ConfigInterface for tests
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
	// Create a mock configuration
	conf := &mockConfig{
		serviceName: "test-service",
		framework:   "spring",
	}

	// Create a watcher
	osPath := filepath.Join("test", "path", "Tiltfile")
	watcher := NewAutoUpdateWatcher(osPath, conf)

	// Verify the watcher properties
	if watcher.tiltfilePath != osPath {
		t.Errorf("tiltfilePath should be '%s', but is '%s'", osPath, watcher.tiltfilePath)
	}

	expected := filepath.Join("test", "path")
	if watcher.projectRoot != expected {
		t.Errorf("projectRoot should be '%s', but is '%s'", expected, watcher.projectRoot)
	}

	if watcher.conf != conf {
		t.Error("conf was not correctly assigned")
	}

	if watcher.isRunning {
		t.Error("isRunning should be false initially")
	}

	if watcher.checkInterval != 5*time.Second {
		t.Errorf("checkInterval should be 5s, but is %v", watcher.checkInterval)
	}
}

func TestAutoUpdateWatcherStart(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	tiltfilePath := filepath.Join(tempDir, "Tiltfile")

	// Create an empty Tiltfile
	if err := os.WriteFile(tiltfilePath, []byte("# Test Tiltfile"), 0644); err != nil {
		t.Fatalf("Unable to create Tiltfile: %v", err)
	}

	// Create a mock configuration
	conf := &mockConfig{
		serviceName: "test-service",
		framework:   "spring",
	}

	// Create a watcher
	watcher := NewAutoUpdateWatcher(tiltfilePath, conf)

	// Start the watcher
	err := watcher.Start()
	if err != nil {
		t.Fatalf("Start() returned an error: %v", err)
	}

	// Verify that the watcher is running
	if !watcher.isRunning {
		t.Error("isRunning should be true after Start()")
	}

	// Test multiple starts
	err = watcher.Start()
	if err == nil {
		t.Error("Start() should return an error if the watcher is already running")
	}

	// Stop the watcher
	watcher.Stop()

	// Verify that the watcher is stopped
	if watcher.isRunning {
		t.Error("isRunning should be false after Stop()")
	}
}

func TestTriggerUpdate(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	tiltfilePath := filepath.Join(tempDir, "Tiltfile")

	// Create an empty Tiltfile
	if err := os.WriteFile(tiltfilePath, []byte("# Test Tiltfile"), 0644); err != nil {
		t.Fatalf("Unable to create Tiltfile: %v", err)
	}

	// Create a mock configuration
	conf := &mockConfig{
		serviceName: "test-service",
		framework:   "spring",
	}

	// Create a watcher
	watcher := NewAutoUpdateWatcher(tiltfilePath, conf)

	// Trigger an update
	watcher.TriggerUpdate()

	// Verify that the update was triggered
	select {
	case triggered := <-watcher.updateTriggered:
		if !triggered {
			t.Error("updateTriggered should be true")
		}
	default:
		t.Error("TriggerUpdate() did not send a signal on the updateTriggered channel")
	}
}
