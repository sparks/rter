package main

import (
	"log"
	"net/http"
	"os"
	"rter/mobile"
	"rter/storage"
	"rter/util"
	"rter/web"
)

func main() {
	setupLogger()

	storage.OpenStorage("root", "", "tcp", "localhost:3306", "rter_v2")
	defer storage.CloseStorage()

	http.HandleFunc("/multiup", mobile.MultiUploadHandler)
	http.HandleFunc("/submit", web.SubmitHandler)

	http.HandleFunc("/ajax/", web.ClientAjax)

	http.HandleFunc("/", web.ClientHandler)

	http.Handle("/uploads/", http.StripPrefix("/uploads", http.FileServer(http.Dir(util.UploadPath))))
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir(util.ResourcePath))))

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
