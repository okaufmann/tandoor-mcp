package bookmarklet_import

import "time"

// BookmarkletImportParam represents the JSON payload to import via bookmarklet URL/HTML
type BookmarkletImportParam struct {
	Url  *string `json:"url,omitempty"`
	Html string  `json:"html"`
}

// BookmarkletImportResponse represents the created BookmarkletImport detail
type BookmarkletImportResponse struct {
	ID        int       `json:"id"`
	Url       *string   `json:"url"`
	Html      string    `json:"html"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// BookmarkletImportListResponse represents metadata returned in list responses
type BookmarkletImportListResponse struct {
	Count    int                            `json:"count"`
	Next     *string                        `json:"next"`
	Previous *string                        `json:"previous"`
	Results  []BookmarkletImportListElement `json:"results"`
}

// BookmarkletImportListElement represents a single summary of bookmarklet import in list responses
type BookmarkletImportListElement struct {
	ID        int       `json:"id"`
	Url       *string   `json:"url"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// ListParams represents query filters for bookmarklet imports list
type ListParams struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}
