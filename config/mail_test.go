package config

import "testing"

func TestMail(t *testing.T) {
	App.LoadConfig("../config.json")
	m := New("renken@tfwno.gf", "Test Mail Subject", "Test Mail body")
	err := m.Send()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("mail sent!")
	}

}
