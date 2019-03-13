package config

import "testing"

func TestConfig(t *testing.T) {
	dsn, err := App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(dsn)
	t.Log(App)

}
