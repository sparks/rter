package utils

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var RterDir = os.Getenv("RTER_DIR")

var UploadPath = filepath.Join(RterDir, "uploads")
var TemplatePath = filepath.Join(RterDir, "templates")
var ResourcePath = filepath.Join(RterDir, "resources")

var PhoneIDValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
