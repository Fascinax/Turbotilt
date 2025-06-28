package runtime

import (
	"os"
	"os/signal"
	"syscall"
)

// SetupCleanup configures the cleanup of temporary files on program exit
func SetupCleanup(tempFiles []string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		CleanupTempFiles(tempFiles)
		os.Exit(0)
	}()
}

// CleanupTempFiles removes the generated temporary files
func CleanupTempFiles(tempFiles []string) {
	for _, file := range tempFiles {
		os.Remove(file)
	}
}
