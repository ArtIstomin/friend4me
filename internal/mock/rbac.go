package mock

import (
	"github.com/labstack/echo"

	"github.com/artistomin/friend4me/internal"
)

// RBAC Mock
type RBAC struct {
	EnforceRoleFn    func(echo.Context, model.AccessRole) error
	EnforceUserFn    func(echo.Context, int) error
	EnforceShelterFn func(echo.Context, int) error
	AccountCreateFn  func(echo.Context, int, int) error
	IsLowerRoleFn    func(echo.Context, model.AccessRole) error
}

// EnforceRole mock
func (a *RBAC) EnforceRole(c echo.Context, role model.AccessRole) error {
	return a.EnforceRoleFn(c, role)
}

// EnforceUser mock
func (a *RBAC) EnforceUser(c echo.Context, id int) error {
	return a.EnforceUserFn(c, id)
}

// EnforceShelter mock
func (a *RBAC) EnforceShelter(c echo.Context, id int) error {
	return a.EnforceShelterFn(c, id)
}

// AccountCreate mock
func (a *RBAC) AccountCreate(c echo.Context, roleID, shelterID int) error {
	return a.AccountCreateFn(c, roleID, shelterID)
}

// IsLowerRole mock
func (a *RBAC) IsLowerRole(c echo.Context, role model.AccessRole) error {
	return a.IsLowerRoleFn(c, role)
}
