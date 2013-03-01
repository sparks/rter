package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterTaxonomy(r *mux.Router) {
	taxonomyRouter := r.PathPrefix("/taxonomy").Subrouter()

	taxonomyRouter.HandleFunc("/", QueryTerm).Methods("GET")
	taxonomyRouter.HandleFunc("/", CreateTerm).Methods("POST")

	taxonomyRouter.HandleFunc("/{id:[0-9]+}/", GetTerm).Methods("GET")
	taxonomyRouter.HandleFunc("/{id:[0-9]+}/", UpdateTerm).Methods("POST")
	taxonomyRouter.HandleFunc("/{id:[0-9]+}/", DeleteTerm).Methods("DELTE")

	taxonomyRouter.HandleFunc("/{id:[0-9]+}/ranking", GetTermRanking).Methods("GET")
	taxonomyRouter.HandleFunc("/{id:[0-9]+}/ranking", UpdateTermRanking).Methods("POST")
}

func QueryTerm(w http.ResponseWriter, r *http.Request) {

}

func CreateTerm(w http.ResponseWriter, r *http.Request) {

}

func GetTerm(w http.ResponseWriter, r *http.Request) {

}

func UpdateTerm(w http.ResponseWriter, r *http.Request) {

}

func DeleteTerm(w http.ResponseWriter, r *http.Request) {

}

func GetTermRanking(w http.ResponseWriter, r *http.Request) {

}

func UpdateTermRanking(w http.ResponseWriter, r *http.Request) {

}
