package runtime

import (
	"fmt"
	"os"
	"os/exec"
)

// Variable pour faciliter les tests unitaires
var execCommand = exec.Command

// Variable pour faciliter les tests unitaires
var isTiltInstalled = checkTiltInstalled

// Options pour l'exécution
type RunOptions struct {
	UseTilt     bool
	Detached    bool
	TempFiles   []string
	ServiceName string // Nom du service à démarrer (pour les projets multi-services)
	DryRun      bool   // Mode simulation sans modifications réelles
	Debug       bool   // Mode débug avec logs détaillés
}

// TiltUp lance Tilt avec les options spécifiées
func TiltUp(opts RunOptions) error {
	// Configurer le nettoyage des fichiers temporaires
	if len(opts.TempFiles) > 0 && !opts.DryRun {
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
	args := []string{"up"}
	
	if opts.Debug {
		args = append(args, "--debug")
	}
	
	if opts.DryRun {
		fmt.Printf("🔍 [DRY-RUN] Commande qui serait exécutée: tilt %s\n", args)
		return nil
	}
	
	cmd := execCommand("tilt", args...)
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

	// Si un service spécifique est demandé
	if opts.ServiceName != "" {
		fmt.Printf("🔍 Démarrage du service spécifique: %s\n", opts.ServiceName)
		args = append(args, opts.ServiceName)
	}
	
	if opts.DryRun {
		fmt.Printf("🔍 [DRY-RUN] Commande qui serait exécutée: docker %s\n", args)
		return nil
	}
	
	cmd := execCommand("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// checkTiltInstalled est l'implémentation réelle de la vérification
// Séparée pour permettre le mocking dans les tests
func checkTiltInstalled() bool {
	cmd := exec.Command("tilt", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
