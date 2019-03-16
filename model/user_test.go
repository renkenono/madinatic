package model

import "testing"

func TestValidateUsername(t *testing.T) {
	for _, user := range []string{"renken", "renken13", "ren-chan", "_renchan", "1renken", "renkenennenenenenennenenenennenenennenenen"} {
		err := ValidateUsername(user)
		t.Log(user)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	for _, email := range []string{"renken@siga.dz", "_renken@si ga.dz"} {
		err := ValidateEmail(email)
		t.Log(email)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestValidatePass(t *testing.T) {
	for _, pass := range []string{"spaces are valid", "length", "\u0000controlcodepoint", "bobandalicédon'texist"} {
		err := ValidatePass(pass)
		t.Log("password: " + pass)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestValidatePhone(t *testing.T) {
	for _, p := range []string{"+213555555555", "+231555555555", "+213455555555", "+21374545"} {
		err := ValidatePhone(p)
		t.Log("phone: " + p)
		if err != nil {
			t.Error(err)
		}
	}
}
