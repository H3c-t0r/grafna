package docker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// AllArchs is a list of all supported Docker image architectures.
var AllArchs = []string{"amd64", "armv7", "arm64"}

// Init initializes the OS for Docker image building.
func Init() error {
	// Necessary for cross-platform builds
	if err := os.Setenv("DOCKER_BUILDKIT", "1"); err != nil {
		log.Println("error setting DOCKER_BUILDKIT environment variable:", err)
	}

	// Enable execution of Docker images for other architectures
	//nolint:gosec
	cmd := exec.Command("docker", "run", "--privileged", "--rm",
		"tonistiigi/binfmt", "--install", "all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable execution of cross-platform Docker images: %w\n%s", err, output)
	}
	fmt.Println("Emulators have been installed successfully!", string(output))

	return nil
}
