package get_recipe_details

import (
	"context"
	"fmt"
	"net/url"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

func Do(ctx context.Context, client *tandoor.Client, recipeID int) (*RecipeResponse, error) {
	path := fmt.Sprintf("/api/recipe/%d/", recipeID)
	
	// Send empty url.Values as we don't need query params
	res, err := tandoor.Request[RecipeResponse](ctx, client, "GET", path, url.Values{}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe details: %w", err)
	}

	return res, nil
}
