package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
	regUpload()
}

func multiUpload() {
	fmt.Println("Ready")
	fmt.Println("Set")

	fi, err := os.Open("images/adf.png")
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

func regUpload() {
	fmt.Println("Ready")
	fmt.Println("Set")

	fi, err := os.Open("images/adf.png")
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
