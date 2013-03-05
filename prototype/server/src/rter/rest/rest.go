package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"rter/data"
	"rter/storage"
	"strconv"
	"strings"
)

var decoder = schema.NewDecoder()

func RegisterRest(r *mux.Router) {
	r.HandleFunc("/{datatype}", Create).Methods("POST")

	r.HandleFunc("/{datatype}", ReadAll).Methods("GET")
	r.HandleFunc("/{datatype}/{key}", Read).Methods("GET")

	r.HandleFunc("/{datatype}/{key}", Update).Methods("PUT")

	r.HandleFunc("/{datatype}/{key}", Delete).Methods("DELETE")

	//TODO: /user/direction  /item/comments  /taxonomy/ranking
}

func Create(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.Split(r.URL.Path, "/")
	var val interface{}

	switch splitPath[1] {
	case "items":
		val = new(data.Item)
	case "users":
		val = new(data.User)
	case "roles":
		val = new(data.Role)
	case "taxonomy":
		val = new(data.Term)
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&val)

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed json.", http.StatusBadRequest)
		return
	}

	err = storage.Insert(val)

	if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(val)

	if err != nil {
		log.Println(err)
	}
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.Split(r.URL.Path, "/")

	var val interface{}

	switch splitPath[1] {
	case "items":
		items := make([]*data.Item, 0)
		val = &items
	case "users":
		users := make([]*data.User, 0)
		val = &users
	case "roles":
		roles := make([]*data.Role, 0)
		val = &roles
	case "taxonomy":
		terms := make([]*data.Term, 0)
		val = &terms
	}

	err := storage.SelectAll(val)

	if err == storage.ErrZeroMatches {
		http.Error(w, "No matches for query", http.StatusNoContent)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(val)

	if err != nil {
		log.Println(err)
	}
}

func Read(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.Split(r.URL.Path, "/")

	var (
		val interface{}
		err error
	)

	switch splitPath[1] {
	case "items":
		item := new(data.Item)
		item.ID, err = strconv.ParseInt(splitPath[2], 10, 64)

		val = item
	case "users":
		user := new(data.User)
		user.ID, err = strconv.ParseInt(splitPath[2], 10, 64)

		val = user
	case "roles":
		role := new(data.Role)
		role.Title = splitPath[2]

		val = role
	case "taxonomy":
		term := new(data.Term)
		term.Term = splitPath[2]

		val = term
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	err = storage.Select(val)

	if err == storage.ErrZeroMatches {
		http.Error(w, "No matches for query", http.StatusNoContent)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(val)

	if err != nil {
		log.Println(err)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.Split(r.URL.Path, "/")
	var val interface{}

	switch splitPath[1] {
	case "items":
		val = new(data.Item)
	case "users":
		val = new(data.User)
	case "roles":
		val = new(data.Role)
	case "taxonomy":
		val = new(data.Term)
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&val)

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed json.", http.StatusBadRequest)
		return
	}

	switch v := val.(type) {
	case (*data.Item):
		v.ID, err = strconv.ParseInt(splitPath[2], 10, 64)
	case (*data.User):
		v.ID, err = strconv.ParseInt(splitPath[2], 10, 64)
	case (*data.Role):
		v.Title = splitPath[2]
	case (*data.Term):
		v.Term = splitPath[2]
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	err = storage.Update(val)

	if err == storage.ErrZeroMatches {
		w.WriteHeader(http.StatusNotModified)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(val)

	if err != nil {
		log.Println(err)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.Split(r.URL.Path, "/")

	var (
		val interface{}
		err error
	)

	switch splitPath[1] {
	case "items":
		item := new(data.Item)
		item.ID, err = strconv.ParseInt(splitPath[2], 10, 64)

		val = item
	case "users":
		user := new(data.User)
		user.ID, err = strconv.ParseInt(splitPath[2], 10, 64)

		val = user
	case "roles":
		role := new(data.Role)
		role.Title = splitPath[2]

		val = role
	case "taxonomy":
		term := new(data.Term)
		term.Term = splitPath[2]

		val = term
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	err = storage.Delete(val)

	if err == storage.ErrZeroMatches {
		http.Error(w, "No matches for query", http.StatusNoContent)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
