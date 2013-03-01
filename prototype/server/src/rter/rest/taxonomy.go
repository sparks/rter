package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterTaxonomy(r *mux.Router) {
	taxonomyRouter := r.PathPrefix("/taxonomy").Subrouter()

	taxonomyRouter.HandleFunc("/", QueryTaxonomy).Methods("GET")
	taxonomyRouter.HandleFunc("/", CreateTaxonomyTerm).Methods("POST")

	taxonomyRouter.HandleFunc("/{id:[0-9]+}/", GetTaxonomyTerm).Methods("GET")
	taxonomyRouter.HandleFunc("/{id:[0-9]+}/", UpdateTaxonomyTerm).Methods("POST")
	taxonomyRouter.HandleFunc("/{id:[0-9]+}/", DeleteTaxonomyTerm).Methods("DELTE")

	taxonomyRouter.HandleFunc("/{id:[0-9]+}/ranking", GetTaxonomyTermRanking).Methods("GET")
	taxonomyRouter.HandleFunc("/{id:[0-9]+}/ranking", UpdateTaxonomyTermRanking).Methods("POST")
}

func QueryTaxonomy(w http.ResponseWriter, r *http.Request) {

}

func CreateTaxonomyTerm(w http.ResponseWriter, r *http.Request) {

}

func GetTaxonomyTerm(w http.ResponseWriter, r *http.Request) {

}

func UpdateTaxonomyTerm(w http.ResponseWriter, r *http.Request) {

}

func DeleteTaxonomyTerm(w http.ResponseWriter, r *http.Request) {

}

func GetTaxonomyTermRanking(w http.ResponseWriter, r *http.Request) {

}

func UpdateTaxonomyTermRanking(w http.ResponseWriter, r *http.Request) {

}
