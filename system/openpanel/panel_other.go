//go:build !darwin

package openpanel

// Available reports that this platform has no open panel of its own to
// present. The caller raises its own chooser instead; nothing here is a
// failure, and no error is reported for one.
func Available() bool { return false }

// ChooseDirectory answers nothing on a platform with no panel. It is here so
// that a caller compiles and links everywhere, not so that it is called:
// [Available] is what the caller branches on.
func ChooseDirectory(view uintptr, dir string) (path string, ok bool) {
	return "", false
}
