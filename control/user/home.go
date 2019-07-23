package user

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
)

func Home(w http.ResponseWriter, r *http.Request) {
	s, err := config.Store.Get(r, "userdata")
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, logErr, err.Error())
		return
	}

	fstr := "welcome guest."
	if !s.IsNew {
		// render guest homepage

		fstr = "welcome " + s.Values["username"].(string)
	}

	// render user homepage

	fstr += s.Flashes()[0].(string)
	
	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"flash":     fstr,
	}

	Render(w, r, data, ViewHome, "home.tmpl")

}
