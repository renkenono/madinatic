package report

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
	"github.com/renkenn/madinatic/model"
)

const (
	viewErr    = "view report failed"
	viewAPIErr = "view report API failed"
)

// ViewReport report
func ViewReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	re, err := model.ReportByID(vars["id"])
	if err != nil {
		if err != model.ErrReportDoesNotExist {
			log.Printf("%s%s: %s", config.ERROR, viewErr, err.Error())
		}
		http.Redirect(w, r, "/error", http.StatusBadRequest)
		return

	}
	u, err := model.UserByIDi(re.UID)
	if err != nil {
		if err != model.ErrUserDoesNotExist {
			log.Printf("%s%s: %s", config.ERROR, viewErr, err.Error())
		}
		http.Redirect(w, r, "/error", http.StatusBadRequest)
		return
	}
	re.Username = u.Username
	w.Write([]byte(fmt.Sprintln(re)))
}

// ViewReportAPI report
func ViewReportAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	re, err := model.ReportByID(vars["id"])
	if err != nil {
		if err != model.ErrReportDoesNotExist {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("%s%s: %s", config.ERROR, viewAPIErr, err.Error())
		return

	}
	u, err := model.UserByIDi(re.UID)
	if err != nil {
		if err != model.ErrUserDoesNotExist {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("%s%s: %s", config.ERROR, viewAPIErr, err.Error())
		return
	}
	re.Username = u.Username
	w.Write([]byte(fmt.Sprintln(re)))
}
