package streaming

import (
	"github.com/gorilla/mux"
	"github.com/igm/sockjs-go/sockjs"
	"net/http"
	"strings"
)

type Debugable interface {
	Debug(bool)
}

type StreamingRouter struct {
	*mux.Router

	debug     bool
	streamers []Debugable
}

func (sr *StreamingRouter) Debug(en bool) {
	sr.debug = en

	for _, s := range sr.streamers {
		s.Debug(en)
	}
}

func NewStreamingRouter() *StreamingRouter {
	r := mux.NewRouter().StrictSlash(true)

	genericStreamer := NewGenericStreamer() //TODO: This is never closed
	r.PathPrefix("/{datatype:items|users|roles|taxonomy}/{key}/{childtype:comments|ranking|direction}").HandlerFunc(
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

	sr := new(StreamingRouter)
	sr.Router = r
	sr.streamers = []Debugable{
		genericStreamer,
	}

	return sr
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
