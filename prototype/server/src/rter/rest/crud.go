// Provide a RESTful CRUD API.
//
// Create, Read, Update, Delete actions can be performed (where appropriate) on all the core data structures given in rter/data. The package creates a router via CRUDRouter() which can be attached to any http prefix
package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net"
	"net/http"
	"rter/auth"
	"rter/data"
	"rter/storage"
	"strconv"
	"strings"
	"time"
	token "videoserver/auth"
)

var decoder = schema.NewDecoder()

// Generate a new CRUD router for RESTful access to the rtER datastructures. Includes support for OPTIONS Method to check what functionality is available
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

	return r
}

// Respond to requests with the OPTIONS Method.
func StateOptions(opts string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-PINGOTHER")
		w.Header().Set("Access-Control-Allow-Methods", opts)

		w.WriteHeader(http.StatusOK)
	}
}

// Generic Create handler
func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, permissions := auth.Challenge(w, r, true)

	if (user == nil || permissions < 1) && vars["datatype"] != "users" { // Allow anyone to create users for now
		http.Error(w, "Please Login", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-PINGOTHER")

	var val interface{} // Generic container for the new object

	// Build a URI like representation of the datatype
	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	// Switch based on that URI like representation and instantiate something in the generic container
	switch strings.Join(types, "/") {
	case "items":
		val = new(data.Item)
	case "items/comments":
		val = new(data.ItemComment)
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

	// Perform the JSON decode
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&val)

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed json.", http.StatusBadRequest)
		return
	}

	// Perform post decode actions, setting automated field, validate values, exectute hooks, etc ...
	switch v := val.(type) {
	case *data.Item:
		v.Author = user.Username
	case *data.ItemComment:
		v.ItemID, err = strconv.ParseInt(vars["key"], 10, 64)
		v.Author = user.Username
	case *data.User:
		v.HashAndSalt()
		v.Role = "public" // TODO: Temporary while anyone can sign up maybe this will change?
	case *data.Term:
		v.Author = user.Username
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	// Perform the DB insert
	err = storage.Insert(val)

	if err != nil {
		log.Println(err)
		http.Error(w, "Insert Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	// Exectute post insert hooks, etc ...
	switch v := val.(type) {
	case *data.Item:
		if v.Type == "streaming-video-v1" {
			v.UploadURI = "http://192.168.30.8:8080/v1/ingest/" + strconv.FormatInt(v.ID, 10)
			v.ThumbnailURI = "http://192.168.30.8:8080/v1/videos/" + strconv.FormatInt(v.ID, 10) + "/thumb/000000001.jpg"
			v.ContentURI = "http://192.168.30.8:8080/v1/videos/" + strconv.FormatInt(v.ID, 10)

			host, _, err := net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				log.Println(err)
				http.Error(w, "Problem building streaming tokens, not remote addresse available.", http.StatusBadRequest)
				return
			}

			t, err := token.GenerateToken(v.UploadURI, host, time.Duration(3600)*time.Second, "1122AABBCCDDEEFF")

			if err != nil {
				log.Println(err)
				http.Error(w, "Problem building streaming tokens, likely due to malformed request.", http.StatusInternalServerError)
				return
			}

			v.Token = t

			err = storage.Update(v) //FIXME: This is awful, but probably not workaroundable?

			if err != nil {
				log.Println(err)
				http.Error(w, "Update Database error, likely due to malformed request.", http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json") // Header are important when GZIP is enabled
	w.WriteHeader(http.StatusCreated)

	// Return the object we've inserted in the database.
	encoder := json.NewEncoder(w)
	err = encoder.Encode(val)

	if err != nil {
		log.Println(err)
	}
}

// Generic Read handler for reading single objects
func Read(w http.ResponseWriter, r *http.Request) {
	//No Auth for the moment

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)

	var (
		val interface{} // Generic container for the read object
		err error
	)

	// Build a URI like representation of the datatype
	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	// Switch based on that URI like representation and instantiate something in the generic container. Also infer the identifier from the vars and perform validation.
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

	// Perform the Select
	err = storage.Select(val)

	if err == storage.ErrZeroAffected {
		http.Error(w, "No matches for query", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Select Database error, likely due to malformed request", http.StatusInternalServerError)
		return
	}

	// Important. Let's never send salt/hash out.
	switch v := val.(type) {
	case *data.User:
		v.Salt = ""
		v.Password = ""
	}

	w.Header().Set("Content-Type", "application/json") // Header are important when GZIP is enabled

	// Return the object we've selected from the database.
	encoder := json.NewEncoder(w)
	err = encoder.Encode(val)

	if err != nil {
		log.Println(err)
	}
}

// Generic Read handler for reading multiple objects possibly with a query
func ReadWhere(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)

	var (
		val interface{} // Generic container for the read objecs
		err error
	)

	whereClause := ""              //We may build up some sort of WHERE clause for the DB request
	args := make([]interface{}, 0) //This WHERE clause would then require some arguments

	// Build a URI like representation of the datatype
	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	// Switch based on that URI like representation and instantiate something in the generic container. Also infer the identifier from the vars and perform validation.
	switch strings.Join(types, "/") {
	case "items":
		items := make([]*data.Item, 0)
		val = &items
	case "items/comments":
		comments := make([]*data.ItemComment, 0)

		// Selecting comments is reliant on the item ID
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

	// Sadly due to the need to load Term.Count value we must have some custom queries here
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
		http.Error(w, "Select2 Database error, likely due to malformed request", http.StatusInternalServerError)
		return
	}

	// Important. Let's never send salt/hash out
	switch v := val.(type) {
	case *[]*data.User:
		for _, user := range *v {
			user.Salt = ""
			user.Password = ""
		}
	}

	w.Header().Set("Content-Type", "application/json") // Header are important when GZIP is enabled

	// Return the selected objects
	encoder := json.NewEncoder(w)
	err = encoder.Encode(val)

	if err != nil {
		log.Println(err)
	}
}

// Generic Update handler
func Update(w http.ResponseWriter, r *http.Request) {
	user, permissions := auth.Challenge(w, r, true)
	if user == nil || permissions < 1 {
		http.Error(w, "Please Login", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)
	var (
		val interface{} // Generic container for the updated object
		err error
	)

	// Build a URI like representation of the datatype
	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	// Switch based on that URI like representation and instantiate something in the generic container. Also infer the identifier from the vars and perform validation.
	switch strings.Join(types, "/") {
	case "items":
		v := new(data.Item)
		v.ID, err = strconv.ParseInt(vars["key"], 10, 64)

		val = v
	case "items/comments":
		v := new(data.ItemComment)
		v.ID, err = strconv.ParseInt(vars["childkey"], 10, 64)

		val = v
	case "users":
		if vars["key"] != user.Username {
			http.Error(w, "Please don't hack other users", http.StatusUnauthorized)
			return
		}

		v := new(data.User)
		v.Username = vars["key"]

		val = v
	case "users/direction":
		v := new(data.UserDirection)
		v.Username = vars["key"]

		val = v
	case "roles":
		v := new(data.Role)
		v.Title = vars["key"]

		val = v
	case "taxonomy":
		v := new(data.Term)
		v.Term = vars["key"]

		val = v
	case "taxonomy/ranking":
		v := new(data.TermRanking)
		v.Term = vars["key"]

		val = v
	default:
		http.NotFound(w, r)
		return
	}

	if err != nil {
		log.Println(err, vars)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	err = storage.Select(val) //Load previous values so that update is non distructive of empty fields

	if err == storage.ErrZeroAffected {
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Select3 Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	// Decode the JSON into our generic object. The decode will leave unscpecified fields untouched.
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&val)

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed json.", http.StatusBadRequest)
		return
	}

	// Validate JSON, run pre-update hooks, etc...
	//We must reset fields we set earlier incase they were changed during the JSON decode
	switch v := val.(type) {
	case (*data.Item):
		v.ID, err = strconv.ParseInt(vars["key"], 10, 64)
		v.Author = user.Username
	case (*data.ItemComment):
		v.ID, err = strconv.ParseInt(vars["childkey"], 10, 64)
		v.Author = user.Username
	case (*data.User):
		v.Username = vars["key"]
	case (*data.UserDirection):
		v.Username = vars["key"]
		v.LockUsername = user.Username
	case (*data.Role):
		v.Title = vars["key"]
	case (*data.Term):
		v.Term = vars["key"]
		v.Author = user.Username
	case (*data.TermRanking):
		v.Term = vars["key"]
	}

	if err != nil {
		log.Println(err, vars)
		http.Error(w, "Malformed key in URI", http.StatusBadRequest)
		return
	}

	// Run the update
	err = storage.Update(val)

	if err == storage.ErrZeroAffected {
		w.WriteHeader(http.StatusNotModified)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Update2 Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Header are important when GZIP is enabled

	// Return the updated item
	encoder := json.NewEncoder(w)
	err = encoder.Encode(val)

	if err != nil {
		log.Println(err)
	}
}

// Generic Delete handler
func Delete(w http.ResponseWriter, r *http.Request) {
	user, permissions := auth.Challenge(w, r, true)
	if user == nil || permissions < 1 {
		http.Error(w, "Please Login", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)

	var (
		val interface{} // Generic container for the deleted object
		err error
	)

	// Build a URI like representation of the datatype
	types := []string{vars["datatype"]}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	// Switch based on that URI like representation and instantiate something in the generic container. Also infer the identifier from the vars and perform validation.
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
		if vars["key"] != user.Username {
			http.Error(w, "Please don't delete other users", http.StatusUnauthorized)
			return
		}

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

	// Perform the delete
	err = storage.Delete(val)

	if err == storage.ErrZeroAffected {
		http.Error(w, "No matches for query", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Delete Database error, likely due to malformed request", http.StatusInternalServerError)
		return
	}

	// Confirm the delete
	w.WriteHeader(http.StatusNoContent)
}
