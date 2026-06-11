package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/step"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_step"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestCreateStepE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath", func(t *testing.T) {
		// Arrange
		args := create_tandoor_step.Args{
			RecipeID:    1,
			Name:        "Prep step",
			Instruction: "Chop the onions",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)

		s := infra.ParseToolResponse[step.StepResponse](t, res)

		if s.ID == 0 {
			t.Errorf("expected step ID > 0, got 0")
		}
		if s.Name != "Prep step" {
			t.Errorf("expected Name='Prep step', got %q", s.Name)
		}
		if s.Instruction != "Chop the onions" {
			t.Errorf("expected Instruction='Chop the onions', got %q", s.Instruction)
		}
	})

	t.Run("ValidationError_MissingInstructionAndName", func(t *testing.T) {
		// Arrange
		args := create_tandoor_step.Args{
			RecipeID: 1,
			// both Name and Instruction intentionally omitted
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}

		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "instruction or name is required") {
			t.Errorf("expected validation message about instruction or name, got %q", errText)
		}
	})

	t.Run("ValidationError_InvalidRecipeID", func(t *testing.T) {
		// Arrange
		args := create_tandoor_step.Args{
			RecipeID:    999999,
			Instruction: "Chop the onions",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true for non-existent recipe, got false")
		}

		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "Error creating step") {
			t.Errorf("expected error message to contain 'Error creating step', got %q", errText)
		}
	})
}
