package report

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
	"github.com/renkenn/madinatic/model"
)

type reportDetail struct {
	ID         uint64
	Title      string
	Address    string
	Desc       string
	Categories []string
	Pictures   []string
	IsSolved   bool
	IsAccepted bool
	IsPending  bool
	Solve      string
	CreatedAt  string // t.Format(time.UnixDate)
	ModifiedAt string
	Username   string
	UserCname  string
	IsAuth     bool
}

const (
	viewErr    = "view report failed"
	viewAPIErr = "view report API failed"
	solveErr   = "solve report failed"
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
	var rd reportDetail
	rd.ID = re.ID
	rd.Title = re.Title
	rd.Desc = re.Desc
	rd.Address = re.Address
	rd.CreatedAt = re.CreatedAt.Format(time.RFC822)
	rd.ModifiedAt = re.ModifiedAt.Format(time.RFC822)
	user, in := IsLoggedIn(w, r)

	pics, err := re.Pics()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
		return
	}
	rd.Pictures = pics

	cats, err := re.Categories()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
		return
	}

	rd.Categories = cats
	u, err := model.UserByIDi(re.UID)
	if err != nil {
		if err != model.ErrUserDoesNotExist {
			log.Printf("%s%s: %s", config.ERROR, viewErr, err.Error())
		}
		http.Redirect(w, r, "/error", http.StatusBadRequest)
		return
	}
	cname, err := u.Cname()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
		return
	}

	re.Username = u.Username
	re.UserCName = cname
	rd.Username = u.Username
	rd.UserCname = cname

	// bools should default to false in golang
	if re.State == model.ReportSolved {
		rd.IsSolved = true
	} else if re.State == model.ReportAccepted {
		rd.IsAccepted = true
	} else {
		rd.IsPending = true
	}

	rd.Solve = fmt.Sprintf("/report/solve/%d", re.ID)
	yes, err := isAuthofReport(user, re)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, reportsViewErr, err.Error())
		return
	}

	rd.IsAuth = yes
	// add guest
	data := map[string]interface{}{
		"report":   rd,
		"guest":    !in,
		"username": user,
	}

	Render(w, r, data, ViewReportDetails, "report.tmpl")
}

func isAuthofReport(username string, r *model.Report) (bool, error) {
	auths, err := r.Auths()
	if err != nil {
		return false, err
	}
	for _, a := range auths {

		if a == username {
			return true, nil
		}
	}
	return false, nil
}

// Solve partially/completely solves the report.
// it expects a call from an auth responsible
// of one subreports
func Solve(w http.ResponseWriter, r *http.Request) {
	user, in := IsLoggedIn(w, r)
	if !in {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	re, err := model.ReportByID(vars["id"])
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, solveErr, err.Error())
		return
	}

	yes, err := isAuthofReport(user, re)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, solveErr, err.Error())
		return
	}
	if !yes {
		http.Redirect(w, r, "/error", http.StatusUnauthorized)
		return
	}

	// here we go again

	err = re.Solve(user)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, solveErr, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/report/view/%d", re.ID), http.StatusFound)
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
