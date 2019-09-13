package model

import (
	"testing"

	"github.com/renkenn/madinatic/config"
)

func TestNewCat(t *testing.T) {
	dsn, err := config.App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}

	_, err = NewCat("renauth category", "renauth")
	if err != nil {
		t.Error(err)
	}
}
