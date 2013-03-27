package streaming

import (
	"encoding/json"
	"github.com/igm/sockjs-go/sockjs"
	"log"
	"rter/data"
	"rter/storage"
)

type ItemBundle struct {
	Action string
	Item   *data.Item
}

type ItemStreamer struct {
	bundleChannels []chan *ItemBundle
}

func NewItemStreamer() *ItemStreamer {
	s := new(ItemStreamer)
	s.bundleChannels = make([]chan *ItemBundle, 0)

	storage.AddListener(s)

	return s
}

func (s *ItemStreamer) InsertEvent(val interface{}) {
	switch v := val.(type) {
	case *data.Item:
		b := new(ItemBundle)
		b.Action = "create"
		b.Item = v

		s.Dispatch(b)
	}
}

func (s *ItemStreamer) UpdateEvent(val interface{}) {
	switch v := val.(type) {
	case *data.Item:
		b := new(ItemBundle)
		b.Action = "update"
		b.Item = v

		s.Dispatch(b)
	}
}

func (s *ItemStreamer) DeleteEvent(val interface{}) {
	switch v := val.(type) {
	case *data.Item:
		b := new(ItemBundle)
		b.Action = "delete"
		b.Item = v

		s.Dispatch(b)
	}
}

func (s *ItemStreamer) Dispatch(b *ItemBundle) {
	for _, l := range s.bundleChannels {
		l <- b
	}
}

func (s *ItemStreamer) Close() {
	storage.RemoveListener(s)
}

func (s *ItemStreamer) SockJSHandler(session sockjs.Conn) {
	localChan := make(chan *ItemBundle)

	s.bundleChannels = append(s.bundleChannels, localChan)

	go func() {
		for {
			bundle, ok := <-localChan
			if !ok { //Chanel was closed
				break
			}
			json, err := json.Marshal(bundle)

			if err != nil {
				log.Println(err)
				continue //Keep trying!
			}
			_, err = session.WriteMessage(json)
			if err != nil {
				log.Println(err)
				break //Assume connection has died
			}
		}
	}()

	for { //This is needed to catch the closure of the sock
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
