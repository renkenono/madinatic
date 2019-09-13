package model

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/renkenn/madinatic/config"
)

// Cat def
type Cat struct {
	ID   uint
	Name string
	Auth uint64
}

// Cat specific errors
var (
	ErrCatDoesNotExist = errors.New("cat does not exist")
)

// NewCat returns a report
func NewCat(name, authname string) (*Cat, error) {
	c := new(Cat)
	c.Name = name
	a, err := AuthByUsername(authname)
	if err != nil {
		return nil, err
	}
	c.Auth = a.ID
	// Insert Report
	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("INSERT INTO categories (cat_name, fk_userid) values(?,?)")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(c.Name, c.Auth)
	if err != nil {
		return nil, err
	}
	// get the last record

	err = config.DB.QueryRow("SELECT pk_catid FROM categories ORDER BY pk_catid DESC LIMIT 1;").Scan(&c.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrReportDoesNotExist
		}
		return nil, err
	}

	return c, nil
}

// CatsStr returns an array of cat strings
func CatsStr() ([]string, error) {
	cs, err := Cats()
	if err != nil {
		return nil, err
	}
	csstr := []string{}
	for _, c := range cs {
		csstr = append(csstr, c.Name)
	}
	return csstr, nil
}

// Cats return list of cats
func Cats() ([]*Cat, error) {
	var cats []*Cat
	config.DB.Lock()
	defer config.DB.Unlock()
	rows, err := config.DB.Query("SELECT pk_catid, cat_name, fk_userid FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var c Cat
	for rows.Next() {
		err := rows.Scan(&c.ID, &c.Name, &c.Auth)
		if err != nil {
			return nil, err
		}
		cats = append(cats, &Cat{
			ID:   c.ID,
			Name: c.Name,
			Auth: c.Auth,
		})
	}

	err = rows.Err()
	return cats, err
}

// CatByIDi uint param
func CatByIDi(id uint) (*Cat, error) {
	c := new(Cat)
	config.DB.Lock()
	defer config.DB.Unlock()
	err := config.DB.QueryRow("SELECT cat_name, fk_userid WHERE pk_catid = ?", id).Scan(c.Name, c.Auth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserDoesNotExist
		}
		return nil, err
	}
	c.ID = id
	return c, nil
}

// CatByID string param
func CatByID(id string) (*Cat, error) {
	nid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}
	c, err := CatByIDi(uint(nid))
	return c, err
}

// CatByName returns cat based on given name
func CatByName(name string) (*Cat, error) {
	c := new(Cat)
	config.DB.Lock()
	defer config.DB.Unlock()
	err := config.DB.QueryRow("SELECT pk_catid, fk_userid FROM categories WHERE cat_name = ?", name).Scan(&c.ID, &c.Auth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserDoesNotExist
		}
		return nil, err
	}
	c.Name = name
	return c, nil
}
