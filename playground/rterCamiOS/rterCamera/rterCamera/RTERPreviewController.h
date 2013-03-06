//
//  previewController.h
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-06.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol RTERPreviewControllerDelegate <NSObject>

@required
- (void)back;

@end

@interface RTERPreviewController : UIViewController

@property (nonatomic, retain) NSObject<RTERPreviewControllerDelegate> *delegate;

- (IBAction)clickedStart:(id)sender;

- (IBAction)clickedBack:(id)sender;

@end
