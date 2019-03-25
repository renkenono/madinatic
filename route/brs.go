package route

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
	"github.com/renkenn/madinatic/control/user"
)

// BrowserRoutes assings html handlers to a given mux
func BrowserRoutes(r *mux.Router) {
	pub := http.FileServer(http.Dir(config.App.Pub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", pub))
	CSRF := csrf.Protect([]byte(config.App.SignKey), csrf.Secure(false))
	gen := r.PathPrefix("/").Subrouter()
	gen.Use(CSRF)
	gen.HandleFunc("/login", user.Login).Methods("POST", "GET")
	gen.HandleFunc("/register", user.Register).Methods("POST", "GET")
	gen.HandleFunc("/confirm/{id:[0-9]+}/{token}", user.Confirm).Methods("GET")
	gen.HandleFunc("/reset/{id:[0-9]+}/{token}", user.ResetPass).Methods("GET", "POST")
	gen.HandleFunc("/reset", user.Reset).Methods("POST")
}
