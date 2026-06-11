package recipe_import

import "time"

// StorageParam represents the nested storage parameters required when creating a recipe import
type StorageParam struct {
	Name     string  `json:"name"`
	Method   *string `json:"method,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Path     string  `json:"path,omitempty"`
	Url      *string `json:"url,omitempty"`
}

// StorageResponse represents storage location metadata returned in responses
type StorageResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Method    string  `json:"method"`
	Username  *string `json:"username"`
	Path      string  `json:"path"`
	Url       *string `json:"url"`
	CreatedBy int     `json:"created_by"`
}

// RecipeImportParam represents the JSON payload to create a recipe import record
type RecipeImportParam struct {
	Storage  StorageParam `json:"storage"`
	Name     string       `json:"name"`
	FileUID  string       `json:"file_uid,omitempty"`
	FilePath string       `json:"file_path,omitempty"`
}

// RecipeImportResponse represents the created RecipeImport record
type RecipeImportResponse struct {
	ID        int             `json:"id"`
	Storage   StorageResponse `json:"storage"`
	Name      string          `json:"name"`
	FileUID   string          `json:"file_uid"`
	FilePath  string          `json:"file_path"`
	CreatedAt time.Time       `json:"created_at"`
}

// RecipeImportListResponse represents a paginated list of RecipeImports
type RecipeImportListResponse struct {
	Count    int                    `json:"count"`
	Next     *string                `json:"next"`
	Previous *string                `json:"previous"`
	Results  []RecipeImportResponse `json:"results"`
}

// ListParams represents query filters for recipe imports list
type ListParams struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}
