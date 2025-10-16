package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/iraqnroll/gochan/context"
	"github.com/iraqnroll/gochan/models"
)

type Template struct {
	htmlTpl *template.Template
}

type NavbarData struct {
	BoardList []models.BoardDto
}

type FooterData struct {
	Sitename string
}

type BasePageData struct {
	Navbar   NavbarData
	Footer   FooterData
	PageData any
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
	}

	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"boardList": func() ([]models.BoardDto, error) {
				return nil, nil
			},
		},
	)

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}

	//Might not be a good idea for large pages, maybe writting directly to responseWriter would be better in the future...
	io.Copy(w, &buf)
}

// For embedded templates
func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	htmlTpl := template.New(patterns[0])
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser not implemented")
			},
		},
	)

	htmlTpl, err := htmlTpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template : %w", err)
	}

	return Template{
		htmlTpl: htmlTpl,
	}, nil
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}
