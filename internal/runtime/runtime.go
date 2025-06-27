package runtime

import (
	"fmt"
	"os"
	"os/exec"
)

// Options pour l'exécution
type RunOptions struct {
	UseTilt   bool
	Detached  bool
	TempFiles []string
}

// TiltUp lance Tilt avec les options spécifiées
func TiltUp(opts RunOptions) error {
	// Configurer le nettoyage des fichiers temporaires
	if len(opts.TempFiles) > 0 {
		SetupCleanup(opts.TempFiles)
	}

	// Vérifier si Tilt est installé
	if !isTiltInstalled() && opts.UseTilt {
		fmt.Println("⚠️ Tilt n'est pas installé. Utilisation de Docker Compose.")
		return ComposeUp(opts)
	}

	if !opts.UseTilt {
		return ComposeUp(opts)
	}

	fmt.Println("🚀 Démarrage avec Tilt...")
	cmd := exec.Command("tilt", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ComposeUp lance Docker Compose avec les options spécifiées
func ComposeUp(opts RunOptions) error {
	fmt.Println("🐳 Démarrage avec Docker Compose...")
	args := []string{"compose", "up"}

	if opts.Detached {
		args = append(args, "-d")
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// isTiltInstalled vérifie si Tilt est installé
func isTiltInstalled() bool {
	cmd := exec.Command("tilt", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
