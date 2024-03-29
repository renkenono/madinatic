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

const (
	// AdminID serves to check if current Auth
	// is the admin or not
	AdminID = 1
)

// Auth custom errors
var (
	ErrNameInvalid = errors.New("name is invalid")
	ErrNotAuth     = errors.New("user is not an authority")
)

// NewAuth auth name is basically the same as username
func NewAuth(id, username, email, pass, phone, name string) (*Auth, []error) {
	var errs []error

	u, uerrs := NewUser(id, username, email, pass, phone, true)
	if len(uerrs) > 0 {
		errs = append(errs, uerrs...)
		return nil, errs
	}
	a := &Auth{
		u, name,
	}

	// Insert Authority
	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("INSERT INTO authorities (pk_userid, name) values(?,?)")
	if err != nil {
		return nil, []error{err}
	}

	_, err = stmt.Exec(a.ID, a.Name)

	if err != nil {
		return nil, []error{err}
	}

	return a, nil
}

// NewAuthID returns last_auth_id + 1
// finish adding new auth <------------
func NewAuthID() (uint64, error) {

	config.DB.Lock()

	defer config.DB.Unlock()
	// default id for first auth is 2
	var lid uint64 = 1
	err := config.DB.QueryRow("SELECT pk_userid FROM authorities ORDER BY pk_userid DESC LIMIT 1;").Scan(&lid)
	if err != nil && err != sql.ErrNoRows {
		return lid, err
	}

	lid++
	return lid, nil
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

// SetName a
func (a *Auth) SetName(f string) error {
	err := ValidateCitizenName(f, "")
	if err != nil {
		return err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("UPDATE authorities SET family_name = ? WHERE pk_userid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(f, a.ID)
	if err != nil {
		return err
	}
	a.Name = f
	return nil

}
