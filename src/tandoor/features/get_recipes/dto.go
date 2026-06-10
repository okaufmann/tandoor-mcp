package get_recipes

// Recipe represents the Tandoor Recipe DTO specific to search feature
type Recipe struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PrepTime    *int   `json:"prep_time,omitempty"`
	CookTime    *int   `json:"cook_time,omitempty"`
	Servings    *int   `json:"servings,omitempty"`
	Created     string `json:"created"`
}

// RecipeListResponse represents the list of recipes returned by the Tandoor API
type RecipeListResponse struct {
	Count    int      `json:"count"`
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []Recipe `json:"results"`
}

// GetRecipesParams represents the query parameters for filtering recipes
type GetRecipesParams struct {
	Query    string `json:"query,omitempty"`
	Search   string `json:"search,omitempty"`
	Foods    []int  `json:"foods,omitempty"`
	Keywords []int  `json:"keywords,omitempty"`
	Limit    *int   `json:"limit,omitempty"`
	Rating   *int   `json:"rating,omitempty"`
}
