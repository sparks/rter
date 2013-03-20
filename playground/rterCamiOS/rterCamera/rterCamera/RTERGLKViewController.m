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
}
@end

@implementation RTERGLKViewController

- (id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil view:(GLKView *)view
{
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self) {
        // Custom initialization
        _curRed = 0.0;
        _increasing = YES;

        _glkView = view;
        _glkView.delegate = self;
        [_glkView setNeedsDisplay];
        self.preferredFramesPerSecond = 15;

        self.view = _glkView;
        self.view.opaque = NO;
        self.view.backgroundColor = [UIColor clearColor];

        lock = [NSNumber numberWithBool:NO];
        
        arrowLeft = [[RTERArrow alloc]initArrow];
        arrowRight = [[RTERArrow alloc]initArrow];
        indicatorFrame = [[RTERIndicatorFrame alloc]initIndicatorFrame];
        

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
    }
    return self;
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


-(void)indicateTurnToDirection:(Indicate) direction withPercentage:(float) percentage {
    @synchronized(lock)
    {
        arrowScale = (arrowScaleMax - arrowScaleMin) * percentage +arrowScaleMin;
        
        switch (direction) {
            case LEFT:
                displayLeft = true;
				displayRight = false;
				[indicatorFrame setColour:RED];
                break;
            case RIGHT:
                displayLeft = false;
				displayRight = true;
                [indicatorFrame setColour:RED];
                break;
            case NONE:
                displayLeft = false;
				displayRight = false;
                [indicatorFrame setColour:GREEN];
            case FREE:
                displayLeft = false;
				displayRight = false;
                [indicatorFrame setColour:BLUE];
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
    printf("Scale: %f",arrowScale_tmp);
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
    /*
    float width = _glkView.frame.size.width;
    float height = _glkView.frame.size.height;
    
    if (height == 0) {
        height = 1;
    }
    aspect = width/height;
    
    distance = 6.0f;
    xTotal = (float) aspect * tanf(rad(45.0/2));
    yTotal = (float) (tanf(rad(45.0/2))*distance *2);
    
    //indicatorFrame.resize(xTotal, yTotal, distance);
    
    // Set the viewport (display area) to cover the entire window
    glViewport(0, 0, width, height);
    
    // Setup perspective projection, with aspect ratio matches viewport
    glMatrixMode(GL_PROJECTION); // Select projection matrix
    glLoadIdentity(); // Reset projection matrix
    // Use perspective projection
    [self gluPerspective:45 :aspect :0.1f :100.f]; //see gluPperspective:::: below
    
    glMatrixMode(GL_MODELVIEW); // Select model-view matrix
    glLoadIdentity(); // Reset
     */
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

/*if (height == 0)
 height = 1; // To prevent divide by zero
 aspect = (float) width / height;
 
 // get the total x and y at distance
 distance = 6.0f;
 xTotal = (float) (aspect * Math.tan(Math.toRadians(45.0 / 2))
 * distance * 2);
 yTotal = (float) (Math.tan(Math.toRadians(45.0 / 2)) * distance * 2);
 
 indicatorFrame.resize(xTotal, yTotal, distance);
 
 // Set the viewport (display area) to cover the entire window
 gl.glViewport(0, 0, width, height);
 
 // Setup perspective projection, with aspect ratio matches viewport
 gl.glMatrixMode(GL10.GL_PROJECTION); // Select projection matrix
 gl.glLoadIdentity(); // Reset projection matrix
 // Use perspective projection
 GLU.gluPerspective(gl, 45, aspect, 0.1f, 100.f);
 
 gl.glMatrixMode(GL10.GL_MODELVIEW); // Select model-view matrix
 gl.glLoadIdentity(); // Reset
 */

- (void)viewDidLoad
{
    [super viewDidLoad];
	// Do any additional setup after loading the view.
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

@end
