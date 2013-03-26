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
	sockHandler := sockjs.NewRouter("/", func(session sockjs.Conn) { itemStreamer.SockJSHandler(session) }, sockjs.DefaultConfig)
	r.PathPrefix("/items").Handler(http.StripPrefix("/items", sockHandler))

	termRankingStreamer := NewTermRankingStreamer()
	r.PathPrefix("/taxonomy/{term}/ranking").HandlerFunc(
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

// 727	func StripPrefix(prefix string, h Handler) Handler {
// 728		return HandlerFunc(func(w ResponseWriter, r *Request) {
// 729			if !strings.HasPrefix(r.URL.Path, prefix) {
// 730				NotFound(w, r)
// 731				return
// 732			}
// 733			r.URL.Path = r.URL.Path[len(prefix):]
// 734			h.ServeHTTP(w, r)
// 735		})
// 736	}
