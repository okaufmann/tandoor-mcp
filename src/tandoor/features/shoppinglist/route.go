package shoppinglist

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// CreateList calls POST /api/shopping-list/
func CreateList(ctx context.Context, c *tandoor.Client, params ShoppingListParam) (*ShoppingListResponse, error) {
	res, err := tandoor.Request[ShoppingListResponse](ctx, c, "POST", "/api/shopping-list/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create shopping list: %w", err)
	}
	return res, nil
}

// ListLists calls GET /api/shopping-list/
func ListLists(ctx context.Context, c *tandoor.Client) (*ShoppingListListResponse, error) {
	res, err := tandoor.Request[ShoppingListListResponse](ctx, c, "GET", "/api/shopping-list/", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list shopping lists: %w", err)
	}
	return res, nil
}

// GetOrCreateDefaultList ensures there is at least one shopping list and returns it
func GetOrCreateDefaultList(ctx context.Context, c *tandoor.Client) (*ShoppingListResponse, error) {
	listsRes, err := ListLists(ctx, c)
	if err != nil {
		return nil, err
	}
	if len(listsRes.Results) > 0 {
		return &listsRes.Results[0], nil
	}
	createRes, err := CreateList(ctx, c, ShoppingListParam{
		Name: "Default",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create default shopping list: %w", err)
	}
	return createRes, nil
}

// CreateEntry calls POST /api/shopping-list-entry/
func CreateEntry(ctx context.Context, c *tandoor.Client, params CreateEntryParam) (*ShoppingListEntryResponse, error) {
	res, err := tandoor.Request[ShoppingListEntryResponse](ctx, c, "POST", "/api/shopping-list-entry/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create shopping list entry: %w", err)
	}
	return res, nil
}

// ListEntries calls GET /api/shopping-list-entry/
func ListEntries(ctx context.Context, c *tandoor.Client, params ListEntriesParams) (*ShoppingListEntryListResponse, error) {
	qb := tandoor.NewQuery()
	if params.Mealplan != nil {
		qb.Add("mealplan", *params.Mealplan)
	}
	if params.UpdatedAfter != nil {
		qb.Add("updated_after", *params.UpdatedAfter)
	}
	if params.Page != nil {
		qb.Add("page", *params.Page)
	}
	if params.PageSize != nil {
		qb.Add("page_size", *params.PageSize)
	}

	res, err := tandoor.Request[ShoppingListEntryListResponse](ctx, c, "GET", "/api/shopping-list-entry/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list shopping list entries: %w", err)
	}
	return res, nil
}

// PatchEntry calls PATCH /api/shopping-list-entry/{id}/
func PatchEntry(ctx context.Context, c *tandoor.Client, id int, params UpdateEntryParam) (*ShoppingListEntryResponse, error) {
	path := fmt.Sprintf("/api/shopping-list-entry/%d/", id)
	res, err := tandoor.Request[ShoppingListEntryResponse](ctx, c, "PATCH", path, nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to patch shopping list entry %d: %w", id, err)
	}
	return res, nil
}

// DeleteEntry calls DELETE /api/shopping-list-entry/{id}/
func DeleteEntry(ctx context.Context, c *tandoor.Client, id int) error {
	path := fmt.Sprintf("/api/shopping-list-entry/%d/", id)
	_, err := tandoor.Request[any](ctx, c, "DELETE", path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete shopping list entry %d: %w", id, err)
	}
	return nil
}
