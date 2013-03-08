package main

const (
	TRANSCODE_TYPE_UNKNOWN 	int = 0
	TRANSCODE_TYPE_AVC 		int = 1
	TRANSCODE_TYPE_TS  		int = 2
)


func GetTranscodeParams(t int) string {
	return "ffmpeg --help"
}
