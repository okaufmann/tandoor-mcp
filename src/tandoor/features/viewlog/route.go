package viewlog

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create registers a viewing event by calling POST /api/view-log/
func Create(ctx context.Context, c *tandoor.Client, params ViewLogParam) (*ViewLogResponse, error) {
	res, err := tandoor.Request[ViewLogResponse](ctx, c, "POST", "/api/view-log/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create view log: %w", err)
	}
	return res, nil
}

// List retrieves a list of view logs by calling GET /api/view-log/
func List(ctx context.Context, c *tandoor.Client, params ListParams) (*ViewLogListResponse, error) {
	qb := tandoor.NewQuery().
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[ViewLogListResponse](ctx, c, "GET", "/api/view-log/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list view logs: %w", err)
	}
	return res, nil
}
