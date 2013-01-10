package server

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".png") {
		http.NotFound(w, r)
		return
	}

	path := strings.Split(r.URL.Path, "/")

	fi, err := os.Open(imagePath + path[len(path)-1])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer fi.Close()

	io.Copy(w, fi)
}
