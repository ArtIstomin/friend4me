package model

// Shelter represents shelter model
type Shelter struct {
	Base
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Owner   User   `json:"owner"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
	About   string `json:"about,omitempty"`
}
