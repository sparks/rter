package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HTMLHandler(w http.ResponseWriter, r *http.Request) bool {
	if strings.HasSuffix(r.URL.Path, ".html") {
		w.Header().Set("Content-Type", "text/html")
	} else {
		return false
	}

	path := strings.Split(r.URL.Path, "/")

	if len(path) < 2 {
		return false
	}

	if !validateFilePath(path[1:]) {
		return false
	}

	page, err := os.Open(filepath.Join(TemplatePath, filepath.Join(path[1:]...)))

	if err != nil {
		return false
	}
	defer page.Close()

	io.Copy(w, page)
	return true
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
