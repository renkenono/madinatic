package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
	"github.com/renkenn/madinatic/route"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	log.Printf("%sloading config", config.INFO)
	err := config.App.LoadConfig()
	if err != nil {
		log.Panicf("%sfailed to load config: %s", config.FATAL, err.Error())
	}

	// close db connection when server is done
	defer config.App.DB.Close()

	log.Printf("%ssuccesfully loaded config %v", config.INFO, config.App)

	r := mux.NewRouter()
	route.BrowserRoutes(r)
	route.APIRoutes(r)
	s := http.Server{
		Addr:    config.App.Adr,
		Handler: r,
	}

	log.Printf("%sstarting server at %s", config.INFO, s.Addr)
	log.Println(s.ListenAndServe())

}
