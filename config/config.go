package config

import (
	"database/sql"
	"encoding/json"
	"os"
	"sync"
)

// Config represents data needed to init server
type Config struct {
	User    string `json:"user"`
	Pass    string `json:"pass"`
	DBName  string `json:"db"`
	Proto   string `json:"proto"`
	BDAdr   string `json:"bd_adr"`
	Adr     string `json:"adr"`
	Pub     string `json:"pub"`
	ESrv    string `json:"esrv"`
	EMail   string `json:"email"`
	EPass   string `json:"epass"`
	SignKey string `json:"sign_key"`
}

// DBC represents lock protected db connection
type DBC struct {
	*sql.DB
	*sync.Mutex
}

const (
	// ConfigFile path to config
	ConfigFile = "./config.json"
	// INFO prefix for logging info
	INFO = "INFO: "
	// WARN prefix for logging warnings
	WARN = "WARN: "
	// ERROR prefix for logging errors
	ERROR = "ERROR: "
	// FATAL prefix for loggin fatal errors
	FATAL = "FATAL: "
	// NAME holds app name
	NAME = "Madina-TIC"
)

var (
	// App holds current config of the server
	App Config
	// DB holds database
	DB DBC
)

// LoadConfig loads config to init server
func (c *Config) LoadConfig(path string) (string, error) {
	conf, err := os.Open(path)
	if err != nil {
		return "", err
	}
	err = json.NewDecoder(conf).Decode(c)
	if err != nil {
		return "", err
	}

	// user:password@tcp(127.0.0.1:3306)/database
	dsn := App.User + ":" + App.Pass
	if App.Proto != "" {
		dsn += "@" + App.Proto + "(" + App.BDAdr + ")"
	}

	dsn += "/" + App.DBName + "?parseTime=true"
	return dsn, nil
}

// InitDB opens a connection and checks if it's working
func (db *DBC) InitDB(dsn string) (err error) {

	db.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	DB.Mutex = new(sync.Mutex)
	// ping db to check to verify conn
	err = db.Ping()
	return err
}
