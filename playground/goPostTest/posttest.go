package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.Println("Launching POST request Debugger")

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", POSTHandler).Methods("POST")
	r.HandleFunc("/", OtherHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":1234", nil))
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New POST request")
	fmt.Println("==================================")

	fmt.Println("Headers:")
	for key, value := range r.Header {
		fmt.Println("\t", key, "->", value)
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Got error reading body", err)
	}

	fmt.Println("Body:")

	fmt.Println("\tBody length (in bytes) is:", len(body))
	fmt.Println("\tBody is:")
	fmt.Println("\t", string(body))
	fmt.Println("==================================")

	w.WriteHeader(http.StatusOK)
}

func OtherHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got non POST Request, rejecting")
	fmt.Println("==================================")

	w.WriteHeader(http.StatusNotFound)

	response := "Sorry I only accept POST requests, you sent a " + r.Method + " request."

	w.Write([]byte(response))
}
