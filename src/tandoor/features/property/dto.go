package property

// PropertyTypeResponse represents a property type in Tandoor
type PropertyTypeResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Unit         *string `json:"unit"`
	Description  *string `json:"description"`
	Order        int     `json:"order"`
	OpenDataSlug *string `json:"open_data_slug"`
	FdcID        *int    `json:"fdc_id"`
}

// PropertyTypeListResponse represents a paginated list of property types
type PropertyTypeListResponse struct {
	Count    int                    `json:"count"`
	Next     *string                `json:"next"`
	Previous *string                `json:"previous"`
	Results  []PropertyTypeResponse `json:"results"`
}

// PropertyTypeParam represents request parameters to create a property type
type PropertyTypeParam struct {
	Name         string  `json:"name"`
	Unit         *string `json:"unit,omitempty"`
	Description  *string `json:"description,omitempty"`
	Order        int     `json:"order,omitempty"`
	OpenDataSlug *string `json:"open_data_slug,omitempty"`
	FdcID        *int    `json:"fdc_id,omitempty"`
}

// PropertyTypeRef is a reference to a property type for nested structures
type PropertyTypeRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PropertyResponse represents a food property in Tandoor
type PropertyResponse struct {
	ID             int                  `json:"id"`
	PropertyAmount *float64             `json:"property_amount"`
	PropertyType   PropertyTypeResponse `json:"property_type"`
}

// PropertyListResponse represents a paginated list of properties
type PropertyListResponse struct {
	Count    int                `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []PropertyResponse `json:"results"`
}

// PropertyParam represents request parameters to create a property
type PropertyParam struct {
	PropertyAmount *float64        `json:"property_amount"`
	PropertyType   PropertyTypeRef `json:"property_type"`
}

// ListPropertyTypesParams represents query parameters for listing property types
type ListPropertyTypesParams struct {
	Category []string `json:"category,omitempty"`
	Page     *int     `json:"page,omitempty"`
	PageSize *int     `json:"page_size,omitempty"`
}

// ListPropertiesParams represents query parameters for listing properties
type ListPropertiesParams struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}
