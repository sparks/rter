//
//  RTERGLKViewController.h
//  rterCamera
//
//  Created by Cameron Bell on 13-03-18.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import <GLKit/GLKit.h>
#import "RTERArrow.h"
#import "RTERIndicatorFrame.h"
#import "Config.h"
#import <CoreLocation/CoreLocation.h>
#import "RTERPreviewController.h"
//@class RTERPreviewController;

@interface RTERGLKViewController : GLKViewController <GLKViewDelegate,CLLocationManagerDelegate,NSURLConnectionDelegate> {
    float _curRed;
    BOOL _increasing;
    GLKView* _glkView;

    //What is application context? OpenGL context?
    EAGLContext *context; // Application's context
    
    RTERArrow *arrowLeft;
    RTERArrow *arrowRight;
    RTERIndicatorFrame *indicatorFrame;

    float aspect;
    float xTotal, yTotal, distance;
    
    float arrowScale;
    float arrowScaleMax;
    float arrowScaleMin;
    
    float arrowPulsateScale;
	float arrowPulsateSpeed;
	float arrowPulsateSpeedMin;
	float arrowPulsateMax;
	float arrowPulsateMin;
	BOOL arrowPulsateIncrease;
    
	BOOL displayLeft;
	BOOL displayRight;
    
    BOOL freeRoam;
    BOOL rightSideUp;
    float desiredOrientation;
    float currentOrientation;
    float orientationTolerance;
    
    NSTimer *getHeadingTimer;
    NSTimer *backgroundUpdateTimer;
    NSURLConnection *currentGetConnection;
    NSURLConnection *currentPutConnection;
    NSMutableData *headingData;
    
    CLLocationManager *locationManager;
    
    UITextView *debugScreen;
}

/*public static enum Indicate {
    LEFT, RIGHT, NONE, FREE
}*/

@property (assign) BOOL streaming;

typedef NS_ENUM(NSInteger, Indicate) {
   LEFT, RIGHT, NONE, FREE
};




- (id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil view:(GLKView *)view previewController:(UIViewController *)prev;
-(void)update;
-(void)indicateTurnToDirection:(Indicate) direction withPercentage:(float) percentage;
-(void)interfaceOrientationDidChange:(UIInterfaceOrientation)orientation;
-(void)setDesiredOrientation:(float)dO;
-(void)startGetPutTimer;
-(void)stopGetPutTimer;
-(void)startBackgroundUpdateTimer;
-(void)stopBackgroundUpdateTimer;
-(void)onSurfaceChange;
-(void)onSurfaceChangedWidth:(float)width Height:(float)height;
-(void)currentFPS:(float)currentFPS;
@end
