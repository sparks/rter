package auth

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"rter/data"
	"rter/storage"
)

var store = sessions.NewCookieStore(
	[]byte("new-authentication-key"),
)

func GetCredentials(w http.ResponseWriter, r *http.Request) (*data.User, int) {
	session, _ := store.Get(r, "rter-credentials")

	user := new(data.User)

	uval, ok := session.Values["username"]
	if !ok {
		return nil, 0
	}

	user.Username, ok = uval.(string)
	if !ok {
		return nil, 0
	}

	pval, ok := session.Values["permissions"]
	if !ok {
		return user, 0
	}

	permissions, ok := pval.(int)
	if !ok {
		return user, 0
	}

	return user, permissions
}

func AuthHandlerFunc(w http.ResponseWriter, r *http.Request) {
	loginInfo := new(data.User)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginInfo)

	if err != nil {
		log.Println("JSON problem")
		log.Println(err)
		http.Error(w, "Malformed json.", http.StatusBadRequest)
		return
	}

	user := new(data.User)
	user.Username = loginInfo.Username

	err = storage.Select(user)

	if err == storage.ErrZeroAffected {
		log.Println("No such user: ", user.Username)
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Database error, likely due to malformed request.", http.StatusInternalServerError)
		return
	}

	if !user.Auth(loginInfo.Password) {
		log.Println("Wrong Password for: ", user.Username)
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		return
	}

	role := new(data.Role)
	role.Title = user.Role

	err = storage.Select(role)

	if err != nil {
		log.Println("Issues loading role during auth", err)
	}

	session, _ := store.Get(r, "rter-credentials")

	session.Values["username"] = user.Username
	session.Values["role"] = user.Role
	session.Values["permissions"] = role.Permissions

	err = session.Save(r, w)

	if err != nil {
		log.Println(err)
	}
}
