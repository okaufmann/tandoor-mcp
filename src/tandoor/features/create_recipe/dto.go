package create_recipe

type CreateRecipeParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Servings    *int        `json:"servings,omitempty"`
	WorkingTime *int        `json:"working_time,omitempty"`
	WaitingTime *int        `json:"waiting_time,omitempty"`
	Steps       []StepParam `json:"steps"`
}

type StepParam struct {
	Instruction string `json:"instruction,omitempty"`
}

type RecipeResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
