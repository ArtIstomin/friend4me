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

// EnforceRole authorizes request by AccessRole
func (s *Service) EnforceRole(c echo.Context, r model.AccessRole) error {
	return checkBool(!(c.Get("role").(int8) > int8(r)))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c echo.Context, ID int) error {
	// TODO: Implement querying db and checking the requested user's shelter_id
	// to allow company admins to view the user
	if s.isAdmin(c) {
		return nil
	}
	return checkBool(c.Get("id").(int) == ID)
}

// EnforceCompany checks whether the request to apply change to company data
// is done by the user belonging to the that company and that the user has role CompanyAdmin.
// If user has admin role, the check for company doesnt need to pass.
func (s *Service) EnforceCompany(c echo.Context, ID int) error {
	if s.isAdmin(c) {
		return nil
	}
	if err := s.EnforceRole(c, model.ShelterAdminRole); err != nil {
		return err
	}
	return checkBool(c.Get("shelter_id").(int) == ID)
}

// UserCreate performs auth check when creating a new user
func (s *Service) UserCreate(c echo.Context, roleID, shelterID int) error {
	return s.IsLowerRole(c, model.AccessRole(roleID))
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for user creation/deletion
func (s *Service) IsLowerRole(c echo.Context, r model.AccessRole) error {
	return checkBool(c.Get("role").(int8) < int8(r))
}
