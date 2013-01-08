package main

import (
	"./server"
	"net/http"
)

func main() {
	// Resume normal operation
	server.SetupMySQL()

	http.HandleFunc("/images/", server.ImageHandler)

	http.HandleFunc("/upload", server.UploadHandler)
	http.HandleFunc("/multiup", server.MultiUploadHandler)

	http.HandleFunc("/view/", server.MakeHandler(server.ViewHandler))
	http.HandleFunc("/edit/", server.MakeHandler(server.EditHandler))
	http.HandleFunc("/save/", server.MakeHandler(server.SaveHandler))

	http.HandleFunc("/", server.MakeHandler(server.MainHandler))

	http.ListenAndServe(":8080", nil)
}
