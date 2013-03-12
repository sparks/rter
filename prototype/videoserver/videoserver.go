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
//   - request authentication and group-based access to endpoints
//   - HTTPS cert/key
//   - rate control: enforce sessions per source (IP) per time -> Redis
//   - rate control: enforce bytes per source (IP) per time -> Redis
//   - insert quota headers into replies
//   - limit bandwidth consumption
//   - implement server status endpoint
// - AVC Transcoding Pipeline
//   - check format compliance (H264 NALU headers, profile/level, SPS/PPS existence)
//   - session uniqueness -> Redis
// - TS Transcoding Pipeline
//   - ... same as AVC, different ffmpeg parameters maybe?
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
	"os"
)

var c ServerConfig

func main() {

	var err error
	ParseConfig(&c)

	// create path for transcoder logfiles
	if err = os.MkdirAll(c.Transcode.Log_file_path, PERM_DIR); err != nil {
		log.Fatal("cannot create log directory %s: %s", c.Transcode.Log_file_path, err)
	}

	// set up HTTP endpoints
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	if c.Ingest.Enable_avc_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/avc", AVCIngestHandler).Methods("POST")
	}
/*
	if c.Ingest.Enable_ts_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/ts", TSIngestHandler).Methods("POST")
	}

	if c.Ingest.Enable_chunk_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/chunk", ChunkIngestHandler).Methods("POST")
	}

	s.HandleFunc("/videos/{id:[0-9]+}/mp4", MP4FileHandler).Methods("GET")
	s.HandleFunc("/videos/{id:[0-9]+}/ogg", OGGFileHandler).Methods("GET")
	s.HandleFunc("/videos/{id:[0-9]+}/webm", WEBMFileHandler).Methods("GET")
	s.HandleFunc("/videos/{id:[0-9]+}/m3u8", M3U8FileHandler).Methods("GET")
	s.HandleFunc("/videos/{id:[0-9]+}/mpd", MPDFileHandler).Methods("GET")


	s.HandleFunc("/videos/{id:[0-9]+}/hls/{segment}", HLSSegmentHandler).Methods("GET")
	s.HandleFunc("/videos/{id:[0-9]+}/dash/{segment}", DASHSegmentHandler).Methods("GET")


	s.HandleFunc("/previews/{id:[0-9]+}/thumb/{thumbid:[0-9]+}", ThumbHandler).Methods("GET")
	s.HandleFunc("/previews/{id:[0-9]+}/poster/{posterid:[0-9]+}", PosterHandler).Methods("GET")
*/

	// catch all (redirect non-registered routes to index '/')
	r.HandleFunc("/", IndexHandler)
	//r.PathPrefix("/").HandlerFunc(RedirectHandler)

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


func AuthenticateRequest(r *http.Request) *ServerError {

	// check auth headers

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

	// extract the video UID from the request
	vars := mux.Vars(r)
	uidstring := vars["id"]

 	// authenticate the request
 	var err *ServerError
 	var session *TranscodeSession

 	// confirm validity and freshness of request
	if err = AuthenticateRequest(r); err != nil {
		// return error response in httpcode
		http.Error(w, err.JSONError(), err.Status())
		return
	}

 	// get or create the session object if quota permits
	if session, err = server.FindOrCreateSession(uidstring, TC_INGEST_AVC); err != nil {
		// return error response in httpcode
		http.Error(w, err.JSONError(), err.Status())
		return
	}

	// forward data
	if session.IsOpen() {
		if err = session.Write(r); err != nil {
			// return error response in httpcode
			http.Error(w, err.JSONError(), err.Status())
		}
	}

	// write response headers
}

func TSIngestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is a MPEG2-TS handler.\n"))
}

func ChunkIngestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is a chunk handler.\n"))
}
