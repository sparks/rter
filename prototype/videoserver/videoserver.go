// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

// Test binary upload to server with
// curl -i --data-binary @videoserver.go http://localhost:6666/v1/ingest/10/avc
//
// Test Http Stream to server
// ffmpeg -v debug -y -re -i file.m4v -vsync 1 -map 0 -codec copy \
//    -bsf h264_mp4toannexb -r 25 -f mpegts -copytb 0 http://localhost:6666/v1/ingest/1/ts
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
	if err = os.MkdirAll(c.Transcode.Log_path, PERM_DIR); err != nil {
		log.Fatal("cannot create log directory %s: %s", c.Transcode.Log_path, err)
	}

	// set up HTTP endpoints
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	if c.Ingest.Enable_avc_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/avc", AVCIngestHandler).Methods("POST")
	}

	if c.Ingest.Enable_ts_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/ts", TSIngestHandler).Methods("POST")
	}
/*
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

	// have a single index handler at server URI root
	r.HandleFunc("/", IndexHandler)

	// catch all (redirect non-registered routes to index '/')
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
	w.Write([]byte("This is the rtER video server.\n"))
}


func AuthenticateRequest(r *http.Request) *ServerError {

	// check auth headers

	return nil
}


func AVCIngestHandler(w http.ResponseWriter, r *http.Request) {
	GenericIngestHandler(w, r, TC_INGEST_AVC)
}

func TSIngestHandler(w http.ResponseWriter, r *http.Request) {
	GenericIngestHandler(w, r, TC_INGEST_TS)
}

func ChunkIngestHandler(w http.ResponseWriter, r *http.Request) {
	GenericIngestHandler(w, r, TC_INGEST_CHUNK)
}

func GenericIngestHandler(w http.ResponseWriter, r *http.Request, t int) {

	//
	// Generic Ingest Loop
	//
	// This function is either called once per video frame (frame-wise HTTP PUSH)
	// or once per source ingest session (multi-part HTTP PUSH)
	//
	// - confirm request validity (signature)
	// - confirm request freshness (time issued)
	// - confirm uniqueness of uid (stream not already at EOS) - keep state only until token expires
	// - enforce server rate limit (max_ingest_sessions)
	// - manage transcoding sessions (create on demand, lookup on request)
	// - forward request to session for write handling
	// - forward response to session for setting response headers

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
	if session, err = server.FindOrCreateSession(uidstring, t); err != nil {
		// return error response in httpcode
		http.Error(w, err.JSONError(), err.Status())
		return
	}

	// forward data
	if session.IsOpen() {
		if err = session.Write(r); err != nil {
			// return error response in httpcode
			http.Error(w, err.JSONError(), err.Status())
			return
		}
		// write response headers
		session.SetResponseHeaders(w)
	}

}
