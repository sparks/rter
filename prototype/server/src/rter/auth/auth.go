// Provides and authentication service for the rtER server. Users can present their credentials (in JSON). Upon valid credentials they are given a signed session cookie (use gorilla/sessions) to keep them signed in. Passwords are sent in plain text so it's not great over http.
//
// Also provides a GetCredentials callback which can be used to check if a user is logged in if so get their permissions.
//
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
	[]byte("new-authentication-key"), // This is the key which is used to sign cookies
)

// Examine the request for a valid session cookie. If one is found, return the username and user permissions
func GetCredentials(r *http.Request) (*data.User, int) {
	session, _ := store.Get(r, "rter-credentials") // Will return an empty map if there isn't a valid session cookie

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

// Handle requests containing JSON with user credentials. If the credentials are valid the response will set a session cookie effectively logging in the user.
//
// Expected JSON format:
// 	{"Username":"theusername", "Password":"userpassword"}
//
// Invalid credentials will result in a 401 StatusUnauthorized response. Malformed JSON will result in a 400 StatusBadRequest or possibly 500 StatusInternalServerError.
func AuthHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON into a user object
	loginInfo := new(data.User)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginInfo)

	if err != nil {
		log.Println("JSON problem")
		log.Println(err)
		http.Error(w, "Malformed json.", http.StatusBadRequest)
		return
	}

	// Try and load the actual user
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

	// Validate
	if !user.Auth(loginInfo.Password) {
		log.Println("Wrong Password for: ", user.Username)
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		return
	}

	// Get role permissions
	role := new(data.Role)
	role.Title = user.Role

	err = storage.Select(role)

	if err != nil {
		log.Println("Issues loading role during auth", err)
		role.Permissions = 0 // Default to no permissions
	}

	// Build a cookie session
	session, _ := store.Get(r, "rter-credentials")

	session.Values["username"] = user.Username
	session.Values["role"] = user.Role
	session.Values["permissions"] = role.Permissions

	err = session.Save(r, w)

	if err != nil {
		log.Println(err)
	}
}
