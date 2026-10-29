package system

import (
	"context"
	"image/color"
	"testing"
	"time"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/theme/a11y"
	"github.com/vibrantgio/theme/tokens"
)

// withLiveAccent replaces the process's live reader for one test. The stub
// reader returns the recorded set for each appearance with ControlAccent
// set to accent, so no test reads this machine's own accent.
func withLiveAccent(t *testing.T, accent color.NRGBA) {
	t.Helper()
	read := func(dark bool) tokens.PlatformColors {
		p := tokens.PlatformLight
		if dark {
			p = tokens.PlatformDark
		}
		p.ControlAccent = color.NRGBA{R: accent.R, G: accent.G, B: accent.B, A: p.ControlAccent.A}
		return p
	}
	saved := livePlatform
	livePlatform = &platformCache{interval: time.Hour, now: time.Now, read: read}
	t.Cleanup(func() { livePlatform = saved })
}

// unpinnedPlatform is the set a stream with no theme colour emits for a.
func unpinnedPlatform(t *testing.T, a Appearance) tokens.PlatformColors {
	t.Helper()
	th := newConfig(nil).theme(rx.Tuple2[Appearance, a11y.A11yPrefs]{First: a})
	var sets []tokens.PlatformColors
	err := th.Platform.Subscribe(context.Background(), func(v tokens.PlatformColors, _ error, done bool) {
		if !done {
			sets = append(sets, v)
		}
	}).Wait()
	if err != nil || len(sets) != 1 {
		t.Fatalf("platform observe: %d sets, err %v", len(sets), err)
	}
	return sets[0]
}

// TestTheMeasuredAccentFillsFollowTheLiveAccent pins that with a live accent
// other than the recorded blue and no theme colour, the four measured
// accent fills equal WithAccent's values for that accent. The NSColor-named
// rows keep what the reader answered.
func TestTheMeasuredAccentFillsFollowTheLiveAccent(t *testing.T) {
	purple := tokens.PlatformLight.SystemPurple
	withLiveAccent(t, purple)
	for _, a := range []Appearance{{}, {Dark: true}} {
		recorded := tokens.PlatformLight
		if a.Dark {
			recorded = tokens.PlatformDark
		}
		want := recorded.WithAccent(purple)
		got := unpinnedPlatform(t, a)
		if got.SidebarSelection != want.SidebarSelection {
			t.Errorf("dark %v: SidebarSelection = %v, want %v", a.Dark, got.SidebarSelection, want.SidebarSelection)
		}
		if got.SidebarSelectionUnemphasizedLabel != want.SidebarSelectionUnemphasizedLabel {
			t.Errorf("dark %v: SidebarSelectionUnemphasizedLabel = %v, want %v", a.Dark, got.SidebarSelectionUnemphasizedLabel, want.SidebarSelectionUnemphasizedLabel)
		}
		if got.DefaultButtonFill != want.DefaultButtonFill {
			t.Errorf("dark %v: DefaultButtonFill = %v, want %v", a.Dark, got.DefaultButtonFill, want.DefaultButtonFill)
		}
		if got.AccentFill != want.AccentFill {
			t.Errorf("dark %v: AccentFill = %v, want %v", a.Dark, got.AccentFill, want.AccentFill)
		}
		if got.SidebarSelection == recorded.SidebarSelection {
			t.Errorf("dark %v: SidebarSelection kept the measured blue under a purple accent", a.Dark)
		}
		read := livePlatform.read(a.Dark)
		if got.SelectedContentBackground != read.SelectedContentBackground {
			t.Errorf("dark %v: SelectedContentBackground = %v, want the reader's %v", a.Dark, got.SelectedContentBackground, read.SelectedContentBackground)
		}
		if got.SidebarSelectionUnemphasized != recorded.SidebarSelectionUnemphasized {
			t.Errorf("dark %v: SidebarSelectionUnemphasized moved; it is grey and does not follow the accent", a.Dark)
		}
	}
}

// TestTheMeasuredAccentFillsKeepTheirValuesOnTheRecordedBlue pins that with
// the recorded blue as the live accent, the four measured accent fills
// equal the measured values.
func TestTheMeasuredAccentFillsKeepTheirValuesOnTheRecordedBlue(t *testing.T) {
	withLiveAccent(t, tokens.PlatformLight.ControlAccent)
	for _, a := range []Appearance{{}, {Dark: true}} {
		recorded := tokens.PlatformLight
		if a.Dark {
			recorded = tokens.PlatformDark
		}
		got := unpinnedPlatform(t, a)
		if got.SidebarSelection != recorded.SidebarSelection {
			t.Errorf("dark %v: SidebarSelection = %v, want the measured %v", a.Dark, got.SidebarSelection, recorded.SidebarSelection)
		}
		if got.SidebarSelectionUnemphasizedLabel != recorded.SidebarSelectionUnemphasizedLabel {
			t.Errorf("dark %v: SidebarSelectionUnemphasizedLabel = %v, want the measured %v", a.Dark, got.SidebarSelectionUnemphasizedLabel, recorded.SidebarSelectionUnemphasizedLabel)
		}
		if got.DefaultButtonFill != recorded.DefaultButtonFill {
			t.Errorf("dark %v: DefaultButtonFill = %v, want the measured %v", a.Dark, got.DefaultButtonFill, recorded.DefaultButtonFill)
		}
		if got.AccentFill != recorded.AccentFill {
			t.Errorf("dark %v: AccentFill = %v, want the measured %v", a.Dark, got.AccentFill, recorded.AccentFill)
		}
		if got != recorded {
			t.Errorf("dark %v: the set on the recorded blue differs from the recorded set", a.Dark)
		}
	}
}

// TestAPinnedThemeColourSkipsTheLiveAccentFollow pins that a pinned theme
// colour rebuilds the set AppKit answered, not a set the live accent has
// already moved. The pinned colour is the recorded blue and the live accent
// is purple.
func TestAPinnedThemeColourSkipsTheLiveAccentFollow(t *testing.T) {
	blue := tokens.PlatformLight.ControlAccent
	withLiveAccent(t, tokens.PlatformLight.SystemPurple)
	for _, dark := range []bool{false, true} {
		th := newConfig([]Option{WithThemeColor(blue)}).theme(rx.Tuple2[Appearance, a11y.A11yPrefs]{First: Appearance{Dark: dark}})
		var got tokens.PlatformColors
		if err := th.Platform.Subscribe(context.Background(), func(v tokens.PlatformColors, _ error, done bool) {
			if !done {
				got = v
			}
		}).Wait(); err != nil {
			t.Fatalf("platform observe: %v", err)
		}
		if want := livePlatform.read(dark).WithAccent(blue); got != want {
			t.Errorf("dark %v: the pinned set is not the reader's set rebuilt for the pinned colour", dark)
		}
	}
}
