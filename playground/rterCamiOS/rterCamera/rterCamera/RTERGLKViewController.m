//
//  RTERGLKViewController.m
//  rterCamera
//
//  Created by Cameron Bell on 13-03-18.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import "RTERGLKViewController.h"



#define rad(X) X*180/M_PI

@interface RTERGLKViewController () {
    NSNumber *lock;
    RTERPreviewController *previewController;
    
    NSOperationQueue *getQueue;
    NSOperationQueue *putQueue;
    
    float latitude;
    float longitude;
}
@end

@implementation RTERGLKViewController

- (id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil view:(GLKView *)view previewController:(RTERPreviewController *)prev
{
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self) {
        // Custom initialization
        
        //reference to previewController to get authString - not sure if authString is dynamic
        previewController = prev;
        
        getQueue = [[NSOperationQueue alloc]init];
        putQueue = [[NSOperationQueue alloc]init];
        
        _curRed = 0.0;
        _increasing = YES;

        _glkView = view;
        _glkView.delegate = self;
        [_glkView setNeedsDisplay];
        self.preferredFramesPerSecond = 60;

        self.view = _glkView;
        self.view.opaque = NO;
        self.view.backgroundColor = [UIColor clearColor];

        lock = [NSNumber numberWithBool:NO];
        
        arrowLeft = [[RTERArrow alloc]initArrow];
        arrowRight = [[RTERArrow alloc]initArrow];
        indicatorFrame = [[RTERIndicatorFrame alloc]initIndicatorFrame];
        
        freeRoam = NO;
        rightSideUp = TRUE;
        orientationTolerance = 30;
        currentOrientation = 0;
        desiredOrientation = 0;
        
        arrowScale = 1.0f;
        arrowScaleMax = 1.2f;
        arrowScaleMin = 0.2f;

        // pulsating variables
        arrowPulsateScale = 1.0f;
        arrowPulsateSpeed = 0.1f;
        arrowPulsateSpeedMin = 0.01f;
        arrowPulsateMax = 1.2f;
        arrowPulsateMin = 0.9f;
        arrowPulsateIncrease = true;

        displayLeft = false;
        displayRight = false;
        
        [self indicateTurnToDirection:RIGHT withPercentage:10.0];

        [self initializeGLView:_glkView];
        
        /*float width = _glkView.frame.size.width;
        float height = _glkView.frame.size.height;
        
        if (height == 0) {
            height = 1;
        }
        aspect = width/height;
        
        distance = 6.0f;
        xTotal = (float) aspect * tanf(rad(45.0/2));
        yTotal = (float) (tanf(rad(45.0/2))*distance *2);
        */
        
        distance = -0.0f;
        
        [self setUpLocationManager];
        
        debugScreen = [[UITextView alloc]initWithFrame:CGRectMake(50, 50, 250, 50)];
        [debugScreen setText:@"Debug Screen"];
        [debugScreen setTextColor:[UIColor redColor]];
        [debugScreen setBackgroundColor:[UIColor clearColor]];
        
        [self.view addSubview:debugScreen];
        
        getHeadingTimer = [NSTimer scheduledTimerWithTimeInterval:4 target:self selector:@selector(getHeadingPutCoordindates) userInfo:nil repeats:YES];
        [getHeadingTimer fire];
        
    }
    return self;
}


-(void)getHeadingPutCoordindates {
    
    
    NSURL *url = [NSURL URLWithString:@"http://rter.cim.mcgill.ca/1.0/users/sparky/direction"];
    NSMutableURLRequest *urlRequest = [[NSMutableURLRequest alloc]initWithURL:url];
    //[urlRequest setHTTPMethod:@"GET"];
    //[urlRequest setValue:@"application/json" forHTTPHeaderField:@"Content-Type"];
    if (!currentGetConnection) {
        currentGetConnection = [[NSURLConnection alloc]initWithRequest:urlRequest delegate:self startImmediately:YES];
        headingData = [[NSMutableData alloc]init];
        [currentGetConnection start];
        NSLog(@"GetConnectionDidStart");
    }
    //put the id of the item in here
   /* NSURL *putUrl = [NSURL URLWithString:@"http://rter.cim.mcgill.ca:8080/1.0/items/<id of item>"];
    //NSMutableURLRequest *putUrlRequest = [[NSMutableURLRequest alloc]initWithURL:url];
    [urlRequest setHTTPMethod:@"PUT"];
    [urlRequest setValue:@"application/json" forHTTPHeaderField:@"Content-Type"];
    if (!currentPutConnection) {
        currentPutConnection = [[NSURLConnection alloc]initWithRequest:putUrlRequest delegate:self startImmediately:YES];
        //headingData = [[NSMutableData alloc]init];
        [currentPutConnection start];
        NSLog(@"PutConnectionDidStart");
    }*/
    
    
    ////PUT REQUEST
    NSMutableURLRequest *putRequest = [NSMutableURLRequest requestWithURL:[NSURL URLWithString:[NSString stringWithFormat:@"http://rter.cim.mcgill.ca:8080/1.0/items/%@",[previewController itemID]]]];
    
    
    // the json string to post
	NSString *jsonString = [NSString stringWithFormat:@"{\"Lat\":\"%f\",\"Lng\":\"%f\",\"Heading\":\"%f\"}",latitude,longitude,currentOrientation];
	NSData *postData = [jsonString dataUsingEncoding:NSUTF8StringEncoding];

    
    [putRequest setHTTPMethod:@"PUT"];
    [putRequest setHTTPBody:postData];
    [putRequest setValue:[[previewController delegate] cookieString] forHTTPHeaderField:@"Set-Cookie"];

    
    

    //[putRequest setValue:[previewController getAuthString] forHTTPHeaderField:@"Authorization"];
    
    [NSURLConnection sendAsynchronousRequest:putRequest
                                       queue:putQueue
                           completionHandler:^(NSURLResponse *response, NSData *data, NSError *error)
     {
         
         NSDictionary *dictionary = [(NSHTTPURLResponse *)response allHeaderFields];
         NSLog(@"%d - %@\n%@", [(NSHTTPURLResponse *)response statusCode], [NSHTTPURLResponse localizedStringForStatusCode:[(NSHTTPURLResponse *)response statusCode]], [dictionary description]);
     }];
    


    
    
    
}

//////NSURLCONNECTION DELEGATE METHODS///////////////////
/*
 this method might be calling more than one times according to incoming data size
 */
-(void)connection:(NSURLConnection *)connection didReceiveData:(NSData *)data{
    
    if ([connection isEqual:currentGetConnection]) {
            [headingData appendData:data];
    }

    
}

/*
 if there is an error occured, this method will be called by connection
 */
-(void)connection:(NSURLConnection *)connection didFailWithError:(NSError *)error{
  
    NSLog(@"CONNECTION_FAILURE");
    
    if (error.code == -1009) {
        NSLog(@"NO INTERNET CONNECTION");
    }else{
        NSLog(@"Unidentified Error");
    }
    currentGetConnection = nil;

    
    
}

/*
 if data is successfully received, this method will be called by connection
 */
-(void)connectionDidFinishLoading:(NSURLConnection *)connection{
    
    if ([connection isEqual:currentGetConnection]) {
        //parse the json here.
        NSLog(@"ConnectionDidFinishLoading");
       // NSString *jsonString = [[NSString alloc] initWithData:headingData encoding:NSUTF8StringEncoding];
        
       // NSLog(@"===JSON_2_B_PARSED===\n%@",jsonString);
        
		NSError *error;
		NSDictionary *jsonDict = [NSJSONSerialization JSONObjectWithData:headingData options:
								  NSJSONReadingMutableContainers error:&error];
       
        NSString *desiredOrientationString = [jsonDict objectForKey:@"Heading"];
        NSLog(@"%@",desiredOrientationString);
        desiredOrientation = [desiredOrientationString floatValue];
        
        
    }
    
  
    
    currentGetConnection = nil;
    
}


//UNTESTED
-(void)initializeGLView:(GLKView *)view {
    glClearColor(0.0f, 0.0f, 0.0f, 0.0f); // Set color's clear-value to
    // black
    glClearDepthf(1.0f); // Set depth's clear-value to farthest
    glEnable(GL_DEPTH_TEST); // Enables depth-buffer for hidden
    // surface removal
    glDepthFunc(GL_LEQUAL); // The type of depth testing to do
    glHint(GL_PERSPECTIVE_CORRECTION_HINT,GL_NICEST);   // nice
                                                        // perspective
                                                        // view
    glShadeModel(GL_SMOOTH); // Enable smooth shading of color
    glDisable(GL_DITHER); // Disable dithering for better performance
    
}

-(void)setDesiredOrientation:(float)dO {
    desiredOrientation = dO;
}

-(void)indicateTurnToDirection:(Indicate) direction withPercentage:(float) percentage {
    @synchronized(lock)
    {
        arrowScale = (arrowScaleMax - arrowScaleMin) * percentage +arrowScaleMin;
        
        switch (direction) {
            case LEFT:{
                displayLeft = true;
				displayRight = false;
				[indicatorFrame setColour:RED];
                break;}
            case RIGHT:{
                displayLeft = false;
				displayRight = true;
                [indicatorFrame setColour:RED];
                break;}
            case NONE:{
                displayLeft = false;
				displayRight = false;
                [indicatorFrame setColour:GREEN];
                break;}
            case FREE:{
                displayLeft = false;
				displayRight = false;
                [indicatorFrame setColour:BLUE];
                break;}
            default:
                break;
        }
        
        
    }
}

-(void)update {
    // pulsate arrows
    if (arrowPulsateIncrease) {
        float speed = (arrowPulsateMax - arrowPulsateScale) * arrowPulsateSpeed;
        if (speed < arrowPulsateSpeedMin)
            speed = arrowPulsateSpeedMin;
        arrowPulsateScale += speed;
        if (arrowPulsateScale >= arrowPulsateMax) {
            arrowPulsateScale = arrowPulsateMax;
            arrowPulsateIncrease = false;
        }
    } else {
        float speed = (arrowPulsateScale - arrowPulsateMin) * arrowPulsateSpeed;
        if (speed < arrowPulsateSpeedMin)
            speed = arrowPulsateSpeedMin;
        arrowPulsateScale -= speed;
        if (arrowPulsateScale <= arrowPulsateMin) {
            arrowPulsateScale = arrowPulsateMin;
            arrowPulsateIncrease = true;
        }
    }
}



- (void)glkView:(GLKView *)view drawInRect:(CGRect)rect {
    // Clear color and depth buffers 
    glClearColor(0.0, 0.0, 0.0, 0.0);
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    
    //[arrowLeft drawInView:view];
   float arrowScale_tmp = arrowPulsateScale * arrowScale;
    //printf("Scale: %f",arrowScale_tmp);
    @synchronized(lock) {
    
        // FRAME
        glLoadIdentity();
        [indicatorFrame drawInView:view];
        
        
        // RIGHT ARROW
        if(displayRight) {
            glLoadIdentity(); // Reset model-view matrix ( NEW )
            glTranslatef(-0.3f/*xTotal / 2.0f - 0.1f*xTotal*/, 0.0f, -distance);
            glScalef(arrowScale_tmp, arrowScale_tmp, 1.0f);
            [arrowRight drawInView:view]; // Draw triangle ( NEW )
        }
        
        // LEFT
        if(displayLeft) {
            glLoadIdentity();
            glTranslatef(0.3f/*-xTotal / 2.0f + 0.1f*xTotal*/, 0.0f, -distance);
            glRotatef(180.0f, 0.0f, 0.0f, 1.0f);
            glScalef(arrowScale_tmp, arrowScale_tmp, 1.0f);
            [arrowLeft drawInView:view]; // Draw quad ( NEW )
        }
    }
    
        
    
    
}

-(void)interfaceOrientationDidChange:(UIInterfaceOrientation)orientation {

}

- (void)gluPerspective:(double)fovy :(double)aspec :(double)zNear :(double)zFar
{
    // Start in projection mode.
    glMatrixMode(GL_PROJECTION);
    glLoadIdentity();
    double xmin, xmax, ymin, ymax;
    ymax = zNear * tan(fovy * M_PI / 360.0);
    ymin = -ymax;
    xmin = ymin * aspec;
    xmax = ymax * aspec;
    glFrustumf(xmin, xmax, ymin, ymax, zNear, zFar);
}



- (void)viewDidLoad
{
    [super viewDidLoad];
	// Do any additional setup after loading the view.
    
}



-(void)setUpLocationManager {
    locationManager = [[CLLocationManager alloc] init];
    locationManager.delegate= self;
    locationManager.desiredAccuracy=kCLLocationAccuracyBestForNavigation;
    // Start heading updates.
    if ([CLLocationManager headingAvailable] && [CLLocationManager locationServicesEnabled])
    {
        locationManager.headingFilter = 5;
        [locationManager startUpdatingHeading];
        [locationManager startUpdatingLocation];
    }

}

- (void)locationManager:(CLLocationManager *)manager didUpdateToLocation:(CLLocation *)newLocation fromLocation:(CLLocation *)oldLocation {
    latitude = newLocation.coordinate.latitude;
    longitude = newLocation.coordinate.longitude;
}
- (void)locationManager:(CLLocationManager *)manager didUpdateHeading:(CLHeading *)newHeading {
    
    assert([self fixAngle:270] == -90);
    assert([self fixAngle:181] == -179);
    

    //get currentOrientation from [-180,180]
    
    //currentOrientation = [self fixAngle:[self fixAngle:tempHeading]-90] ;
    
    UIInterfaceOrientation deviceOrientation = self.interfaceOrientation;
    if (deviceOrientation == UIDeviceOrientationLandscapeRight) {
        rightSideUp = NO;
    }else{
        rightSideUp = YES;
    }
    
    if (rightSideUp) {
        currentOrientation = [self fixAngle:[self addAngleOffset:newHeading.trueHeading positive:YES]];
    }else{
        currentOrientation = [self fixAngle:[self addAngleOffset:newHeading.trueHeading positive:NO]];
    }
    
    
    
    
    //Log heading
    //NSLog(@"currentOrientation: %f",currentOrientation);
    [debugScreen setText:[NSString stringWithFormat:@"CurrentOrientation: %f \nDesired Orientation:%f",currentOrientation,desiredOrientation]];
    
    
    
    //if free roam indicate FREE and return
    if (freeRoam) {
        [self indicateTurnToDirection:FREE withPercentage:0.0f];
        return;
    }

    //get device orientation
//    UIInterfaceOrientation deviceOrientation = self.interfaceOrientation;
//    
//    if (deviceOrientation == UIDeviceOrientationLandscapeRight) {
//        rightSideUp = YES;
//    }else{
//        rightSideUp = NO;
//    }
    
    if (rightSideUp) {
        [debugScreen setText:[debugScreen.text stringByAppendingString:@"\n^"]];
    }else{
        [debugScreen setText:[debugScreen.text stringByAppendingString:@"\nv"]];
    }
    
    //graphics logic
    BOOL rightArrow = YES;
    float differenceAngle = [self fixAngle:desiredOrientation - currentOrientation];
    if (abs(differenceAngle) > orientationTolerance) {
        if (differenceAngle > 0) {
            // turn right
            rightArrow = YES;
        } else {
            // turn left
            rightArrow = NO;
        }
        
        // flip arrow incase device is flipped
        if (rightSideUp) {
            rightArrow = !rightArrow;
        }
        
        if (rightArrow) {
            [self indicateTurnToDirection:RIGHT withPercentage:abs(differenceAngle)/180.f];
        } else {
            [self indicateTurnToDirection:LEFT withPercentage:abs(differenceAngle)/180.0f];
        }

    }else{
        [self indicateTurnToDirection:NONE withPercentage:0.0f];
    }
}

-(float)addAngleOffset:(float)angle positive:(BOOL)pos{
    
    if (pos) {
        angle += 90;
    }else{
        angle -= 90;
    }
    if (angle <0) {
        angle += 360;
    }
    if (angle >= 360) {
        angle -= 360;
    }
    
    return angle;
}

-(float)fixAngle:(float)angle {
    if (angle > 180.0f) {
        angle = -180.0f + fmod(angle,180);
    } else if (angle < -180.0f) {
        angle = 180.0f - abs(angle) % 180;
    }
    return angle;
}

/*-(BOOL)shouldAutorotateToInterfaceOrientation:(UIInterfaceOrientation)toInterfaceOrientation {
    return UIInterfaceOrientationLandscapeRight;
}

-(NSUInteger)supportedInterfaceOrientations{
    return UIInterfaceOrientationMaskLandscapeRight;
}*/


- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

@end
