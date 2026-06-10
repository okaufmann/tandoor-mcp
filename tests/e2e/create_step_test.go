package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/step"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_step"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestCreateStepE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath", func(t *testing.T) {
		res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", create_tandoor_step.Args{
			RecipeID:    1,
			Name:        "Prep step",
			Instruction: "Chop the onions",
		})
		AssertToolSuccess(t, res, err)
		
		t.Logf("HappyPath Response JSON: %s", res.Content[0].(*mcp_sdk.TextContent).Text)

		s := ParseToolResponse[step.StepResponse](t, res)
		if s.Recipe != 1 {
			t.Errorf("expected Recipe 1, got %v", s.Recipe)
		}
		if s.Instruction != "Chop the onions" {
			t.Errorf("expected instruction 'Chop the onions', got %q", s.Instruction)
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		// Missing recipe and instruction
		res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", create_tandoor_step.Args{})
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		
		if len(res.Content) > 0 {
			if tc, ok := res.Content[0].(*mcp_sdk.TextContent); ok {
				t.Logf("Validation Error Response JSON: %s", tc.Text)
			}
		}
		
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}

		var errMsgs []string
		for _, c := range res.Content {
			if tc, ok := c.(*mcp_sdk.TextContent); ok {
				errMsgs = append(errMsgs, tc.Text)
			}
		}
		
		errText := strings.Join(errMsgs, " ")
		if !strings.Contains(errText, "Error creating step") {
			t.Errorf("expected error message to contain 'Error creating step', got %q", errText)
		}
		// Tandoor will probably throw validation errors for Recipe and Instruction
		// We assert it propagates down correctly.
	})
}
