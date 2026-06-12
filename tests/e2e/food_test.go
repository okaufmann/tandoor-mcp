package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/food"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_food"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_food_inherit_fields"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_foods"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestFoodE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateFood", func(t *testing.T) {
		// Arrange
		name := "Organic Blueberry"
		pluralName := "Organic Blueberries"
		description := "Fresh handpicked blueberries"
		ignoreShopping := true

		args := create_food.Args{
			Name:           name,
			PluralName:     &pluralName,
			Description:    description,
			IgnoreShopping: ignoreShopping,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_food", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		fd := infra.ParseToolResponse[food.FoodResponse](t, res)

		if fd.ID == 0 {
			t.Errorf("expected food ID > 0, got 0")
		}
		if fd.Name != name {
			t.Errorf("expected name %q, got %q", name, fd.Name)
		}
		if fd.PluralName == nil || *fd.PluralName != pluralName {
			t.Errorf("expected plural name %q, got %v", pluralName, fd.PluralName)
		}
		if fd.Description != description {
			t.Errorf("expected description %q, got %q", description, fd.Description)
		}
		if fd.IgnoreShopping != ignoreShopping {
			t.Errorf("expected ignore_shopping %v, got %v", ignoreShopping, fd.IgnoreShopping)
		}
	})

	t.Run("HappyPath_GetFoods", func(t *testing.T) {
		// Arrange: Create a food first
		name := "Fresh Strawberry"
		args := create_food.Args{
			Name: name,
		}
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_food", args)
		infra.AssertToolSuccess(t, createRes, createErr)
		fd := infra.ParseToolResponse[food.FoodResponse](t, createRes)

		// Act
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_foods", get_foods.Args{
			Query: &name,
		})

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		list := infra.ParseToolResponse[food.FoodListResponse](t, listRes)

		found := false
		for _, item := range list.Results {
			if item.ID == fd.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find food ID %d in get_foods response", fd.ID)
		}
	})

	t.Run("HappyPath_GetFoodInheritFields", func(t *testing.T) {
		// Arrange: None required (uses standard Tandoor default metadata fields)

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "get_food_inherit_fields", get_food_inherit_fields.Args{})

		// Assert
		infra.AssertToolSuccess(t, res, err)
		list := infra.ParseToolResponse[[]food.FoodInheritFieldResponse](t, res)

		// Tandoor provides standard inherit fields by default
		if len(*list) == 0 {
			t.Errorf("expected at least one food inheritance field, got 0")
		}
	})

	t.Run("ValidationError_MissingName", func(t *testing.T) {
		// Arrange
		args := create_food.Args{
			Name: "", // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_food", args)

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
