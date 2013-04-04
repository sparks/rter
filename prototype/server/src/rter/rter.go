// The rtER web Server provides the RESTful API, Streaming API, Authentication and the web client interface.
//
// RESTful API access to all the rtER data structures and user account with authentication support. The Streaming API provides real-time stream of certain data structures. Secure authentication with cookies is provided. A interactive and collaborative web client is served.is
//
// Options:
// 	-gzip=false: enable gzip compression
// 	-http=true: enable http
// 	-http-port=8080: set the http port to use
// 	-https=false: enable https
// 	-https-port=10433: set the https port to use
// 	-log-file="": set server logfile
// 	-probe=0: perform logging on requests
//	-serve-log-file=true: serve log file over http
//
// Env Variable:
// 	RTER_LOGFILE: set server log file (flag takes precedence)
// 	RTER_DIR: set the dir where the 'www' and 'uploads' directories are located
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rter/auth"
	"rter/compressor"
	"rter/legacy"
	"rter/rest"
	"rter/storage"
	"rter/streaming"
)

// TODO: Make a flag for the secret token (video server)
// TODO: Make a flag for the cookie signing (auth)
// TODO: Make a flag for the video server URI
var (
	probeLevel    = flag.Int("probe", 0, "perform logging on requests")
	httpsFlag     = flag.Bool("https", false, "enable https")
	httpFlag      = flag.Bool("http", true, "enable http")
	gzipFlag      = flag.Bool("gzip", false, "enable gzip compression")
	logfile       = flag.String("log-file", "", "set server logfile")
	serveLogFlag  = flag.Bool("serve-log-file", true, "serve logfile over http")
	httpPort      = flag.Int("http-port", 8080, "set the http port to use")
	httpsPort     = flag.Int("https-port", 10433, "set the https port to use")
	rterDir       = flag.String("rter-dir", "", "sets the dir 'www' and 'uploads' will be")
	sockDebugFlag = flag.Bool("sock-debug", false, "debug the websocket connections")
)

func main() {
	flag.Parse()

	setupLogger()
	setupRterDir()

	var uploadPath = filepath.Join(*rterDir, "uploads")
	var wwwPath = filepath.Join(*rterDir, "www")

	log.Println("Launching rtER Server")

	err := storage.OpenStorage("rter", "j2pREch8", "tcp", "localhost:3306", "rter")

	if err != nil {
		log.Fatalf("Failed to open connection to database %v", err)
	}
	defer storage.CloseStorage()

	// First setup the subrouters

	r := mux.NewRouter().StrictSlash(true)

	sr := streaming.NewStreamingRouter()
	sr.Debug(*sockDebugFlag)
	r.PathPrefix("/1.0/streaming").Handler(http.StripPrefix("/1.0/streaming", sr)) // Must register more specific paths first

	r.PathPrefix("/1.0").Handler(http.StripPrefix("/1.0", rest.CRUDRouter())) // Less specific paths later

	// Hand static files

	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads", http.FileServer(http.Dir(uploadPath))))

	r.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir(filepath.Join(wwwPath, "css")))))
	r.PathPrefix("/js").Handler(http.StripPrefix("/js", http.FileServer(http.Dir(filepath.Join(wwwPath, "js")))))
	r.PathPrefix("/vendor").Handler(http.StripPrefix("/vendor", http.FileServer(http.Dir(filepath.Join(wwwPath, "vendor")))))
	r.PathPrefix("/asset").Handler(http.StripPrefix("/asset", http.FileServer(http.Dir(filepath.Join(wwwPath, "asset")))))
	r.PathPrefix("/template").Handler(http.StripPrefix("/template", http.FileServer(http.Dir(filepath.Join(wwwPath, "template")))))

	r.HandleFunc("/favicon.ico",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(wwwPath, "asset", "favicon.ico"))
		},
	).Methods("GET")

	if *serveLogFlag { // Should be run after setupLogger() since it depends on setting up logfile
		if *logfile == "" {
			log.Println("\t-Serve Log Disable (No Log File)")
		} else {
			log.Println("\t-Serve Log Enabled")
			r.HandleFunc("/log",
				func(w http.ResponseWriter, r *http.Request) {
					http.ServeFile(w, r, *logfile)
				},
			).Methods("GET")
		}
	}

	// Specific Handlers

	r.HandleFunc("/auth", auth.AuthHandlerFunc).Methods("POST") // Authentication service
	r.HandleFunc("/multiup",                                    // Legacy support for android prototype app
		func(w http.ResponseWriter, r *http.Request) {
			legacy.MultiUploadHandler(*rterDir, uploadPath, w, r)
		},
	)

	r.HandleFunc("/", // Web client
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(wwwPath, "index.html"))
		},
	)

	// Server final setup and adjustement

	r.NotFoundHandler = http.HandlerFunc(Debug404) // Catch all 404s

	var rootHandler http.Handler = r

	if *gzipFlag { // Wrap rootHandler with on the fly gzip compressor
		log.Print("\t-GZIP Enabled. Warning Websockets are flaky with gzip")
		rootHandler = compressor.GzipHandler(rootHandler)
	}

	if *probeLevel > 0 { // Wrap rootHandler with debugging probe
		log.Print("\t-Probe Enabled, Level ", *probeLevel)
		rootHandler = ProbeHandler(*probeLevel, rootHandler)
	}

	http.Handle("/", rootHandler)

	// Launch Server

	waits := make([]chan bool, 0) // Prevent from quitting till server routines finish

	if *httpsFlag { // HTTPS
		httpsChan := make(chan bool)
		waits = append(waits, httpsChan)

		go func() {
			log.Println(fmt.Sprintf("\t-Using HTTPS on port %v", *httpsPort))
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%v", *httpsPort), "cert.pem", "key.pem", nil))

			httpsChan <- true
		}()
	}

	if *httpFlag { // HTTP
		httpChan := make(chan bool)
		waits = append(waits, httpChan)

		go func() {
			log.Println(fmt.Sprintf("\t-Using HTTP on port %v", *httpPort))
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *httpPort), nil))

			httpChan <- true
		}()
	}

	for _, w := range waits {
		<-w // Wait for all the ListenAndServe routines to finish
	}
}

// Handler to catch 404s. Notes the 404 in the log.
func Debug404(w http.ResponseWriter, r *http.Request) {
	log.Println("404 Served")
	http.NotFound(w, r)
}

// Returns the same handler, but intercepts the request first and logs the Method and URL.
func ProbeHandler(level int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/log" {
			if level > 1 {
				log.Println("Headers:")
				for key, value := range r.Header {
					log.Println("\t", key, "->", value)
				}
			}
			if level > 0 {
				log.Println(r.Method, r.URL)
			}
		}
		h.ServeHTTP(w, r)
	})
}

// Set the log output file based on flag or env variable if available. (Flag takes precedence).
func setupLogger() {
	if *logfile == "" { // flag takes precendence over ENV variable
		*logfile = os.Getenv("RTER_LOGFILE")
	}

	if *logfile != "" {
		file, err := os.OpenFile(*logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)

		if err == nil {
			log.SetOutput(file)
		} else {
			log.Println(err)
		}
	}
}

// Set the log output file based on flag or env variable if available. (Flag takes precedence).
func setupRterDir() {
	if *rterDir == "" { // flag takes precendence over ENV variable
		*rterDir = os.Getenv("RTER_DIR")
	}
}
