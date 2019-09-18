package user

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
	"github.com/renkenn/madinatic/model"
)

const (
	regAPIErr = "register API failed"
	regErr    = "register failed"
)

type resp struct {
	Errors []uint `json:"errors"`
	*model.Citizen
}

type req struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Pass       string `json:"password"`
	FirstName  string `json:"first_name"`
	FamilyName string `json:"family_name"`
}

// RegisterAPI allows registration over JSON
func RegisterAPI(w http.ResponseWriter, r *http.Request) {

	var re req
	err := DecodeJSON(w, r, &re, regAPIErr)
	if err != nil {
		return
	}

	var rep resp
	c, cerrs := model.NewCitizen(re.ID, re.Username, re.Email, re.Pass, re.Phone, re.FirstName, re.FamilyName)
	if len(cerrs) > 0 {
		for _, err := range cerrs {
			switch err {
			case model.ErrUserIDInvalid:
				rep.Errors = append(rep.Errors, ErrUserIDInvalid)
			case model.ErrUserIDExists:
				rep.Errors = append(rep.Errors, ErrUserIDExists)
			case model.ErrUsernameInvalid:
				rep.Errors = append(rep.Errors, ErrUsernameInvalid)
			case model.ErrUsernameExists:
				rep.Errors = append(rep.Errors, ErrUsernameExists)
			case model.ErrEmailInvalid:
				rep.Errors = append(rep.Errors, ErrEmailInvalid)
			case model.ErrEmailExists:
				rep.Errors = append(rep.Errors, ErrEmailExists)
			case model.ErrPhoneInvalid:
				rep.Errors = append(rep.Errors, ErrPhoneInvalid)
			case model.ErrPhoneExists:
				rep.Errors = append(rep.Errors, ErrPhoneExists)
			case model.ErrPassInvalid:
				rep.Errors = append(rep.Errors, ErrPassInvalid)
			}

		}
		if len(rep.Errors) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.ERROR, regAPIErr, cerrs[0].Error())
			return
		}

	}

	// send results
	rep.Citizen = c
	MarshalJSON(w, rep, regAPIErr)
}

// Register GET+POST
func Register(w http.ResponseWriter, r *http.Request) {
	// return webpage if GET
	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"csrfField": csrf.TemplateField(r),
		}

		Render(w, r, data, ViewRegister, "register.tmpl")
		return
	}

	// handle form
	log.Printf("%sregister POST", config.INFO)
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}

	su := req{
		r.FormValue("id"),
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("phone"),
		r.FormValue("password"),
		r.FormValue("first_name"),
		r.FormValue("family_name"),
	}

	errs := make(map[string]interface{})

	// dummy check for pass verification
	if r.FormValue("password") != r.FormValue("repassword") {
		errs["RePassErr"] = Out[ErrRePassInvalid]
	}

	c, cerrs := model.NewCitizen(su.ID, su.Username, su.Email, su.Pass, su.Phone, su.FirstName, su.FamilyName)
	if len(cerrs) > 0 {
		for _, err := range cerrs {
			uerr := false
			switch err {
			case model.ErrUserIDInvalid:
				errs["IDErr"] = Out[ErrUserIDInvalid]
			case model.ErrUserIDExists:
				errs["IDErr"] = Out[ErrUserIDExists]
			case model.ErrUsernameInvalid:
				errs["UsernameErr"] = Out[ErrUsernameInvalid]
			case model.ErrUsernameExists:
				errs["UsernameErr"] = Out[ErrUsernameExists]
			case model.ErrEmailInvalid:
				errs["EmailErr"] = Out[ErrEmailInvalid]
			case model.ErrEmailExists:
				errs["EmailErr"] = Out[ErrEmailExists]
			case model.ErrPhoneInvalid:
				errs["PhoneErr"] = Out[ErrPhoneInvalid]
			case model.ErrPhoneExists:
				errs["PhoneErr"] = Out[ErrPhoneExists]
			case model.ErrPassInvalid:
				errs["PassErr"] = Out[ErrPassInvalid]
			case model.ErrFirstNameInvalid:
				errs["FirstNameErr"] = Out[ErrFirstNameInvalid]
			case model.ErrFamilyNameInvalid:
				errs["FamilyNameErr"] = Out[ErrFamilyNameInvalid]
			default:
				uerr = true
			}
			if uerr {
				log.Printf("%s%s: %s", config.ERROR, regErr, cerrs[0].Error())
				http.Redirect(w, r, "/error", http.StatusInternalServerError)
				return
			}
		}

	} else {

		// successfully registered the user
		http.Redirect(w, r, "/", http.StatusFound)

		// hardcoded, change it later
		link := "Confirm your account\nhttp://localhost:8080/confirm/"
		link += strconv.FormatUint(c.ID, 10)
		t, err := c.Confirmed()
		if err != nil {
			if err == model.ErrUserNotConfirmed {
				link = link + "/" + t

				// send confirmation email on a seperate
				// goroutine
				go func(email, link string) {
					m := config.NewMail(email, "Madina-TIC account confirmation", link)
					err := m.Send()
					if err != nil {
						log.Printf("%s%s: %s", config.ERROR, regErr, err.Error())
					}
				}(c.Email, link)
			}
		}
		return
	}

	// report back errors
	errs["ID"] = r.FormValue("id")
	errs["Username"] = r.FormValue("username")
	errs["Email"] = r.FormValue("email")
	errs["Phone"] = r.FormValue("phone")
	errs["FirstName"] = r.FormValue("first_name")
	errs["FamilyName"] = r.FormValue("family_name")

	errs["csrfField"] = csrf.TemplateField(r)
	Render(w, r, errs, ViewRegister, "register.tmpl")
}

// Confirm user's confirmation
func Confirm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u, err := model.UserByID(vars["id"])
	if err != nil {
		if err != model.ErrUserDoesNotExist {
			log.Printf("%s%s: %s", config.ERROR, regErr, err.Error())
		}
		http.Redirect(w, r, "/error", http.StatusBadRequest)
		return

	}
	err = u.Confirm(vars["token"])
	if err != nil {
		if err != model.ErrTokenInvalid {
			log.Printf("%s%s: %s", config.ERROR, regErr, err.Error())
		}
		http.Redirect(w, r, "/error", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Your account was successfully activated."))
}
