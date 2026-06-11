package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/bookmarklet_import"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_bookmarklet_import"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_bookmarklet_imports"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestBookmarkletImportE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateAndGet", func(t *testing.T) {
		// Arrange
		urlVal := "https://example.com/recipe"
		htmlVal := "<html><body>Recipe HTML</body></html>"

		createArgs := create_bookmarklet_import.Args{
			Url:  &urlVal,
			Html: htmlVal,
		}

		// Act
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_bookmarklet_import", createArgs)

		// Assert
		infra.AssertToolSuccess(t, createRes, createErr)
		logEntry := infra.ParseToolResponse[bookmarklet_import.BookmarkletImportResponse](t, createRes)

		if logEntry.ID == 0 {
			t.Errorf("expected bookmarklet import ID > 0, got 0")
		}
		if logEntry.Url == nil || *logEntry.Url != urlVal {
			t.Errorf("expected Url %q, got %v", urlVal, logEntry.Url)
		}
		if logEntry.Html != htmlVal {
			t.Errorf("expected Html %q, got %q", htmlVal, logEntry.Html)
		}

		// Act
		listArgs := get_bookmarklet_imports.Args{}
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_bookmarklet_imports", listArgs)

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		logsList := infra.ParseToolResponse[bookmarklet_import.BookmarkletImportListResponse](t, listRes)

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
			t.Errorf("expected to find bookmarklet import ID %d in get_bookmarklet_imports response", logEntry.ID)
		}
	})

	t.Run("ValidationError_MissingHtml", func(t *testing.T) {
		// Arrange
		args := create_bookmarklet_import.Args{
			Html: "",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_bookmarklet_import", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "html content is required") {
			t.Errorf("expected validation message, got %q", errText)
		}
	})
}
