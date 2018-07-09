package request

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/artistomin/friend4me/internal"
)

// Register contains registration request
type Register struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Username        string `json:"username" validate:"required,min=3,alphanum"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	Email           string `json:"email" validate:"required,email"`

	CompanyID int `json:"company_id" validate:"required"`
	RoleID    int `json:"role_id" validate:"required"`
}

// UserCreate validates user creation request
func UserCreate(c echo.Context) (*Register, error) {
	r := new(Register)
	if err := c.Bind(r); err != nil {
		return nil, err
	}
	if r.Password != r.PasswordConfirm {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
	}
	if r.RoleID < int(model.SuperAdminRole) || r.RoleID > int(model.UserRole) {
		return nil, echo.NewHTTPError(http.StatusBadRequest)
	}
	return r, nil
}

// Password contains password change request
type Password struct {
	ID                 int    `json:"-"`
	OldPassword        string `json:"old_password" validate:"required,min=8"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required"`
}

// PasswordChange validates password change request
func PasswordChange(c echo.Context) (*Password, error) {
	id, err := ID(c)
	if err != nil {
		return nil, err
	}
	p := new(Password)
	if err := c.Bind(p); err != nil {
		return nil, err
	}
	if p.NewPassword != p.NewPasswordConfirm {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
	}
	p.ID = id
	return p, nil
}

// UpdateUser contains user update data from json request
type UpdateUser struct {
	ID        int     `json:"-"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2"`
	Mobile    *string `json:"mobile,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Address   *string `json:"address,omitempty"`
}

// UserUpdate validates user update request
func UserUpdate(c echo.Context) (*UpdateUser, error) {
	id, err := ID(c)
	if err != nil {
		return nil, err
	}
	u := new(UpdateUser)
	if err := c.Bind(u); err != nil {
		return nil, err
	}
	u.ID = id
	return u, nil
}
