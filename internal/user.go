package model

import (
	"time"
)

// User represents user domain model
type User struct {
	Base
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Password  string     `json:"-"`
	Email     string     `json:"email"`
	Mobile    string     `json:"mobile,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Address   string     `json:"address,omitempty"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	Active    bool       `json:"active"`
	Token     string     `json:"-"`

	Role *Role `json:"role,omitempty"`

	RoleID    int `json:"-"`
	ShelterID int `json:"shelter_id,omitempty"`
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID        int
	ShelterID int
	Email     string
	Role      AccessRole
}

// UpdateLastLogin updates last login field
func (u *User) UpdateLastLogin() {
	t := time.Now()
	u.LastLogin = &t
}

// UserDB represents user database interface (repository)
type UserDB interface {
	Create(User) (*User, error)
	ChangePassword(*User) error
	View(int) (*User, error)
	FindByEmail(string) (*User, error)
	FindByToken(string) (*User, error)
	List(*ListQuery, *Pagination) ([]User, error)
	Delete(*User) error
	Update(*User) (*User, error)
}
