// Package user contains user application services
package user

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/artistomin/friend4me/internal"
	"github.com/artistomin/friend4me/internal/auth"
	"github.com/artistomin/friend4me/internal/platform/query"
	"github.com/artistomin/friend4me/internal/platform/structs"
)

// New creates new user application service
func New(udb model.UserDB, rbac model.RBACService, auth model.AuthService) *Service {
	return &Service{udb: udb, rbac: rbac, auth: auth}
}

// Service represents user application service
type Service struct {
	udb  model.UserDB
	rbac model.RBACService
	auth model.AuthService
}

// Create creates a new user account
func (s *Service) Create(c echo.Context, req model.User) (*model.User, error) {
	if err := s.rbac.UserCreate(c, req.RoleID, req.CompanyID, req.LocationID); err != nil {
		return nil, err
	}
	req.Password = auth.HashPassword(req.Password)
	return s.udb.Create(req)
}

// ChangePassword changes user's password
func (s *Service) ChangePassword(c echo.Context, oldPass, newPass string, id int) error {
	if err := s.rbac.EnforceUser(c, id); err != nil {
		return err
	}
	u, err := s.udb.View(id)
	if err != nil {
		return err
	}
	if !auth.HashMatchesPassword(u.Password, oldPass) {
		return echo.NewHTTPError(http.StatusBadRequest, "old password is not correct")
	}
	u.Password = auth.HashPassword(newPass)
	return s.udb.ChangePassword(u)
}

// List returns list of users
func (s *Service) List(c echo.Context, p *model.Pagination) ([]model.User, error) {
	u := s.auth.User(c)
	q, err := query.List(u)
	if err != nil {
		return nil, err
	}
	return s.udb.List(q, p)
}

// View returns single user
func (s *Service) View(c echo.Context, id int) (*model.User, error) {
	if err := s.rbac.EnforceUser(c, id); err != nil {
		return nil, err
	}
	return s.udb.View(id)
}

// Delete deletes a user
func (s *Service) Delete(c echo.Context, id int) error {
	u, err := s.udb.View(id)
	if err != nil {
		return err
	}
	if err := s.rbac.IsLowerRole(c, u.Role.AccessLevel); err != nil {
		return err
	}
	return s.udb.Delete(u)
}

// Update contains user's information used for updating
type Update struct {
	ID        int
	FirstName *string
	LastName  *string
	Mobile    *string
	Phone     *string
	Address   *string
}

// Update updates user's contact information
func (s *Service) Update(c echo.Context, u *Update) (*model.User, error) {
	if err := s.rbac.EnforceUser(c, u.ID); err != nil {
		return nil, err
	}
	usr, err := s.udb.View(u.ID)
	if err != nil {
		return nil, err
	}
	structs.Merge(usr, u)
	return s.udb.Update(usr)
}
