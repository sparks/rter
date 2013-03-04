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

	// itemsRouter.HandleFunc("/{id:[0-9]+}/comments/{id:[0-9]+}", GetComment).Methods("GET")
	itemsRouter.HandleFunc("/{id:[0-9]+}/comments/{id:[0-9]+}", UpdateComment).Methods("POST")
	itemsRouter.HandleFunc("/{id:[0-9]+}/comments/{id:[0-9]+}", DeleteComment).Methods("DELETE")
}

func QueryItems(w http.ResponseWriter, r *http.Request) {
	items, err := storage.SelectAllItems()

	if err == storage.ErrZeroMatches {
		http.Error(w, "No Items", http.StatusNoContent)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(items)

	if err != nil {
		log.Println(err)
	}
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

	err = storage.Insert(item)

	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(item)

	if err != nil {
		log.Println(err)
	}
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	item := new(data.Item)
	item.ID = parseID(2, w, r)

	err := storage.SelectItem(item)

	if err == storage.ErrZeroMatches {
		http.Error(w, "No such Item", http.StatusNoContent)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(item)

	if err != nil {
		log.Println(err)
	}
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	item := new(data.Item)
	err := decoder.Decode(&item)

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed json for Item", http.StatusBadRequest)
		return
	}

	item.ID = parseID(2, w, r)

	err = storage.UpdateItem(item)

	if err == storage.ErrZeroMatches {
		w.WriteHeader(http.StatusNotModified)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	item := new(data.Item)
	item.ID = parseID(2, w, r)

	err := storage.Delete(item)

	if err == storage.ErrZeroMatches {
		http.Error(w, "No such Item", http.StatusNoContent)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func QueryComments(w http.ResponseWriter, r *http.Request) {

}

func CreateComment(w http.ResponseWriter, r *http.Request) {

}

func UpdateComment(w http.ResponseWriter, r *http.Request) {

}

func DeleteComment(w http.ResponseWriter, r *http.Request) {

}
