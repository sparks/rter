#ifndef _MEDIASTREAMER_H_
#define _MEDIASTREAMER_H_

#include "talk/base/thread.h"
#include "talk/base/messagequeue.h"

namespace talk_base {
    class Thread;
};

class MediaBuffer;

class MediaStreamer : public sigslot::has_slots<>, public talk_base::MessageHandler {  
public:
    MediaStreamer(int ifd, int ofd);
    ~MediaStreamer();
    void Start();
    void Stop();

protected:    
    virtual void OnMessage(talk_base::Message *msg);

    void doCapture();
    void doCapture2();
    int checkSingleSliceNAL(const std::deque<unsigned char> &pattern , int &slice_type, unsigned int &frame_num);
    int fillBuffer(unsigned char *buf, unsigned int len);
    int fillBuffer(unsigned int len);
    int flushBuffer(unsigned char *buf, unsigned int len);
    
    void doStreaming();

private:
    enum {
        MSG_BEGIN_CAPTURE_TASK,
        MSG_BEGIN_STREAMING_TASK,
    };

    int infd;
    int outfd;    
    int frame_num_length;

    talk_base::Thread *captureThread;
    talk_base::Thread *streamingThread;
public:
    static MediaBuffer *mediaBuffer;
    static MediaStreamer *mediaStreamer;    
};

#endif
