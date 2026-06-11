package ingredient_parser

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Parse calls POST /api/ingredient-parser/post/ to parse raw ingredient string(s)
func Parse(ctx context.Context, c *tandoor.Client, params IngredientParserRequest) (*IngredientParserResponse, error) {
	res, err := tandoor.Request[IngredientParserResponse](ctx, c, "POST", "/api/ingredient-parser/post/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ingredients: %w", err)
	}
	return res, nil
}
