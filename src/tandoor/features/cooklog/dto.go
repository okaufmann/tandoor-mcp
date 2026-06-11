package cooklog

import "time"

// UserResponse represents nested user information in CookLog
type UserResponse struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
}

// CookLogParam represents the JSON payload to log a cooking event
type CookLogParam struct {
	Recipe    int        `json:"recipe"`
	Servings  *int       `json:"servings,omitempty"`
	Rating    *int       `json:"rating,omitempty"`
	Comment   *string    `json:"comment,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// CookLogResponse represents the CookLog returned by Tandoor
type CookLogResponse struct {
	ID        int          `json:"id"`
	Recipe    int          `json:"recipe"`
	Servings  *int         `json:"servings"`
	Rating    *int         `json:"rating"`
	Comment   *string      `json:"comment"`
	CreatedBy UserResponse `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// CookLogListResponse represents a paginated list of CookLogs
type CookLogListResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []CookLogResponse `json:"results"`
}

// ListParams represents query filters for fetching cook logs
type ListParams struct {
	Recipe   *int `json:"recipe,omitempty"`
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}
