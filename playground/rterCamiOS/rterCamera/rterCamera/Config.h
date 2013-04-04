//
//  Config.h
//  rterCamera
//
//  Created by Cameron Bell on 13-03-19.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#ifndef rterCamera_Config_h
#define rterCamera_Config_h

#define SERVER @"rter.cim.mcgill.ca"
//@"142.157.58.153:8080"
//

//Camera FPS
#define DESIRED_FPS 7
#define DESIRED_FPS_IPHONE5 15

#define OPENGL_FPS 60

#define SERVER_GEO_UPDATE_PERIOD 1 //in seconds, period at which we will poll the server and send our heading

typedef NS_ENUM(NSInteger,Colour) {
    RED,
    GREEN,
    BLUE,
    CLEAR
};

// for checking if we have an iPhone 5
#define HEIGHT_IPHONE_5 568
#define IS_IPHONE_5 ([[UIScreen mainScreen] bounds ].size.height == HEIGHT_IPHONE_5 )


#endif
