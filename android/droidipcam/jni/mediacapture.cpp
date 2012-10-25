#include <sys/time.h>
#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <list>
#include <deque>
#include "ipcamera.h"
#include "mediastreamer.h"
#include "mediabuffer.h"
#include "mediapak.h"

// this method support any camera, but will generating error video package.
void MediaStreamer::doCapture2() {
    std::deque<unsigned char> video_check_pattern;
    video_check_pattern.resize(9, 0x00);
    
    unsigned char *buf;
    buf = new unsigned char[1024*512];

    /*
    FILE *fp = fopen ("/sdcard/streaming.flv", "wb");
    flvPackager->setParameter(640, 480, 30);
    flvPackager->addVideoHeader(&mediaInfo.sps_data[0], mediaInfo.sps_data.size(), &mediaInfo.pps_data[0], mediaInfo.pps_data.size());
    fwrite(flvPackager->getBuffer(), flvPackager->bufferLength(), 1, fp);
    flvPackager->resetBuffer();
    */

    LOGD("Native: Begin capture");

    unsigned int last_frame_num = 0;
    int frame_count = 0;
    while(1) {
        if ( infd < 0)
            break;

        // find video slice data from es streaming 
        unsigned char current_byte;
        if ( read(infd, &current_byte, 1) < 0)
            break;
        video_check_pattern.pop_front();
        video_check_pattern.push_back(current_byte);
        
        int slice_type;
        unsigned int frame_num;
        int nal_length = checkSingleSliceNAL( video_check_pattern, slice_type, frame_num );
        if ( nal_length > 0) {
            if ( (slice_type == 0) && (frame_num != (last_frame_num + 1) ) ) {
                LOGD("Error, wrong number, FIXME FIXME");
                {
                    char temp[512];
                    snprintf(temp, 512, "ST=%d,  FN=%d, NAL=%d, LFN=%d, FNL=%d, 0x%02x%02x%02x%02x", 
                            slice_type, frame_num, nal_length, last_frame_num, frame_num_length, 
                            video_check_pattern[5],video_check_pattern[6],video_check_pattern[7],video_check_pattern[8] );
                    LOGD(temp);
                }
                //last_frame_num = frame_num;
                continue;
            }
            last_frame_num = frame_num;
            
                       
            for(int i = 0; i < (int)video_check_pattern.size(); i++) {
                buf[i] = video_check_pattern[i];
            }
            if ( fillBuffer( &buf[video_check_pattern.size()] , nal_length - (video_check_pattern.size() - 4) ) < 0)
                break;
            /*
            flvPackager->addVideoFrame( buf, nal_length + 4, slice_type, frame_count*30);
            fwrite(flvPackager->getBuffer(), flvPackager->bufferLength(), 1, fp);
            flvPackager->resetBuffer();
            */
            mediaBuffer->PushBuffer( buf, nal_length + 4, frame_count*88, slice_type ? MEDIA_TYPE_VIDEO_KEYFRAME : MEDIA_TYPE_VIDEO);

            frame_count++;
        }
    }
    delete buf;

}

static int64_t getCurrentTime() {
    
    struct timeval tv;

    gettimeofday (&tv, NULL);

    return (INT64_C(1000000) * tv.tv_sec) + tv.tv_usec;
}

// this method only support H.264 + AMR_NB
void MediaStreamer::doCapture() {
    unsigned char *buf;
    buf = new unsigned char[MAX_VIDEO_PACKAGE];
    unsigned int aseq = 0;
    unsigned int vseq = 0;
   
    const unsigned int STACK_SIZE = 128;
    std::list<int64_t> timestack;
    unsigned long last_ts = 0;

    // skip none useable heaer bytes
    fillBuffer( buf, mediaInfo.begin_skip);

    // fectching real time video data from camera
    while(1) {
        if ( fillBuffer(buf, 4) < 0)
            break;

checking_buffer:        
        if ( buf[0] == 0x00  ) {
            unsigned int vpkg_len = (buf[1] << 16) + (buf[2] << 8) + buf[3];
            if ( fillBuffer(&buf[4], vpkg_len ) < 0)
              break;
            vpkg_len += 4;
            
            if ( vpkg_len > (unsigned int)MAX_VIDEO_PACKAGE ) {
              LOGD("ERROR: Drop big video frame....");
              vseq++;
              fillBuffer(vpkg_len);
              continue; 
            }

            int slice_type = 0;
            if ( (buf[5] & 0xF8 ) == 0xB8) {
                slice_type = 1;
            } else if ( ((buf[5] & 0xFF) == 0x88) 
                    && ((buf[5] & 0x80) == 0x80) ) {
                slice_type = 1;
            } else if ( (buf[5] & 0xE0) == 0xE0) {
                slice_type = 0;
            } else if ( (buf[5] & 0xFE) == 0x9A) {
                slice_type = 0;
            }
            buf[0] = 0x00;
            buf[1] = 0x00;
            buf[2] = 0x00;
            buf[3] = 0x01; 
           
#if 1        
            // computing the current package's timestamp 
            int64_t cts = getCurrentTime() / 1000;

            if( timestack.size() >= STACK_SIZE) {
                timestack.pop_back();
                timestack.push_front(cts);
            } else {
                timestack.push_front(cts);
            }

            if ( timestack.size() < STACK_SIZE) {
                cts = (timestack.size() - 1) * 100;      // default = 10 fps
                last_ts = (unsigned long)cts;
            } else {
                unsigned long total_ms;
                total_ms = timestack.front() - timestack.back();
                cts = last_ts + total_ms / (STACK_SIZE - 1);
                last_ts = cts;
            }
            mediaBuffer->PushBuffer( buf, vpkg_len, last_ts, slice_type ? MEDIA_TYPE_VIDEO_KEYFRAME : MEDIA_TYPE_VIDEO);
#else
            vseq ++;
            mediaBuffer->PushBuffer( buf, vpkg_len, vseq * 1000 / mediaInfo.video_frame_rate, slice_type ? MEDIA_TYPE_VIDEO_KEYFRAME : MEDIA_TYPE_VIDEO);
#endif
        } else {
            // fetching AMR_NB audio package
            static const unsigned char packed_size[16] = {12, 13, 15, 17, 19, 20, 26, 31, 5, 0, 0, 0, 0, 0, 0, 0};
            unsigned int mode = (buf[0]>>3) & 0x0F;
            unsigned int size = packed_size[mode] + 1;
            if ( size > 4) {
                if ( fillBuffer(&buf[4], size - 4) < 0)
                    break;
                aseq ++;
                //SignalNewPackage(buf, 32, ats, MEDIA_TYPE_AUDIO);
            } else {
                fillBuffer(&buf[4], size );
                for(int i = 0; i < 4; i++) 
                    buf[i] = buf[size+i];
                //SignalNewPackage(buf, 32, ats, MEDIA_TYPE_AUDIO);
                goto checking_buffer;
            }
        }
    }
    delete buf;
}


int MediaStreamer::checkSingleSliceNAL(const std::deque<unsigned char> &pattern , int &slice_type, unsigned int &frame_num) {   
    
    // 1. first we check NAL's size, valid size should less than 192K 
    if ( pattern[0] != 0x00)
        return -1;                          
    if ( pattern[1] != 0x00)   
        return -1;

    // 2. check NAL header including NAL start and type,
    //    only nal_unit_type = 1 and 5 are selected
    //    nal_ref_idc > 0
    if (   (pattern[4] != 0x21)
        && (pattern[4] != 0x25)
        && (pattern[4] != 0x41)
        && (pattern[4] != 0x45)
        && (pattern[4] != 0x61)
        && (pattern[4] != 0x65) ) {
        return -1; 
    }  

    // 3. checking fist_mb (should be 0), slice type should be I or P, 
    //    frame_num should be continued. 
    // Only following pattens are supported. 
    //  
    // I Frame: b 1011 1***   (first_mb = 0, slice_type = 2, pps_id = 0)   
    // I Frame: b 1000 1000 1 (first_mb = 0, slice_type = 7, pps_id = 0)
    // P Frame: b 111* ****   (first_mb = 0, slice_type = 0, pps_id = 0)
    // P Frame: b 1001 101*   (first_mb = 0, slice_type = 5, pps_id = 0)
    // 
    int frame_num_skip = -1;
    if ( (pattern[5] & 0xF8 ) == 0xB8) {
        slice_type = 1;
        frame_num_skip = 5;
    } else if ( ((pattern[5] & 0xFF) == 0x88) 
                && ((pattern[6] & 0x80) == 0x80) ) {
        slice_type = 1;
        frame_num_skip = 9;
    } else if ( (pattern[5] & 0xE0) == 0xE0) {
        slice_type = 0;
        frame_num_skip = 3;
    } else if ( (pattern[5] & 0xFE) == 0x9A) {
        slice_type = 0;
        frame_num_skip = 7;
    }
    if ( slice_type == -1)
        return -1;

    if ( frame_num_length == -1) {
        frame_num_length = 0;
        frame_num = 0;
    } else if ( frame_num_length == 0) {
        unsigned int bits = (pattern[5] << 24) + (pattern[6] << 16) + (pattern[7] << 8) + pattern[8];
        bits = bits << frame_num_skip;
        for(int i = 0; i < (31 - frame_num_skip); i++) {
            if ( bits & 0x80000000 ) {
                frame_num_length = i + 1; 
                break;
            }
            bits = bits << 1;
        }
    }
    
    if ( frame_num_length > 0 ) {
        unsigned int bits = (pattern[5] << 24) + (pattern[6] << 16) + (pattern[7] << 8) + pattern[8];
        bits = bits << frame_num_skip;
        bits = bits >> ( 32 - frame_num_length );
        frame_num =  bits;
    }

    int nal_length = (pattern[1] << 16)  + (pattern[2] << 8) + pattern[3];
    return nal_length;
}

int MediaStreamer::fillBuffer(unsigned char *buf, unsigned int len) {
  while(len > 0) {
    int ret = recv(infd, buf, len, 0);
    
    if ( ret < 0)
        return -1;
    if ( ret == 0)
        continue;
    
    if ( infd < 0)
        return -1;

    len -= ret;
    buf += ret;
  }
  return len;
}

int MediaStreamer::fillBuffer(unsigned int len) {
  unsigned char temp[1024];
  int ret;
  while(len > 0) {
    if ( len > 1024)
        ret = recv(infd, temp, 1024, 0);
    else
        ret = recv(infd, temp, len, 0);
    
    if ( ret < 0)
        return -1;
    if ( ret == 0)
        continue;
    
    if ( infd < 0)
        return -1;

    len -= ret;
  }
  return len;
}


