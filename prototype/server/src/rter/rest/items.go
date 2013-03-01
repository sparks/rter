package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterItems(r *mux.Router) {
	itemsRouter := r.PathPrefix("/items").Subrouter()

	itemsRouter.HandleFunc("/", QueryItems).Methods("GET")
	itemsRouter.HandleFunc("/", CreateItem).Methods("POST")

	itemsRouter.HandleFunc("/{id:[0-9]+}", GetItem).Methods("GET")
	itemsRouter.HandleFunc("/{id:[0-9]+}", UpdateItem).Methods("POST")
	itemsRouter.HandleFunc("/{id:[0-9]+}", DeleteItem).Methods("DELETE")

	itemsRouter.HandleFunc("/{id:[0-9]+}/comments", QueryComments).Methods("GET")
	itemsRouter.HandleFunc("/{id:[0-9]+}/comments", CreateComment).Methods("POST")

	itemsRouter.HandleFunc("/{id:[0-9]+}/comments/{id:[0-9]+}", GetComment).Methods("GET")
	itemsRouter.HandleFunc("/{id:[0-9]+}/comments/{id:[0-9]+}", UpdateComment).Methods("POST")
	itemsRouter.HandleFunc("/{id:[0-9]+}/comments/{id:[0-9]+}", DeleteComment).Methods("DELETE")
}

func QueryItems(w http.ResponseWriter, r *http.Request) {

}

func CreateItem(w http.ResponseWriter, r *http.Request) {

}

func GetItem(w http.ResponseWriter, r *http.Request) {

}

func UpdateItem(w http.ResponseWriter, r *http.Request) {

}

func DeleteItem(w http.ResponseWriter, r *http.Request) {

}

func QueryComments(w http.ResponseWriter, r *http.Request) {

}

func CreateComment(w http.ResponseWriter, r *http.Request) {

}

func GetComment(w http.ResponseWriter, r *http.Request) {

}

func UpdateComment(w http.ResponseWriter, r *http.Request) {

}

func DeleteComment(w http.ResponseWriter, r *http.Request) {

}
