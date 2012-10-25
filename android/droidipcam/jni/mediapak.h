#ifndef _MEDIAPAK_H_
#define _MEDIAPAK_H_

#include <stdint.h>
#include "ipcamera.h"

class FlashVideoPackager
{
	public:
		FlashVideoPackager();
		~FlashVideoPackager();

		//Byte align buffer
		inline void resetBuffer()
		{
			mediaLength = 0;
		}
		inline unsigned char *getBuffer()
		{
			return &mediaBuffer[0];
		}
		unsigned int bufferLength()
		{
			return mediaLength;
		}

		void setParameter(int width, int height, int vfr);
		void addVideoHeader(unsigned char *sps, unsigned int sps_size, unsigned char *pps, unsigned int pps_size);
        void addVideoFrame(unsigned char *p, unsigned int length, int intraFlag, uint32_t ts);
		void addAudioFrame(unsigned char *p, unsigned int length, uint32_t ts);

	private:
		void putByte(uint8_t val);
		void putBE32(uint32_t val);
		void putBE64(uint64_t val);
		void putBE24(uint32_t val);
		void putBE16(uint16_t val);
	
		void putTag(const char *tag);
		void putString(const char *str);
		void putDouble(double d);
		void appendData(uint8_t *data, unsigned int size);

	private:
		unsigned char mediaBuffer[MAX_VIDEO_PACKAGE];		//for our hardware the buffer is OK.
		unsigned int mediaLength;
};

#endif
