#include "ipcamera.h"
#include "mediastreamer.h"
#include "mediabuffer.h"
#include "mediapak.h"

MediaBuffer* MediaStreamer::mediaBuffer = NULL;
MediaStreamer* MediaStreamer::mediaStreamer = NULL;

int StartStreamingMedia(int infd, int outfd) {

    if ( MediaStreamer::mediaBuffer == NULL)
        MediaStreamer::mediaBuffer = new MediaBuffer(32, 120, MAX_VIDEO_PACKAGE, 1024); 
    MediaStreamer::mediaBuffer->Reset(); 

    if ( MediaStreamer::mediaStreamer != NULL) {
        LOGD("Before delete mediaStreamer");
        delete MediaStreamer::mediaStreamer;
    }
    MediaStreamer::mediaStreamer = new MediaStreamer(infd, outfd);
    MediaStreamer::mediaStreamer->Start();

    return 1;
}

void StopStreamingMedia() {
    // release object in the begin of netxt request
    if ( MediaStreamer::mediaStreamer != NULL)  
        MediaStreamer::mediaStreamer->Stop();
}

MediaStreamer::MediaStreamer(int ifd, int ofd) {
    frame_num_length =1;
    infd = ifd;
    outfd = ofd;

    captureThread = NULL;
    streamingThread = NULL;    
}

MediaStreamer::~MediaStreamer() {
    if ( streamingThread != NULL) {
        delete streamingThread;
    } 
    if ( captureThread != NULL) {
        delete captureThread;
    }
}

void MediaStreamer::Start() {
    captureThread = new talk_base::Thread();
    captureThread->Start();
    captureThread->Post(this, MSG_BEGIN_CAPTURE_TASK);

    streamingThread = new talk_base::Thread();
    streamingThread->Start();
    streamingThread->Post(this, MSG_BEGIN_STREAMING_TASK);
}

void MediaStreamer::Stop() {
    infd = -1;
    outfd = -1;

    if ( streamingThread != NULL)
        streamingThread->Quit();

    if ( captureThread != NULL)
        captureThread->Quit();
}

void MediaStreamer::OnMessage(talk_base::Message *msg) {
    switch( msg->message_id) {
        case MSG_BEGIN_CAPTURE_TASK:
            doCapture();        
            break;

        case MSG_BEGIN_STREAMING_TASK:
            doStreaming();
            break;

        default:
            break;
    }
}

