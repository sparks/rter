// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

// Command Line Options
var configfile = flag.String("config", "config.json", "server config file")

// Server Configuration
type ServerConfig struct {
	// server
	Server struct {
		Addr            string `json:"addr"`
		Port            uint64 `json:"port"`
		Production_mode bool   `json:"production_mode"`
		Secure_mode     bool   `json:"secure_mode"`
		Cert_file       string `json:"cert_file"`
		Key_file        string `json:"key_file"`
		Session_timeout uint64 `json:"session_timeout"`
		Session_maxage  uint64 `json:"session_maxage"`
	}
	// limits
	Limits struct {
		Max_cpu                                int    `json:"max_cpu"`
		Max_memory_mbytes                      uint64 `json:"max_memory_mbytes"`
		Max_ingest_sessions                    uint64 `json:"max_ingest_sessions"`
		Max_ingest_bandwidth_kbit              uint64 `json:"max_ingest_bandwidth_kbit"`
		Rate_limit_enable                      bool   `json:"rate_limit_enable"`
		Rate_limit_ingest_window               uint64 `json:"rate_limit_ingest_window"`
		Rate_limit_ingest_sessions_per_source  uint64 `json:"rate_limit_ingest_sessions_per_source"`
		Rate_limit_ingest_requests_per_session uint64 `json:"rate_limit_ingest_requests_per_session"`
		Rate_limit_ingest_bytes_per_source     uint64 `json:"rate_limit_ingest_bytes_per_source"`
	}
	Auth struct {
		// auth
		Enabled      bool   `json:"enabled"`
		Token_secret string `json:"token_secret"`
	}
	// ingest
	Ingest struct {
		Enable_avc_ingest   bool `json:"avc"`
		Enable_ts_ingest    bool `json:"ts"`
		Enable_chunk_ingest bool `json:"chunk"`
	}
	// transcode
	Transcode struct {
		Command     string `json:"command"`
		Log_path    string `json:"log_path"`
		Output_path string `json:"output_path"`
		Hls         struct {
			Enabled        bool   `json:"enabled"`
			Segment_length uint64 `json:"segment_length"`
		}
		Dash struct {
			Enabled        bool   `json:"enabled"`
			Segment_length uint64 `json:"segment_length"`
		}
		Mp4 struct {
			Enabled bool `json:"enabled"`
		}
		Ogg struct {
			Enabled bool `json:"enabled"`
		}
		Webm struct {
			Enabled bool `json:"enabled"`
		}
		Thumb struct {
			Enabled bool   `json:"enabled"`
			Size    string `json:"size"`
			Step    uint64 `json:"step"`
		}
		Poster struct {
			Enabled bool   `json:"enabled"`
			Size    string `json:"size"`
			Count   uint64 `json:"count"`
			Skip    uint64 `json:"skip"`
			Step    uint64 `json:"step"`
		}
	}
}

func (c *ServerConfig) ParseConfig() {

	// parse command line
	flag.Parse()

	// set default values
	c.Server.Addr = "127.0.0.1"                             // bind server to IP address
	c.Server.Port = 8080                                    // bind server to http port
	c.Server.Production_mode = false                        // run in production or develop mode
	c.Server.Secure_mode = false                            // use SSL mode
	c.Server.Cert_file = ""                                 // SSL CA certificate
	c.Server.Key_file = ""                                  // SSL private key
	c.Server.Session_timeout = 10                           // close after 10 seconds inactivity
	c.Server.Session_maxage = 3600                          // keep state for at most 1 hour
	c.Limits.Max_cpu = 1                                    // max number of CPUs used
	c.Limits.Max_memory_mbytes = 128                        // max amount of memory used
	c.Limits.Max_ingest_sessions = 10                       // max active sessions
	c.Limits.Max_ingest_bandwidth_kbit = 10000              // max bandwidth for ingest data
	c.Limits.Rate_limit_enable = false                      // enable per-client rate limit
	c.Limits.Rate_limit_ingest_window = 15                  // time window for resetting rate limits in minutes
	c.Limits.Rate_limit_ingest_sessions_per_source = 100    // new session limit per client
	c.Limits.Rate_limit_ingest_requests_per_session = 30000 // max requests (frames) per session
	c.Limits.Rate_limit_ingest_bytes_per_source = 134217728 // data volume limit per client
	c.Auth.Enabled = false                                  // disable token auth
	c.Auth.Token_secret = ""                                // private request authentication key
	c.Ingest.Enable_avc_ingest = true                       // enable H264AVC ingest endpoint
	c.Ingest.Enable_ts_ingest = true                        // enable MPEG2-TS ingest endpoint
	c.Ingest.Enable_chunk_ingest = false                    // enable chunked file transfer endpoint
	c.Transcode.Command = "ffmpeg"                          // transcoder command
	c.Transcode.Log_path = "./data/log"                     // transcoder logfile path
	c.Transcode.Output_path = "./data"                      // transcoder output path root
	c.Transcode.Hls.Enabled = false
	c.Transcode.Hls.Segment_length = 2
	c.Transcode.Dash.Enabled = false
	c.Transcode.Dash.Segment_length = 2
	c.Transcode.Mp4.Enabled = true
	c.Transcode.Ogg.Enabled = false
	c.Transcode.Webm.Enabled = false

	c.Transcode.Thumb.Enabled = false // save live thumbnails
	c.Transcode.Thumb.Size = "160x90" // thumbnail scaling dimensions
	c.Transcode.Thumb.Step = 2        // interval between thumbnails in sec

	c.Transcode.Poster.Enabled = false // save poster image
	c.Transcode.Poster.Size = "auto"   // auto = same size as source video
	c.Transcode.Poster.Count = 1       // number of poster images to store
	c.Transcode.Poster.Skip = 10       // skip number of seconds at start
	c.Transcode.Poster.Step = 0        // interval between poster images in sec

	// read config
	jsonconfig, err := ioutil.ReadFile(*configfile)
	if err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}

	// unpack config from JSON into Go struct (sets only the defined values)
	err = json.Unmarshal(jsonconfig, &c)
	if err != nil {
		log.Fatalf("Error parsing config file: %s\n", err)
	}
}

func (c *ServerConfig) Print() {
	if b, err := json.MarshalIndent(c, "", "  "); err != nil {
		log.Printf("ServerConfig error: %s", err)
	} else {
		log.Printf("ServerConfig:")
		os.Stdout.Write(b)
		os.Stdout.WriteString("\n")
	}
}

func (c *ServerConfig) SanityCheck() {

	// server IP and port
	// - warn when localhost is set in production
	if c.Server.Production_mode && c.Server.Addr == "127.0.0.1" {
		log.Printf("Warning: server is reachable via localhost only!")
	}

	// - warn when 0.0.0.0 is set in production and dev
	if c.Server.Addr == "0.0.0.0" {
		log.Printf("Warning: server is reachable on all interfaces (0.0.0.0)!")
	}

	// - warn when port is blocked by Safari/iOS (6666, 6667, 6000)
	if c.Server.Port == 6666 || c.Server.Port == 6667 || c.Server.Port == 6000 {
		log.Printf("Warning: server port may be blocked by Apple Webkit browsers!")
	}

	// secure mode warning
	if c.Server.Production_mode && !c.Server.Secure_mode {
		log.Printf("Warning: HTTPS is strongly recommended for production mode!")
	}

	// secure mode requires cert and key files
	if c.Server.Secure_mode &&
		(c.Server.Cert_file == "" || c.Server.Key_file == "") {
		log.Fatal("Error: secure mode requires cert_file and key_file!")
	}

	// token auth waring
	if c.Server.Production_mode && !c.Auth.Enabled {
		log.Printf("Warning: Token authorization is strongly recommended for production mode!")
	}

	// token auth requires server secret
	if c.Auth.Enabled && c.Auth.Token_secret == "" {
		log.Fatal("Error: token auth requires secret key!")
	}

	// sane resource limits
	if c.Limits.Max_cpu <= 0 {
		log.Printf("Warning: max_cpu must be >= 0! Set to 1.")
		c.Limits.Max_cpu = 1
	}

	available_cpu := runtime.NumCPU()
	if c.Limits.Max_cpu > available_cpu {
		log.Printf("Warning: max_cpu reduced to %d", available_cpu)
		c.Limits.Max_cpu = available_cpu
	}

	// meaningful rate limits for given reset interval
	if c.Limits.Rate_limit_enable {
		// - requests per session assuming frame-wise upload at 30fps
		if c.Limits.Rate_limit_ingest_requests_per_session < c.Limits.Rate_limit_ingest_window*60*30 {
			log.Printf("Warning: request rate limit may not be sufficient for frame-wise POST")
		}

		// - bytes per source sufficient for at least one 600kbit/sec stream (=75000 bytes/sec)
		if c.Limits.Rate_limit_ingest_bytes_per_source < c.Limits.Rate_limit_ingest_window*60*75000 {
			log.Printf("Warning: volume rate limit may not be sufficient for 600kbit streams")
		}
	}

	// test output directory permissions (?wx??????) 0300
	info, err := os.Stat(c.Transcode.Output_path)
	if err != nil || !info.IsDir() || info.Mode().Perm()&os.FileMode(0300) != 0300 {
		log.Fatal("Error: output path is no writable directory!")
	}

	// test transcoder executable exist (r?xr?xr?x) 555 0x16D
	info, err = os.Stat(c.Transcode.Command)
	if err != nil || info.Mode().Perm()&os.FileMode(0500) != 0500 {
		log.Fatal("Error: transcoder command not executable!")
	}
}

func (c *ServerConfig) CheckTranscoderCapabilities() {

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
