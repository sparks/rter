package rest

import (
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var decoder = schema.NewDecoder()

func parseID(position int, w http.ResponseWriter, r *http.Request) int64 {
	splitPath := strings.Split(r.URL.Path, "/")

	if len(splitPath) < position+1 {
		http.Error(w, "Invalid URI", http.StatusBadRequest)
	}

	ID, err := strconv.Atoi(splitPath[position])

	if err != nil {
		log.Println(err)
		http.Error(w, "Malformed Item ID", http.StatusBadRequest)
	}

	return int64(ID)
}
