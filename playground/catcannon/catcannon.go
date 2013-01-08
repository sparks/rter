package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	multiUpload()
}

func multiUpload() {
	fmt.Println("Ready")
	fmt.Println("Set")

	pread, pwrite := io.Pipe()
	defer pread.Close()
	mpwrite := multipart.NewWriter(pwrite)

	content_type := mpwrite.FormDataContentType()
	fmt.Println(content_type)

	respchan := make(chan *http.Response)

	go func() {
		resp, err := http.Post("http://localhost:8080/multiup", content_type, pread)
		if err != nil {
			panic(err)
		}
		respchan <- resp
	}()

	mpfilewrite, err := mpwrite.CreateFormFile("image", "tomato.png")
	if err != nil {
		panic(err)
	}
	fi, err := os.Open("cat.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	io.Copy(mpfilewrite, fi)
	mpwrite.Close()
	pwrite.Close()

	resp := <-respchan
	fmt.Println(resp)
	fmt.Println("Fire")
}

func regUpload() {
	fmt.Println("Ready")
	fmt.Println("Set")

	fi, err := os.Open("cat.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	resp, err := http.Post("http://localhost:8080/upload", "image/png", fi)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	fmt.Println("Fire")
}
