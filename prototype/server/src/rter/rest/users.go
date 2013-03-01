package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterUsers(r *mux.Router) {
	usersRouter := r.PathPrefix("/users").Subrouter()

	usersRouter.HandleFunc("/", QueryUsers).Methods("GET")
	usersRouter.HandleFunc("/", CreateUser).Methods("POST")

	usersRouter.HandleFunc("/{id:[0-9]+}", GetUser).Methods("GET")
	usersRouter.HandleFunc("/{id:[0-9]+}", UpdateUser).Methods("POST")
	usersRouter.HandleFunc("/{id:[0-9]+}", DeleteUser).Methods("DELETE")

	usersRouter.HandleFunc("/{id:[0-9]+}/direction", GetUserDirection).Methods("GET")
	usersRouter.HandleFunc("/{id:[0-9]+}/direction", UpdateUserDirection).Methods("POST")

	// rolesRouter := r.PathPrefix("/roles").Subrouter()

	// rolesRouter.HandleFunc("/", QueryRoles).Methods("GET")
	// rolesRouter.HandleFunc("/", QueryRoles).Methods("POST")

	// rolesRouter.HandleFunc("/{role}").Methods("GET")
	// rolesRouter.HandleFunc("/{role}").Methods("POST")
	// rolesRouter.HandleFunc("/{role}").Methods("DELETE")
}

func QueryUsers(w http.ResponseWriter, r *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

}

func GetUser(w http.ResponseWriter, r *http.Request) {

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func GetUserDirection(w http.ResponseWriter, r *http.Request) {

}

func UpdateUserDirection(w http.ResponseWriter, r *http.Request) {

}
