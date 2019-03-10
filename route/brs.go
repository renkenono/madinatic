package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
)

// BrowserRoutes assings html handlers to a given mux
func BrowserRoutes(r *mux.Router) {
	pub := http.FileServer(http.Dir(config.App.Pub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", pub))
}
