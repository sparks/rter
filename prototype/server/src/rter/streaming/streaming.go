package streaming

import (
	"github.com/gorilla/mux"
	"github.com/igm/sockjs-go/sockjs"
	"log"
	"net/http"
	"strings"
)

func StreamingRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	genericStreamer := NewGenericStreamer() //TODO: This is never closed
	r.PathPrefix("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:ranking|direction}").HandlerFunc(
		// TODO: Is there a better less weird way of doing this dynamic binding of the websockets
		func(w http.ResponseWriter, r *http.Request) {
			GenericStreamingHandler(genericStreamer, w, r)
		},
	)

	r.PathPrefix("/{datatype:items|users|roles|taxonomy}").HandlerFunc(
		// TODO: Is there a better less weird way of doing this dynamic binding of the websockets
		func(w http.ResponseWriter, r *http.Request) {
			GenericStreamingHandler(genericStreamer, w, r)
		},
	)

	return r
}

func GenericStreamingHandler(g *GenericStreamer, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Build a URI like representation of the datatype
	types := []string{vars["datatype"]}

	if key, ok := vars["key"]; ok {
		types = append(types, key)
	}

	if childtype, ok := vars["childtype"]; ok {
		types = append(types, childtype)
	}

	crudstring := strings.Join(types, "/")

	sockHandler := sockjs.NewRouter("/", func(session sockjs.Conn) { g.SockJSHandler(crudstring, session) }, sockjs.DefaultConfig)
	(http.StripPrefix("/"+crudstring, sockHandler)).ServeHTTP(w, r)

}
