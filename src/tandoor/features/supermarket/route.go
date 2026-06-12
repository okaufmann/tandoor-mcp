package supermarket

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// CreateSupermarket calls POST /api/supermarket/
func CreateSupermarket(ctx context.Context, c *tandoor.Client, params SupermarketParam) (*SupermarketResponse, error) {
	res, err := tandoor.Request[SupermarketResponse](ctx, c, "POST", "/api/supermarket/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create supermarket: %w", err)
	}
	return res, nil
}

// ListSupermarkets calls GET /api/supermarket/
func ListSupermarkets(ctx context.Context, c *tandoor.Client, params ListParams) (*SupermarketListResponse, error) {
	qb := tandoor.NewQuery().
		Add("query", params.Query).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[SupermarketListResponse](ctx, c, "GET", "/api/supermarket/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list supermarkets: %w", err)
	}
	return res, nil
}

// CreateCategory calls POST /api/supermarket-category/
func CreateCategory(ctx context.Context, c *tandoor.Client, params SupermarketCategoryParam) (*SupermarketCategoryResponse, error) {
	res, err := tandoor.Request[SupermarketCategoryResponse](ctx, c, "POST", "/api/supermarket-category/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create supermarket category: %w", err)
	}
	return res, nil
}

// ListCategories calls GET /api/supermarket-category/
func ListCategories(ctx context.Context, c *tandoor.Client, params ListParams) (*SupermarketCategoryListResponse, error) {
	qb := tandoor.NewQuery().
		Add("query", params.Query).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[SupermarketCategoryListResponse](ctx, c, "GET", "/api/supermarket-category/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list supermarket categories: %w", err)
	}
	return res, nil
}

// GetCategory calls GET /api/supermarket-category/{id}/
func GetCategory(ctx context.Context, c *tandoor.Client, id int) (*SupermarketCategoryResponse, error) {
	path := fmt.Sprintf("/api/supermarket-category/%d/", id)
	res, err := tandoor.Request[SupermarketCategoryResponse](ctx, c, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get supermarket category %d: %w", id, err)
	}
	return res, nil
}

// CreateRelation calls POST /api/supermarket-category-relation/
func CreateRelation(ctx context.Context, c *tandoor.Client, params SupermarketCategoryRelationParam) (*SupermarketCategoryRelationResponse, error) {
	res, err := tandoor.Request[SupermarketCategoryRelationResponse](ctx, c, "POST", "/api/supermarket-category-relation/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create supermarket category relation: %w", err)
	}
	return res, nil
}

// ListRelations calls GET /api/supermarket-category-relation/
func ListRelations(ctx context.Context, c *tandoor.Client, params ListParams) (*SupermarketCategoryRelationListResponse, error) {
	qb := tandoor.NewQuery().
		Add("query", params.Query).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[SupermarketCategoryRelationListResponse](ctx, c, "GET", "/api/supermarket-category-relation/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list supermarket category relations: %w", err)
	}
	return res, nil
}

// DeleteRelation calls DELETE /api/supermarket-category-relation/{id}/
func DeleteRelation(ctx context.Context, c *tandoor.Client, id int) error {
	path := fmt.Sprintf("/api/supermarket-category-relation/%d/", id)
	_, err := tandoor.Request[any](ctx, c, "DELETE", path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete supermarket category relation %d: %w", id, err)
	}
	return nil
}
