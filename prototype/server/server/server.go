package server

const imagePath = "images/"
const templatePath = "templates/"

func checkError(error error) {
	if error != nil {
		panic(error)
	}
}
