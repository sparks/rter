package server

import (
	"os"
	"path/filepath"
	"regexp"
)

var rterDir = os.Getenv("RTER_DIR")

var ImagePath = filepath.Join(rterDir, "images")
var TemplatePath = filepath.Join(rterDir, "templates")
var ResourcePath = filepath.Join(rterDir, "templates", "resources")

var phoneIDValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

var filenameValidator = regexp.MustCompile("^[a-zA-Z0-9_]*\\.?[a-zA-Z0-9_]+\\.[a-zA-Z0-9]+$")
var folderNameValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
