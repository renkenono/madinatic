package model

import (
	"testing"

	"github.com/renkenn/madinatic/config"
)

func TestReport(t *testing.T) {
	dsn, err := config.App.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
	}

	err = config.DB.InitDB(dsn)
	if err != nil {
		t.Error(err)
	}

	_, err = NewReport("admin", "Report non-approuvé 3", "Description of Report 1", "Welhaṣa", "35.235330", "-1.506706")
	if err != nil {
		t.Error(err)
	}
	// err = r.EditTitle("Edited Report 0")
	// if err != nil {
	// 	t.Error(err)
	// }
	// err = r.EditDesc("Edited Description of Report 0")
	// if err != nil {
	// 	t.Error(err)
	// }
	// err = r.SetState(ReportAccepted)
	// if err != nil {
	// 	t.Error(err)
	// }

	// rs, err := ReportsByUser("renken")
	// if err != nil {
	// 	t.Error(err)
	// }

	// rs, err = ReportsByState(ReportAccepted)
	// if err != nil {
	// 	t.Error(err)
	// }

	// for _, r := range rs {
	// 	err = r.Delete()
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// }

	// _, err = ReportByID("50")
	// if err != nil {
	// 	t.Error(err)
	// }
}
