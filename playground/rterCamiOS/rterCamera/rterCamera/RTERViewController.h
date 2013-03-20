//
//  rterCameraViewController.h
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-05.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import <UIKit/UIKit.h>
#import "RTERPreviewController.h"


@interface RTERViewController : UIViewController<RTERPreviewControllerDelegate> {

}

- (IBAction)startCamera:(id)sender;

@end
