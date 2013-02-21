package server

import (
	"os"
	"path/filepath"
	"regexp"
)

var rter_dir = os.Getenv("RTER_DIR")

var ImagePath = filepath.Join(rter_dir, "images")
var TemplatePath = filepath.Join(rter_dir, "templates")
var ResourcePath = filepath.Join(rter_dir, "templates", "resources")

var phoneIDValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

var filenameValidator = regexp.MustCompile("^[a-zA-Z0-9_]*\\.?[a-zA-Z0-9_]+\\.[a-zA-Z0-9]+$")
var folderNameValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func checkError(error error) {
	if error != nil {
		panic(error)
	}
}
