package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

type Message struct {
	Name string
	Body string
	Time int64
}

type Page struct {
	Title string
	Body  []byte
}

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
var templates = template.Must(template.ParseFiles(templatePath + "main.html"))

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	// title := r.URL.Path
	// if !titleValidator.MatchString(title) {
	// 	http.NotFound(w, r)
	// 	return
	// }

	p := &Page{"", []byte{}}

	err := templates.ExecuteTemplate(w, "main.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ClientAjax(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application.json")

	fmt.Println("URL: ", r.URL)

	req := make([]byte, 100)
	n, ok := r.Body.Read(req)

	fmt.Println("Body Result: ", n, ok)
	fmt.Println("Body", string(req))

	socks := r.FormValue("socks")
	apples := r.FormValue("apples")

	fmt.Println("Form Values: ", socks, apples)

	m := Message{"Alice", "Hello", 1294706395881547000}

	b, _ := json.Marshal(m)

	fmt.Println(string(b))

	w.Write(b)
}
