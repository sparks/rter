package streaming

import (
	"encoding/json"
	"github.com/igm/sockjs-go/sockjs"
	"log"
	"rter/data"
	"rter/storage"
)

type TermRankingStreamer struct {
	termRankingChannels []chan *data.TermRanking
}

func NewTermRankingStreamer() *TermRankingStreamer {
	t := new(TermRankingStreamer)
	t.termRankingChannels = make([]chan *data.TermRanking, 0)

	storage.AddListener(t)

	return t
}

func (t *TermRankingStreamer) InsertEvent(val interface{}) {}

func (t *TermRankingStreamer) UpdateEvent(val interface{}) {
	switch v := val.(type) {
	case *data.TermRanking:
		t.Dispatch(v)
	}
}

func (t *TermRankingStreamer) DeleteEvent(val interface{}) {}

func (t *TermRankingStreamer) Dispatch(r *data.TermRanking) {
	for _, l := range t.termRankingChannels {
		l <- r
	}
}

func (t *TermRankingStreamer) Close() {
	storage.RemoveListener(t)
}

func (t *TermRankingStreamer) SockJSHandler(term string, session sockjs.Conn) {
	localChan := make(chan *data.TermRanking)

	t.termRankingChannels = append(t.termRankingChannels, localChan)

	go func() {
		for {
			termRanking, ok := <-localChan
			if !ok { // Chanel was closed
				break
			}

			if termRanking.Term != term {
				continue // Ignore
			}

			json, err := json.Marshal(termRanking)

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

	for i, sliceChan := range t.termRankingChannels {
		if sliceChan == localChan {
			t.termRankingChannels = append(t.termRankingChannels[0:i], t.termRankingChannels[i+1:len(t.termRankingChannels)]...)
			break
		}
	}

	close(localChan)
}
