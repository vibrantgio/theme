package system

import (
	"sync"
	"time"

	"github.com/vibrantgio/theme/tokens"
)

// platformColors answers the live set for the appearance: what AppKit
// reports right now for every name in the set, under aqua or darkAqua.
// The accent rows are the platform's own — the reader asks AppKit for
// controlAccentColor and the selection colours by name, so nothing on macOS
// derives them from the accent the shim reads out of the defaults, and
// tokens.PlatformColors.WithAccent is not on this path.
//
// It is deliberately independent of the palette precedence in the theme
// bridge: this set is the platform's answer, not the application's brand,
// so [WithSeed] and [WithPalette] do not reach it. It is equally
// independent of the Appearance's accent, which is the same setting read a
// second way; only Dark selects a side, so a stubbed Source still chooses
// which appearance is emitted.
func platformColors(a Appearance) tokens.PlatformColors {
	light, dark := livePlatform.sets()
	if a.Dark {
		return dark
	}
	return light
}

// livePlatform is the process's one cached reading of the two sets.
var livePlatform = &platformCache{
	interval: slowReadInterval,
	now:      time.Now,
	read:     liveSet,
}

// platformCache holds the two live sets and re-reads them at most once per
// interval, the cadence the accent key is read at and for the same reason:
// the colour set follows settings a user changes by hand, so a change must
// arrive within a poll or two, and asking AppKit for eighty-odd colours on
// every emission would pay for a reading nothing changed.
type platformCache struct {
	interval time.Duration
	now      func() time.Time                      // injectable clock for tests
	read     func(dark bool) tokens.PlatformColors // injectable reader for tests

	mu          sync.Mutex
	light, dark tokens.PlatformColors
	readAt      time.Time
	haveRead    bool
}

// sets returns the two cached sets, refreshing them when the interval has
// elapsed. The first call always reads, so the first emission is accurate.
func (c *platformCache) sets() (light, dark tokens.PlatformColors) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	if !c.haveRead || now.Sub(c.readAt) >= c.interval {
		c.light, c.dark = c.read(false), c.read(true)
		c.haveRead = true
		c.readAt = now
	}
	return c.light, c.dark
}
