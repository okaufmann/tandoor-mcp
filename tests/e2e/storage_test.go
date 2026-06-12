package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/storage"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_storage"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_storages"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestStorageE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateStorage", func(t *testing.T) {
		// Arrange
		name := "My Local Storage"
		method := "LOCAL"
		path := "/var/tandoor/media"

		args := create_storage.Args{
			Name:   name,
			Method: method,
			Path:   path,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_storage", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		store := infra.ParseToolResponse[storage.StorageResponse](t, res)

		if store.ID == 0 {
			t.Errorf("expected storage ID > 0, got 0")
		}
		if store.Name != name {
			t.Errorf("expected name %q, got %q", name, store.Name)
		}
		if store.Method != method {
			t.Errorf("expected method %q, got %q", method, store.Method)
		}
		if store.Path != path {
			t.Errorf("expected path %q, got %q", path, store.Path)
		}
	})

	t.Run("HappyPath_GetStorages", func(t *testing.T) {
		// Arrange: Create a storage integration first
		name := "Nextcloud Storage"
		method := "NEXTCLOUD"
		urlVal := "https://nextcloud.example.com"
		args := create_storage.Args{
			Name:   name,
			Method: method,
			URL:    &urlVal,
		}
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_storage", args)
		infra.AssertToolSuccess(t, createRes, createErr)
		store := infra.ParseToolResponse[storage.StorageResponse](t, createRes)

		// Act
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_storages", get_storages.Args{})

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		list := infra.ParseToolResponse[storage.StorageListResponse](t, listRes)

		found := false
		for _, item := range list.Results {
			if item.ID == store.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find storage ID %d in get_storages response", store.ID)
		}
	})

	t.Run("ValidationError_MissingName", func(t *testing.T) {
		// Arrange
		args := create_storage.Args{
			Name:   "", // invalid
			Method: "LOCAL",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_storage", args)

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
