// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

// Test binary upload to server with
// curl -i --data-binary @videoserver.go http://localhost:6666/v1/ingest/10/avc
//
// Test MPEG2TS stream to server
// ffmpeg -v debug -y -re -i file.m4v -vsync 1 -map 0 -codec copy \
//    -bsf h264_mp4toannexb -r 25 -f mpegts -copytb 0 http://localhost:6660/v1/ingest/1/ts
//
// Test Raw H264/AVC stream (also reencodes the file to adhere to our format specs)
// ffmpeg -v debug -y -re -i file.m4v -f h264 -c:v libx264 -preset ultrafast \
//    -tune zerolatency -crf 20 -x264opts keyint=50:bframes=0:ratetol=1.0:ref=1:repeat-headers=1 \
//    -profile baseline -maxrate 1200k -bufsize 1200k -an http://localhost:6660/v1/ingest/1/avc
//
// Todo/Implement: Features
// - Common Server Management
//   - check :id for >0
//   - check transcoder capabilities (transcode.go:CheckTranscoderCapabilities())
//   - HTTPS cert/key
//   - rate control: enforce sessions per source (IP? or user?) per time -> Redis
//   - rate control: enforce bytes per source (IP? or user?) per time -> Redis
//   - rate control: reset quota after deadline
//   - insert quota headers into replies
//   - limit bandwidth consumption
//   - enforce token lifetime when session is estabished
//   - implement server status endpoint (Package expvar)
//   - total server statistics (per session this is already accounted for)
//   - transcoder user/sys time is incorrect
//   - session uniqueness and EOS -> Redis
//   - ingest server redirect when quota limit reached -> Redis
// - Transcoding Pipelines
//   - select best x264 options for chosen level (3.0)
//   - check whether x264 options are correct for chosen level (3.0)
//   - check format compliance (H264 NALU headers, profile/level, SPS/PPS existence)
// - File Download (built-in file server for videos, images, segments)
//   - check if cache headers are properly set
//   - enable text file compression (m3u8, mpd)
// - HTTP chunk mode upload endpoint for files (simple, how to allow continuations?)
// - Interactive Chunk Upload (like DropBox: reorder-safe, byte-range, file multiplexing)
// - Websocket for AVC/TS upload
// - server bandwidth testing (ab or httperf): ab -c 500 -n 500 http://localhost:1234/

// API key generation
// http://stackoverflow.com/questions/1448455/php-api-key-generator

// OAuth 1.0a
// http://oauth.net/core/1.0a

// Limit bandwidth (max_ingest_bandwidth_kbit) and memory (max_memory_mbytes)
// http://stackoverflow.com/questions/14582471/golang-memory-consumption-management
// http://lwn.net/Articles/428100/
// http://evanfarrer.blogspot.ca/2012/05/making-friendly-elastic-applications-in.html
// https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/9JxgtOJqqRU

package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"videoserver/config"
	"videoserver/server"
	"videoserver/utils"
)

var C config.ServerConfig
var S *server.State

func main() {

	C.ParseConfig()

	// redirect logging if logfile is specified and can be created
	if C.Server.Logfile != "" {
		file, err := os.OpenFile(C.Server.Logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, utils.PERM_FILE)

		if err == nil {
			log.SetOutput(file)
		}
	}

	var err error

	if !C.Server.Production_mode {
		C.Print()
	}
	C.SanityCheck()
	C.CheckTranscoderCapabilities()

	S = server.NewServer(&C)

	// set resource limits
	runtime.GOMAXPROCS(C.Limits.Max_cpu)

	// create path for transcoder logfiles
	if err = os.MkdirAll(C.Transcode.Log_path, utils.PERM_DIR); err != nil {
		log.Fatal("cannot create log directory %s: %s", C.Transcode.Log_path, err)
	}

	// set up HTTP endpoints
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	if C.Ingest.Enable_avc_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/avc", AVCIngestHandler).Methods("POST")
	}

	if C.Ingest.Enable_ts_ingest {
		s.HandleFunc("/ingest/{id:[0-9]+}/ts", TSIngestHandler).Methods("POST")
	}

	// playback handler used for development only
	if !C.Server.Production_mode {
		s.HandleFunc("/videos/{id:[0-9]+}/play", PlaybackHandler).Methods("GET")
	}

	/*
		if C.Ingest.Enable_chunk_ingest {
			s.HandleFunc("/ingest/{id:[0-9]+}/chunk", ChunkIngestHandler).Methods("POST")
		}

		s.HandleFunc("/videos/{id:[0-9]+}/video.mp4", MP4FileHandler).Methods("GET")
		s.HandleFunc("/videos/{id:[0-9]+}/video.ogv", OGGFileHandler).Methods("GET")
		s.HandleFunc("/videos/{id:[0-9]+}/video.webm", WEBMFileHandler).Methods("GET")
		s.HandleFunc("/videos/{id:[0-9]+}/index.m3u8", M3U8FileHandler).Methods("GET")
		s.HandleFunc("/videos/{id:[0-9]+}/index.mpd", MPDFileHandler).Methods("GET")
		s.HandleFunc("/videos/{id:[0-9]+}/hls/{segment}.ts", HLSSegmentHandler).Methods("GET")
		s.HandleFunc("/videos/{id:[0-9]+}/dash/{segment}.ts", DASHSegmentHandler).Methods("GET")
		s.HandleFunc("/videos/{id:[0-9]+}/thumb/{thumbid:[0-9]+}.jpg", ThumbHandler).Methods("GET")
		s.HandleFunc("/videos/{id:[0-9]+}/poster/{posterid:[0-9]+}.jpg", PosterHandler).Methods("GET")
	*/

	// have a single index handler at server URI root
	r.HandleFunc("/", IndexHandler)

	if !C.Server.Production_mode {
		s.PathPrefix("/videos/").Handler(http.StripPrefix("/v1/videos/",
			http.FileServer(http.Dir(C.Transcode.Output_path))))
	}

	// catch all (redirect non-registered routes to index '/')
	//r.PathPrefix("/").HandlerFunc(RedirectHandler)

	// attach router to HTTP(s) server
	http.Handle("/", r)

	// run the server
	serverAddr := C.Server.Addr + ":" + strconv.FormatUint(C.Server.Port, 10)
	if C.Server.Secure_mode {
		log.Printf("HTTPS Server running at %s\n", serverAddr)
		err = http.ListenAndServeTLS(serverAddr, C.Server.Cert_file, C.Server.Key_file, nil)
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

func AVCIngestHandler(w http.ResponseWriter, r *http.Request) {
	GenericIngestHandler(w, r, server.TC_INGEST_AVC)
}

func TSIngestHandler(w http.ResponseWriter, r *http.Request) {
	GenericIngestHandler(w, r, server.TC_INGEST_TS)
}

func ChunkIngestHandler(w http.ResponseWriter, r *http.Request) {
	GenericIngestHandler(w, r, server.TC_INGEST_CHUNK)
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
	var err *server.Error
	var session *server.Session

	if C.Auth.Enabled {
		// confirm validity and freshness of request
		if err = server.AuthenticateRequest(r, C.Auth.Token_secret); err != nil {
			// return error response in httpcode
			server.ServeError(w, err.JSONError(), err.Status())
			return
		}
	}

	// get or create the session object if quota permits
	if session, err = S.FindOrCreateSession(uidstring, t); err != nil {
		// return error response in httpcode
		server.ServeError(w, err.JSONError(), err.Status())
		return
	}

	// forward data
	if session.IsOpen() {
		if err = session.Write(r, t); err != nil {
			// return error response in httpcode
			server.ServeError(w, err.JSONError(), err.Status())
			return
		}
		// write response headers
		session.SetResponseHeaders(w)
	}

}

// generates and returns a simple HTML5 website containing a video element
const (
	PLAY_TMPL_BEGIN    string = `<!doctype html><html lang=en><head><meta charset=utf-8><title>rtER Video Demo Stream {{.}} -- [Dev Mode]</title></head><body><video controls autoplay poster="/v1/videos/{{.}}/poster/000000001.jpg" x-webkit-airplay="allow">`
	PLAY_TMPL_SRC_HLS  string = `<source src="/v1/videos/{{.}}/index.m3u8" type="application/x-mpegURL">`
	PLAY_TMPL_SRC_MP4  string = `<source src="/v1/videos/{{.}}/video.mp4" type="video/mp4; codecs=avc1.42E01E,mp4a.40.2">`
	PLAY_TMPL_SRC_WEBM string = `<source src="/v1/videos/{{.}}/video.webm" type="video/webm; codecs=vp8,vorbis">`
	PLAY_TMPL_SRC_OGG  string = `<source src="/v1/videos/{{.}}/video.ogv" type="video/ogg; codecs=theora,vorbis">`
	PLAY_TMPL_END      string = `</video></body></html>`
)

func PlaybackHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uidstring := vars["id"]

	// construct the template
	tpl := PLAY_TMPL_BEGIN
	if C.Transcode.Hls.Enabled {
		tpl += PLAY_TMPL_SRC_HLS
	}
	if C.Transcode.Mp4.Enabled {
		tpl += PLAY_TMPL_SRC_MP4
	}
	if C.Transcode.Webm.Enabled {
		tpl += PLAY_TMPL_SRC_WEBM
	}
	if C.Transcode.Ogg.Enabled {
		tpl += PLAY_TMPL_SRC_OGG
	}
	tpl += PLAY_TMPL_END

	t, err := template.New("player").Parse(tpl)
	if err != nil {
		log.Fatalf("Error parsing player template template: %s\n", err)
	}

	err = t.Execute(w, uidstring)
	if err != nil {
		log.Fatalf("Error generating player template string: %s\n", err)
	}
}
