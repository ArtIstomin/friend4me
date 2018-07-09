package query

import (
	"github.com/labstack/echo"

	"github.com/artistomin/friend4me/internal"
)

// List prepares data for list queries
func List(u *model.AuthUser) (*model.ListQuery, error) {
	switch true {
	case int(u.Role) <= 2: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == model.ShelterAdminRole:
		return &model.ListQuery{Query: "shelter_id = ?", ID: u.ShelterID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
