// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

import (
	"text/template"
	"bytes"
	"log"
)

const (
	TC_INGEST_UNKNOWN 	int = 0
	TC_INGEST_AVC 		int = 1
	TC_INGEST_TS  		int = 2
	TC_INGEST_CHUNK		int = 3
)

const (
	TC_TARGET_HLS 		int = 0
	TC_TARGET_DASH 		int = 1
	TC_TARGET_MP4  		int = 3
	TC_TARGET_OGG  		int = 4
	TC_TARGET_WBEM  	int = 5
	TC_TARGET_THUMB  	int = 6
	TC_TARGET_POSTER  	int = 7
)

// mix server config and transcode session to make it accessible for template expansion
type TemplateData struct {
	C *ServerConfig
	S *TranscodeSession
}

const (
	TC_ARG_TS2HLS string = " -vsync 2 -copyts -copytb 1 -codec copy -map 0 -f segment -segment_time 2 -segment_format mpegts -segment_list_flags +live -segment_list {{.C.Transcode.Hls.Path}}/{{.S.UID}}.m3u8  {{.C.Transcode.Hls.Path}}/{{.S.UID}}-%09d.ts"
	TC_ARG_AVC2HLS string = " -vsync 0 -copyts -copytb 1 -codec copy -map 0 -f segment -segment_time 2 -segment_format mpegts -segment_list_flags +live -segment_list {{.C.Transcode.Hls.Path}}/{{.S.UID}}.m3u8  {{.C.Transcode.Hls.Path}}/{{.S.UID}}-%09d.ts"
	TC_ARG_DASH string = ""
	TC_ARG_MP4 string = " -codec copy {{.C.Transcode.Mp4.Path}}/{{.S.UID}}.mp4 "
	TC_ARG_OGG string = " -codec:v libtheora -b:v 600k -codec:a libvorbis -b:a 128k {{.C.Transcode.Ogg.Path}}/{{.S.UID}}.ogv "
	TC_ARG_WBEM string = " -codec:v libvpx -quality realtime -cpu-used 0 -b:v 600k -qmin 10 -qmax 42 -maxrate 600k -bufsize 1000k -threads 1 -codec:a libvorbis -b:a 128k -f webm {{.C.Transcode.Webm.Path}}/{{.S.UID}}.webm "
	TC_ARG_THUMB string = " -vsync 1 -r 0.5 -f image2 -s 160x90 {{.C.Transcode.Thumb.Path}}/thumb-{{.S.UID}}-%09d.jpg "
	TC_ARG_POSTER string = " -vsync 1 -r 0.5 -f image2 {{.C.Transcode.Poster.Path}}/poster-{{.S.UID}}-%09d.jpg "
)
// -ss 00:00:10 -vframes 1

const (
	TC_CMD_START_PROD string = "-y -re -v quiet -fflags nobuffer -i pipe:0 "
	TC_CMD_START_DEV string = "-y -re -v debug -fflags nobuffer -i pipe:0 "
	TC_CMD_END_PROD string = ""
	TC_CMD_END_DEV string = ""
)


func IsMimeTypeValid(t int, m string) bool {
	// first letter and letter after hyphen uppercase, rest lowercase
	//contentType := http.CanonicalHeaderKey(m)

	return true


}


// assemble transcode command (for now we run a single transcoder 'ffmpeg')
func BuildTranscodeCommand(s *TranscodeSession) string {
	//return "ffmpeg --help"
	var cmd string

	// transcoder command (generate debug output when in dev mode)
	if c.Server.Production_mode { cmd = TC_CMD_START_PROD
	} else { cmd = TC_CMD_START_DEV }

	// segment file formats
	if c.Transcode.Hls.Enabled { cmd += TC_ARG_AVC2HLS }
	if c.Transcode.Dash.Enabled { cmd += TC_ARG_DASH }

	// full file formats
	if c.Transcode.Mp4.Enabled { cmd += TC_ARG_MP4 }
	if c.Transcode.Ogg.Enabled { cmd += TC_ARG_OGG }
	if c.Transcode.Webm.Enabled { cmd += TC_ARG_WBEM }

	// image formats
	if c.Transcode.Thumb.Enabled { cmd += TC_ARG_THUMB }
	if c.Transcode.Poster.Enabled { cmd += TC_ARG_POSTER }

	// end trancode command line
	if c.Server.Production_mode { cmd += TC_CMD_END_PROD
	} else { cmd += TC_CMD_END_DEV }

	// combine session and server config for access by template matcher
	var cmd_writer bytes.Buffer
	var d = TemplateData{&c, s}

	// replace placeholders with config strings
	t, err := template.New("cmd").Parse(cmd)
	if err != nil { log.Fatalf("Error parsing cmd template: %s\n", err) }

	err = t.Execute(&cmd_writer, d)
	if err != nil { log.Fatalf("Error generating cmd string: %s\n", err) }

	return cmd_writer.String()
}
