//
//  previewController.m
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-06.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import "RTERPreviewController.h"
#import <math.h>

#define DESIRED_FPS 15

@interface RTERPreviewController ()
{
    AVCaptureSession *captureSession;
    AVCaptureVideoPreviewLayer *previewLayer;
    AVCaptureVideoDataOutput *outputDevice;
    
    BOOL sendingData;
    
    // save default frame rate
    CMTime defaultMaxFrameDuration;
    CMTime defaultMinFrameDuration;
    
    // desired frame rate
    CMTime desiredFrameDuration;
}

@end

@implementation RTERPreviewController

@synthesize toobar;
@synthesize previewView;

- (id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil
{
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self) {
        // init stuff
        sendingData = NO;
        
        // desired FPS
        desiredFrameDuration = CMTimeMake(1, DESIRED_FPS);
        
        // listen for notifications
        [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(appWillResignActive) name:UIApplicationWillResignActiveNotification object:nil];
        [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(appDidBecomeActive) name:UIApplicationDidBecomeActiveNotification object:nil];
        
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    // Do any additional setup after loading the view from its nib.
    // capture session
    captureSession = [[AVCaptureSession alloc] init];
    
    // video session settings
    
    /* possible resolution settings:
     AVCaptureSessionPresetPhoto;
     AVCaptureSessionPresetHigh;
     AVCaptureSessionPresetMedium;
     AVCaptureSessionPresetLow;
     AVCaptureSessionPreset320x240;
     AVCaptureSessionPreset352x288;
     AVCaptureSessionPreset640x480;
     AVCaptureSessionPreset960x540;
     AVCaptureSessionPreset1280x720;
     */
    if ([captureSession canSetSessionPreset:AVCaptureSessionPreset640x480]) {
        captureSession.sessionPreset = AVCaptureSessionPreset640x480;
        NSLog(@"640x480");
    }
    else {
        // Handle the failure.
    }
    
    previewLayer = [AVCaptureVideoPreviewLayer layerWithSession:captureSession];
    
    
    AVCaptureDevice *videoDevice = [AVCaptureDevice defaultDeviceWithMediaType:AVMediaTypeVideo];
    if (videoDevice) {
        NSError *error;
        AVCaptureDeviceInput *videoIn = [AVCaptureDeviceInput deviceInputWithDevice:videoDevice error:&error];
        if (!error) {
            if ([captureSession canAddInput:videoIn])
                [captureSession addInput:videoIn];
            else {
                NSLog(@"Couldn't add video input");
                [self onExit];
            }
        } else {
            NSLog(@"Couldn't create video input");
            [self onExit];
        }
    } else {
        NSLog(@"Couldn't create video capture device");
        [self onExit];
    }
    
    //init output
    outputDevice = [[AVCaptureVideoDataOutput alloc] init];
    
    // set pixel buffer format
    outputDevice.videoSettings = [NSDictionary dictionaryWithObjectsAndKeys:[NSNumber numberWithUnsignedInt:kCVPixelFormatType_420YpCbCr8BiPlanarFullRange], (id)kCVPixelBufferPixelFormatTypeKey,
                                 nil];
    // set self as the delegate for the output for now
    [outputDevice setSampleBufferDelegate:self queue:dispatch_get_main_queue()];
    
    // add preview layer to preview view
    [previewView.layer addSublayer:previewLayer];
    
    // set the location and size of teh preview layer to that of the preview view
    [previewLayer setFrame:previewView.bounds];

    // resize preview to fit within the view, but retain its original aspect ration
    [previewLayer setVideoGravity:AVLayerVideoGravityResizeAspectFill];
    
    // make sure the preview stays within the bounds
    // (otherwise it will take up the whole screen)
    previewView.clipsToBounds = YES;
    
    // get the default FPS
    defaultMaxFrameDuration = previewLayer.connection.videoMaxFrameDuration;
    defaultMinFrameDuration = previewLayer.connection.videoMinFrameDuration;
    
    // start the capture session so that the preview shows up
    [captureSession startRunning];
    
}

- (void)willAnimateRotationToInterfaceOrientation:(UIInterfaceOrientation)toInterfaceOrientation duration:(NSTimeInterval)duration
{
    // this happens in the middle of the orientation animation
    // the bounds of all the auto rotated views have already been set
    
    // rotate the video
    switch (toInterfaceOrientation) {
        case UIInterfaceOrientationLandscapeLeft:
            [[previewLayer connection] setVideoOrientation:AVCaptureVideoOrientationLandscapeLeft];
            break;
        case UIInterfaceOrientationLandscapeRight:
            [[previewLayer connection] setVideoOrientation:AVCaptureVideoOrientationLandscapeRight];
            break;
        case UIInterfaceOrientationPortraitUpsideDown:
            // not supporting this orientation
            break;
        default:
            [[previewLayer connection] setVideoOrientation:AVCaptureVideoOrientationPortrait];
            break;
    }

    // the bounds have changed
    [previewLayer setFrame: [previewView bounds]];
}

- (void)appWillResignActive {
    if (captureSession && [captureSession isRunning]) {
        [captureSession stopRunning];
    }
}

- (void)appDidBecomeActive {
    if (captureSession) {
        [captureSession startRunning];
    }
}

- (void)onExit {
    // stop listening for notifications
    [[NSNotificationCenter defaultCenter] removeObserver:self];
    
    if (captureSession && [captureSession isRunning]) {
        if(sendingData) {
            [captureSession removeOutput:outputDevice];
        }
        [captureSession stopRunning];
    }
    [[self delegate] back];
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (IBAction)clickedStart:(id)sender {
    if(!sendingData) {
        // start recording
        sendingData = YES;
        [self startRecording];
        
        [(UIBarButtonItem *) sender setTitle:@"stop"];
    } else {
        // stop recording
        sendingData = NO;
        [self stopRecording];
        
        [(UIBarButtonItem *) sender setTitle:@"start"];
    }
    
}

- (void) startRecording {
    
    [captureSession addOutput:outputDevice];
    
    /* set the frame rate
     * for some reason have to set both the max and the min for it to work properly */
        
    AVCaptureConnection *conn = [previewLayer connection]; //[outputDevice connectionWithMediaType:AVMediaTypeVideo];
    
    CMTimeShow(conn.videoMinFrameDuration);
    CMTimeShow(conn.videoMaxFrameDuration);
    
    if (conn.isVideoMinFrameDurationSupported)
        conn.videoMinFrameDuration = desiredFrameDuration;
    if (conn.isVideoMaxFrameDurationSupported)
        conn.videoMaxFrameDuration = desiredFrameDuration;
    
    CMTimeShow(conn.videoMinFrameDuration);
    CMTimeShow(conn.videoMaxFrameDuration);
    
//    for (NSString *codec in [outputDevice availableVideoCodecTypes]) {
//        NSLog(@"%@", codec);
//    }
}

- (void) stopRecording {
    [captureSession removeOutput:outputDevice];
    
    /* restore to default frame rate when not "recording"
     * for some reason have to set both the max and the min for it to work properly */
    AVCaptureConnection *conn = [previewLayer connection]; //[outputDevice connectionWithMediaType:AVMediaTypeVideo];
    
    CMTimeShow(conn.videoMinFrameDuration);
    CMTimeShow(conn.videoMaxFrameDuration);
    
    if (conn.isVideoMinFrameDurationSupported)
        conn.videoMinFrameDuration = defaultMinFrameDuration;
    if (conn.isVideoMaxFrameDurationSupported)
        conn.videoMaxFrameDuration = defaultMaxFrameDuration;
    
    CMTimeShow(conn.videoMinFrameDuration);
    CMTimeShow(conn.videoMaxFrameDuration);
}

-(void) captureOutput:(AVCaptureOutput *)captureOutput didOutputSampleBuffer:(CMSampleBufferRef)sampleBuffer fromConnection:(AVCaptureConnection *)connection
{
    CVImageBufferRef imageBuffer = CMSampleBufferGetImageBuffer( sampleBuffer );
    CGSize imageSize = CVImageBufferGetEncodedSize( imageBuffer );
    // also in the 'mediaSpecific' dict of the sampleBuffer
    
    NSLog( @"frame captured at %.fx%.f", imageSize.width, imageSize.height );
    
}

- (IBAction)clickedBack:(id)sender {
    [self onExit];
}
@end
