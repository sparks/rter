// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"videoserver/auth"
	"videoserver/config"
)

// --------------------------------------------------------------------------
// server state

type State struct {
	c              *config.ServerConfig
	activeSessions map[uint64]*Session
	closedSessions map[uint64]*time.Timer
}

// instantiate new server
func NewServer(conf *config.ServerConfig) *State {
	return &State{
		c:              conf,
		activeSessions: make(map[uint64]*Session),
		closedSessions: make(map[uint64]*time.Timer),
	}
}

type Error struct {
	code   int
	status int
	msg    string
}

func NewError(c int, s int, m string) *Error {
	return &Error{code: c, status: s, msg: m}
}

func (e *Error) Error() string { return e.msg }
func (e *Error) Code() int     { return e.code }
func (e *Error) Status() int   { return e.status }

func (e *Error) JSONError() string {
	return "{\n  \"errors\": [\n    {\n      \"code\": " +
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
	ErrorBadID             = NewError(1, http.StatusForbidden, "malformed id")
	ErrorQuotaExceeded     = NewError(2, http.StatusForbidden, "too many active server sessions")
	ErrorWrongMimetype     = NewError(3, http.StatusUnsupportedMediaType, "wrong MIME type for enpoint")
	ErrorTranscodeFailed   = NewError(4, http.StatusForbidden, "transcoder write on closed pipe")
	ErrorEOS               = NewError(5, http.StatusForbidden, "already at end-of-stream state")
	ErrorIO                = NewError(6, http.StatusForbidden, "storage failed")
	ErrorWrongEndpointType = NewError(7, http.StatusUnsupportedMediaType, "type mismatch between open session and ingest endpoint")
	ErrorRequestInProgress = NewError(7, http.StatusForbidden, "a request for this session is already in progress")
	ErrorInvalidClient     = NewError(7, http.StatusForbidden, "endpoint is already locked to another consumer")

	ErrorAuthTokenRequired = NewError(100, http.StatusUnauthorized, "authorization token required for this endpoint")
	ErrorAuthTokenExpired  = NewError(101, http.StatusUnauthorized, "authorization token expired")
	ErrorAuthTokenInvalid  = NewError(102, http.StatusUnauthorized, "authorization token invalid")
	ErrorAuthNoPermission  = NewError(103, http.StatusForbidden, "no permission on this endpoint")
	ErrorAuthUrlMismatch   = NewError(104, http.StatusForbidden, "request and token URL mismatch")
	ErrorAuthBadSignature  = NewError(105, http.StatusForbidden, "bad signature in authorization token")
)

//	ErrSessionQuotaExceded

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
func (s *State) FindOrCreateSession(idstr string, t int) (*Session, *Error) {

	// todo: lock? are http handlers called concurrently? maybe use channel
	// what happens if a handler is called while another is running on the same video

	id, err := strconv.ParseUint(idstr, 10, 64)

	if err != nil {
		log.Printf("Malformed id: expected number, got `%s`", idstr)
		return nil, ErrorBadID
	}

	// ensure uniqueness (session id is non-closed and non-failed)
	if _, found := s.closedSessions[id]; found {
		log.Printf("Session %d already at EOS", id)
		return nil, ErrorEOS
	}

	// look up session id in map of active sessions
	session, found := s.activeSessions[id]

	// for new sessions check quota before creating an entry
	if !found {
		if uint64(len(s.activeSessions)) < s.c.Limits.Max_ingest_sessions {
			log.Printf("Creating New Session for id=%d", id)
			session = NewTranscodeSession(s, s.c, id)
			s.activeSessions[id] = session
		} else {
			log.Printf("Too many active sessions")
			return nil, ErrorQuotaExceeded
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

func (s *State) SessionUpdate(id uint64, state int) {

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
			time.AfterFunc(time.Duration(s.c.Server.Session_maxage)*time.Second,
				func() { delete(s.closedSessions, id) })
		delete(s.activeSessions, id)
		return
	}
}

func AuthenticateRequest(r *http.Request, key string) *Error {

	// parse token from HTTP request Authorization header
	t, err := auth.NewTokenFromHttpRequest(r)
	if err != nil {
		if t == nil {
			return ErrorAuthTokenRequired
		} else {
			return ErrorAuthTokenInvalid
		}
	}

	// tokens issued for a resource are valid for all sub-resources
	if !strings.HasPrefix(r.URL.String(), t.Resource) {
		return ErrorAuthUrlMismatch
	}

	// tokens are only valid up to a specified liftime
	if t.VerifyLifetime() != nil {
		return ErrorAuthTokenExpired
	}

	// tokens are only valid when their signature is correct
	if t.VerifySignature(key) != nil {
		return ErrorAuthBadSignature
	}

	return nil
}
