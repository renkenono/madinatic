package model

import (
	"io"
	"os"
	"path"

	"github.com/renkenn/madinatic/config"
)

// Pic struct
type Pic struct {
	Name string
	RID  uint64
}

// NewPic stores pic
func NewPic(rid uint64, f io.Reader, fn string) (*Pic, error) {
	fn1 := path.Join("web", "public", "storage", fn)
	out, err := os.Create(fn1)
	if err != nil {
		return nil, err
	}

	// this is chink and will obviously overwrite pic
	// if filename already exists thus data loss
	// who cares anyway
	defer out.Close()
	_, err = io.Copy(out, f)
	if err != nil {
		return nil, err
	}

	// store in db
	p := new(Pic)
	p.RID = rid
	p.Name = path.Join("/static", "storage", fn)

	// Insert Pic
	// name could be SQL injection so yeah
	// make sure to properly check for that
	// TODO: if error 1062: duplicate entry
	// then add `_` at the end of the name
	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("INSERT INTO pictures (pk_name, fk_reportid) values(?,?)")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(p.Name, p.RID)
	return nil, err
}

// PicsByReport .
func PicsByReport(id uint64) ([]string, error) {
	var ps []string
	config.DB.Lock()
	defer config.DB.Unlock()
	rows, err := config.DB.Query("SELECT pk_name FROM pictures WHERE fk_reportid = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var s string
	for rows.Next() {
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		ps = append(ps, s)
	}

	err = rows.Err()
	return ps, err
}
