package server

import (
	"encoding/json"
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

		for _, tile := range layout {
			if phoneIDValidator.MatchString(tile.ContentID) {
				tile.Sanitize()
				_, _, err = database.Query("INSERT INTO layout (content_id, col, row, size_x, size_y) VALUES(\"%s\", %d, %d, %d, %d) ON DUPLICATE KEY UPDATE col=%d, row=%d, size_x=%d, size_y=%d;", tile.ContentID, tile.Col, tile.Row, tile.SizeX, tile.SizeY, tile.Col, tile.Row, tile.SizeX, tile.SizeY)
				checkError(err)
			}
		}

	} else if r.URL.Path == "/ajax/getlayout" {
		rows, _, err := database.Query("SELECT content_id, col, row, size_x, size_y FROM layout ORDER BY col, row;")

		layout := make([]*LayoutTile, len(rows))

		for i, row := range rows {
			layout[i] = &LayoutTile{
				row.Str(0),
				row.Int(1),
				row.Int(2),
				row.Int(3),
				row.Int(4),
			}
		}

		layoutJSON, err := json.Marshal(layout)
		checkError(err)

		w.Write(layoutJSON)
	}
}

func (tile *LayoutTile) Sanitize() {
	if tile.SizeX < 1 {
		tile.SizeX = 1
	}
	if tile.SizeY < 1 {
		tile.SizeY = 1
	}
	if tile.Col < 1 {
		tile.Col = 1
	}
	if tile.Row < 1 {
		tile.Row = 1
	}
}
