package report

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/renkenn/madinatic/config"
	"github.com/renkenn/madinatic/model"

	"github.com/gorilla/csrf"
	. "github.com/renkenn/madinatic/control"
)

const (
	rcreateAPIErr = "report create API failed"
	rcreateErr    = "report create failed"
)

type repAPI struct {
	ID uint64 `json:"id"`
}

// CreateAPI recv HTTP POST Form
func CreateAPI(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.ERROR, rcreateAPIErr, err.Error())
		return
	}

	// chink token into username :DDDDDDDDDDDDD
	report, err := crep(r, r.MultipartForm.Value["token"][0], rcreateAPIErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.ERROR, rcreateAPIErr, err.Error())
		return
	}
	var rep repAPI
	rep.ID = report.ID
	MarshalJSON(w, rep, rcreateAPIErr)

}

// Create a report
func Create(w http.ResponseWriter, r *http.Request) {
	s, err := config.Store.Get(r, "userdata")
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		log.Printf("%s%s: %s", config.INFO, rcreateErr, err.Error())
		return
	}

	_, ok := s.Values["username"]
	if s.IsNew || !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == http.MethodGet {
		cats, err := model.Cats()
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusInternalServerError)
			log.Printf("%s%s: %s", config.INFO, rcreateErr, err.Error())
			return
		}

		var cs []*model.Cat
		for i := 0; i < len(cats)-1; i++ {
			in := false
			for j := i + 1; j < len(cats); j++ {
				if cats[i].Name == cats[j].Name {
					in = true
					break
				}
			}
			if !in || i == len(cats)-2 {
				cs = append(cs, cats[i])
			}
		}
		data := map[string]interface{}{
			"csrfField": csrf.TemplateField(r),
			"Cats":      cs,
		}

		Render(w, r, data, ViewReportCreate, "report_create.tmpl")
		return
	}

	// handle form

	log.Printf("%screate report POST", config.INFO)
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}

	username := s.Values["username"].(string)

	// TODO: update this
	nr, err := crep(r, username, rcreateErr)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}

	str := fmt.Sprintf("/report/view/%d", nr.ID)

	http.Redirect(w, r, str, http.StatusFound)

}

func crep(r *http.Request, username, errstr string) (*model.Report, error) {
	selectedCs := r.MultipartForm.Value["cat"]
	if len(selectedCs) == 0 {
		return nil, errors.New("no cats selected")
	}

	re, err := model.NewReport(username, r.MultipartForm.Value["title"][0], r.MultipartForm.Value["desc"][0], r.MultipartForm.Value["address"][0], r.MultipartForm.Value["lat"][0], r.MultipartForm.Value["long"][0])
	if err != nil {
		return nil, err
	}

	for _, sc := range selectedCs {
		err := re.NewSubReport(sc)
		if err != nil {

			log.Printf("%s%s: %s", config.INFO, errstr, err.Error())
			err = re.Delete()
			if err != nil {
				log.Printf("%s%s: %s", config.INFO, errstr, err.Error())
			}
			return nil, err
		}
	}

	pics := r.MultipartForm.File["pics"]

	for _, p := range pics {
		tmp, err := p.Open()
		if err != nil {
			log.Printf("%s%s: %s", config.INFO, errstr, err.Error())
			return nil, err
		}
		_, err = re.NewPic(tmp, p.Filename)
		if err != nil {
			log.Printf("%s%s: %s", config.INFO, errstr, err.Error())
			return nil, err
		}
	}
	return re, nil
}
