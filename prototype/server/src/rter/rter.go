package main

import (
	"log"
	"net/http"
	"os"
	"rter/server"
)

func main() {
	setupLogger()

	server.SetupMySQL()
	defer server.CloseMySQL()

	http.HandleFunc("/multiup", server.MultiUploadHandler)
	http.HandleFunc("/submit", server.SubmitHandler)

	http.HandleFunc("/ajax/", server.ClientAjax)

	http.HandleFunc("/", server.ClientHandler)

	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir(server.ImagePath))))
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
