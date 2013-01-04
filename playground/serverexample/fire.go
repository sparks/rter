package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Ready")
	fmt.Println("Set")

	fi, err := os.Open("images/adf.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	reader := bufio.NewReader(fi)

	resp, err := http.Post("http://localhost:8080/upload", "image/png", reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

	fmt.Println("Fire")
}
