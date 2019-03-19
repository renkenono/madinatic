package config

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestConfig(t *testing.T) {
	dsn, err := App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(dsn)
	t.Log(App)

	err = DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}
}
