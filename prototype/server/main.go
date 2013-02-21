package main

import (
	"./server"
	"log"
	"net/http"
)

func main() {
	server.SetupMySQL()

	http.HandleFunc("/multiup", server.MultiUploadHandler)
	http.HandleFunc("/submit", server.SubmitHandler)

	http.HandleFunc("/ajax/", server.ClientAjax)

	http.HandleFunc("/", server.ClientHandler)

	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir(server.ImagePath))))
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir(server.ResourcePath))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
