//go:build darwin

#import <Foundation/Foundation.h>

#include "naming_darwin.h"

int vgio_naming_compare(const char *a, const char *b) {
	@autoreleasepool {
		NSString *lhs = [NSString stringWithUTF8String:a];
		NSString *rhs = [NSString stringWithUTF8String:b];
		if (lhs == nil) {
			lhs = @"";
		}
		if (rhs == nil) {
			rhs = @"";
		}
		switch ([lhs localizedStandardCompare:rhs]) {
			case NSOrderedAscending:  return -1;
			case NSOrderedDescending: return 1;
			default:                  return 0;
		}
	}
}
