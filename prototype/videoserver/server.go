// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

import (
	"net/http"
	"log"
)

// --------------------------------------------------------------------------
// server state

type ServerState struct {
	activeSessions map[uint64]*TranscodeSession
}

// Server Error Reasons
//var (
//	ErrSessionQuotaExceded = errors.New()
//)
//const (
//	SOURCE_TIMEOUT
//   STREAM_EOS // no error
//    TRANSCODER_FAILED
//)

// AUTH
// auth failure: token expired, token invalid, no permissions on endpoint
// INGEST
// too many active server sessions (total and per source)
// bandwidth rate limit exceded for source
// bitstream invalid
// ingest stream already ended
// DOWNLOAD
// unknown resource id (stream, segment, thumb, poster)
//

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

