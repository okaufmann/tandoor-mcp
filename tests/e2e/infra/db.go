package infra

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// PurgeAndSeedDatabase uses the Docker socket to find the tandoor-setup container
// and runs it again. Since the setup container's CMD is to execute the python
// seed script, this cleanly drops all tables, restarts identities (IDs), and
// seeds the required test fixture data (User, Space, Token, Recipe #1).
func PurgeAndSeedDatabase() error {
	// Find the container ID of tandoor-setup
	cmd := exec.Command("docker", "ps", "-a", "-q", "--filter", "label=com.docker.compose.service=tandoor-setup")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil || strings.TrimSpace(out.String()) == "" {
		return fmt.Errorf("failed to find tandoor-setup container: %v", err)
	}
	
	containerID := strings.Split(strings.TrimSpace(out.String()), "\n")[0]

	// Restart the container and attach to wait for it to finish
	startCmd := exec.Command("docker", "start", "-a", containerID)
	if err := startCmd.Run(); err != nil {
		return fmt.Errorf("failed to run tandoor-setup container: %v", err)
	}

	return nil
}
