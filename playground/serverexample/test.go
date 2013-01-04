package main

import (
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

const lenPath = len("/view/")
const imgPath = "images/"

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "main.html"))
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
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

type Image struct{}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, 200, 200)
}

func (i Image) At(x, y int) color.Color {
	xf := float64(x)
	yf := float64(y)

	a := []float64{.0000001, -0.00000012, -0.00000011, 0.0000009, 0.0000008, -0.0000007, -0.0000009, -0.0000001}
	b := []float64{.0000005, -0.00000016, -0.000001, 0.0000005, -0.0000008, 0.0000007, -0.0000009, -0.0000003}
	c := []float64{-.0000003, -0.00000012, -0.00000014, -0.00000013, 0.0000006, 0.0000007, -0.0000009, 0.000000111}

	return color.RGBA{
		uint8(a[0]*xf + a[1]*math.Pow(xf, 2) + a[2]*math.Pow(xf, 3) + a[3]*math.Pow(xf, 4) + a[4]*yf + a[5]*math.Pow(yf, 2) + a[6]*math.Pow(yf, 3) + a[7]*math.Pow(yf, 4)),
		uint8(b[0]*xf + b[1]*math.Pow(xf, 2) + b[2]*math.Pow(xf, 3) + b[3]*math.Pow(xf, 4) + b[4]*yf + b[5]*math.Pow(yf, 2) + b[6]*math.Pow(yf, 3) + b[7]*math.Pow(yf, 4)),
		uint8(c[0]*xf + c[1]*math.Pow(xf, 2) + c[2]*math.Pow(xf, 3) + c[3]*math.Pow(xf, 4) + c[4]*yf + c[5]*math.Pow(yf, 2) + c[6]*math.Pow(yf, 3) + c[7]*math.Pow(yf, 4)),
		255,
	}
}

func pngHandler(w http.ResponseWriter, r *http.Request) {
	i := &Image{}
	png.Encode(w, i)
}

func imgHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".png") {
		http.NotFound(w, r)
		return
	}

	path := strings.Split(r.URL.Path, "/")

	buff, err := ioutil.ReadFile(imgPath + path[len(path)-1])
	if err != nil {
		http.NotFound(w, r)
		return
		// panic(err)
	}

	w.Write(buff)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	ioutil.WriteFile("test.png", p, 0600)
}

func main() {
	http.HandleFunc("/test.png", pngHandler)
	http.HandleFunc("/images/", imgHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/", makeHandler(mainHandler))
	http.ListenAndServe(":8080", nil)
}
