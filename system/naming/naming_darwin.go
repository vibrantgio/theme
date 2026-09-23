//go:build darwin

package naming

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation

#include <stdlib.h>

#include "naming_darwin.h"
*/
import "C"

import "unsafe"

// platformCompare asks Foundation for the order, so a list draws the order
// the reader already sees in the platform's own file browser.
func platformCompare(a, b string) int {
	ca := C.CString(a)
	defer C.free(unsafe.Pointer(ca))
	cb := C.CString(b)
	defer C.free(unsafe.Pointer(cb))
	return int(C.vgio_naming_compare(ca, cb))
}
