#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <deque>
#include "ipcamera.h"
#include "mediastreamer.h"
#include "mediabuffer.h"
#include "mediapak.h"

int MediaStreamer::flushBuffer(unsigned char *buf, unsigned int len) {
  while(len > 0) {
    int ret = send(outfd, buf, len, 0);
    
    if ( ret < 0)
        return -1;
    if ( ret == 0) {
        usleep(10);
        continue;
    }

    if ( outfd < 0)
        return -1;

    len -= ret;
    buf += ret;
  }
  return len;
}

void MediaStreamer::doStreaming() {

    LOGD("Native: Begin streaming");

    MediaPackage *media_package;
    FlashVideoPackager *flvPackager = new FlashVideoPackager();
    flvPackager->setParameter(mediaInfo.video_width, mediaInfo.video_height, 30);
    flvPackager->addVideoHeader(&mediaInfo.sps_data[0], mediaInfo.sps_data.size(), &mediaInfo.pps_data[0], mediaInfo.pps_data.size());
 
    while(1) {
        if ( outfd < 0)
            break;

        media_package = NULL;
        if ( mediaBuffer->PullBuffer(&media_package, MEDIA_TYPE_VIDEO) == false) {
            talk_base::Thread::SleepMs(50);             // wait for 1/20 second
            continue;
        }
        
        flvPackager->addVideoFrame( media_package->data, 
                                    media_package->length,
                                    (media_package->media_type == MEDIA_TYPE_VIDEO_KEYFRAME),
                                    media_package->ts);
        int ret = flushBuffer(flvPackager->getBuffer(), flvPackager->bufferLength());
        if ( ret < 0)
            break;

        flvPackager->resetBuffer();
    }

    delete flvPackager;    
}
