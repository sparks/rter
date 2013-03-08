// rtER Project - SRL, McGill University, 2013
//
// Author: echa@cim.mcgill.ca

package main

const (
	TRANSCODE_TYPE_UNKNOWN 	int = 0
	TRANSCODE_TYPE_AVC 		int = 1
	TRANSCODE_TYPE_TS  		int = 2
)

const (
	TRANSCODE_PARAM_HLS string = ""
	TRANSCODE_PARAM_DASH string = ""
	TRANSCODE_PARAM_MP4 string = ""
	TRANSCODE_PARAM_OGG string = ""
	TRANSCODE_PARAM_THUMB string = ""
	TRANSCODE_PARAM_POSTER string = ""
)


func IsMimeTypeValid(t int, m string) bool {
	// first letter and letter after hyphen uppercase, rest lowercase
	//contentType := http.CanonicalHeaderKey(m)

	return true


}



// TS input
// ffmpeg -v quiet -fflags nobuffer -i pipe:0 -vsync 2 -copyts -copytb 1
//  -codec copy -map 0 -f segment -segment_time 2 -segment_format mpegts
//  -segment_list_flags +live -segment_list test.m3u8  teststream-%09d.ts

// AVC input


func GetTranscodeParams(t int) string {
	return "ffmpeg --help"
}
