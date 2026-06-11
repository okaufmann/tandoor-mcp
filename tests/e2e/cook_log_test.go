package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/cooklog"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_cook_log"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_cook_logs"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestCookLogE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateAndGet", func(t *testing.T) {
		// Arrange
		recipeID := 1
		servings := 4
		rating := 5
		comment := "Delicious meal"

		createArgs := create_cook_log.Args{
			Recipe:   recipeID,
			Servings: &servings,
			Rating:   &rating,
			Comment:  &comment,
		}

		// Act
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_cook_log", createArgs)

		// Assert
		infra.AssertToolSuccess(t, createRes, createErr)
		logEntry := infra.ParseToolResponse[cooklog.CookLogResponse](t, createRes)

		if logEntry.ID == 0 {
			t.Errorf("expected cook log ID > 0, got 0")
		}
		if logEntry.Recipe != recipeID {
			t.Errorf("expected Recipe ID %d, got %d", recipeID, logEntry.Recipe)
		}
		if logEntry.Servings == nil || *logEntry.Servings != servings {
			t.Errorf("expected Servings %d, got %v", servings, logEntry.Servings)
		}
		if logEntry.Rating == nil || *logEntry.Rating != rating {
			t.Errorf("expected Rating %d, got %v", rating, logEntry.Rating)
		}
		if logEntry.Comment == nil || *logEntry.Comment != comment {
			t.Errorf("expected Comment %q, got %v", comment, logEntry.Comment)
		}

		// Act
		listArgs := get_cook_logs.Args{
			Recipe: &recipeID,
		}
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_cook_logs", listArgs)

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		logsList := infra.ParseToolResponse[cooklog.CookLogListResponse](t, listRes)

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
			t.Errorf("expected to find cook log ID %d in get_cook_logs response", logEntry.ID)
		}
	})

	t.Run("ValidationError_MissingRecipe", func(t *testing.T) {
		// Arrange
		args := create_cook_log.Args{
			Recipe: 0, // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_cook_log", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "recipe is required") {
			t.Errorf("expected validation message about recipe, got %q", errText)
		}
	})
}
