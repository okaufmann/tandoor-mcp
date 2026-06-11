package cooklog

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create registers a cooking event by calling POST /api/cook-log/
func Create(ctx context.Context, c *tandoor.Client, params CookLogParam) (*CookLogResponse, error) {
	res, err := tandoor.Request[CookLogResponse](ctx, c, "POST", "/api/cook-log/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create cook log: %w", err)
	}
	return res, nil
}

// List retrieves a list of cook logs by calling GET /api/cook-log/
func List(ctx context.Context, c *tandoor.Client, params ListParams) (*CookLogListResponse, error) {
	qb := tandoor.NewQuery().
		Add("recipe", params.Recipe).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[CookLogListResponse](ctx, c, "GET", "/api/cook-log/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list cook logs: %w", err)
	}
	return res, nil
}
