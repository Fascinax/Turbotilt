package config

import (
	"fmt"
	"turbotilt/internal/render"
	"turbotilt/internal/scan"
)

// GenerateFilesFromMemory génère des fichiers Dockerfile, docker-compose.yml et Tiltfile
// à partir de la configuration stockée en mémoire
func GenerateFilesFromMemory() error {
	store := GetMemoryStore()
	if !store.HasSelectedServices() {
		return fmt.Errorf("no services found in memory")
	}

	// Récupérer le manifeste depuis la mémoire
	manifest := GetManifestFromMemory()

	// Cas d'un seul service
	if len(manifest.Services) == 1 {
		service := manifest.Services[0]
		opts, err := ConvertManifestToRenderOptions(service)
		if err != nil {
			return fmt.Errorf("error converting service to render options: %w", err)
		}

		// Génération du Dockerfile
		if err := render.GenerateDockerfile(*opts); err != nil {
			return fmt.Errorf("error generating Dockerfile: %w", err)
		}

		// Génération du docker-compose.yml
		if err := render.GenerateCompose(*opts); err != nil {
			return fmt.Errorf("error generating docker-compose.yml: %w", err)
		}

		// Génération du Tiltfile
		if err := render.GenerateTiltfile(*opts); err != nil {
			return fmt.Errorf("error generating Tiltfile: %w", err)
		}

		return nil
	}

	// Cas multi-services
	// Convertir les services du Manifest en options de rendu
	serviceList := render.ServiceList{
		Services: []render.Options{},
	}

	// Identifier les services d'application vs les services dépendants
	appServices := []ManifestService{}
	depServices := []scan.ServiceConfig{}

	for _, service := range manifest.Services {
		// Vérifier si c'est un service d'application (avec runtime) ou un service dépendant (avec type)
		if service.Runtime != "" {
			appServices = append(appServices, service)
		} else if service.Type != "" {
			// Convertir en scan.ServiceConfig
			depService := scan.ServiceConfig{
				Type:        scan.ServiceType(service.Type),
				Version:     service.Version,
				Port:        service.Port,
				Credentials: service.Env,
			}
			depServices = append(depServices, depService)
		}
	}

	// Si aucun service d'application n'est trouvé, retourner une erreur
	if len(appServices) == 0 {
		return fmt.Errorf("no application services found in memory")
	}

	// Préparer la liste des services pour le rendu multi-services
	for _, service := range appServices {
		opts, err := ConvertManifestToRenderOptions(service)
		if err != nil {
			return fmt.Errorf("error converting service %s to render options: %w", service.Name, err)
		}

		// Ajouter les services dépendants à chaque service d'application
		opts.Services = depServices
		serviceList.Services = append(serviceList.Services, *opts)
	}

	// Génération du Tiltfile multi-services
	if err := render.GenerateMultiServiceTiltfile(serviceList); err != nil {
		return fmt.Errorf("error generating multi-service Tiltfile: %w", err)
	}

	// TODO: Génération de docker-compose.yml multi-services
	// Cette partie dépend de l'implémentation existante des templates multi-services
	// Générer chaque Dockerfile dans le dossier approprié
	for _, service := range appServices {
		opts, err := ConvertManifestToRenderOptions(service)
		if err != nil {
			continue
		}

		// Génération du Dockerfile dans le dossier du service
		if err := render.GenerateDockerfile(*opts); err != nil {
			return fmt.Errorf("error generating Dockerfile for service %s: %w", service.Name, err)
		}
	}

	return nil
}
