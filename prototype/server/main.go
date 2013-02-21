package main

import (
	"./server"
	"log"
	"net/http"
	"os"
)

func main() {

	logOutputFile := os.Getenv("RTER_LOGFILE")

	if logOutputFile != "" {
		logFile, err := os.Create(logOutputFile)
		if err == nil {
			log.SetOutput(logFile)
		} else {
			log.Println(err)
		}
	}

	log.Println("hello all")

	server.SetupMySQL()

	http.HandleFunc("/multiup", server.MultiUploadHandler)
	http.HandleFunc("/submit", server.SubmitHandler)

	http.HandleFunc("/ajax/", server.ClientAjax)

	http.HandleFunc("/", server.ClientHandler)

	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir(server.ImagePath))))
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir(server.ResourcePath))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
