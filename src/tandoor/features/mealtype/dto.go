package mealtype

// MealTypeParam represents the JSON payload to create or update a meal type
type MealTypeParam struct {
	Name  string  `json:"name"`
	Order *int    `json:"order,omitempty"`
	Time  *string `json:"time,omitempty"`
	Color *string `json:"color,omitempty"`
}

// MealTypeResponse represents the MealType returned by Tandoor
type MealTypeResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Order     int     `json:"order"`
	Time      *string `json:"time"`
	Color     *string `json:"color"`
	CreatedBy int     `json:"created_by"`
}

// MealTypeListResponse represents a paginated list of MealTypes
type MealTypeListResponse struct {
	Count    int                `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []MealTypeResponse `json:"results"`
}

// ListParams represents query filters for meal types list
type ListParams struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}
