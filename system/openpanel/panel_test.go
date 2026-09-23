package openpanel

import (
	"runtime"
	"testing"
)

// TestOnlyMacOSOffersAPanel pins the one answer a caller branches on. The
// panel is macOS's own and is presented through the platform's framework, so
// a build that cannot reach one — any other platform, or macOS built without
// cgo — has no panel to offer. The panel itself is not driven here:
// presenting it would need a window and a reader to answer it.
func TestOnlyMacOSOffersAPanel(t *testing.T) {
	want := runtime.GOOS == "darwin" && builtWithCgo
	if got := Available(); got != want {
		t.Errorf("Available() = %v on %s, want %v", got, runtime.GOOS, want)
	}
}

// TestWithNoPanelNothingIsChosen requires the stub to answer nothing rather
// than an empty path that reads as a choice, so a caller that ignores
// [Available] still cannot open a vault at "".
func TestWithNoPanelNothingIsChosen(t *testing.T) {
	if Available() {
		t.Skip("this platform has a panel of its own")
	}
	if path, ok := ChooseDirectory(0, "/vaults/Second Brain"); ok || path != "" {
		t.Errorf("ChooseDirectory answered %q, %v; want nothing chosen", path, ok)
	}
}
