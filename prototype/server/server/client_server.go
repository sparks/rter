package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
	"sync"
)

type PageContent struct {
	Phones []*PhoneContent `json:"phones"`
	Layout []*LayoutTile   `json:"layout"`
}

type PhoneContent struct {
	ContentID string  `json:"content_id"`
	Filepath  string  `json:"filepath"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
}

type LayoutTile struct {
	ContentID string `json:"content_id"`
	Col       int    `json:"col"`
	Row       int    `json:"row"`
	SizeX     int    `json:"size_x"`
	SizeY     int    `json:"size_y"`
}

var templates = template.Must(template.ParseFiles(templatePath + "main.html"))

var writeLock sync.Mutex

var rowsMatchedValidator = regexp.MustCompile(".*0.*0.*0")

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	templates = template.Must(template.ParseFiles(templatePath + "main.html")) // TODO: For dev only, remove if deployed. Reloads HTML every request instead of caching

	err := templates.ExecuteTemplate(w, "main.html", fetchPageContent())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ClientAjax(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path == "/ajax/pushlayout" {
		decoder := json.NewDecoder(r.Body)

		var layout []*LayoutTile
		err := decoder.Decode(&layout)
		checkError(err)

		writeLock.Lock()
		for _, tile := range layout {
			if phoneIDValidator.MatchString(tile.ContentID) {
				tile.Sanitize()
				
				_, res, err := database.Query("UPDATE layout SET col=%d, row=%d, size_x=%d, size_y=%d WHERE content_id=\"%s\";", tile.Col, tile.Row, tile.SizeX, tile.SizeY, tile.ContentID) 
				checkError(err)
				
				// Check that no rows were matched via a ghetto regex, since the Result object returned contains no public field for matched rows, only affected rows.
				if rowsMatchedValidator.MatchString(res.Message()) {
					_, res, err = database.Query("INSERT INTO layout (content_id, col, row, size_x, size_y) VALUES(\"%s\", %d, %d, %d, %d);", tile.ContentID, tile.Col, tile.Row, tile.SizeX, tile.SizeY)
					checkError(err)
				}
				
			}
		}
		writeLock.Unlock()

	} else if r.URL.Path == "/ajax/getlayout" {
		layoutJSON, err := json.Marshal(fetchPageContent())
		checkError(err)

		w.Write(layoutJSON)
	}
}

func fetchPageContent() *PageContent {
	// rows, _, err := database.Query("SELECT content.content_id, content.filepath, content.geolat, content.geolng, layout.col, layout.row, layout.size_x, layout.size_y FROM content LEFT JOIN layout ON (layout.content_id = content.content_id) WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	phoneRows, _, err := database.Query("SELECT content.content_id, content.filepath, content.geolat, content.geolng FROM content WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")
	checkError(err)

	phones := make([]*PhoneContent, len(phoneRows))

	for i, row := range phoneRows {
		phones[i] = &PhoneContent{
			row.Str(0),
			row.Str(1),
			row.Float(2),
			row.Float(3),
		}
	}

	layoutRows, _, err := database.Query("SELECT content_id, col, row, size_x, size_y FROM layout ORDER BY col, row;")
	checkError(err)

	layout := make([]*LayoutTile, len(layoutRows))

	for i, row := range layoutRows {
		layout[i] = &LayoutTile{
			row.Str(0),
			row.Int(1),
			row.Int(2),
			row.Int(3),
			row.Int(4),
		}
	}

	return &PageContent{phones, layout}
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
