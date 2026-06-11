package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/recipe_import"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_recipe_import"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_recipe_imports"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestRecipeImportE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateAndGet", func(t *testing.T) {
		// Arrange
		storageName := "Local E2E Storage"
		storagePath := "/tmp/media"
		importName := "E2E Recipe Import File"

		createArgs := create_recipe_import.Args{
			StorageName: storageName,
			StoragePath: storagePath,
			ImportName:  importName,
		}

		// Act
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_recipe_import", createArgs)

		// Assert
		infra.AssertToolSuccess(t, createRes, createErr)
		logEntry := infra.ParseToolResponse[recipe_import.RecipeImportResponse](t, createRes)

		if logEntry.ID == 0 {
			t.Errorf("expected recipe import ID > 0, got 0")
		}
		if logEntry.Name != importName {
			t.Errorf("expected Name %q, got %q", importName, logEntry.Name)
		}
		if logEntry.Storage.Name != storageName {
			t.Errorf("expected Storage Name %q, got %q", storageName, logEntry.Storage.Name)
		}
		if logEntry.Storage.Path != storagePath {
			t.Errorf("expected Storage Path %q, got %q", storagePath, logEntry.Storage.Path)
		}

		// Act
		listArgs := get_recipe_imports.Args{}
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_recipe_imports", listArgs)

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		logsList := infra.ParseToolResponse[recipe_import.RecipeImportListResponse](t, listRes)

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
			t.Errorf("expected to find recipe import ID %d in get_recipe_imports response", logEntry.ID)
		}
	})

	t.Run("ValidationError_MissingImportName", func(t *testing.T) {
		// Arrange
		args := create_recipe_import.Args{
			StorageName: "Local E2E Storage",
			ImportName:  "", // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_recipe_import", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "import_name is required") {
			t.Errorf("expected validation message, got %q", errText)
		}
	})
}
