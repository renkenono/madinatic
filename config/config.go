package config

import (
	"database/sql"
	"encoding/json"
	"os"
)

// Config represents data needed to init server
type Config struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	DBName string `json:"db"`
	Proto  string `json:"proto"`
	BDAdr  string `json:"bd_adr"`
	Adr    string `json:"adr"`
	Pub    string `json:"pub"`
	DB     *sql.DB
}

const (
	config = "./config.json"
	// INFO prefix for logging info
	INFO = "INFO: "
	// WARN prefix for logging warnings
	WARN = "WARN: "
	// ERROR prefix for logging errors
	ERROR = "ERROR: "
	// FATAL prefix for loggin fatal errors
	FATAL = "FATAL: "
)

var (
	// App holds current config of the server
	App Config
)

// LoadConfig loads config to init server
func (c *Config) LoadConfig() error {
	conf, err := os.Open(config)
	if err != nil {
		return err
	}
	err = json.NewDecoder(conf).Decode(c)
	if err != nil {
		return err
	}

	// user:password@tcp(127.0.0.1:3306)/database
	dsn := App.User + ":" + App.Pass
	if App.Proto != "" {
		dsn += "@" + App.Proto + "(" + App.BDAdr + ")"
	}

	dsn += "/" + App.DBName

	App.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// ping db to check to verify conn
	err = App.DB.Ping()
	return err

}
