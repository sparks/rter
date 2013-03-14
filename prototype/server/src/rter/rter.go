package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rter/mobile"
	"rter/rest"
	"rter/storage"
	"rter/utils"
	"rter/web"
)

func main() {
	setupLogger()

	err := storage.OpenStorage("rter", "j2pREch8", "tcp", "localhost:3306", "rter")
	if err != nil {
		log.Fatalf("Failed to open connection to database %v", err)
	}
	defer storage.CloseStorage()

	r := mux.NewRouter().StrictSlash(true)

	crud := rest.CRUDRouter()
	r.PathPrefix("/1.0").Handler(http.StripPrefix("/1.0", crud))

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

	// r.NotFoundHandler = http.HandlerFunc(rootRedirect)

	http.Handle("/", r)

	log.Println("Launching rtER Server")
	// log.Fatal(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
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
