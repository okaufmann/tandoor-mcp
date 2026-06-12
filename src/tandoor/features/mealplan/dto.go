package mealplan

import "time"

// RecipeOverview represents the simple recipe metadata nested inside MealPlanResponse
type RecipeOverview struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// MealTypeResponse represents the meal type metadata nested inside MealPlanResponse
type MealTypeResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Order     int     `json:"order"`
	Time      *string `json:"time"`
	Color     *string `json:"color"`
	CreatedBy int     `json:"created_by"`
}

// MealPlanResponse represents a meal plan returned by Tandoor
type MealPlanResponse struct {
	ID           int              `json:"id"`
	Title        *string          `json:"title"`
	Recipe       *RecipeOverview  `json:"recipe"`
	Servings     float64          `json:"servings"`
	Note         *string          `json:"note"`
	NoteMarkdown string           `json:"note_markdown"`
	FromDate     time.Time        `json:"from_date"`
	ToDate       *time.Time       `json:"to_date"`
	MealType     MealTypeResponse `json:"meal_type"`
	CreatedBy    int              `json:"created_by"`
	RecipeName   string           `json:"recipe_name"`
	MealTypeName string           `json:"meal_type_name"`
	Shopping     bool             `json:"shopping"`
}

// MealPlanListResponse represents a paginated list of MealPlans
type MealPlanListResponse struct {
	Count    int                `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []MealPlanResponse `json:"results"`
}

// RecipeOverviewParam represents the nested recipe parameter
type RecipeOverviewParam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MealTypeParam represents the nested meal type parameter
type MealTypeParam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MealPlanParam represents the JSON payload to create or update a meal plan
type MealPlanParam struct {
	Title       *string              `json:"title,omitempty"`
	Recipe      *RecipeOverviewParam `json:"recipe,omitempty"`
	Servings    float64              `json:"servings"`
	Note        *string              `json:"note,omitempty"`
	FromDate    time.Time            `json:"from_date"`
	ToDate      *time.Time           `json:"to_date,omitempty"`
	MealType    MealTypeParam        `json:"meal_type"`
	AddShopping *bool                `json:"addshopping,omitempty"`
}

// ListParams represents query filters for fetching meal plans
type ListParams struct {
	FromDate *string `json:"from_date,omitempty"` // YYYY-MM-DD
	ToDate   *string `json:"to_date,omitempty"`   // YYYY-MM-DD
	MealType []int   `json:"meal_type,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"page_size,omitempty"`
}

// AutoMealPlanParam represents the JSON payload to auto-generate a meal plan
type AutoMealPlanParam struct {
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	MealTypeID  int       `json:"meal_type_id"`
	Keywords    []int     `json:"keywords,omitempty"`
	KeywordMode *string   `json:"keyword_mode,omitempty"` // "and" or "or"
	Servings    float64   `json:"servings"`
	AddShopping bool      `json:"addshopping"`
}

// AutoMealPlanResponse represents the response returned by the auto-plan endpoint
type AutoMealPlanResponse struct {
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	MealTypeID  int       `json:"meal_type_id"`
	Keywords    []int     `json:"keywords"`
	KeywordMode string    `json:"keyword_mode"`
	Servings    float64   `json:"servings"`
	AddShopping bool      `json:"addshopping"`
}
