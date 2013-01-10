package server

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var filenameValidator = regexp.MustCompile("^[a-zA-Z0-9]*\\.?[a-zA-Z0-9]+\\.[a-zA-Z0-9]+$")
var folderNameValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	if !(strings.HasSuffix(r.URL.Path, ".png") || strings.HasSuffix(r.URL.Path, ".jpeg") || strings.HasSuffix(r.URL.Path, ".jpg")) {
		// fmt.Println("Wrong Suffix")
		http.NotFound(w, r)
		return
	}

	path := strings.Split(r.URL.Path, "/")

	if len(path) < 3 {
		// fmt.Println("Path too short")
		http.NotFound(w, r)
		return
	}

	if !validateFilePath(path[2:]) {
		// fmt.Println("Invalid Path")
		http.NotFound(w, r)
		return
	}

	fi, err := os.Open(imagePath + strings.Join(path[2:], "/"))
	if err != nil {
		// fmt.Println("No such file")
		http.NotFound(w, r)
		return
	}
	defer fi.Close()

	io.Copy(w, fi)
}

func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else {
		// fmt.Println("Wrong Suffix")
		http.NotFound(w, r)
		return
	}

	path := strings.Split(r.URL.Path, "/")

	if len(path) < 3 {
		// fmt.Println("Path too short")
		http.NotFound(w, r)
		return
	}

	if !validateFilePath(path[2:]) {
		// fmt.Println("Invalid Path")
		http.NotFound(w, r)
		return
	}

	fi, err := os.Open(resourcePath + strings.Join(path[2:], "/"))
	if err != nil {
		// fmt.Println("No such file")
		http.NotFound(w, r)
		return
	}
	defer fi.Close()

	io.Copy(w, fi)
}

func validateFilePath(path []string) bool {
	for i := 0; i < len(path)-1; i++ {
		if !folderNameValidator.MatchString(path[i]) {
			// fmt.Println("Invalid Folder Name")
			return false
		}
	}

	if !filenameValidator.MatchString(path[len(path)-1]) {
		// fmt.Println("Invalid Filename")
		return false
	}

	return true
}
