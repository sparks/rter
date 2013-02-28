package server

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type PageContent struct {
	Content []*ContentChunk `json:"content"`
}

type ContentChunk struct {
	ContentID  string `json:"content_id"`
	ConentType string `json:"content_type"`

	Filepath string `json:"filepath"`
	URL      string `json:"url"`

	Description string `json:"description"`

	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
	Heading       float64 `json:"heading"`
	TargetHeading float64 `json:"target_heading"`

	Col   int `json:"col"`
	Row   int `json:"row"`
	SizeX int `json:"size_x"`
	SizeY int `json:"size_y"`
}

var templates = template.Must(template.ParseFiles(filepath.Join(TemplatePath, "index.html")))

var writeLock sync.Mutex

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		http.ServeFile(w, r, filepath.Join(TemplatePath, r.URL.Path))
	} else {
		err := templates.ExecuteTemplate(w, "index.html", fetchPageContent())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	url := r.FormValue("url")
	imageFile, header, err := r.FormFile("image")

	t := time.Now()

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%v", t.UnixNano())))
	id := fmt.Sprintf("%x", hasher.Sum(nil))

	if err != nil {
		path := filepath.Join(ImagePath, "client", "default.png")
		path = path[len(rterDir):]
		_, err = db.Exec("INSERT INTO content (content_id, content_type, filepath, description, url) VALUES(?, ?, ?, ?);", id, path, description, url)
		checkError(err)
	} else {
		os.Mkdir(filepath.Join(ImagePath, "client"), os.ModeDir|0755)

		path := filepath.Join(ImagePath, "client")

		if strings.HasSuffix(header.Filename, ".png") {
			path = filepath.Join(path, fmt.Sprintf("%v.png", id))
		} else if strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, "jpeg") {
			path = filepath.Join(path, fmt.Sprintf("%v.jpg", id))
		}

		outputFile, err := os.Create(path)
		checkError(err)
		defer outputFile.Close()

		io.Copy(outputFile, imageFile)

		path = path[len(rterDir):]

		_, err = db.Exec("INSERT INTO content (content_id, content_type, filepath, description, url) VALUES(?, ?, ?, ?);", id, path, description, url)
		checkError(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func ClientAjax(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	// w.Header().Set("Content-Type", "application/json")

	if r.URL.Path == "/ajax/pushlayout" {
		decoder := json.NewDecoder(r.Body)

		var layout []*ContentChunk
		err := decoder.Decode(&layout)
		checkError(err)

		writeLock.Lock()
		for _, chunk := range layout {
			if phoneIDValidator.MatchString(chunk.ContentID) {
				chunk.SanitizeLayout()

				rows, err := db.Query("SELECT uid FROM layout WHERE content_id=?", chunk.ContentID)
				checkError(err)

				if rows.Next() {
					_, err := db.Exec("UPDATE layout SET col=?, row=?, size_x=?, size_y=? WHERE content_id=?;", chunk.Col, chunk.Row, chunk.SizeX, chunk.SizeY, chunk.ContentID)
					checkError(err)
				} else {
					_, err = db.Exec("INSERT INTO layout (content_id, col, row, size_x, size_y) VALUES(?, ?, ?, ?, ?);", chunk.ContentID, chunk.Col, chunk.Row, chunk.SizeX, chunk.SizeY)
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

		var headingChunk *ContentChunk
		err := decoder.Decode(&headingChunk)
		checkError(err)

		_, err = db.Query("UPDATE phones SET target_heading=? WHERE phone_id=?;", headingChunk.TargetHeading, headingChunk.ContentID)
		checkError(err)
	} else if r.URL.Path == "/ajax/pushdescription" {
		decoder := json.NewDecoder(r.Body)

		var descChunk *ContentChunk
		err := decoder.Decode(&descChunk)
		checkError(err)

		_, err = db.Query("UPDATE content as c1 INNER JOIN (select c2.uid from content as c2 where c2.content_id=?  ORDER by c2.timestamp DESC LIMIT 1) as x ON x.uid = c1.uid  SET c1.description=?;", descChunk.ContentID, descChunk.Description)
		checkError(err)
	}
}

func fetchPageContent() *PageContent {
	// rows, _, err := database.Query("SELECT content.content_id, content.content_type, content.filepath, content.geolat, content.geolng, content.heading, phones.target_heading, layout.col, layout.row, layout.size_x, layout.size_y FROM content  LEFT JOIN layout ON layout.content_id = content.content_id LEFT JOIN phones ON phones.phone_id = content.content_id WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	// phoneRows, _, err := database.Query("SELECT content.content_id, content.filepath, content.geolat, content.geolng, content.heading FROM content WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	rows, err := db.Query("select c1.content_id, c1.content_type, c1.filepath, c1.url, c1.description, c1.geolat, c1.geolng, c1.heading, phones.target_heading, layout.col, layout.row, layout.size_x, layout.size_y from (select content_id, max(timestamp) as maxtimestamp from content group by content_id) as c2 inner join content as c1 on c1.content_id = c2.content_id and c1.timestamp = c2.maxtimestamp LEFT JOIN layout ON layout.content_id = c1.content_id LEFT JOIN phones ON phones.phone_id = c1.content_id;")

	checkError(err)

	content := make([]*ContentChunk, 0)

	for i := 0; rows.Next(); i++ {
		chunk := &ContentChunk{}

		rows.Scan(
			&chunk.ContentID,
			&chunk.ConentType,
			&chunk.Filepath,
			&chunk.URL,
			&chunk.Description,
			&chunk.Lat,
			&chunk.Lng,
			&chunk.Heading,
			&chunk.TargetHeading,
			&chunk.Col,
			&chunk.Row,
			&chunk.SizeX,
			&chunk.SizeY,
		)
		content = append(content, chunk)
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
