package model

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/renkenn/madinatic/config"
)

func TestNewCitizen(t *testing.T) {
	dsn, err := config.App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}

	_, errs := NewCitizen("109891379001120119", "renken", "renkennate@gmail.com", "renken123", "213555555555", "Renken", "family name")
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestCitizens(t *testing.T) {
	dsn, err := config.App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}

	cs, err := Citizens()
	if err != nil {
		t.Error(err)
	}

	for _, c := range cs {
		t.Log(c)
	}
	t.Error("success")
}
