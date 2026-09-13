// The contrast floors a foreground is gated at.
package tokens

// The floors a foreground is gated at, named for the consumers that draw
// with them. Each is a floor on |Lc| — color.Magnitude — and each is APCA's
// own published level for the job, so a consumer that imports one and a test
// that gates on it are reading the same number: color.APCA is the system's
// one contrast measure and these are the levels it is read against.
const (
	// TextFloor is APCA's minimum for body text, Lc 75 — what a link, a
	// label or any run of words owes the surface it is set on.
	TextFloor = 75.0

	// GraphicFloor is APCA's minimum for a non-text graphic that has to be
	// resolved — an icon, a check, a rule, a control's edge: Lc 45. A
	// shape carries its meaning in its outline as well as its value, so it
	// is held lower than text and never lower than this.
	GraphicFloor = 45.0
)
