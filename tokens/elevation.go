// Shadow depth: what a floating surface casts, in device-independent pixels.
package tokens

// ElevationScale carries each level's shadow depth in device-independent
// pixels, and nothing else. A level's FILL is the platform's own — the
// window's plane, the chrome material, the content, the card — read off
// [PlatformColors] by name; the shadow is the separate, opt-in cue a
// floating surface casts, which effects/depth renders.
//
// The six levels run from the backdrop up, anchored on the content: the
// backdrop is the bare window plane, the chrome the sidebars and toolbars,
// Level0 the content, Level1 a raised inset on it, and Level2 and Level3 the
// two floating levels — a dialog and a menu. Neither level under the content
// casts anything: the backdrop is what everything stands on, and chrome lies
// flat on it.
type ElevationScale struct {
	Backdrop float32 // 0 dp
	Chrome   float32 // 0 dp
	Level0   float32 // 0 dp
	Level1   float32 // 1 dp
	Level2   float32 // 3 dp
	Level3   float32 // 6 dp
}

// Elevation is the default scale instance.
var Elevation = ElevationScale{
	Backdrop: 0,
	Chrome:   0,
	Level0:   0,
	Level1:   1,
	Level2:   3,
	Level3:   6,
}
