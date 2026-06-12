package storage

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create calls POST /api/storage/
func Create(ctx context.Context, c *tandoor.Client, params StorageParam) (*StorageResponse, error) {
	res, err := tandoor.Request[StorageResponse](ctx, c, "POST", "/api/storage/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}
	return res, nil
}

// List calls GET /api/storage/
func List(ctx context.Context, c *tandoor.Client, params ListParams) (*StorageListResponse, error) {
	qb := tandoor.NewQuery().
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[StorageListResponse](ctx, c, "GET", "/api/storage/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list storages: %w", err)
	}
	return res, nil
}
