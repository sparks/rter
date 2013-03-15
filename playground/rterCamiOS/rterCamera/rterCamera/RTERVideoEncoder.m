//
//  RTERVideoEncoder.m
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-12.
//  Copyright (c) 2013 rtER. All rights reserved.
//
//  based on https://github.com/chrisballinger/FFmpeg-iOS-Encoder
//  by Christopher Ballinger

#import "RTERVideoEncoder.h"

@implementation RTERVideoEncoder

@synthesize readyToEncode, formatDescription;

/*** super ***/
- (id) init {
    if (self = [super init]) {
        readyToEncode = NO;
    }
    return self;
}

//- (void) setupEncoderWithFormatDescription:(CMFormatDescriptionRef)newFormatDescription {
//    formatDescription = newFormatDescription;
//    CFRetain(formatDescription);
//    readyToEncode = YES;
//}
//- (void) finishEncoding {
//    CFRelease(formatDescription);
//    formatDescription = NULL;
//    readyToEncode = NO;
//}
//- (void) encodeSampleBuffer:(CMSampleBufferRef)sampleBuffer {}

/*** new ***/
//- (void) setupEncoderWithFormatDescription:(CMFormatDescriptionRef)newFormatDescription {
//    CMVideoDimensions dimensions;
//    dimensions.width = 320;
//    dimensions.height = 240;
//    [self setupEncoderWithFormatDescription:newFormatDescription desiredOutputSize:dimensions];
//}
//
//- (void) setupEncoderWithFormatDescription:(CMFormatDescriptionRef)newFormatDescription desiredOutputSize:(CMVideoDimensions)desiredOutputSize {
//    inputSize = CMVideoFormatDescriptionGetDimensions(newFormatDescription);
//    outputSize = desiredOutputSize;
//    c = NULL;
//    frameNumber = 0;
//    int codec_id = AV_CODEC_ID_MPEG4;
//    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
//    NSString *basePath = ([paths count] > 0) ? [paths objectAtIndex:0] : nil;
//    NSString *movieName = [NSString stringWithFormat:@"%f.mpg",[[NSDate date] timeIntervalSince1970]];
//    const char *filename = [[NSString stringWithFormat:@"%@/%@", basePath, movieName] UTF8String];
//    printf("Encode video file %s\n", filename);
//    
//    /* find the mpeg1 video encoder */
//    codec = avcodec_find_encoder(codec_id);
//    if (!codec) {
//        fprintf(stderr, "Codec not found\n");
//        exit(1);
//    }
//    
//    c = avcodec_alloc_context3(codec);
//    
//    /* put sample parameters */
//    c->bit_rate = 400000;
//    /* resolution must be a multiple of two */
//    c->width = outputSize.width;
//    c->height = outputSize.height;
//    /* frames per second */
//    c->time_base= (AVRational){1,25};
//    c->gop_size = 10; /* emit one intra frame every ten frames */
//    c->max_b_frames=1;
//    c->pix_fmt = PIX_FMT_YUV420P;
//    
//    if(codec_id == AV_CODEC_ID_H264)
//        av_opt_set(c->priv_data, "preset", "slow", 0);
//    
//    /* open it */
//    if (avcodec_open2(c, codec, NULL) < 0) {
//        fprintf(stderr, "Could not open codec\n");
//        exit(1);
//    }
//    
//    f = fopen(filename, "wb");
//    if (!f) {
//        fprintf(stderr, "Could not open %s\n", filename);
//        exit(1);
//    }
//    
//    frame = avcodec_alloc_frame();
//    if (!frame) {
//        fprintf(stderr, "Could not allocate video frame\n");
//        exit(1);
//    }
//    frame->format = c->pix_fmt;
//    frame->width  = inputSize.width;
//    frame->height = inputSize.height;
//    
//    scaledFrame = avcodec_alloc_frame();
//    if (!scaledFrame) {
//        fprintf(stderr, "Could not allocate video frame\n");
//        exit(1);
//    }
//    scaledFrame->format = c->pix_fmt;
//    scaledFrame->width  = c->width;
//    scaledFrame->height = c->height;
//    
//    /* create scaling context */
//    sws_ctx = sws_getContext(inputSize.width, inputSize.height, PIX_FMT_YUV420P, outputSize.width, outputSize.height, PIX_FMT_YUV420P, SWS_BILINEAR, NULL, NULL, NULL);
//    if (!sws_ctx) {
//        fprintf(stderr,
//                "Impossible to create scale context for the conversion "
//                "fmt:%s s:%dx%d -> fmt:%s s:%dx%d\n",
//                av_get_pix_fmt_name(PIX_FMT_YUV420P), inputSize.width, inputSize.height,
//                av_get_pix_fmt_name(PIX_FMT_YUV420P), outputSize.width, outputSize.height);
//        ret = AVERROR(EINVAL);
//        exit(1);
//    }
//    
//    /* the image can be allocated by any means and av_image_alloc() is
//     * just the most convenient way if av_malloc() is to be used */
//    ret = av_image_alloc(frame->data, frame->linesize, inputSize.width, inputSize.height,
//                         c->pix_fmt, 32);
//    if (ret < 0) {
//        fprintf(stderr, "Could not allocate raw picture buffer\n");
//        exit(1);
//    }
//    
//    /* the image can be allocated by any means and av_image_alloc() is
//     * just the most convenient way if av_malloc() is to be used */
//    ret = av_image_alloc(scaledFrame->data, scaledFrame->linesize, c->width, c->height,
//                         c->pix_fmt, 32);
//    if (ret < 0) {
//        fprintf(stderr, "Could not allocate raw picture buffer\n");
//        exit(1);
//    }
//    
//    //[super setupEncoderWithFormatDescription:newFormatDescription];
//    formatDescription = newFormatDescription;
//    CFRetain(formatDescription);
//    readyToEncode = YES;
//}

- (void) setupEncoderWithDimesions:(CMVideoDimensions)dimensions
{
    outputSize = dimensions;
    c = avcodec_alloc_context3(0); //NULL;
    frameNumber = 0;
    int codec_id = AV_CODEC_ID_H264; //  CODEC_ID_H264;
    
    
    
    /* register all the codecs */
    avcodec_register_all();
            
    /* find the video encoder */
    //codec = avcodec_find_encoder(codec_id);
    codec = avcodec_find_encoder_by_name("libx264");
    if (!codec) {
        NSLog( @"Codec not found");
        exit(1);
    }
    
    c = avcodec_alloc_context3(codec);
    
//    /* put sample parameters */
//    c->bit_rate = 400000;
//    /* resolution must be a multiple of two */
//    c->width = 352;
//    c->height = 288;
//    /* frames per second */
//    c->time_base= (AVRational){1,25};
//    c->gop_size = 10; /* emit one intra frame every ten frames */
//    c->max_b_frames=1;
//    c->pix_fmt = AV_PIX_FMT_YUV420P;
    
    
    c->width = outputSize.width;
    c->height = outputSize.height;
    //c->bit_rate = c->width * c->height * 4;
    c->max_b_frames=0;
    c->pix_fmt = PIX_FMT_YUV420P;
    //c->time_base= (AVRational){1,15};
    
    // not sure about these
    c->refs = 1; //ref = 1
    c->keyint_min = 15*2;   //keyint=<framerate * segment length>
    c->level = 30;  // level=30
    c->bit_rate_tolerance = 1; //ratetol=1.0
    
    
    
    
//    int64_t crf = -1;
    
    if(codec_id == AV_CODEC_ID_H264) {
        av_opt_set(c->priv_data, "preset", "ultrafast", AV_OPT_SEARCH_CHILDREN);
        av_opt_set(c->priv_data, "tune", "zerolatency", AV_OPT_SEARCH_CHILDREN);
        
        // not sure about these
        //av_opt_set(c->priv_data, "crf", "20", AV_OPT_SEARCH_CHILDREN);
        av_opt_set_int(c->priv_data, "crf", 20, AV_OPT_SEARCH_CHILDREN);    //crf=20
        av_opt_set(c->priv_data, "hrd", "crf", AV_OPT_SEARCH_CHILDREN);     //hrd=cbr
        av_opt_set_int(c->priv_data, "lookahead", 0, AV_OPT_SEARCH_CHILDREN);    //lookahead=0
    }
    
    
    /* open it */
    if (avcodec_open2(c, codec, NULL) < 0) {
        NSLog(@"Could not open codec");
        exit(1);
    }
    
    //av_opt_get_int(c->priv_data, "crf", 1, &crf);
    char preset[50];
    uint_fast8_t **preset_ptr = (uint_fast8_t**)&preset;
    
    av_opt_get(c->priv_data, "preset", 1, preset_ptr);
    NSLog(@"preset %s", (char *)*preset_ptr);
    
    frame = avcodec_alloc_frame();
    if (!frame) {
        NSLog(@"Could not allocate video frame");
        exit(1);
    }
    frame->format = c->pix_fmt;
    frame->width  = c->width;
    frame->height = c->height;
    
    /* the image can be allocated by any means and av_image_alloc() is
     * just the most convenient way if av_malloc() is to be used */
    ret = av_image_alloc(frame->data, frame->linesize, c->width, c->height,
                         c->pix_fmt, 1);
    
    if (ret < 0) {
        NSLog(@"Could not allocate raw picture buffer");
        exit(1);
    }
    
    readyToEncode = YES;
}

- (void) scaleVideoToOutputSize {
    
}

- (void) finishEncoding {
//    uint8_t endcode[] = { 0, 0, 1, 0xb7 };
//    
//    for (got_output = 1; got_output; frameNumber++) {
//        fflush(stdout);
//        
//        ret = avcodec_encode_video2(c, &pkt, NULL, &got_output);
//        if (ret < 0) {
//            fprintf(stderr, "Error encoding frame\n");
//            exit(1);
//        }
//        
//        if (got_output) {
//            //printf("Write frame %3d (size=%5d)\n", frameNumber, pkt.size);
//            fwrite(pkt.data, 1, pkt.size, f);
//            av_free_packet(&pkt);
//        }
//    }
//    
//    /* add sequence end code to have a real mpeg file */
//    fwrite(endcode, 1, sizeof(endcode), f);
//    fclose(f);
    
    avcodec_close(c);
    av_free(c);
    av_freep(&frame->data[0]);
//    av_freep(&scaledFrame->data[0]);
    avcodec_free_frame(&frame);
//    avcodec_free_frame(&scaledFrame);
//    sws_freeContext(sws_ctx);
    printf("\n");
    //[super finishEncoding];
    CFRelease(formatDescription);
    formatDescription = NULL;
    readyToEncode = NO;
}

- (int) encodeSampleBuffer:(CMSampleBufferRef)sampleBuffer
                         output:(AVPacket *)pkt
{
    CVImageBufferRef pixelBuffer = CMSampleBufferGetImageBuffer(sampleBuffer);
    CVPixelBufferLockBaseAddress( pixelBuffer, 0 );
    
	int bufferWidth = 0;
	int bufferHeight = 0;
	uint8_t *pixel = NULL;
        
    if (CVPixelBufferIsPlanar(pixelBuffer)) {
        //int planeCount = CVPixelBufferGetPlaneCount(pixelBuffer);
        int basePlane = 0;
        pixel = (uint8_t *)CVPixelBufferGetBaseAddressOfPlane(pixelBuffer, basePlane);
        bufferHeight = CVPixelBufferGetHeightOfPlane(pixelBuffer, basePlane);
        bufferWidth = CVPixelBufferGetWidthOfPlane(pixelBuffer, basePlane);
    } else {
        pixel = (uint8_t *)CVPixelBufferGetBaseAddress(pixelBuffer);
        bufferWidth = CVPixelBufferGetWidth(pixelBuffer);
        bufferHeight = CVPixelBufferGetHeight(pixelBuffer);
    }
    
    av_init_packet(pkt);
    pkt->data = NULL;    // packet data will be allocated by the encoder
    pkt->size = 0;
    
    //unsigned char y_pixel = pixel[0];
    
    
    fflush(stdout);
    for (int y = 0; y < bufferHeight; y++) {
        for (int x = 0; x < bufferWidth; x++) {
            frame->data[0][y * frame->linesize[0] + x] = pixel[0];
            pixel++;
        }
    }
    
    /* Cb and Cr */
    
    for (int y = 0; y < bufferHeight / 2; y++) {
        for (int x = 0; x < bufferWidth / 2; x++) {
            frame->data[1][y * frame->linesize[1] + x] = pixel[0];
            frame->data[2][y * frame->linesize[2] + x] = pixel[1];
            pixel+=2;
        }
    }
    
    /* convert to destination format */
//    sws_scale(sws_ctx, (const uint8_t * const*)frame->data,
//              frame->linesize, 0, inputSize.height, scaledFrame->data, scaledFrame->linesize);
    
    frame->pts = frameNumber;
    
    /* encode the image */
    ret = avcodec_encode_video2(c, pkt, frame, &got_output);
    if (ret < 0) {
        fprintf(stderr, "Error encoding frame\n");
        exit(1);
    }
    
//    if (got_output) {
//        NSLog(@"encoded frame");
//    } else {
//        NSLog(@"empty output");
//    }
    
    frameNumber++;
    CVPixelBufferUnlockBaseAddress( pixelBuffer, 0 );
    
    return got_output;
}

-(void) freePacket:(AVPacket *)pkt
{
    av_free_packet(pkt);
}



@end
