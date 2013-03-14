// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

import (
	"text/template"
	"bytes"
	"log"
	"fmt"
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
	Thumb_size string
	Thumb_rate uint64
	Poster_size string
	Poster_corrected_count uint64
	Poster_rate uint64
	Poster_skip string
}

// Codec Support
//
// Make sure to compile ffmpeg with support for the following codecs if you intend to
// use them as target:
//
//   - theora: Theora video encoder for OGG
//   - vorbis: Vorbis audio encoder for OGG
//   - libogg: OGG stream format
//   - libvpx: Google's VP8 encoder
//
// OSX example
// brew install libvpx libogg libvorbis theora
// brew install ffmpeg --with-theora --with-libogg --with-libvorbis --with-libvpx

const (
	TC_ARG_HLS string = " -f segment -codec copy -map 0 -segment_time {{.C.Transcode.Hls.Segment_length}} -segment_format mpegts -segment_list_flags +live -segment_list_type hls -individual_header_trailer 1 -segment_list index.m3u8 hls/%09d.ts"
	TC_ARG_DASH string = " "
	TC_ARG_MP4 string = " -codec copy video.mp4 "
	TC_ARG_OGG string = " -codec:v libtheora -b:v 600k -codec:a vorbis -b:a 128k video.ogv "
	TC_ARG_WBEM string = " -f webm -codec:v libvpx -quality realtime -cpu-used 0 -b:v 600k -qmin 10 -qmax 42 -maxrate 600k -bufsize 1000k -threads 1 -codec:a libvorbis -b:a 128k video.webm "
	TC_ARG_THUMB string = " -f image2 {{.Thumb_size}} -vsync 1 -vf fps=fps=1/{{.Thumb_rate}} thumb/%09d.jpg "
	TC_ARG_POSTER string = " -f image2 {{.Poster_size}} -vsync 1 -vf fps=fps=1/{{.Poster_rate}} {{.Poster_skip}} -vframes {{.Poster_corrected_count}} poster/%09d.jpg "
)

// Used Ingest options
//
// -fflags +genpts+igndts+nobuffer
// -err_detect compliant
// -avoid_negative_ts 1 [bool]
// -correct_ts_overflow 1 [bool]
// -max_delay 500000 [microsec]
// -analyzeduration 500000 [microsec]
// -f mpegts -c:0 h264
// -vsync 0
// -copyts
// -copytb 1

// Failed options
//
// -probesize 2048 [bytes]					-- not needed since it already worked with analyzeduration
// -avioflags direct 						-- broke format detection
// -f mpegtsraw -compute_pcr 0 				-- created an invalid MPEGTS bitstream
// -use_wallclock_as_timestamps 1 [bool] 	-- broke TS timing

// Unused Options
//
// -copyinkf:0
// -fflags +discardcorrupt
// -fpsprobesize 2

const (
	TC_ARG_TSIN string = " -fflags +genpts+igndts+nobuffer -err_detect compliant -avoid_negative_ts 1 -correct_ts_overflow 1 -max_delay 500000 -analyzeduration 500000 -f mpegts -c:0 h264 -vsync 0 -copyts -copytb 1 "
	TC_ARG_AVCIN string = " -fflags +genpts+igndts -max_delay 0 -analyzeduration 0 -f h264 -c:0 h264 -copytb 0 "
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
	var d = TemplateData{&c, s, "", 1, "", 1, 1, ""}

	// Poster
	if c.Transcode.Poster.Enabled {
		// scaling parameters
		if c.Transcode.Poster.Size == "auto" || c.Transcode.Poster.Size == "" {
			d.Poster_size = ""
		} else {
			d.Poster_size = "-s " + c.Transcode.Poster.Size
		}

		// skip parameter
		if c.Transcode.Poster.Skip > 0 {
			hou := c.Transcode.Poster.Skip / 3600
			min := (c.Transcode.Poster.Skip - hou*3600) / 60
			sec := (c.Transcode.Poster.Skip - hou*3600 - min*60)
			d.Poster_skip = fmt.Sprintf("-ss %02d:%02d:%02d.0", hou, min, sec)
		}

		// step interval
		if c.Transcode.Poster.Step > 0 { d.Poster_rate = c.Transcode.Poster.Step }

		// increase the count of poster frames to output since ffmpeg flushes
		// the image2 pipeline at the end only, hence a single image would not be
		// written before the stream ends
		//
		// on the downside, this solution writes one image more than expected by the user
		//
		d.Poster_corrected_count = c.Transcode.Poster.Count + 1
	}

	// Thumbnail
	if c.Transcode.Thumb.Enabled {
		// scaling parameter
		if c.Transcode.Thumb.Size != "" {
			d.Thumb_size = "-s " + c.Transcode.Thumb.Size
		}

		// step interval
		if c.Transcode.Thumb.Step > 0 { d.Thumb_rate = c.Transcode.Thumb.Step }
	}

	// replace placeholders with config strings
	t, err := template.New("cmd").Parse(cmd)
	if err != nil { log.Fatalf("Error parsing cmd template: %s\n", err) }

	err = t.Execute(&cmd_writer, d)
	if err != nil { log.Fatalf("Error generating cmd string: %s\n", err) }

	return cmd_writer.String()
}

func CheckTranscoderCapabilities() {

	// HLS format
 	//ffmpeg -formats -v quiet | grep " segment  "
 	// '  E segment         segment'

	// MPEG2TS muxer and demuxer
	//ffmpeg -formats -v quiet | grep "mpegts "
	// ' DE mpegts          MPEG-TS (MPEG-2 Transport Stream)'

	// MP4 format
	// ffmpeg -formats -v quiet | grep " mp4 "
	// '  E mp4             MP4 (MPEG-4 Part 14)''

	// OGG format (Theora, Vorbis encoders)
	// ffmpeg -formats -v quiet | grep "ogg"
	// ' DE ogg             Ogg'
	// ffmpeg -codecs -v quiet | grep theora
 	// ' DEV.L. theora               Theora (encoders: libtheora )'
	// ffmpeg -codecs -v quiet | grep vorbis
 	// ' DEA.L. vorbis               Vorbis (decoders: vorbis libvorbis ) (encoders: vorbis libvorbis )''

	// Webm format (libvpx, Vorbis encoders)
	//ffmpeg -formats -v quiet | grep webm | grep E
	// '  E webm            WebM'
 	// ffmpeg -codecs -v quiet  | grep vp8
 	// ' DEV.L. vp8                  On2 VP8 (decoders: vp8 libvpx ) (encoders: libvpx )'
	// ffmpeg -codecs -v quiet | grep vorbis
 	// ' DEA.L. vorbis               Vorbis (decoders: vorbis libvorbis ) (encoders: vorbis libvorbis )''

	// JPEG format (jpeg encoder)
	// ffmpeg -formats -v quiet | grep " mjpeg"
	// ' DE mjpeg           raw MJPEG video'
	// ffmpeg -codecs -v quiet | grep " mjpeg "
	// ' DEVIL. mjpeg                Motion JPEG'

	// H264 format and decoder
	// ffmpeg -formats -v quiet  | grep h264
 	// ' DE h264            raw H.264 video'
	// ffmpeg -codecs -v quiet  | grep h264
 	// ' DEV.LS h264                 H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10 (decoders: h264 h264_vda ) (encoders: libx264 libx264rgb )'

}
