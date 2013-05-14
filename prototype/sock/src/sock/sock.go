package main

import (
	"fmt"
	"net"
	"os"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"io"
	"strings"
	"sync"
	"net/url"
	"bytes"
	"time"
	"strconv"
)

type Jar struct {
    lk      sync.Mutex
    cookies map[string][]*http.Cookie
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	if jar.cookies == nil {
		jar.cookies = make(map[string][]*http.Cookie)
	}
	jar.cookies[u.Host] = cookies
//	for _, val := range jar.cookies[u.Host] {
//		fmt.Println(val.String() + "\n")
//	}
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}


func main() {

	service := ":1200"
	listener, err := net.Listen("tcp", service)
	checkError(err)

	
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(30000 * time.Millisecond))

//	fo, err := os.Create("output.ts")
//	if err != nil {
//		panic(err)
//	}
//	defer func() {
//		if err := fo.Close(); err != nil {
//			panic(err)
//		}
//	}()
	
	host := "http://localhost:8000"
	jar := new(Jar)
	client := http.Client{nil, nil, jar}
	
	req, err := http.NewRequest("POST", host + "/auth", strings.NewReader("{\"Username\":\"anonymous\",\"Password\":\"anonymous\"}"))
	if err != nil {
		fmt.Printf("%#v\n", err);
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%#v\n", err);
		return
	}
	cookies := strings.Split(resp.Header.Get("Set-Cookie"), ";")
	var auth string = ""
	for _, val := range cookies {
		val = strings.TrimSpace(val)
		if strings.HasPrefix(val, "rter-credentials") {
			auth = strings.SplitN(val, "=", 2)[1]
		}
	}
	if auth == "" {
		return
	}
	resp.Body.Close()
	
	reqBody := `{"Type":"streaming-video-v1","Live":true,"StartTime":"` + time.Now().Add(-4 * time.Hour).Format("2006-01-02T15:04:05Z") + `","HasGeo":true,"HasHeading":true}`
	req, err = http.NewRequest("POST", host + "/1.0/items", strings.NewReader(reqBody))
	req.Header.Set("Cookie", "rter-credentials=" + auth)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("%#v\n", err);
		return
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%#v\n", err);
		return
	}
	var j interface{}
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		fmt.Printf("%#v\n", err);
		return
	}
	data := j.(map[string]interface{})
	token := data["Token"].(map[string]interface{})
	itemId := strconv.FormatFloat(data["ID"].(float64), 'f', 0, 64)
	fmt.Println(token["rter_resource"].(string))
	resp.Body.Close()
	
	defer func() {
		reqBody := `{"Live":false,"StopTime":"` + time.Now().Add(-4 * time.Hour - 3 * time.Second).Format("2006-01-02T15:04:05Z") + `"}`
		req, err = http.NewRequest("PUT", host + "/1.0/items/" + itemId, strings.NewReader(reqBody))
		req.Header.Set("Cookie", "rter-credentials=" + auth)
		fmt.Println(reqBody);
		resp, err = client.Do(req)
	}();
	
	const chunkSize int = 192000
	var bbuf [chunkSize]byte
	var buffer = new(bytes.Buffer)
	
	fo, err := os.Create("output.ts")
	defer func() {
		if err := fo.Close(); err != nil {
		    panic(err)
	}	
	}()	
	if err != nil { 
		panic(err)
	}
	
	for {
		n, err := conn.Read(bbuf[0:])
		if err != nil && err != io.EOF {
			fmt.Printf("%#v\n", err);
			return
		}
		buffer.Write(bbuf[0:n])
		if buffer.Len() > chunkSize || err == io.EOF {
			chunk := buffer.Next(chunkSize)
			req, err2 := http.NewRequest("POST", token["rter_resource"].(string) + "/ts", bytes.NewReader(chunk))
			if err2 != nil {
				fmt.Printf("%#v\n", err2);
				return
			}
			req.Header.Set("Cookie", "rter-credentials=" + auth)
			req.Header.Set("Authorization", "rtER rter_resource=\"" + token["rter_resource"].(string) + "\", rter_valid_until=\"" + token["rter_valid_until"].(string) + "\", rter_signature=\"" + token["rter_signature"].(string) + "\"")
			resp, err2 = client.Do(req)
			fmt.Printf("%#v\n", "Sent packet");
			if err2 != nil {
				fmt.Printf("%#v\n", err2);
				return
			}
			if _, err := fo.Write(chunk); err != nil {
			    panic(err)
			}
			fmt.Printf("%#v\n", "Chunk written");
			
		}
		if err == io.EOF {
			fmt.Printf("%#v\n", "End Of Stream");
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
