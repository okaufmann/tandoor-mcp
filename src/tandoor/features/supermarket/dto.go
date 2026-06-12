package supermarket

import "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/shoppinglist"

// SupermarketCategoryResponse represents a category
type SupermarketCategoryResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	OpenDataSlug *string `json:"open_data_slug"`
}

// SupermarketCategoryListResponse represents a paginated list of SupermarketCategories
type SupermarketCategoryListResponse struct {
	Count    int                           `json:"count"`
	Next     *string                       `json:"next"`
	Previous *string                       `json:"previous"`
	Results  []SupermarketCategoryResponse `json:"results"`
}

// SupermarketCategoryParam represents request params to create/update a category
type SupermarketCategoryParam struct {
	Name         string  `json:"name"`
	Description  string  `json:"description,omitempty"`
	OpenDataSlug *string `json:"open_data_slug,omitempty"`
}

// CategoryRef is used to reference a category nested inside a relation
type CategoryRef struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// SupermarketCategoryRelationResponse represents a relation between supermarket and category
type SupermarketCategoryRelationResponse struct {
	ID          int                         `json:"id"`
	Category    SupermarketCategoryResponse `json:"category"`
	Supermarket int                         `json:"supermarket"`
	Order       int                         `json:"order"`
}

// SupermarketCategoryRelationListResponse represents a paginated list of relations
type SupermarketCategoryRelationListResponse struct {
	Count    int                                   `json:"count"`
	Next     *string                               `json:"next"`
	Previous *string                               `json:"previous"`
	Results  []SupermarketCategoryRelationResponse `json:"results"`
}

// SupermarketCategoryRelationParam represents request params to create a relation
type SupermarketCategoryRelationParam struct {
	Category    CategoryRef `json:"category"`
	Supermarket int         `json:"supermarket"`
	Order       int         `json:"order"`
}

// SupermarketResponse represents a supermarket
type SupermarketResponse struct {
	ID                    int                                   `json:"id"`
	Name                  string                                `json:"name"`
	Description           string                                `json:"description"`
	ShoppingLists         []shoppinglist.ShoppingListResponse   `json:"shopping_lists"`
	CategoryToSupermarket []SupermarketCategoryRelationResponse `json:"category_to_supermarket"`
	OpenDataSlug          *string                               `json:"open_data_slug"`
}

// SupermarketListResponse represents a paginated list of Supermarkets
type SupermarketListResponse struct {
	Count    int                   `json:"count"`
	Next     *string               `json:"next"`
	Previous *string               `json:"previous"`
	Results  []SupermarketResponse `json:"results"`
}

// SupermarketParam represents request params to create/update a supermarket
type SupermarketParam struct {
	Name          string                        `json:"name"`
	Description   string                        `json:"description,omitempty"`
	ShoppingLists []shoppinglist.ShoppingListID `json:"shopping_lists,omitempty"`
	OpenDataSlug  *string                       `json:"open_data_slug,omitempty"`
}

// ListParams represents query filters for fetching supermarkets, categories, or relations
type ListParams struct {
	Query    *string `json:"query,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"page_size,omitempty"`
}
