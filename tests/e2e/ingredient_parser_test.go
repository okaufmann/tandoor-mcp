package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient_parser"
	"github.com/compilercomplied/tandoor-mcp/src/tools/parse_ingredients"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestIngredientParserE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_Single", func(t *testing.T) {
		// Arrange
		args := parse_ingredients.Args{
			Ingredient: "2 cups flour",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "parse_ingredients", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		parsed := infra.ParseToolResponse[ingredient_parser.IngredientParserResponse](t, res)

		if parsed.Ingredient == nil {
			t.Fatalf("expected parsed ingredient to not be nil")
		}
		if parsed.Ingredient.Amount != 2 {
			t.Errorf("expected amount 2, got %v", parsed.Ingredient.Amount)
		}
		if parsed.Ingredient.Food == nil || parsed.Ingredient.Food.Name != "flour" {
			t.Errorf("expected food name 'flour', got %v", parsed.Ingredient.Food)
		}
		if parsed.Ingredient.Unit == nil || parsed.Ingredient.Unit.Name != "cups" {
			t.Errorf("expected unit name 'cups', got %v", parsed.Ingredient.Unit)
		}
	})

	t.Run("ValidationError_Empty", func(t *testing.T) {
		// Arrange
		args := parse_ingredients.Args{}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "parse_ingredients", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "either ingredient or ingredients must be provided") {
			t.Errorf("expected validation message, got %q", errText)
		}
	})
}
