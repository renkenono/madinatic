package user

import (
	"log"
	"net/http"
	"strconv"

	"github.com/renkenn/madinatic/model"

	"github.com/gorilla/csrf"

	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
)

const (
	setAPIErr = "settings API failed"
	setErr    = "settings failed"
)

type respset struct {
	Errors []uint `json:"errors"`
}

type reqset struct {
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Pass        string `json:"password"`
	OldPass     string `json:"old_password"`
}

// SettingsAPI returns an array of errors
// empty array if no errors
// else array containing numbers corresponding to the errors
// chinked token :DDDDDDDDDdd
func SettingsAPI(w http.ResponseWriter, r *http.Request) {
	var re reqset
	err := DecodeJSON(w, r, &re, setAPIErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.ERROR, setAPIErr, err.Error())
		return
	}

	u, err := model.UserByUsername(re.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.ERROR, setAPIErr, err.Error())
		return
	}

	var rep respset

	err = model.CheckPasswordHash(re.OldPass, u.Pass())
	if err != nil {
		rep.Errors = append(rep.Errors, ErrWrongPassword)
		MarshalJSON(w, rep, setAPIErr)
		return
	}

	err = u.EditPass(re.Pass)
	if err != nil {
		rep.Errors = append(rep.Errors, ErrPassInvalid)
		MarshalJSON(w, rep, setAPIErr)
		return
	}

	if u.Email != re.Email {
		err = u.EditEmail(re.Email)
		if err != nil {
			if err == model.ErrEmailExists {
				rep.Errors = append(rep.Errors, ErrEmailExists)
			} else if err == model.ErrEmailInvalid {
				rep.Errors = append(rep.Errors, ErrEmailInvalid)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("%s%s: %s", config.ERROR, setAPIErr, err.Error())
				return
			}
		}
	}

	nphone, err := strconv.ParseUint(re.Phone, 10, 64)
	if err != nil && u.Phone != nphone {
		err = u.EditPhone(re.Phone)
		if err != nil {
			if err == model.ErrPhoneExists {
				rep.Errors = append(rep.Errors, ErrPhoneExists)
			} else if err == model.ErrPhoneInvalid {
				rep.Errors = append(rep.Errors, ErrPhoneInvalid)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("%s%s: %s", config.ERROR, setAPIErr, err.Error())
				return
			}
		}

	}

	MarshalJSON(w, rep, setAPIErr)

}

func Settings(w http.ResponseWriter, r *http.Request) {
	s, err := config.Store.Get(r, "userdata")
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, setErr, err.Error())
		return
	}

	if s.IsNew {
		// login successful, redirect to "/"
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	u, err := model.UserByUsername(s.Values["username"].(string))
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Email":     u.Email,
		"Phone":     u.Phone,
		"csrfField": csrf.TemplateField(r),
	}
	if r.Method == http.MethodGet {
		Render(w, r, data, ViewSettings, "settings.tmpl")
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}

	err = model.CheckPasswordHash(r.FormValue("oldpassword"), u.Pass())
	if err != nil {
		data["OldPassErr"] = Out[ErrWrongPassword]
		Render(w, r, data, ViewSettings, "settings.tmpl")
		return
	}

	if r.FormValue("password") != r.FormValue("repassword") {
		data["RePassErr"] = Out[ErrRePassInvalid]
		Render(w, r, data, ViewSettings, "settings.tmpl")
		return
	}

	err = u.EditPass(r.FormValue("password"))
	if err != nil {
		data["PassErr"] = Out[ErrPassInvalid]
	}

	if u.Email != r.FormValue("email") {
		err = u.EditEmail(r.FormValue("email"))
		data["Email"] = r.FormValue("email")
		if err != nil {
			if err == model.ErrEmailExists {
				data["EmailErr"] = Out[ErrEmailExists]
			} else if err == model.ErrEmailInvalid {
				data["EmailErr"] = Out[ErrEmailInvalid]
			} else {
				http.Redirect(w, r, "/error", http.StatusInternalServerError)
				return
			}
		}
	}

	nphone, err := strconv.ParseUint(r.FormValue("phone"), 10, 64)
	if err != nil && u.Phone != nphone {
		data["Phone"] = r.FormValue("phone")
		err = u.EditPhone(r.FormValue("phone"))
		if err != nil {
			if err == model.ErrPhoneExists {
				data["PhoneErr"] = Out[ErrPhoneExists]
			} else if err == model.ErrPhoneInvalid {
				data["PhoneErr"] = Out[ErrPhoneInvalid]
			} else {
				http.Redirect(w, r, "/error", http.StatusInternalServerError)
				return
			}
		}

	}
	s.AddFlash("successfully updated info.")
	data["flash"] = s.Flashes()[0]
	Render(w, r, data, ViewSettings, "settings.tmpl")
}
