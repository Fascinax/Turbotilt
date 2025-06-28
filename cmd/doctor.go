package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"turbotilt/internal/config"
	"turbotilt/internal/logger"
)

// Commented out diagnostic result structure - not used in current implementation
// type diagResult struct {
// 	installed bool
// 	version   string
// 	detail    string
// 	weight    int // Importance: 3 = critical, 2 = important, 1 = optional
// 	required  bool
// }

var (
	verbose          bool
	logToFile        bool
	showAllInfo      bool
	showSummary      bool
	logFilePath      string
	validateManifest bool
	// Commented out variables - not used in current implementation
	// checkTiltfile    bool
	// fixMode          bool
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Checks installation & config",
	Long: `Checks the installation and configuration of the Turbotilt environment.
Performs a complete diagnostic of dependencies and necessary tools.
Provides a health score and recommendations to fix problems.

Examples:
  turbotilt doctor             # Complete analysis
  turbotilt doctor --verbose   # Complete analysis with details
  turbotilt doctor --debug     # Debug mode for detailed information`,
	Run: func(cmd *cobra.Command, args []string) {
		// Record start time to calculate execution duration
		startTime := time.Now()

		// Configure logger
		if verbose {
			logger.SetLevel(logger.DEBUG)
			logger.Debug("Verbose mode enabled")
		}

		// If validate-manifest option is enabled, only validate the manifest without running the full diagnostic
		if validateManifest {
			// Import config package
			validateManifestFile()
			return
		}

		// Configure log file path if not specified
		if logToFile && logFilePath == "" {
			logFilePath = fmt.Sprintf("turbotilt-doctor-%s.log",
				time.Now().Format("20060102-150405"))
		}

		if logToFile {
			if err := logger.EnableFileLogging(logFilePath); err != nil {
				fmt.Printf("⚠️ Unable to create log file: %v\n", err)
			} else {
				defer logger.DisableFileLogging()
				logger.Info("Log file created: %s", logFilePath)
				if verbose {
					fmt.Printf("📄 Logs saved to: %s\n\n", logFilePath)
				}
			}
		}

		logger.Info("=== Turbotilt Environment Diagnostic ===")
		logger.Info("Started on %s", time.Now().Format("2006-01-02 15:04:05"))
		logger.Info("Operating system: %s", runtime.GOOS)
		fmt.Println("🔍 Checking Turbotilt environment...")

		// Structure to store results
		type diagResult struct {
			installed bool
			version   string
			detail    string
			weight    int // Importance: 3 = critical, 2 = important, 1 = optional
			required  bool
		}
		results := make(map[string]diagResult)

		fmt.Println("\n📋 Checking required dependencies:")
		logger.Debug("Checking required dependencies...")

		// Check Docker (critical)
		fmt.Print("⏳ Docker : ")
		logger.Debug("Checking Docker installation...")
		if version, err := execCommand("docker", "--version"); err == nil {
			truncVersion := strings.TrimSpace(version)
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Docker installed: %s", truncVersion)

			// Additional check for Docker daemon
			if _, err := execCommand("docker", "info"); err != nil {
				fmt.Println("   ⚠️ Docker daemon doesn't appear to be running")
				logger.Warning("Docker daemon not responding")
				results["docker"] = diagResult{true, truncVersion, "Daemon not responding", 3, true}
			} else {
				results["docker"] = diagResult{true, truncVersion, "OK", 3, true}
			}
		} else {
			fmt.Println("❌ Not installed or not accessible")
			logger.Warning("Docker not installed or not accessible")
			results["docker"] = diagResult{false, "", "Not installed", 3, true}
		}

		// Check Docker Compose (critical)
		fmt.Print("⏳ Docker Compose : ")
		logger.Debug("Checking Docker Compose installation...")
		if version, err := execCommand("docker", "compose", "version"); err == nil {
			truncVersion := strings.Split(version, "\n")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Docker Compose installed: %s", truncVersion)
			results["docker-compose"] = diagResult{true, truncVersion, "OK", 3, true}
		} else {
			fmt.Println("❌ Not installed or not accessible")
			logger.Warning("Docker Compose not installed or not accessible")
			results["docker-compose"] = diagResult{false, "", "Not installed", 3, true}
		}

		// Check Tilt (important)
		fmt.Print("⏳ Tilt : ")
		logger.Debug("Checking Tilt installation...")
		if version, err := execCommand("tilt", "version"); err == nil {
			truncVersion := strings.Split(version, "\n")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Tilt installed: %s", truncVersion)
			results["tilt"] = diagResult{true, truncVersion, "OK", 2, false}
		} else {
			fmt.Println("❌ Not installed or not accessible")
			fmt.Println("   👉 To install Tilt: https://docs.tilt.dev/install.html")
			logger.Warning("Tilt not installed or not accessible")
			results["tilt"] = diagResult{false, "", "Not installed", 2, false}
		}

		fmt.Println("\n📋 Checking development tools:")
		logger.Debug("Checking development tools...")

		// Check Java (optional)
		fmt.Print("⏳ Java : ")
		logger.Debug("Checking Java installation...")
		if version, err := execCommand("java", "-version"); err == nil {
			truncVersion := strings.Split(version, "\n")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Java installed: %s", truncVersion)
			results["java"] = diagResult{true, truncVersion, "OK", 1, false}
		} else {
			fmt.Println("❌ Not installed or not accessible")
			logger.Warning("Java not installed or not accessible")
			results["java"] = diagResult{false, "", "Not installed", 1, false}
		}

		// Check Maven (optional)
		logger.Debug("Checking Maven installation...")
		fmt.Print("⏳ Maven : ")
		if version, err := execCommand("mvn", "--version"); err == nil {
			truncVersion := strings.Split(version, "\n")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Maven installed: %s", truncVersion)
			results["maven"] = diagResult{true, truncVersion, "OK", 1, false}
		} else {
			fmt.Println("⚠️ Not installed (required for some projects)")
			logger.Warning("Maven not installed")
			results["maven"] = diagResult{false, "", "Not installed", 1, false}
		}

		logger.Debug("Checking Gradle installation...")
		// Check Gradle
		fmt.Print("⏳ Gradle : ")
		if version, err := execCommand("gradle", "--version"); err == nil {
			truncVersion := strings.Split(strings.Split(version, "\n")[0], "----")[0]
			fmt.Printf("✅ %s\n", truncVersion)
			logger.Info("Gradle installed: %s", truncVersion)
			results["gradle"] = diagResult{true, truncVersion, "OK", 1, false}
		} else {
			fmt.Println("⚠️ Not installed (required for some projects)")
			logger.Warning("Gradle not installed")
			results["gradle"] = diagResult{false, "", "Not installed", 1, false}
		}

		logger.Debug("Checking project files...")
		// Check current project
		fmt.Println("\n📋 Current project:")

		// Check if required files exist
		projectFiles := checkProjectFiles()

		// Analyze project health locally
		projectHealth := func(diagnostics map[string]diagResult) int {
			total := 0
			max := 0

			// Sum the weights of installed components
			for tool, result := range diagnostics {
				weight := result.weight
				max += weight

				if result.installed {
					if result.detail == "OK" {
						total += weight
					} else if !result.required {
						// If not required, we still count partial points
						total += weight / 2
					}
				}

				logger.Debug("Evaluation of %s: installed=%t, weight=%d, detail=%s, partial score=%d/%d",
					tool, result.installed, weight, result.detail, total, max)
			}

			// Convert to percentage
			if max == 0 {
				return 0
			}
			return total * 100 / max
		}(results)

		// Display recommendations
		fmt.Println("\n📋 Recommendations:")
		if !results["docker"].installed {
			fmt.Println("❗ Docker is required: https://docs.docker.com/get-docker/")
			logger.Error("Docker missing - installation required")
		} else if results["docker"].detail != "OK" {
			fmt.Println("⚠️ Make sure the Docker daemon is running")
			logger.Warning("Problem with Docker: %s", results["docker"].detail)
		}

		if !results["docker-compose"].installed {
			fmt.Println("❗ Docker Compose is required: https://docs.docker.com/compose/install/")
			logger.Error("Docker Compose missing - installation required")
		}

		if !results["tilt"].installed {
			fmt.Println("❗ Tilt is strongly recommended: https://docs.tilt.dev/install.html")
			logger.Warning("Tilt missing - installation recommended")
		}

		// Specific recommendations for Java developers
		if len(projectFiles["java"]) > 0 && !results["java"].installed {
			fmt.Println("⚠️ Java files were detected but Java is not installed")
			logger.Warning("Java required for this project but not installed")
		}

		if len(projectFiles["maven"]) > 0 && !results["maven"].installed {
			fmt.Println("⚠️ Maven files were detected but Maven is not installed")
			logger.Warning("Maven required for this project but not installed")
		}

		if len(projectFiles["gradle"]) > 0 && !results["gradle"].installed {
			fmt.Println("⚠️ Gradle files were detected but Gradle is not installed")
			logger.Warning("Gradle required for this project but not installed")
		}

		// Display the overall project health
		fmt.Println("\n📊 Overall health:", healthToEmoji(projectHealth))
		logger.Info("Overall project health: %d%%", projectHealth)

		fmt.Println("\n🔧 Available commands:")
		fmt.Println("▶️ turbotilt init   : Initialize a project")
		fmt.Println("▶️ turbotilt up     : Start the environment")
		fmt.Println("▶️ turbotilt stop   : Stop the environment")
		fmt.Println("▶️ turbotilt doctor : Check configuration")

		// Calculate and display execution time
		duration := time.Since(startTime)
		fmt.Printf("\n⏱️ Diagnostic completed in %.2f seconds\n", duration.Seconds())

		if logToFile {
			fmt.Printf("📄 Log saved to: %s\n", logFilePath)
		}

		logger.Info("Diagnostic terminé en %.2f secondes", duration.Seconds())
		logger.Debug("Doctor command completed")
	},
}

// Function removed because it was replaced by a local implementation

// healthToEmoji converts a health score to an emoji representation with a progress bar
func healthToEmoji(health int) string {
	var emoji, grade, barGraph string

	// Determine the grade
	switch {
	case health >= 90:
		emoji = "✅"
		grade = "Excellent"
	case health >= 70:
		emoji = "🟢"
		grade = "Good"
	case health >= 50:
		emoji = "🟡"
		grade = "Average"
	case health >= 30:
		emoji = "🟠"
		grade = "Problematic"
	default:
		emoji = "🔴"
		grade = "Critical"
	}

	// Create a visual progress bar
	completed := health / 10
	remaining := 10 - completed

	barGraph = strings.Repeat("█", completed) + strings.Repeat("░", remaining)

	return fmt.Sprintf("%s %s (%d%%) %s", emoji, grade, health, barGraph)
}

// execCommand executes a command and returns its output
func execCommand(command string, args ...string) (string, error) {
	logger.Debug("Executing command: %s %s", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil && verbose {
		logger.Debug("Execution error: %v", err)
	}

	// For commands like 'java -version' that write to stderr
	outputStr := string(output)
	if outputStr == "" && err == nil {
		// Some commands might not produce output but succeed
		outputStr = "OK (no output)"
	}

	return outputStr, err
}

func checkProjectFiles() map[string][]string {
	// Map to store files found by category
	foundFiles := map[string][]string{
		"java":   {},
		"maven":  {},
		"gradle": {},
		"docker": {},
		"tilt":   {},
	}

	fileChecks := []struct {
		path     string
		desc     string
		category string
	}{
		{"pom.xml", "Maven", "maven"},
		{"build.gradle", "Gradle", "gradle"},
		{"build.gradle.kts", "Gradle Kotlin", "gradle"},
		{"Dockerfile", "Docker", "docker"},
		{"docker-compose.yml", "Docker Compose", "docker"},
		{"Tiltfile", "Tilt", "tilt"},
		{"src/main/java", "Java Source", "java"},
		{".mvn", "Maven Wrapper", "maven"},
		{"gradlew", "Gradle Wrapper", "gradle"},
	}

	fmt.Println("Checking project files:")

	for _, check := range fileChecks {
		if _, err := os.Stat(check.path); err == nil {
			fmt.Printf("  ✅ %s found (%s)\n", check.path, check.desc)
			logger.Info("File found: %s (%s)", check.path, check.desc)
			foundFiles[check.category] = append(foundFiles[check.category], check.path)
		} else if verbose {
			fmt.Printf("  ❌ %s not found\n", check.path)
			logger.Debug("Missing file: %s", check.path)
		}
	}

	if len(foundFiles["maven"]) > 0 {
		fmt.Println("  📄 Maven project detected")
	}
	if len(foundFiles["gradle"]) > 0 {
		fmt.Println("  📄 Gradle project detected")
	}
	if len(foundFiles["docker"]) > 0 {
		fmt.Println("  📦 Docker configuration detected")
	}
	if len(foundFiles["tilt"]) > 0 {
		fmt.Println("  🚀 Tilt configuration detected")
	}
	return foundFiles
}

// validateManifestFile validates the turbotilt.yaml manifest file
func validateManifestFile() {
	fmt.Println("🔍 Validating manifest...")

	// Find the manifest
	configPath, isManifest, err := config.FindConfiguration()
	if err != nil || !isManifest {
		fmt.Println("❌ turbotilt.yaml manifest not found")
		return
	}

	fmt.Printf("📄 Manifest found: %s\n", configPath)

	// Load and validate the manifest
	manifest, err := config.LoadManifest(configPath)
	if err != nil {
		fmt.Printf("❌ Validation error: %v\n", err)
		return
	}

	// Display a summary of the manifest
	fmt.Println("✅ The manifest is valid!")
	fmt.Printf("📊 Contains %d service(s):\n", len(manifest.Services))

	// Display service details
	for i, service := range manifest.Services {
		fmt.Printf("   [%d] %s\n", i+1, service.Name)

		if service.Runtime != "" {
			fmt.Printf("       - Type: Application (%s)\n", service.Runtime)
			fmt.Printf("       - Path: %s\n", service.Path)
			fmt.Printf("       - Port: %s\n", service.Port)
			fmt.Printf("       - Java: %s\n", service.Java)
		} else if service.Type != "" {
			fmt.Printf("       - Type: Dependent service (%s)\n", service.Type)
			if service.Version != "" {
				fmt.Printf("       - Version: %s\n", service.Version)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(doctorCmd)

	doctorCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display detailed information")
	doctorCmd.Flags().BoolVarP(&logToFile, "log", "l", false, "Save results to a log file")
	doctorCmd.Flags().StringVar(&logFilePath, "log-file", "", "Path to the log file to use")
	doctorCmd.Flags().BoolVar(&showAllInfo, "all", false, "Show all information")
	doctorCmd.Flags().BoolVar(&showSummary, "summary", false, "Show summary only")
	doctorCmd.Flags().BoolVar(&validateManifest, "validate-manifest", false, "Validate the syntax of turbotilt.yaml manifest")
}
