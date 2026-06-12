package mealtype

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create creates a new meal type by calling POST /api/meal-type/
func Create(ctx context.Context, c *tandoor.Client, params MealTypeParam) (*MealTypeResponse, error) {
	res, err := tandoor.Request[MealTypeResponse](ctx, c, "POST", "/api/meal-type/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create meal type: %w", err)
	}
	return res, nil
}

// List lists meal types by calling GET /api/meal-type/
func List(ctx context.Context, c *tandoor.Client, params ListParams) (*MealTypeListResponse, error) {
	qb := tandoor.NewQuery().
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[MealTypeListResponse](ctx, c, "GET", "/api/meal-type/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list meal types: %w", err)
	}
	return res, nil
}
