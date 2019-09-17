package main

import (
	"log"
	"net/http"

	"github.com/renkenn/madinatic/model"

	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
	"github.com/renkenn/madinatic/route"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	log.Printf("%sloading config", config.INFO)
	dsn, err := config.App.LoadConfig(config.ConfigFile)
	if err != nil {
		log.Panicf("%sfailed to load config: %s", config.FATAL, err.Error())
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		log.Panicf("%sfailed to init db connection: %s", config.FATAL, err.Error())
	}

	// create the admin's account
	model.NewAuth("1", "admin", "admin@admin.com", "admin", "213566643311", "Admin")
	// _, errs := model.NewAuth("1", "admin", "admin@admin.com", "admin", "213566643311", "Admin")
	// if len(errs) > 0 {
	// 	for _, e := range errs {
	// 		fmt.Println(e.Error())
	// 	}
	// }

	// close db connection when server is done
	defer config.DB.Close()

	log.Printf("%ssuccesfully loaded config %v", config.INFO, config.App)

	r := mux.NewRouter()
	route.APIRoutes(r)
	route.BrowserRoutes(r)
	s := http.Server{
		Addr:    config.App.Adr,
		Handler: r,
	}

	log.Printf("%sstarting server at %s", config.INFO, s.Addr)
	log.Println(s.ListenAndServe())

}
