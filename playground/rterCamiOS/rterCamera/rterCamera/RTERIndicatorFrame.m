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
    GLfloat frameVertices[48];
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
    
    [self resizeWithX:1.0f Y:1.0f distance:0.0f];
    
    return self;
        
        
}
-(void)setColour:(NSInteger)colour {
    currentColour = colour;
}

-(void)drawInView:(GLKView *)view {
    glLoadIdentity();
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    glEnableClientState(GL_VERTEX_ARRAY);
    
    glVertexPointer(3, GL_FLOAT, 0, &frameVertices);
    
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
//        glDrawArrays(GL_TRIANGLES, face * 4, 4);
    }
    glDisableClientState(GL_VERTEX_ARRAY);
}

-(void)resizeWithX:(float)xTotal Y:(float)yTotal distance:(float)distance {
    float framePercentWidth = 0.025f;
    float frameWidth;
    //use the smallest dimension to determine frame width
    if (xTotal < yTotal) {
        frameWidth = framePercentWidth*xTotal;
    } else {
        frameWidth = framePercentWidth*yTotal;
    }
    float right = xTotal/2.0f;
    float left = -right;
    float top = yTotal/2.0f;
    float bottom = -top;
    distance = -distance;
    
    NSLog(@"frame right %f, left %f, top %f, bottom %f", right, left, top, bottom);
    
    
    float vertices_tmp[48] =
    { // Vertices for the arrow
        // TOP
        left, top-frameWidth, distance, // 0. left-bottom-top
        right, top-frameWidth, distance, // 1. right-bottom-top
        left, top, distance, // 2. left-top-top
        right, top, distance, // 3. right-top-top
        // LEFT
        left, bottom+frameWidth, distance, // 0. left-bottom-left
        left+frameWidth, bottom+frameWidth, distance, // 1. right-bottom-left
        left, top-frameWidth, distance, // 2. left-top-left
        left+frameWidth, top-frameWidth, distance, // 3. right-top-left
        // BOTTOM
        left, bottom, distance, // 0. left-bottom-bottom
        right, bottom, distance, // 1. right-bottom-bottom
        left, bottom+frameWidth, distance, // 2. left-top-bottom
        right, bottom+frameWidth, distance, // 3. right-top-bottom
        // RIGHT
        right-frameWidth, bottom+frameWidth, distance, // 0. left-bottom-left
        right, bottom+frameWidth, distance, // 1. right-bottom-left
        right-frameWidth, top-frameWidth, distance, // 2. left-top-left
        right, top-frameWidth, distance // 3. right-top-left
    };
    
    memcpy(frameVertices, vertices_tmp, sizeof(float)*48);
}


@end
