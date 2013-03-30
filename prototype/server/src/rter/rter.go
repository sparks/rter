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
	"rter/utils"
)

var (
	probeLevel   = flag.Int("probe", 0, "perform logging on requests")
	httpsFlag    = flag.Bool("https", false, "enable https")
	httpFlag     = flag.Bool("http", true, "enable http")
	gzipFlag     = flag.Bool("gzip", false, "enable gzip compression")
	logfile      = flag.String("log-file", "", "set server logfile")
	serveLogFlag = flag.Bool("serve-log-file", true, "serve logfile over http")
	httpPort     = flag.Int("http-port", 8080, "set the http port to use")
	httpsPort    = flag.Int("https-port", 10433, "set the https port to use")
)

func main() {
	flag.Parse()

	setupLogger()

	log.Println("Launching rtER Server")

	err := storage.OpenStorage("rter", "j2pREch8", "tcp", "localhost:3306", "rter")

	if err != nil {
		log.Fatalf("Failed to open connection to database %v", err)
	}
	defer storage.CloseStorage()

	r := mux.NewRouter().StrictSlash(true)

	s := streaming.StreamingRouter()
	r.PathPrefix("/1.0/streaming").Handler(http.StripPrefix("/1.0/streaming", s)) //Must register more specific paths first

	crud := rest.CRUDRouter()
	r.PathPrefix("/1.0").Handler(http.StripPrefix("/1.0", crud))

	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads", http.FileServer(http.Dir(utils.UploadPath))))

	r.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir(filepath.Join(utils.WWWPath, "css")))))
	r.PathPrefix("/js").Handler(http.StripPrefix("/js", http.FileServer(http.Dir(filepath.Join(utils.WWWPath, "js")))))
	r.PathPrefix("/vendor").Handler(http.StripPrefix("/vendor", http.FileServer(http.Dir(filepath.Join(utils.WWWPath, "vendor")))))
	r.PathPrefix("/asset").Handler(http.StripPrefix("/asset", http.FileServer(http.Dir(filepath.Join(utils.WWWPath, "asset")))))
	r.PathPrefix("/template").Handler(http.StripPrefix("/template", http.FileServer(http.Dir(filepath.Join(utils.WWWPath, "template")))))

	r.HandleFunc("/favicon.ico",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(utils.WWWPath, "asset", "favicon.ico"))
		},
	)

	if *serveLogFlag && *logfile != "" { //Sorta hacky this also depends on the setupLogger running before in case envvar was set
		log.Println("\t-Serve Log Enabled")
		r.HandleFunc("/log",
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, *logfile)
			},
		).Methods("GET")
	}

	r.HandleFunc("/auth", auth.AuthHandlerFunc).Methods("POST")
	r.HandleFunc("/multiup", legacy.MultiUploadHandler)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(utils.WWWPath, "new.html"))
	})

	r.NotFoundHandler = http.HandlerFunc(debug404)

	var rootHandler http.Handler = r

	if *gzipFlag {
		log.Print("\t-GZIP Enabled")
		rootHandler = compressor.GzipHandler(rootHandler)
	}

	if *probeLevel > 0 {
		log.Print("\t-Probe Enabled, Level ", *probeLevel)
		rootHandler = ProbeHandler(*probeLevel, rootHandler)
	}

	http.Handle("/", rootHandler)

	waits := make([]chan bool, 0)

	if *httpsFlag {
		httpsChan := make(chan bool)
		waits = append(waits, httpsChan)

		go func() {
			log.Println(fmt.Sprintf("\t-Using HTTPS on port %v", *httpsPort))
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%v", *httpsPort), "cert.pem", "key.pem", nil))

			httpsChan <- true
		}()
	}

	if *httpFlag {
		httpChan := make(chan bool)
		waits = append(waits, httpChan)

		go func() {
			log.Println(fmt.Sprintf("\t-Using HTTP on port %v", *httpPort))
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *httpPort), nil))

			httpChan <- true
		}()
	}

	for _, w := range waits {
		<-w
	}
}

func debug404(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

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

func setupLogger() {
	if *logfile == "" { //flag takes precendence over ENV variable
		*logfile = os.Getenv("RTER_LOGFILE")
	}

	if *logfile != "" {
		file, err := os.Create(*logfile)

		if err == nil {
			log.SetOutput(file)
		} else {
			log.Println(err)
		}
	}
}
