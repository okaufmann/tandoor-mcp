package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_ingredient"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestCreateIngredientE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			FoodName: "Onion",
			UnitName: "pieces",
			Amount:   2,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)

		ing := infra.ParseToolResponse[ingredient.IngredientResponse](t, res)

		if ing.ID == 0 {
			t.Errorf("expected ingredient ID > 0, got 0")
		}
		if ing.Food.Name != "Onion" {
			t.Errorf("expected Food.Name='Onion', got %q", ing.Food.Name)
		}
		if ing.Unit.Name != "pieces" {
			t.Errorf("expected Unit.Name='pieces', got %q", ing.Unit.Name)
		}
		if ing.Amount != 2 {
			t.Errorf("expected Amount=2, got %v", ing.Amount)
		}
	})

	t.Run("HappyPath_WithNote", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			FoodName: "Garlic",
			UnitName: "cloves",
			Amount:   3,
			Note:     "finely chopped",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)

		ing := infra.ParseToolResponse[ingredient.IngredientResponse](t, res)

		if ing.Note != "finely chopped" {
			t.Errorf("expected Note='finely chopped', got %q", ing.Note)
		}
	})

	t.Run("ValidationError_MissingFoodName", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			UnitName: "pieces",
			Amount:   1,
			// FoodName intentionally omitted
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "food_name is required") {
			t.Errorf("expected validation message about food_name, got %q", errText)
		}
	})

	t.Run("ValidationError_MissingUnitName", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			FoodName: "Salt",
			Amount:   1,
			// UnitName intentionally omitted
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "unit_name is required") {
			t.Errorf("expected validation message about unit_name, got %q", errText)
		}
	})

	t.Run("ValidationError_ZeroAmount", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			FoodName: "Pepper",
			UnitName: "grams",
			Amount:   0, // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true for zero amount, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "amount must be greater than 0") {
			t.Errorf("expected validation message about amount, got %q", errText)
		}
	})
}
