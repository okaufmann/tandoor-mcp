package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/mealtype"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_meal_type"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_meal_types"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestMealTypeE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateAndGet", func(t *testing.T) {
		// Arrange
		name := "Breakfast"

		createArgs := create_meal_type.Args{
			Name: name,
		}

		// Act
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_meal_type", createArgs)

		// Assert
		infra.AssertToolSuccess(t, createRes, createErr)
		logEntry := infra.ParseToolResponse[mealtype.MealTypeResponse](t, createRes)

		if logEntry.ID == 0 {
			t.Errorf("expected meal type ID > 0, got 0")
		}
		if logEntry.Name != name {
			t.Errorf("expected Name %q, got %q", name, logEntry.Name)
		}

		// Act
		listArgs := get_meal_types.Args{}
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_meal_types", listArgs)

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		logsList := infra.ParseToolResponse[mealtype.MealTypeListResponse](t, listRes)

		if logsList.Count == 0 {
			t.Errorf("expected count > 0")
		}
		found := false
		for _, entry := range logsList.Results {
			if entry.ID == logEntry.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find meal type ID %d in get_meal_types response", logEntry.ID)
		}
	})

	t.Run("ValidationError_MissingName", func(t *testing.T) {
		// Arrange
		args := create_meal_type.Args{
			Name: "", // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_meal_type", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "name is required") {
			t.Errorf("expected validation message, got %q", errText)
		}
	})
}
