//
//  previewController.m
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-06.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import "RTERPreviewController.h"
#import <math.h>
#import "RTERVideoEncoder.h"
#import "RTERGLKViewController.h"

#import "RTERArrow.h"

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
    
    // encoder
    RTERVideoEncoder *encoder;
    
    int x;
}

@end

@implementation RTERPreviewController

@synthesize toobar;
@synthesize previewView;
@synthesize glkView = _glkView;

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
        
        // init dispatch queues
        encoderQueue = dispatch_queue_create("com.rterCamera.encoderQueue", DISPATCH_QUEUE_SERIAL);
        postQueue = dispatch_queue_create("com.rterCamera.postQueue", DISPATCH_QUEUE_SERIAL);
        
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    // Do any additional setup after loading the view from its nib.
    // capture session
    captureSession = [[AVCaptureSession alloc] init];
    
    // encoder
    encoder = [[RTERVideoEncoder alloc] init];

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
    if ([captureSession canSetSessionPreset:AVCaptureSessionPreset352x288]) {
        captureSession.sessionPreset = AVCaptureSessionPreset352x288;
        NSLog(@"352x288");
        
        CMVideoDimensions dimensions;
        dimensions.width = 352;
        dimensions.height = 288;
                
        [encoder setupEncoderWithDimesions:dimensions];
    }
//    if ([captureSession canSetSessionPreset:AVCaptureSessionPreset640x480]) {
//        captureSession.sessionPreset = AVCaptureSessionPreset640x480;
//        NSLog(@"640x480");
//        
//        CMVideoDimensions dimensions;
//        dimensions.width = 640;
//        dimensions.height = 480;
//        
//        [encoder setupEncoderWithDimesions:dimensions];
//    }
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
    /* possible ones to ues for h.264:
     * kCVPixelFormatType_420YpCbCr8BiPlanarVideoRange
     * kCVPixelFormatType_420YpCbCr8BiPlanarFullRange
     */
    outputDevice.videoSettings = [NSDictionary dictionaryWithObjectsAndKeys:[NSNumber numberWithUnsignedInt:kCVPixelFormatType_420YpCbCr8BiPlanarVideoRange], (id)kCVPixelBufferPixelFormatTypeKey,
                                 nil];
    // set self as the delegate for the output for now
    [outputDevice setSampleBufferDelegate:self queue:encoderQueue];
    
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
    
    


    EAGLContext *context = [[EAGLContext alloc]initWithAPI:kEAGLRenderingAPIOpenGLES1];
    //_glkView = [[GLKView alloc]initWithFrame:screenBounds context:context];
    _glkView.context = context;
    //_glkView.delegate = _glkVC;
    [self.glkView setNeedsDisplay];

       
   
    
    _glkVC = [[RTERGLKViewController alloc]initWithNibName:nil bundle:nil view:_glkView];


    NSTimer *nst_Timer = [NSTimer scheduledTimerWithTimeInterval:4 target:self selector:@selector(showTime) userInfo:nil repeats:YES];
    [nst_Timer fire];
    x = 1;
        

    
}

-(void)showTime {
    if (x == 0) {
        x = 1;
        [_glkVC indicateTurnToDirection:RIGHT withPercentage:1];
        
    }else if(x==1){
        x = 2;
        [_glkVC indicateTurnToDirection:NONE withPercentage:1];
    }else if(x==2){
        x = 3;
        [_glkVC indicateTurnToDirection:LEFT withPercentage:1];
    }else{
        x = 0;
        [_glkVC indicateTurnToDirection:FREE withPercentage:1];
    }
}

#pragma mark - GLKViewDelegate

/*- (void)glkView:(GLKView *)view drawInRect:(CGRect)rect {
    
    glClearColor(_curRed, 0.0, 0.0, 1.0);
    glClear(GL_COLOR_BUFFER_BIT);
    
}
*/
#pragma mark - GLKViewControllerDelegate
/*
- (void)glkViewControllerUpdate:(GLKViewController *)controller {
    if (_increasing) {
        _curRed += 1.0 * controller.timeSinceLastUpdate;
    } else {
        _curRed -= 1.0 * controller.timeSinceLastUpdate;
    }
    if (_curRed >= 1.0) {
        _curRed = 1.0;
        _increasing = NO;
    }
    if (_curRed <= 0.0) {
        _curRed = 0.0;
        _increasing = YES;
    }
    
    //[_glkView display];
}*/

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

/* process the frames here */

-(void) captureOutput:(AVCaptureOutput *)captureOutput didOutputSampleBuffer:(CMSampleBufferRef)sampleBuffer fromConnection:(AVCaptureConnection *)connection
{
    CVImageBufferRef imageBuffer = CMSampleBufferGetImageBuffer( sampleBuffer );
    CGSize imageSize = CVImageBufferGetEncodedSize( imageBuffer );    
    //NSLog( @"frame captured at %.fx%.f", imageSize.width, imageSize.height );
    
    //CVPixelBufferLockBaseAddress(imageBuffer,0); // lock buffe
        
        
    AVPacket pkt;   // encoder output
    if([encoder encodeSampleBuffer:sampleBuffer output:&pkt]) {
        NSLog(@"encoded frame");
    
        dispatch_async(postQueue, ^{
            NSMutableURLRequest *postRequest = [NSMutableURLRequest requestWithURL:[NSURL URLWithString:@"http://142.157.34.160:8080/v1/ingest/0/avc"]];
            [postRequest setHTTPMethod:@"POST"];
            [postRequest setHTTPBody:[NSData dataWithBytes:pkt.data length:pkt.size]];
            NSHTTPURLResponse *response;
            NSError *err;
            NSData *responseData = [NSURLConnection sendSynchronousRequest:postRequest returningResponse:&response error:&err];
            //        if ([response respondsToSelector:@selector(allHeaderFields)]) {
            NSDictionary *dictionary = [response allHeaderFields];
            NSLog([dictionary description]);
            //        }
            [encoder freePacket:&pkt];
            //        NSLog(@"finished sending frame");
        });
    }
}
- (void)willRotateToInterfaceOrientation:(UIInterfaceOrientation)toInterfaceOrientation duration:(NSTimeInterval)duration {

    [_glkVC interfaceOrientationDidChange:toInterfaceOrientation];
}

- (IBAction)clickedBack:(id)sender {
    [self onExit];
}
@end
