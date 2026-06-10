package e2e_test

import (
	"context"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/get_recipes"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGetRecipesE2E(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	res, err := infra.CallTool(ctx, fixture.Client, "get_recipes", get_recipes.GetRecipesParams{
		Query: "Tandoori Chicken",
	})

	// Assert
	AssertToolSuccess(t, res, err)

	if len(res.Content) == 0 {
		t.Fatalf("expected at least 1 recipe, but got empty content")
	}

	hasText := false
	for _, c := range res.Content {
		if textContent, ok := c.(*mcp_sdk.TextContent); ok {
			hasText = true
			if textContent.Text == "No recipes found matching the criteria." {
				t.Fatalf("expected at least 1 recipe, but found none")
			}
		}
	}

	if !hasText {
		t.Fatalf("expected text content in response")
	}
}
