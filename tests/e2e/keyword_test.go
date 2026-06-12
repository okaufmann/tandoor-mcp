package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/keyword"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_keyword"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_keywords"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestKeywordE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateKeyword", func(t *testing.T) {
		// Arrange
		name := "Italian"
		description := "Italian cuisine and recipes"

		args := create_keyword.Args{
			Name:        name,
			Description: description,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_keyword", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		kw := infra.ParseToolResponse[keyword.KeywordResponse](t, res)

		if kw.ID == 0 {
			t.Errorf("expected keyword ID > 0, got 0")
		}
		if kw.Name != name {
			t.Errorf("expected name %q, got %q", name, kw.Name)
		}
		if kw.Description != description {
			t.Errorf("expected description %q, got %q", description, kw.Description)
		}
	})

	t.Run("HappyPath_GetKeywords", func(t *testing.T) {
		// Arrange: Create a keyword first
		name := "Mexican"
		args := create_keyword.Args{
			Name: name,
		}
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_keyword", args)
		infra.AssertToolSuccess(t, createRes, createErr)
		kw := infra.ParseToolResponse[keyword.KeywordResponse](t, createRes)

		// Act
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_keywords", get_keywords.Args{
			Query: &name,
		})

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		list := infra.ParseToolResponse[keyword.KeywordListResponse](t, listRes)

		found := false
		for _, item := range list.Results {
			if item.ID == kw.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find keyword ID %d in get_keywords response", kw.ID)
		}
	})

	t.Run("ValidationError_MissingName", func(t *testing.T) {
		// Arrange
		args := create_keyword.Args{
			Name: "", // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_keyword", args)

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
