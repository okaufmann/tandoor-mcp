package keyword

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create calls POST /api/keyword/
func Create(ctx context.Context, c *tandoor.Client, params KeywordParam) (*KeywordResponse, error) {
	res, err := tandoor.Request[KeywordResponse](ctx, c, "POST", "/api/keyword/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyword: %w", err)
	}
	return res, nil
}

// List calls GET /api/keyword/
func List(ctx context.Context, c *tandoor.Client, params ListParams) (*KeywordListResponse, error) {
	qb := tandoor.NewQuery().
		Add("query", params.Query).
		Add("root", params.Root).
		Add("tree", params.Tree).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[KeywordListResponse](ctx, c, "GET", "/api/keyword/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list keywords: %w", err)
	}
	return res, nil
}
