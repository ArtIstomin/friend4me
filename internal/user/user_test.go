package user_test

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/artistomin/friend4me/internal"
	"github.com/artistomin/friend4me/internal/mock"
	"github.com/artistomin/friend4me/internal/mock/mockdb"
	"github.com/artistomin/friend4me/internal/user"
)

func TestCreate(t *testing.T) {
	type args struct {
		c   echo.Context
		req model.User
	}
	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData *model.User
		udb      *mockdb.User
		rbac     *mock.RBAC
	}{{
		name: "Fail on is lower role",
		rbac: &mock.RBAC{
			UserCreateFn: func(echo.Context, int, int, int) error {
				return model.ErrGeneric
			}},
		wantErr: true,
		args: args{req: model.User{
			FirstName: "John",
			LastName:  "Doe",
			Username:  "JohnDoe",
			RoleID:    1,
			Password:  "Thranduil8822",
		}},
	},
		{
			name: "Success",
			args: args{req: model.User{
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
				RoleID:    1,
				Password:  "Thranduil8822",
			}},
			udb: &mockdb.User{
				CreateFn: func(u model.User) (*model.User, error) {
					u.CreatedAt = mock.TestTime(2000)
					u.UpdatedAt = mock.TestTime(2000)
					u.Base.ID = 1
					return &u, nil
				},
			},
			rbac: &mock.RBAC{
				UserCreateFn: func(echo.Context, int, int, int) error {
					return nil
				}},
			wantData: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
				RoleID:    1,
			}}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(tt.udb, tt.rbac, nil)
			usr, err := s.Create(tt.args.c, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				tt.wantData.Password = usr.Password
				assert.Equal(t, tt.wantData, usr)
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	type args struct {
		c       echo.Context
		oldpass string
		newpass string
		id      int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
		udb     *mockdb.User
		rbac    *mock.RBAC
	}{
		{
			name: "Fail on EnforceUser",
			args: args{id: 1},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return model.ErrGeneric
				}},
			wantErr: true,
		},
		{
			name:    "Fail on ViewUser",
			args:    args{id: 1},
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, model.ErrGeneric
				},
			},
		},
		{
			name: "Fail on PasswordMatch",
			args: args{id: 1, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Password: "IncorrectHashedPassword",
					}, nil
				},
			},
		},
		{
			name: "Success",
			args: args{id: 1, oldpass: "hunter123", newpass: "password"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Password: "$2a$10$udRBroNGBeOYwSWCVzf6Lulg98uAoRCIi4t75VZg84xgw6EJbFNsG",
					}, nil
				},
				// Check whether password was hashed correctly
				ChangePasswordFn: func(usr *model.User) error {
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(tt.udb, tt.rbac, nil)
			err := s.ChangePassword(tt.args.c, tt.args.oldpass, tt.args.newpass, tt.args.id)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestView(t *testing.T) {
	type args struct {
		c  echo.Context
		id int
	}
	cases := []struct {
		name     string
		args     args
		wantData *model.User
		wantErr  error
		udb      *mockdb.User
		rbac     *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{id: 5},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return model.ErrGeneric
				}},
			wantErr: model.ErrGeneric,
		},
		{
			name: "Success",
			args: args{id: 1},
			wantData: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
			},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					if id == 1 {
						return &model.User{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(2000),
								UpdatedAt: mock.TestTime(2000),
							},
							FirstName: "John",
							LastName:  "Doe",
							Username:  "JohnDoe",
						}, nil
					}
					return nil, nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(tt.udb, tt.rbac, nil)
			usr, err := s.View(tt.args.c, tt.args.id)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestList(t *testing.T) {
	type args struct {
		c   echo.Context
		pgn *model.Pagination
	}
	cases := []struct {
		name     string
		args     args
		wantData []model.User
		wantErr  bool
		udb      *mockdb.User
		auth     *mock.Auth
	}{
		{
			name: "Fail on query List",
			args: args{c: nil, pgn: &model.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			wantErr: true,
			auth: &mock.Auth{
				UserFn: func(c echo.Context) *model.AuthUser {
					return &model.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       model.UserRole,
					}
				}}},
		{
			name: "Success",
			args: args{c: nil, pgn: &model.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			auth: &mock.Auth{
				UserFn: func(c echo.Context) *model.AuthUser {
					return &model.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       model.AdminRole,
					}
				}},
			udb: &mockdb.User{
				ListFn: func(*model.ListQuery, *model.Pagination) ([]model.User, error) {
					return []model.User{
						{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(1999),
								UpdatedAt: mock.TestTime(2000),
							},
							FirstName: "John",
							LastName:  "Doe",
							Email:     "johndoe@gmail.com",
							Username:  "johndoe",
						},
						{
							Base: model.Base{
								ID:        2,
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
							},
							FirstName: "Hunter",
							LastName:  "Logan",
							Email:     "logan@aol.com",
							Username:  "hunterlogan",
						},
					}, nil
				}},
			wantData: []model.User{
				{
					Base: model.Base{
						ID:        1,
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
					},
					FirstName: "John",
					LastName:  "Doe",
					Email:     "johndoe@gmail.com",
					Username:  "johndoe",
				},
				{
					Base: model.Base{
						ID:        2,
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
					},
					FirstName: "Hunter",
					LastName:  "Logan",
					Email:     "logan@aol.com",
					Username:  "hunterlogan",
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(tt.udb, nil, tt.auth)
			usrs, err := s.List(tt.args.c, tt.args.pgn)
			assert.Equal(t, tt.wantData, usrs)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

}

func TestDelete(t *testing.T) {
	type args struct {
		c  echo.Context
		id int
	}
	cases := []struct {
		name    string
		args    args
		wantErr error
		udb     *mockdb.User
		rbac    *mock.RBAC
	}{
		{
			name:    "Fail on ViewUser",
			args:    args{id: 1},
			wantErr: model.ErrGeneric,
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, model.ErrGeneric
				},
			},
		},
		{
			name: "Fail on RBAC",
			args: args{id: 1},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Role: &model.Role{
							AccessLevel: model.UserRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, model.AccessRole) error {
					return model.ErrGeneric
				}},
			wantErr: model.ErrGeneric,
		},
		{
			name: "Success",
			args: args{id: 1},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Role: &model.Role{
							AccessLevel: model.AdminRole,
							ID:          2,
							Name:        "Admin",
						},
					}, nil
				},
				DeleteFn: func(usr *model.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, model.AccessRole) error {
					return nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(tt.udb, tt.rbac, nil)
			err := s.Delete(tt.args.c, tt.args.id)
			if err != tt.wantErr {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		c   echo.Context
		upd *user.Update
	}
	cases := []struct {
		name     string
		args     args
		wantData *model.User
		wantErr  error
		udb      *mockdb.User
		rbac     *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{upd: &user.Update{
				ID: 1,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return model.ErrGeneric
				}},
			wantErr: model.ErrGeneric,
		},
		{
			name: "Fail on ViewUser",
			args: args{upd: &user.Update{
				ID: 1,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantErr: model.ErrGeneric,
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, model.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{upd: &user.Update{
				ID:        1,
				FirstName: mock.Str2Ptr("John"),
				LastName:  mock.Str2Ptr("Doe"),
				Mobile:    mock.Str2Ptr("123456"),
				Phone:     mock.Str2Ptr("234567"),
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantData: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(1990),
					UpdatedAt: mock.TestTime(2000),
				},
				CompanyID:  1,
				LocationID: 2,
				RoleID:     3,
				FirstName:  "John",
				LastName:   "Doe",
				Mobile:     "123456",
				Phone:      "234567",
				Address:    "Work Address",
				Email:      "golang@go.org",
			},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					if id == 1 {
						return &model.User{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(1990),
								UpdatedAt: mock.TestTime(1991),
							},
							CompanyID:  1,
							LocationID: 2,
							RoleID:     3,
							FirstName:  "Joanna",
							LastName:   "Doep",
							Mobile:     "334455",
							Phone:      "444555",
							Address:    "Work Address",
							Email:      "golang@go.org",
						}, nil
					}
					return nil, model.ErrGeneric
				},
				UpdateFn: func(usr *model.User) (*model.User, error) {
					usr.UpdatedAt = mock.TestTime(2000)
					return usr, nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(tt.udb, tt.rbac, nil)
			usr, err := s.Update(tt.args.c, tt.args.upd)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
