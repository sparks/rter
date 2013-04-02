//
//  RTERIndicatorFrame.m
//  rterCamera
//
//  Created by Cameron Bell on 13-03-19.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import "RTERIndicatorFrame.h"


@interface RTERIndicatorFrame () {
    Colour currentColour;
}
@end

@implementation RTERIndicatorFrame

static const GLfloat vertices[48] = {
    // TOP
    -1.0f, 0.9f, 0.0f, // 0. left-bottom-top
    1.0f, 0.9f, 0.0f, // 1. right-bottom-top
    -1.0f, 1.0f, 0.0f, // 2. left-top-top
    1.0f, 1.0f, 0.0f, // 3. right-top-top
    // LEFT
    -1.0f, -0.9f, 0.0f, // 0. left-bottom-left
    -0.9f, -0.9f, 0.0f, // 1. right-bottom-left
    -1.0f, 0.9f, 0.0f, // 2. left-top-left
    -0.9f, 0.9f, 0.0f, // 3. right-top-left
    // BOTTOM
    -1.0f, -1.0f, 0.0f, // 0. left-bottom-bottom
    1.0f, -1.0f, 0.0f, // 1. right-bottom-bottom
    -1.0f, -0.9f, 0.0f, // 2. left-top-bottom
    1.0f, -0.9f, 0.0f, // 3. right-top-bottom
    // RIGHT
    0.9f, -0.9f, 0.0f, // 0. left-bottom-left
    1.0f, -0.9f, 0.0f, // 1. right-bottom-left
    0.9f, 0.9f, 0.0f, // 2. left-top-left
    1.0f, 0.9f, 0.0f // 3. right-top-left
};

-(id)initIndicatorFrame{
    if (self = [super init]) {
        
        
//        for (int i = 0; i<48; i++) {
//            NSLog(@"Vert: %f",vertices[i]);
//        }
       
        currentColour = BLUE;

    }
    
    return self;
        
        
}
-(void)setColour:(NSInteger)colour {
    currentColour = colour;
}

-(void)drawInView:(GLKView *)view {
    glLoadIdentity();
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    glEnableClientState(GL_VERTEX_ARRAY);
    
    glVertexPointer(3, GL_FLOAT, 0, &vertices);
    
    // Render all the faces
    for (int face = 0; face < 4; face++) {
        // Set the color for each of the faces
        switch (currentColour) {
			case RED:
				glColor4f(0.9f, 0.0f, 0.0f, 1.0f); //0.8f);
				break;
			case BLUE:
                glColor4f(0.0f, 0.0f, 0.9f, 1.0f); //0.8f);
				break;
			case GREEN:
                glColor4f(0.0f, 0.9f, 0.0f, 1.0f); //0.8f);
				break;
        }
        
        // Draw the primitive from the vertex-array directly
        glDrawArrays(GL_TRIANGLE_STRIP, face * 4, 4);
    }
    glDisableClientState(GL_VERTEX_ARRAY);
   }


@end
