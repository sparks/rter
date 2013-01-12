package server

import (
	"regexp"
)

const imagePath = "images/"
const resourcePath = "templates/resources/"
const templatePath = "templates/"

var phoneIDValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

var filenameValidator = regexp.MustCompile("^[a-zA-Z0-9_]*\\.?[a-zA-Z0-9_]+\\.[a-zA-Z0-9]+$")
var folderNameValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func checkError(error error) {
	if error != nil {
		panic(error)
	}
}
