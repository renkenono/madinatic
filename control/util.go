package control

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/renkenn/madinatic/config"
)

const (
	OutUserDoesNotExist = 0
	OutWrongPassword    = iota
	OutUserNotConfirmed = iota
	ViewAuth            = 0
	ViewHomePage        = iota
)

// Out holds messages sent to the end user
// consider making it a map of arrays of strings
// for supporting more than one language
// example: Out["en"][OutUserNotConfirmed]
// and Out["fr"][OutUserNotConfirmed]
var (
	Out = []string{
		"Utilisateur n'existe pas",
		"Mot de passe incorrecte",
		"Utilisateur n'est pas confirmé",
	}

	View = []string{
		"auth",
	}

	viewsPath = path.Join("web", "views")
)

// Render a view. failing results in panic
func Render(w http.ResponseWriter, r *http.Request, data interface{}, v int, fns ...string) {
	p := []string{viewsPath}
	p = append(p, fns...)
	tmpls := []string{path.Join(p...)}
	var tmpl *template.Template
	tmpl = template.Must(template.ParseFiles(tmpls...))

	// map[string]interface{} {} for readability
	err := tmpl.ExecuteTemplate(w, "auth", data)
	if err != nil {
		log.Fatalf("%stemplate error: %s", config.FATAL, err.Error())
	}

}
