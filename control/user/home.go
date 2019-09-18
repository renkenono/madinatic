package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/renkenn/madinatic/model"

	"github.com/gorilla/csrf"

	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
)

type reportSlide struct {
	Title      string
	Address    string
	Categories []string
	Picture    string
	Link       string
}

const (
	homeViewErr = "homepage view failed"
	faqViewErr  = "faq view failed"
)

// Home page
func Home(w http.ResponseWriter, r *http.Request) {
	s, err := config.Store.Get(r, "userdata")
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, homeViewErr, err.Error())
		return
	}

	guest := false
	_, ok := s.Values["username"]
	if s.IsNew || !ok {
		// render guest homepage

		guest = true
	}

	// render user homepage
	rs, err := model.ReportsLatest(4)
	fmt.Println(rs)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, homeViewErr, err.Error())
		return
	}

	var rsslides [][]reportSlide
	j := 0
	var slide []reportSlide
	var done bool
	for i, re := range rs {
		done = false
		pics, err := re.Pics()
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.INFO, homeViewErr, err.Error())
			return
		}

		cats, err := re.Categories()
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.INFO, homeViewErr, err.Error())
			return
		}

		slide = append(slide, reportSlide{
			Title:      re.Title,
			Address:    re.Address,
			Categories: cats,
			Picture:    pics[0],
			Link:       fmt.Sprintf("/report/view/%d", re.ID),
		})
		if i%2 == 1 {
			rsslides = append(rsslides, slide)
			slide = nil
			j++
			done = true
		}
	}
	if !done {
		rsslides = append(rsslides, slide)
	}

	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"guest":     guest,
		"slides":    rsslides,
	}

	Render(w, r, data, ViewHome, "home.tmpl")

}

// FAQ page
func FAQ(w http.ResponseWriter, r *http.Request) {
	s, err := config.Store.Get(r, "userdata")
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, faqViewErr, err.Error())
		return
	}

	guest := false
	_, ok := s.Values["username"]

	if s.IsNew || !ok {
		// render guest homepage
		guest = true
	}

	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"guest":     guest,
	}

	Render(w, r, data, ViewFAQ, "faq.tmpl")

}

// Err page
func Err(w http.ResponseWriter, r *http.Request) {

	Render(w, r, nil, ViewErr, "err.tmpl")

}
