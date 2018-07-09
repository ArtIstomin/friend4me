package rbac_test

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/artistomin/friend4me/internal"
	"github.com/artistomin/friend4me/internal/mock"
	"github.com/artistomin/friend4me/internal/rbac"
)

func TestNew(t *testing.T) {
	rbacService := rbac.New(nil)
	if rbacService == nil {
		t.Error("RBAC Service not initialized")
	}
}

func TestEnforceRole(t *testing.T) {
	type args struct {
		ctx  echo.Context
		role model.AccessRole
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Not authorized",
			args: args{
				ctx:  mock.EchoCtxWithKeys([]string{"role"}, int8(3)),
				role: model.SuperAdminRole,
			},
			wantErr: true,
		},
		{
			name: "Authorized",
			args: args{
				ctx:  mock.EchoCtxWithKeys([]string{"role"}, int8(0)),
				role: model.ShelterAdminRole,
			},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(nil)
			res := rbacSvc.EnforceRole(tt.args.ctx, tt.args.role)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceUser(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same user, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 15, int8(3)), id: 122},
			wantErr: true,
		},
		{
			name:    "Not same user, but admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 22, int8(0)), id: 44},
			wantErr: false,
		},
		{
			name:    "Same user",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 8, int8(3)), id: 8},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(nil)
			res := rbacSvc.EnforceUser(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceShelter(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same shelter, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 7, int8(5)), id: 9},
			wantErr: true,
		},
		{
			name:    "Same shelter, not shelter admin or admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 22, int8(5)), id: 22},
			wantErr: true,
		},
		{
			name:    "Same shelter, shelter admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 5, int8(3)), id: 5},
			wantErr: false,
		},
		{
			name:    "Not same shelter but admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 8, int8(2)), id: 9},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(nil)
			res := rbacSvc.EnforceShelter(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestAccountCreate(t *testing.T) {
	type args struct {
		ctx        echo.Context
		roleID     int
		shelter_id int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Different shelter, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 2, int8(5)), roleID: 5, shelter_id: 7},
			wantErr: true,
		},
		{
			name:    "Same not shelter, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 2, int8(5)), roleID: 5, shelter_id: 2},
			wantErr: true,
		},
		{
			name:    "Different shelter, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 2, int8(3)), roleID: 4, shelter_id: 2},
			wantErr: false,
		},
		{
			name:    "Same shelter, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 2, int8(3)), roleID: 5, shelter_id: 2},
			wantErr: false,
		},
		{
			name:    "Same shelter, creating user role, admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 2, int8(3)), roleID: 5, shelter_id: 2},
			wantErr: false,
		},
		{
			name:    "Different everything, admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"shelter_id", "role"}, 2, int8(1)), roleID: 2, shelter_id: 7},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(nil)
			res := rbacSvc.AccountCreate(tt.args.ctx, tt.args.roleID, tt.args.shelter_id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestIsLowerRole(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{"role"}, int8(3))
	rbacSvc := rbac.New(nil)
	if rbacSvc.IsLowerRole(ctx, model.AccessRole(4)) != nil {
		t.Error("The requested user is higher role than the user requesting it")
	}
	if rbacSvc.IsLowerRole(ctx, model.AccessRole(2)) == nil {
		t.Error("The requested user is lower role than the user requesting it")
	}
}
