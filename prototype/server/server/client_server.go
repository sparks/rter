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
	ContentID string  `json:"content_id"`
	Filepath  string  `json:"filepath"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	Col       int     `json:"col"`
	Row       int     `json:"row"`
	SizeX     int     `json:"size_x"`
	SizeY     int     `json:"size_y"`
}

var templates = template.Must(template.ParseFiles(templatePath + "main.html"))

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	rows, _, err := database.Query("SELECT content.content_id, content.filepath, content.geolat, content.geolong, layout.col, layout.row, layout.size_x, layout.size_y FROM content LEFT JOIN layout ON (layout.content_id = content.content_id) WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	phones := make([]*Phone, len(rows))

	for i, row := range rows {
		phones[i] = &Phone{
			row.Str(0),
			row.Str(1),
			row.Float(2),
			row.Float(3),
			row.Int(4),
			row.Int(5),
			row.Int(6),
			row.Int(7),
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
