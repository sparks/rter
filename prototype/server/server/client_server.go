package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type PageContent struct {
	Content []*ContentChunk `json:"content"`
}

type ContentChunk struct {
	ContentID     string  `json:"content_id"`
	ConentType    string  `json:"content_type"`
	Filepath      string  `json:"filepath"`
	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
	Heading       float64 `json:"heading"`
	TargetHeading float64 `json:"target_heading"`
	Col           int     `json:"col"`
	Row           int     `json:"row"`
	SizeX         int     `json:"size_x"`
	SizeY         int     `json:"size_y"`
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

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	url := r.FormValue("url")
	imageFile, header, error := r.FormFile("image")
	
	id := "client"
	
	if error != nil {
		path := imagePath + "client/default.png"
		_, _, error = database.Query("INSERT INTO content (content_id, content_type, filepath, description, url) VALUES(\"%s\", \"web\", \"%s\", \"%s\", \"%s\");", id, path, description, url)
		checkError(error)
	} else {
		os.Mkdir(imagePath + id, os.ModeDir | 0755)
	
		t := time.Now()
		path := imagePath
		
		if strings.HasSuffix(header.Filename, ".png") {
			path += fmt.Sprintf("%v/%v.png", id, t.UnixNano())
		} else if strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, "jpeg") {
			path += fmt.Sprintf("%v/%v.jpg", id, t.UnixNano())
		}
		
		outputFile, error := os.Create(path)
		checkError(error)
		defer outputFile.Close()
		
		io.Copy(outputFile, imageFile)
		
		_, _, error = database.Query("INSERT INTO content (content_id, content_type, filepath, description, url) VALUES(\"%s\", \"web\", \"%s\", \"%s\", \"%s\");", id, path, description, url)
		checkError(error)
	}
	
	http.Redirect(w, r, "/", http.StatusFound)
}

func ClientAjax(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path == "/ajax/pushlayout" {
		decoder := json.NewDecoder(r.Body)

		var layout []*ContentChunk
		err := decoder.Decode(&layout)
		checkError(err)

		writeLock.Lock()
		for _, chunk := range layout {
			if phoneIDValidator.MatchString(chunk.ContentID) {
				chunk.SanitizeLayout()

				_, res, err := database.Query("UPDATE layout SET col=%d, row=%d, size_x=%d, size_y=%d WHERE content_id=\"%s\";", chunk.Col, chunk.Row, chunk.SizeX, chunk.SizeY, chunk.ContentID)
				checkError(err)

				// Check that no rows were matched via a ghetto regex, since the Result object returned contains no public field for matched rows, only affected rows.
				if rowsMatchedValidator.MatchString(res.Message()) {
					_, _, err = database.Query("INSERT INTO layout (content_id, col, row, size_x, size_y) VALUES(\"%s\", %d, %d, %d, %d);", chunk.ContentID, chunk.Col, chunk.Row, chunk.SizeX, chunk.SizeY)
					checkError(err)
				}

			}
		}
		writeLock.Unlock()

	} else if r.URL.Path == "/ajax/getlayout" {
		layoutJSON, err := json.Marshal(fetchPageContent().Content)
		checkError(err)

		w.Write(layoutJSON)
	} else if r.URL.Path == "/ajax/pushheading" {
		decoder := json.NewDecoder(r.Body)

		var targetHeading *ContentChunk
		err := decoder.Decode(&targetHeading)
		checkError(err)

		_, _, err = database.Query("UPDATE phone_id SET target_heading=%v WHERE phone_id=\"%s\";", targetHeading.TargetHeading, targetHeading.ContentID)
		checkError(err)
	}
}

func fetchPageContent() *PageContent {
	rows, _, err := database.Query("SELECT content.content_id, content.content_type, content.filepath, content.geolat, content.geolng, content.heading, phones.target_heading, layout.col, layout.row, layout.size_x, layout.size_y FROM content LEFT JOIN (layout, phones) ON (layout.content_id = content.content_id AND phones.phone_id = content.content_id) WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	// phoneRows, _, err := database.Query("SELECT content.content_id, content.filepath, content.geolat, content.geolng, content.heading FROM content WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	checkError(err)

	content := make([]*ContentChunk, len(rows))

	for i, row := range rows {
		content[i] = &ContentChunk{
			row.Str(0),   //content_id
			row.Str(1),   //content_type
			row.Str(2),   //filepath
			row.Float(3), //geolat
			row.Float(4), //geolng
			row.Float(5), //heading
			row.Float(6), //target_heading
			row.Int(7),   //col
			row.Int(8),   //row
			row.Int(9),   //size_x
			row.Int(10),  //size_y
		}
	}

	return &PageContent{content}
}

func (chunk *ContentChunk) SanitizeLayout() {
	if chunk.SizeX < 1 {
		chunk.SizeX = 1
	}
	if chunk.SizeY < 1 {
		chunk.SizeY = 1
	}
	if chunk.Col < 1 {
		chunk.Col = 1
	}
	if chunk.Row < 1 {
		chunk.Row = 1
	}
}
