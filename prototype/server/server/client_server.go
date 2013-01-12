package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Page struct {
	Phones []*PhoneContent
}

type PhoneContent struct {
	ContentID string  `json:"content_id"`
	Filepath  string  `json:"filepath"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
}

type LayoutTile struct {
	ContentID string `json:"content_id"`
	Col       int    `json:"col"`
	Row       int    `json:"row"`
	SizeX     int    `json:"size_x"`
	SizeY     int    `json:"size_y"`
}

var templates = template.Must(template.ParseFiles(templatePath + "main.html"))

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// rows, _, err := database.Query("SELECT content.content_id, content.filepath, content.geolat, content.geolong, layout.col, layout.row, layout.size_x, layout.size_y FROM content LEFT JOIN layout ON (layout.content_id = content.content_id) WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	rows, _, err := database.Query("SELECT content.content_id, content.filepath, content.geolat, content.geolong FROM content WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	phones := make([]*PhoneContent, len(rows))

	for i, row := range rows {
		phones[i] = &PhoneContent{
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

	if r.URL.Path == "/ajax/pushlayout" {
		decoder := json.NewDecoder(r.Body)

		var layout []*LayoutTile
		err := decoder.Decode(&layout)
		checkError(err)

		v, err := json.Marshal(layout)
		checkError(err)

		fmt.Println(string(v))
	} else if r.URL.Path == "/ajax/getlayout" {
		// rows, _, err := database.Query("SELECT content_id, col, row, size_x, size_y FROM layout;")

	}
}
