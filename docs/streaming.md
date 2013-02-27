# Live Streaming to HTML5 Browsers

Here is a collection of resources about live video streaming from a networked anywhere to HTML5 web browsers. The technologies considered are

* DASH (Dynamic Adaptive Streaming over HTTP, an MPEG Standard)
* HLS (HTML Live Streaming, Apple's implementation in Quicktime and iOS)

Terminology used throughout this document

* __Streaming Source__, a device/app capturing a single live video feed
* __Streaming Server__, a central entity distributing multiple live video feeds
* __Streaming Client__, a HTML5 browser displaying one or multiple video feeds

### DASH, Dynamic Adaptive Streaming over HTTP

[DASH]() is a recent standard for segmented video streaming supporting adaptation through bitstream switching. DASH specifies a description format (MPD files), requires a certain packetization  format (MPEG2-TS or ISO-MBFF), a subset of codecs (H264, AAC, VP8?), and HTTP as transport, however, DASH does __not__ specify transport or adaptation strategies, buffering or playback interactions (which all are left to implementations). A normal Web Server can be used to deliver segmented files to browsers which makes elegant use of already deployed infrastructure including CDNs and caches. DASH further defines different profiles, e.g. for live and on-demand streaming which basically refine applicable limits on configuration settings.

To use DASH one nees a server-side implementation for generating compliant segmented video files plus the MPD description. At a Web client an implementation is required to fetch and parse the MPD and download segments as required.

#### Software

* Segmentation and MPD generation: [GPAC](2) in a recent version from [SVN](3), for more info on DASH support see [here](4) and [here](5)
* Delivery: any HTTP server
* Playback: Google Chrome with Media Source Extention plus a 3rd party DASH Javascript Engine, e.g. [DASH-JS](6)
* Firefox with built-in DASH client support (version?)
* [Decoder Test Sequences](http://dash-mse-test.appspot.com/decoder-test.html)
* [DASH Test Sequences](http://gpac.wp.mines-telecom.fr/2012/02/23/dash-sequences/)


#### Live Streaming of DASH using GPAC

[GPAC DASH Options](http://gpac.wp.mines-telecom.fr/mp4box/dash/)

Call `MP4Box` on regular basis with the segments you need to add to the live session:

```
MP4Box [DASH_OPTIONS] -dash-ctx dasher.cfg fileX.mp4
```

#### GPAC DASH Example from File
```
~/opt/bin/MP4Box -dash 2000 -out aow.mpd -dash-profile live -rap -frag-rap -segment-name aow-dash-Num$Number$-BW$Bandwidth$-Time$Time$ -base-url http://tsunamijs.kidtsunami.com/assets/dash -mpd-title "AOW Documentary" -mpd-source "Some Source" -mpd-info-url "InfoURL" -dash-ctx aow-context.bin ../video/aow-docu-2011.m4v
```
Todo: try options `-dash-live=aow-context.bin 2000` and `-noprog`

#### GPAC Live DASH Example from 'Capture to File'
```
~/opt/bin/MP4Box -dash 20000 -out aow.mpd -dash-profile live -rap -frag-rap -segment-name aow-dash-Num$Number$-BW$Bandwidth$-Time$Time$ -base-url http://tsunamijs.kidtsunami.com/assets/dash -mpd-title "AOW Documentary" -mpd-source "Some Source" -mpd-info-url "InfoURL" -dash-live=aow-context.bin 2000 vlc-output.h264
```
__Problem:__ MP4Box sleep interval doubled each time, probably because the file kept expanding


### HLS, HTML Live Streaming

[HLS](1) is a proprietary 'standard' by Apple to feed live video content to iOS devices and Quicktime media player. Like DASH, HLS separates video files into segments for individual download. HLS requires MPEG2-TS and H264/AAC. The description format is an extended version of the MP3 playlist format M3U called M3U8.


#### Software

* Segmentation and M3U8 generation: FFmpeg, FFServer
* Playback: Safari, iOS (see also [HLS support](7))

#### FFMPEG HLS Examples
```
./ffmpeg -v 9 -loglevel 99 -re -i sourcefile.avi -an \
-c:v libx264 -b:v 128k -vpre ipod320 \
-flags -global_header -map 0 -f segment -segment_time 4 \
-segment_list test.m3u8 -segment_format mpegts stream%05d.ts
```

```
# not sure if its working
ffmpeg -i encoded.mp4 -c copy -map 0 -vbsf h264_mp4toannexb -f segment -segment_time 10 -segment_format mpegts stream%d.ts
```

### FFSERVER as streaming server
* [FFMpeg Streaming Guide](http://ffmpeg.org/trac/ffmpeg/wiki/StreamingGuide)
* [FFServer](http://ffmpeg.org/ffserver.html)
* [FFServer Example](http://ffmpeg.org/trac/ffmpeg/wiki/Streaming%20media%20with%20ffserver)
* [Live Streaming with FFMpeg](http://sonnati.wordpress.com/2012/07/02/ffmpeg-the-swiss-army-knife-of-internet-streaming-part-v/)
* [Another live streaming example](http://www.onvos.com/http-live-streaming-howto.html)



[1](http://developer.apple.com/library/ios/#documentation/NetworkingInternet/Conceptual/StreamingMediaGuide/UsingHTTPLiveStreaming/UsingHTTPLiveStreaming.html)
[2](http://gpac.wp.mines-telecom.fr/mp4box/dash/)
[3](http://sourceforge.net/projects/gpac/develop)
[4](http://gpac.wp.mines-telecom.fr/2011/02/02/mp4box-fragmentation-segmentation-splitting-and-interleaving/)
[5](http://gpac.wp.mines-telecom.fr/2012/02/01/dash-support/)
[6](http://www-itec.uni-klu.ac.at/dash/?p=792)
[7](http://www.longtailvideo.com/html5/hls)
