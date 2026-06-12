package keyword

// KeywordResponse represents a keyword in Tandoor
type KeywordResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FullName    string `json:"full_name"`
	Parent      *int   `json:"parent"`
	NumChild    int    `json:"numchild"`
}

// KeywordListResponse represents a paginated list of keywords
type KeywordListResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []KeywordResponse `json:"results"`
}

// KeywordParam represents request parameters to create a keyword
type KeywordParam struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parent      *int   `json:"parent,omitempty"`
}

// ListParams represents query parameters for listing keywords
type ListParams struct {
	Query    *string `json:"query,omitempty"`
	Root     *int    `json:"root,omitempty"`
	Tree     *int    `json:"tree,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"page_size,omitempty"`
}
