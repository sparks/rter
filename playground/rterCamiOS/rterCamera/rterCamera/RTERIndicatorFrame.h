//
//  RTERIndicatorFrame.h
//  rterCamera
//
//  Created by Cameron Bell on 13-03-19.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <GLKit/GLKit.h>
#import "Config.h"



@interface RTERIndicatorFrame : NSObject {

    
    
}


-(void)drawInView:(GLKView *)view;
-(id)initIndicatorFrame;
-(void)setColour:(NSInteger)colour;

@end
