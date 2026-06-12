package infra

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Fixture manages the lifecycle of the Tandoor MCP server for E2E testing.
type Fixture struct {
	Client *Client
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

// SetupFixture starts the MCP server process and connects a new Client to it,
// or connects to an external server if EXTERNAL_MCP_URL is provided.
func SetupFixture() (*Fixture, error) {
	if extURL := os.Getenv("EXTERNAL_MCP_URL"); extURL != "" {
		client, err := NewClient(context.Background(), extURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create external MCP client: %w", err)
		}
		return &Fixture{
			Client: client,
		}, nil
	}

	// Find the repository root to run the main.go
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working dir: %w", err)
	}

	// Traverse up to find go.mod to ensure we are at the root
	repoRoot := wd
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(filepath.Join(repoRoot, "go.mod")); err == nil {
			break
		}
		repoRoot = filepath.Dir(repoRoot)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Build the binary
	tmpDir, err := os.MkdirTemp("", "mcp-e2e")
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	tmpBin := filepath.Join(tmpDir, "tandoor-mcp")
	buildCmd := exec.Command("go", "build", "-o", tmpBin, "./src/")
	buildCmd.Dir = repoRoot
	if err := buildCmd.Run(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to build MCP server: %w", err)
	}

	// Start the MCP server using SSE transport
	cmd := exec.CommandContext(ctx, tmpBin, "-transport", "sse", "-port", "8081")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(),
		"LOG_FORMAT=plain",
		"LOG_HTTP_BODY=false",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start MCP server: %w", err)
	}

	// Wait for the server to bind to the port
	time.Sleep(3 * time.Second)

	client, err := NewClient(context.Background(), "http://localhost:8081/sse")
	if err != nil {
		cancel()
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}

	return &Fixture{
		Client: client,
		cmd:    cmd,
		cancel: cancel,
	}, nil
}

// Teardown closes the client and terminates the MCP server process.
func (f *Fixture) Teardown() {
	if f.Client != nil {
		f.Client.Close()
	}
	if f.cancel != nil {
		f.cancel()
	}
	if f.cmd != nil && f.cmd.Process != nil {
		f.cmd.Process.Kill()
		f.cmd.Wait()
	}
}
