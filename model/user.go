package model

import (
	"errors"
	"net/mail"
	"regexp"
	"time"
	"unicode"
)

// User reprepsents a basic user model
// password is hashed
type User struct {
	ID         uint64 `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	pass       string
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAT time.Time `json:"modified_at"`
}

var (
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
	err := ValidateUserData(username, email, pass, phone)
	if err != nil {
		return err
	}

	return nil
}

// ValidateUserData validates given data following the given rules
// although there must be RFC-based verification
// or at least strong validation, I'll keep it simple
// thus not for production
func ValidateUserData(username, email, pass, phone string) error {
	err := ValidateUsername(username)
	if err != nil {
		return err
	}

	err = ValidateEmail(email)
	if err != nil {
		return err
	}

	err = ValidatePass(pass)
	if err != nil {
		return err
	}

	err = ValidatePhone(phone)
	if err != nil {
		return err
	}
	return nil
}

// ValidateUsername returns an error if following rules are not met
// username must be of length > 4 and < 30
func ValidateUsername(username string) error {
	n := len(username)
	if n <= 4 || n >= 30 {
		return ErrUsernameInvalid
	}

	// match with regexp
	ok, err := regexp.MatchString(`^[A-Za-z]+([ _-]?[A-Za-z0-9]+)*$`, username)
	if err != nil || !ok {
		return ErrUsernameInvalid
	}

	return nil
}

// ValidateEmail returns an error if following rules are not met
// email must be valid RFC 5322
func ValidateEmail(email string) error {

	// name could be empty
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrEmailInvalid
	}
	return nil
}

// ValidatePass returns an error if following rules are not met
// password must be of length > 8
func ValidatePass(pass string) error {
	if len(pass) < 8 {
		return ErrPassInvalid
	}
	for _, c := range pass {
		if unicode.IsControl(c) {
			return ErrPassInvalid
		}
	}
	return nil
}

// ValidatePhone returns an error if following rules are not met
// phone must be of length 12
func ValidatePhone(phone string) error {

	// hardcoded values but who cares
	// TODO: need to add FIX area codes
	ok, err := regexp.MatchString(`[\+]?213[5|6|7][0-9]{6}`, phone)
	if err != nil || !ok {
		return ErrPhoneInvalid
	}
	return nil
}

// ExistsUserID checks if email already exists
func ExistsUserID(email string) error {
	return nil
}

// ExistsUsername checks if username already exists
func ExistsUsername(email string) error {
	return nil
}

// ExistsPhone checks if phone already exists
func ExistsPhone(email string) error {
	return nil
}

// ExistsEmail checks if email already exists
func ExistsEmail(email string) error {
	return nil
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
