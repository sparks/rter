package streaming

import (
	"github.com/gorilla/mux"
	"github.com/igm/sockjs-go/sockjs"
)

func StreamingRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	itemStreamer := NewItemStreamer()
	sockHandler := sockjs.NewRouter("/items", func(session sockjs.Conn) { itemStreamer.SockJSHandler(session) }, sockjs.DefaultConfig)
	r.PathPrefix("/items").Handler(sockHandler)

	return r
}
