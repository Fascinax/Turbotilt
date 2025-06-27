package runtime

import (
	"os"
	"os/signal"
	"syscall"
)

// SetupCleanup configure le nettoyage des fichiers temporaires à la fermeture
func SetupCleanup(tempFiles []string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c
		CleanupTempFiles(tempFiles)
		os.Exit(0)
	}()
}

// CleanupTempFiles supprime les fichiers temporaires générés
func CleanupTempFiles(tempFiles []string) {
	for _, file := range tempFiles {
		os.Remove(file)
	}
}
