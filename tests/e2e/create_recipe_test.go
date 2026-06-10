package e2e_test

import (
	"context"
	"testing"

	api_create_recipe "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/create_recipe"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_recipe"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestCreateRecipeE2E(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_recipe", create_tandoor_recipe.Args{
		Name:        "Test Recipe 1",
		Description: "A tasty recipe",
	})

	// Assert
	AssertToolSuccess(t, res, err)

	recipe := ParseToolResponse[api_create_recipe.RecipeResponse](t, res)

	if recipe.Name != "Test Recipe 1" {
		t.Errorf("expected name 'Test Recipe 1', got %q", recipe.Name)
	}
	if recipe.Description != "A tasty recipe" {
		t.Errorf("expected description 'A tasty recipe', got %q", recipe.Description)
	}
}
