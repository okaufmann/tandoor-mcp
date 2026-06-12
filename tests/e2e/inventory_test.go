package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/inventory"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_inventory_location"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_inventory_entry"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_inventory_entries"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_inventory_locations"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_inventory_logs"
	"github.com/compilercomplied/tandoor-mcp/src/tools/update_inventory_entry"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestInventoryE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateLocation", func(t *testing.T) {
		// Arrange
		name := "Kitchen Pantry"
		isFreezer := false
		args := create_inventory_location.Args{
			Name:      name,
			IsFreezer: isFreezer,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_inventory_location", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		loc := infra.ParseToolResponse[inventory.InventoryLocationResponse](t, res)

		if loc.ID == 0 {
			t.Errorf("expected location ID > 0, got 0")
		}
		if loc.Name != name {
			t.Errorf("expected name %q, got %q", name, loc.Name)
		}
		if loc.IsFreezer != isFreezer {
			t.Errorf("expected is_freezer %v, got %v", isFreezer, loc.IsFreezer)
		}
	})

	t.Run("HappyPath_GetLocations", func(t *testing.T) {
		// Arrange: Create a location first
		args := create_inventory_location.Args{
			Name:      "Garage Freezer",
			IsFreezer: true,
		}
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_inventory_location", args)
		infra.AssertToolSuccess(t, createRes, createErr)
		loc := infra.ParseToolResponse[inventory.InventoryLocationResponse](t, createRes)

		// Act
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_inventory_locations", get_inventory_locations.Args{})

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		list := infra.ParseToolResponse[inventory.InventoryLocationListResponse](t, listRes)

		found := false
		for _, item := range list.Results {
			if item.ID == loc.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find location ID %d in get_inventory_locations response", loc.ID)
		}
	})

	t.Run("HappyPath_CreateEntry", func(t *testing.T) {
		// Arrange: Create location first
		locArgs := create_inventory_location.Args{
			Name: "Cabinets",
		}
		locRes, locErr := infra.CallTool(ctx, fixture.Client, "create_inventory_location", locArgs)
		infra.AssertToolSuccess(t, locRes, locErr)
		loc := infra.ParseToolResponse[inventory.InventoryLocationResponse](t, locRes)

		foodName := "Flour"
		unitName := "grams"
		amount := "1000"
		note := "Baking flour"

		entryArgs := create_inventory_entry.Args{
			InventoryLocationID: loc.ID,
			FoodNameOrID:        foodName,
			Amount:              amount,
			UnitNameOrID:        unitName,
			Note:                &note,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_inventory_entry", entryArgs)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		entry := infra.ParseToolResponse[inventory.InventoryEntryResponse](t, res)

		if entry.ID == 0 {
			t.Errorf("expected entry ID > 0, got 0")
		}
		if entry.InventoryLocation.ID != loc.ID {
			t.Errorf("expected location ID %d, got %d", loc.ID, entry.InventoryLocation.ID)
		}
		if entry.Food.Name != foodName {
			t.Errorf("expected food name %q, got %q", foodName, entry.Food.Name)
		}
		if entry.Unit.Name != unitName {
			t.Errorf("expected unit name %q, got %q", unitName, entry.Unit.Name)
		}
		if entry.Amount != 1000.0 {
			t.Errorf("expected amount 1000.0, got %f", entry.Amount)
		}
	})

	t.Run("HappyPath_GetEntries", func(t *testing.T) {
		// Arrange: Create location and entry first
		locRes, locErr := infra.CallTool(ctx, fixture.Client, "create_inventory_location", create_inventory_location.Args{
			Name: "Shelves",
		})
		infra.AssertToolSuccess(t, locRes, locErr)
		loc := infra.ParseToolResponse[inventory.InventoryLocationResponse](t, locRes)

		entryRes, entryErr := infra.CallTool(ctx, fixture.Client, "create_inventory_entry", create_inventory_entry.Args{
			InventoryLocationID: loc.ID,
			FoodNameOrID:        "Sugar",
			Amount:              "500",
			UnitNameOrID:        "grams",
		})
		infra.AssertToolSuccess(t, entryRes, entryErr)
		entry := infra.ParseToolResponse[inventory.InventoryEntryResponse](t, entryRes)

		// Act
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_inventory_entries", get_inventory_entries.Args{
			InventoryLocationID: &loc.ID,
		})

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		list := infra.ParseToolResponse[inventory.InventoryEntryListResponse](t, listRes)

		found := false
		for _, item := range list.Results {
			if item.ID == entry.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find entry ID %d in get_inventory_entries response", entry.ID)
		}
	})

	t.Run("HappyPath_UpdateEntry", func(t *testing.T) {
		// Arrange: Create location and entry first
		locRes, locErr := infra.CallTool(ctx, fixture.Client, "create_inventory_location", create_inventory_location.Args{
			Name: "Basement",
		})
		infra.AssertToolSuccess(t, locRes, locErr)
		loc := infra.ParseToolResponse[inventory.InventoryLocationResponse](t, locRes)

		entryRes, entryErr := infra.CallTool(ctx, fixture.Client, "create_inventory_entry", create_inventory_entry.Args{
			InventoryLocationID: loc.ID,
			FoodNameOrID:        "Potatoes",
			Amount:              "10",
			UnitNameOrID:        "pieces",
		})
		infra.AssertToolSuccess(t, entryRes, entryErr)
		entry := infra.ParseToolResponse[inventory.InventoryEntryResponse](t, entryRes)

		newAmount := "15"
		updateArgs := update_inventory_entry.Args{
			EntryID: entry.ID,
			Amount:  &newAmount,
		}

		// Act
		updateRes, updateErr := infra.CallTool(ctx, fixture.Client, "update_inventory_entry", updateArgs)

		// Assert
		infra.AssertToolSuccess(t, updateRes, updateErr)
		updated := infra.ParseToolResponse[inventory.InventoryEntryResponse](t, updateRes)

		if updated.Amount != 15.0 {
			t.Errorf("expected amount 15.0, got %f", updated.Amount)
		}
	})

	t.Run("HappyPath_GetLogs", func(t *testing.T) {
		// Arrange: Create location and entry to generate transaction logs
		locRes, locErr := infra.CallTool(ctx, fixture.Client, "create_inventory_location", create_inventory_location.Args{
			Name: "Main Store",
		})
		infra.AssertToolSuccess(t, locRes, locErr)
		loc := infra.ParseToolResponse[inventory.InventoryLocationResponse](t, locRes)

		entryRes, entryErr := infra.CallTool(ctx, fixture.Client, "create_inventory_entry", create_inventory_entry.Args{
			InventoryLocationID: loc.ID,
			FoodNameOrID:        "Rice",
			Amount:              "5",
			UnitNameOrID:        "kg",
		})
		infra.AssertToolSuccess(t, entryRes, entryErr)
		entry := infra.ParseToolResponse[inventory.InventoryEntryResponse](t, entryRes)

		// Act
		logRes, logErr := infra.CallTool(ctx, fixture.Client, "get_inventory_logs", get_inventory_logs.Args{
			EntryID: &entry.ID,
		})

		// Assert
		infra.AssertToolSuccess(t, logRes, logErr)
		logs := infra.ParseToolResponse[inventory.InventoryLogListResponse](t, logRes)

		if logs.Count == 0 {
			t.Errorf("expected at least one inventory log for the created entry")
		}
	})

	t.Run("ValidationError_MissingLocationName", func(t *testing.T) {
		// Arrange
		args := create_inventory_location.Args{
			Name: "", // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_inventory_location", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true for empty name")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "name is required") {
			t.Errorf("expected error message to contain 'name is required', got %q", errText)
		}
	})
}
