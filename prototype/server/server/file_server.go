package server

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
var templates = template.Must(template.ParseFiles(templatePath+"edit.html", templatePath+"view.html", templatePath+"main.html"))

func (p *Page) save() error {
	filename := dataPath + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := dataPath + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func MainHandler(w http.ResponseWriter, r *http.Request, title string) {
	if title == "" {
		p := &Page{"Main", []byte("Welcome!")}
		renderTemplate(w, "main", p)
	} else {
		p, err := loadPage(title)
		if err != nil {
			http.Redirect(w, r, "/edit/"+title, http.StatusFound)
			return
		}
		renderTemplate(w, "view", p)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := ""

		if len(r.URL.Path) >= lenPath {
			title = r.URL.Path[lenPath:]
			if !titleValidator.MatchString(title) {
				http.NotFound(w, r)
				return
			}
		} else if len(r.URL.Path) > 1 {
			title = r.URL.Path[1:]
			if !titleValidator.MatchString(title) {
				http.NotFound(w, r)
				return
			}
		}

		fn(w, r, title)
	}
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".png") {
		http.NotFound(w, r)
		return
	}

	path := strings.Split(r.URL.Path, "/")

	fi, err := os.Open(imagePath + path[len(path)-1])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer fi.Close()

	io.Copy(w, fi)
}
