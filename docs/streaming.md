# Live Streaming to HTML5 Browsers

Here is a collection of resources about live video streaming from a networked anywhere to HTML5 web browsers. The technologies considered are

* HLS (HTML Live Streaming, Apple's implementation in Quicktime and iOS)
* DASH (Dynamic Adaptive Streaming over HTTP, an MPEG Standard)
* WebRTC (Real-Time point-to-point video, a W3C standard)

Terminology used throughout this document

* __Streaming Source__, a device/app capturing a single live video feed
* __Streaming Server__, a central entity distributing multiple live video feeds
* __Streaming Client__, a HTML5 browser displaying one or multiple video feeds


## Classical Video-on-Demand

Some resouces and code from traditional file-based video playback can be repurposed for live video in web browsers. Here's a non-exhaustive list of recources on that issue:

* Streaming Servers
  * [Erlyvideo][8] (limited formats, HLS, RTSP/RTP, RTMP, MPEG-TS over HTTP/UDP, RTP multicast)
  * [FFServer][10] (multi-format RTSP/RTP, HTTP, from file/pipe/local capture)
  * [FFMpeg][9] supports streaming as well, see [Streaming Guide][24]
  * [LiveMedia][11] (multi-format RTSP/RTP, RTP over HTTP, RTP/UDP multicast)
  * [VLC][20], see [detailed features][23]
  * [Wowza][26], (commercial)
* Video Playback in web browsers
  * see [format support][16] in browsers
  * support according to [HTML5 standard][17]
  * on iOS devices, see [Apple documentation][18] and [ZEncoder recommendations][19]
  * on Android, see ([supported formats][13])
* Javascript Players (all using Flash as fallback if native format support is missing)
  * [Video.js](http://videojs.com/)
  * [MediaElement.js](http://mediaelementjs.com/)
  * [Projekktor](http://www.projekktor.com)

## Classical Live Video Streaming

The most robust and widespread technologies used to stream live video between any type of equipment are not accessible from within current web browsers. However, some technologies are available as 3rd party library on mobile platforms. I briefly list them here for

* RTSP/RTP/RTCP, a set of IETF protocols for live video communication
  * .. over UDP and multicast UDP
  * .. over HTTP
  * .. over TCP
* MPEG2 Transport Streams, a MPEG standard for packaging compressed video for transport
  * used in professional production and broadcasting (DVB)
  * also transport framing for HLS and DASH

Resources
* [FFMpeg Streaming Guide](http://ffmpeg.org/trac/ffmpeg/wiki/StreamingGuide)
* [FFServer Example](http://ffmpeg.org/trac/ffmpeg/wiki/Streaming%20media%20with%20ffserver)
* [Live Streaming with FFMpeg](http://sonnati.wordpress.com/2012/07/02/ffmpeg-the-swiss-army-knife-of-internet-streaming-part-v/)
* [Another live streaming example](http://www.onvos.com/http-live-streaming-howto.html)
* [VLC Advanced Streaming Howto][22]


## HLS: HTTP Live Streaming

[HLS][1] is a proprietary 'standard' by Apple to feed live video content to iOS devices and Quicktime media player. Like DASH, HLS separates video files into segments for individual download. HLS requires MPEG2-TS and H264/AAC. The description format is an extended version of the MP3 playlist format M3U called M3U8.


#### Software

* Encoding at Source
  * Android H264 Baseline Profile in MPEG2-TS (see [supported formats][13])
  * iOS should support H264 Baseline, for code examples see hints on [Stackoverflow][14]
  * YouTube Live encoding [recommendations][15]
* Segmentation and M3U8 generation:
  * [FFmpeg][9], see [segmenter doc][12] for details
  * [VLC][20], see [Wiki][21] for details
  * [Erlyvideo][8] (commercial version only!)
  * [Live HLS Streamer][28], Open Source and based on FFMPeg (also on [GitHub][29])
* HLS Playback: on HTML5/JS web browsers
  * native support on Safari and iOS only (see also [HLS support][7], and [HLS Streaming from the iOS perspective][27])


#### FFMPEG HLS Examples

```
# video only transcoding into H264 for iPod/iPhone
./ffmpeg -v 9 -loglevel info -re -i sourcefile.avi -an \
-c:v libx264 -b:v 128k -vpre ipod320 \
-flags -global_header -map 0 -f segment -segment_time 4 \
-segment_list test.m3u8 -segment_format mpegts stream%05d.ts
```


```
# not sure if its working
ffmpeg -i encoded.mp4 -c copy -map 0 -vbsf h264_mp4toannexb -f segment -segment_time 10 -segment_format mpegts stream%d.ts
```


## DASH, Dynamic Adaptive Streaming over HTTP

[DASH](25) is a recent standard for segmented video streaming supporting adaptation through bitstream switching. DASH specifies a description format (MPD files), requires a certain packetization  format (MPEG2-TS or ISO-MBFF), a subset of codecs (H264, AAC, VP8?), and HTTP as transport, however, DASH does __not__ specify transport or adaptation strategies, buffering or playback interactions (which all are left to implementations). A normal Web Server can be used to deliver segmented files to browsers which makes elegant use of already deployed infrastructure including CDNs and caches. DASH further defines different profiles, e.g. for live and on-demand streaming which basically refine applicable limits on configuration settings.

To use DASH one nees a server-side implementation for generating compliant segmented video files plus the MPD description. At a Web client an implementation is required to fetch and parse the MPD and download segments as required.

#### Software

* Segmentation and MPD generation: [GPAC][2] in a recent version from [SVN][3], for more info on DASH support see [here][4] and [here][5]
* Delivery: any HTTP server
* Playback: Google Chrome with Media Source Extention plus a 3rd party DASH Javascript Engine, e.g. [DASH-JS][6]
* Firefox with built-in DASH client support (version?)

#### Examples

* [DASH Player Test in Chrome](http://dash-mse-test.appspot.com/dash-player.html)
* [Decoder Test Sequences](http://dash-mse-test.appspot.com/decoder-test.html)
* [DASH Test Sequences](http://gpac.wp.mines-telecom.fr/2012/02/23/dash-sequences/)


### DASH Live Streaming using GPAC

[GPAC DASH Options](http://gpac.wp.mines-telecom.fr/mp4box/dash/)

Call `MP4Box` on regular basis with the segments you need to add to the live session:

```
MP4Box [DASH_OPTIONS] -dash-ctx dasher.cfg fileX.mp4
```

#### GPAC DASH Example from File
```
MP4Box -dash 2000 -out index.mpd -dash-profile live -rap -frag-rap \
  -segment-name dash-Num$Number$-BW$Bandwidth$-Time$Time$ \
  -base-url http://localhost:8080/assets/dash -mpd-title "Testvideo" \
  -mpd-source "Some Source" -mpd-info-url "InfoURL" -dash-ctx dash-context.bin \
  ../video/file.m4v
```

Todo: try options `-dash-live=aow-context.bin 2000` and `-noprog`

#### GPAC Live DASH Example from 'Capture to File'
```
MP4Box -dash 20000 -out index.mpd -dash-profile live -rap -frag-rap \
  -segment-name dash-Num$Number$-BW$Bandwidth$-Time$Time$ \
  -base-url http://localhost:8080/assets/dash -mpd-title "Title" \
  -mpd-source "Some Source" -mpd-info-url "InfoURL" -dash-live=context.bin 2000 \
  video.h264
```

__Problem:__ MP4Box sleep interval doubled each time, probably because the file kept expanding


[1]: http://developer.apple.com/library/ios/#documentation/NetworkingInternet/Conceptual/StreamingMediaGuide/UsingHTTPLiveStreaming/UsingHTTPLiveStreaming.html
[2]: http://gpac.wp.mines-telecom.fr/mp4box/dash/
[3]: http://sourceforge.net/projects/gpac/develop
[4]: http://gpac.wp.mines-telecom.fr/2011/02/02/mp4box-fragmentation-segmentation-splitting-and-interleaving/
[5]: http://gpac.wp.mines-telecom.fr/2012/02/01/dash-support/
[6]: http://www-itec.uni-klu.ac.at/dash/?p=792
[7]: http://www.longtailvideo.com/html5/hls
[8]: http://erlyvideo.org
[9]: http://www.ffmpeg.org/
[10]: http://ffmpeg.org/ffserver.html
[11]: http://www.live555.com/liveMedia/
[12]: http://www.ffmpeg.org/ffmpeg-formats.html#segment_002c-stream_005fsegment_002c-ssegment
[13]: http://developer.android.com/guide/appendix/media-formats.html
[14]: http://stackoverflow.com/questions/10313308/how-can-i-perform-hardware-accelerated-h-264-encoding-and-decoding-for-streaming
[15]: http://support.google.com/youtube/bin/answer.py?hl=en&answer=1723080
[16]: http://wiki.whatwg.org/wiki/Video_type_parameters#Browser_Support
[17]: http://www.w3.org/TR/2011/WD-html5-author-20110809/the-source-element.html
[18]: http://developer.apple.com/library/ios/#technotes/tn2224/_index.html
[19]: http://blog.zencoder.com/2012/01/24/encoding-settings-for-perfect-ipadiphone-video/
[20]: http://www.videolan.org/vlc/index.html
[21]: http://wiki.videolan.org/Main_Page
[22]: http://wiki.videolan.org/Documentation:Streaming_HowTo/Advanced_Streaming_Using_the_Command_Line
[23]: http://www.videolan.org/streaming-features.html
[24]: http://ffmpeg.org/trac/ffmpeg/wiki/StreamingGuide
[25]: http://dashif.org/mpeg-dash/
[26]: http://www.wowza.com/
[27]: http://blog.refractalize.org/post/8171724662/implementing-http-live-streaming
[28]: http://www.ioncannon.net/projects/http-live-video-stream-segmenter-and-distributor/
[29]: https://github.com/carsonmcdonald/HTTP-Live-Video-Stream-Segmenter-and-Distributor
