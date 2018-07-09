package pgsql

import (
	"net/http"

	"github.com/go-pg/pg"
	"github.com/labstack/echo"

	"github.com/artistomin/friend4me/internal"
)

// NewUserDB returns a new UserDB instance
func NewUserDB(c *pg.DB, l echo.Logger) *UserDB {
	return &UserDB{c, l}
}

// UserDB represents the client for user table
type UserDB struct {
	cl  *pg.DB
	log echo.Logger
}

// Create creates a new user on database
func (u *UserDB) Create(usr model.User) (*model.User, error) {
	var user = new(model.User)
	res, err := u.cl.Query(user, "select id from users where username = ? or email = ? and deleted_at is null", usr.Username, usr.Email)
	if err != nil {
		u.log.Error("UserDB Error: %v", err)
		return nil, err
	}
	if res.RowsReturned() != 0 {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
	}
	if err := u.cl.Insert(&usr); err != nil {
		u.log.Error("UserDB Error: %v", err)
		return nil, err
	}
	return &usr, nil
}

// ChangePassword changes user's password
func (u *UserDB) ChangePassword(usr *model.User) error {
	_, err := u.cl.Model(usr).Column("password", "updated_at").WherePK().Update()
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return err
}

// View returns single user by ID
func (u *UserDB) View(id int) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	_, err := u.cl.QueryOne(user, sql, id)
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return user, err
}

// FindByUsername queries for single user by username
func (u *UserDB) FindByUsername(uname string) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	_, err := u.cl.QueryOne(user, sql, uname)
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return user, err
}

// FindByToken queries for single user by token
func (u *UserDB) FindByToken(token string) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`
	_, err := u.cl.QueryOne(user, sql, token)
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return user, err
}

// List returns list of all users retreivable for the current user, depending on role
func (u *UserDB) List(qp *model.ListQuery, p *model.Pagination) ([]model.User, error) {
	var users []model.User
	q := u.cl.Model(&users).Column("user.*", "Role").Limit(p.Limit).Offset(p.Offset).Where(notDeleted).Order("user.id desc")
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		u.log.Warnf("UserDB Error: %v", err)
		return nil, err
	}
	return users, nil
}

// Delete sets deleted_at for a user
func (u *UserDB) Delete(user *model.User) error {
	user.Delete()
	_, err := u.cl.Model(user).Column("deleted_at").WherePK().Update()
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return err
}

// Update updates user's contact info
func (u *UserDB) Update(user *model.User) (*model.User, error) {
	_, err := u.cl.Model(user).WherePK().UpdateNotNull()
	if err != nil {
		u.log.Warnf("UserDB Error: %v", err)
	}
	return user, err
}
