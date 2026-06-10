package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Client is a generic client for interacting with the MCP server in E2E tests.
type Client struct {
	session *mcp_sdk.ClientSession
}

// NewClient establishes an SSE connection to the MCP server at the given endpoint.
func NewClient(ctx context.Context, endpoint string) (*Client, error) {
	impl := &mcp_sdk.Implementation{
		Name:    "tandoor-mcp-e2e-test-client",
		Version: "1.0.0",
	}

	mcpClient := mcp_sdk.NewClient(impl, nil)
	transport := &mcp_sdk.SSEClientTransport{
		Endpoint:   endpoint,
		HTTPClient: http.DefaultClient,
	}

	session, err := mcpClient.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MCP server: %w", err)
	}

	return &Client{
		session: session,
	}, nil
}

// Close closes the MCP session.
func (c *Client) Close() error {
	return c.session.Close()
}

// CallTool executes an MCP tool with a strongly typed generic payload.
func CallTool[Req any](ctx context.Context, c *Client, toolName string, args Req) (*mcp_sdk.CallToolResult, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}

	var arguments map[string]interface{}
	if err := json.Unmarshal(b, &arguments); err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments into map: %w", err)
	}

	return c.session.CallTool(ctx, &mcp_sdk.CallToolParams{
		Name:      toolName,
		Arguments: arguments,
	})
}
