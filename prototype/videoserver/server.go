// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

import (
	"net/http"
	"strconv"
	"log"
)

// --------------------------------------------------------------------------
// server state

type ServerState struct {
	activeSessions map[uint64]*TranscodeSession
}

// instantiate global state variable
var server = ServerState {
	activeSessions: make(map[uint64]*TranscodeSession),
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


// new HTTP status codes not defined in net/http
const (
	StatusTooManyRequests = 429
)

// Server Error Reasons
var (
	ServerErrorBadID           = NewServerError(1, http.StatusForbidden, "malformed id")
	ServerErrorQuotaExceeded   = NewServerError(2, http.StatusForbidden, "too many active server sessions")
	ServerErrorWrongMimetype   = NewServerError(3, http.StatusUnsupportedMediaType, "wrong MIME type for enpoint")
	ServerErrorTranscodeFailed = NewServerError(4, http.StatusForbidden, "Transcoder write on closed pipe")
)

//	ErrSessionQuotaExceded
//	ServerErrorClientTimeout
//  ServerErrorEos

// AUTH
// auth failure: token expired, token invalid, no permissions on endpoint
// INGEST
// too many active server sessions (total and per source)
// bandwidth rate limit exceded for source
// bitstream invalid
// ingest stream already ended
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
func (s *ServerState) FindOrCreateSession(idstr string) (*TranscodeSession, *ServerError) {

	// todo: lock? are http handlers called concurrently? maybe use channel
	// what happens if a handler is called while another is running on the same video

 	id, err := strconv.ParseUint(idstr, 10, 64)

 	if err != nil {
 		log.Printf("Malformed id: expected number, got `%s`", idstr)
 		return nil, ServerErrorBadID
 	}

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
			return nil, ServerErrorQuotaExceeded
		}
	}

	return session, nil
}

