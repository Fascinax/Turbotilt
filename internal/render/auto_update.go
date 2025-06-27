package render

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"turbotilt/internal/logger"
)

// Configuration interface to avoid import cycle
type ConfigInterface interface {
	GetServiceName() string
	GetFramework() string
	// Add other methods as needed
}

// AutoUpdateWatcher watches a Tiltfile for changes and updates it automatically
type AutoUpdateWatcher struct {
	tiltfilePath    string
	projectRoot     string
	conf            ConfigInterface
	lastCheckTime   time.Time
	checkInterval   time.Duration
	isRunning       bool
	stopChan        chan struct{}
	updateTriggered chan bool
}

// NewAutoUpdateWatcher creates a new watcher for auto-updating Tiltfiles
func NewAutoUpdateWatcher(tiltfilePath string, conf ConfigInterface) *AutoUpdateWatcher {
	return &AutoUpdateWatcher{
		tiltfilePath:    tiltfilePath,
		projectRoot:     filepath.Dir(tiltfilePath),
		conf:            conf,
		lastCheckTime:   time.Now(),
		checkInterval:   5 * time.Second,
		stopChan:        make(chan struct{}),
		updateTriggered: make(chan bool, 1),
	}
}

// Start begins the auto-update watcher
func (w *AutoUpdateWatcher) Start() error {
	if w.isRunning {
		return fmt.Errorf("auto-update watcher already running")
	}
	
	logger.Info("Starting auto-update watcher for %s", w.tiltfilePath)
	w.isRunning = true
	
	go w.watchLoop()
	
	return nil
}

// Stop halts the auto-update watcher
func (w *AutoUpdateWatcher) Stop() {
	if !w.isRunning {
		return
	}
	
	w.stopChan <- struct{}{}
	w.isRunning = false
	logger.Info("Stopped auto-update watcher for %s", w.tiltfilePath)
}

// TriggerUpdate manually triggers an update check
func (w *AutoUpdateWatcher) TriggerUpdate() bool {
	if !w.isRunning {
		return false
	}
	
	w.updateTriggered <- true
	return true
}

// WaitForUpdate blocks until an update is performed or timeout occurs
func (w *AutoUpdateWatcher) WaitForUpdate(timeout time.Duration) bool {
	if !w.isRunning {
		return false
	}
	
	select {
	case <-time.After(timeout):
		return false
	case <-w.updateTriggered:
		return true
	}
}

// watchLoop is the main loop that checks for changes
func (w *AutoUpdateWatcher) watchLoop() {
	ticker := time.NewTicker(w.checkInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-w.stopChan:
			return
		case <-ticker.C:
			w.checkForChanges()
		}
	}
}

// checkForChanges checks if source files have changed and updates the Tiltfile if needed
func (w *AutoUpdateWatcher) checkForChanges() {
	// Scan for file changes that would affect the Tiltfile
	changed := w.detectFileChanges()
	
	if changed {
		logger.Info("Detected changes, updating Tiltfile...")
		err := w.updateTiltfile()
		if err != nil {
			logger.Error("Failed to update Tiltfile: %v", err)
		} else {
			logger.Info("Tiltfile updated successfully")
			// Notify any waiting goroutines
			select {
			case w.updateTriggered <- true:
			default:
			}
		}
	}
}

// detectFileChanges checks if any source files or configuration have changed
func (w *AutoUpdateWatcher) detectFileChanges() bool {
	// Check if pom.xml, build.gradle, or docker-compose.yml have changed
	sources := []string{
		filepath.Join(w.projectRoot, "pom.xml"),
		filepath.Join(w.projectRoot, "build.gradle"),
		filepath.Join(w.projectRoot, "docker-compose.yml"),
		filepath.Join(w.projectRoot, "Dockerfile"),
	}
	
	for _, src := range sources {
		if fileChanged(src, w.lastCheckTime) {
			w.lastCheckTime = time.Now()
			return true
		}
	}
	
	return false
}

// updateTiltfile regenerates the Tiltfile
func (w *AutoUpdateWatcher) updateTiltfile() error {
	// Read the existing Tiltfile template
	templatePath := filepath.Join(w.projectRoot, "turbotilt", "Tiltfile.template")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}
	
	// Generate new Tiltfile content (here we'd normally call the renderer)
	content := string(templateContent)
	
	// Write the updated content to the file
	err = os.WriteFile(w.tiltfilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Tiltfile: %w", err)
	}
	
	return nil
}

// fileChanged checks if a file has been modified since the last check time
func fileChanged(filepath string, since time.Time) bool {
	info, err := os.Stat(filepath)
	if err != nil {
		// If the file doesn't exist or can't be accessed, consider it unchanged
		return false
	}
	
	return info.ModTime().After(since)
}
