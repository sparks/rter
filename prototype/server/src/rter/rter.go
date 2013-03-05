package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"rter/mobile"
	"rter/rest"
	"rter/storage"
	"rter/utils"
	"rter/web"
)

func main() {
	setupLogger()

	storage.OpenStorage("root", "", "tcp", "localhost:3306", "rter_test")
	defer storage.CloseStorage()

	r := mux.NewRouter()

	rest.RegisterCRUD(r)

	r.HandleFunc("/multiup", mobile.MultiUploadHandler)
	r.HandleFunc("/submit", web.SubmitHandler)

	r.HandleFunc("/ajax/", web.ClientAjax)

	r.HandleFunc("/", web.ClientHandler)

	r.Handle("/uploads/", http.StripPrefix("/uploads", http.FileServer(http.Dir(utils.UploadPath))))
	r.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir(utils.ResourcePath))))

	http.Handle("/", r)

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
