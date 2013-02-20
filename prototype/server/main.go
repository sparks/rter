package main

import (
	"./server"
	"log"
	"net/http"
)

func main() {
	server.SetupMySQL()

	// http.HandleFunc("/upload", server.UploadHandler)
	http.HandleFunc("/multiup", server.MultiUploadHandler)
	http.HandleFunc("/submit", server.SubmitHandler)

	http.HandleFunc("/ajax/", server.ClientAjax)

	http.HandleFunc("/", server.ClientHandler)

	http.HandleFunc("/images/", server.ImageHandler)
	http.HandleFunc("/resources/", server.ResourceHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
