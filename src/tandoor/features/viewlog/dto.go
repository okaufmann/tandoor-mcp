package viewlog

import "time"

// ViewLogParam represents the JSON payload to log a view event
type ViewLogParam struct {
	Recipe int `json:"recipe"`
}

// ViewLogResponse represents the ViewLog returned by Tandoor
type ViewLogResponse struct {
	ID        int       `json:"id"`
	Recipe    int       `json:"recipe"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// ViewLogListResponse represents a paginated list of ViewLogs
type ViewLogListResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []ViewLogResponse `json:"results"`
}

// ListParams represents query filters for fetching view logs
type ListParams struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}
