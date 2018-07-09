package model

// AccessRole represents access role type
type AccessRole int8

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = iota + 1

	// AdminRole has admin specific permissions
	AdminRole

	// ShelterAdminRole can edit shelter specific things
	ShelterAdminRole

	// AdopterRole is a standard user
	AdopterRole
)

// Role model
type Role struct {
	ID          int        `json:"id"`
	AccessLevel AccessRole `json:"access_level"`
	Name        string     `json:"name"`
}
