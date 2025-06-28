package runtime

import (
	"fmt"
	"os"
	"os/exec"
)

// Variable to facilitate unit testing
var execCommand = exec.Command

// Variable to facilitate unit testing
var isTiltInstalled = checkTiltInstalled

// Options for execution
type RunOptions struct {
	UseTilt     bool
	Detached    bool
	TempFiles   []string
	ServiceName string // Name of the service to start (for multi-service projects)
	DryRun      bool   // Simulation mode without actual changes
	Debug       bool   // Debug mode with detailed logs
}

// TiltUp launches Tilt with the specified options
func TiltUp(opts RunOptions) error {
	// Set up cleanup for temporary files
	if len(opts.TempFiles) > 0 && !opts.DryRun {
		SetupCleanup(opts.TempFiles)
	}

	// Check if Tilt is installed
	if !isTiltInstalled() && opts.UseTilt {
		fmt.Println("⚠️ Tilt is not installed. Using Docker Compose.")
		return ComposeUp(opts)
	}

	if !opts.UseTilt {
		return ComposeUp(opts)
	}

	fmt.Println("🚀 Starting with Tilt...")
	args := []string{"up"}

	if opts.Debug {
		args = append(args, "--debug")
	}

	if opts.DryRun {
		fmt.Printf("🔍 [DRY-RUN] Command that would be executed: tilt %s\n", args)
		return nil
	}

	cmd := execCommand("tilt", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ComposeUp launches Docker Compose with the specified options
func ComposeUp(opts RunOptions) error {
	fmt.Println("🐳 Starting with Docker Compose...")
	args := []string{"compose", "up"}

	if opts.Detached {
		args = append(args, "-d")
	}

	// If a specific service is requested
	if opts.ServiceName != "" {
		fmt.Printf("🔍 Starting specific service: %s\n", opts.ServiceName)
		args = append(args, opts.ServiceName)
	}

	if opts.DryRun {
		fmt.Printf("🔍 [DRY-RUN] Command that would be executed: docker %s\n", args)
		return nil
	}

	cmd := execCommand("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// checkTiltInstalled is the actual implementation of the check
// Separated to allow mocking in tests
func checkTiltInstalled() bool {
	cmd := exec.Command("tilt", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
