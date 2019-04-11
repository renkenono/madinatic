package model

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/renkenn/madinatic/config"
)

// Auth superset of a User
// might create an issue when sending reports feed
// since citizen has two names whereas auth has only one
// the easiest solution is just to concat the two with a space
type Auth struct {
	*User
	Name string `json:"name"`
}

// Auth custom errors
var (
	ErrNameInvalid = errors.New("name is invalid")
	ErrNotAuth     = errors.New("user is not an authority")
)

// NewAuth auth name is basically the same as username
func NewAuth(id, username, email, pass, phone, name string) (*Auth, error) {
	err := ValidateCitizenName(name, "")
	if err != nil {
		return nil, err
	}

	u, err := NewUser(id, username, email, pass, phone, false)
	if err != nil {
		return nil, err
	}
	a := &Auth{
		u, name,
	}

	// Insert Authority
	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("INSERT INTO authorities (pk_userid, name) values(?,?)")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(a.ID, a.Name)

	return a, err
}

// Auths returns list of authorities
func Auths() ([]*Auth, error) {
	var as []*Auth
	config.DB.Lock()
	rows, err := config.DB.Query("SELECT pk_userid, name FROM authorities")
	if err != nil {
		return nil, err
	}

	config.DB.Unlock()
	defer rows.Close()
	var a Auth
	a.User = new(User)

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.Name)
		if err != nil {
			return nil, err
		}

		u, err := UserByID(strconv.FormatUint(a.ID, 10))
		if err != nil {
			return nil, err
		}

		as = append(as, &Auth{
			u, a.Name,
		})
	}

	err = rows.Err()
	return as, err
}

// AuthByID returns Auth of given ID
func AuthByID(id string) (*Auth, error) {
	nid, err := ValidateUserID(id)
	if err != nil {
		return nil, err
	}

	a := new(Auth)
	a.User = new(User)
	a.User, err = UserByID(id)
	if err != nil {
		return nil, err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT name FROM authorities WHERE pk_userid = ?", nid).Scan(&a.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotAuth
		}
		return nil, err
	}
	return a, nil

}

// AuthByUsername returns Auth of given username
func AuthByUsername(username string) (*Auth, error) {
	_, err := ValidateUsername(username)
	if err != nil {
		return nil, err
	}

	a := new(Auth)
	a.User = new(User)
	a.User, err = UserByUsername(username)
	if err != nil {
		return nil, err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT name FROM authorities WHERE pk_userid = ?", a.ID).Scan(&a.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotAuth
		}
		return nil, err
	}
	return a, nil

}
