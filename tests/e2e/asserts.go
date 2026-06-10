package e2e_test

import (
	"encoding/json"
	"testing"

	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ParseToolResponse extracts the text content from a CallToolResult and unmarshals it into T.
func ParseToolResponse[T any](t *testing.T, res *mcp_sdk.CallToolResult) *T {
	t.Helper()

	if len(res.Content) == 0 {
		t.Fatalf("expected at least 1 content item")
	}

	for _, c := range res.Content {
		if tc, ok := c.(*mcp_sdk.TextContent); ok {
			if len(tc.Text) == 0 {
				t.Fatalf("expected text content")
			}
			var result T
			if err := json.Unmarshal([]byte(tc.Text), &result); err != nil {
				t.Fatalf("failed to parse tool response json: %v\nJSON payload: %s", err, tc.Text)
			}
			return &result
		}
	}

	t.Fatalf("expected text content in response")
	return nil
}

// AssertToolSuccess is a common testing utility to check if a tool call was successful.
func AssertToolSuccess(t *testing.T, res *mcp_sdk.CallToolResult, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("failed to call tool: %v", err)
	}

	if res.IsError {
		var errMsgs []string
		for _, c := range res.Content {
			if tc, ok := c.(*mcp_sdk.TextContent); ok {
				errMsgs = append(errMsgs, tc.Text)
			}
		}
		t.Fatalf("tool returned error: %v", errMsgs)
	}
}
