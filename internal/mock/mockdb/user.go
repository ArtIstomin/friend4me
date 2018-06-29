package mockdb

import (
	"github.com/artistomin/friend4me/internal"
)

// User database mock
type User struct {
	CreateFn         func(model.User) (*model.User, error)
	ChangePasswordFn func(*model.User) error
	ViewFn           func(int) (*model.User, error)
	FindByUsernameFn func(string) (*model.User, error)
	FindByTokenFn    func(string) (*model.User, error)
	ListFn           func(*model.ListQuery, *model.Pagination) ([]model.User, error)
	DeleteFn         func(*model.User) error
	UpdateFn         func(*model.User) (*model.User, error)
}

// Create mock
func (u *User) Create(usr model.User) (*model.User, error) {
	return u.CreateFn(usr)
}

// ChangePassword mock
func (u *User) ChangePassword(usr *model.User) error {
	return u.ChangePasswordFn(usr)
}

// View mock
func (u *User) View(id int) (*model.User, error) {
	return u.ViewFn(id)
}

// FindByUsername mock
func (u *User) FindByUsername(username string) (*model.User, error) {
	return u.FindByUsernameFn(username)
}

// FindByToken mock
func (u *User) FindByToken(token string) (*model.User, error) {
	return u.FindByTokenFn(token)
}

// List mock
func (u *User) List(lq *model.ListQuery, p *model.Pagination) ([]model.User, error) {
	return u.ListFn(lq, p)
}

// Delete mock
func (u *User) Delete(usr *model.User) error {
	return u.DeleteFn(usr)
}

// Update mock
func (u *User) Update(usr *model.User) (*model.User, error) {
	return u.UpdateFn(usr)
}
