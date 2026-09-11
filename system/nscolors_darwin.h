// The C surface of the package's Objective-C (nscolors_darwin.m). The Go
// side (nscolors_darwin.go) explains which names it asks for and why the
// answer is cached.

#ifndef VGIO_THEME_SYSTEM_NSCOLORS_DARWIN_H
#define VGIO_THEME_SYSTEM_NSCOLORS_DARWIN_H

// vgio_theme_nscolor resolves one of AppKit's semantic colours by the name
// of its NSColor class property — "windowBackgroundColor", "systemRed" —
// under the aqua appearance (dark == 0) or darkAqua (dark != 0), converts it
// to sRGB and writes it to rgba as four non-premultiplied bytes, r g b a.
//
// It returns 1 on success and 0 when NSColor has no property of that name
// (an older system than the name) or the colour cannot be converted to sRGB
// (a pattern or a catalogue colour without a component representation); on 0
// nothing is written. The caller keeps its recorded value for that name.
//
// No window and no NSApplication are required, and none is created: the
// drawing appearance this reads under is thread-local and set for the
// duration of the call, so any thread may call it and no application state
// is touched.
extern int vgio_theme_nscolor(const char *name, int dark, unsigned char *rgba);

#endif
