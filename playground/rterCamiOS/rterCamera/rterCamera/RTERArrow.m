//
//  RTERArrow.m
//  rterCamera
//
//  Created by Cameron Bell on 13-03-18.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import "RTERArrow.h"

#define f(X) [NSNumber numberWithFloat:X]

@interface RTERArrow () {

}
@end

@implementation RTERArrow

static const GLfloat triangle[12] = {
    -0.5f, -1.0f, 0.0f, // 0. left-bottom
    0.5f, 0.0f, 0.0f, // 1. right-bottom
    -0.5f, 1.0f, 0.0f // 2. left-top
};


-(id)initArrow {
    
    if (self = [super init]) {

        Vertex3D    vertex1 = Vertex3DMake(.7, -0.2, 0.0);
        Vertex3D    vertex2 = Vertex3DMake(0.85, 0.0, 0.0);
        Vertex3D    vertex3 = Vertex3DMake(0.7, 0.2, 0.0);
        _triangle = Triangle3DMake(vertex1, vertex2, vertex3);
        
        
        
        
        
  
    }
    return self;
}

static inline GLfloat Vertex3DCalculateDistanceBetweenVertices (Vertex3D first, Vertex3D second)
{
    GLfloat deltaX = second.x - first.x;
    GLfloat deltaY = second.y - first.y;
    GLfloat deltaZ = second.z - first.z;
    return sqrtf(deltaX*deltaX + deltaY*deltaY + deltaZ*deltaZ );
};
static inline Vertex3D Vertex3DMake(CGFloat inX, CGFloat inY, CGFloat inZ)
{
    Vertex3D ret;
    ret.x = inX;
    ret.y = inY;
    ret.z = inZ;
    return ret;
}
static inline Triangle3D Triangle3DMake(Vertex3D inX, Vertex3D inY, Vertex3D inZ)
{
    Triangle3D ret;
    ret.v1 = inX;
    ret.v2 = inY;
    ret.v3 = inZ;
    return ret;
}

-(void)drawInView:(GLKView *)view {
    //float arrowScale_tmp = 2.0f;
    //glLoadIdentity();
    //glClearColor(0.0, 0.0, 0.0,0.0);
    //glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    glEnableClientState(GL_VERTEX_ARRAY);
    glColor4f(1.0, 0.0, 0.0, 1.0);
    glVertexPointer(3, GL_FLOAT, 0, triangle); //&_triangle);
    //glRotatef(180.0f, 0.0f, 0.0f, 1.0f);
    //glScalef(1.0f, arrowScale_tmp, 1.0f);
    glDrawArrays(GL_TRIANGLES, 0, 3);
    glDisableClientState(GL_VERTEX_ARRAY);
    
}
 
@end
