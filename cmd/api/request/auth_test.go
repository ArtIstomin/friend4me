package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/artistomin/friend4me/cmd/api/request"
	"github.com/artistomin/friend4me/internal/mock"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		wantErr  bool
		wantData *request.Credentials
	}{
		{
			name: "Success",
			req:  `{"email":"johndoe@gmail.com","password":"hunter123"}`,
			wantData: &request.Credentials{
				Email:    "johndoe@gmail.com",
				Password: "hunter123",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "", bytes.NewBufferString(tt.req))
			c := mock.EchoCtx(req, w)
			resp, err := request.Login(c)
			assert.Equal(t, tt.wantData, resp)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
