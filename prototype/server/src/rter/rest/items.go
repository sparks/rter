package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rter/data"
	"rter/storage"
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
	decoder := json.NewDecoder(r.Body)

	item := new(data.Item)
	err := decoder.Decode(&item)

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed json for Item", http.StatusBadRequest)
		return
	}

	err = storage.InsertItem(item)

	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)

	encoder.Encode(item)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	item := new(data.Item)
	err := decoder.Decode(&item)

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed json for Item", http.StatusBadRequest)
		return
	}

	err = storage.GetItem(item)

	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	encoder := json.NewEncoder(w)

	encoder.Encode(item)
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
