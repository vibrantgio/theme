// Package openpanel presents the platform's own chooser for a folder.
//
// Choosing a folder is the platform's open panel where the platform offers
// one, and the library's modal where it does not, so this package answers
// two questions: whether a panel exists at all ([Available]), and what the
// reader chose in it ([ChooseDirectory]). An application asks the first
// before it decides which chooser to raise, and the second only where the
// first said yes.
//
// It is a package of its own rather than another file beside the colour
// shim in theme/system. That shim reads values and touches no application
// state — its own header states that it needs no window and no
// NSApplication, and that any thread may call it. A panel is the opposite
// on both counts: it is a window presented over the application's own, and
// AppKit will only present it from the main thread. Keeping the two apart
// keeps that promise readable.
//
// The library draws nothing of the panel. What it looks like, which
// shortcuts it carries and where it opens scrolled to are the platform's,
// and no token reaches it.
package openpanel
