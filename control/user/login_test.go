package user

import (
	"testing"
)

func TestNewAccessToken(t *testing.T) {

	token, err := newAccessToken("renken")
	t.Log(token)
	t.Error(token, err)
}
