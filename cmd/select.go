package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/your-username/turbotilt/internal/config"
	"github.com/your-username/turbotilt/internal/i18n"
	"github.com/your-username/turbotilt/internal/logger"
	"github.com/your-username/turbotilt/internal/scan"
)

var (
	selectOutputFile string
	selectCreateConfig bool
	selectLaunchSelected bool
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select [directory]",
	Short: "Detect microservices and select which ones to launch",
	Long: `Scan a directory for microservices and let you select which ones to launch.
For example:
  turbotilt select ./my-project
  turbotilt select ./my-project --output turbotilt.yaml
  turbotilt select ./my-project --create-config --launch`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		t := i18n.GetTranslator()
		log := logger.GetLogger()
		
		// Get the directory to scan
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		
		// Get absolute path
		absDir, err := filepath.Abs(dir)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		
		log.Infof(t.Tr("Scanning directory: %s"), absDir)
		
		// Scan for services
		detector := scan.NewDetector()
		services, err := detector.ScanDirectory(absDir)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		
		if len(services) == 0 {
			log.Info(t.Tr("No microservices detected in the directory"))
			os.Exit(0)
		}
		
		log.Infof(t.Tr("Found %d microservices:"), len(services))
		for i, service := range services {
			fmt.Printf("[%d] %s (%s) - %s\n", i+1, service.Name, service.Type, service.Path)
		}
		
		// Ask user to select services
		fmt.Println(t.Tr("\nSelect services to include (comma-separated numbers, 'all' for all services):"))
		var input string
		fmt.Scanln(&input)
		
		selectedServices := []scan.Service{}
		if input == "all" {
			selectedServices = services
		} else {
			// Parse selection
			selections := strings.Split(input, ",")
			for _, s := range selections {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				
				idx, err := strconv.Atoi(s)
				if err != nil || idx < 1 || idx > len(services) {
					log.Errorf(t.Tr("Invalid selection: %s"), s)
					continue
				}
				
				selectedServices = append(selectedServices, services[idx-1])
			}
		}
		
		if len(selectedServices) == 0 {
			log.Info(t.Tr("No services selected, exiting"))
			os.Exit(0)
		}
		
		log.Infof(t.Tr("Selected %d services:"), len(selectedServices))
		for _, service := range selectedServices {
			log.Infof("- %s (%s)", service.Name, service.Type)
		}
		
		// Create config if requested
		if selectCreateConfig || selectOutputFile != "" {
			// Create a manifest from selected services
			manifest := config.Manifest{
				Version: "1",
				Services: map[string]config.Service{},
			}
			
			for _, service := range selectedServices {
				servicePath, err := filepath.Rel(absDir, service.Path)
				if err != nil {
					servicePath = service.Path
				}
				
				manifest.Services[service.Name] = config.Service{
					Path: servicePath,
					Type: service.Type,
				}
			}
			
			// Marshal to YAML
			yamlData, err := yaml.Marshal(manifest)
			if err != nil {
				log.Errorf(t.Tr("Error creating manifest: %v"), err)
				os.Exit(1)
			}
			
			// Determine output file
			outputFile := "turbotilt.yaml"
			if selectOutputFile != "" {
				outputFile = selectOutputFile
			}
			
			// Write to file
			outputPath := filepath.Join(absDir, outputFile)
			err = os.WriteFile(outputPath, yamlData, 0644)
			if err != nil {
				log.Errorf(t.Tr("Error writing manifest to %s: %v"), outputPath, err)
				os.Exit(1)
			}
			
			log.Infof(t.Tr("Created manifest at %s"), outputPath)
		}
		
		// Launch selected services if requested
		if selectLaunchSelected {
			// If we created a config, use it; otherwise create a temporary one
			configFile := ""
			tempFile := false
			
			if selectCreateConfig || selectOutputFile != "" {
				configFile = filepath.Join(absDir, selectOutputFile)
				if configFile == "" {
					configFile = filepath.Join(absDir, "turbotilt.yaml")
				}
			} else {
				// Create a temporary config
				manifest := config.Manifest{
					Version: "1",
					Services: map[string]config.Service{},
				}
				
				for _, service := range selectedServices {
					servicePath, err := filepath.Rel(absDir, service.Path)
					if err != nil {
						servicePath = service.Path
					}
					
					manifest.Services[service.Name] = config.Service{
						Path: servicePath,
						Type: service.Type,
					}
				}
				
				// Marshal to YAML
				yamlData, err := yaml.Marshal(manifest)
				if err != nil {
					log.Errorf(t.Tr("Error creating temporary manifest: %v"), err)
					os.Exit(1)
				}
				
				// Create temp file
				tmpFile, err := os.CreateTemp(absDir, "turbotilt-*.yaml")
				if err != nil {
					log.Errorf(t.Tr("Error creating temporary manifest file: %v"), err)
					os.Exit(1)
				}
				
				_, err = tmpFile.Write(yamlData)
				if err != nil {
					log.Errorf(t.Tr("Error writing to temporary manifest file: %v"), err)
					os.Exit(1)
				}
				
				tmpFile.Close()
				configFile = tmpFile.Name()
				tempFile = true
				
				// Register cleanup handler for the temp file
				log.Infof(t.Tr("Created temporary manifest at %s"), configFile)
				defer func() {
					os.Remove(configFile)
					log.Infof(t.Tr("Removed temporary manifest"))
				}()
			}
			
			// Change directory to the scanned directory
			originalDir, err := os.Getwd()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			
			err = os.Chdir(absDir)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			
			// Run up command
			log.Infof(t.Tr("Launching selected services..."))
			upCmd.Run(cmd, []string{"-f", configFile})
			
			// Return to original directory
			os.Chdir(originalDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
	
	// Flags
	selectCmd.Flags().StringVarP(&selectOutputFile, "output", "o", "", "Output file for the generated turbotilt.yaml")
	selectCmd.Flags().BoolVarP(&selectCreateConfig, "create-config", "c", false, "Create a turbotilt.yaml file with the selected services")
	selectCmd.Flags().BoolVarP(&selectLaunchSelected, "launch", "l", false, "Launch the selected services after selection")
}
