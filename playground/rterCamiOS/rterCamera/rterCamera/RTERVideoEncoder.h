//
//  RTERVideoEncoder.h
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-12.
//  Copyright (c) 2013 rtER. All rights reserved.
//
//  based on https://github.com/chrisballinger/FFmpeg-iOS-Encoder
//  by Christopher Ballinger

#import <Foundation/Foundation.h>
#import <AVFoundation/AVFoundation.h>

#include <libavutil/opt.h>
#include <libavcodec/avcodec.h>
#include <libavutil/channel_layout.h>
#include <libavutil/common.h>
#include <libavutil/imgutils.h>
#include <libavutil/mathematics.h>
#include <libavutil/samplefmt.h>


@interface RTERVideoEncoder : NSObject {
    CMFormatDescriptionRef formatDescription;
    AVCodec *codec;
    AVCodecContext *c;
    AVFrame *frame;
    
    int frameNumber, ret, got_output;
    FILE *f;
    CMVideoDimensions outputSize;
    CMVideoDimensions inputSize;
    struct SwsContext *sws_ctx;
    AVFrame *scaledFrame;
}

- (void) setupEncoderWithFormatDescription:(CMFormatDescriptionRef)formatDescription;
- (void) setupEncoderWithFormatDescription:(CMFormatDescriptionRef)newFormatDescription desiredOutputSize:(CMVideoDimensions)desiredOutputSize;
- (void) setupEncoderWithDimesions:(CMVideoDimensions)dimensions;
- (void) finishEncoding;
- (int) encodeSampleBuffer:(CMSampleBufferRef)sampleBuffer
                    output:(AVPacket *)pkt;
-(void) freePacket:(AVPacket *)pkt;
-(void) freeEncoder;


@property (nonatomic) BOOL readyToEncode;
@property (nonatomic) CMFormatDescriptionRef formatDescription;

@end
