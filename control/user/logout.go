package user

import (
	"log"
	"net/http"

	"github.com/renkenn/madinatic/config"
)

// Logout and delete session
func Logout(w http.ResponseWriter, r *http.Request) {
	s, err := config.Store.Get(r, "userdata")
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, logErr, err.Error())
		return
	}

	s.Values["username"] = ""
	delete(s.Values, "username")
	s.Options.MaxAge = -1
	http.Redirect(w, r, "/", http.StatusFound)
}
