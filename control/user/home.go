package user

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
)

// Home page
func Home(w http.ResponseWriter, r *http.Request) {
	s, err := config.Store.Get(r, "userdata")
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, logErr, err.Error())
		return
	}

	fstr := "welcome guest."
	_, ok := s.Values["username"]
	if !s.IsNew && !ok {
		// render guest homepage

		fstr = "welcome " + s.Values["username"].(string)
		// fstr += s.Flashes()[0].(string)
	}

	// render user homepage

	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"flash":     fstr,
	}

	Render(w, r, data, ViewHome, "home.tmpl")

}
