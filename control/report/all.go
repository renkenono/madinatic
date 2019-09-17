package report

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
	"github.com/renkenn/madinatic/model"
)

type reportSlide struct {
	I          int
	Title      string
	Address    string
	Categories []string
	Picture    string
	Link       string
}

const (
	reportsViewErr = "reports view failed"
)

// ReportsView returns list of reports
func ReportsView(w http.ResponseWriter, r *http.Request) {

	s, err := config.Store.Get(r, "userdata")
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
		return
	}

	guest := false
	_, ok := s.Values["username"]

	if s.IsNew || !ok {
		// render guest homepage
		guest = true
	}

	rs, err := model.ReportsByState(model.ReportAccepted)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
		return
	}

	/* move inside the loop in case needed in the future
	u, err := model.UserByIDi(re.UID)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
		return
	}

	re.UserCName, err = u.Cname()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
		return
	}
	*/

	var rsslides [][]reportSlide
	var slide []reportSlide
	done := false
	for i, re := range rs {
		pics, err := re.Pics()
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
			return
		}

		cats, err := re.Categories()
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
			return
		}

		slide = append(slide, reportSlide{
			Title:      re.Title,
			Address:    re.Address,
			Categories: cats,
			Picture:    pics[0],
			Link:       fmt.Sprintf("/report/view/%d", re.ID),
		})
		if i%3 == 2 {
			rsslides = append(rsslides, slide)
			done = true
		}
	}
	if !done {
		rsslides = append(rsslides, slide)

	}

	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"slides":    rsslides,
		"guest":     guest,
	}

	Render(w, r, data, ViewReports, "reports.tmpl")

}
