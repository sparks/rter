// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

import (
	"text/template"
	"bytes"
	"log"
	"os"
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
	TC_ARG_HLS string = " -f segment -codec copy -map 0 -segment_time {{.C.Transcode.Hls.Segment_length}} -segment_format mpegts -segment_list_flags +live -segment_list_type hls -individual_header_trailer 1 -segment_list index.m3u8 hls/%09d.ts"
	TC_ARG_DASH string = " "
	TC_ARG_MP4 string = " -codec copy video.mp4 "
	TC_ARG_OGG string = " -codec:v libtheora -b:v 600k -codec:a libvorbis -b:a 128k video.ogv "
	TC_ARG_WBEM string = " -codec:v libvpx -quality realtime -cpu-used 0 -b:v 600k -qmin 10 -qmax 42 -maxrate 600k -bufsize 1000k -threads 1 -codec:a libvorbis -b:a 128k -f webm video.webm "
	TC_ARG_THUMB string = " -vsync 1 -r 0.5 -f image2 -s 160x90 thumb/%09d.jpg "
	TC_ARG_POSTER string = " -vsync 1 -r 0.5 -f image2 poster/%09d.jpg "
)
// -ss 00:00:10 -vframes 1
// -f mpegtsraw -compute_pcr 0 ?
// -copyinkf:0
// -fflags +genpts // create PTS values
// -fflags 'discardcorrupt'
// -probesize 2048 (bytes)
// - analyzeduration uS
const (
	TC_ARG_TSIN string = " -fflags +nobuffer+genpts -analyzeduration 500k -f mpegts -c:0 h264 -vsync 0 -copyts -copytb 1 "
	TC_ARG_AVCIN string = " -fflags +nobuffer+genpts -probesize 1024 -f h264 -c:0 h264 -copytb 0 "
)

const (
	TC_CMD_START_PROD string = "-y -v quiet "
	TC_CMD_START_DEV string = "-y -v debug "
	TC_CMD_INPUT string = " -i pipe:0 "
	TC_CMD_END_PROD string = ""
	TC_CMD_END_DEV string = ""
)


func IsMimeTypeValid(t int, m string) bool {
	// first letter and letter after hyphen uppercase, rest lowercase
	//contentType := http.CanonicalHeaderKey(m)
	switch t {
	case TC_INGEST_TS:
		return true  // video/x-mpegts

	case TC_INGEST_AVC:
		return true //

	default:
		return false
	}
	return true
}

func createOutputDirectories(idstr string) *ServerError {

	// HLS: <hls-data-path>/<id>/hls
	if c.Transcode.Hls.Enabled {
		p := c.Transcode.Output_path + "/" + idstr + "/hls"
		if err := os.MkdirAll(p, PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return ServerErrorIO
		}
	}

	// DASH: <dash-data-path>/<id>/dash
	if c.Transcode.Dash.Enabled {
		p := c.Transcode.Output_path + "/" + idstr + "/dash"
		if err := os.MkdirAll(p, PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return ServerErrorIO
		}
	}

	// MP4: <*-data-path>/<id>
	if c.Transcode.Mp4.Enabled {
		p := c.Transcode.Output_path + "/" + idstr
		if err := os.MkdirAll(p, PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return ServerErrorIO
		}
	}

	// OGG: <*-data-path>/<id>
	if c.Transcode.Ogg.Enabled {
		p := c.Transcode.Output_path + "/" + idstr
		if err := os.MkdirAll(p, PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return ServerErrorIO
		}
	}

	// WEBM: <*-data-path>/<id>
	if c.Transcode.Webm.Enabled {
		p := c.Transcode.Output_path + "/" + idstr
		if err := os.MkdirAll(p, PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return ServerErrorIO
		}
	}

	// Thumb: <*-data-path>/<id>/thumb
	if c.Transcode.Thumb.Enabled {
		p := c.Transcode.Output_path + "/" + idstr + "/thumb"
		if err := os.MkdirAll(p, PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return ServerErrorIO
		}
	}

	// Poster: <*-data-path>/<id>/poster
	if c.Transcode.Poster.Enabled {
		p := c.Transcode.Output_path + "/" + idstr + "/poster"
		if err := os.MkdirAll(p, PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return ServerErrorIO
		}
	}

	return nil
}


// assemble transcode command (for now we run a single transcoder 'ffmpeg')
func BuildTranscodeCommand(s *TranscodeSession) string {
	//return "ffmpeg --help"
	var cmd string

	// transcoder command (generate debug output when in dev mode)
	if c.Server.Production_mode { cmd = TC_CMD_START_PROD
	} else { cmd = TC_CMD_START_DEV }

	// input spec
	switch s.Type {
	case TC_INGEST_AVC: cmd += TC_ARG_AVCIN + TC_CMD_INPUT
	case TC_INGEST_TS: cmd += TC_ARG_TSIN + TC_CMD_INPUT
	}

	// segment file formats
	if c.Transcode.Hls.Enabled { cmd += TC_ARG_HLS }
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
