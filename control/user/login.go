package user

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	// "html/template"
	"log"
	"net/http"
	// "path"

	"github.com/dgrijalva/jwt-go"
	"github.com/renkenn/madinatic/config"
	"github.com/renkenn/madinatic/model"
)

// errors and expire time
const (
	ExpireTime          = 5097600
	ErrUserDoesNotExist = 1
	ErrWrongPassword    = iota
	ErrUserNotConfirmed = iota
	ErrOther            = iota
)

type Credentials struct {
	Username string `json:"username"`
	Pass     string `json:"password"`
}

// LoginResp is what the client requesting LoginAPI
// will recieve, it holds an access token and
// an error number that refers to the following
// 0 - no error
// 1 - ErrUserDoesNotExist
// 2 - ErrUserNotConfirmed
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
	r.ParseForm()
	var cred Credentials
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%slogin API failed: %s", config.ERROR, err.Error())
		return
	}

	var resp LoginResp
	var token string
	// validate credentials
	u, err := model.Login(cred.Username, cred.Pass)
	if err != nil {
		if err != model.ErrUserDoesNotExist && err != bcrypt.ErrMismatchedHashAndPassword {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%slogin API failed: %s", config.ERROR, err.Error())
			return
		}
		if err == model.ErrUserDoesNotExist {
			resp.Error = ErrUserDoesNotExist
		} else {
			resp.Error = ErrWrongPassword
		}
		goto end
	}

	err = u.Confirmed()
	if err != nil {
		if err != model.ErrUserNotConfirmed {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%slogin API failed: %s", config.ERROR, err.Error())
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
		log.Printf("%slogin API failed: %s", config.ERROR, err.Error())
		return
	}
	resp.AccessToken = token

end:
	// marshal response
	respjson, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%slogin API failed: %s", config.ERROR, err.Error())
		return
	}

	// send the token back
	w.WriteHeader(http.StatusOK)
	w.Write(respjson)
}

// Login handles HTML login Forms
func Login(w http.ResponseWriter, r *http.Request) {
}

// newAccessToken returns a JWT valid token made for user given
// it expires after 1 year (might change)
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
