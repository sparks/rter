package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	checkError(error)

	phone_id := r.FormValue("phone_id")

	lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
	checkError(err)
	long, err := strconv.ParseFloat(r.FormValue("long"), 64)
	checkError(err)

	insert, error := database.Prepare("INSERT INTO content (phone_id, filepath, geolat, geolong) VALUES(?, ?, ?, ?)")
	checkError(error)

	path := imagePath + header.Filename

	outputFile, error := os.Create(path)
	checkError(error)
	defer outputFile.Close()

	io.Copy(outputFile, imageFile)

	_, error = insert.Run([]byte(phone_id), []byte(path), lat, long)
	checkError(error)

	queryDatabase()
}
