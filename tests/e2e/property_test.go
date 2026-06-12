package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/property"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_property"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_property_type"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_properties"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_property_types"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestPropertyE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreatePropertyType", func(t *testing.T) {
		// Arrange
		name := "Vitamin C"
		unit := "mg"
		desc := "Ascorbic acid"

		args := create_property_type.Args{
			Name:        name,
			Unit:        &unit,
			Description: &desc,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_property_type", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		pt := infra.ParseToolResponse[property.PropertyTypeResponse](t, res)

		if pt.ID == 0 {
			t.Errorf("expected property type ID > 0, got 0")
		}
		if pt.Name != name {
			t.Errorf("expected name %q, got %q", name, pt.Name)
		}
		if pt.Unit == nil || *pt.Unit != unit {
			t.Errorf("expected unit %q, got %v", unit, pt.Unit)
		}
		if pt.Description == nil || *pt.Description != desc {
			t.Errorf("expected description %q, got %v", desc, pt.Description)
		}
	})

	t.Run("HappyPath_GetPropertyTypes", func(t *testing.T) {
		// Arrange: Create a property type first
		name := "Calcium"
		args := create_property_type.Args{
			Name: name,
		}
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_property_type", args)
		infra.AssertToolSuccess(t, createRes, createErr)
		pt := infra.ParseToolResponse[property.PropertyTypeResponse](t, createRes)

		// Act
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_property_types", get_property_types.Args{})

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		list := infra.ParseToolResponse[property.PropertyTypeListResponse](t, listRes)

		found := false
		for _, item := range list.Results {
			if item.ID == pt.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find property type ID %d in get_property_types response", pt.ID)
		}
	})

	t.Run("HappyPath_CreateProperty", func(t *testing.T) {
		// Arrange: Create a property type
		ptRes, err1 := infra.CallTool(ctx, fixture.Client, "create_property_type", create_property_type.Args{Name: "Iron"})
		infra.AssertToolSuccess(t, ptRes, err1)
		pt := infra.ParseToolResponse[property.PropertyTypeResponse](t, ptRes)

		args := create_property.Args{
			PropertyAmount: 14.5,
			PropertyTypeID: pt.ID,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_property", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		prop := infra.ParseToolResponse[property.PropertyResponse](t, res)

		if prop.ID == 0 {
			t.Errorf("expected property ID > 0, got 0")
		}
		if prop.PropertyAmount == nil || *prop.PropertyAmount != 14.5 {
			t.Errorf("expected amount 14.5, got %v", prop.PropertyAmount)
		}
		if prop.PropertyType.ID != pt.ID {
			t.Errorf("expected property type ID %d, got %d", pt.ID, prop.PropertyType.ID)
		}
	})

	t.Run("HappyPath_GetProperties", func(t *testing.T) {
		// Arrange: Create a property type and property first
		ptRes, err1 := infra.CallTool(ctx, fixture.Client, "create_property_type", create_property_type.Args{Name: "Zinc"})
		infra.AssertToolSuccess(t, ptRes, err1)
		pt := infra.ParseToolResponse[property.PropertyTypeResponse](t, ptRes)

		propRes, err2 := infra.CallTool(ctx, fixture.Client, "create_property", create_property.Args{
			PropertyAmount: 8.0,
			PropertyTypeID: pt.ID,
		})
		infra.AssertToolSuccess(t, propRes, err2)
		prop := infra.ParseToolResponse[property.PropertyResponse](t, propRes)

		// Act
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_properties", get_properties.Args{})

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		list := infra.ParseToolResponse[property.PropertyListResponse](t, listRes)

		found := false
		for _, item := range list.Results {
			if item.ID == prop.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find property ID %d in get_properties response", prop.ID)
		}
	})

	t.Run("ValidationError_MissingName", func(t *testing.T) {
		// Arrange
		args := create_property_type.Args{
			Name: "", // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_property_type", args)

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

	t.Run("ValidationError_ZeroAmount", func(t *testing.T) {
		// Arrange
		args := create_property.Args{
			PropertyAmount: 0, // invalid
			PropertyTypeID: 1,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_property", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true for zero amount")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "property_amount must be greater than 0") {
			t.Errorf("expected error message to contain 'property_amount must be greater than 0', got %q", errText)
		}
	})
}
