//
//  rterCameraViewController.m
//  rterCamera
//
//  Created by Stepan Salenikovich on 2013-03-05.
//  Copyright (c) 2013 rtER. All rights reserved.
//

#import "RTERViewController.h"

@interface RTERViewController ()

@end

@implementation RTERViewController

- (void)viewDidLoad
{
    [super viewDidLoad];
	// Do any additional setup after loading the view, typically from a nib.
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (IBAction)startCamera:(id)sender {
    RTERPreviewController *preview = [[RTERPreviewController alloc] init];
    
    preview.delegate = self;
    preview.modalTransitionStyle = UIModalTransitionStyleFlipHorizontal;
    
    [self presentViewController:preview
                       animated:YES
                     completion:nil];
    
}

- (void)back {
    [self dismissViewControllerAnimated:YES completion: nil];
}

@end
