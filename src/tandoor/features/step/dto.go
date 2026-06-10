package step

type StepParam struct {
	Name        string `json:"name,omitempty"`
	Instruction string `json:"instruction,omitempty"`
	Time        *int   `json:"time,omitempty"`
	RecipeID    int    `json:"step_recipe,omitempty"`
	Order       *int   `json:"order,omitempty"`
	Ingredients []int  `json:"ingredients"`
}

type StepResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Instruction string `json:"instruction"`
	Time        int    `json:"time"`
	Recipe      int    `json:"step_recipe"`
	Order       int    `json:"order"`
}
