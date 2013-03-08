// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

// Test binary upload to server with
// curl -i --data-binary @videoserver.go http://localhost:6666/v1/ingest/10/avc
//
// Unsure
// - is body already complete when handler is called? if not, how to deal with broken TCP connection
//
// Todo/Implement: Features
// - Common Server Management
//   - ServerError id -> text translation
//   - request authentication
//   - HTTPS cert/key
//   - rate control: sessions per source (IP) per time
//   - rate control: bytes per source (IP) per time
//   - quota headers
//   - limit memory/bandwidth consumption
// - AVC Transcoding Pipeline
//   - ffmpeg parameter assembly (HLS, Thumb, Poster, MP4, OGG)
//	 - process management (start/stop/monitor)
//   - EOS/timeout/close-session handling (and uniqueness assumption)
//   - check format compliance (H264 NALU headers, profile/level, SPS/PPS existence)
// - TS Transcoding Pipeline
// - File Download
//   - cache headers, mime-types, text file compression
//   - byte range support
// - Chunk Upload (reorder, single file multiplexing)
// - Websocket for AVC/TS frame-wise upload


// Limit bandwidth (max_ingest_bandwidth_kbit) and memory (max_memory_mbytes)
// http://stackoverflow.com/questions/14582471/golang-memory-consumption-management
// http://lwn.net/Articles/428100/
// http://evanfarrer.blogspot.ca/2012/05/making-friendly-elastic-applications-in.html
// https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/9JxgtOJqqRU


package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"log"
)

var c ServerConfig

func main() {

	var err error
	ParseConfig(&c)

	// set up endpoints
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()


	if c.Ingest.Enable_avc_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/avc", AVCIngestHandler).Methods("POST")
	}

	if c.Ingest.Enable_ts_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/ts", TSIngestHandler).Methods("POST")
	}

	if c.Ingest.Enable_chunk_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/chunk", ChunkIngestHandler).Methods("POST")
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
