package streaming

import (
	"github.com/gorilla/mux"
	"github.com/igm/sockjs-go/sockjs"
	"log"
	"net/http"
)

func StreamingRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	itemStreamer := NewItemStreamer()
	itemSockHandler := sockjs.NewRouter("/", func(session sockjs.Conn) { itemStreamer.SockJSHandler(session) }, sockjs.DefaultConfig)
	r.PathPrefix("/items").Handler(http.StripPrefix("/items", itemSockHandler))

	termRankingStreamer := NewTermRankingStreamer()
	r.PathPrefix("/taxonomy/{term}/ranking").HandlerFunc( // TODO: Is there a better less weird way of doing this dynamic binding of the websockets
		func(w http.ResponseWriter, r *http.Request) {
			HandleTermRanking(termRankingStreamer, w, r)
		},
	)

	return r
}

func HandleTermRanking(termRankingStreamer *TermRankingStreamer, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sockHandler := sockjs.NewRouter("/", func(session sockjs.Conn) { termRankingStreamer.SockJSHandler(vars["term"], session) }, sockjs.DefaultConfig)
	(http.StripPrefix("/taxonomy/"+vars["term"]+"/ranking", sockHandler)).ServeHTTP(w, r)
}

func probe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Println("Tax")
	log.Println(vars)

	log.Println(r.Method, r.URL)
}
