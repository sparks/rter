package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	ioutil.WriteFile(imagePath+"test.png", p, 0600)
}

func MultiUploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	checkError(err)

	ins, err := db.Prepare("INSERT INTO content (phone_id, filepath, geolat, geolong) VALUES(?, ?, ?, ?)")
	checkError(err)

	path := imagePath + header.Filename
	fo, err := os.Create(path)
	checkError(err)
	defer fo.Close()

	io.Copy(fo, file)

	fakeLat := 45.129848
	fakeLong := 40.357694
	_, err = ins.Run([]byte("look_a_phone"), []byte(path), fakeLat, fakeLong)
	checkError(err)

	queryDB()
}
