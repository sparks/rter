package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	p, error := ioutil.ReadAll(r.Body)
	if error != nil {
		return
	}

	ioutil.WriteFile(imagePath+"test.png", p, 0600)
}

func MultiUploadHandler(w http.ResponseWriter, r *http.Request) {
	imageFile, header, error := r.FormFile("image")
	// checkError(error)
	if error != nil {
		return
	}

	phoneID := r.FormValue("phone_id")

	if !phoneIDValidator.MatchString(phoneID) {
		http.Error(w, "Malformed Request: Invalid phone_id", http.StatusBadRequest)
		fmt.Println("upload failed, phone_id malformed:", phoneID)
		return
	}

	rows, _, err := database.Query("SELECT * FROM phones where phone_id = \"%v\";", phoneID)

	if len(rows) == 0 {
		http.Error(w, "Malformed Request: Invalid phone_id", http.StatusBadRequest)
		fmt.Println("upload failed, phone_id invalid:", phoneID)
		return
	}

	os.Mkdir(imagePath+phoneID, os.ModeDir|0755)

	valid_pos := true
	valid_heading := true

	lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
	if err != nil {
		valid_pos = false
	}

	lng, err := strconv.ParseFloat(r.FormValue("lng"), 64)
	if err != nil {
		valid_pos = false
	}

	heading, err := strconv.ParseFloat(r.FormValue("heading"), 64)
	if err != nil {
		valid_heading = false
	}

	t := time.Now()
	path := imagePath

	if strings.HasSuffix(header.Filename, ".png") {
		path += fmt.Sprintf("%v/%v.png", phoneID, t.UnixNano())
	} else if strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, "jpeg") {
		path += fmt.Sprintf("%v/%v.jpg", phoneID, t.UnixNano())
	}

	outputFile, error := os.Create(path)
	checkError(error)
	defer outputFile.Close()

	io.Copy(outputFile, imageFile)

	if valid_pos && valid_heading {
		_, _, error = database.Query("INSERT INTO content (content_id, content_type, filepath, geolat, geolng, heading) VALUES(\"%s\", \"mobile\", \"%s\", %v, %v, %v);", phoneID, path, lat, lng, heading)
	} else if valid_pos {
		_, _, error = database.Query("INSERT INTO content (content_id, content_type, filepath, geolat, geolng) VALUES(\"%s\", \"mobile\", \"%s\", %v, %v);", phoneID, path, lat, lng)
	} else {
		_, _, error = database.Query("INSERT INTO content (content_id, content_type, filepath) VALUES(\"%s\", \"mobile\", \"%s\");", phoneID, path)
	}
	checkError(error)

	rows, _, err = database.Query("SELECT target_heading from phones where phone_id=\"%s\"", phoneID)
	checkError(err)

	if len(rows) > 0 {
		switch v := rows[0][0].(type) {
		case []byte:
			w.Write(v)
		default:
			w.Write([]byte("0.0"))
		}
	}

	fmt.Println("upload complete, phone_id", phoneID)
}
