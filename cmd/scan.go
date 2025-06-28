package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"turbotilt/internal/logger"
	"turbotilt/internal/scan"

	"github.com/spf13/cobra"
)

var (
	scanPath          string
	maxDepth          int
	nonInteractive    bool
	scanAutoGenerate  bool
	scanOutputFile    string
	scanOutputFormat  string
	scanIncludeAll    bool
	scanSkipFramework bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan directory for microservices",
	Long:  `Scan a directory for microservices and allow selection of which ones to run.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debugMode {
			logger.SetLevel(logger.DEBUG)
			logger.Debug("Debug mode enabled")
		}

		fmt.Println("🔍 Scanning for microservices...")

		// Use current directory if no path is specified
		if scanPath == "" {
			var err error
			scanPath, err = os.Getwd()
			if err != nil {
				fmt.Printf("❌ Error getting current directory: %v\n", err)
				return
			}
		}

		// Convert relative path to absolute
		if !filepath.IsAbs(scanPath) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("❌ Error getting current directory: %v\n", err)
				return
			}
			scanPath = filepath.Join(cwd, scanPath)
		}

		fmt.Printf("📂 Scanning directory: %s\n", scanPath)
		
		// Find potential microservice directories
		microservices, err := findMicroservices(scanPath, maxDepth)
		if err != nil {
			fmt.Printf("❌ Error scanning for microservices: %v\n", err)
			return
		}

		if len(microservices) == 0 {
			fmt.Println("❌ No microservices found in the specified directory.")
			return
		}

		fmt.Printf("✅ Found %d potential microservices\n", len(microservices))

		// User selection if not in non-interactive mode
		selectedServices := microservices
		if !nonInteractive && !scanIncludeAll {
			selectedServices = selectMicroservices(microservices)
		}

		if len(selectedServices) == 0 {
			fmt.Println("ℹ️ No microservices selected. Exiting.")
			return
		}

		fmt.Printf("✅ Selected %d microservices\n", len(selectedServices))

		// Generate manifest if requested
		if scanAutoGenerate {
			if err := generateManifest(selectedServices, scanOutputFile, scanOutputFormat); err != nil {
				fmt.Printf("❌ Error generating manifest: %v\n", err)
				return
			}
			fmt.Printf("✅ Generated manifest: %s\n", scanOutputFile)
		}

		// Print next steps
		fmt.Println("\n📋 Next steps:")
		if scanAutoGenerate {
			fmt.Printf("1. Review the generated manifest: %s\n", scanOutputFile)
			fmt.Println("2. Run 'turbotilt up' to start the selected services")
		} else {
			fmt.Println("1. Run 'turbotilt init' to generate configurations for the selected services")
			fmt.Println("2. Run 'turbotilt up' to start the services")
		}
	},
}

// findMicroservices recursively scans directories up to maxDepth looking for potential microservices
func findMicroservices(rootPath string, maxDepth int) ([]MicroserviceInfo, error) {
	var result []MicroserviceInfo

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == rootPath {
			return nil
		}

		// Check if we've exceeded max depth
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}
		depth := len(strings.Split(relPath, string(os.PathSeparator)))
		if depth > maxDepth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process directories
		if !info.IsDir() {
			return nil
		}

		// Check if directory has indicators of a microservice
		indicators := detectMicroserviceIndicators(path)
		if len(indicators) > 0 {
			ms := MicroserviceInfo{
				Path:       path,
				Name:       filepath.Base(path),
				Indicators: indicators,
				Framework:  "unknown",
			}

			// Detect framework if not skipping
			if !scanSkipFramework {
				// Change to the directory temporarily to use the framework detector
				currentDir, err := os.Getwd()
				if err == nil {
					os.Chdir(path)
					defer os.Chdir(currentDir)
					
					framework, err := scan.DetectFramework()
					if err == nil && framework != "" {
						ms.Framework = framework
					}
				}
			}

			result = append(result, ms)
			
			// Skip subdirectories of a detected microservice
			return filepath.SkipDir
		}

		return nil
	})

	return result, err
}

// detectMicroserviceIndicators checks a directory for files that indicate it might be a microservice
func detectMicroserviceIndicators(path string) []string {
	var indicators []string

	// Common files that indicate a microservice
	indicatorFiles := map[string]string{
		"pom.xml":                    "Maven project",
		"build.gradle":               "Gradle project",
		"package.json":               "Node.js project",
		"requirements.txt":           "Python project",
		"go.mod":                     "Go project",
		"Dockerfile":                 "Docker project",
		"docker-compose.yml":         "Docker Compose project",
		"application.properties":     "Spring Boot configuration",
		"application.yml":            "Spring Boot configuration",
		"application.yaml":           "Spring Boot configuration",
		"quarkus.properties":         "Quarkus configuration",
		"micronaut-application.yml":  "Micronaut configuration",
		"angular.json":               "Angular project",
	}

	for file, description := range indicatorFiles {
		if _, err := os.Stat(filepath.Join(path, file)); err == nil {
			indicators = append(indicators, description)
		}
	}

	return indicators
}

// MicroserviceInfo contains information about a detected microservice
type MicroserviceInfo struct {
	Path       string
	Name       string
	Indicators []string
	Framework  string
}

// selectMicroservices presents a menu to select which microservices to include
func selectMicroservices(microservices []MicroserviceInfo) []MicroserviceInfo {
	fmt.Println("\n📋 Select microservices to include (comma-separated numbers, or 'all'):")
	
	for i, ms := range microservices {
		indicators := strings.Join(ms.Indicators, ", ")
		fmt.Printf("%d. %s [%s] - %s\n", i+1, ms.Name, ms.Framework, indicators)
	}
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nSelection: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	// Return all if "all" is specified
	if strings.ToLower(input) == "all" {
		return microservices
	}
	
	// Parse selection
	var selected []MicroserviceInfo
	
	// Split by comma and process each number
	for _, numStr := range strings.Split(input, ",") {
		numStr = strings.TrimSpace(numStr)
		if numStr == "" {
			continue
		}
		
		num, err := strconv.Atoi(numStr)
		if err != nil || num < 1 || num > len(microservices) {
			fmt.Printf("⚠️ Invalid selection: %s - skipping\n", numStr)
			continue
		}
		
		selected = append(selected, microservices[num-1])
	}
	
	return selected
}

// generateManifest creates a turbotilt.yaml manifest for the selected microservices
func generateManifest(services []MicroserviceInfo, outputFile string, format string) error {
	// If no output file specified, use default
	if outputFile == "" {
		outputFile = "turbotilt.yaml"
	}
	
	// Basic template for the manifest
	manifestContent := "# Generated by turbotilt scan\n"
	manifestContent += "services:\n"
	
	for _, service := range services {
		manifestContent += fmt.Sprintf("  - name: %s\n", service.Name)
		manifestContent += fmt.Sprintf("    path: %s\n", getRelativePath(service.Path))
		
		// Add framework-specific configuration
		switch service.Framework {
		case "spring":
			manifestContent += "    runtime: spring\n"
			manifestContent += "    build: maven\n"
			manifestContent += "    java: \"11\"\n"
			manifestContent += "    port: \"8080\"\n"
		case "quarkus":
			manifestContent += "    runtime: quarkus\n"
			manifestContent += "    build: maven\n"
			manifestContent += "    java: \"11\"\n"
			manifestContent += "    port: \"8080\"\n"
		case "micronaut":
			manifestContent += "    runtime: micronaut\n"
			manifestContent += "    build: maven\n"
			manifestContent += "    java: \"11\"\n"
			manifestContent += "    port: \"8080\"\n"
		case "node":
			manifestContent += "    type: node\n"
			manifestContent += "    port: \"3000\"\n"
		case "angular":
			manifestContent += "    type: angular\n"
			manifestContent += "    port: \"4200\"\n"
		case "python":
			manifestContent += "    type: python\n"
			manifestContent += "    port: \"5000\"\n"
		case "go":
			manifestContent += "    type: go\n"
			manifestContent += "    port: \"8000\"\n"
		default:
			manifestContent += "    # Unknown framework - please configure manually\n"
			manifestContent += "    #runtime: \n"
			manifestContent += "    #build: \n"
			manifestContent += "    #port: \n"
		}
		
		manifestContent += "    devMode: true\n"
		manifestContent += "\n"
	}
	
	// Write to file
	return os.WriteFile(outputFile, []byte(manifestContent), 0644)
}

// getRelativePath returns a path relative to the current directory
func getRelativePath(path string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return path
	}
	
	relPath, err := filepath.Rel(cwd, path)
	if err != nil {
		return path
	}
	
	return relPath
}

func init() {
	rootCmd.AddCommand(scanCmd)
	
	scanCmd.Flags().StringVarP(&scanPath, "path", "p", "", "Directory to scan (default: current directory)")
	scanCmd.Flags().IntVarP(&maxDepth, "depth", "d", 3, "Maximum directory depth to scan")
	scanCmd.Flags().BoolVarP(&nonInteractive, "non-interactive", "n", false, "Non-interactive mode (select all services)")
	scanCmd.Flags().BoolVarP(&scanAutoGenerate, "generate", "g", false, "Auto-generate turbotilt.yaml manifest")
	scanCmd.Flags().StringVarP(&scanOutputFile, "output", "o", "turbotilt.yaml", "Output file for generated manifest")
	scanCmd.Flags().StringVarP(&scanOutputFormat, "format", "f", "yaml", "Output format (yaml or json)")
	scanCmd.Flags().BoolVarP(&scanIncludeAll, "all", "a", false, "Include all detected services")
	scanCmd.Flags().BoolVarP(&scanSkipFramework, "skip-framework", "s", false, "Skip framework detection (faster)")
}
