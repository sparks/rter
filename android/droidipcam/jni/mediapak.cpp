#include <stdio.h>
#include <string.h>
#include "mediapak.h"

#define AMF_END_OF_OBJECT         0x09
#define FLV_VIDEO_FRAMETYPE_OFFSET   4

#define VIDEO_ONLY	1

typedef enum
{
	AMF_DATA_TYPE_NUMBER      = 0x00,
	AMF_DATA_TYPE_BOOL        = 0x01,
	AMF_DATA_TYPE_STRING      = 0x02,
	AMF_DATA_TYPE_OBJECT      = 0x03,
	AMF_DATA_TYPE_NULL        = 0x05,
	AMF_DATA_TYPE_UNDEFINED   = 0x06,
	AMF_DATA_TYPE_REFERENCE   = 0x07,
	AMF_DATA_TYPE_MIXEDARRAY  = 0x08,
	AMF_DATA_TYPE_OBJECT_END  = 0x09,
	AMF_DATA_TYPE_ARRAY       = 0x0a,
	AMF_DATA_TYPE_DATE        = 0x0b,
	AMF_DATA_TYPE_LONG_STRING = 0x0c,
	AMF_DATA_TYPE_UNSUPPORTED = 0x0d,
} AMFDataType;

enum
{
	FLV_TAG_TYPE_AUDIO = 0x08,
	FLV_TAG_TYPE_VIDEO = 0x09,
	FLV_TAG_TYPE_META  = 0x12,
};

enum
{
	FLV_FRAME_KEY   = 1 << FLV_VIDEO_FRAMETYPE_OFFSET | 7,
	FLV_FRAME_INTER = 2 << FLV_VIDEO_FRAMETYPE_OFFSET | 7,
};

enum
{
	FLV_CODECID_PCM = 0,
	FLV_CODECID_H264 = 7,
};

static uint64_t dbl2int( double value )
{
    void *p = &value;
	return *(uint64_t*)p;
}

static uint16_t endian_fix16( uint16_t x )
{
	return (x<<8)|(x>>8);
}

static uint32_t endian_fix32( uint32_t x )
{
	return (x<<24) + ((x<<8)&0xff0000) + ((x>>8)&0xff00) + (x>>24);
}

FlashVideoPackager::FlashVideoPackager()
{
	resetBuffer();

    putTag( "FLV" ); 	// Signature
	putByte( 1 );    	// Version

#ifdef VIDEO_ONLY	
	putByte( 1 );   	// Video only
#else
	putByte( 1 + 4 );   // Video and Audio
#endif

	putBE32( 9 );    	// DataOffset
	putBE32( 0 );    	// PreviousTagSize0
}

FlashVideoPackager::~FlashVideoPackager()
{

}

inline void FlashVideoPackager::appendData(uint8_t *data, unsigned int size)
{
	memcpy(&mediaBuffer[mediaLength], data, size);
	mediaLength += size;
}

void FlashVideoPackager::putByte(uint8_t val)
{
	mediaBuffer[mediaLength] = val;
	mediaLength++;
}

void FlashVideoPackager::putBE16(uint16_t val)
{
	val = endian_fix16(val);
	appendData((unsigned char *)&val, 2);
}

void FlashVideoPackager::putBE24(uint32_t val)
{
	putBE16( val >> 8 );
	putByte( val );	
}

void FlashVideoPackager::putBE32(uint32_t val)
{
	val = endian_fix32(val);
	appendData((uint8_t *)&val, 4);
}

void FlashVideoPackager::putBE64(uint64_t val)
{
	putBE32( val >> 32 );
	putBE32( val );
}

void FlashVideoPackager::putTag(const char *tag)
{
	while( *tag )
		putByte( *tag++ );
}

void FlashVideoPackager::putString(const char *str)
{
	uint16_t len = strlen( str );
	putBE16( len );
	appendData( (uint8_t*)str, len );
}

void FlashVideoPackager::putDouble(double d)
{
	putByte( AMF_DATA_TYPE_NUMBER );
	putBE64( dbl2int( d ) );
}

void FlashVideoPackager::setParameter(int width, int height, int vfr)
{
	putByte(FLV_TAG_TYPE_META ); // Tag Type "script data"

	//the header data length should be tagLength which will computed later
#ifdef VIDEO_ONLY
	putBE24( 182 ); 			// header data length 
#else
	putBE24( 290 );				// header data length
#endif
	putBE24( 0 ); 				// timestamp
	putBE32( 0 ); 				// reserved
	
	int startPos = mediaLength;

	//header data begin
	putByte( AMF_DATA_TYPE_STRING );
	putString( "onMetaData" );

	putByte( AMF_DATA_TYPE_MIXEDARRAY );

#ifdef VIDEO_ONLY
	putBE32( 5 + 2 );		// +2 for duration and file size
#else
	putBE32( 5 + 5 + 2 );			// +2 for duration and file size
#endif

	putString( "duration" );
	putDouble( 0 ); 			//FIXME written at end of encoding

	//Video Information
	{
		putString( "width" );
		putDouble( width );

		putString( "height" );
		putDouble( height );

		putString( "framerate" );
		//putDouble( (double)vfr*1.0 );
		putDouble(0.0);

		putString( "videocodecid" );
		putDouble( FLV_CODECID_H264 );

		putString( "videodatarate" );
		putDouble( 0 ); 				// written at end of encoding
	}

#ifndef VIDEO_ONLY
	//Audio Information
	{
		putString( "audiodatarate");	
		putDouble( 64000 / 1024.0);		//6.4K 

		putString( "audiosamplerate");
		//putDouble( 8000);				//8K HZ PCM
		putDouble( 5112);				//5112 HZ PCM

		putString( "audiosamplesize");
		putDouble( 8);					//8 bit PCM_U8

		putString( "stereo");			//Only One Channel
		putByte( 0); 

		putString( "audiocodecid");
		putDouble( 0);					//PCM_U8 = 0
	}
#endif
	
	putString( "encoder" );
	putByte( AMF_DATA_TYPE_STRING);
	putString( "Lavf52.87.1");

	putString( "filesize" );
	putDouble( 0 ); 					//FIXME written at end of encoding

	putString( "" );
	putByte( AMF_END_OF_OBJECT );
	//header data end

	int endPos = mediaLength;
	int tagLength = endPos - startPos;

	//fprintf(stderr,"tagLength = %d\n", tagLength );
	putBE32( tagLength + 11); 				// total tag length
}


//Because first frame including SPS PPS information
void FlashVideoPackager::addVideoHeader(unsigned char *sps, unsigned int sps_size, unsigned char *pps, unsigned int pps_size)
{
    unsigned int frame_size = sps_size + pps_size + 16;

    //frame header infomation
	putByte( FLV_TAG_TYPE_VIDEO );
	putBE24( frame_size ); // frame data szie
	putBE24( 0 );  // timestamp
	putByte( 0 );  // timestamp extended
	putBE24( 0 );  // StreamID - Always 0

	//frame data : SPS
	putByte( 7 | FLV_FRAME_KEY ); // Frametype and CodecID
	putByte( 0 );      // AVC sequence header
	putBE24( 0 );      // composition time
	putByte( 1 );      // version
	putByte( sps[1] ); // profile
	putByte( sps[2] ); // profile
	putByte( sps[3] ); // level
	putByte( 0xff );   // 6 bits reserved (111111) + 2 bits nal size length - 1 (11)
	putByte( 0xe1 );   // 3 bits reserved (111) + 5 bits number of sps (00001)
	putBE16( sps_size );	//SPS size = 9
	appendData( sps, sps_size );	

	//frame data : PPS
	putByte( 1 ); 		// number of pps
	putBE16( pps_size);
	appendData( pps, pps_size);

	putBE32( frame_size + 11 );	 	// Last tag size
}

//Because first frame including SPS PPS information
void FlashVideoPackager::addVideoFrame(unsigned char *p, unsigned int length, int intraFlag, uint32_t ts)
{
	p[0] = ((length-4) & 0xFF000000)>>24;
	p[1] = ((length-4) & 0x00FF0000)>>16;
	p[2] = ((length-4) & 0x0000FF00)>>8;
	p[3] = (length-4) & 0x000000FF;

	// A new frame - write packet header
	putByte( FLV_TAG_TYPE_VIDEO );
	putBE24( 5 + length ); 			// calculated later
	putBE24( ts );
	putByte( ts >> 24 );
	putBE24( 0 );

	putByte( intraFlag ? FLV_FRAME_KEY : FLV_FRAME_INTER );
	putByte( 1 ); 					// AVC NALU, data
	putBE24( 0 );					// offset  = 0

	appendData( p, length);						

	putBE32( 11 + 5 + length); 		// Last tag size
}

void FlashVideoPackager::addAudioFrame(unsigned char *p, unsigned int length, uint32_t ts)
{
	putByte(FLV_TAG_TYPE_AUDIO); 
	putBE24(length + 1);
	putBE24(ts);			//TS
	putByte(ts>>24);		//TS
	putBE24(0x00);			//flv->reserved

	putByte(0x00);			//flags
	appendData( p, length );

    putBE32( 12 + length ); // Last tag size
}

