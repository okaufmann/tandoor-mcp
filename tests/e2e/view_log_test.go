package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/viewlog"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_view_log"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_view_logs"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestViewLogE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateAndGet", func(t *testing.T) {
		// Arrange
		recipeID := 1

		createArgs := create_view_log.Args{
			Recipe: recipeID,
		}

		// Act
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_view_log", createArgs)

		// Assert
		infra.AssertToolSuccess(t, createRes, createErr)
		logEntry := infra.ParseToolResponse[viewlog.ViewLogResponse](t, createRes)

		if logEntry.ID == 0 {
			t.Errorf("expected view log ID > 0, got 0")
		}
		if logEntry.Recipe != recipeID {
			t.Errorf("expected Recipe ID %d, got %d", recipeID, logEntry.Recipe)
		}
		if logEntry.CreatedBy == 0 {
			t.Errorf("expected CreatedBy ID > 0, got 0")
		}

		// Act
		listArgs := get_view_logs.Args{}
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_view_logs", listArgs)

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		logsList := infra.ParseToolResponse[viewlog.ViewLogListResponse](t, listRes)

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
			t.Errorf("expected to find view log ID %d in get_view_logs response", logEntry.ID)
		}
	})

	t.Run("ValidationError_MissingRecipe", func(t *testing.T) {
		// Arrange
		args := create_view_log.Args{
			Recipe: 0, // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_view_log", args)

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
