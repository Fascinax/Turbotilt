package runtime

import (
	"os/exec"
)

func TiltUp() error {
	cmd := exec.Command("tilt", "up")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func ComposeUp() error {
	cmd := exec.Command("docker-compose", "up")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
