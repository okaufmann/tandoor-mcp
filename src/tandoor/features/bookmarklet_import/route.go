package bookmarklet_import

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create creates a new bookmarklet import by calling POST /api/bookmarklet-import/
func Create(ctx context.Context, c *tandoor.Client, params BookmarkletImportParam) (*BookmarkletImportResponse, error) {
	res, err := tandoor.Request[BookmarkletImportResponse](ctx, c, "POST", "/api/bookmarklet-import/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create bookmarklet import: %w", err)
	}
	return res, nil
}

// List lists bookmarklet imports by calling GET /api/bookmarklet-import/
func List(ctx context.Context, c *tandoor.Client, params ListParams) (*BookmarkletImportListResponse, error) {
	qb := tandoor.NewQuery().
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[BookmarkletImportListResponse](ctx, c, "GET", "/api/bookmarklet-import/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookmarklet imports: %w", err)
	}
	return res, nil
}
