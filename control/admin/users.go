package admin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
	"github.com/renkenn/madinatic/model"
)

const (
	usersViewErr = "users view failed"
	userBanErr   = "user ban failed"
)

type replyc struct {
	I int
	*model.Citizen
	Link string
}

// UsersView returns a list of users
func UsersView(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	// us
	us, err := model.Citizens()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, usersViewErr, err.Error())
		return
	}
	// check if admin/auth/client

	var replies []replyc

	for i, u := range us {
		replies = append(replies, replyc{
			i,
			u,
			fmt.Sprintf("/user/delete/%d", u.ID),
		})
	}

	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"Users":     replies,
	}

	Render(w, r, data, ViewDashboardUsers, "d_users.tmpl")
}

// UserDelete deletes a user
// /user/delete/[id]
func UserDelete(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	vars := mux.Vars(r)
	u, err := model.UserByID(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, userBanErr, err.Error())
		return
	}

	err = u.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, userBanErr, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
