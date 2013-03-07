package main

import (
//	"io/ioutil"
	"net/http"
	"runtime"
	"errors"
	"strconv"
	"log"
	"os"
	"io"
)

// Transcoder states
const (
	TC_INIT		int = 0
	TC_RUNNING  int = 1
	TC_EOS      int = 2
	TC_FAILED	int = 3
)

type TranscodeSession struct {
	UID			uint64		// video UID
	State     	int         // state
	PID			uint64		// id of transcoder (ffmpeg) process
	//Pipe  		?  			// IO channel to transcoder
	Command		string  	// command line used to start transcoder

	// Statistics
	BytesIn 	uint64      // total number of bytes received in requests
	BytesOut 	uint64      // total number of bytes forwarded to transcoder
	CallsIn 	uint64      // total number of times the ingest handler was called
	//CPUusage
	//ExitStatus
}

// transcode session errors
var (
	ErrTranscodeFailed = errors.New("Transcode process failed")
)

func NewTranscodeSession(id uint64) *TranscodeSession {
	log.Printf("Session constructor")

	s := TranscodeSession{
		UID: id,
	}

	// make sure Close is properly called
	runtime.SetFinalizer(&s, (*TranscodeSession).Close)
	return &s
}

func (s *TranscodeSession) IsOpen() bool {
	return s.State == TC_RUNNING
}

func (s *TranscodeSession) Open(params string) *ServerError {

	if s.IsOpen() { return nil }
	log.Printf("Opening session: %s", params)
	// create pipe

	// start transcode process

	// set timeout

	// set state
	s.State = TC_RUNNING
	return nil
}

func (s *TranscodeSession) Close() *ServerError {

	if !s.IsOpen() { return nil }
	log.Printf("Closing session")

	// kill transcode process

	// close pipe

	// set state
	s.State = TC_EOS
	return nil
}

func (s *TranscodeSession) Write(d io.ReadCloser) *ServerError {

	if !s.IsOpen() { return NewServerError("Transcoder write on closed pipe", 1, http.StatusForbidden) }

	log.Printf("Writing data to session %d", s.UID)

	// push data into pipe
	idstr := strconv.FormatUint(s.UID, 10)
	filename := c.Paths.Data_storage_path + "/" + idstr + ".h264"
    f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    io.Copy(f, d)
    d.Close()
    f.Close()

	return nil
}

func (s *TranscodeSession) HandleTimeout() {
	log.Printf("Session timeout")
	s.Close()
}



// --------------------------------------------------------------------------
// server state

type ServerState struct {
	activeSessions map[uint64]*TranscodeSession
}

// server errors
//var (
//	ErrSessionQuotaExceded = errors.New()
//)

// global server state
var server = ServerState {
	activeSessions: make(map[uint64]*TranscodeSession),
}

// Returns an active transcoding session for the requested video id
//
// ensures
// - session quota limit is kept
// - session is unique (later across a cluster of ingest servers)
// - session has not been closed before (stream is already in EOS state)
// - session is active and transcoder is running
func (s *ServerState) FindOrCreateSession(id uint64) (*TranscodeSession, *ServerError) {

	// todo: lock? are http handlers called concurrently? maybe use channel
	// what happens if a handler is called while another is running on the same video

	// look up session id in map of active sessions
	session, found := s.activeSessions[id]

	// for new sessions check quota before creating an entry
	if !found {
		if uint64(len(s.activeSessions)) < c.Limits.Max_ingest_sessions {
			log.Printf("Creating New Session for id=%d", id)
			session = NewTranscodeSession(id)
			s.activeSessions[id] = session
		} else {
			log.Printf("Too many active sessions")
			return nil, NewServerError("Session: too many active server sessions",
								       2, http.StatusForbidden)
		}
	}

	return session, nil
}

