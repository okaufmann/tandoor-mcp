package inventory

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// ListHouseholds calls GET /api/household/
func ListHouseholds(ctx context.Context, c *tandoor.Client) (*HouseholdListResponse, error) {
	res, err := tandoor.Request[HouseholdListResponse](ctx, c, "GET", "/api/household/", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list households: %w", err)
	}
	return res, nil
}

// CreateHousehold calls POST /api/household/
func CreateHousehold(ctx context.Context, c *tandoor.Client, name string) (*HouseholdResponse, error) {
	params := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	res, err := tandoor.Request[HouseholdResponse](ctx, c, "POST", "/api/household/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create household: %w", err)
	}
	return res, nil
}

// GetDefaultHousehold retrieves the first available household or creates one if none exists
func GetDefaultHousehold(ctx context.Context, c *tandoor.Client) (*HouseholdResponse, error) {
	res, err := ListHouseholds(ctx, c)
	if err != nil {
		return nil, err
	}
	if len(res.Results) > 0 {
		return &res.Results[0], nil
	}
	return CreateHousehold(ctx, c, "Default Household")
}

// CreateLocation calls POST /api/inventory-location/
func CreateLocation(ctx context.Context, c *tandoor.Client, params InventoryLocationParam) (*InventoryLocationResponse, error) {
	res, err := tandoor.Request[InventoryLocationResponse](ctx, c, "POST", "/api/inventory-location/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory location: %w", err)
	}
	return res, nil
}

// ListLocations calls GET /api/inventory-location/
func ListLocations(ctx context.Context, c *tandoor.Client, params ListLocationsParams) (*InventoryLocationListResponse, error) {
	qb := tandoor.NewQuery().
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[InventoryLocationListResponse](ctx, c, "GET", "/api/inventory-location/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory locations: %w", err)
	}
	return res, nil
}

// GetLocation calls GET /api/inventory-location/{id}/
func GetLocation(ctx context.Context, c *tandoor.Client, id int) (*InventoryLocationResponse, error) {
	path := fmt.Sprintf("/api/inventory-location/%d/", id)
	res, err := tandoor.Request[InventoryLocationResponse](ctx, c, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory location %d: %w", id, err)
	}
	return res, nil
}

// CreateEntry calls POST /api/inventory-entry/
func CreateEntry(ctx context.Context, c *tandoor.Client, params InventoryEntryParam) (*InventoryEntryResponse, error) {
	res, err := tandoor.Request[InventoryEntryResponse](ctx, c, "POST", "/api/inventory-entry/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory entry: %w", err)
	}
	return res, nil
}

// ListEntries calls GET /api/inventory-entry/
func ListEntries(ctx context.Context, c *tandoor.Client, params ListEntriesParams) (*InventoryEntryListResponse, error) {
	qb := tandoor.NewQuery().
		Add("inventory_location_id", params.InventoryLocationID).
		Add("food_id", params.FoodID).
		Add("empty", params.Empty).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[InventoryEntryListResponse](ctx, c, "GET", "/api/inventory-entry/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory entries: %w", err)
	}
	return res, nil
}

// GetEntry calls GET /api/inventory-entry/{id}/
func GetEntry(ctx context.Context, c *tandoor.Client, id int) (*InventoryEntryResponse, error) {
	path := fmt.Sprintf("/api/inventory-entry/%d/", id)
	res, err := tandoor.Request[InventoryEntryResponse](ctx, c, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory entry %d: %w", id, err)
	}
	return res, nil
}

// PatchEntry calls PATCH /api/inventory-entry/{id}/
func PatchEntry(ctx context.Context, c *tandoor.Client, id int, params InventoryEntryUpdateParam) (*InventoryEntryResponse, error) {
	path := fmt.Sprintf("/api/inventory-entry/%d/", id)
	res, err := tandoor.Request[InventoryEntryResponse](ctx, c, "PATCH", path, nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory entry %d: %w", id, err)
	}
	return res, nil
}

// DeleteEntry calls DELETE /api/inventory-entry/{id}/
func DeleteEntry(ctx context.Context, c *tandoor.Client, id int) error {
	path := fmt.Sprintf("/api/inventory-entry/%d/", id)
	_, err := tandoor.Request[any](ctx, c, "DELETE", path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete inventory entry %d: %w", id, err)
	}
	return nil
}

// ListLogs calls GET /api/inventory-log/
func ListLogs(ctx context.Context, c *tandoor.Client, params ListLogsParams) (*InventoryLogListResponse, error) {
	qb := tandoor.NewQuery().
		Add("entry_id", params.EntryID).
		Add("food_id", params.FoodID).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[InventoryLogListResponse](ctx, c, "GET", "/api/inventory-log/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory logs: %w", err)
	}
	return res, nil
}
