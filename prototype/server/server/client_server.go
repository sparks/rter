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
	"regexp"
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

var rowsMatchedValidator = regexp.MustCompile(".*0.*0.*0")

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
	imageFile, header, error := r.FormFile("image")

	t := time.Now()

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%v", t.UnixNano())))
	id := fmt.Sprintf("%x", hasher.Sum(nil))

	if error != nil {
		path := filepath.Join(ImagePath, "client", "default.png")
		_, _, error = database.Query("INSERT INTO content (content_id, content_type, filepath, description, url) VALUES(\"%s\", \"web\", \"%s\", \"%s\", \"%s\");", id, path, description, url)
		checkError(error)
	} else {
		os.Mkdir(filepath.Join(ImagePath, "client"), os.ModeDir|0755)

		path := filepath.Join(ImagePath, "client")

		if strings.HasSuffix(header.Filename, ".png") {
			path = filepath.Join(path, fmt.Sprintf("%v.png", id))
		} else if strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, "jpeg") {
			path = filepath.Join(path, fmt.Sprintf("%v.jpg", id))
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

		var headingChunk *ContentChunk
		err := decoder.Decode(&headingChunk)
		checkError(err)

		_, _, err = database.Query("UPDATE phones SET target_heading=%v WHERE phone_id=\"%s\";", headingChunk.TargetHeading, headingChunk.ContentID)
		checkError(err)
	} else if r.URL.Path == "/ajax/pushdescription" {
		decoder := json.NewDecoder(r.Body)

		var descChunk *ContentChunk
		err := decoder.Decode(&descChunk)
		checkError(err)

		_, _, err = database.Query("UPDATE content as c1 INNER JOIN (select c2.uid from content as c2 where c2.content_id=\"%s\"  ORDER by c2.timestamp DESC LIMIT 1) as x ON x.uid = c1.uid  SET c1.description=\"%s\";", descChunk.ContentID, descChunk.Description)
		checkError(err)
	}
}

func fetchPageContent() *PageContent {
	// rows, _, err := database.Query("SELECT content.content_id, content.content_type, content.filepath, content.geolat, content.geolng, content.heading, phones.target_heading, layout.col, layout.row, layout.size_x, layout.size_y FROM content  LEFT JOIN layout ON layout.content_id = content.content_id LEFT JOIN phones ON phones.phone_id = content.content_id WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	// phoneRows, _, err := database.Query("SELECT content.content_id, content.filepath, content.geolat, content.geolng, content.heading FROM content WHERE (SELECT COUNT(*) FROM content AS c WHERE c.content_id = content.content_id AND c.timestamp >= content.timestamp) <= 1;")

	rows, _, err := database.Query("select c1.content_id, c1.content_type, c1.filepath, c1.url, c1.description, c1.geolat, c1.geolng, c1.heading, phones.target_heading, layout.col, layout.row, layout.size_x, layout.size_y from (select content_id, max(timestamp) as maxtimestamp from content group by content_id) as c2 inner join content as c1 on c1.content_id = c2.content_id and c1.timestamp = c2.maxtimestamp LEFT JOIN layout ON layout.content_id = c1.content_id LEFT JOIN phones ON phones.phone_id = c1.content_id;")

	checkError(err)

	content := make([]*ContentChunk, len(rows))

	for i, row := range rows {
		content[i] = &ContentChunk{
			row.Str(0),   //content_id
			row.Str(1),   //content_type
			row.Str(2),   //filepath
			row.Str(3),   //url
			row.Str(4),   //description
			row.Float(5), //geolat
			row.Float(6), //geolng
			row.Float(7), //heading
			row.Float(8), //target_heading
			row.Int(9),   //col
			row.Int(10),  //row
			row.Int(11),  //size_x
			row.Int(12),  //size_y
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
