package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var phone_ids = []string{
	"1e7f033bfc7b3625fa07c9a3b6b54d2c81eeff98",
	"fe7f033bfc7b3625fa06c9a3b6b54b2c81eeff98",
	"b6200c5cc15cfbddde2874c40952a7aa25a869dd",
	"852decd1fbc083cf6853e46feebb08622d653602",
	"e1830fcefc3f47647ffa08350348d7e34b142b0b",
	"48ad32292ff86b4148e0f754c2b9b55efad32d1e",
	"acb519f53a55d9dea06efbcc804eda79d305282e",
	"ze7f033bfc7b3625fa06c5a316b54b2c81eeff98",
	"t6200c5cc15cfbddde2875c41952a7aa25a869dd",
	"952decd1fbc083cf6853e56f1ebb08622d653602",
	"y1830fcefc3f47647ffa05351348d7e34b142b0b",
	"x8ad32292ff86b4148e0f55412b9b55efad32d1e",
	"qcb519f53a55d9dea06ef5cc104eda79d305282e",
}

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

func multipartUpload(image image.Image, phone_id string, lat, lng float64, heading float64) {
	// fmt.Println("Performing Multipart Image Upload:")
	// fmt.Println("==================================")

	pipeReader, pipeWriter := io.Pipe()

	multipartWriter := multipart.NewWriter(pipeWriter)
	contentType := multipartWriter.FormDataContentType()
	// fmt.Println(contentType)

	responseChan := make(chan *http.Response)

	go func() {
		response, error := http.Post("http://rter.cim.mcgill.ca:8080/multiup", contentType, pipeReader)
		// response, error := http.Post("http://localhost:8080/multiup", contentType, pipeReader)
		checkError(error)
		responseChan <- response
	}()

	imageWriter, error := multipartWriter.CreateFormFile("image", "image.png")
	checkError(error)
	error = png.Encode(imageWriter, image)
	checkError(error)

	idWriter, error := multipartWriter.CreateFormField("phone_id")
	checkError(error)
	io.WriteString(idWriter, phone_id)

	latWriter, error := multipartWriter.CreateFormField("lat")
	checkError(error)
	io.WriteString(latWriter, fmt.Sprintf("%v", lat))

	lngWriter, error := multipartWriter.CreateFormField("lng")
	checkError(error)
	io.WriteString(lngWriter, fmt.Sprintf("%v", lng))

	headingWriter, error := multipartWriter.CreateFormField("heading")
	checkError(error)
	io.WriteString(headingWriter, fmt.Sprintf("%v", heading))

	pipeWriter.Close()
	multipartWriter.Close()
	pipeReader.Close()

	response := <-responseChan
	fmt.Println(response.Status)
	r, _ := ioutil.ReadAll(response.Body)
	fmt.Println("Response Body", string(r))
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
	rand.Seed(time.Now().Unix())

	exit := make(chan bool)
	num_clients := len(phone_ids)

	go func() {
		n := 0
		for {
			for i := 0; i < num_clients; i++ {
				multipartUpload(fetchStockImage(200, 200), phone_ids[i], 45.50745+rand.Float64()/5-0.1, -73.5793+rand.Float64()/5-0.1, -180.0+rand.Float64()*360.0)
				fmt.Println(45.4+rand.Float64()/5, -73.4+rand.Float64()/5)
				fmt.Println(n, i)
				select {
				case <-exit:
					return
				default:
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()

	var input string
	fmt.Scanf("%s", &input)
	fmt.Println("Waiting to exit ...")
	exit <- true
}
