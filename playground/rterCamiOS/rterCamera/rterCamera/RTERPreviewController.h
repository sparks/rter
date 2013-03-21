//
//  previewController.h
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-06.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import <UIKit/UIKit.h>
#import <CoreMedia/CoreMedia.h>
#import <AVFoundation/AVFoundation.h>
#import <dispatch/dispatch.h>

#import <QuartzCore/QuartzCore.h>

@protocol RTERPreviewControllerDelegate <NSObject>

@required
- (void)back;

@end

@interface RTERPreviewController : UIViewController<AVCaptureAudioDataOutputSampleBufferDelegate>
{
    // dispatch queue for encoding and POSTing
    dispatch_queue_t postQueue;
    dispatch_queue_t encoderQueue;
    NSOperationQueue *postOpQueue;
}

@property (nonatomic, retain) NSObject<RTERPreviewControllerDelegate> *delegate;

@property (strong, nonatomic) IBOutlet UIView *previewView;

@property (strong, nonatomic) IBOutlet UIToolbar *toobar;

- (IBAction)clickedStart:(id)sender;

- (IBAction)clickedBack:(id)sender;

- (void) captureOutput:(AVCaptureOutput *)captureOutput didOutputSampleBuffer:(CMSampleBufferRef)sampleBuffer fromConnection:(AVCaptureConnection *)connection;


@end
