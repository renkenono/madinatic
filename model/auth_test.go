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

	_, errs := NewAuth("231321", "renauth", "renken@auth.dz", "renkenpass", "213555655555", "RenkenAuth")
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}
