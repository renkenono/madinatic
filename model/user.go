package model

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/renkenn/madinatic/config"
)

// User reprepsents a basic user model
// password is hashed
type User struct {
	ID         uint64 `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      uint64 `json:"phone"`
	pass       string
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAT time.Time `json:"modified_at"`
}

var (
	// ErrUserIDInvalid error
	ErrUserIDInvalid = errors.New("userid is invalid")
	// ErrUserIDExists error
	ErrUserIDExists = errors.New("userid already exists")
	// ErrUsernameInvalid error
	ErrUsernameInvalid = errors.New("username is invalid")
	// ErrPassInvalid error
	ErrPassInvalid = errors.New("password is invalid")
	// ErrPhoneInvalid error
	ErrPhoneInvalid = errors.New("phone is invalid")
	// ErrEmailInvalid error
	ErrEmailInvalid = errors.New("email is invalid")
	// ErrUsernameExists error
	ErrUsernameExists = errors.New("username already exists")
	// ErrPhoneExists error
	ErrPhoneExists = errors.New("phone already exists")
	// ErrEmailExists error
	ErrEmailExists = errors.New("email already exists")
)

// New creates a new user after validating data
// assuming that id is validating by calling type
// MySQL does NOT return multi value errors when inserting
// if two users post the exact same data
// INSERT will only return error 1062 duplicate values key userid
// assuming that it is unlikely for a user to insert the same data
// while processing the other's data
func (u *User) New(id, username, email, pass, phone string) error {
	return nil
}

// ValidateUserID returns an error if following rules are not met
// id must be valid uint64
// always check for error before using nid
func ValidateUserID(id string) (uint64, error) {
	nid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, ErrUserIDInvalid
	}
	if nid > 999999999999999999 {
		return 0, ErrUserIDInvalid
	}
	return nid, nil
}

// ValidateUsername returns an error if following rules are not met
// username must be of length > 4 and < 30
func ValidateUsername(name string) (string, error) {
	n := len(name)
	if n <= 4 || n >= 30 {
		return "", ErrUsernameInvalid
	}

	// match with regexp
	ok, err := regexp.MatchString(`^[A-Za-z]+([ _-]?[A-Za-z0-9]+)*$`, name)
	if err != nil || !ok {
		return "", ErrUsernameInvalid
	}

	return name, nil
}

// ValidateEmail returns an error if following rules are not met
// email must be valid RFC 5322
func ValidateEmail(email string) (string, error) {

	// name could be empty
	_, err := mail.ParseAddress(email)
	if err != nil {
		return "", ErrEmailInvalid
	}
	return email, nil
}

// ValidatePass returns an error if following rules are not met
// password must be of length > 8
func ValidatePass(pass string) (string, error) {
	if len(pass) < 8 {
		return "", ErrPassInvalid
	}
	for _, c := range pass {
		if unicode.IsControl(c) {
			return "", ErrPassInvalid
		}
	}
	return pass, nil
}

// ValidatePhone returns an error if following rules are not met
// phone must be of length 12
func ValidatePhone(phone string) (uint64, error) {

	// hardcoded values but who cares
	// TODO: need to add FIX area codes
	ok, err := regexp.MatchString(`213[5|6|7][0-9]{6}`, phone)
	if err != nil || !ok {
		return 0, ErrPhoneInvalid
	}

	nphone, err := strconv.ParseUint(phone, 10, 64)
	fmt.Println(nphone)
	if err != nil {
		return 0, ErrPhoneInvalid
	}
	return nphone, nil
}

// ExistsUserID checks if email already exists
func ExistsUserID(id string) error {
	nid, err := ValidateUserID(id)
	if err != nil {
		return err
	}

	var qid uint64
	err = config.DB.QueryRow("SELECT pk_userid FROM users where pk_userid = ?", nid).Scan(&qid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return ErrUserIDExists
}

// ExistsUsername checks if username already exists
func ExistsUsername(name string) error {
	_, err := ValidateUsername(name)
	if err != nil {
		return err
	}
	var rid uint64
	err = config.DB.QueryRow("SELECT pk_userid FROM users where username = ?", name).Scan(&rid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return ErrUsernameExists
}

// ExistsPhone checks if phone already exists
func ExistsPhone(phone string) error {
	nphone, err := ValidatePhone(phone)
	if err != nil {
		return err
	}

	var rid uint64
	err = config.DB.QueryRow("SELECT pk_userid FROM users where phone = ?", nphone).Scan(&rid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return ErrPhoneExists
}

// ExistsEmail checks if email already exists
func ExistsEmail(email string) error {
	var rid uint64
	err := config.DB.QueryRow("SELECT pk_userid FROM users where email = ?", email).Scan(&rid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return ErrEmailExists
}

// EditEmail edits user's email after validation
func (u *User) EditEmail(email string) error {
	return nil
}

// EditPassword edits user's password after validation
func (u *User) EditPassword(pass string) error {
	return nil
}

// EditPhone edits user's phone after validation
func (u *User) EditPhone(phone string) error {
	return nil
}

// Delete a user
func (u *User) Delete() error {
	return nil
}
