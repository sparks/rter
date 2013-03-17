// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Transcoder states
const (
	TC_INIT    int = 0
	TC_RUNNING int = 1
	TC_EOS     int = 2
	TC_FAILED  int = 3
)

type TranscodeSession struct {
	Server  *ServerState     // link to server used for signalling session state
	UID     uint64           // video UID
	idstr   string           // stringyfied UID
	Type    int              // ingest type id TRANSCODE_TYPE_XX
	state   int              // our state (not the state of the external process)
	Args    string           // command line arguments for transcoder
	Proc    *os.Process      // process management
	Pipe    *os.File         // IO channel to transcoder
	Pstate  *os.ProcessState // set when transcoder finished
	LogFile *os.File         // transcoder logfile
	Timer   *time.Timer      // session inactivity timer

	Consumer string // API consumer linked to the session
	live     bool   // true when request is in progress

	// Statistics
	BytesIn   int64         // total number of bytes received in requests
	BytesOut  int64         // total number of bytes forwarded to transcoder
	CallsIn   int64         // total number of times the ingest handler was called
	CpuUser   time.Duration // user-space CPU time used
	CpuSystem time.Duration // system CPU time used
}

func NewTranscodeSession(srv *ServerState, id uint64) *TranscodeSession {

	s := TranscodeSession{
		Server: srv,
		UID:    id,
		state:  TC_INIT,
		live:   false,
	}

	// stringify ID
	s.idstr = strconv.FormatUint(s.UID, 10)

	// register with server
	srv.SessionUpdate(id, TC_INIT)

	// make sure Close is properly called
	runtime.SetFinalizer(&s, (*TranscodeSession).Close)
	return &s
}

func (s *TranscodeSession) setState(state int) {
	// EOS and FAILED are final
	if s.state == TC_EOS || s.state == TC_FAILED {
		return
	}

	// set state and inform server
	s.state = state
	s.Server.SessionUpdate(s.UID, s.state)
}

func (s *TranscodeSession) IsOpen() bool {
	return s.state == TC_RUNNING
}

func (s *TranscodeSession) Open(t int) *ServerError {

	if s.IsOpen() {
		return nil
	}

	s.Type = t
	s.Args = BuildTranscodeCommand(s)
	log.Printf("Opening transcoder session: %s", s.Args)

	// create output directory structure
	if err := createOutputDirectories(s.idstr); err != nil {
		return err
	}

	// create pipe
	pr, pw, err := os.Pipe()
	if err != nil {
		s.setState(TC_FAILED)
		return ServerErrorTranscodeFailed
	}
	s.Pipe = pw

	// create logfile
	logname := c.Transcode.Log_path + "/" + s.idstr + ".log"
	s.LogFile, _ = os.OpenFile(logname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	// start transcode process
	var attr os.ProcAttr
	attr.Dir = c.Transcode.Output_path + "/" + s.idstr
	attr.Files = []*os.File{pr, s.LogFile, s.LogFile}
	s.Proc, err = os.StartProcess(c.Transcode.Command, strings.Fields(s.Args), &attr)

	if err != nil {
		log.Printf("Error starting process: %s", err)
		s.setState(TC_FAILED)
		pr.Close()
		pw.Close()
		s.LogFile.Close()
		s.Pipe = nil
		s.Type = 0
		s.Args = ""
		return ServerErrorTranscodeFailed
	}

	// close read-end of pipe and logfile after successful start
	pr.Close()
	s.LogFile.Close()

	// set timeout for session cleanup
	s.Timer = time.AfterFunc(time.Duration(c.Server.Session_timeout)*time.Second,
		func() { s.HandleTimeout() })

	// set state
	s.setState(TC_RUNNING)
	return nil
}

func (s *TranscodeSession) Close() *ServerError {

	// cancel close timeout (todo: potential race condition?)
	s.Timer.Stop()

	if !s.IsOpen() {
		return nil
	}

	// set state
	s.setState(TC_EOS)

	// close pipe
	s.Pipe.Close()

	// gracefully shut down transcode process (SIGINT, 2)
	var err error
	if err = s.Proc.Signal(syscall.SIGINT); err != nil {
		log.Printf("Sending signal to transcoder failed: %s", err)
		// assuming the transcoder process has finished
	}

	log.Printf("Waiting for transcoder shutdown")
	s.Pstate, err = s.Proc.Wait()
	if err != nil {
		log.Printf("Transcoder exited with error: %s and state %s", err, s.Pstate.String())
		return nil
	}

	log.Printf("Transcoder exit state is %s", s.Pstate.String())

	// get process statistics
	s.CpuSystem = s.Pstate.SystemTime()
	s.CpuUser = s.Pstate.UserTime()

	log.Printf("Session %d closed: %d calls, %d bytes in, %d bytes out, %s user, %s sys",
		s.UID, s.CallsIn, s.BytesIn, s.BytesOut, s.CpuUser.String(), s.CpuSystem.String())

	return nil
}

func (s *TranscodeSession) ValidateRequest(r *http.Request, t int) *ServerError {

	// on the first call store caller IP
	if s.Consumer == "" {
		s.Consumer = r.RemoteAddr
	} else if s.Consumer != r.RemoteAddr {
		// check if this is the same consumer
		return ServerErrorInvalidClient
	}

	// check if this is the only active request for this resource
	if s.live {
		return ServerErrorRequestInProgress
	}

	// check for proper mime type
	if !IsMimeTypeValid(s.Type, r.Header.Get("Content-Type")) {
		return ServerErrorWrongMimetype
	}

	// cannot mix endpoints / stream types once session is open
	if t != s.Type {
		return ServerErrorWrongEndpointType
	}

	// check content

	return nil
}

// Write Handling
//
// This function is called every time the client issues a new POST request. There
// are two alternative:
//
// - Chunked-Transfer: the client sets `Transfer-Encoding: chunked` and keeps
//     pushing new video data. In this case the function does not return and
//     handles the client request as a single transaction until the client closes
//     its transport connection or the connection fails
// - normal POST: the client sets `Content-Length: xxx` and pushes a single chunk
//     of video data (usually a frame) per request. In this case the function does
//     return after each individual POST request has been forwarded. The client
//     usually sets `Connection: keep-alive` to leave the connection open for
//     subsequent requests.
//
// Signalling End-Of-Stream condition
//
// In chunked-mode EOS is signalled by the client by dropping the transport
// connection. In normal POST mode the client can signal EOS by pushing an empty
// request (Content-Length: 0).
//
// Handling connection or client failure
//
// Golang's http framework signals a failed connection by returning EOF from
// a Reader interface which is not considered an error condition. The code below
// uses a session timeout to define when a stream is considered broken, thus
// implicitly reaching its EOS state. Before that timeout a client can try reconnecting
// and continuing stream upload.
//
// Kown Issues
// Handling timeout in chunked mode is currently not supported by Golang's http
// framework. We have to rely on cooperative clients who close their connections.

func (s *TranscodeSession) Write(r *http.Request, t int) *ServerError {

	// session must be active to perform write
	if !s.IsOpen() {
		return ServerErrorTranscodeFailed
	}

	// reset session close timeout (potential race condition with timeout handler)
	s.Timer.Stop()

	log.Printf("Writing data to session %d", s.UID)

	// check request compatibility (mime type, content)
	if err := s.ValidateRequest(r, t); err != nil {
		return err
	}

	// go live
	s.live = true

	// leave live state on exit
	defer func() { s.live = false }()

	// push data into pipe until body us empty or EOF (broken pipe)
	written, err := io.Copy(s.Pipe, r.Body)
	log.Printf("Written %d bytes to session %d", written, s.UID)

	// statitsics
	s.CallsIn++
	s.BytesIn += written
	s.BytesOut += written

	// error handling
	if err == nil && written == 0 {
		// normal session close on source request (empty body)
		log.Printf("Closing session %d down for good", s.UID)
		// close the http session
		r.Close = true
		s.Close()

	} else if err != nil {
		// session close due to broken pipe (transcoder)
		log.Printf("Closing session %d on broken pipe.", s.UID)
		s.Close()
		return ServerErrorTranscodeFailed
	}

	r.Body.Close()

	// restart timer on success
	s.Timer = time.AfterFunc(time.Duration(c.Server.Session_timeout)*time.Second,
		func() { s.HandleTimeout() })

	return nil
}

func (s *TranscodeSession) HandleTimeout() {
	s.Close()
}

func (s *TranscodeSession) SetResponseHeaders(w http.ResponseWriter) {

	// quota/rate headers
	// - available session requests
	// - available bytes for all sessions
	// - time to rate reset

}
