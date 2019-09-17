package model

import (
	"crypto/md5"
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
	ModifiedAt time.Time `json:"modified_at"`
}

// User model custom errors
var (
	ErrUserIDInvalid    = errors.New("userid is invalid")
	ErrUserIDExists     = errors.New("userid already exists")
	ErrUsernameInvalid  = errors.New("username is invalid")
	ErrUsernameExists   = errors.New("username already exists")
	ErrEmailInvalid     = errors.New("email is invalid")
	ErrEmailExists      = errors.New("email already exists")
	ErrPhoneInvalid     = errors.New("phone is invalid")
	ErrPhoneExists      = errors.New("phone already exists")
	ErrPassInvalid      = errors.New("password is invalid")
	ErrUserDoesNotExist = errors.New("user does not exist")
	ErrUserNotConfirmed = errors.New("user account is not confirmed")
	ErrTokenInvalid     = errors.New("token is invalid")
)

// NewUser creates a new user after validating data
// assuming that id is validating by calling type
// MySQL does NOT return multi value errors when inserting
// if two users post the exact same data
// INSERT will only return error 1062 duplicate values key userid
// assuming that it is unlikely for a user to insert the same data
// while processing the other's data
func NewUser(id, username, email, pass, phone string, confirmed bool) (*User, []error) {

	// validate data
	var errs []error
	var err error
	u := new(User)
	u.ID, err = ExistsUserID(id)
	if err != nil {
		errs = append(errs, err)
	}

	u.Username, err = ExistsUsername(username)
	if err != nil {
		errs = append(errs, err)
	}

	u.Email, err = ExistsEmail(email)
	if err != nil {
		errs = append(errs, err)
	}

	u.Phone, err = ExistsPhone(phone)
	if err != nil {
		errs = append(errs, err)
	}

	_, err = ValidatePass(pass)
	if err != nil {
		errs = append(errs, err)
	}

	// store the hashed pass
	u.pass, err = HashPassword(pass)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	// create token
	u.CreatedAt = time.Now()
	u.ModifiedAt = u.CreatedAt
	token := ""
	if !confirmed {
		token = fmt.Sprintf("%x", md5.Sum([]byte(u.CreatedAt.String())))
	}

	passToken := fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))

	// Insert User
	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("INSERT INTO users (pk_userid, username, email, password, phone, confirm_token, reset_token, created_at, modified_at) values(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return nil, []error{err}
	}

	_, err = stmt.Exec(u.ID, u.Username, u.Email, u.pass, u.Phone, token, passToken, u.CreatedAt, u.ModifiedAt)
	if err != nil {
		return nil, []error{err}
	}

	return u, nil

}

// Login checks user's credentials
func Login(username, pass string) (*User, error) {
	u, err := UserByUsername(username)
	if err != nil {
		return nil, err
	}

	err = CheckPasswordHash(pass, u.pass)
	return u, err
}

// Cname returns full name
func (u *User) Cname() (string, error) {
	err := u.IsCitizen()
	if err != nil {
		if err == ErrNotCitizen {
			// Auth
			a, err := AuthByUsername(u.Username)
			if err != nil {
				return "", err
			}
			return a.Name, nil
		}
		// unexpected error
		return "", err
	}
	// citizen
	c, err := CitizenByUsername(u.Username)
	if err != nil {
		return "", err
	}

	return c.FamilyName + " " + c.FamilyName, nil

}

// Pass returns password
func (u *User) Pass() string {
	return u.pass
}

// ValidateUserID returns an error if following rules are not met
// id must be valid uint64
// always check for error before using nid
func ValidateUserID(id string) (uint64, error) {
	nid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || nid > 999999999999999999 {
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
	ok, err := regexp.MatchString(`213[5|6|7][0-9]{8}$`, phone)
	if err != nil || !ok {
		return 0, ErrPhoneInvalid
	}

	nphone, err := strconv.ParseUint(phone, 10, 64)
	if err != nil {
		return 0, ErrPhoneInvalid
	}
	return nphone, nil
}

// ExistsUserID checks if email already exists
func ExistsUserID(id string) (uint64, error) {
	nid, err := ValidateUserID(id)
	if err != nil {
		return 0, err
	}

	var qid uint64
	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT pk_userid FROM users WHERE pk_userid = ?", nid).Scan(&qid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nid, nil
		}
		return 0, err
	}

	return 0, ErrUserIDExists
}

// ExistsUsername checks if username already exists
func ExistsUsername(name string) (string, error) {
	sn, err := ValidateUsername(name)
	if err != nil {
		return "", err
	}
	var rid uint64
	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT pk_userid FROM users WHERE username = ?", name).Scan(&rid)
	if err != nil {
		if err == sql.ErrNoRows {
			return sn, nil
		}
		return "", err
	}
	return "", ErrUsernameExists
}

// ExistsPhone checks if phone already exists
func ExistsPhone(phone string) (uint64, error) {
	nphone, err := ValidatePhone(phone)
	if err != nil {
		return 0, err
	}

	var rid uint64
	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT pk_userid FROM users WHERE phone = ?", nphone).Scan(&rid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nphone, nil
		}
		return 0, err
	}
	return 0, ErrPhoneExists
}

// ExistsEmail checks if email already exists
func ExistsEmail(email string) (string, error) {
	_, err := ValidateEmail(email)
	if err != nil {
		return "", err
	}
	var rid uint64
	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT pk_userid FROM users WHERE email = ?", email).Scan(&rid)
	if err != nil {
		if err == sql.ErrNoRows {
			return email, nil
		}
		return "", err
	}
	return "", ErrEmailExists
}

// UserByID returns user based on given ID(pk_userid)
func UserByID(id string) (*User, error) {
	nid, err := ValidateUserID(id)
	if err != nil {
		return nil, err
	}

	u := new(User)
	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT pk_userid, username, email, phone, password, created_at, modified_at FROM users WHERE pk_userid = ?", nid).Scan(&u.ID, &u.Username, &u.Email, &u.Phone, &u.pass, &u.CreatedAt, &u.ModifiedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserDoesNotExist
		}
		return nil, err
	}
	return u, nil
}

// UserByIDi uint64 param
func UserByIDi(id uint64) (*User, error) {
	u := new(User)
	config.DB.Lock()
	defer config.DB.Unlock()
	err := config.DB.QueryRow("SELECT pk_userid, username, email, phone, password, created_at, modified_at FROM users WHERE pk_userid = ?", id).Scan(&u.ID, &u.Username, &u.Email, &u.Phone, &u.pass, &u.CreatedAt, &u.ModifiedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserDoesNotExist
		}
		return nil, err
	}
	return u, nil
}

// UserByEmail returns user based on given email
func UserByEmail(email string) (*User, error) {
	_, err := ValidateEmail(email)
	if err != nil {
		return nil, err
	}

	u := new(User)
	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT pk_userid, username, email, phone, password, created_at, modified_at FROM users WHERE email = ?", email).Scan(&u.ID, &u.Username, &u.Email, &u.Phone, &u.pass, &u.CreatedAt, &u.ModifiedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserDoesNotExist
		}
		return nil, err
	}
	return u, nil
}

// UserByUsername returns user based on given username
func UserByUsername(username string) (*User, error) {
	_, err := ValidateUsername(username)
	if err != nil {
		return nil, err
	}

	u := new(User)
	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT pk_userid, username, email, phone, password, created_at, modified_at FROM users WHERE username = ?", username).Scan(&u.ID, &u.Username, &u.Email, &u.Phone, &u.pass, &u.CreatedAt, &u.ModifiedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserDoesNotExist
		}
		return nil, err
	}
	return u, nil
}

// EditEmail edits user's email after validation
func (u *User) EditEmail(email string) error {
	var err error
	_, err = ExistsEmail(email)
	if err != nil {
		return err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("UPDATE users SET email = ? WHERE pk_userid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(email, u.ID)
	if err != nil {
		return err
	}
	u.Email = email
	return nil
}

// EditPhone edits user's phone after validation
func (u *User) EditPhone(phone string) error {
	np, err := ExistsPhone(phone)
	if err != nil {
		return err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("UPDATE users SET phone = ? WHERE pk_userid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(np, u.ID)
	if err != nil {
		return err
	}
	u.Phone = np
	return nil
}

// ResetToken returns the reset token
func (u *User) ResetToken() (string, error) {
	var token string
	config.DB.Lock()
	defer config.DB.Unlock()
	err := config.DB.QueryRow("SELECT reset_token FROM users WHERE pk_userid = ?", u.ID).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserDoesNotExist
		}
		return "", err
	}
	return token, nil

}

// EditPass edits user's password after validation
func (u *User) EditPass(pass string) error {
	_, err := ValidatePass(pass)
	if err != nil {
		return err
	}

	// store the hashed pass
	h, err := HashPassword(pass)
	if err != nil {
		return err
	}
	passToken := fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))

	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("UPDATE users SET password = ?, reset_token = ? WHERE pk_userid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(h, passToken, u.ID)
	if err != nil {
		return err
	}
	u.pass = h
	return nil
}

// Confirmed returns nil if confirmed
// ErrUserNotConfirmed if not
// else error
func (u *User) Confirmed() (string, error) {
	config.DB.Lock()
	var token string
	defer config.DB.Unlock()
	err := config.DB.QueryRow("SELECT confirm_token FROM users WHERE pk_userid = ?", u.ID).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserDoesNotExist
		}
		return "", err
	}
	if token != "" {
		return token, ErrUserNotConfirmed
	}
	return "", nil
}

// Confirm user
func (u *User) Confirm(t string) error {
	config.DB.Lock()
	defer config.DB.Unlock()
	var token string
	err := config.DB.QueryRow("SELECT confirm_token FROM users WHERE pk_userid = ?", u.ID).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserDoesNotExist
		}
		return err
	}

	if t == token {
		stmt, err := config.DB.Prepare("UPDATE users SET confirm_token = \"\" WHERE pk_userid = ?")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(u.ID)
	} else {
		err = ErrTokenInvalid
	}
	return err
}

// Delete a user
func (u *User) Delete() error {

	rs, err := ReportsByUser(u.Username)
	if err != nil {
		return err
	}

	for _, r := range rs {
		// admin id 1
		err := r.EditUser(1)
		if err != nil {
			return err
		}
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("DELETE FROM users WHERE pk_userid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(u.ID)
	if err != nil {
		return err
	}
	return nil
}

// IsCitizen or not
func (u *User) IsCitizen() error {
	config.DB.Lock()
	defer config.DB.Unlock()
	err := config.DB.QueryRow("SELECT pk_userid FROM citizens WHERE pk_userid = ?", u.ID).Scan(&u.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotCitizen
		}
	}
	return nil

}

// IsAdmin or not
func (u *User) IsAdmin() bool {
	return u.ID == AdminID
}

// Users return list of users
func Users() ([]*User, error) {
	var users []*User
	config.DB.Lock()
	defer config.DB.Unlock()
	rows, err := config.DB.Query("SELECT pk_userid, username, email, phone, password, created_at, modified_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u User
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Phone, &u.pass, &u.CreatedAt, &u.ModifiedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &User{
			ID:         u.ID,
			Username:   u.Username,
			Email:      u.Email,
			Phone:      u.Phone,
			pass:       u.pass,
			CreatedAt:  u.CreatedAt,
			ModifiedAt: u.ModifiedAt,
		})
	}

	err = rows.Err()
	return users, err
}
