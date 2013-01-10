package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func fetchStockImage(x, y int) image.Image {
	resp, err := http.Get(fmt.Sprintf("http://lorempixel.com/%v/%v/", x, y))
	checkError(err)

	image, err := jpeg.Decode(resp.Body)
	checkError(err)

	return image
}

func loadFile(filename string) image.Image {
	imageFile, error := os.Open(filename)
	checkError(error)
	defer imageFile.Close()

	image, err := png.Decode(imageFile)
	checkError(err)

	return image
}

func multipartUpload(image image.Image) {
	fmt.Println("Performing Multipart Image Upload:")
	fmt.Println("==================================")

	pipeReader, pipeWriter := io.Pipe()

	multipartWriter := multipart.NewWriter(pipeWriter)
	contentType := multipartWriter.FormDataContentType()
	fmt.Println(contentType)

	go func() {
		response, error := http.Post("http://localhost:8080/multiup", contentType, pipeReader)
		checkError(error)

		fmt.Println("Response:", response)
	}()

	multipartImageWriter, error := multipartWriter.CreateFormFile("image", "image.png")
	checkError(error)

	error = png.Encode(multipartImageWriter, image)

	checkError(error)

	multipartNameWriter, error := multipartWriter.CreateFormField("name")
	checkError(error)

	io.WriteString(multipartNameWriter, "phone_identifier")

	pipeWriter.Close()
	multipartWriter.Close()
	pipeReader.Close()
}

func regularPNGUpload(filename string) {
	fmt.Println("Performing Regular PNG Upload")
	fmt.Println("=============================")

	imageFile, error := os.Open(filename)
	checkError(error)
	defer imageFile.Close()

	response, error := http.Post("http://localhost:8080/upload", "image/png", imageFile)
	checkError(error)

	fmt.Println("Response:", response)
}

func checkError(error error) {
	if error != nil {
		panic(error)
	}
}

func main() {
	multipartUpload(fetchStockImage(200, 200))
}
