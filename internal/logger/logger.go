package logger

import (
	"fmt"
	"os"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// Global logger configuration
var (
	currentLevel LogLevel = INFO
	logFile      *os.File
	useColors    bool = true
	useEmojis    bool = true
)

// SetLevel sets the minimum log level
func SetLevel(level LogLevel) {
	currentLevel = level
}

// EnableFileLogging enables logging to a file
func EnableFileLogging(path string) error {
	var err error
	logFile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return nil
}

// DisableFileLogging disables file logging and closes the file
func DisableFileLogging() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

// SetUseColors enables/disables colors in log output
func SetUseColors(use bool) {
	useColors = use
}

// SetUseEmojis enables/disables emojis in log output
func SetUseEmojis(use bool) {
	useEmojis = use
}

// Debug logs a message at DEBUG level
func Debug(format string, args ...interface{}) {
	if currentLevel <= DEBUG {
		log(DEBUG, format, args...)
	}
}

// Info logs a message at INFO level
func Info(format string, args ...interface{}) {
	if currentLevel <= INFO {
		log(INFO, format, args...)
	}
}

// Warning logs a message at WARNING level
func Warning(format string, args ...interface{}) {
	if currentLevel <= WARNING {
		log(WARNING, format, args...)
	}
}

// Error logs a message at ERROR level
func Error(format string, args ...interface{}) {
	if currentLevel <= ERROR {
		log(ERROR, format, args...)
	}
}

// Logger represents a logger instance
type Logger struct{}

// GetLogger returns a logger instance
func GetLogger() *Logger {
	return &Logger{}
}

// Debug logs a message at DEBUG level
func (l *Logger) Debug(format string, args ...interface{}) {
	if currentLevel <= DEBUG {
		log(DEBUG, format, args...)
	}
}

// Info logs a message at INFO level
func (l *Logger) Info(format string, args ...interface{}) {
	if currentLevel <= INFO {
		log(INFO, format, args...)
	}
}

// Infof logs a message at INFO level with formatting
func (l *Logger) Infof(format string, args ...interface{}) {
	if currentLevel <= INFO {
		log(INFO, format, args...)
	}
}

// Warning logs a message at WARNING level
func (l *Logger) Warning(format string, args ...interface{}) {
	if currentLevel <= WARNING {
		log(WARNING, format, args...)
	}
}

// Warningf logs a message at WARNING level with formatting
func (l *Logger) Warningf(format string, args ...interface{}) {
	if currentLevel <= WARNING {
		log(WARNING, format, args...)
	}
}

// Error logs a message at ERROR level
func (l *Logger) Error(format string, args ...interface{}) {
	if currentLevel <= ERROR {
		log(ERROR, format, args...)
	}
}

// Errorf logs a message at ERROR level with formatting
func (l *Logger) Errorf(format string, args ...interface{}) {
	if currentLevel <= ERROR {
		log(ERROR, format, args...)
	}
}

// Fatal logs a message at FATAL level and terminates the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	log(FATAL, format, args...)
	os.Exit(1)
}

// log is the internal logging function
func log(level LogLevel, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var levelStr, colorPrefix, colorSuffix, emoji string

	switch level {
	case DEBUG:
		levelStr = "DEBUG"
		colorPrefix = "\033[36m" // Cyan
		emoji = "🔍"
	case INFO:
		levelStr = "INFO"
		colorPrefix = "\033[32m" // Green
		emoji = "ℹ️"
	case WARNING:
		levelStr = "WARN"
		colorPrefix = "\033[33m" // Yellow
		emoji = "⚠️"
	case ERROR:
		levelStr = "ERROR"
		colorPrefix = "\033[31m" // Red
		emoji = "❌"
	case FATAL:
		levelStr = "FATAL"
		colorPrefix = "\033[35m" // Magenta
		emoji = "💥"
	}

	colorSuffix = "\033[0m"
	if !useColors {
		colorPrefix = ""
		colorSuffix = ""
	}

	if !useEmojis {
		emoji = ""
	}

	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("%s %s[%s]%s %s%s",
		timestamp,
		colorPrefix,
		levelStr,
		colorSuffix,
		emoji,
		message)

	fmt.Println(logLine)

	if logFile != nil {
		fileLogLine := fmt.Sprintf("%s [%s] %s\n", timestamp, levelStr, message)
		if _, err := logFile.WriteString(fileLogLine); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
		}
	}
}
