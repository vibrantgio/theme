//go:build darwin

#import <AppKit/AppKit.h>

#include "panel_darwin.h"

// The Go side of the answer (panel_darwin.go), declared the way the other
// Objective-C in the organization declares its Go exports rather than through
// the generated header. The path is owned by the panel's autorelease pool and
// is valid only for the duration of this call, which is why the Go side copies
// it there and then.
extern void vgioOpenPanelChosen(uintptr_t handle, char *path);

void vgio_openpanel_directory(uintptr_t view, const char *dir, uintptr_t handle) {
	@autoreleasepool {
		// The start directory is turned into a URL here, on the calling
		// thread, rather than inside the block: the C string belongs to the
		// caller and is freed once the answer is in, which is after the block
		// has run but not before it was scheduled.
		NSURL *start = nil;
		if (dir != NULL && dir[0] != '\0') {
			NSString *path = [NSString stringWithUTF8String:dir];
			if (path != nil) {
				start = [NSURL fileURLWithPath:path isDirectory:YES];
			}
		}
		dispatch_async(dispatch_get_main_queue(), ^{
			@autoreleasepool {
				NSOpenPanel *panel = [NSOpenPanel openPanel];
				panel.canChooseDirectories = YES;
				panel.canChooseFiles = NO;
				panel.allowsMultipleSelection = NO;
				panel.prompt = @"Open";
				if (start != nil) {
					panel.directoryURL = start;
				}
				void (^answer)(NSModalResponse) = ^(NSModalResponse response) {
					NSURL *url = nil;
					if (response == NSModalResponseOK) {
						url = [[panel URLs] firstObject];
					}
					if (url == nil) {
						vgioOpenPanelChosen(handle, NULL);
						return;
					}
					vgioOpenPanelChosen(handle, (char *)[[url path] UTF8String]);
				};
				NSWindow *host = [(__bridge NSView *)(void *)view window];
				if (host != nil) {
					[panel beginSheetModalForWindow:host completionHandler:answer];
					return;
				}
				[panel beginWithCompletionHandler:answer];
			}
		});
	}
}
