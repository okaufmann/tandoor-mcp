package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/shoppinglist"
	"github.com/compilercomplied/tandoor-mcp/src/tools/add_shopping_list_item"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_shopping_list"
	"github.com/compilercomplied/tandoor-mcp/src/tools/remove_shopping_list_item"
	"github.com/compilercomplied/tandoor-mcp/src/tools/update_shopping_list_item"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestShoppingListE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_AddShoppingListItem", func(t *testing.T) {
		// Arrange
		foodName := "Organic Apples"
		unitName := "pieces"
		amount := "5"
		note := "Get the red ones"

		addArgs := add_shopping_list_item.Args{
			FoodNameOrID: foodName,
			Amount:       amount,
			UnitNameOrID: unitName,
			Note:         note,
		}

		// Act
		addRes, addErr := infra.CallTool(ctx, fixture.Client, "add_shopping_list_item", addArgs)

		// Assert
		infra.AssertToolSuccess(t, addRes, addErr)
		entry := infra.ParseToolResponse[shoppinglist.ShoppingListEntryResponse](t, addRes)

		if entry.ID == 0 {
			t.Errorf("expected entry ID > 0, got 0")
		}
		if entry.Food == nil || entry.Food.Name != foodName {
			t.Errorf("expected food name %q, got %v", foodName, entry.Food)
		}
		if entry.Unit == nil || entry.Unit.Name != unitName {
			t.Errorf("expected unit name %q, got %v", unitName, entry.Unit)
		}
		if entry.Amount != 5.0 {
			t.Errorf("expected amount 5.0, got %f", entry.Amount)
		}
		if entry.Checked {
			t.Errorf("expected checked to be false initially")
		}
	})

	t.Run("HappyPath_GetShoppingList_Unchecked", func(t *testing.T) {
		// Arrange: Add an item first
		addArgs := add_shopping_list_item.Args{
			FoodNameOrID: "Bananas",
			Amount:       "3",
			UnitNameOrID: "bunches",
		}
		addRes, addErr := infra.CallTool(ctx, fixture.Client, "add_shopping_list_item", addArgs)
		infra.AssertToolSuccess(t, addRes, addErr)
		entry := infra.ParseToolResponse[shoppinglist.ShoppingListEntryResponse](t, addRes)

		// Act
		getRes, getErr := infra.CallTool(ctx, fixture.Client, "get_shopping_list", get_shopping_list.Args{
			Checked: "false",
		})

		// Assert
		infra.AssertToolSuccess(t, getRes, getErr)
		list := infra.ParseToolResponse[[]shoppinglist.ShoppingListEntryResponse](t, getRes)

		found := false
		for _, item := range *list {
			if item.ID == entry.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find entry ID %d in get_shopping_list response", entry.ID)
		}
	})

	t.Run("HappyPath_UpdateShoppingListItem", func(t *testing.T) {
		// Arrange: Add an item first
		addArgs := add_shopping_list_item.Args{
			FoodNameOrID: "Milk",
			Amount:       "1",
			UnitNameOrID: "liter",
		}
		addRes, addErr := infra.CallTool(ctx, fixture.Client, "add_shopping_list_item", addArgs)
		infra.AssertToolSuccess(t, addRes, addErr)
		entry := infra.ParseToolResponse[shoppinglist.ShoppingListEntryResponse](t, addRes)

		newAmount := "2.5"
		checkedVal := true
		updateArgs := update_shopping_list_item.Args{
			ItemID:  entry.ID,
			Amount:  &newAmount,
			Checked: &checkedVal,
		}

		// Act
		updateRes, updateErr := infra.CallTool(ctx, fixture.Client, "update_shopping_list_item", updateArgs)

		// Assert
		infra.AssertToolSuccess(t, updateRes, updateErr)
		updatedEntry := infra.ParseToolResponse[shoppinglist.ShoppingListEntryResponse](t, updateRes)

		if updatedEntry.Amount != 2.5 {
			t.Errorf("expected amount 2.5, got %f", updatedEntry.Amount)
		}
		if !updatedEntry.Checked {
			t.Errorf("expected checked to be true after update")
		}
	})

	t.Run("HappyPath_GetShoppingList_Checked", func(t *testing.T) {
		// Arrange: Add an item first and mark it checked
		addArgs := add_shopping_list_item.Args{
			FoodNameOrID: "Eggs",
			Amount:       "12",
			UnitNameOrID: "pcs",
		}
		addRes, addErr := infra.CallTool(ctx, fixture.Client, "add_shopping_list_item", addArgs)
		infra.AssertToolSuccess(t, addRes, addErr)
		entry := infra.ParseToolResponse[shoppinglist.ShoppingListEntryResponse](t, addRes)

		checkedVal := true
		updateArgs := update_shopping_list_item.Args{
			ItemID:  entry.ID,
			Checked: &checkedVal,
		}
		updateRes, updateErr := infra.CallTool(ctx, fixture.Client, "update_shopping_list_item", updateArgs)
		infra.AssertToolSuccess(t, updateRes, updateErr)

		// Act
		getCheckedRes, getCheckedErr := infra.CallTool(ctx, fixture.Client, "get_shopping_list", get_shopping_list.Args{
			Checked: "true",
		})

		// Assert
		infra.AssertToolSuccess(t, getCheckedRes, getCheckedErr)
		checkedList := infra.ParseToolResponse[[]shoppinglist.ShoppingListEntryResponse](t, getCheckedRes)

		foundChecked := false
		for _, item := range *checkedList {
			if item.ID == entry.ID {
				foundChecked = true
				break
			}
		}
		if !foundChecked {
			t.Errorf("expected to find checked entry ID %d in get_shopping_list (checked=true)", entry.ID)
		}
	})

	t.Run("HappyPath_RemoveShoppingListItem", func(t *testing.T) {
		// Arrange: Add an item first
		addArgs := add_shopping_list_item.Args{
			FoodNameOrID: "Bread",
			Amount:       "1",
			UnitNameOrID: "loaf",
		}
		addRes, addErr := infra.CallTool(ctx, fixture.Client, "add_shopping_list_item", addArgs)
		infra.AssertToolSuccess(t, addRes, addErr)
		entry := infra.ParseToolResponse[shoppinglist.ShoppingListEntryResponse](t, addRes)

		removeArgs := remove_shopping_list_item.Args{
			ItemID: entry.ID,
		}

		// Act
		removeRes, removeErr := infra.CallTool(ctx, fixture.Client, "remove_shopping_list_item", removeArgs)

		// Assert
		infra.AssertToolSuccess(t, removeRes, removeErr)

		// Verify deletion
		getBothRes, getBothErr := infra.CallTool(ctx, fixture.Client, "get_shopping_list", get_shopping_list.Args{
			Checked: "both",
		})
		infra.AssertToolSuccess(t, getBothRes, getBothErr)
		bothList := infra.ParseToolResponse[[]shoppinglist.ShoppingListEntryResponse](t, getBothRes)

		foundDeleted := false
		for _, item := range *bothList {
			if item.ID == entry.ID {
				foundDeleted = true
				break
			}
		}
		if foundDeleted {
			t.Errorf("expected entry ID %d to be deleted and not returned", entry.ID)
		}
	})

	t.Run("ValidationError_InvalidAmountFormat", func(t *testing.T) {
		// Arrange
		addArgs := add_shopping_list_item.Args{
			FoodNameOrID: "Bananas",
			Amount:       "invalid-amount",
			UnitNameOrID: "bunches",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "add_shopping_list_item", addArgs)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true for invalid amount")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "invalid amount format") {
			t.Errorf("expected error message to contain 'invalid amount format', got %q", errText)
		}
	})
}
