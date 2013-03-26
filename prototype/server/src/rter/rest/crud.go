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

func CRUDRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}", StateOptions("GET, POST")).Methods("OPTIONS")

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}", ReadWhere).Methods("GET")
	r.HandleFunc("/{datatype:items|users|roles|taxonomy}", Create).Methods("POST")

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}", StateOptions("GET, PUT, DELETE")).Methods("OPTIONS")

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}", Read).Methods("GET")
	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}", Update).Methods("PUT")
	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}", Delete).Methods("DELETE")

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:ranking|direction}", StateOptions("GET, PUT")).Methods("OPTIONS")

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:ranking|direction}", Read).Methods("GET")
	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:ranking|direction}", Update).Methods("PUT")

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:comments}", StateOptions("GET, POST")).Methods("OPTIONS")

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:comments}", ReadWhere).Methods("GET")
	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:comments}", Create).Methods("POST")

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:comments}/{childkey}", StateOptions("GET, PUT, DELETE")).Methods("OPTIONS")

	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:comments}/{childkey}", Read).Methods("GET")
	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:comments}/{childkey}", Update).Methods("PUT")
	r.HandleFunc("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:comments}/{childkey}", Delete).Methods("DELETE")

	r.PathPrefix("/").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			probe("Unkown CRUD request", r)
		},
	)

	return r
}

func StateOptions(opts string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		probe("Options Request", r)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-PINGOTHER")
		w.Header().Set("Access-Control-Allow-Methods", opts)

		w.WriteHeader(http.StatusOK)
	}
}

func probe(message string, r *http.Request) {
	// log.Println(message)
	// log.Println(r.Method, r.URL)
	// e, _ := json.Marshal(r)
	// log.Println(string(e))
}

func Create(w http.ResponseWriter, r *http.Request) {
	probe("Create Request", r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-PINGOTHER")

	vars := mux.Vars(r)

	var val interface{}

	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	switch strings.Join(types, "/") {
	case "items":
		val = new(data.Item)
	case "items/comments":
		val = new(data.ItemComment)
	case "users":
		user := new(data.User)
		user.HashAndSalt()

		val = user
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
	case *data.ItemComment:
		v.ItemID, err = strconv.ParseInt(vars["key"], 10, 64)
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
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
	probe("Read Request", r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)

	var (
		val interface{}
		err error
	)

	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	switch strings.Join(types, "/") {
	case "items":
		item := new(data.Item)
		item.ID, err = strconv.ParseInt(vars["key"], 10, 64)

		val = item
	case "items/comments":
		comment := new(data.ItemComment)
		comment.ID, err = strconv.ParseInt(vars["childkey"], 10, 64)

		val = comment
	case "users":
		user := new(data.User)
		user.Username = vars["key"]

		val = user
	case "users/direction":
		direction := new(data.UserDirection)
		direction.Username = vars["key"]

		val = direction
	case "roles":
		role := new(data.Role)
		role.Title = vars["key"]

		val = role
	case "taxonomy":
		term := new(data.Term)
		term.Term = vars["key"]

		val = term
	case "taxonomy/ranking":
		ranking := new(data.TermRanking)
		ranking.Term = vars["key"]

		val = ranking
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

	if err == storage.ErrZeroAffected {
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

func ReadWhere(w http.ResponseWriter, r *http.Request) {
	probe("ReadWhere Request", r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)

	var (
		val interface{}
		err error
	)

	whereClause := ""
	args := make([]interface{}, 0)

	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	switch strings.Join(types, "/") {
	case "items":
		items := make([]*data.Item, 0)
		val = &items
	case "items/comments":
		comments := make([]*data.ItemComment, 0)

		whereClause = "WHERE ItemID=?"
		var itemID int64
		itemID, err = strconv.ParseInt(vars["key"], 10, 64)
		args = append(args, itemID)

		val = &comments
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

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	switch val.(type) {
	case *[]*data.Term:
		err = storage.SelectQuery(val, "SELECT t.*, count(r.Term) FROM Terms AS t, TermRelationships AS r WHERE t.Term = r.Term GROUP BY t.Term")
	default:
		err = storage.SelectWhere(val, whereClause, args...)
	}

	if err == storage.ErrZeroAffected {
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
	probe("Update Request", r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)
	var val interface{}

	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	switch strings.Join(types, "/") {
	case "items":
		val = new(data.Item)
	case "items/comments":
		val = new(data.ItemComment)
	case "users":
		val = new(data.User)
	case "users/direction":
		val = new(data.UserDirection)
	case "roles":
		val = new(data.Role)
	case "taxonomy":
		val = new(data.Term)
	case "taxonomy/ranking":
		val = new(data.TermRanking)
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
	case (*data.ItemComment):
		v.ID, err = strconv.ParseInt(vars["childkey"], 10, 64)
	case (*data.User):
		v.Username = vars["key"]
	case (*data.UserDirection):
		v.Username = vars["key"]
	case (*data.Role):
		v.Title = vars["key"]
	case (*data.Term):
		v.Term = vars["key"]
	case (*data.TermRanking):
		v.Term = vars["key"]
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	err = storage.Update(val)

	if err == storage.ErrZeroAffected {
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
	probe("Delete Request", r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)

	var (
		val interface{}
		err error
	)

	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	switch strings.Join(types, "/") {
	case "items":
		item := new(data.Item)
		item.ID, err = strconv.ParseInt(vars["key"], 10, 64)

		val = item
	case "items/comments":
		comment := new(data.ItemComment)
		comment.ID, err = strconv.ParseInt(vars["childkey"], 10, 64)

		val = comment
	case "users":
		user := new(data.User)
		user.Username = vars["key"]

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

	if err == storage.ErrZeroAffected {
		http.Error(w, "No matches for query", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
