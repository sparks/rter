// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

import (
	"encoding/json"
	"io/ioutil"
	"flag"
	"log"
)

// Command Line Options
var configfile = flag.String("config", "config.json", "server config file")

// Server Configuration
type ServerConfig struct {
	// server
	Server struct {
		Addr string `json:"addr"`
		Port uint64 `json:"port"`
		Secure_mode bool `json:"secure_mode"`
		Cert_file string `json:"cert_file"`
		Key_file string `json:"key_file"`
	}
	// limits
	Limits struct {
		Max_memory_mbytes uint64 `json:"max_memory_mbytes"`
		Max_ingest_sessions uint64 `json:"max_ingest_sessions"`
		Max_ingest_bandwidth_kbit uint64 `json:"max_ingest_bandwidth_kbit"`
		Rate_limit_enable bool `json:"rate_limit_enable"`
		Rate_limit_ingest_window uint64 `json:"rate_limit_ingest_window"`
		Rate_limit_ingest_sessions_per_source uint64 `json:"rate_limit_ingest_sessions_per_source"`
		Rate_limit_ingest_bytes_per_source uint64 `json:"rate_limit_ingest_bytes_per_source"`
	}
	Auth struct {
	// auth
		Signkey string `json:"signkey"`
	}
	// ingest
	Ingest struct {
		Enable_avc_ingest bool `json:"avc"`
		Enable_ts_ingest bool `json:"ts"`
		Enable_chunk_ingest bool `json:"chunk"`
	}
	// paths
	Paths struct {
		Data_storage_path string `json:"storage"`
	}
	// transcode
	Transcode struct {
		Enable_hls_transcode bool `json:"hls"`
		Enable_mp4_transcode bool `json:"mp4"`
		Enable_ogg_transcode bool `json:"ogg"`
		Enable_dash_transcode bool `json:"dash"`
		Enable_thumb_transcode bool `json:"thumb"`
		Enable_poster_transcode bool `json:"poster"`
	}
}

func ParseConfig(c *ServerConfig) {

	// parse command line
	flag.Parse()

	// set default values
	c.Server.Addr = "127.0.0.1"
	c.Server.Port = 8080
	c.Server.Secure_mode = false
	c.Server.Cert_file = ""
	c.Server.Key_file = ""
	c.Limits.Max_memory_mbytes = 128
	c.Limits.Max_ingest_sessions = 10
	c.Limits.Max_ingest_bandwidth_kbit = 10000
	c.Limits.Rate_limit_enable = false
	c.Limits.Rate_limit_ingest_window = 15
	c.Limits.Rate_limit_ingest_sessions_per_source = 100
	c.Limits.Rate_limit_ingest_bytes_per_source = 134217728
	c.Auth.Signkey = "none"
	c.Ingest.Enable_avc_ingest = true
	c.Ingest.Enable_ts_ingest = true
	c.Ingest.Enable_chunk_ingest = false
	c.Paths.Data_storage_path = "./data"
	c.Transcode.Enable_hls_transcode = true
	c.Transcode.Enable_mp4_transcode = false
	c.Transcode.Enable_ogg_transcode = false
	c.Transcode.Enable_dash_transcode = false
	c.Transcode.Enable_thumb_transcode = true
	c.Transcode.Enable_poster_transcode = true

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

    log.Printf("ServerConfig: %+v\n", c)
}
