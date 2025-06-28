package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"turbotilt/internal/config"
	"turbotilt/internal/render"
	"turbotilt/internal/scan"
)

var (
	// Options for the init command
	forceFramework   string
	port             string
	jdkVersion       string
	devMode          bool
	detectServices   bool
	generateManifest bool
	fromManifest     bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scan and generate Tiltfile & Compose",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Initializing Turbotilt...")

		// Look for an existing manifest
		configPath, isManifest, _ := config.FindConfiguration()

		// If --from-manifest is requested or if a manifest exists and --generate-manifest is not requested
		if fromManifest || (isManifest && !generateManifest) {
			if configPath == "" {
				fmt.Println("❌ No manifest found. Use --generate-manifest to create one.")
				return
			}

			fmt.Printf("📄 Using manifest %s\n", configPath)
			manifest, err := config.LoadManifest(configPath)
			if err != nil {
				fmt.Printf("❌ Error loading manifest: %v\n", err)
				return
			}

			fmt.Printf("✅ Manifest loaded with %d service(s)\n", len(manifest.Services))

			// Convert manifest services to render options
			serviceList := render.ServiceList{
				Services: []render.Options{},
			}

			for _, service := range manifest.Services {
				// Ignore dependent services (without runtime)
				if service.Runtime == "" {
					continue
				}

				opts, err := config.ConvertManifestToRenderOptions(service)
				if err != nil {
					fmt.Printf("⚠️ Warning: %v\n", err)
					continue
				}

				serviceList.Services = append(serviceList.Services, *opts)
			}

			// Generate files for a multi-service project
			if len(serviceList.Services) > 0 {
				fmt.Println("🔧 Generating configurations for a multi-service project...")

				if err := render.GenerateMultiServiceCompose(serviceList); err != nil {
					fmt.Printf("❌ Error generating docker-compose.yml: %v\n", err)
					return
				}

				if err := render.GenerateMultiServiceTiltfile(serviceList); err != nil {
					fmt.Printf("❌ Error generating Tiltfile: %v\n", err)
					return
				}

				fmt.Println("✨ Turbotilt configuration completed!")
				fmt.Println("📋 Files generated from manifest:")
				fmt.Println("   - docker-compose.yml")
				fmt.Println("   - Tiltfile")
				fmt.Println("\n▶️ To start the environment: turbotilt up")
				return
			}
		}

		// If we get here, proceed with auto-detection or CLI options

		// Detect framework or use the specified one
		framework := forceFramework
		var err error

		if framework == "" {
			framework, err = scan.DetectFramework()
			if err != nil {
				fmt.Printf("❌ Error detecting framework: %v\n", err)
				return
			}
		}

		fmt.Printf("✅ Framework detected/selected: %s\n", framework)

		// Detect services if requested
		var services []scan.ServiceConfig
		if detectServices {
			fmt.Println("🔍 Detecting dependent services...")
			services, err = scan.DetectServices()
			if err != nil {
				fmt.Printf("⚠️ Warning during service detection: %v\n", err)
			}

			// Display detected services
			if len(services) > 0 {
				fmt.Println("✅ Detected services:")
				for _, service := range services {
					fmt.Printf("   - %s\n", service.Type)
				}
			} else {
				fmt.Println("ℹ️ No dependent services detected")
			}
		}

		// Determine application name (current folder by default)
		appName := "app"
		cwd, err := os.Getwd()
		if err == nil {
			appName = filepath.Base(cwd)
		}

		// Prepare render options
		renderOpts := render.Options{
			ServiceName: appName, // Use the name to identify it in a multi-service context
			Framework:   framework,
			AppName:     appName,
			Port:        port,
			JDKVersion:  jdkVersion,
			DevMode:     devMode,
			Path:        ".",
			Services:    services,
		}

		// Generate manifest if requested
		if generateManifest {
			fmt.Println("📝 Generating turbotilt.yaml manifest...")

			// Create a configuration based on detection results
			cfg := config.Config{
				Project: config.ProjectConfig{
					Name:        appName,
					Description: "Turbotilt Project",
					Version:     "1.0.0",
				},
				Framework: config.FrameworkConfig{
					Type:       framework,
					JdkVersion: jdkVersion,
				},
				Docker: config.DockerConfig{
					Port: port,
				},
				Development: config.DevelopmentConfig{
					EnableLiveReload: devMode,
				},
				Services: []config.ServiceConfig{},
			}

			// Convert scan.ServiceConfig to config.ServiceConfig
			for _, svc := range services {
				// Generate a name based on type
				serviceName := strings.ToLower(string(svc.Type))

				configSvc := config.ServiceConfig{
					Name:        serviceName,
					Type:        string(svc.Type),
					Version:     svc.Version,
					Port:        svc.Port,
					Environment: svc.Credentials,
				}
				cfg.Services = append(cfg.Services, configSvc)
			}

			// Generate manifest from configuration
			manifest := config.GenerateManifestFromConfig(cfg)

			// Save manifest
			if err := config.SaveManifest(manifest, config.ManifestFileName); err != nil {
				fmt.Printf("❌ Error saving manifest: %v\n", err)
			} else {
				fmt.Printf("✅ Manifest %s generated successfully!\n", config.ManifestFileName)
			}
		}

		// Generate files
		if err := render.GenerateDockerfile(renderOpts); err != nil {
			fmt.Printf("❌ Error generating Dockerfile: %v\n", err)
			return
		}

		// Use the new docker-compose generator with service support
		if len(services) > 0 {
			if err := render.GenerateComposeWithServices(renderOpts); err != nil {
				fmt.Printf("❌ Error generating docker-compose.yml: %v\n", err)
				return
			}
		} else {
			if err := render.GenerateCompose(renderOpts); err != nil {
				fmt.Printf("❌ Error generating docker-compose.yml: %v\n", err)
				return
			}
		}

		if err := render.GenerateTiltfile(renderOpts); err != nil {
			fmt.Printf("❌ Error generating Tiltfile: %v\n", err)
			return
		}

		fmt.Println("✨ Turbotilt configuration completed!")
		fmt.Println("📋 Generated files:")
		fmt.Println("   - Dockerfile")
		fmt.Println("   - docker-compose.yml")
		fmt.Println("   - Tiltfile")
		if generateManifest {
			fmt.Printf("   - %s\n", config.ManifestFileName)
		}
		fmt.Println("\n▶️ To start the environment: turbotilt up")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags for the init command
	initCmd.Flags().StringVarP(&forceFramework, "framework", "f", "", "Manually specify the framework (spring, quarkus, java)")
	initCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to expose for the application")
	initCmd.Flags().StringVarP(&jdkVersion, "jdk", "j", "11", "JDK version to use")
	initCmd.Flags().BoolVarP(&devMode, "dev", "d", true, "Enable development configurations")
	initCmd.Flags().BoolVarP(&detectServices, "services", "s", true, "Detect and configure dependent services (MySQL, PostgreSQL, etc.)")
	initCmd.Flags().BoolVarP(&generateManifest, "generate-manifest", "g", false, "Generate a turbotilt.yaml manifest from detection")
	initCmd.Flags().BoolVarP(&fromManifest, "from-manifest", "m", false, "Initialize project from an existing manifest")
}
