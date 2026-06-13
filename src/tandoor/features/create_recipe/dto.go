package create_recipe

import (
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
)

type CreateRecipeParams struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Servings    *int        `json:"servings,omitempty"`
	WorkingTime *int        `json:"working_time,omitempty"`
	WaitingTime *int        `json:"waiting_time,omitempty"`
	Steps       []StepParam `json:"steps"`
}

type StepParam struct {
	Name        string                          `json:"name,omitempty"`
	Instruction string                          `json:"instruction,omitempty"`
	Time        *int                            `json:"time,omitempty"`
	Order       *int                            `json:"order,omitempty"`
	Ingredients []ingredient.IngredientResponse `json:"ingredients"`
}

type RecipeResponse struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Steps       []StepResponse `json:"steps,omitempty"`
}

type StepResponse struct {
	ID          int                             `json:"id"`
	Name        string                          `json:"name"`
	Instruction string                          `json:"instruction"`
	Time        int                             `json:"time"`
	Order       int                             `json:"order"`
	Ingredients []ingredient.IngredientResponse `json:"ingredients"`
}
