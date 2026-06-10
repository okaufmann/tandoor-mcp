package step

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

type RecipeSteps struct {
	Steps []StepParam `json:"steps"`
}

// Create executes a create step request by fetching the existing recipe steps and patching them
func Create(ctx context.Context, c *tandoor.Client, params StepParam) (*StepResponse, error) {
	// First fetch the existing recipe
	recipeEndpoint := fmt.Sprintf("/api/recipe/%d/", params.RecipeID)
	existingRecipe, err := tandoor.Request[RecipeSteps](ctx, c, "GET", recipeEndpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recipe: %w", err)
	}

	// Append the new step
	existingRecipe.Steps = append(existingRecipe.Steps, params)

	// Patch the recipe
	updatedRecipe, err := tandoor.Request[RecipeSteps](ctx, c, "PATCH", recipeEndpoint, nil, existingRecipe)
	if err != nil {
		return nil, fmt.Errorf("failed to update recipe with new step: %w", err)
	}

	// Return the newly created step (the last one in the array)
	if len(updatedRecipe.Steps) > 0 {
		lastStep := updatedRecipe.Steps[len(updatedRecipe.Steps)-1]
		return &StepResponse{
			Instruction: lastStep.Instruction,
			Name:        lastStep.Name,
			Order:       0, // Or whatever
			Recipe:      params.RecipeID,
		}, nil
	}

	return nil, fmt.Errorf("failed to extract created step")
}
