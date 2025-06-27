package render

import (
	"os"
)

func GenerateTiltfile() error {
	f, err := os.Create("Tiltfile")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("docker_build('app', '.', dockerfile='Dockerfile')\nlive_update = [ sync('./src', '/app/src') ]\nk8s_yaml('docker-compose.yml')\n")
	return err
}
