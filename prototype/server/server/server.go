package server

const lenPath = len("/view/")

const imagePath = "images/"
const templatePath = "templates/"
const dataPath = "data/"

func checkError(error error) {
	if error != nil {
		panic(error)
	}
}
