// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

// Test binary upload to server with
// curl -i --data-binary @videoserver.go http://localhost:6666/v1/ingest/10/avc
//
// Test MPEG2TS stream to server
// ffmpeg -v debug -y -re -i file.m4v -vsync 1 -map 0 -codec copy \
//    -bsf h264_mp4toannexb -r 25 -f mpegts -copytb 0 http://localhost:6666/v1/ingest/1/ts
//
// Test Raw H264/AVC stream (also reencodes the file to adhere to our format specs)
// ffmpeg -v debug -y -re -i file.m4v -f h264 -c:v libx264 -preset ultrafast \
//    -tune zerolatency -crf 20 -x264opts keyint=50:bframes=0:ratetol=1.0:ref=1:repeat-headers=1 \
//    -profile baseline -maxrate 1200k -bufsize 1200k -an http://localhost:6666/v1/ingest/1/avc
//
// Todo/Implement: Features
// - Common Server Management
//   - check transcoder capabilities (transcode.go:CheckTranscoderCapabilities())
//   - configuration sanity check (config.go:SanityCheck())
//   - request authentication: API_KEY, REQUEST_TOKEN
//   - HTTPS cert/key
//   - rate control: enforce sessions per source (IP) per time -> Redis
//   - rate control: enforce bytes per source (IP) per time -> Redis
//   - insert quota headers into replies
//   - limit bandwidth consumption
//   - implement server status endpoint
//   - session uniqueness and EOS -> Redis
//   - total server statistics (per session this is already accounted for)
//   - transcoder user/sys time is incorrect
//   - ingest server redirect when quota limit reached -> Redis
// - Transcoding Pipelines
//   - check format compliance (H264 NALU headers, profile/level, SPS/PPS existence)
// - File Download (built-in file server for videos, images, segments)
//   - play endpoint generating HTML <video> in response to a GET request
//   - cache headers, mime-types, text file compression (m3u8, mpd)
//   - byte range support
// - HTTP chunk mode upload endpoint (simple, how to allow continuations?)
// - Interactive Chunk Upload (like DropBox: reorder-safe, byte-range, file multiplexing)
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
	"runtime"
	"log"
	"os"
)

var c ServerConfig

func main() {

	var err error
	c.ParseConfig()
	c.SanityCheck()

    // print config in dev mode
    if !c.Server.Production_mode { c.Print() }

	// set resource limits
	runtime.GOMAXPROCS(c.Limits.Max_cpu)

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
