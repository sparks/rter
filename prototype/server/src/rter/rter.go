package main

import (
	"log"
	"net/http"
	"os"
	"rter/server"
)

var rterDir = os.Getenv("RTER_DIR")

var UploadPath = filepath.Join(rterDir, "uploads")
var TemplatePath = filepath.Join(rterDir, "templates")
var ResourcePath = filepath.Join(rterDir, "resources")

var phoneIDValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

var filenameValidator = regexp.MustCompile("^[a-zA-Z0-9_]*\\.?[a-zA-Z0-9_]+\\.[a-zA-Z0-9]+$")
var folderNameValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func main() {
	setupLogger()

	server.SetupMySQL()
	defer server.CloseMySQL()

	http.HandleFunc("/multiup", server.MultiUploadHandler)
	http.HandleFunc("/submit", server.SubmitHandler)

	http.HandleFunc("/ajax/", server.ClientAjax)

	http.HandleFunc("/", server.ClientHandler)

	http.Handle("/uploads/", http.StripPrefix("/uploads", http.FileServer(http.Dir(server.UploadPath))))
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir(server.ResourcePath))))

	log.Println("Launching rtER Server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupLogger() {
	logOutputFile := os.Getenv("RTER_LOGFILE")

	if logOutputFile != "" {
		logFile, err := os.Create(logOutputFile)

		if err == nil {
			log.SetOutput(logFile)
		} else {
			log.Println(err)
		}
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
