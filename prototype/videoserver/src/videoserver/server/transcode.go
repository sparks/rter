// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package server

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
	"videoserver/config"
	"videoserver/utils"
)

const (
	TC_INGEST_UNKNOWN int = 0
	TC_INGEST_AVC     int = 1
	TC_INGEST_TS      int = 2
	TC_INGEST_CHUNK   int = 3
)

const (
	TC_TARGET_HLS      int = 0 // original Apple HLS streaming (H264, MPEG2-TS)
	TC_TARGET_DASH     int = 1 // DASH standard (not implemented yet)
	TC_TARGET_MP4      int = 3 // MP4 (H264, AAC) files for download
	TC_TARGET_OGG      int = 4 // OGG (Theora, Vorbis) files for download
	TC_TARGET_WBEM     int = 5 // WebM (VP8, Vorbis) files for download
	TC_TARGET_WEBM_HLS int = 6 // WebM (VP8, Vorbis) HLS streaming
	TC_TARGET_THUMB    int = 7 // JPG thumbnail images
	TC_TARGET_POSTER   int = 8 // JPG poster images
)

// mix server config and transcode session to make it accessible for template expansion
type TemplateData struct {
	C                      *config.ServerConfig
	S                      *Session
	Webm_gop_size          uint64
	Thumb_size             string
	Thumb_rate             uint64
	Poster_size            string
	Poster_corrected_count uint64
	Poster_rate            uint64
	Poster_skip            string
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
	TC_ARG_HLS      string = " -f segment -codec copy -map 0 -segment_time {{.C.Transcode.Hls.Segment_length}} -segment_format mpegts -segment_list_flags +live -segment_list_type m3u8 -individual_header_trailer 1 -segment_list index.m3u8 hls/%09d.ts"
	TC_ARG_DASH     string = " "
	TC_ARG_MP4      string = " -codec copy video.mp4 "
	TC_ARG_OGG      string = " -codec:v libtheora -b:v 1200k -codec:a vorbis -b:a 128k video.ogv "
	TC_ARG_WEBM     string = " -f webm -codec:v libvpx -quality realtime -cpu-used 0 -b:v 1200k -qmin 10 -qmax 42 -minrate 1200k -maxrate 1200k -bufsize 1500k -threads 1 -codec:a libvorbis -b:a 128k video.webm "
	TC_ARG_WEBM_HLS string = " -f webm -codec:v libvpx -quality realtime -cpu-used 0 -keyint_min {{.Webm_gop_size}} -g {{.Webm_gop_size}} -b:v 1200k -qmin 10 -qmax 42 -maxrate 1200k -bufsize 1500k -lag-in-frames 0 -rc_lookahead 0 -flags +global_header -codec:a libvorbis -b:a 128k -flags +global_header -map 0 -f segment -segment_list_flags +live -segment_time {{.C.Transcode.Hls.Segment_length}} -segment_format webm -flags +global_header -segment_list webm_index.m3u8 webm/%09d.webm "
	TC_ARG_THUMB    string = " -f image2 {{.Thumb_size}} -vsync 1 -vf fps=fps=1/{{.Thumb_rate}} thumb/%09d.jpg "
	TC_ARG_POSTER   string = " -f image2 {{.Poster_size}} -vsync 1 -vf fps=fps=1/{{.Poster_rate}} {{.Poster_skip}} -vframes {{.Poster_corrected_count}} poster/%09d.jpg "
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
// -fflags +nobuffer                        -- breaks SPS/PPS detection on some TS streams from Android ffmpeg

// Unused Options
//
// -copyinkf:0
// -fflags +discardcorrupt
// -fpsprobesize 2

const (
	TC_ARG_TSIN  string = " -fflags +genpts+igndts -err_detect compliant -avoid_negative_ts 1 -correct_ts_overflow 1 -max_delay 500000 -analyzeduration 500000 -f mpegts -c:0 h264 -vsync 0 -copyts -copytb 1 "
	TC_ARG_AVCIN string = " -fflags +genpts+igndts -max_delay 0 -analyzeduration 0 -f h264 -c:0 h264 -copytb 0 "
)

const (
	TC_CMD_START_PROD string = "-y -v quiet "
	TC_CMD_START_DEV  string = "-y -v debug "
	TC_CMD_INPUT      string = " -i pipe:0 "
	TC_CMD_END_PROD   string = ""
	TC_CMD_END_DEV    string = ""
)

func (s *Session) IsMimeTypeValid(m string) bool {
	// first letter and letter after hyphen uppercase, rest lowercase
	//contentType := http.CanonicalHeaderKey(m)
	switch s.Type {
	case TC_INGEST_TS:
		return true // video/x-mpegts

	case TC_INGEST_AVC:
		return true //

	default:
		return false
	}
	return true
}

func (s *Session) createOutputDirectories() error {

	// HLS: <hls-data-path>/<id>/hls
	if s.c.Transcode.Hls.Enabled {
		p := s.c.Transcode.Output_path + "/" + s.idstr + "/hls"
		if err := os.MkdirAll(p, utils.PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return err
		}
	}

	// DASH: <dash-data-path>/<id>/dash
	if s.c.Transcode.Dash.Enabled {
		p := s.c.Transcode.Output_path + "/" + s.idstr + "/dash"
		if err := os.MkdirAll(p, utils.PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return err
		}
	}

	// MP4: <*-data-path>/<id>
	if s.c.Transcode.Mp4.Enabled {
		p := s.c.Transcode.Output_path + "/" + s.idstr
		if err := os.MkdirAll(p, utils.PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return err
		}
	}

	// OGG: <*-data-path>/<id>
	if s.c.Transcode.Ogg.Enabled {
		p := s.c.Transcode.Output_path + "/" + s.idstr
		if err := os.MkdirAll(p, utils.PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return err
		}
	}

	// WEBM: <*-data-path>/<id>
	if s.c.Transcode.Webm.Enabled {
		p := s.c.Transcode.Output_path + "/" + s.idstr
		if err := os.MkdirAll(p, utils.PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return err
		}
	}

	// WEBM-HLS: <*-data-path>/<id>/webm
	if s.c.Transcode.Webm_hls.Enabled {
		p := s.c.Transcode.Output_path + "/" + s.idstr + "/webm"
		if err := os.MkdirAll(p, utils.PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return err
		}
	}

	// Thumb: <*-data-path>/<id>/thumb
	if s.c.Transcode.Thumb.Enabled {
		p := s.c.Transcode.Output_path + "/" + s.idstr + "/thumb"
		if err := os.MkdirAll(p, utils.PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return err
		}
	}

	// Poster: <*-data-path>/<id>/poster
	if s.c.Transcode.Poster.Enabled {
		p := s.c.Transcode.Output_path + "/" + s.idstr + "/poster"
		if err := os.MkdirAll(p, utils.PERM_DIR); err != nil {
			log.Printf("Error: cannot create directory %s: %s", p, err)
			return err
		}
	}

	return nil
}

// assemble transcode command (for now we run a single transcoder 'ffmpeg')
func (s *Session) BuildTranscodeCommand() string {
	//return "ffmpeg --help"
	var cmd string

	// transcoder command (generate debug output when in dev mode)
	if s.c.Server.Production_mode {
		cmd = TC_CMD_START_PROD
	} else {
		cmd = TC_CMD_START_DEV
	}

	// input spec
	switch s.Type {
	case TC_INGEST_AVC:
		cmd += TC_ARG_AVCIN + TC_CMD_INPUT
	case TC_INGEST_TS:
		cmd += TC_ARG_TSIN + TC_CMD_INPUT
	}

	// segment file formats
	if s.c.Transcode.Hls.Enabled {
		cmd += TC_ARG_HLS
	}
	if s.c.Transcode.Dash.Enabled {
		cmd += TC_ARG_DASH
	}
	if s.c.Transcode.Webm_hls.Enabled {
		cmd += TC_ARG_WEBM_HLS
	}

	// full file formats
	if s.c.Transcode.Mp4.Enabled {
		cmd += TC_ARG_MP4
	}
	if s.c.Transcode.Ogg.Enabled {
		cmd += TC_ARG_OGG
	}
	if s.c.Transcode.Webm.Enabled {
		cmd += TC_ARG_WEBM
	}

	// image formats
	if s.c.Transcode.Thumb.Enabled {
		cmd += TC_ARG_THUMB
	}
	if s.c.Transcode.Poster.Enabled {
		cmd += TC_ARG_POSTER
	}

	// end trancode command line
	if s.c.Server.Production_mode {
		cmd += TC_CMD_END_PROD
	} else {
		cmd += TC_CMD_END_DEV
	}

	// combine session and server config for access by template matcher
	var cmd_writer bytes.Buffer
	var d = TemplateData{s.c, s, 0, "", 1, "", 1, 1, ""}

	// Poster
	if s.c.Transcode.Poster.Enabled {
		// scaling parameters
		if s.c.Transcode.Poster.Size == "auto" || s.c.Transcode.Poster.Size == "" {
			d.Poster_size = ""
		} else {
			d.Poster_size = "-s " + s.c.Transcode.Poster.Size
		}

		// skip parameter
		if s.c.Transcode.Poster.Skip > 0 {
			hou := s.c.Transcode.Poster.Skip / 3600
			min := (s.c.Transcode.Poster.Skip - hou*3600) / 60
			sec := (s.c.Transcode.Poster.Skip - hou*3600 - min*60)
			d.Poster_skip = fmt.Sprintf("-ss %02d:%02d:%02d.0", hou, min, sec)
		}

		// step interval
		if s.c.Transcode.Poster.Step > 0 {
			d.Poster_rate = s.c.Transcode.Poster.Step
		}

		// increase the count of poster frames to output since ffmpeg flushes
		// the image2 pipeline at the end only, hence a single image would not be
		// written before the stream ends
		//
		// on the downside, this solution writes one image more than expected by the user
		//
		d.Poster_corrected_count = s.c.Transcode.Poster.Count + 1
	}

	// Thumbnail
	if s.c.Transcode.Thumb.Enabled {
		// scaling parameter
		if s.c.Transcode.Thumb.Size != "" {
			d.Thumb_size = "-s " + s.c.Transcode.Thumb.Size
		}

		// step interval
		if s.c.Transcode.Thumb.Step > 0 {
			d.Thumb_rate = s.c.Transcode.Thumb.Step
		}
	}

	// GOP size (25 fps is a guess)
	if s.c.Transcode.Webm_hls.Enabled {
		d.Webm_gop_size = s.c.Transcode.Webm_hls.Gop_size
	}

	// replace placeholders with config strings
	t, err := template.New("cmd").Parse(cmd)
	if err != nil {
		log.Fatalf("Error parsing cmd template: %s\n", err)
	}

	err = t.Execute(&cmd_writer, d)
	if err != nil {
		log.Fatalf("Error generating cmd string: %s\n", err)
	}

	return cmd_writer.String()
}
