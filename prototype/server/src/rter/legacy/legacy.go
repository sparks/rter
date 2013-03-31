// Provides support for the early prototype Android App.strconv
//
// Mimics the behavior of the rtER server v1 and accepts streams of images from android phones running the prototype app. This relies on users with Username matching the phone_id already existing in the DB (much like in the rtER server v1
package legacy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rter/data"
	"rter/storage"
	"rter/utils"
	"strconv"
	"strings"
	"time"
)

// Handle a POST request with an image, phone_id, lat, lng and heading. Return a target heading from the web UI
func MultiUploadHandler(w http.ResponseWriter, r *http.Request) {
	imageFile, header, err := r.FormFile("image")
	if err != nil {
		return
	}

	user := new(data.User)
	user.Username = r.FormValue("phone_id")

	err = storage.Select(user)

	if err == storage.ErrZeroAffected {
		log.Println("upload failed, phone_id invalid:", user.Username)
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	os.Mkdir(filepath.Join(utils.UploadPath, user.Username), os.ModeDir|0755)

	matchingItems := make([]*data.Item, 0)
	err = storage.SelectWhere(&matchingItems, "WHERE Type=\"streaming-video-v0\" AND Author=?", user.Username)

	exists := true

	if err == storage.ErrZeroAffected {
		exists = false
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	var item *data.Item

	if exists {
		item = matchingItems[0]
	} else {
		item = new(data.Item)
	}

	item.Author = user.Username

	item.Type = "streaming-video-v0"
	item.Live = true
	item.HasGeo = true
	item.HasHeading = true

	item.Lat, err = strconv.ParseFloat(r.FormValue("lat"), 64)
	if err != nil {
		item.HasGeo = false
	}

	item.Lng, err = strconv.ParseFloat(r.FormValue("lng"), 64)
	if err != nil {
		item.HasGeo = false
	}

	item.Heading, err = strconv.ParseFloat(r.FormValue("heading"), 64)
	if err != nil {
		item.HasHeading = false
	}

	item.StartTime = time.Now()
	item.StopTime = item.StartTime
	path := utils.UploadPath

	if strings.HasSuffix(header.Filename, ".png") {
		path = filepath.Join(path, fmt.Sprintf("%v/%v.png", user.Username, item.StopTime.UnixNano()))
	} else if strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, "jpeg") {
		path = filepath.Join(path, fmt.Sprintf("%v/%v.jpg", user.Username, item.StopTime.UnixNano()))
	}

	outputFile, err := os.Create(path)
	utils.Must(err)
	defer outputFile.Close()

	io.Copy(outputFile, imageFile)

	path = path[len(utils.RterDir):]

	item.ContentURI = path
	item.ThumbnailURI = path

	if exists {
		storage.Update(item)
	} else {
		storage.Insert(item)
	}

	userDirection := new(data.UserDirection)
	userDirection.Username = user.Username

	err = storage.Select(userDirection)

	if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(strconv.FormatFloat(userDirection.Heading, 'f', 6, 64)))
}
