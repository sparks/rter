package streaming

import (
	"encoding/json"
	"github.com/igm/sockjs-go/sockjs"
	"log"
	"rter/data"
	"rter/storage"
)

type Bundle struct {
	Action string
	Val    data.CRUDable
}

type GenericStreamer struct {
	bundleChannels []chan *Bundle
}

func NewGenericStreamer() *GenericStreamer {
	s := new(GenericStreamer)
	s.bundleChannels = make([]chan *Bundle, 0)

	storage.AddListener(s)

	return s
}

func (s *GenericStreamer) InsertEvent(val interface{}) {
	c, ok := val.(data.CRUDable)

	if !ok {
		return
	}

	b := new(Bundle)
	b.Action = "create"
	b.Val = c

	s.Dispatch(b)
}

func (s *GenericStreamer) UpdateEvent(val interface{}) {
	c, ok := val.(data.CRUDable)

	if !ok {
		return
	}

	b := new(Bundle)
	b.Action = "update"
	b.Val = c

	s.Dispatch(b)
}

func (s *GenericStreamer) DeleteEvent(val interface{}) {
	c, ok := val.(data.CRUDable)

	if !ok {
		return
	}

	b := new(Bundle)
	b.Action = "delete"
	b.Val = c

	s.Dispatch(b)
}

func (s *GenericStreamer) Dispatch(b *Bundle) {
	for _, l := range s.bundleChannels {
		l <- b
	}
}

func (s *GenericStreamer) Close() {
	storage.RemoveListener(s)
}

func (s *GenericStreamer) SockJSHandler(crudpath string, session sockjs.Conn) {
	localChan := make(chan *Bundle)

	s.bundleChannels = append(s.bundleChannels, localChan)

	go func() {
		for {
			bundle, ok := <-localChan
			if !ok { // Chanel was closed
				break
			}

			if crudpath != bundle.Val.CRUDPrefix() {
				continue
			}

			json, err := json.Marshal(bundle)

			if err != nil {
				log.Println(err)
				continue // Keep trying!
			}
			_, err = session.WriteMessage(json)
			if err != nil {
				log.Println(err)
				break // Assume connection has died
			}
		}
	}()

	for { // This is needed to catch the closure of the sock
		_, err := session.ReadMessage()
		if err != nil {
			break
		}
	}

	for i, sliceChan := range s.bundleChannels {
		if sliceChan == localChan {
			s.bundleChannels = append(s.bundleChannels[0:i], s.bundleChannels[i+1:len(s.bundleChannels)]...)
			break
		}
	}

	close(localChan)
}
