// The interaction states a component surface can be in.
package tokens

// State enumerates the interaction states a component surface can be in. It
// names the state and nothing else: what each one paints is the platform's
// answer, which the component reads off [PlatformColors] — the hover and the
// press overlays over the fill beneath, the selection rows for a selected
// one, the accent for the default action.
type State int

const (
	StateNormal State = iota
	StateHover
	StatePressed
	StateSelected
	StateDisabled
	StateFocus
	StateDragged
)
