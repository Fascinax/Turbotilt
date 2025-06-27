package logger

import (
	"fmt"
	"os"
	"time"
)

// LogLevel représente le niveau de log
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// Configuration globale du logger
var (
	currentLevel LogLevel = INFO
	logFile      *os.File
	useColors    bool     = true
	useEmojis    bool     = true
)

// SetLevel définit le niveau minimum de log
func SetLevel(level LogLevel) {
	currentLevel = level
}

// EnableFileLogging active la journalisation dans un fichier
func EnableFileLogging(path string) error {
	var err error
	logFile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return nil
}

// DisableFileLogging désactive la journalisation dans un fichier
func DisableFileLogging() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

// SetUseColors active/désactive les couleurs dans les logs
func SetUseColors(use bool) {
	useColors = use
}

// SetUseEmojis active/désactive les emojis dans les logs
func SetUseEmojis(use bool) {
	useEmojis = use
}

// Debug log un message au niveau DEBUG
func Debug(format string, args ...interface{}) {
	if currentLevel <= DEBUG {
		log(DEBUG, format, args...)
	}
}

// Info log un message au niveau INFO
func Info(format string, args ...interface{}) {
	if currentLevel <= INFO {
		log(INFO, format, args...)
	}
}

// Warning log un message au niveau WARNING
func Warning(format string, args ...interface{}) {
	if currentLevel <= WARNING {
		log(WARNING, format, args...)
	}
}

// Error log un message au niveau ERROR
func Error(format string, args ...interface{}) {
	if currentLevel <= ERROR {
		log(ERROR, format, args...)
	}
}

// Fatal log un message et termine le programme
func Fatal(format string, args ...interface{}) {
	log(FATAL, format, args...)
	os.Exit(1)
}

// log est la fonction interne de journalisation
func log(level LogLevel, format string, args ...interface{}) {
	// Timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	// Préfixes selon le niveau
	var levelStr, colorPrefix, colorSuffix, emoji string
	
	switch level {
	case DEBUG:
		levelStr = "DEBUG"
		colorPrefix = "\033[36m" // Cyan
		emoji = "🔍"
	case INFO:
		levelStr = "INFO"
		colorPrefix = "\033[32m" // Vert
		emoji = "ℹ️"
	case WARNING:
		levelStr = "WARN"
		colorPrefix = "\033[33m" // Jaune
		emoji = "⚠️"
	case ERROR:
		levelStr = "ERROR"
		colorPrefix = "\033[31m" // Rouge
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
	
	// Formater le message
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("%s %s[%s]%s %s%s", 
		timestamp, 
		colorPrefix, 
		levelStr, 
		colorSuffix,
		emoji,
		message)
	
	// Afficher sur la console
	fmt.Println(logLine)
	
	// Enregistrer dans le fichier si configuré
	if logFile != nil {
		// Sans couleurs dans le fichier
		fileLogLine := fmt.Sprintf("%s [%s] %s\n", timestamp, levelStr, message)
		logFile.WriteString(fileLogLine)
	}
}
