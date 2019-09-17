package model

import (
	"testing"

	"github.com/renkenn/madinatic/config"
)

func TestNewAuth(t *testing.T) {
	dsn, err := config.App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}

	_, errs := NewAuth("1", "admin", "admin@admin.com", "adminpass", "213555665555", "Admin")
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}
