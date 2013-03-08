//  rtER Project
//
//  Author: echa@cim.mcgill.ca

// Test binary upload to server with
// curl -i --data-binary @videoserver.go http://localhost:6666/v1/ingest/10/avc

package main

import (
	"github.com/gorilla/mux"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"flag"
	"log"
)

// Command Line Options
var configfile = flag.String("config", "config.json", "server config file")

// Server Configuration
type ServerConfig struct {
	// server
	Server struct {
		Addr string `json:"addr"`
		Port uint64 `json:"port"`
		Secure_mode bool `json:"secure_mode"`
		Cert_file string `json:"cert_file"`
		Key_file string `json:"key_file"`
	}
	// limits
	Limits struct {
		Max_memory_mbytes uint64 `json:"max_memory_mbytes"`
		Max_ingest_sessions uint64 `json:"max_ingest_sessions"`
		Max_ingest_bandwidth_kbit uint64 `json:"max_ingest_bandwidth_kbit"`
		Rate_limit_enable bool `json:"rate_limit_enable"`
		Rate_limit_ingest_window uint64 `json:"rate_limit_ingest_window"`
		Rate_limit_ingest_sessions_per_source uint64 `json:"rate_limit_ingest_sessions_per_source"`
		Rate_limit_ingest_bytes_per_source uint64 `json:"rate_limit_ingest_bytes_per_source"`
	}
	Auth struct {
	// auth
		Signkey string `json:"signkey"`
	}
	// ingest
	Ingest struct {
		Enable_avc_ingest bool `json:"avc"`
		Enable_ts_ingest bool `json:"ts"`
		Enable_chunk_ingest bool `json:"chunk"`
	}
	// paths
	Paths struct {
		Data_storage_path string `json:"storage"`
	}
	// transcode
	Transcode struct {
		Enable_hls_transcode bool `json:"hls"`
		Enable_mp4_transcode bool `json:"mp4"`
		Enable_ogg_transcode bool `json:"ogg"`
		Enable_dash_transcode bool `json:"dash"`
		Enable_thumb_transcode bool `json:"thumb"`
		Enable_poster_transcode bool `json:"poster"`
	}
}

func ParseConfig(c *ServerConfig) {

	// set default values
	c.Server.Addr = "127.0.0.1"
	c.Server.Port = 8080
	c.Server.Secure_mode = false
	c.Server.Cert_file = ""
	c.Server.Key_file = ""
	c.Limits.Max_memory_mbytes = 128
	c.Limits.Max_ingest_sessions = 10
	c.Limits.Max_ingest_bandwidth_kbit = 10000
	c.Limits.Rate_limit_enable = false
	c.Limits.Rate_limit_ingest_window = 15
	c.Limits.Rate_limit_ingest_sessions_per_source = 100
	c.Limits.Rate_limit_ingest_bytes_per_source = 134217728
	c.Auth.Signkey = "none"
	c.Ingest.Enable_avc_ingest = true
	c.Ingest.Enable_ts_ingest = true
	c.Ingest.Enable_chunk_ingest = false
	c.Paths.Data_storage_path = "./data"
	c.Transcode.Enable_hls_transcode = true
	c.Transcode.Enable_mp4_transcode = false
	c.Transcode.Enable_ogg_transcode = false
	c.Transcode.Enable_dash_transcode = false
	c.Transcode.Enable_thumb_transcode = true
	c.Transcode.Enable_poster_transcode = true

	// read config
    jsonconfig, err := ioutil.ReadFile(*configfile)
	if err != nil {
    	log.Fatalf("Error reading config file: %s\n", err)
    }

    // unpack config from JSON into Go struct (sets only the defined values)
	err = json.Unmarshal(jsonconfig, &c)
	if err != nil {
    	log.Fatalf("Error parsing config file: %s\n", err)
    }

    log.Printf("ServerConfig: %+v\n", c)
}

var c ServerConfig

func main() {

	var err error
	flag.Parse()
	ParseConfig(&c)

	// set up endpoints
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()


	if c.Ingest.Enable_avc_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/avc", AVCIngestHandler).Methods("POST")
	}

	if c.Ingest.Enable_ts_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/ts", TSIngestHandler).Methods("POST", "GET")
	}

	if c.Ingest.Enable_chunk_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/chunk", ChunkIngestHandler).Methods("POST", "GET")
	}
/*
	s.HandleFunc("/videos/{id:[0-9]+}/mp4", MP4FileHandler).Methods("GET")
	s.HandleFunc("/videos/{id:[0-9]+}/ogg", OGGFileHandler).Methods("GET")
	s.HandleFunc("/videos/{id:[0-9]+}/m3u8", M3U8FileHandler).Methods("GET")
	s.HandleFunc("/videos/{id:[0-9]+}/mpd", MPDFileHandler).Methods("GET")


	s.HandleFunc("/videos/{id:[0-9]+}/hls/{segment}", HLSSegmentHandler).Methods("GET")
	s.HandleFunc("/videos/{id:[0-9]+}/dash/{segment}", DASHSegmentHandler).Methods("GET")


	s.HandleFunc("/previews/{id:[0-9]+}/{thumbid:[0-9]+}", ThumbHandler).Methods("GET")
	s.HandleFunc("/previews/{id:[0-9]+}/{posterid:[0-9]+}", PosterHandler).Methods("GET")
*/

	// catch all (redirect non-registered routes to index '/')
	r.HandleFunc("/", IndexHandler)
	r.PathPrefix("/").HandlerFunc(RedirectHandler)

	// attach router to HTTP(s) server
	http.Handle("/", r)

	// run the server
	serverAddr := c.Server.Addr + ":" + strconv.FormatUint(c.Server.Port, 10)
	if c.Server.Secure_mode {
	    log.Printf("HTTPS Server running at %s\n", serverAddr)
		err = http.ListenAndServeTLS(serverAddr, c.Server.Cert_file, c.Server.Key_file, nil)
	} else {
		log.Printf("HTTP Server running at %s\n", serverAddr)
		err = http.ListenAndServe(serverAddr, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}


// Internal Error Reasons
//const (
//	SOURCE_TIMEOUT
//   STREAM_EOS // no error
//    TRANSCODER_FAILED
//)

// Signalled errors to client (HTTP response codes)

// rate limit exceded (calls, bytes?)
// auth failure: token expired, token invalid
// bitstream invalid
// stale stream id

type ServerError struct {
	msg    string
	code   int
	status int
}

func NewServerError(m string, c int, s int) *ServerError {
	return &ServerError{msg: m, code: c, status: s}
}

func (e *ServerError) Error() string { return e.msg }
func (e *ServerError) Code() int { return e.code }
func (e *ServerError) Status() int { return e.status }

func (e *ServerError) JSONError() string {
	return  "{\"errors\":[{\"code\": " +
	        strconv.Itoa(e.code) +
	        "\"message\": \"" +
	        e.msg +
	        "\"}]}"
}


// new HTTP status codes not defined in net/http
const (
	StatusTooManyRequests = 429
)


//func CheckGlobalRateLimit() {}

func AuthenticateRequest(r *http.Request) *ServerError {
	return nil
}


func AVCIngestHandler(w http.ResponseWriter, r *http.Request) {
// todo
// - find a way to limit bandwidth (max_ingest_bandwidth_kbit)
//
// Ingest Loop (this function is called once per video frame, also for new sessions)
// - confirm request validity (signature)
// - confirm request freshness (time issued)
// - confirm uniqueness of uid (stream not already at EOS) - keep state only until token expires


// - [LOCK] check rate limit (max_ingest_sessions)
// - if this is the first request for this video uid
//   - check quota (rate_limit_ingest_sessions_per_source)
//   - launch transcoder
//	 - set disconnect timeout
// - if this is a continuation request
//   - find transcoder pipe
//   - update disconnect timeout
// - check quota (rate_limit_ingest_bytes_per_source)
// - strip headers from data (check for multipart, ascii, form, etc.)
// - check bitstream validity
// - forward binary data
// - update statistics (rate quota)
// - prepare response headers

 	// authenticate the request
 	err := AuthenticateRequest(r)
	if err != nil {
		// return error response in httpcode
		http.Error(w, err.JSONError(), err.Status())
		return
	}

	// extract the video UID from the request
	vars := mux.Vars(r)
	uidstring := vars["id"]
 	uid, _ := strconv.ParseUint(uidstring, 10, 64)

 	// get the session object (atomic)
	session, err := server.FindOrCreateSession(uid)

	if err != nil {
		// fail return error response in httpcode
		http.Error(w, err.JSONError(), err.Status())
		return
	}

	// open new sessions first
	if !session.IsOpen() {
		session.Open(GetTranscodeParams(TRANSCODE_TYPE_AVC))
	}

	// forward data
	if session.IsOpen() {
		session.Write(r.Body)
	}
}

func TSIngestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is a MPEG2-TS handler.\n"))
}

func ChunkIngestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is a chunk handler.\n"))
}
