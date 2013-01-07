package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Page struct {
	Title string
	Body  []byte
}

const lenPath = len("/view/")

const imagePath = "images/"
const templatePath = "templates/"
const dataPath = "data/"

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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func mainHandler(w http.ResponseWriter, r *http.Request, title string) {
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

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
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

func imageHandler(w http.ResponseWriter, r *http.Request) {
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	ioutil.WriteFile(imagePath+"test.png", p, 0600)
}

func multiUploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		panic(err)
	}

	fo, err := os.Create(imagePath + header.Filename)
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	io.Copy(fo, file)
}

func main() {
	
	// Testing database connection, queries and output
	db := mysql.New("tcp", "", "localhost:3306", "root", "", "rter")
	
	err := db.Connect()
	if err != nil {
		panic(err)
	}
	
	rows, _, err := db.Query("select * from content")
	if err != nil {
		panic(err)
	}
	
	for _, row := range rows {
		for _, col := range row {
			if col == nil {
				null := []byte("NULL")
				os.Stdout.Write(null)
			} else {
				// Type assertion required because []interface{} "type" is entirely unknown
				val := col.([]byte)
				os.Stdout.Write(append(val, []byte("  |  ")...))
			}
		}
	}
	
	// Resume normal operation
	http.HandleFunc("/images/", imageHandler)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/multiup", multiUploadHandler)

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.HandleFunc("/", makeHandler(mainHandler))

	http.ListenAndServe(":8080", nil)
}
