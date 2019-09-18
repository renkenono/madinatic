package admin

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
	"github.com/renkenn/madinatic/model"
)

const (
	catCreateErr = "cat create failed"
	catsViewErr  = "cat view failed"
)

type replycat struct {
	I     int
	Name  string
	Auths []string
	Link  string
}

// CatsView returns a list of cats
func CatsView(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	cs, err := model.Cats()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, catsViewErr, err.Error())
		return
	}

	var replies []replycat

	i := 0
	for _, c := range cs {
		in := false
		for _, r := range replies {
			if r.Name == c.Name {
				in = true
			}
		}
		if !in {
			replies = append(replies, replycat{
				I:    i,
				Name: c.Name,
				Link: fmt.Sprintf("/cat/delete/%d", c.ID),
			})
			i++
		}
	}

	for i := 0; i < len(replies); i++ {
		for _, c := range cs {
			if replies[i].Name == c.Name {

				a, err := model.AuthByID(strconv.FormatUint(c.Auth, 10))
				if err != nil {
					http.Redirect(w, r, "/error", http.StatusInternalServerError)
					log.Printf("%s%s: %s", config.INFO, catsViewErr, err.Error())
					return
				}
				s := a.Name
				replies[i].Auths = append(replies[i].Auths, s)
			}
		}
	}

	// creating a category needs list of auths
	as, err := model.Auths()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, catsViewErr, err.Error())
		return
	}
	var auths []string
	for _, a := range as {
		auths = append(auths, a.Username)
	}

	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"Cats":      replies,
		"Auths":     auths,
	}

	Render(w, r, data, ViewDashboardCats, "d_cats.tmpl")
}

// CreateCat creates a new category
func CreateCat(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	// handle form

	log.Printf("%screate cat POST", config.INFO)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}

	selectedAuths := r.MultipartForm.Value["auth"]
	if len(selectedAuths) == 0 {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, catCreateErr, "no auth selected")
		return
	}

	name := r.MultipartForm.Value["cat"][0]

	for _, sa := range selectedAuths {
		_, err := model.NewCat(name, sa)
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.INFO, catCreateErr, err.Error())
			return
		}
	}

	http.Redirect(w, r, "/dashboard/auths", http.StatusFound)
}
