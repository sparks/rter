//
//  rterCameraAppDelegate.h
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-05.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import <UIKit/UIKit.h>

@class RTERViewController;

@interface RTERAppDelegate : UIResponder <UIApplicationDelegate>

@property (strong, nonatomic) UIWindow *window;

@property (strong, nonatomic) RTERViewController *viewController;

@end
