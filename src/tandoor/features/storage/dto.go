package storage

// StorageResponse represents a storage integration (Dropbox, Nextcloud, Local) returned by Tandoor
type StorageResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Method    string  `json:"method"`
	Username  *string `json:"username"`
	URL       *string `json:"url"`
	Path      string  `json:"path"`
	CreatedBy int     `json:"created_by"`
}

// StorageListResponse represents a paginated list of storages
type StorageListResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []StorageResponse `json:"results"`
}

// StorageParam represents the payload to create or update a storage integration
type StorageParam struct {
	Name     string  `json:"name"`
	Method   string  `json:"method,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Token    *string `json:"token,omitempty"`
	URL      *string `json:"url,omitempty"`
	Path     string  `json:"path,omitempty"`
}

// ListParams represents query filters for listing storage integrations
type ListParams struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}
