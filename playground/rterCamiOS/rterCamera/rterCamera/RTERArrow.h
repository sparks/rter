//
//  RTERArrow.h
//  rterCamera
//
//  Created by Cameron Bell on 13-03-18.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <GLKit/GLKit.h>


typedef struct {
    GLfloat x;
    GLfloat y;
    GLfloat z;
} Vertex3D;

typedef struct {
    Vertex3D v1;
    Vertex3D v2;
    Vertex3D v3;
} Triangle3D;

@interface RTERArrow : NSObject {
    //const GLfloat triangleVertices[];
    Triangle3D _triangle;
    
}







-(id)initArrow;
-(void)drawInView:(GLKView *)view;

@end
