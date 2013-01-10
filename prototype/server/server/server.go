package server

const imagePath = "images/"
const resourcePath = "templates/resources/"
const templatePath = "templates/"

func checkError(error error) {
	if error != nil {
		panic(error)
	}
}
