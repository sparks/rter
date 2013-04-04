//
//  rterCameraViewController.h
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-05.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import <UIKit/UIKit.h>
#import "RTERPreviewController.h"

@interface RTERViewController : UIViewController<RTERPreviewControllerDelegate,NSURLConnectionDataDelegate,UITextFieldDelegate>

@property(nonatomic, strong) IBOutlet UITextField *userField;
@property(nonatomic, strong) IBOutlet UITextField *passField;
@property(nonatomic, retain) NSString *cookieString;
@property(nonatomic, retain) NSString *userName;

- (IBAction)startCamera:(id)sender;
- (IBAction)login:(id)sender;

@end
