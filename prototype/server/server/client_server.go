package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Name string
	Body string
	Time int64
}

func ClientAjax(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application.json")

	fmt.Println("URL: ", r.URL)

	req := make([]byte, 100)
	n, ok := r.Body.Read(req)

	fmt.Println("Body Result: ", n, ok)
	fmt.Println("Body", string(req))

	socks := r.FormValue("socks")
	apples := r.FormValue("apples")

	fmt.Println("Form Values: ", socks, apples)

	m := Message{"Alice", "Hello", 1294706395881547000}

	b, _ := json.Marshal(m)

	fmt.Println(string(b))

	w.Write(b)
}
