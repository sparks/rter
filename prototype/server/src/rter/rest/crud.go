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
)

var decoder = schema.NewDecoder()

func RegisterCRUD(r *mux.Router) {
	r.HandleFunc("/{datatype}", Create).Methods("POST")

	r.HandleFunc("/{datatype}", ReadAll).Methods("GET")
	r.HandleFunc("/{datatype}/{key}", Read).Methods("GET")

	r.HandleFunc("/{datatype}/{key}", Update).Methods("PUT")

	r.HandleFunc("/{datatype}/{key}", Delete).Methods("DELETE")

	//TODO: /user/direction  /item/comments  /taxonomy/ranking
}

func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var val interface{}

	switch vars["datatype"] {
	case "items":
		val = new(data.Item)
	case "users":
		val = new(data.User)
	case "roles":
		val = new(data.Role)
	case "taxonomy":
		val = new(data.Term)
	default:
		http.NotFound(w, r)
		return
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

func Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var (
		val interface{}
		err error
	)

	switch vars["datatype"] {
	case "items":
		item := new(data.Item)
		item.ID, err = strconv.ParseInt(vars["key"], 10, 64)

		val = item
	case "users":
		user := new(data.User)
		user.ID, err = strconv.ParseInt(vars["key"], 10, 64)

		val = user
	case "roles":
		role := new(data.Role)
		role.Title = vars["key"]

		val = role
	case "taxonomy":
		term := new(data.Term)
		term.Term = vars["key"]

		val = term
	default:
		http.NotFound(w, r)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	err = storage.Select(val)

	if err == storage.ErrZeroMatches {
		http.Error(w, "No matches for query", http.StatusNotFound)
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

func ReadAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var val interface{}

	switch vars["datatype"] {
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
	default:
		http.NotFound(w, r)
		return

	}

	err := storage.SelectAll(val)

	if err == storage.ErrZeroMatches {
		http.Error(w, "No matches for query", http.StatusNotFound)
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
	vars := mux.Vars(r)
	var val interface{}

	switch vars["datatype"] {
	case "items":
		val = new(data.Item)
	case "users":
		val = new(data.User)
	case "roles":
		val = new(data.Role)
	case "taxonomy":
		val = new(data.Term)
	default:
		http.NotFound(w, r)
		return
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
		v.ID, err = strconv.ParseInt(vars["key"], 10, 64)
	case (*data.User):
		v.ID, err = strconv.ParseInt(vars["key"], 10, 64)
	case (*data.Role):
		v.Title = vars["key"]
	case (*data.Term):
		v.Term = vars["key"]
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
	vars := mux.Vars(r)

	var (
		val interface{}
		err error
	)

	switch vars["datatype"] {
	case "items":
		item := new(data.Item)
		item.ID, err = strconv.ParseInt(vars["key"], 10, 64)

		val = item
	case "users":
		user := new(data.User)
		user.ID, err = strconv.ParseInt(vars["key"], 10, 64)

		val = user
	case "roles":
		role := new(data.Role)
		role.Title = vars["key"]

		val = role
	case "taxonomy":
		term := new(data.Term)
		term.Term = vars["key"]

		val = term
	default:
		http.NotFound(w, r)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	err = storage.Delete(val)

	if err == storage.ErrZeroMatches {
		http.Error(w, "No matches for query", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
