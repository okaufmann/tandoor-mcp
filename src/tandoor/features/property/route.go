package property

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// CreatePropertyType calls POST /api/property-type/
func CreatePropertyType(ctx context.Context, c *tandoor.Client, params PropertyTypeParam) (*PropertyTypeResponse, error) {
	res, err := tandoor.Request[PropertyTypeResponse](ctx, c, "POST", "/api/property-type/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create property type: %w", err)
	}
	return res, nil
}

// ListPropertyTypes calls GET /api/property-type/
func ListPropertyTypes(ctx context.Context, c *tandoor.Client, params ListPropertyTypesParams) (*PropertyTypeListResponse, error) {
	qb := tandoor.NewQuery().
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	for _, cat := range params.Category {
		qb.Add("category", &cat)
	}

	res, err := tandoor.Request[PropertyTypeListResponse](ctx, c, "GET", "/api/property-type/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list property types: %w", err)
	}
	return res, nil
}

// GetPropertyType calls GET /api/property-type/{id}/
func GetPropertyType(ctx context.Context, c *tandoor.Client, id int) (*PropertyTypeResponse, error) {
	path := fmt.Sprintf("/api/property-type/%d/", id)
	res, err := tandoor.Request[PropertyTypeResponse](ctx, c, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get property type %d: %w", id, err)
	}
	return res, nil
}

// CreateProperty calls POST /api/property/
func CreateProperty(ctx context.Context, c *tandoor.Client, params PropertyParam) (*PropertyResponse, error) {
	res, err := tandoor.Request[PropertyResponse](ctx, c, "POST", "/api/property/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create property: %w", err)
	}
	return res, nil
}

// ListProperties calls GET /api/property/
func ListProperties(ctx context.Context, c *tandoor.Client, params ListPropertiesParams) (*PropertyListResponse, error) {
	qb := tandoor.NewQuery().
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[PropertyListResponse](ctx, c, "GET", "/api/property/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list properties: %w", err)
	}
	return res, nil
}

// GetProperty calls GET /api/property/{id}/
func GetProperty(ctx context.Context, c *tandoor.Client, id int) (*PropertyResponse, error) {
	path := fmt.Sprintf("/api/property/%d/", id)
	res, err := tandoor.Request[PropertyResponse](ctx, c, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get property %d: %w", id, err)
	}
	return res, nil
}
