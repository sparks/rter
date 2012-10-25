#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <string>
#include <vector>
#include <deque>
#include "ipcamera.h"


MediaCheckInfo mediaInfo;

/*
   put_byte(pb, 1);      // version 
   put_byte(pb, sps[1]); // profile 
   put_byte(pb, sps[2]); // profile compat 
   put_byte(pb, sps[3]); // level 
   put_byte(pb, 0xff);   // 6 bits reserved (111111) + 2 bits nal size length - 1 (11)

   put_byte(pb, 0xe1);   // 3 bits reserved (111) + 5 bits number of sps (00001) 


   put_be16(pb, sps_size);
   put_buffer(pb, sps, sps_size);
   put_byte(pb, 1);      // number of pps 
   put_be16(pb, pps_size);
   put_buffer(pb, pps, pps_size);      
*/

void getHeader(FILE *fp) {
    unsigned char temp_buffer[1024];
    long current_pos = ftell(fp);
    fread(temp_buffer, 1, 1024, fp);
    fseek(fp, current_pos, SEEK_SET);    

    int sps_size = (temp_buffer[5] << 8) + temp_buffer[6];
    for(int i = 0; i < sps_size; i++) {
        mediaInfo.sps_data.push_back( temp_buffer[7+i] );
    }

    int pps_size = (temp_buffer[8 + sps_size] << 8) + temp_buffer[9 + sps_size];
    for(int i = 0; i < pps_size; i++) {
        mediaInfo.pps_data.push_back( temp_buffer[10 + sps_size + i] );
    }
}

/*
713     if( p_box->data.p_mdhd->i_version )
714     {
715         MP4_GET8BYTES( p_box->data.p_mdhd->i_creation_time );
716         MP4_GET8BYTES( p_box->data.p_mdhd->i_modification_time );
717         MP4_GET4BYTES( p_box->data.p_mdhd->i_timescale );
718         MP4_GET8BYTES( p_box->data.p_mdhd->i_duration );
719     }
720     else
721     {
722         MP4_GET4BYTES( p_box->data.p_mdhd->i_creation_time );
723         MP4_GET4BYTES( p_box->data.p_mdhd->i_modification_time );
724         MP4_GET4BYTES( p_box->data.p_mdhd->i_timescale );
725         MP4_GET4BYTES( p_box->data.p_mdhd->i_duration );
726     }
*/
unsigned int getTimeScale(FILE *fp, unsigned char version) {
    unsigned char temp_buffer[1024];
    unsigned int ret;
    
    if ( version ) {
        fread(temp_buffer, 1, 3 + 8 + 8 + 4, fp);
    
        ret = temp_buffer[22] +
            (temp_buffer[21] << 8) +
            (temp_buffer[20] << 16) +
            (temp_buffer[19] << 24);
    } else {
        fread(temp_buffer, 1, 3 + 4 + 4 + 4, fp);
        ret = temp_buffer[14] +
            (temp_buffer[13] << 8) +
            (temp_buffer[12] << 16) +
            (temp_buffer[11] << 24);
    } 

    return ret;
}
/*
MP4_GETVERSIONFLAGS( p_box->data.p_stts );
MP4_GET4BYTES( p_box->data.p_stts->i_entry_count );
MP4_GET4BYTES( p_box->data.p_stts->i_sample_count[i] );
MP4_GET4BYTES( p_box->data.p_stts->i_sample_delta[i] );
*/
unsigned int getFrameRate(FILE *fp, unsigned int time_scale) {
    unsigned char temp_buffer[1024];
    fread(temp_buffer, 1, 1024, fp);
    
    unsigned int entry_count = temp_buffer[7] +
                              (temp_buffer[6] << 8) +
                              (temp_buffer[5] << 16) +
                              (temp_buffer[4] << 24);

    if ( entry_count > 0) {
        unsigned int duration = temp_buffer[15] +
                              (temp_buffer[14] << 8) +
                              (temp_buffer[13] << 16) +
                              (temp_buffer[12] << 24);
        
        std::cout << "duration = " << duration << std::endl;
        mediaInfo.video_frame_rate = (int)( time_scale * 1.0 / duration + 0.5);
        return duration;
    }
    
    return 0;
}

int CheckMedia(const int wid, const int hei, const std::string mp4_file) {
    mediaInfo.video_width = wid;
    mediaInfo.video_height = hei;
    mediaInfo.audio_codec = -1;
    mediaInfo.video_frame_rate = -1;

    std::deque<unsigned char> mdat;
    std::deque<unsigned char> avcC;
    std::deque<unsigned char> stts;

    for(int i = 0; i < 4; i++) {
        mdat.push_back(0);
        avcC.push_back(0);
        stts.push_back(0);
    }
    avcC.push_back(0);

    mediaInfo.begin_skip = -1;
    FILE *fp = fopen( mp4_file.c_str(), "rb");
    unsigned char c;

    // 0. find av data(mdat) start position
    int pos = 0;
    while( !feof(fp) ) {
        c = fgetc(fp);
        pos ++;
        mdat.push_back(c);
        mdat.pop_front();
        if (    mdat[0] == 'm'
                && mdat[1] == 'd'
                && mdat[2] == 'a'
                && mdat[3] == 't') {
            mediaInfo.begin_skip = pos;
            std::cout << "Found MDAT skipped = "<< pos << std::endl;
            break; 
        }
    } 

    if ( mediaInfo.begin_skip < 0)
        return 0;

    mediaInfo.sps_data.clear();
    mediaInfo.pps_data.clear();
    fseek(fp, 0l, SEEK_SET);
    unsigned int time_scale = 0;
    unsigned int frame_duration = 0;
    bool avcFind = false;

    // 1. get sps&pps and time scale
    while( !feof(fp) ) {
        c = fgetc(fp);
        avcC.push_back(c);
        avcC.pop_front();
        if (    avcC[0] == 'a'
                && avcC[1] == 'v'
                && avcC[2] == 'c'
                && avcC[3] == 'C'
                && avcC[4] == 0x01) {
            avcFind = true;
            getHeader(fp);
            break;
        } else if (avcC[0] == 'm'
                && avcC[1] == 'd'
                && avcC[2] == 'h'
                && avcC[3] == 'd') {
            time_scale = getTimeScale(fp, avcC[4]);
        }
    }  

    if ( !avcFind ) {
        return 0;
    }
    
    std::cout << "time_scale = " << time_scale << std::endl;
    
    // 2. trying to find the stss atom 
    while( !feof(fp) ) {
        c = fgetc(fp);
        stts.push_back(c);
        stts.pop_front();
        if ( stts[0] == 's' &&
             stts[1] == 't' && 
             stts[2] == 't' && 
             stts[3] == 's') {
            frame_duration = getFrameRate(fp, time_scale);
            break;
        }

    }    

    std::cout << "frame rate = " << mediaInfo.video_frame_rate << std::endl;

    {
        char temp[512];
        sprintf(temp, "skip = %d, time_scale = %d, duration = %d, frame_rate = %d\n", 
                        mediaInfo.begin_skip, 
                        time_scale,
                        frame_duration, 
                        mediaInfo.video_frame_rate);
        LOGD(temp);
    }

    return 1;
}

#if 0
int main(int argc, const char *argv[]) {
    if ( argc > 1) {
        std::string fname = argv[1];
        CheckMedia(0, 0, fname);
    }
}
#endif

