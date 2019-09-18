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

type req struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Pass     string `json:"password"`
	Name     string `json:"name"`
}

type replya struct {
	I int
	*model.Auth
	Link string
}

const (
	authCreateErr = "auth create failed"
)

// AuthCreate creates an auth
func AuthCreate(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	if r.Method == http.MethodGet {
		cats, err := model.Cats()
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.INFO, authCreateErr, err.Error())
			return
		}
		var cs []*model.Cat
		for i := 0; i < len(cats)-1; i++ {
			in := false
			for j := i + 1; j < len(cats); j++ {
				if cats[i].Name == cats[j].Name {
					in = true
					break
				}
			}
			if !in || i == len(cats)-2 {
				cs = append(cs, cats[i])
			}
		}

		data := map[string]interface{}{
			"csrfField": csrf.TemplateField(r),
			"Cats":      cs,
		}

		Render(w, r, data, ViewAuthCreate, "d_new_auth.tmpl")
		return
	}

	log.Printf("%screate auth POST", config.INFO)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}

	lid, err := model.NewAuthID()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.ERROR, authCreateErr, err.Error())
		return
	}

	su := req{
		strconv.FormatUint(lid, 10),
		r.MultipartForm.Value["username"][0],
		r.MultipartForm.Value["email"][0],
		r.MultipartForm.Value["phone"][0],
		r.MultipartForm.Value["password"][0],
		r.MultipartForm.Value["name"][0],
	}

	fmt.Println("su ", su)

	errs := make(map[string]interface{})

	// dummy check for pass verification
	// if r.FormValue("password") != r.FormValue("repassword") {
	// errs["RePassErr"] = Out[ErrRePassInvalid]
	// }

	// update su.ID
	a, cerrs := model.NewAuth(su.ID, su.Username, su.Email, su.Pass, su.Phone, su.Name)
	if len(cerrs) > 0 {
		for _, err := range cerrs {
			uerr := false
			switch err {
			case model.ErrUserIDInvalid:
				errs["IDErr"] = Out[ErrUserIDInvalid]
			case model.ErrUserIDExists:
				errs["IDErr"] = Out[ErrUserIDExists]
			case model.ErrUsernameInvalid:
				errs["UsernameErr"] = Out[ErrUserIDInvalid]
			case model.ErrUsernameExists:
				errs["UsernameErr"] = Out[ErrUsernameExists]
			case model.ErrEmailInvalid:
				errs["EmailErr"] = Out[ErrUsernameInvalid]
			case model.ErrEmailExists:
				errs["EmailErr"] = Out[ErrEmailExists]
			case model.ErrPhoneInvalid:
				errs["PhoneErr"] = Out[ErrEmailInvalid]
			case model.ErrPhoneExists:
				errs["PhoneErr"] = Out[ErrPhoneExists]
			case model.ErrPassInvalid:
				errs["PassErr"] = Out[ErrPassInvalid]
			case model.ErrNameInvalid:
				errs["NameErr"] = Out[ErrFirstNameInvalid]
			default:
				uerr = true
			}
			if uerr {
				log.Printf("%s%s: %s", config.ERROR, authCreateErr, cerrs[0].Error())
				http.Redirect(w, r, "/error", http.StatusInternalServerError)
				return
			}
		}

	} else {
		fmt.Println("created auth ", a.Username)
		// hostile
		err := linkAuthCat(r, a.Username, authCreateErr)
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			return
		}

		// successfully registered the user
		http.Redirect(w, r, "/dashboard/auths", http.StatusFound)
		return
	}

	for _, e := range cerrs {
		fmt.Println("error :", e.Error())
	}
	// report back errors
	errs["Username"] = r.MultipartForm.Value["username"][0]
	errs["Email"] = r.MultipartForm.Value["email"][0]
	errs["Phone"] = r.MultipartForm.Value["phone"][0]
	errs["Name"] = r.MultipartForm.Value["name"][0]

	errs["csrfField"] = csrf.TemplateField(r)
	cats, err := model.Cats()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, authCreateErr, err.Error())
		return
	}
	errs["Cats"] = cats
	Render(w, r, errs, ViewAuthCreate, "d_new_auth.tmpl")
}

// AuthsView returns a list of Auth
func AuthsView(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(w, r) {
		return
	}

	// us
	as, err := model.Auths()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, usersViewErr, err.Error())
		return
	}
	// check if admin/auth/client

	var replies []replya

	for i, a := range as {
		replies = append(replies, replya{
			i,
			a,
			fmt.Sprintf("/user/delete/%d", a.ID),
		})
	}

	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"Users":     replies,
	}

	Render(w, r, data, ViewDashboardAuth, "d_auth.tmpl")
}

func linkAuthCat(r *http.Request, username, errstr string) error {
	selectedCs := r.MultipartForm.Value["cat"]
	// if len(selectedCs) == 0 {
	// 	return errors.New("no cats selected")
	// }

	for _, c := range selectedCs {
		_, err := model.NewCat(c, username)
		if err != nil {
			return err
		}
	}
	return nil
}
