package mobile

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rter/storage"
	"rter/utils"
	"strconv"
	"strings"
	"time"
)

func MultiUploadHandler(w http.ResponseWriter, r *http.Request) {
	imageFile, header, err := r.FormFile("image")
	if err != nil {
		return
	}

	phoneID := r.FormValue("phone_id")

	if !utils.PhoneIDValidator.MatchString(phoneID) {
		http.Error(w, "Malformed Request: Invalid phone_id", http.StatusBadRequest)
		log.Println("upload failed, phone_id malformed:", phoneID)
		return
	}

	rows := storage.MustQuery("SELECT * FROM phones where phone_id = ?;", phoneID)

	if !rows.Next() {
		http.Error(w, "Malformed Request: Invalid phone_id", http.StatusBadRequest)
		log.Println("upload failed, phone_id invalid:", phoneID)
		return
	}

	os.Mkdir(filepath.Join(utils.UploadPath, phoneID), os.ModeDir|0755)

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
	path := utils.UploadPath

	if strings.HasSuffix(header.Filename, ".png") {
		path = filepath.Join(path, fmt.Sprintf("%v/%v.png", phoneID, t.UnixNano()))
	} else if strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, "jpeg") {
		path = filepath.Join(path, fmt.Sprintf("%v/%v.jpg", phoneID, t.UnixNano()))
	}

	outputFile, err := os.Create(path)
	utils.Must(err)
	defer outputFile.Close()

	io.Copy(outputFile, imageFile)

	path = path[len(utils.RterDir):]

	if valid_pos && valid_heading {
		storage.MustQuery("INSERT INTO content (content_id, content_type, filepath, geolat, geolng, heading) VALUES(?, \"mobile\", ?, ?, ?, ?);", phoneID, path, lat, lng, heading)
	} else if valid_pos {
		storage.MustQuery("INSERT INTO content (content_id, content_type, filepath, geolat, geolng) VALUES(?, \"mobile\", ?, ?, ?);", phoneID, path, lat, lng)
	} else {
		storage.MustQuery("INSERT INTO content (content_id, content_type, filepath) VALUES(?, \"mobile\",  ?);", phoneID, path)
	}

	rows = storage.MustQuery("SELECT target_heading from phones where phone_id=?", phoneID)

	if rows.Next() {
		var target_heading []byte
		err := rows.Scan(&target_heading)
		utils.Must(err)
		w.Write(target_heading)
	}

	log.Println("upload complete, phone_id", phoneID, ", heading", valid_heading, ", position", valid_pos)
}
