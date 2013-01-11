package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Message struct {
	Name string
	Body string
	Time int64
}

type Page struct {
	Phones []*Phone
}

type Phone struct {
	ID        string
	Filepath  string
	Lat, Long float64
}

var templates = template.Must(template.ParseFiles(templatePath + "main.html"))

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	rows, _, err := database.Query("SELECT phone_id, filepath, geolat, geolong FROM content where (select count(*) from content as c where c.phone_id = content.phone_id and c.timestamp <= content.timestamp) <= 1;")

	phones := make([]*Phone, len(rows))

	for i, row := range rows {
		phones[i] = &Phone{
			row.Str(0),
			row.Str(1),
			row.Float(2),
			row.Float(3),
		}
	}

	p := &Page{phones}

	templates = template.Must(template.ParseFiles(templatePath + "main.html")) //TODO: For dev only, remove if deployed. Reloads HTML every request instead of caching

	err = templates.ExecuteTemplate(w, "main.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ClientAjax(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")

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

	// fmt.Println(string(b))

	w.Write(b)
}
