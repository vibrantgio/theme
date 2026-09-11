//go:build darwin

#import <AppKit/AppKit.h>
#include <math.h>

#include "nscolors_darwin.h"

int vgio_theme_nscolor(const char *name, int dark, unsigned char *rgba) {
	@autoreleasepool {
		NSString *property = [NSString stringWithUTF8String:name];
		if (property == nil) {
			return 0;
		}
		SEL sel = NSSelectorFromString(property);
		// Every semantic colour is a class property of NSColor, so asking
		// whether the class responds is the version check: a name this
		// system does not carry answers no and the caller keeps its
		// recorded value rather than a guess.
		if (sel == NULL || ![NSColor respondsToSelector:sel]) {
			return 0;
		}
		NSAppearance *appearance = [NSAppearance appearanceNamed:
			(dark ? NSAppearanceNameDarkAqua : NSAppearanceNameAqua)];
		if (appearance == nil) {
			return 0;
		}
		__block int ok = 0;
		// A dynamic colour answers for whichever appearance is current on
		// this thread, so the resolution happens inside the appearance
		// rather than around it.
		[appearance performAsCurrentDrawingAppearance:^{
			NSColor *c = [(id)[NSColor class] performSelector:sel];
			NSColor *srgb = [c colorUsingColorSpace:[NSColorSpace sRGBColorSpace]];
			if (srgb == nil) {
				return;
			}
			rgba[0] = (unsigned char)lround([srgb redComponent] * 255.0);
			rgba[1] = (unsigned char)lround([srgb greenComponent] * 255.0);
			rgba[2] = (unsigned char)lround([srgb blueComponent] * 255.0);
			rgba[3] = (unsigned char)lround([srgb alphaComponent] * 255.0);
			ok = 1;
		}];
		return ok;
	}
}
