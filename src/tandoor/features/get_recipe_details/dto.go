package get_recipe_details

type RecipeResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Steps       []Step `json:"steps"`
	Created     string `json:"created"`
}

type Step struct {
	ID          int    `json:"id"`
	Instruction string `json:"instruction"`
}
