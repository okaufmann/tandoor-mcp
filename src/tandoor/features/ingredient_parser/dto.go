package ingredient_parser

// FoodSimple represents simple food metadata returned by the parser
type FoodSimple struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	PluralName *string `json:"plural_name"`
}

// UnitResponse represents unit metadata returned by the parser
type UnitResponse struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	PluralName *string `json:"plural_name"`
}

// IngredientSimple represents a parsed ingredient
type IngredientSimple struct {
	ID           int           `json:"id"`
	Food         *FoodSimple   `json:"food"`
	Unit         *UnitResponse `json:"unit"`
	Amount       float64       `json:"amount"`
	Note         *string       `json:"note"`
	Order        int           `json:"order"`
	IsHeader     bool          `json:"is_header"`
	NoAmount     bool          `json:"no_amount"`
	OriginalText *string       `json:"original_text"`
}

// IngredientParserRequest represents the payload to parse ingredients
type IngredientParserRequest struct {
	Ingredient  string   `json:"ingredient,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
}

// IngredientParserResponse represents the parsed ingredient(s) returned by Tandoor
type IngredientParserResponse struct {
	Ingredient  *IngredientSimple  `json:"ingredient"`
	Ingredients []IngredientSimple `json:"ingredients"`
}
