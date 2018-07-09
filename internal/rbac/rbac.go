package rbac

import (
	"github.com/labstack/echo"

	"github.com/artistomin/friend4me/internal"
)

// New creates new RBAC service
func New(udb model.UserDB) *Service {
	return &Service{udb}
}

// Service is RBAC application service
type Service struct {
	udb model.UserDB
}

func checkBool(b bool) error {
	if b {
		return nil
	}
	return echo.ErrForbidden
}

func (s *Service) isAdmin(c echo.Context) bool {
	return !(c.Get("role").(int8) > int8(model.AdminRole))
}

func (s *Service) isShelterAdmin(c echo.Context) bool {
	// Must query shelter ID in database for the given user
	return !(c.Get("role").(int8) > int8(model.ShelterAdminRole))
}

// EnforceRole authorizes request by AccessRole
func (s *Service) EnforceRole(c echo.Context, r model.AccessRole) error {
	return checkBool(!(c.Get("role").(int8) > int8(r)))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c echo.Context, ID int) error {
	// TODO: Implement querying db and checking the requested user's shelter_id
	// to allow shelter admins to view the user
	if s.isAdmin(c) {
		return nil
	}
	return checkBool(c.Get("id").(int) == ID)
}

// EnforceShelter checks whether the request to apply change to shelter data
// is done by the user belonging to the that shelter and that the user has role ShelterAdmin.
// If user has admin role, the check for shelter doesnt need to pass.
func (s *Service) EnforceShelter(c echo.Context, ID int) error {
	if s.isAdmin(c) {
		return nil
	}

	if err := s.EnforceRole(c, model.ShelterAdminRole); err != nil {
		return err
	}

	return checkBool(c.Get("shelter_id").(int) == ID)
}

// AccountCreate performs auth check when creating a new account
func (s *Service) AccountCreate(c echo.Context, roleID, shelterID int) error {
	return s.IsLowerRole(c, model.AccessRole(roleID))
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for user creation/deletion
func (s *Service) IsLowerRole(c echo.Context, r model.AccessRole) error {
	return checkBool(c.Get("role").(int8) < int8(r))
}
