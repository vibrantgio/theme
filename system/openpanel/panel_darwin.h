// The C surface of the package's Objective-C (panel_darwin.m). The Go side
// (panel_darwin.go) states what the panel is configured to ask for and why
// the call blocks.

#ifndef VGIO_THEME_SYSTEM_OPENPANEL_PANEL_DARWIN_H
#define VGIO_THEME_SYSTEM_OPENPANEL_PANEL_DARWIN_H

#include <stdint.h>

// vgio_openpanel_directory presents an NSOpenPanel that chooses one
// directory, starting at dir when dir names one, and answers by calling
// vgioOpenPanelChosen with handle and the chosen path — or with NULL where
// the reader cancelled.
//
// view is the CFTypeRef of the application's NSView, which the panel is
// presented as a sheet on the window of. Where view is zero, or its window
// has gone, the panel is presented standalone instead: a chooser the reader
// can answer is better than none.
//
// It returns at once. The panel is put up on the main queue, which is the
// only thread AppKit may present from, and the answer arrives on that queue
// whenever the reader gives one — so the caller must not wait on the thread
// this is called from being the one that finishes the work.
extern void vgio_openpanel_directory(uintptr_t view, const char *dir, uintptr_t handle);

#endif
