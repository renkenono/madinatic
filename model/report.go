package model

import (
	"database/sql"
	"errors"
	"io"
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
	UserCName       string
}

// SubReport struct
type SubReport struct {
	RID   uint64
	CID   uint
	State uint8
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

// NewReport returns a report
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

// NewSubReport fn
func (r *Report) NewSubReport(catname string) error {
	sb := new(SubReport)
	sb.RID = r.ID
	c, err := CatByName(catname)
	if err != nil {
		return err
	}
	sb.CID = c.ID
	// Insert Report
	config.DB.Lock()
	defer config.DB.Unlock()
	stmt, err := config.DB.Prepare("INSERT INTO subreports (fk_reportid, fk_catid, curr_state) values(?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(sb.RID, sb.CID, sb.State)
	return err
}

// Categories returns a string array of
// the report's categories
func (r *Report) Categories() ([]string, error) {
	var cats []string
	config.DB.Lock()
	defer config.DB.Unlock()
	rows, err := config.DB.Query("SELECT (cat_name) FROM subreports JOIN categories ON fk_catid=pk_catid WHERE fk_reportid=?;", r.ID)

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

		// just want to avoid slicing issues
		// which is unavoidable with this most likely
		// will see in run-time
		s1 := s
		cats = append(cats, s1)
	}

	err = rows.Err()
	return cats, err
}

// Auths retuns concerned auths
// didn't even bother changing var names
func (r *Report) Auths() ([]string, error) {
	var cats []string
	config.DB.Lock()
	defer config.DB.Unlock()
	rows, err := config.DB.Query("SELECT username FROM (SELECT pk_userid FROM (SELECT categories.fk_userid FROM subreports JOIN categories ON fk_catid=pk_catid WHERE fk_reportid=? ) as a JOIN authorities ON pk_userid=fk_userid) as b JOIN users ON b.pk_userid = users.pk_userid;", r.ID)

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

		// just want to avoid slicing issues
		// which is unavoidable with this most likely
		// will see in run-time
		s1 := s
		cats = append(cats, s1)
	}

	err = rows.Err()
	return cats, err
}

// NewPic help
func (r *Report) NewPic(f io.Reader, fn string) (*Pic, error) {
	p, err := NewPic(r.ID, f, fn)
	return p, err
}

// Pics of report
func (r *Report) Pics() ([]string, error) {
	ps, err := PicsByReport(r.ID)
	return ps, err
}

// ReportByID h
func ReportByID(id string) (*Report, error) {
	idi, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.New("id not a valid uint64")
	}

	r := new(Report)
	config.DB.Lock()
	defer config.DB.Unlock()
	err = config.DB.QueryRow("SELECT pk_reportid, title, descr, created_at, modified_at, latitude, longitude, addr, curr_state, state_modified_at, fk_userid FROM reports WHERE pk_reportid = ?", idi).Scan(&r.ID, &r.Title, &r.Desc, &r.CreatedAt, &r.ModifiedAt, &r.Latitude, &r.Longtitude, &r.Address, &r.State, &r.StateModifiedAt, &r.UID)
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

// Reports returns list of reports
// Update Cname for each report after this
// this impacts performance somewhat
func Reports() ([]*Report, error) {
	var reports []*Report

	config.DB.Lock()
	defer config.DB.Unlock()

	rows, err := config.DB.Query("SELECT pk_reportid, title, descr, reports.created_at, addr, curr_state, fk_userid, users.username FROM reports JOIN users ON users.pk_userid = reports.fk_userid;")
	if err != nil {
		return nil, err
	}

	var r Report
	for rows.Next() {
		err := rows.Scan(&r.ID, &r.Title, &r.Desc, &r.CreatedAt, &r.Address, &r.State, &r.UID, &r.Username)
		if err != nil {
			return nil, err
		}

		reports = append(reports, &Report{
			ID:        r.ID,
			Title:     r.Title,
			Desc:      r.Desc,
			CreatedAt: r.CreatedAt,
			Address:   r.Address,
			State:     r.State,
			UID:       r.UID,
			Username:  r.Username,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return reports, nil
}

// ReportsLatest returns the n latest reports
func ReportsLatest(n uint64) ([]*Report, error) {
	var reports []*Report

	config.DB.Lock()
	defer config.DB.Unlock()

	rows, err := config.DB.Query("SELECT pk_reportid, title, descr, reports.created_at, addr, curr_state, fk_userid, users.username FROM reports JOIN users ON users.pk_userid = reports.fk_userid WHERE curr_state > 1 ORDER BY reports.pk_reportid DESC LIMIT ?;", n)
	if err != nil {
		return nil, err
	}

	var r Report
	for rows.Next() {
		err := rows.Scan(&r.ID, &r.Title, &r.Desc, &r.CreatedAt, &r.Address, &r.State, &r.UID, &r.Username)
		if err != nil {
			return nil, err
		}

		reports = append(reports, &Report{
			ID:        r.ID,
			Title:     r.Title,
			Desc:      r.Desc,
			CreatedAt: r.CreatedAt,
			Address:   r.Address,
			State:     r.State,
			UID:       r.UID,
			Username:  r.Username,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return reports, nil
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

// EditUser gives ownership to a given user
func (r *Report) EditUser(id uint64) error {

	config.DB.Lock()
	defer config.DB.Unlock()
	now := time.Now()
	stmt, err := config.DB.Prepare("UPDATE reports SET fk_userid = ?, modified_at = ? WHERE pk_reportid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id, now, r.ID)
	if err != nil {
		return err
	}

	r.UID = id
	r.ModifiedAt = now
	return nil
}

// Solve partially/completely solves the report.
// it expects a call from an auth responsible
// of one subreports
func (r *Report) Solve(auth string) error {
	//SELECT pk_catid FROM categories JOIN authorities ON pk_userid=1 AND pk_userid=fk_userid;
	/*
		SELECT fk_reportid, fk_catid FROM (SELECT pk_catid FROM categories JOIN authorities ON pk_userid=231321 AND pk_userid=fk_userid) as a JOIN subreports ON pk_catid=fk_catid AND fk_reportid=2;
	*/

	a, err := AuthByUsername(auth)
	if err != nil {
		return err
	}

	config.DB.Lock()
	defer config.DB.Unlock()
	rows, err := config.DB.Query("SELECT fk_catid FROM (SELECT pk_catid FROM categories JOIN authorities ON pk_userid=? AND pk_userid=fk_userid) as a JOIN subreports ON pk_catid=fk_catid AND fk_reportid=?;", a.ID, r.ID)
	if err != nil {
		return err
	}

	var cats []uint
	var cat uint

	for rows.Next() {
		err := rows.Scan(&cat)
		if err != nil {
			return err
		}
		cats = append(cats, cat)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	rows.Close()

	for _, cat := range cats {
		stmt, err := config.DB.Prepare("UPDATE subreports SET curr_state = ?  WHERE fk_catid = ? AND fk_reportid = ? ")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(ReportSolved, cat, r.ID)
		if err != nil {
			return err
		}
	}

	var count int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM subreports WHERE fk_reportid = ? ;", r.ID).Scan(&count)
	if err != nil {
		return err
	}
	if len(cats) == count {
		err := r.SetState(ReportSolved)
		if err != nil {
			return err
		}
	}
	return nil

}
