package system

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit

#include <stdlib.h>

#include "nscolors_darwin.h"
*/
import "C"

import (
	"image/color"
	"reflect"
	"strings"
	"unsafe"

	"github.com/vibrantgio/theme/tokens"
)

// liveSet asks AppKit for every name in the platform set under one
// appearance — aqua or darkAqua — and returns the set it answers. It starts
// from the recorded set for that appearance and overwrites each field the
// platform answers for, so a name this system does not carry keeps the
// recorded value rather than a zero.
//
// The names come from the token set itself: a field is AppKit's own name
// with the trailing "Color" dropped ([appKitName] inverts that), so a field
// added to tokens.PlatformColors is read live without a list here to
// maintain.
//
// Two kinds of field are not asked for, and neither is a list here either.
// FindHighlight carries the find highlight as Mail paints it, not what
// AppKit answers for findHighlightColor, which the platform's own
// applications do not paint. The measured materials — the chrome material,
// the card's fill, the two state overlays, the floating shadow — carry the
// struct tag `appkit:"-"`, which says the platform gives that fill no
// NSColor name at all: they were read off the stored captures, and asking
// AppKit for a name it does not answer would leave them at the recorded
// value anyway. Skipping them keeps the naming rule honest rather than
// relying on that.
//
// Every value is non-premultiplied sRGB with the alpha AppKit reports, to
// the byte: the recorded sets quantize the catalogue's two-decimal alpha,
// and this does not.
func liveSet(dark bool) tokens.PlatformColors {
	p := tokens.PlatformLight
	if dark {
		p = tokens.PlatformDark
	}
	set := reflect.ValueOf(&p).Elem()
	typ := set.Type()
	for i := range typ.NumField() {
		if typ.Field(i).Tag.Get("appkit") == "-" {
			continue
		}
		field := typ.Field(i).Name
		if field == "FindHighlight" {
			continue
		}
		if c, ok := nsColor(appKitName(field), dark); ok {
			set.Field(i).Set(reflect.ValueOf(c))
		}
	}
	return p
}

// nsColor resolves one NSColor class property under one appearance. ok is
// false where the platform has no such name or cannot express it in sRGB.
func nsColor(name string, dark bool) (c color.NRGBA, ok bool) {
	property := C.CString(name)
	defer C.free(unsafe.Pointer(property))
	var rgba [4]C.uchar
	appearance := C.int(0)
	if dark {
		appearance = 1
	}
	if C.vgio_theme_nscolor(property, appearance, &rgba[0]) == 0 {
		return color.NRGBA{}, false
	}
	return color.NRGBA{R: uint8(rgba[0]), G: uint8(rgba[1]), B: uint8(rgba[2]), A: uint8(rgba[3])}, true
}

// appKitName is the token set's naming rule read backwards: a field is
// AppKit's name with the trailing "Color" dropped and nothing else changed,
// and the system colours carry no such suffix.
func appKitName(field string) string {
	name := strings.ToLower(field[:1]) + field[1:]
	if strings.HasPrefix(field, "System") {
		return name
	}
	return name + "Color"
}
