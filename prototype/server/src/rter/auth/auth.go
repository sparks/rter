package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"rter/data"
	"rter/storage"
)

func AuthHandlerFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Auth")
	loginInfo := new(data.User)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginInfo)

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed json.", http.StatusBadRequest)
		return
	}

	user := new(data.User)
	user.Username = loginInfo.Username

	err = storage.Select(user)

	if err == storage.ErrZeroAffected {
		log.Println("No such user")
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	if !user.Auth(loginInfo.Password) {
		log.Println("Wrong Password")
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You're cool"))
}
