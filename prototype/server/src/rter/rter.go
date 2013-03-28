package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rter/auth"
	"rter/compressor"
	"rter/mobile"
	"rter/rest"
	"rter/storage"
	"rter/streaming"
	"rter/utils"
	"rter/web"
)

func main() {
	setupLogger()

	probe := flag.Int("probe", 0, "probe and log Method and URL for every request")
	https := flag.Bool("https", false, "use https")
	gzip := flag.Bool("gzip", false, "enable gzip compression")

	flag.Parse()

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

	r.HandleFunc("/auth", auth.AuthHandlerFunc).Methods("POST")

	r.HandleFunc("/multiup", mobile.MultiUploadHandler)

	r.HandleFunc("/submit", web.SubmitHandler).Methods("POST")
	r.HandleFunc("/submit",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(utils.WWWPath, "submit.html"))
		},
	).Methods("GET")

	r.PathPrefix("/ajax").HandlerFunc(web.ClientAjax)

	r.HandleFunc("/", web.ClientHandler)
	r.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(utils.WWWPath, "new.html"))
	})

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

	r.NotFoundHandler = http.HandlerFunc(debug404)

	log.Println("Launching rtER Server")

	var rootHandler http.Handler = r

	if *gzip {
		log.Print("\t-GZIP Enabled")
		rootHandler = compressor.GzipHandler(rootHandler)
	}

	if *probe > 0 {
		log.Print("\t-Probe Enabled, Level ", *probe)
		rootHandler = ProbeHandler(*probe, rootHandler)
	}

	http.Handle("/", rootHandler)

	if *https {
		log.Println("\t-Using HTTPS")
		log.Fatal(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
	} else {
		log.Println("\t-Using HTTP")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}

func debug404(w http.ResponseWriter, r *http.Request) {
	log.Println("404 Served")
	log.Println(r.Method, r.URL)
	http.NotFound(w, r)
}

func ProbeHandler(level int, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if level > 1 {
			log.Println("Headers:")
			for key, value := range r.Header {
				log.Println("\t", key, "->", value)
			}
		}
		if level > 0 {
			log.Println(r.Method, r.URL)
		}
		h.ServeHTTP(w, r)
	})
}

func setupLogger() {
	logOutputFile := os.Getenv("RTER_LOGFILE")

	if logOutputFile != "" {
		logFile, err := os.Create(logOutputFile)

		if err == nil {
			log.SetOutput(logFile)
		} else {
			log.Println(err)
		}
	}
}
