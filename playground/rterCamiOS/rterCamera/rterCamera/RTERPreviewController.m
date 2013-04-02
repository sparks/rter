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
	
	NSURLConnection *streamingAuthConnection;
	NSString *authString; //If authString doesn't have to be private it should be a property CB
    
    GLKView* _glkView;
    RTERGLKViewController* _glkVC;
}

@end

@implementation RTERPreviewController

@synthesize toobar;
@synthesize previewView;
@synthesize streamingToken;
@synthesize streamingEndpoint;
@synthesize glkView = _glkView;
@synthesize itemID;


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
        
        postOpQueue = [[NSOperationQueue alloc] init];
        
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    // Do any additional setup after loading the view from its nib.
    streamingToken = @"";
	
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
    
    //Create OpenGL Layer
    
    //create eaglcontext
    EAGLContext *context = [[EAGLContext alloc]initWithAPI:kEAGLRenderingAPIOpenGLES1];
    
    //assign context to synthesized GLKView
    _glkView.context = context;
    [self.glkView setNeedsDisplay];
    
    //initialize View Controller for the GLKView
    _glkVC = [[RTERGLKViewController alloc]initWithNibName:nil bundle:nil view:_glkView previewController:self];

    
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
    
	// get rid of tokens and authentication data for streaming
	streamingAuthConnection = nil;
	streamingEndpoint = nil;
	streamingToken = nil;
	
	[[self delegate] back];
	
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (IBAction)clickedStart:(id)sender {
    if(!sendingData) {
		
        sendingData = YES;
        
		// get token for video streaming
		[self getStreamingToken];
		
        // start recording
//        
//        [self startRecording];
		
        
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


-(void) getStreamingToken {
	NSLog(@"Attempting to get Streaming token:");
	
	// the json string to post
	NSString *jsonString = [NSString stringWithFormat:@"{\"Type\":\"streaming-video-v1\",\"StartTime\":\"0001-01-01T00:00:00Z\"}"];
	NSData *postData = [jsonString dataUsingEncoding:NSUTF8StringEncoding];
	
	// setup the request
	NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:[NSURL URLWithString:@"http://rter.cim.mcgill.ca:80/1.0/items"]];
	
	//NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:[NSURL URLWithString:@"http://142.157.58.36:8080/1.0/items"]];
	
	[request setHTTPMethod:@"POST"];
	[request setHTTPShouldHandleCookies:YES];
	[request setHTTPBody:postData];
	[request setAllowsCellularAccess:YES];
	[request setValue:@"application/json" forHTTPHeaderField:@"Content-Type"];
	[request setValue:[NSString stringWithFormat:@"%d",[postData length]] forHTTPHeaderField:@"Content-Length"];
	[request setValue:[[self delegate] cookieString] forHTTPHeaderField:@"Set-Cookie"];
	
	streamingAuthConnection = [NSURLConnection connectionWithRequest:request delegate:[self delegate]];
    //streamingAuthConnection = [[NSURLConnection alloc]initWithRequest:request delegate:self startImmediately:YES];
}



-(NSURLConnection*)getAuthConnection{
	return streamingAuthConnection;
}



-(void)setAuthString:(NSString*)newAuth {
	authString = newAuth;
}

-(NSString *)getAuthString {
    return authString;
}

/* process the frames here */

-(void) captureOutput:(AVCaptureOutput *)captureOutput didOutputSampleBuffer:(CMSampleBufferRef)sampleBuffer fromConnection:(AVCaptureConnection *)connection
{
        
    AVPacket pkt;   // encoder output
    if([encoder encodeSampleBuffer:sampleBuffer output:&pkt]) {
        NSLog(@"encoded frame");
        
        NSMutableURLRequest *postRequest = [NSMutableURLRequest requestWithURL:[NSURL URLWithString:[NSString stringWithFormat:@"%@/avc", streamingEndpoint]]];
		
		//NSMutableURLRequest *postRequest = [NSMutableURLRequest requestWithURL:[NSURL URLWithString:@"http://142.157.46.36:1234"]];
		
        [postRequest setHTTPMethod:@"POST"];
        [postRequest setHTTPBody:[NSData dataWithBytes:pkt.data length:pkt.size]];
		[postRequest setValue:[[self delegate] cookieString] forHTTPHeaderField:@"Set-Cookie"];
		[postRequest setValue:authString forHTTPHeaderField:@"Authorization"];
        
        dispatch_async(postQueue, ^{
            NSHTTPURLResponse *response;
            NSError *err;
            NSData *responseData = [NSURLConnection sendSynchronousRequest:postRequest returningResponse:&response error:&err];
            //        if ([response respondsToSelector:@selector(allHeaderFields)]) {
            NSDictionary *dictionary = [response allHeaderFields];
            NSLog( @"%@", [dictionary description]);
        });
		
//        [NSURLConnection sendAsynchronousRequest:postRequest
//                                           queue:postOpQueue
//                               completionHandler:^(NSURLResponse *response, NSData *data, NSError *error)
//        {
//            
//            NSDictionary *dictionary = [(NSHTTPURLResponse *)response allHeaderFields];
//            NSLog(@"%d - %@\n%@", [(NSHTTPURLResponse *)response statusCode], [NSHTTPURLResponse localizedStringForStatusCode:[(NSHTTPURLResponse *)response statusCode]], [dictionary description]);
//        }];
        
        [encoder freePacket:&pkt];
    
    }
}

- (void)willRotateToInterfaceOrientation:(UIInterfaceOrientation)toInterfaceOrientation duration:(NSTimeInterval)duration {
    [_glkVC interfaceOrientationDidChange:toInterfaceOrientation];
}

- (IBAction)clickedBack:(id)sender {
    [self onExit];
}
@end
