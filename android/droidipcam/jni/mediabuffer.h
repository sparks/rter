#ifndef _MEDIABUFFER_H_
#define _MEDIABUFFER_H_

#include <list>
#include <vector>
#include "talk/base/sigslot.h"
#include "talk/base/criticalsection.h"

enum MEDIA_TYPE {
  MEDIA_TYPE_AUDIO,
  MEDIA_TYPE_VIDEO, 
  MEDIA_TYPE_VIDEO_KEYFRAME,
};

const int MUX_OFFSET = 32;
struct MediaPackage{
public:
  MediaPackage(const unsigned int size) {
	data = new unsigned char[size + MUX_OFFSET];
    data += MUX_OFFSET;
  };
  ~MediaPackage() {
    if ( data != NULL) {
      data = data - MUX_OFFSET;
      delete data;
    }
  }
  
  unsigned char *data;
  unsigned int length;
  unsigned int ts;
  unsigned int seq;
  MEDIA_TYPE media_type;	
};

class MediaBuffer: public sigslot::has_slots<> {
public:
  MediaBuffer(const unsigned int vnum, const unsigned int anum, const unsigned int vsize, const unsigned int asize);
  ~MediaBuffer();
  void Reset();

  // access from diffrent threads, they are safe.
  bool PushBuffer(const unsigned char *d, const unsigned int len, const unsigned int ts, const MEDIA_TYPE mt);
  bool PullBuffer(MediaPackage **ppkg, const MEDIA_TYPE mt);

  unsigned int VBufferSize(){
    return vbuffer_.size();
  }
  unsigned int ABufferSize() {
    return abuffer_.size();
  }

private:
  bool pushAudioPackage(const unsigned char *d, const unsigned int len, const unsigned int ts);
  bool pushVideoPackage(const unsigned char *d, const unsigned int len, const unsigned int ts, const unsigned int isIntra);

private:
  unsigned int vpkg_size_;
  unsigned int apkg_size_;

  unsigned int vpkg_seq_;
  
  std::list<MediaPackage*> vbuffer_;
  std::list<MediaPackage*> abuffer_;
  MediaPackage *vpkg_released;
  MediaPackage *apkg_released;

  std::vector<MediaPackage*> vpkg_pool_;
  std::vector<MediaPackage*> apkg_pool_;

  talk_base::CriticalSection mutex_; 
};

#endif
