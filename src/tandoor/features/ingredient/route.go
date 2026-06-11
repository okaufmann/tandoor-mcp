package ingredient

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create calls POST /api/ingredient/ to create a standalone ingredient.
func Create(ctx context.Context, c *tandoor.Client, params IngredientParam) (*IngredientResponse, error) {
	res, err := tandoor.Request[IngredientResponse](ctx, c, "POST", "/api/ingredient/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create ingredient: %w", err)
	}
	return res, nil
}
