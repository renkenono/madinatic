package route

import (
	"net/http"

	"github.com/renkenn/madinatic/control/admin"
	"github.com/renkenn/madinatic/control/report"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
	"github.com/renkenn/madinatic/control/user"
)

// BrowserRoutes assings html handlers to a given mux
func BrowserRoutes(r *mux.Router) {
	pub := http.FileServer(http.Dir(config.App.Pub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", pub))
	CSRF := csrf.Protect([]byte(config.App.SignKey), csrf.Secure(false))
	gen := r.PathPrefix("/").Subrouter()
	gen.Use(CSRF)
	gen.HandleFunc("/faq", user.FAQ).Methods("GET")
	gen.HandleFunc("/error", user.Err).Methods("GET")
	gen.HandleFunc("/login", user.Login).Methods("POST", "GET")
	gen.HandleFunc("/register", user.Register).Methods("POST", "GET")
	gen.HandleFunc("/confirm/{id:[0-9]+}/{token}", user.Confirm).Methods("GET")
	gen.HandleFunc("/reset/{id:[0-9]+}/{token}", user.ResetPass).Methods("GET", "POST")
	gen.HandleFunc("/reset", user.Reset).Methods("POST")
	gen.HandleFunc("/", user.Home).Methods("GET")
	gen.HandleFunc("/settings", user.Settings).Methods("GET", "POST")
	gen.HandleFunc("/logout", user.Logout).Methods("GET", "POST")

	gen.HandleFunc("/reports", report.ReportsView).Methods("GET")
	gen.HandleFunc("/report/create", report.Create).Methods("GET", "POST")
	gen.HandleFunc("/report/view/{id:[0-9]+}", report.ViewReport).Methods("GET")
	gen.HandleFunc("/report/accept/{id:[0-9]+}", admin.Accept).Methods("GET", "POST")
	gen.HandleFunc("/report/solve/{id:[0-9]+}", report.Solve).Methods("GET")

	gen.HandleFunc("/dashboard", admin.PendingReportsView).Methods("GET")
	gen.HandleFunc("/dashboard/users", admin.UsersView).Methods("GET")
	gen.HandleFunc("/dashboard/auths", admin.AuthsView).Methods("GET")
	gen.HandleFunc("/dashboard/cats", admin.CatsView).Methods("GET")
	gen.HandleFunc("/dashboard/auth/create", admin.AuthCreate).Methods("GET", "POST")
	gen.HandleFunc("/dashboard/cat/create", admin.CreateCat).Methods("POST")
	gen.HandleFunc("/user/delete/{id:[0-9]+}", admin.UserDelete).Methods("GET")
	gen.HandleFunc("/dashboard/reports/accepted", admin.ApprovedReportsView).Methods("GET")
	gen.HandleFunc("/dashboard/reports/pending", admin.PendingReportsView).Methods("GET")
	gen.HandleFunc("/report/delete/{id:[0-9]+}", admin.ReportDelete).Methods("GET")
}
