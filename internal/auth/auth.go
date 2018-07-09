package auth

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"

	"github.com/artistomin/friend4me/internal"
)

// New creates new auth service
func New(udb model.UserDB, j JWT) *Service {
	return &Service{
		udb: udb,
		jwt: j,
	}
}

// Service represents auth application service
type Service struct {
	udb model.UserDB
	jwt JWT
}

// JWT represents jwt interface
type JWT interface {
	GenerateToken(*model.User) (string, string, error)
}

// Authenticate tries to authenticate the user provided by username and password
func (s *Service) Authenticate(c echo.Context, user, pass string) (*model.AuthToken, error) {
	u, err := s.udb.FindByUsername(user)
	if err != nil {
		return nil, err
	}
	if !HashMatchesPassword(u.Password, pass) {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Username or password does not exist")
	}

	if !u.Active {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}
	token, expire, err := s.jwt.GenerateToken(u)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}

	u.UpdateLastLogin()
	u.Token = xid.New().String()
	_, err = s.udb.Update(u)
	if err != nil {
		return nil, err
	}

	return &model.AuthToken{Token: token, Expires: expire, RefreshToken: u.Token}, nil
}

// Refresh refreshes jwt token and puts new claims inside
func (s *Service) Refresh(c echo.Context, token string) (*model.RefreshToken, error) {
	user, err := s.udb.FindByToken(token)
	if err != nil {
		return nil, err
	}
	token, expire, err := s.jwt.GenerateToken(user)
	if err != nil {
		return nil, model.ErrGeneric
	}
	return &model.RefreshToken{Token: token, Expires: expire}, nil
}

// Me returns info about currently logged user
func (s *Service) Me(c echo.Context) (*model.User, error) {
	au := s.User(c)
	return s.udb.View(au.ID)
}

// User returns user data stored in jwt token
func (s *Service) User(c echo.Context) *model.AuthUser {
	id := c.Get("id").(int)
	companyID := c.Get("company_id").(int)
	user := c.Get("username").(string)
	email := c.Get("email").(string)
	role := c.Get("role").(int8)
	return &model.AuthUser{
		ID:        id,
		Username:  user,
		CompanyID: companyID,
		Email:     email,
		Role:      model.AccessRole(role),
	}
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// HashMatchesPassword matches hash with password. Returns true if hash and password match.
func HashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
