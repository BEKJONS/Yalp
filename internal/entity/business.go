package entity

// Business category options
const (
	CategoryRestaurant    = "Restaurant"
	CategoryRetail        = "Retail"
	CategoryService       = "Service"
	CategoryHealthcare    = "Healthcare"
	CategoryEntertainment = "Entertainment"
)

// Business entity
type Business struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Location           Location `json:"location"`
	Category           string   `json:"category"`
	Description        string   `json:"description"`
	ContactInformation string   `json:"contact_information"`
	Attachments        []string `json:"attachments"`
	CreatedBy          string   `json:"created_by"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
	DeletedAt          string   `json:"deleted_at,omitempty"` // can be null
}

// Location entity for latitude and longitude
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Request parameters for single business entity actions
type BusinessSingleRequest struct {
	ID string `json:"id"`
}

// Response structure for a list of businesses
type BusinessList struct {
	Items []Business `json:"businesses"`
	Count int        `json:"count"`
}
