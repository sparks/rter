package main

import (
	"./server"
	"net/http"
)

func main() {
	server.SetupMySQL()

	http.HandleFunc("/upload", server.UploadHandler)
	http.HandleFunc("/multiup", server.MultiUploadHandler)
	// http.HandleFunc("/nehil", server.Nehil)

	http.HandleFunc("/ajax/", server.ClientAjax)

	http.HandleFunc("/", server.ClientHandler)

	http.HandleFunc("/images/", server.ImageHandler)
	http.HandleFunc("/resources/", server.ResourceHandler)

	http.ListenAndServe(":8080", nil)
}
