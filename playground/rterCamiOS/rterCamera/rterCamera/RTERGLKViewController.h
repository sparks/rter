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

@interface RTERGLKViewController : GLKViewController <GLKViewDelegate> {
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

}

/*public static enum Indicate {
    LEFT, RIGHT, NONE, FREE
}*/

typedef NS_ENUM(NSInteger, Indicate) {
   LEFT, RIGHT, NONE, FREE
};




-(id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil view:(GLKView *)view;
-(void)update;
-(void)indicateTurnToDirection:(Indicate) direction withPercentage:(float) percentage;
-(void)interfaceOrientationDidChange:(UIInterfaceOrientation)orientation;
@end
