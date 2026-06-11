package recipe_import

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create creates a new recipe import config by calling POST /api/recipe-import/
func Create(ctx context.Context, c *tandoor.Client, params RecipeImportParam) (*RecipeImportResponse, error) {
	res, err := tandoor.Request[RecipeImportResponse](ctx, c, "POST", "/api/recipe-import/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create recipe import: %w", err)
	}
	return res, nil
}

// List lists recipe imports by calling GET /api/recipe-import/
func List(ctx context.Context, c *tandoor.Client, params ListParams) (*RecipeImportListResponse, error) {
	qb := tandoor.NewQuery().
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[RecipeImportListResponse](ctx, c, "GET", "/api/recipe-import/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list recipe imports: %w", err)
	}
	return res, nil
}
