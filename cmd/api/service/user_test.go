package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/artistomin/friend4me/cmd/api/server"
	"github.com/artistomin/friend4me/cmd/api/service"
	"github.com/artistomin/friend4me/internal"
	"github.com/artistomin/friend4me/internal/auth"
	"github.com/artistomin/friend4me/internal/mock"
	"github.com/artistomin/friend4me/internal/mock/mockdb"
	"github.com/artistomin/friend4me/internal/user"
)

func TestListUsers(t *testing.T) {
	type listResponse struct {
		Users []model.User `json:"users"`
		Page  int          `json:"page"`
	}
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *listResponse
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			req:        `?limit=2222&page=-1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on query list",
			req:  `?limit=100&page=1`,
			auth: &mock.Auth{
				UserFn: func(c echo.Context) *model.AuthUser {
					return &model.AuthUser{
						ID:        1,
						ShelterID: 2,
						Role:      model.AdopterRole,
						Email:     "john@mail.com",
					}
				}},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `?limit=100&page=1`,
			auth: &mock.Auth{
				UserFn: func(c echo.Context) *model.AuthUser {
					return &model.AuthUser{
						ID:        1,
						ShelterID: 2,
						Role:      model.SuperAdminRole,
						Email:     "john@mail.com",
					}
				}},
			udb: &mockdb.User{
				ListFn: func(q *model.ListQuery, p *model.Pagination) ([]model.User, error) {
					if p.Limit == 100 && p.Offset == 100 {
						return []model.User{
							{
								Base: model.Base{
									ID:        10,
									CreatedAt: mock.TestTime(2001),
									UpdatedAt: mock.TestTime(2002),
								},
								FirstName: "John",
								LastName:  "Doe",
								Email:     "john@mail.com",
								ShelterID: 2,
								Role: &model.Role{
									ID:          1,
									AccessLevel: 1,
									Name:        "SUPER_ADMIN",
								},
							},
							{
								Base: model.Base{
									ID:        11,
									CreatedAt: mock.TestTime(2004),
									UpdatedAt: mock.TestTime(2005),
								},
								FirstName: "Joanna",
								LastName:  "Dye",
								Email:     "joanna@mail.com",
								ShelterID: 1,
								Role: &model.Role{
									ID:          2,
									AccessLevel: 2,
									Name:        "ADMIN",
								},
							},
						}, nil
					}
					return nil, model.ErrGeneric
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &listResponse{
				Users: []model.User{
					{
						Base: model.Base{
							ID:        10,
							CreatedAt: mock.TestTime(2001),
							UpdatedAt: mock.TestTime(2002),
						},
						FirstName: "John",
						LastName:  "Doe",
						Email:     "john@mail.com",
						ShelterID: 2,
						Role: &model.Role{
							ID:          1,
							AccessLevel: 1,
							Name:        "SUPER_ADMIN",
						},
					},
					{
						Base: model.Base{
							ID:        11,
							CreatedAt: mock.TestTime(2004),
							UpdatedAt: mock.TestTime(2005),
						},
						FirstName: "Joanna",
						LastName:  "Dye",
						Email:     "joanna@mail.com",
						ShelterID: 1,
						Role: &model.Role{
							ID:          2,
							AccessLevel: 2,
							Name:        "ADMIN",
						},
					},
				}, Page: 1},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("/v1/users")
			service.NewUser(user.New(tt.udb, tt.rbac, tt.auth), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(listResponse)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestViewUser(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *model.User
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			req:        `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, int) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, int) error {
					return nil
				},
			},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Email:     "johndoe@gmail.com",
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Email:     "johndoe@gmail.com",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("/v1/users")
			service.NewUser(user.New(tt.udb, tt.rbac, tt.auth), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(model.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		id         string
		wantStatus int
		wantResp   *model.User
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","mobile":"123456","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, int) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, int) error {
					return nil
				},
			},
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Email:     "johndoe@gmail.com",
						Address:   "Work",
						Phone:     "332223",
					}, nil
				},
				UpdateFn: func(usr *model.User) (*model.User, error) {
					usr.UpdatedAt = mock.TestTime(2010)
					usr.Mobile = "991991"
					return usr, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2010),
				},
				FirstName: "jj",
				LastName:  "okocha",
				Email:     "johndoe@gmail.com",
				Phone:     "321321",
				Address:   "home",
				Mobile:    "991991",
			},
		},
	}

	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("/v1/users")
			service.NewUser(user.New(tt.udb, tt.rbac, tt.auth), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.id
			req, _ := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(model.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		wantStatus int
		udb        *mockdb.User
		rbac       *mock.RBAC
		auth       *mock.Auth
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Role: &model.Role{
							AccessLevel: model.ShelterAdminRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, model.AccessRole) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Role: &model.Role{
							AccessLevel: model.ShelterAdminRole,
						},
					}, nil
				},
				DeleteFn: func(*model.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, model.AccessRole) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("/v1/users")
			service.NewUser(user.New(tt.udb, tt.rbac, tt.auth), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.id
			req, _ := http.NewRequest("DELETE", path, nil)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestCreate(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *model.User
		rbac       *mock.RBAC
		udb        *mockdb.User
	}{
		{
			name:       "Invalid request",
			req:        `{"first_name":"John","last_name":"Doe","password":"hunter123","password_confirm":"hunter1234","email":"johndoe@gmail.com","shelter_id":1,"role_id":3}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on userSvc",
			req:  `{"first_name":"John","last_name":"Doe","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","shelter_id":1,"role_id":2}`,
			rbac: &mock.RBAC{
				UserCreateFn: func(c echo.Context, roleID, shelterID int) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"first_name":"John","last_name":"Doe","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","shelter_id":1,"role_id":2}`,
			rbac: &mock.RBAC{
				UserCreateFn: func(c echo.Context, roleID, shelterID int) error {
					return nil
				},
			},
			udb: &mockdb.User{
				CreateFn: func(usr model.User) (*model.User, error) {
					usr.ID = 1
					usr.CreatedAt = mock.TestTime(2018)
					usr.UpdatedAt = mock.TestTime(2018)
					return &usr, nil
				},
			},
			wantResp: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2018),
					UpdatedAt: mock.TestTime(2018),
				},
				FirstName: "John",
				LastName:  "Doe",
				Email:     "johndoe@gmail.com",
				ShelterID: 1,
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("/v1/users")
			service.NewUser(user.New(tt.udb, tt.rbac, nil), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(model.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestChangePassword(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		id         string
		udb        *mockdb.User
		rbac       *mock.RBAC
	}{
		{
			name:       "Invalid request",
			req:        `{"new_password":"new_password","old_password":"my_old_password", "new_password_confirm":"new_password_cf"}`,
			wantStatus: http.StatusBadRequest,
			id:         "1",
		},
		{
			name: "Fail on RBAC",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return echo.ErrForbidden
				},
			},
			id:         "1",
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				},
			},
			id: "1",
			udb: &mockdb.User{
				ViewFn: func(id int) (*model.User, error) {
					return &model.User{
						Password: auth.HashPassword("oldpassw"),
					}, nil
				},
				ChangePasswordFn: func(usr *model.User) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("/v1/users")
			service.NewUser(user.New(tt.udb, tt.rbac, nil), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users/" + tt.id + "/password"
			req, err := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
