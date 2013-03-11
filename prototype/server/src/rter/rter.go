package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
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
	r.HandleFunc("/submit", web.SubmitHandler)

	r.PathPrefix("/ajax").HandlerFunc(web.ClientAjax)

	r.HandleFunc("/", web.ClientHandler)

	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads", http.FileServer(http.Dir(utils.UploadPath))))
	r.PathPrefix("/resources").Handler(http.StripPrefix("/resources", http.FileServer(http.Dir(utils.ResourcePath))))
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
