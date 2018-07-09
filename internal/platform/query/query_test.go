package query_test

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/artistomin/friend4me/internal"
	"github.com/artistomin/friend4me/internal/platform/query"
)

func TestList(t *testing.T) {
	type args struct {
		user *model.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *model.ListQuery
		wantErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: &model.AuthUser{
				Role: model.SuperAdminRole,
			}},
		},
		{
			name: "Shelter admin user",
			args: args{user: &model.AuthUser{
				Role:      model.ShelterAdminRole,
				ShelterID: 1,
			}},
			wantData: &model.ListQuery{
				Query: "shelter_id = ?",
				ID:    1},
		},
		{
			name: "Adopter user",
			args: args{user: &model.AuthUser{
				Role: model.AdopterRole,
			}},
			wantErr: echo.ErrForbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.List(tt.args.user)
			assert.Equal(t, tt.wantData, q)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
