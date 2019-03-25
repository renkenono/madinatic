package control

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/renkenn/madinatic/config"
)

const (
	ExpireTime = 5097600

	ErrUserDoesNotExist  = 1
	ErrWrongPassword     = iota
	ErrUserNotConfirmed  = iota
	ErrUserIDInvalid     = iota
	ErrUserIDExists      = iota
	ErrUsernameInvalid   = iota
	ErrUsernameExists    = iota
	ErrEmailInvalid      = iota
	ErrEmailExists       = iota
	ErrPhoneInvalid      = iota
	ErrPhoneExists       = iota
	ErrPassInvalid       = iota
	ErrRePassInvalid     = iota
	ErrFirstNameInvalid  = iota
	ErrFamilyNameInvalid = iota
)

const (
	ViewLogin    = 0
	ViewRegister = iota
	ViewReset    = iota
	ViewHomePage = iota
)

// Out holds messages sent to the end user
// consider making it a map of arrays of strings
// for supporting more than one language
// example: Out["en"][OutUserNotConfirmed]
// and Out["fr"][OutUserNotConfirmed]
var (
	Out = []string{
		"",
		"Utilisateur n'existe pas",
		"Mot de passe incorrecte",
		"Utilisateur n'est pas confirmé",
		"ID est invalid",
		"ID existe déja",
		"Username est invalid",
		"Username existe déja",
		"Email est invalid",
		"Email existe déja",
		"Phone est invalid",
		"Phone existe déja",
		"Mot de passe est invalid",
		"Mot de passe n'est pas le même",
		"Prénom est invalid",
		"Nom est invalid",
	}

	View = []string{
		"login",
		"register",
		"reset",
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
	err := tmpl.ExecuteTemplate(w, View[v], data)
	if err != nil {
		log.Fatalf("%stemplate error: %s", config.FATAL, err.Error())
	}

}

func DecodeJSON(w http.ResponseWriter, r *http.Request, data interface{}, suff string) error {
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s%s: %s", config.ERROR, suff, err.Error())
		return err
	}
	return nil
}

func MarshalJSON(w http.ResponseWriter, data interface{}, suff string) {
	djson, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.FATAL, suff, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(djson)
}
