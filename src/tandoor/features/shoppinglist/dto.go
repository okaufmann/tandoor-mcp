package shoppinglist

import "time"

// ShoppingListResponse represents a shopping list returned by Tandoor
type ShoppingListResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Color       *string `json:"color"`
}

// ShoppingListListResponse represents a paginated list of ShoppingLists
type ShoppingListListResponse struct {
	Count    int                    `json:"count"`
	Next     *string                `json:"next"`
	Previous *string                `json:"previous"`
	Results  []ShoppingListResponse `json:"results"`
}

// ShoppingListParam represents the parameters to create or update a shopping list
type ShoppingListParam struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// FoodShopping represents the food detail nested inside shopping list entries
type FoodShopping struct {
	ID         int     `json:"id,omitempty"`
	Name       string  `json:"name"`
	PluralName *string `json:"plural_name,omitempty"`
}

// UnitResponse represents the unit detail nested inside shopping list entries
type UnitResponse struct {
	ID         int     `json:"id,omitempty"`
	Name       string  `json:"name"`
	PluralName *string `json:"plural_name,omitempty"`
}

// ShoppingListRecipe represents the recipe association in a shopping list entry
type ShoppingListRecipe struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Recipe   *int    `json:"recipe"`
	MealPlan *int    `json:"mealplan"`
	Servings float64 `json:"servings"`
}

// UserResponse represents the user detail nested inside shopping list entries
type UserResponse struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

// ShoppingListEntryResponse represents a shopping list entry/item
type ShoppingListEntryResponse struct {
	ID             int                    `json:"id"`
	ListRecipe     *int                   `json:"list_recipe"`
	ShoppingLists  []ShoppingListResponse `json:"shopping_lists"`
	Food           *FoodShopping          `json:"food"`
	Unit           *UnitResponse          `json:"unit"`
	Amount         float64                `json:"amount"`
	Order          int                    `json:"order"`
	Checked        bool                   `json:"checked"`
	Ingredient     *int                   `json:"ingredient"`
	ListRecipeData *ShoppingListRecipe    `json:"list_recipe_data"`
	CreatedBy      *UserResponse          `json:"created_by"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	CompletedAt    *time.Time             `json:"completed_at"`
	DelayUntil     *time.Time             `json:"delay_until"`
}

// ShoppingListEntryListResponse represents a paginated list of ShoppingListEntries
type ShoppingListEntryListResponse struct {
	Count    int                         `json:"count"`
	Next     *string                     `json:"next"`
	Previous *string                     `json:"previous"`
	Results  []ShoppingListEntryResponse `json:"results"`
}

// ShoppingListID is a helper for nested IDs
type ShoppingListID struct {
	ID int `json:"id"`
}

// CreateEntryParam represents the request body to create a new entry
type CreateEntryParam struct {
	ShoppingLists []ShoppingListID `json:"shopping_lists"`
	Food          FoodShopping     `json:"food"`
	Unit          *UnitResponse    `json:"unit,omitempty"`
	Amount        float64          `json:"amount"`
	Checked       *bool            `json:"checked,omitempty"`
}

// UpdateEntryParam represents the request body to patch an entry
type UpdateEntryParam struct {
	Amount  *float64      `json:"amount,omitempty"`
	Checked *bool         `json:"checked,omitempty"`
	Unit    *UnitResponse `json:"unit,omitempty"`
}

// ListEntriesParams represents query filters for fetching shopping list entries
type ListEntriesParams struct {
	Mealplan     *int    `json:"mealplan,omitempty"`
	UpdatedAfter *string `json:"updated_after,omitempty"`
	Page         *int    `json:"page,omitempty"`
	PageSize     *int    `json:"page_size,omitempty"`
}
