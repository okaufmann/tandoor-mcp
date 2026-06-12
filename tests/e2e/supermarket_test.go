package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/supermarket"
	"github.com/compilercomplied/tandoor-mcp/src/tools/add_category_to_supermarket"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_supermarket"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_supermarket_category"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_supermarket_categories"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_supermarkets"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestSupermarketE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateSupermarketCategory", func(t *testing.T) {
		// Arrange
		catName := "Produce"
		catDesc := "Fresh fruits and vegetables"
		catArgs := create_supermarket_category.Args{
			Name:        catName,
			Description: catDesc,
		}

		// Act
		catRes, catErr := infra.CallTool(ctx, fixture.Client, "create_supermarket_category", catArgs)

		// Assert
		infra.AssertToolSuccess(t, catRes, catErr)
		category := infra.ParseToolResponse[supermarket.SupermarketCategoryResponse](t, catRes)

		if category.ID == 0 {
			t.Errorf("expected category ID > 0, got 0")
		}
		if category.Name != catName {
			t.Errorf("expected category name %q, got %q", catName, category.Name)
		}
		if category.Description != catDesc {
			t.Errorf("expected category description %q, got %q", catDesc, category.Description)
		}
	})

	t.Run("HappyPath_GetSupermarketCategories", func(t *testing.T) {
		// Arrange: Create a category first
		catArgs := create_supermarket_category.Args{
			Name:        "Bakery",
			Description: "Fresh bread and pastries",
		}
		catRes, catErr := infra.CallTool(ctx, fixture.Client, "create_supermarket_category", catArgs)
		infra.AssertToolSuccess(t, catRes, catErr)
		category := infra.ParseToolResponse[supermarket.SupermarketCategoryResponse](t, catRes)

		// Act
		listCatRes, listCatErr := infra.CallTool(ctx, fixture.Client, "get_supermarket_categories", get_supermarket_categories.Args{})

		// Assert
		infra.AssertToolSuccess(t, listCatRes, listCatErr)
		catList := infra.ParseToolResponse[supermarket.SupermarketCategoryListResponse](t, listCatRes)

		foundCat := false
		for _, cat := range catList.Results {
			if cat.ID == category.ID {
				foundCat = true
				break
			}
		}
		if !foundCat {
			t.Errorf("expected to find category ID %d in get_supermarket_categories list", category.ID)
		}
	})

	t.Run("HappyPath_CreateSupermarket", func(t *testing.T) {
		// Arrange
		smName := "Green Grocery"
		smDesc := "Local organic grocery store"
		smArgs := create_supermarket.Args{
			Name:        smName,
			Description: smDesc,
		}

		// Act
		smRes, smErr := infra.CallTool(ctx, fixture.Client, "create_supermarket", smArgs)

		// Assert
		infra.AssertToolSuccess(t, smRes, smErr)
		market := infra.ParseToolResponse[supermarket.SupermarketResponse](t, smRes)

		if market.ID == 0 {
			t.Errorf("expected supermarket ID > 0, got 0")
		}
		if market.Name != smName {
			t.Errorf("expected supermarket name %q, got %q", smName, market.Name)
		}
		if market.Description != smDesc {
			t.Errorf("expected supermarket description %q, got %q", smDesc, market.Description)
		}
	})

	t.Run("HappyPath_GetSupermarkets", func(t *testing.T) {
		// Arrange: Create a supermarket first
		smArgs := create_supermarket.Args{
			Name:        "SuperFood",
			Description: "Large supermarket chain",
		}
		smRes, smErr := infra.CallTool(ctx, fixture.Client, "create_supermarket", smArgs)
		infra.AssertToolSuccess(t, smRes, smErr)
		market := infra.ParseToolResponse[supermarket.SupermarketResponse](t, smRes)

		// Act
		listSmRes, listSmErr := infra.CallTool(ctx, fixture.Client, "get_supermarkets", get_supermarkets.Args{})

		// Assert
		infra.AssertToolSuccess(t, listSmRes, listSmErr)
		smList := infra.ParseToolResponse[supermarket.SupermarketListResponse](t, listSmRes)

		foundSm := false
		for _, sm := range smList.Results {
			if sm.ID == market.ID {
				foundSm = true
				break
			}
		}
		if !foundSm {
			t.Errorf("expected to find supermarket ID %d in get_supermarkets list", market.ID)
		}
	})

	t.Run("HappyPath_AddCategoryToSupermarket", func(t *testing.T) {
		// Arrange: Create a supermarket and a category first
		catArgs := create_supermarket_category.Args{
			Name: "Dairy",
		}
		catRes, catErr := infra.CallTool(ctx, fixture.Client, "create_supermarket_category", catArgs)
		infra.AssertToolSuccess(t, catRes, catErr)
		category := infra.ParseToolResponse[supermarket.SupermarketCategoryResponse](t, catRes)

		smArgs := create_supermarket.Args{
			Name: "Dairy Plaza",
		}
		smRes, smErr := infra.CallTool(ctx, fixture.Client, "create_supermarket", smArgs)
		infra.AssertToolSuccess(t, smRes, smErr)
		market := infra.ParseToolResponse[supermarket.SupermarketResponse](t, smRes)

		orderVal := 2
		relationArgs := add_category_to_supermarket.Args{
			SupermarketID: market.ID,
			CategoryID:    category.ID,
			Order:         orderVal,
		}

		// Act
		relRes, relErr := infra.CallTool(ctx, fixture.Client, "add_category_to_supermarket", relationArgs)

		// Assert
		infra.AssertToolSuccess(t, relRes, relErr)
		relation := infra.ParseToolResponse[supermarket.SupermarketCategoryRelationResponse](t, relRes)

		if relation.ID == 0 {
			t.Errorf("expected relation ID > 0, got 0")
		}
		if relation.Supermarket != market.ID {
			t.Errorf("expected supermarket ID %d, got %d", market.ID, relation.Supermarket)
		}
		if relation.Category.ID != category.ID {
			t.Errorf("expected category ID %d, got %d", category.ID, relation.Category.ID)
		}
		if relation.Order != orderVal {
			t.Errorf("expected order %d, got %d", orderVal, relation.Order)
		}
	})

	t.Run("ValidationError_MissingSupermarketName", func(t *testing.T) {
		// Arrange
		smArgs := create_supermarket.Args{
			Name: "", // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_supermarket", smArgs)

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
