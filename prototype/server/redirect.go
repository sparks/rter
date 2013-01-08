package main

import (
	"net/http"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.cim.mcgill.ca/sre/projects/rter/", http.StatusTemporaryRedirect)
}

func main() {
	http.HandleFunc("/", redirect)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}
