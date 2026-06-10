package create_recipe

import (
	"context"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

func Create(ctx context.Context, c *tandoor.Client, params CreateRecipeParams) (*RecipeResponse, error) {
	return tandoor.Request[RecipeResponse](ctx, c, "POST", "/api/recipe/", nil, params)
}
