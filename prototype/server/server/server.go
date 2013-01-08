package server

const lenPath = len("/view/")

const imagePath = "images/"
const templatePath = "templates/"
const dataPath = "data/"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
