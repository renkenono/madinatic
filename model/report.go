package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/renkenn/madinatic/config"
)

// Report struct
type Report struct {
	ID              uint64
	Title           string
	Desc            string
	CreatedAt       time.Time
	ModifiedAt      time.Time
	Latitude        float64
	Longtitude      float64
	Address         string
	State           uint8
	StateModifiedAt time.Time
	Username        string
	UID             uint64
}

// Report state values
const (
	ReportPending  = 0
	ReportRejected = iota
	ReportAccepted = iota
	ReportSolved   = iota
)

// Report specific errors
var (
	ErrTitleInvalid       = errors.New("title is invalid")
	ErrDescInvalid        = errors.New("desc is invalid")
	ErrLatInvalid         = errors.New("latitude is invalid")
	ErrLongInvalid        = errors.New("longtitude is invalid")
	ErrOpStateInvalid     = errors.New("current state does not allow such an operation")
	ErrReportDoesNotExist = errors.New("report does not exist")
)

// NewReport returns a report duhhhhhahdisaudhasuidhuasidhai
func NewReport(username, title, desc, addr, lat, long string) (*Report, error) {
	r := new(Report)
	r.Title = title
	r.Desc = desc
	r.Address = addr
	f, err := strconv.ParseFloat(strings.TrimSpace(lat), 64)
	if err != nil {
		return nil, ErrLatInvalid
	}
	r.Latitude = f

	f, err = strconv.ParseFloat(strings.TrimSpace(long), 64)
	if err != nil {
		return nil, ErrLongInvalid
	}
	r.Longtitude = f

	r.CreatedAt = time.Now()
	r.ModifiedAt = r.CreatedAt
	r.StateModifiedAt = r.ModifiedAt
	r.State = 0

	u, err := UserByUsername(username)
	if err != nil {
		return nil, err
	}
	r.Username = username

	// Insert Report
	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("INSERT INTO reports (title, descr, created_at, modified_at, latitude, longitude, addr, curr_state, state_modified_at, fk_userid) values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(r.Title, r.Desc, r.CreatedAt, r.ModifiedAt, r.Latitude, r.Longtitude, r.Address, r.State, r.StateModifiedAt, u.ID)
	if err != nil {
		return nil, err
	}
	// get the last record

	err = config.DB.QueryRow("SELECT pk_reportid FROM reports ORDER BY pk_reportid DESC LIMIT 1;").Scan(&r.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrReportDoesNotExist
		}
		return nil, err
	}

	return r, nil
}

// ReportsByUser returns reports made by user (useless I think)
func ReportsByUser(username string) ([]*Report, error) {
	var reports []*Report
	u, err := UserByUsername(username)
	if err != nil {
		return nil, err
	}
	config.DB.Lock()
	defer config.DB.Unlock()
	rows, err := config.DB.Query("SELECT pk_reportid, title, descr, created_at, modified_at, latitude, longitude, addr, curr_state, state_modified_at FROM reports WHERE fk_userid = ?", u.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var r Report
	for rows.Next() {
		err := rows.Scan(&r.ID, &r.Title, &r.Desc, &r.CreatedAt, &r.ModifiedAt, &r.Latitude, &r.Longtitude, &r.Address, &r.State, &r.StateModifiedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &Report{
			ID:              r.ID,
			Title:           r.Title,
			Desc:            r.Desc,
			CreatedAt:       r.CreatedAt,
			ModifiedAt:      r.ModifiedAt,
			Latitude:        r.Latitude,
			Longtitude:      r.Longtitude,
			Address:         r.Address,
			State:           r.State,
			StateModifiedAt: r.StateModifiedAt,
			Username:        username,
			UID:             r.UID,
		})
	}

	err = rows.Err()
	return reports, err
}

// ReportsByState returns reports made by user
// make sure to get the username of u and update report
// can't be performed here due to double mutex lock
func ReportsByState(s uint8) ([]*Report, error) {
	var reports []*Report

	config.DB.Lock()
	defer config.DB.Unlock()

	rows, err := config.DB.Query("SELECT pk_reportid, title, descr, created_at, modified_at, latitude, longitude, addr, curr_state, state_modified_at, fk_userid FROM reports WHERE curr_state = ?", s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var r Report
	for rows.Next() {
		err := rows.Scan(&r.ID, &r.Title, &r.Desc, &r.CreatedAt, &r.ModifiedAt, &r.Latitude, &r.Longtitude, &r.Address, &r.State, &r.StateModifiedAt, &r.UID)
		if err != nil {
			return nil, err
		}
		fmt.Println(r)

		reports = append(reports, &Report{
			ID:              r.ID,
			Title:           r.Title,
			Desc:            r.Desc,
			CreatedAt:       r.CreatedAt,
			ModifiedAt:      r.ModifiedAt,
			Latitude:        r.Latitude,
			Longtitude:      r.Longtitude,
			Address:         r.Address,
			State:           r.State,
			StateModifiedAt: r.StateModifiedAt,
			UID:             r.UID,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	println("something ", reports)
	return reports, nil
}

// SetState sets a state of a report, a bit chink
func (r *Report) SetState(state uint8) error {
	if r.State > state {
		return ErrOpStateInvalid
	}
	stmt, err := config.DB.Prepare("UPDATE reports SET curr_state = ? WHERE pk_reportid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(state, r.ID)
	if err != nil {
		return nil
	}
	r.State = state
	return nil
}

// Delete a report
func (r *Report) Delete() error {
	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("DELETE FROM reports WHERE pk_reportid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r.ID)
	if err != nil {
		return err
	}
	return nil
}

// EditTitle sdasda
func (r *Report) EditTitle(title string) error {
	if r.State != ReportPending {
		return ErrOpStateInvalid
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	now := time.Now()
	stmt, err := config.DB.Prepare("UPDATE reports SET title = ?, modified_at = ? WHERE pk_reportid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(title, now, r.ID)
	if err != nil {
		return err
	}

	r.Title = title
	r.ModifiedAt = now
	return nil
}

// EditDesc sdasda
func (r *Report) EditDesc(desc string) error {
	if r.State != ReportPending {
		return ErrOpStateInvalid
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	now := time.Now()
	stmt, err := config.DB.Prepare("UPDATE reports SET descr = ?, modified_at = ? WHERE pk_reportid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(desc, now, r.ID)
	if err != nil {
		return err
	}

	r.Desc = desc
	r.ModifiedAt = now
	return nil
}
