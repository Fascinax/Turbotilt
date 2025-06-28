package logger

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureOutput captures standard output for testing
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
	// Save and restore log level
	oldLevel := currentLevel
	defer func() { currentLevel = oldLevel }()

	// Test level modification
	SetLevel(DEBUG)
	if currentLevel != DEBUG {
		t.Errorf("Level should be DEBUG, but it is %v", currentLevel)
	}

	SetLevel(ERROR)
	if currentLevel != ERROR {
		t.Errorf("Level should be ERROR, but it is %v", currentLevel)
	}
}

func TestFileLogging(t *testing.T) {
	// Disable file logging to start
	DisableFileLogging()
	if logFile != nil {
		t.Error("Log file should be nil after DisableFileLogging")
	}

	// Create a temporary file for logs
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	// Enable file logging
	err := EnableFileLogging(logPath)
	if err != nil {
		t.Fatalf("EnableFileLogging returned an error: %v", err)
	}
	if logFile == nil {
		t.Error("Log file should not be nil after EnableFileLogging")
	}

	// Write a log message
	Info("Test message")

	// Close the log file
	DisableFileLogging()
	if logFile != nil {
		t.Error("Log file should be nil after DisableFileLogging")
	}

	// Verify that the message was written to the file
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Unable to read log file: %v", err)
	}
	if !strings.Contains(string(content), "Test message") {
		t.Error("Log message was not written to the file")
	}
}

func TestSetUseColors(t *testing.T) {
	// Save and restore parameter
	oldUseColors := useColors
	defer func() { useColors = oldUseColors }()

	// Test modification
	SetUseColors(false)
	if useColors != false {
		t.Error("useColors should be false")
	}

	SetUseColors(true)
	if useColors != true {
		t.Error("useColors should be true")
	}
}

func TestSetUseEmojis(t *testing.T) {
	// Save and restore parameter
	oldUseEmojis := useEmojis
	defer func() { useEmojis = oldUseEmojis }()

	// Test modification
	SetUseEmojis(false)
	if useEmojis != false {
		t.Error("useEmojis should be false")
	}

	SetUseEmojis(true)
	if useEmojis != true {
		t.Error("useEmojis should be true")
	}
}

func TestLogLevels(t *testing.T) {
	// Save and restore parameters
	oldLevel := currentLevel
	oldUseColors := useColors
	oldUseEmojis := useEmojis
	defer func() {
		currentLevel = oldLevel
		useColors = oldUseColors
		useEmojis = oldUseEmojis
	}()

	// Disable colors and emojis to simplify tests
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
			// Replace Fatal with a function that doesn't exit the program
			log(FATAL, format, args...)
		}, true, "[FATAL]"},
	}

	for _, test := range tests {
		// Set the minimum level for the message to be displayed
		SetLevel(test.level)

		output := captureOutput(func() {
			test.logFunc("Test message")
		})

		if test.shouldPrint && !strings.Contains(output, test.contains) {
			t.Errorf("Log of level %v should contain '%s', but produced: %s",
				test.level, test.contains, output)
		}

		// Test filtering by setting the minimum level above
		if test.level < FATAL {
			SetLevel(test.level + 1)
			output = captureOutput(func() {
				test.logFunc("Test message")
			})
			if output != "" {
				t.Errorf("Log of level %v should not be displayed when minimum level is %v",
					test.level, test.level+1)
			}
		}
	}
}

func TestFormatting(t *testing.T) {
	// Save and restore parameters
	oldLevel := currentLevel
	oldUseColors := useColors
	oldUseEmojis := useEmojis
	defer func() {
		currentLevel = oldLevel
		useColors = oldUseColors
		useEmojis = oldUseEmojis
	}()

	// Disable colors and emojis
	SetUseColors(false)
	SetUseEmojis(false)
	SetLevel(INFO)

	output := captureOutput(func() {
		Info("Test %s with %d parameters", "message", 2)
	})

	if !strings.Contains(output, "Test message with 2 parameters") {
		t.Errorf("Message formatting is not correct: %s", output)
	}
}
