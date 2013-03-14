// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

import (
	"net/http"
	"strconv"
	"time"
	"fmt"
	"log"
)

// --------------------------------------------------------------------------
// server state

type ServerState struct {
	activeSessions map[uint64]*TranscodeSession
	closedSessions map[uint64]*time.Timer
}

// instantiate global state variable
var server = ServerState {
	activeSessions: make(map[uint64]*TranscodeSession),
	closedSessions: make(map[uint64]*time.Timer),
}


type ServerError struct {
	code   int
	status int
	msg    string
}

func NewServerError(c int, s int, m string) *ServerError {
	return &ServerError{code: c, status: s, msg: m}
}

func (e *ServerError) Error() string { return e.msg }
func (e *ServerError) Code() int { return e.code }
func (e *ServerError) Status() int { return e.status }

func (e *ServerError) JSONError() string {
	return  "{\n  \"errors\": [\n    {\n      \"code\": " +
	        strconv.Itoa(e.code) +
	        ",\n      \"message\": \"" +
	        e.msg +
	        "\"\n    }\n  ]\n}"
}

func ServeError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
   	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}

// new HTTP status codes not defined in net/http
const (
	StatusTooManyRequests = 429
)

// Server Error Reasons
var (
	ServerErrorBadID           = NewServerError(1, http.StatusForbidden, "malformed id")
	ServerErrorQuotaExceeded   = NewServerError(2, http.StatusForbidden, "too many active server sessions")
	ServerErrorWrongMimetype   = NewServerError(3, http.StatusUnsupportedMediaType, "wrong MIME type for enpoint")
	ServerErrorTranscodeFailed = NewServerError(4, http.StatusForbidden, "transcoder write on closed pipe")
	ServerErrorEOS             = NewServerError(5, http.StatusForbidden, "already at end-of-stream state")
	ServerErrorIO              = NewServerError(6, http.StatusForbidden, "storage failed")
)

//	ErrSessionQuotaExceded
//	ServerErrorClientTimeout

// AUTH
// auth failure: token expired, token invalid, no permissions on endpoint
// INGEST
// too many active server sessions (total and per source)
// bandwidth rate limit exceded for source
// bitstream invalid
// DOWNLOAD
// unknown resource id (stream, segment, thumb, poster) -> 404
//


// Returns an active transcoding session for the requested video id
//
// ensures
// - session quota limit is kept
// - session is unique (later across a cluster of ingest servers)
// - session has not been closed before (stream is already in EOS state)
// - session is active and transcoder is running
//
func (s *ServerState) FindOrCreateSession(idstr string, t int) (*TranscodeSession, *ServerError) {

	// todo: lock? are http handlers called concurrently? maybe use channel
	// what happens if a handler is called while another is running on the same video

 	id, err := strconv.ParseUint(idstr, 10, 64)

 	if err != nil {
 		log.Printf("Malformed id: expected number, got `%s`", idstr)
 		return nil, ServerErrorBadID
 	}

 	// ensure uniqueness (session id is non-closed and non-failed)
 	if _, found := s.closedSessions[id]; found {
			log.Printf("Session %d already at EOS", id)
			return nil, ServerErrorEOS
 	}

	// look up session id in map of active sessions
	session, found := s.activeSessions[id]

	// for new sessions check quota before creating an entry
	if !found {
		if uint64(len(s.activeSessions)) < c.Limits.Max_ingest_sessions {
			log.Printf("Creating New Session for id=%d", id)
			session = NewTranscodeSession(s, id)
			s.activeSessions[id] = session
		} else {
			log.Printf("Too many active sessions")
			return nil, ServerErrorQuotaExceeded
		}
	}

	// open new sessions with specified type of request endpoint
	if !session.IsOpen() {
		if err := session.Open(t); err != nil {
			// return error
			return nil, err
		}
	}

	return session, nil
}

func (s *ServerState) SessionUpdate(id uint64, state int) {

	switch state {
		default:
			// fail on unknonw states
			log.Fatal("Unhandled Session State %d", state)
		case TC_INIT, TC_RUNNING:
			// session create is already handled in FindOrCreateSession()
			return
		case TC_FAILED, TC_EOS:
			// here we only have to deal with session shutdown

			// store self-deleting entry
			s.closedSessions[id] =
				time.AfterFunc(time.Duration(c.Server.Session_maxage) * time.Second,
	 						   func() { delete(s.closedSessions, id) })
			delete(s.activeSessions, id)
			return
	}
}
