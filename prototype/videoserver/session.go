// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

import (
//	"io/ioutil"
	"net/http"
	"runtime"
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
	Type        int         // ingest type id TRANSCODE_TYPE_XX
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

func (s *TranscodeSession) Open(t int) *ServerError {

	if s.IsOpen() { return nil }

	s.Type = t
	s.Command = BuildTranscodeCommand(s)
	log.Printf("Opening session: %s", s.Command)
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

func (s *TranscodeSession) ValidateRequest(r *http.Request) *ServerError {

	// check for proper mime type
	if !IsMimeTypeValid(s.Type, r.Header.Get("Content-Type")) {
		return ServerErrorWrongMimetype
	}

	// check content

	return nil
}

func (s *TranscodeSession) Write(r *http.Request) *ServerError {

	if !s.IsOpen() { return ServerErrorTranscodeFailed }

	log.Printf("Writing data to session %d", s.UID)

	// check request compatibility (mime type, content)
	if err := s.ValidateRequest(r); err != nil { return err }


	// TODO: push data into pipe


	// old code below (appending to file)
	idstr := strconv.FormatUint(s.UID, 10)
	filename := c.Transcode.Mp4.Path + "/" + idstr + ".h264"
    f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    io.Copy(f, r.Body)
    r.Body.Close()
    f.Close()

	return nil
}

func (s *TranscodeSession) HandleTimeout() {
	log.Printf("Session timeout")
	s.Close()
}

