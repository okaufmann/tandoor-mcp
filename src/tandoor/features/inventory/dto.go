package inventory

import "time"

// HouseholdResponse represents a household
type HouseholdResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// HouseholdListResponse represents a list of households
type HouseholdListResponse struct {
	Count   int                 `json:"count"`
	Results []HouseholdResponse `json:"results"`
}

// InventoryLocationResponse represents an inventory location
type InventoryLocationResponse struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	IsFreezer bool              `json:"is_freezer"`
	Household HouseholdResponse `json:"household"`
}

// InventoryLocationListResponse represents a paginated list of locations
type InventoryLocationListResponse struct {
	Count    int                         `json:"count"`
	Next     *string                     `json:"next"`
	Previous *string                     `json:"previous"`
	Results  []InventoryLocationResponse `json:"results"`
}

// InventoryLocationParam represents parameters to create or update an inventory location
type InventoryLocationParam struct {
	Name      string            `json:"name"`
	IsFreezer bool              `json:"is_freezer"`
	Household HouseholdResponse `json:"household"`
}

// FoodRef is a simple reference to a food item
type FoodRef struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// UnitRef is a simple reference to a unit
type UnitRef struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// InventoryEntryResponse represents an entry in the inventory
type InventoryEntryResponse struct {
	ID                int                       `json:"id"`
	InventoryLocation InventoryLocationResponse `json:"inventory_location"`
	SubLocation       *string                   `json:"sub_location"`
	Code              *string                   `json:"code"`
	Food              FoodRef                   `json:"food"`
	Unit              UnitRef                   `json:"unit"`
	Amount            float64                   `json:"amount"`
	Expires           *string                   `json:"expires"` // YYYY-MM-DD
	Note              *string                   `json:"note"`
	Label             string                    `json:"label"`
	CreatedAt         time.Time                 `json:"created_at"`
	CreatedBy         int                       `json:"created_by"`
}

// InventoryEntryListResponse represents a paginated list of inventory entries
type InventoryEntryListResponse struct {
	Count    int                      `json:"count"`
	Next     *string                  `json:"next"`
	Previous *string                  `json:"previous"`
	Results  []InventoryEntryResponse `json:"results"`
}

// InventoryEntryParam represents parameters to create an inventory entry
type InventoryEntryParam struct {
	InventoryLocation InventoryLocationResponse `json:"inventory_location"`
	Food              FoodRef                   `json:"food"`
	Unit              UnitRef                   `json:"unit"`
	Amount            float64                   `json:"amount"`
	SubLocation       *string                   `json:"sub_location,omitempty"`
	Code              *string                   `json:"code,omitempty"`
	Expires           *string                   `json:"expires,omitempty"`
	Note              *string                   `json:"note,omitempty"`
}

// InventoryEntryUpdateParam represents parameters to update/patch an inventory entry
type InventoryEntryUpdateParam struct {
	Amount            *float64                   `json:"amount,omitempty"`
	InventoryLocation *InventoryLocationResponse `json:"inventory_location,omitempty"`
	Note              *string                    `json:"note,omitempty"`
}

// InventoryLogResponse represents an inventory transaction log
type InventoryLogResponse struct {
	ID                   int                       `json:"id"`
	Entry                InventoryEntryResponse    `json:"entry"`
	BookingType          string                    `json:"booking_type"`
	OldAmount            float64                   `json:"old_amount"`
	NewAmount            float64                   `json:"new_amount"`
	OldInventoryLocation InventoryLocationResponse `json:"old_inventory_location"`
	NewInventoryLocation InventoryLocationResponse `json:"new_inventory_location"`
	Note                 *string                   `json:"note"`
	CreatedAt            time.Time                 `json:"created_at"`
}

// InventoryLogListResponse represents a paginated list of inventory logs
type InventoryLogListResponse struct {
	Count    int                    `json:"count"`
	Next     *string                `json:"next"`
	Previous *string                `json:"previous"`
	Results  []InventoryLogResponse `json:"results"`
}

// ListLocationsParams represents query parameters for listing locations
type ListLocationsParams struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}

// ListEntriesParams represents query parameters for listing entries
type ListEntriesParams struct {
	InventoryLocationID *int  `json:"inventory_location_id,omitempty"`
	FoodID              *int  `json:"food_id,omitempty"`
	Empty               *bool `json:"empty,omitempty"`
	Page                *int  `json:"page,omitempty"`
	PageSize            *int  `json:"page_size,omitempty"`
}

// ListLogsParams represents query parameters for listing logs
type ListLogsParams struct {
	EntryID  *int `json:"entry_id,omitempty"`
	FoodID   *int `json:"food_id,omitempty"`
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}
