package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

	insert, error := database.Prepare("INSERT INTO content (phone_id, filepath, geolat, geolong) VALUES(?, ?, ?, ?)")
	checkError(error)

	path := imagePath + header.Filename
	outputFile, error := os.Create(path)
	checkError(error)
	defer outputFile.Close()

	io.Copy(outputFile, imageFile)

	fakeLat := 45.129848
	fakeLong := 40.357694
	_, error = insert.Run([]byte("look_a_phone"), []byte(path), fakeLat, fakeLong)
	checkError(error)

	queryDatabase()
}
