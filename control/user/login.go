package user

import (
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/renkenn/madinatic/config"
	. "github.com/renkenn/madinatic/control"
	"github.com/renkenn/madinatic/model"
)

const (
	logAPIErr = "login API failed"
	logErr    = "login failed"
	resetErr  = "reset password failed"
)

type credentials struct {
	Username string `json:"username"`
	Pass     string `json:"password"`
}

// LoginResp is what the client requesting LoginAPI
// will recieve, it holds an access token and
// an error number that refers to the following
// 0 - no error
// 1 - ErrUserDoesNotExist
// 2 - ErrWrongPassword
// 3 - ErrUserNotConfirmed
type LoginResp struct {
	Error       uint   `json:"error"`
	AccessToken string `json:"access_token"`
}

// LoginAPI allows a user to enter their credentials
// and login to the webservice
// they *must* provide a username and a password
// the response will contain an error code and
// one JWT bearer (access token)
// one *must* check the HTTP status code before doing
// any processing if it's different to StatusOK
// check error which value is of the following
// 0 - no error
// 1 - ErrUserDoesNotExist
// 2 - ErrWrongPassword
// 3 - ErrUserNotConfirmed
// access token is valid for a limited time -3 months-
// once the access token expires, the user is forced to
// login again in order to generate a new token
func LoginAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var cred credentials
	err := DecodeJSON(w, r, &cred, logAPIErr)
	if err != nil {
		return
	}

	var resp LoginResp
	var token string
	// validate credentials
	u, err := model.Login(cred.Username, cred.Pass)
	if err != nil {
		if err != model.ErrUserDoesNotExist && err != bcrypt.ErrMismatchedHashAndPassword && err != model.ErrUsernameInvalid {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.ERROR, logAPIErr, err.Error())
			return
		}
		if err == model.ErrUserDoesNotExist || err == model.ErrUsernameInvalid {
			resp.Error = ErrUserDoesNotExist
		} else {
			resp.Error = ErrWrongPassword
		}
		goto end
	}

	_, err = u.Confirmed()
	if err != nil {
		if err != model.ErrUserNotConfirmed {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.ERROR, logAPIErr, err.Error())
			return
		}
		resp.Error = ErrUserNotConfirmed
		goto end
	}

	// user exists and is confirmed
	// create tokens
	token, err = newAccessToken(u.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.ERROR, logAPIErr, err.Error())
		return
	}
	resp.AccessToken = token

end:
	MarshalJSON(w, resp, logAPIErr)
}

// Login handles HTML login Forms
func Login(w http.ResponseWriter, r *http.Request) {

	s, err := config.Store.Get(r, "userdata")
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, logErr, err.Error())
		return
	}
	if !s.IsNew && s.Values["username"] != "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// return webpage if GET
	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"csrfField": csrf.TemplateField(r),
		}

		Render(w, r, data, ViewLogin, "login.tmpl")
		return
	}

	// handle form
	log.Printf("%slogin: POST Request", config.INFO)
	err = r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, logErr, err.Error())
		return
	}

	cred := credentials{r.FormValue("username"), r.FormValue("password")}
	var errstr string
	u, err := model.Login(cred.Username, cred.Pass)
	if err != nil {
		if err != model.ErrUserDoesNotExist && err != bcrypt.ErrMismatchedHashAndPassword && err != model.ErrUsernameInvalid {
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.INFO, logErr, err.Error())
			return
		} else if err == model.ErrUserDoesNotExist || err == model.ErrUsernameInvalid {
			errstr = Out[ErrUserDoesNotExist]
		} else {
			errstr = Out[ErrWrongPassword]
		}
		goto logerr
	}

	_, err = u.Confirmed()
	if err != nil {
		if err != model.ErrUserNotConfirmed {
			log.Printf("%s%s: %s", config.INFO, logErr, err.Error())
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			return
		}
		errstr = Out[ErrUserNotConfirmed]

		goto logerr
	}

	s.Values["username"] = u.Username
	s.AddFlash("Logged in successfully.")
	err = s.Save(r, w)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, logErr, err.Error())
		return
	}

	// login successful, redirect to "/"
	http.Redirect(w, r, "/", http.StatusFound)
	return
logerr:
	// errstr != ""
	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"error":     errstr,
	}

	Render(w, r, data, ViewLogin, "login.tmpl")
}

// ResetPass of user
func ResetPass(w http.ResponseWriter, r *http.Request) {
	// validate if id exists
	vars := mux.Vars(r)
	u, err := model.UserByID(vars["id"])
	if err != nil {
		return
	}

	// validate if reset_token is correct
	token, err := u.ResetToken()
	if err != nil {
		log.Printf("%s%s: %s", config.ERROR, logErr, err.Error())
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}
	if token != vars["token"] {
		http.Redirect(w, r, "/error", http.StatusBadRequest)
		return
	}

	// valid url
	// return webpage
	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
		"Username":  u.Username,
		"URL":       r.URL.Path,
	}

	if r.Method == http.MethodGet {
		Render(w, r, data, ViewReset, "reset.tmpl")
		return
	}

	// process change
	err = r.ParseForm()
	if err != nil {
		log.Printf("%s%s: %s", config.ERROR, logErr, err.Error())
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}

	err = u.EditPass(r.FormValue("password"))
	if err != nil {
		if err == model.ErrPassInvalid {
			// return view with invalid password error
			data["error"] = Out[ErrPassInvalid]

			Render(w, r, data, ViewReset, "reset.tmpl")
		} else {
			log.Printf("%s%s: %s", config.ERROR, logErr, err.Error())
			http.Redirect(w, r, "/error", http.StatusBadRequest)

		}
		return
	}

	w.Write([]byte("Password successfully updated."))
}

func Reset(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("%s%s: %s", config.ERROR, resetErr, err.Error())
		return
	}
	u, err := model.UserByEmail(r.FormValue("email"))
	if err != nil && err != model.ErrUserDoesNotExist {
		log.Printf("%s%s: %s", config.ERROR, resetErr, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	go func(u *model.User) {
		link := "Reset Your Account's Password\nhttp://localhost:8080/reset/"
		link += strconv.FormatUint(u.ID, 10)
		t, err := u.ResetToken()
		if err != nil {
			log.Printf("%s%s: %s", config.ERROR, resetErr, err.Error())
			return
		}
		link = link + "/" + t
		m := config.NewMail(u.Email, "Madina-TIC reset account password", link)
		err = m.Send()
		if err != nil {
			log.Printf("%s%s: %s", config.ERROR, resetErr, err.Error())
		}

	}(u)
}

// newAccessToken returns a JWT valid token made for user given
// it expires after a period of time (might change)
// you get to regenerated it by simply logging in
func newAccessToken(username string) (ss string, err error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: ExpireTime,
		Issuer:    username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err = token.SignedString([]byte(config.App.SignKey))
	return ss, err
}
