package model

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/renkenn/madinatic/config"
)

func TestValidateUserID(t *testing.T) {
	for _, id := range []string{"0", "32142142155151515151515", "213213214", "109891379001120119", "999999999999999999", "9999999999999999999", "h"} {
		nid, err := ValidateUserID(id)
		t.Log(id)
		if err != nil {
			t.Error(err)
		}
		t.Log(nid)
	}
}
func TestValidateUsername(t *testing.T) {
	for _, user := range []string{"renken", "renken13", "ren-chan", "_renchan", "1renken", "renkenennenenenenennenenenennenenennenenen"} {
		_, err := ValidateUsername(user)
		t.Log(user)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	for _, email := range []string{"renken@siga.dz", "_renken@si ga.dz"} {
		_, err := ValidateEmail(email)
		t.Log(email)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestValidatePass(t *testing.T) {
	for _, pass := range []string{"spaces are valid", "length", "\u0000controlcodepoint", "bobandalicédon'texist"} {
		_, err := ValidatePass(pass)
		t.Log("password: " + pass)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestValidatePhone(t *testing.T) {
	for _, p := range []string{"213555555555", "231555555555", "213455555555", "21374545"} {
		np, err := ValidatePhone(p)
		t.Log("phone: " + p)
		t.Log(np)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestExistsUserID(t *testing.T) {
	dsn, err := config.App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}

	for _, id := range []string{"101", "102"} {
		err := ExistsUserID(id)
		t.Log(id)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestExistsUsername(t *testing.T) {
	dsn, err := config.App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}

	for _, name := range []string{"renkenn", "admin"} {
		err := ExistsUsername(name)
		t.Log(name)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestExistsPhone(t *testing.T) {
	dsn, err := config.App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}

	for _, phone := range []string{"213555555555", "72324324"} {
		err := ExistsPhone(phone)
		t.Log(phone)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestExistsEmail(t *testing.T) {
	dsn, err := config.App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}

	for _, email := range []string{"renken@email.com", "renken@siga.dz"} {
		err := ExistsEmail(email)
		t.Log(email)
		if err != nil {
			t.Error(err)
		}
	}
}
