package model_test

import (
	"testing"

	"github.com/artistomin/friend4me/internal"
)

func TestUpdateLastLogin(t *testing.T) {
	user := &model.User{
		FirstName: "TestGuy",
	}
	user.UpdateLastLogin()
	if user.LastLogin.IsZero() {
		t.Errorf("Last login time was not changed")
	}
}
