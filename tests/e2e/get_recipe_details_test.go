package e2e_test

import (
	"context"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/get_recipe_details"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/create_recipe"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_recipe"
	get_details_tool "github.com/compilercomplied/tandoor-mcp/src/tools/get_recipe_details"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestGetRecipeDetailsE2E(t *testing.T) {
	// Arrange: Create a recipe first
	ctx := context.Background()

	createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_tandoor_recipe", create_tandoor_recipe.Args{
		Name:        "Details Test Recipe",
		Description: "A recipe for testing details",
	})
	AssertToolSuccess(t, createRes, createErr)

	createdRecipe := ParseToolResponse[create_recipe.RecipeResponse](t, createRes)

	// Act: Get its details
	res, err := infra.CallTool(ctx, fixture.Client, "get_recipe_details", get_details_tool.Args{
		RecipeID: createdRecipe.ID,
	})
	
	// Assert
	AssertToolSuccess(t, res, err)
	recipeDetails := ParseToolResponse[get_recipe_details.RecipeResponse](t, res)

	if recipeDetails.ID != createdRecipe.ID {
		t.Errorf("expected ID %d, got %d", createdRecipe.ID, recipeDetails.ID)
	}
	if recipeDetails.Name != "Details Test Recipe" {
		t.Errorf("expected Name 'Details Test Recipe', got %q", recipeDetails.Name)
	}
}
