package admin

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
	"github.com/renkenn/madinatic/model"
)

type reportSlide struct {
	I          int
	ID         uint64
	Title      string
	Address    string
	Categories []string
	Picture    string
	Link       string
	Accept     string
	Delete     string
}

const (
	reportsViewErr   = "reports view failed"
	reportsDeleteErr = "report delete failed"
	reportStateErr   = "report state failed"
)

// ApprovedReportsView returns list of reports
func ApprovedReportsView(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
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
	}

	Render(w, r, data, ViewDashboardAccReports, "d_reports_accepted.tmpl")

}

// PendingReportsView returns list of reports
func PendingReportsView(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	rs, err := model.ReportsByState(model.ReportPending)
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
			I:          i,
			ID:         re.ID,
			Title:      re.Title,
			Address:    re.Address,
			Categories: cats,
			Picture:    pics[0],
			Link:       fmt.Sprintf("/report/view/%d", re.ID),
			Accept:     fmt.Sprintf("/report/accept/%d", re.ID),
			Delete:     fmt.Sprintf("/report/delete/%d", re.ID),
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
	}

	Render(w, r, data, ViewDashboardPenReports, "d_reports_pending.tmpl")

}

// ReportDelete deletes given report
// user must be the admin
// /report/[id]/delete
func ReportDelete(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	vars := mux.Vars(r)
	re, err := model.ReportByID(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsDeleteErr, err.Error())
		return
	}

	err = re.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsDeleteErr, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ReportState sets a report's state
// to the given state
// /report/[id]/state/[state]
func ReportState(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	vars := mux.Vars(r)
	state, err := strconv.ParseUint(vars["state"], 10, 8)
	if err != nil && state >= model.ReportSolved {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportStateErr, err.Error())
		return
	}

	re, err := model.ReportByID(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportStateErr, err.Error())
		return
	}
	err = re.SetState(uint8(state))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportStateErr, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Accept reorients the report
func Accept(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	vars := mux.Vars(r)
	re, err := model.ReportByID(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, userBanErr, err.Error())
		return
	}

	cats, err := re.Categories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, userBanErr, err.Error())
		return
	}

	newCats := r.MultipartForm.Value["cats"]

	for _, c := range cats {
		rm := true
		for _, nc := range newCats {
			if c == nc {
				rm = false
			}
		}
		if rm {
			currc, err := model.CatByName(c)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("%s%s: %s", config.INFO, userBanErr, err.Error())
				return
			}

			// sacrificing the architecture
			config.DB.Lock()

			stmt, err := config.DB.Prepare("DELETE FROM subreports WHERE fk_reportid = ? AND fk_catid = ?")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("%s%s: %s", config.INFO, userBanErr, err.Error())
				return
			}

			_, err = stmt.Exec(re.ID, currc.ID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("%s%s: %s", config.INFO, userBanErr, err.Error())
				return
			}

			config.DB.Unlock()

		}
	}

	err = re.SetState(model.ReportAccepted)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, userBanErr, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
