//go:build darwin

package openpanel

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework AppKit

#include <stdint.h>
#include <stdlib.h>

#include "panel_darwin.h"
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

// Available reports that macOS has an open panel of its own, which it has
// had for as long as there has been an AppKit.
func Available() bool { return true }

// ChooseDirectory presents the platform's open panel over the window of the
// given view and returns the directory the reader chose. ok is false where
// they cancelled, and the caller changes nothing then.
//
// view is the CFTypeRef Gio hands out in an app.AppKitViewEvent; the panel
// is a sheet on that view's window. A zero view, or one whose window has
// gone, gets a standalone panel rather than none.
//
// The panel is directories only, one of them, opened at dir, with Cancel
// and Open, and the library draws nothing of it.
//
// It blocks until the reader answers, so it belongs off the render
// goroutine — in an mvu command, not in an update. AppKit presents from the
// main thread, which is a thread this call is not on and must not become:
// the work is put on the main queue and the answer comes back on it.
func ChooseDirectory(view uintptr, dir string) (path string, ok bool) {
	// One buffered slot, so the main queue hands the answer over and
	// returns whether this goroutine has reached the receive yet or not.
	answer := make(chan string, 1)
	handle := cgo.NewHandle(answer)
	defer handle.Delete()

	start := C.CString(dir)
	defer C.free(unsafe.Pointer(start))

	C.vgio_openpanel_directory(C.uintptr_t(view), start, C.uintptr_t(handle))
	chosen := <-answer
	return chosen, chosen != ""
}

//export vgioOpenPanelChosen
func vgioOpenPanelChosen(handle C.uintptr_t, path *C.char) {
	answer := cgo.Handle(handle).Value().(chan string)
	var chosen string
	if path != nil {
		chosen = C.GoString(path)
	}
	answer <- chosen
}
