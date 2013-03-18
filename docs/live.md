# Live Transcoding

This is a collection of working commandline examples to show how one could use FFMpeg and VLC for live transcoding of video streams. All examples have been tested on OSX 10.7.5 with FFMPeg 1.1.3 and VLC 2.0.5 in early 2013.

Documentation links

- [FFMpeg](https://ffmpeg.org/ffmpeg.html), [Muxers](https://ffmpeg.org/ffmpeg-formats.html), [Encoders](https://ffmpeg.org/ffmpeg-codecs.html), [Protocols](https://ffmpeg.org/ffmpeg-protocols.html)
- FFMPEG [multiple outputs](http://ffmpeg.org/trac/ffmpeg/wiki/Creating%20multiple%20outputs)
- [VLC](http://www.videolan.org/) [Docu](http://wiki.videolan.org/Documentation:Documentation), [Advanced Use](http://wiki.videolan.org/Documentation:Play_HowTo/Advanced_Use_of_VLC)


## Creating a TCP (HTTP) Video Streaming Connection

#### Running a fake source
```
ffmpeg -v debug -y -re -i file.mp4 -vsync 1 -codec copy -bsf h264_mp4toannexb -f mpegts http://localhost:6666@listen
```

#### Running a fake server
```
ffmpeg -v debug -fflags nobuffer -i tcp://127.0.0.1:6666 -r 15 -vsync 2 -copyts -copytb 1 -codec copy -map 0 -f segment -segment_time 2 -segment_list stream.m3u8 -segment_format mpegts -segment_list_flags +live stream-%09d.ts
```


## Testing the rtER Video Streaming Server

Data path: Fake Source -> Server -> Web Clients

Start the video server with the command `./videoserver` and once a stream is recording visit the player page at `server:/v1/videos/:id/play` where the id is the video streams unique identifier.


#### Stream an MPEG2 Transport Stream from a file source

```
ffmpeg -v debug -y -re -i ../../../trimmmr/ingest/data/video/aow-docu-2011.m4v -vsync 1 -map 0 -codec copy -bsf h264_mp4toannexb -f mpegts -copytb 0 http://localhost:6666/v1/ingest/2/avc
```

#### Live transcode a file and send it as MPEG2 Transport Stream

```
ffmpeg -v debug -y -re -i ../../../trimmmr/ingest/data/video/aow-docu-2011.m4v -f mpegts -c:v libx264 -preset ultrafast -tune zerolatency -crf 20 -x264opts keyint=50:bframes=0:ratetol=1.0:ref=1 -profile baseline -maxrate 1200k -bufsize 1200k  -c:a copy http://localhost:6666/v1/ingest/17/ts
```

#### Capture and encode live video with VLC, send as MPEG2 Transport Stream, and bridge through FFMPEG

```
/Applications/VLC.app/Contents/MacOS/VLC qtcapture:// -vvvv --no-drop-late-frames --no-skip-frames --sout='#transcode{vcodec=h264,fps=15,venc=x264{preset=ultrafast,tune=zerolatency,keyint=30,bframes=0,ref=1,level=30,profile=baseline,hrd=cbr,crf=20,ratetol=1.0,vbv-maxrate=1200,vbv-bufsize=1200,lookahead=0}}:standard{access=http{mime="video/MP2T"},mux=ts,dst=127.0.0.1:5555}' --qtcapture-width=640 --qtcapture-height=480 --live-caching=200 --intf=macosx

ffmpeg -v debug -y -i http://127.0.0.1:5555 -vsync 1 -map 0 -codec copy -r 15 -f mpegts -copytb 0 http://localhost:6666/v1/ingest/18/ts
```

#### Live transcode a file and send it as raw H264/AVC stream

```
ffmpeg -v debug -y -re -i ../../../trimmmr/ingest/data/video/aow-docu-2011.m4v -f h264 -c:v libx264 -preset ultrafast -tune zerolatency -crf 20 -x264opts keyint=50:bframes=0:ratetol=1.0:ref=1:repeat-headers=1 -profile baseline -maxrate 1200k -bufsize 1200k  -an http://localhost:6666/v1/ingest/30/avc
```


## VLC Live Stream Examples

#### General Parameters
```
-vvvv
--no-drop-late-frames
--no-skip-frames
--sout='#transcode{<params>}:standard{<params>}'
--mtu 8192
```

#### QT capture on OSX (must use `--intf=macosx` to access QT capture)
```
--qtcapture-width 640
--qtcapture-height 480
--live-caching 200 [ms] !important
--intf=macosx
```

#### X264 Parameters
```
--sout-x264-keyint=30
--sout-x264-bframes=0
--sout-x264-ref=1
--sout-x264-level=30
--sout-x264-profile=baseline
--sout-x264-hrd=cbr
--sout-x264-crf=20
--sout-x264-ratetol=1.0
--sout-x264-vbv-maxrate=1200  [kbit/s]
--sout-x264-vbv-bufsize=1200  [kbit/s]
--sout-x264-preset=ultrafast
--sout-x264-tune=zerolatency
--sout-x264-aud
--sout-x264-lookahead=0
```

#### Video Encoding Parameters
```
--sout-transcode-venc x264
--sout-transcode-vcodec h264
--sout-transcode-vb 2000
--sout-transcode-fps 25
--sout-transcode-width W
--sout-transcode-height H
```

#### Audio Encoding Parameters
```
--sout-transcode-aenc
--sout-transcode-acodec mp4a
--sout-transcode-ab 64 [kbit/s]
--sout-transcode-channels 2
--sout-transcode-samplerate 44100
```

#### Transport Stream Muxing
```
--sout-standard-mux=ts
--sout-ts-shaping 2000 [ms]      minimum interval with constant bitrate in VBR stream
--sout-ts-use-key-frames      limit shaping interval at key frames
--sout-ts-dts-delay=0
--sout-mux-caching=0
--clock-synchro=1
```

#### UDP output
```
--sout-standard-access=udp
--sout-standard-dst IP:PORT
--sout-mux-caching=20 [ms]
--sout-udp-caching=0 [ms]
--sout-udp-group=10            send groups of 10 packets at a time
--sout-udp-late=100 [ms]    ?  drop packets arriving later than N ms
--sout-udp-raw              ?  dont wait until MTU is filled before sending
```

#### HTTP output
```
--sout-http-mime="video/MP2T"
```

### VLC - HTTP Live Stream (this is not HLS or DASH!)

VLC implements a streaming server and FFMpeg pulls raw data via HTTP. Note that VLC, unlike FFMpeg, is not able to send a live stream towards a HTTP server.

#### 1  Run capture outputting stream via built-in HTTP server
```
/Applications/VLC.app/Contents/MacOS/VLC qtcapture:// -vvvv --no-drop-late-frames --no-skip-frames --sout='#transcode{vcodec=h264,fps=15,venc=x264{preset=ultrafast,tune=zerolatency,keyint=30,bframes=0,ref=1,level=30,profile=baseline,hrd=cbr,crf=20,ratetol=1.0,vbv-maxrate=1200,vbv-bufsize=1200,lookahead=0}}:standard{access=http{mime="video/MP2T"},mux=ts,dst=127.0.0.1:5555}' --qtcapture-width=640 --qtcapture-height=480 --live-caching=200 --intf=macosx
```

#### 2  Run HLS segmenter
```
ffmpeg -v debug -fflags nobuffer -i http://127.0.0.1:5555 -r 15 -vsync 2 -copyts -copytb 1 -codec copy -map 0 -f segment -segment_time 2 -segment_list test.m3u8 -segment_format mpegts -segment_list_flags +live teststream-%09d.ts
```

#### 3  Teardown

`ffmpeg` automatically dies on ingest socket close.



### VLC - UDP Live Stream

VLC pushes TS stream to FFMpeg listening on UDP address/port.

#### 1  Run segmenter listening for stream on private port
```
ffmpeg -v debug -fflags nobuffer -i 'udp://127.0.0.1:5555?fifo_size=1000000&overrun_nonfatal=1' -r 15 -vsync 0 -copyts -copytb 1 -codec copy -map 0 -f segment -segment_time 2 -segment_list test.m3u8 -segment_format mpegts -segment_list_flags +live teststream-%09d.ts
```

#### 2  Run capture pushing data via UDP to segmenter socket
```
/Applications/VLC.app/Contents/MacOS/VLC qtcapture://1 -vvvv --no-drop-late-frames --no-skip-frames --sout='#transcode{vcodec=h264,fps=15,venc=x264{preset=ultrafast,tune=zerolatency,keyint=30,bframes=0,ref=1,level=30,profile=baseline,hrd=cbr,crf=20,ratetol=1.0,vbv-maxrate=1200,vbv-bufsize=1200,lookahead=0}}:standard{access=udp,mux=ts,dst=127.0.0.1:5555}' --sout-mux-caching=0 --sout-udp-caching=0 --sout-udp-group=10 --clock-synchro=1 --sout-ts-shaping=2000 --sout-ts-use-key-frames --qtcapture-width=640 --qtcapture-height=480 --live-caching=200 --intf=macosx
```

#### 3  Teardown

FFMpeg finishes writing the last segment and finalises the m3u8 file on close. Signal 2 is recommended, but `ctrl-c` worked for me as well.

```
# FFMpeg must be killed with signal 2 (I had to send the signal multiple times to be sure the correct thread received it)
kill -2 $PID
```


## RTSP/RTP Live Streaming

In this scenario VLC implements an RTSP server and FFMpeg connects to fetch the stream. Other interactions are possibel with RTSP/RTP, for example, having VLC send RTP over UDP to a multicast address and make FFMpeg join the multicast group. Signalling information in this scenario is passed via SDP (Session Description Protocol), a text based representation of technical details.

Status: __FAILS__

#### 1  Run capture and RTSP/RTP server
```
/Applications/VLC.app/Contents/MacOS/VLC -I dummy --ttl 12 qtcapture://1 -vvvv --no-drop-late-frames --no-skip-frames --sout='#transcode{vcodec=h264,vb=1200,fps=15,venc=x264{preset=ultrafast,tune=zerolatency,keyint=30,bframes=0,ref=1,level=30,profile=baseline,hrd=cbr,crf=20,ratetol=1.0,vbv-maxrate=1200,vbv-bufsize=1200,aud,lookahead=0,repeat-headers=1}}:rtp{name="Teststream1",sdp=rtsp://127.0.0.1:5555/teststream.sdp}' --sout-mux-caching=0 --sout-rtp-caching=0 --clock-synchro=1 --sout-ts-shaping=2000 --sout-ts-use-key-frames --qtcapture-width=640 --qtcapture-height=480 --live-caching=200 --intf=macosx --rtsp-host=127.0.0.1 --rtsp-port=5555
```

#### 2  Run segmenter as RTSP/RTP client connecting to the server
```
ffmpeg -v debug -fflags nobuffer -i rtsp://127.0.0.1:5555/teststream.sdp -r 15 -copyts -copytb 1 -codec copy -map 0 -f segment -segment_time 2 -segment_list test.m3u8 -segment_format mpegts -segment_list_flags +live teststream-%09d.ts
```

When packing raw H264 into RTP VLC does not insert SPS/PPS before every keyframe, so the resulting stream becomes unplayable. When requesting VLC to mux H264 into TS before packaging RTP, FFMpeg fails due to a RTSP handshake problem (SETUP failed: 459 Aggregate operation not allowed).

#### 3  Teardown

_Not tested due to issue above_ FFMpeg should shut down as soon as the RTSP connection is closed (by streaming server sending a TEARDOWN message) or by TCP connection reset resulting from unintended server disconnect.




## Other examples

#### Use FFMpeg to create an HLS compatible segmented stream
```
ffmpeg -v debug -fflags nobuffer -i pipe:0 -vsync 0 -copyts -copytb 1 -codec copy -map 0 -f segment -segment_time 2 -segment_list test.m3u8 -segment_format mpegts -segment_list_flags +live teststream-%09d.ts
```

#### Streaming recorded video files from VLC

Status: __OK__
```
/Applications/VLC.app/Contents/MacOS/VLC ../video/aow-docu-2011.m4v -vvvv --intf=rc --sout '#duplicate{dst=display,dst="transcode{vcodec=h264,vb=2000,fps=25,scale=1,width=640,height=480,acodec=mp4a,ab=128,channels=2,samplerate=44100,venc=x264{keyint=50,ref=1,ratetol=1.0}}:standard{access=udp,mux=ts,dst=127.0.0.1:5555}"'
```

#### Sending VLC Output to stdout
```
# in a stream splitting example
duplicate{dst=file{mux=ts,dst='-'},dst=display}
``

#### Stream to local and pipe to media stream segmenter, which creates index file and can be streamed to iOS
vlc --ttl 12 qtcapture:// -vvv input_stream --sout="#transcode{venc=x264{keyint=60,idrint=2},vcodec=h264,vb=300,acodec=mp4a,ab=32,channels=2, samplerate=22050}:\
duplicate{dst=file{mux=ts,dst='-'},dst=display}" | mediastreamsegmenter -s 10 -D -f /Users/fernyb/hls/live

