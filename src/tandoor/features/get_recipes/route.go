package get_recipes

import (
	"context"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Search executes a recipe search request against the Tandoor API client
func Search(ctx context.Context, c *tandoor.Client, params GetRecipesParams) (*RecipeListResponse, error) {
	query := tandoor.NewQuery().
		Add("query", params.Query).
		Add("search", params.Search).
		Add("foods", params.Foods).
		Add("keywords", params.Keywords).
		Add("limit", params.Limit).
		Add("rating", params.Rating).
		Values()

	return tandoor.Request[RecipeListResponse](ctx, c, "GET", "/api/recipe/", query, nil)
}
