package route

import (
	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/control/user"
)

// APIRoutes assings API handlers to a given mux
func APIRoutes(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/login", user.LoginAPI).Methods("POST")
	api.HandleFunc("/register", user.RegisterAPI).Methods("POST")
	api.HandleFunc("/settings", user.SettingsAPI).Methods("POST")
}
