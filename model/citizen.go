package model

import (
	"database/sql"
	"errors"
	"regexp"
	"strconv"

	"github.com/renkenn/madinatic/config"
)

// Citizen superset of a User
type Citizen struct {
	*User
	FirstName  string `json:"first_name"`
	FamilyName string `json:"family_name"`
}

// Citizen custom errors
var (
	ErrFirstNameInvalid  = errors.New("first name is not valid")
	ErrFamilyNameInvalid = errors.New("family name is not valid")
	ErrNotCitizen        = errors.New("user is not a citizen")
)

// NewCitizen creates a new Citizen account
// Essentially it is just a user
// with first name and family name added
// a Citizen account is never confirmed by default
func NewCitizen(id, username, email, pass, phone, first, family string) (*Citizen, []error) {
	var errs []error
	err := ValidateCitizenName(first, family)
	if err != nil {
		errs = append(errs, err)
	}

	u, uerrs := NewUser(id, username, email, pass, phone, false)
	if len(uerrs) > 0 {
		errs = append(errs, uerrs...)
		return nil, errs
	}
	c := &Citizen{
		u, first, family,
	}

	// Insert User
	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("INSERT INTO citizens (pk_userid, first_name, family_name) values(?,?,?)")
	if err != nil {
		return nil, []error{err}
	}

	_, err = stmt.Exec(c.ID, c.FirstName, c.FamilyName)
	if err != nil {
		return nil, []error{err}
	}

	return c, nil
}

// ValidateCitizenName returns an error if the following rules are not met
// first length [3, 30]
// maybe proper checking should be done
func ValidateCitizenName(first, family string) error {
	n := len(first)
	if n < 3 || n > 30 {
		return ErrFirstNameInvalid
	}

	// match with regexp
	ok, err := regexp.MatchString(`^[[:alpha:]]+[[:alpha:] ]*$`, first)
	if err != nil || !ok {
		return ErrFirstNameInvalid
	}

	if family == "" {
		return nil
	}

	n = len(family)
	if n < 3 || n > 30 {
		return ErrFamilyNameInvalid
	}

	// match with regexp
	ok, err = regexp.MatchString(`^[[:alpha:]]+[[:alpha:] ]*$`, family)
	if err != nil || !ok {
		return ErrFamilyNameInvalid
	}

	return nil
}

// Citizens returns all citizens
func Citizens() ([]*Citizen, error) {
	var cs []*Citizen
	config.DB.Lock()
	rows, err := config.DB.Query("SELECT pk_userid, first_name, family_name FROM citizens")
	if err != nil {
		return nil, err
	}
	// avoid dead lock
	config.DB.Unlock()
	defer rows.Close()

	var c Citizen

	// don't forget malloc
	c.User = new(User)
	for rows.Next() {
		err := rows.Scan(&c.ID, &c.FirstName, &c.FamilyName)
		if err != nil {
			return nil, err
		}
		u, err := UserByID(strconv.FormatUint(c.ID, 10))
		if err != nil {
			return nil, err
		}
		cs = append(cs, &Citizen{
			u, c.FirstName, c.FamilyName,
		})
	}

	err = rows.Err()
	return cs, err
}

// CitizenByID returns user based on given id
func CitizenByID(id string) (*Citizen, error) {
	nid, err := ValidateUserID(id)
	if err != nil {
		return nil, err
	}

	c := new(Citizen)
	c.User = new(User)
	c.User, err = UserByID(id)
	if err != nil {
		return nil, err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT first_name, family_name FROM citizens WHERE pk_userid = ?", nid).Scan(&c.FirstName, &c.FamilyName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotCitizen
		}
		return nil, err
	}
	return c, nil
}

// CitizenByUsername returns user based on given username
func CitizenByUsername(username string) (*Citizen, error) {
	_, err := ValidateUsername(username)
	if err != nil {
		return nil, err
	}

	c := new(Citizen)
	c.User = new(User)
	c.User, err = UserByUsername(username)
	if err != nil {
		return nil, err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT first_name, family_name FROM citizens WHERE pk_userid = ?", c.ID).Scan(&c.FirstName, &c.FamilyName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotCitizen
		}
		return nil, err
	}
	return c, nil
}

// SetFirstName c
func (c *Citizen) SetFirstName(f string) error {
	err := ValidateCitizenName(f, "")
	if err != nil {
		return err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("UPDATE citizens SET first_name = ? WHERE pk_userid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(f, c.ID)
	if err != nil {
		return err
	}
	c.FirstName = f
	return nil

}

// SetFamilyName c
func (c *Citizen) SetFamilyName(f string) error {
	err := ValidateCitizenName(f, "")
	if err != nil {
		return err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("UPDATE citizens SET family_name = ? WHERE pk_userid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(f, c.ID)
	if err != nil {
		return err
	}
	c.FamilyName = f
	return nil
}
