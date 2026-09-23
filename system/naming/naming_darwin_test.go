package naming

import "testing"

// The Go order is the platform's order, measured rather than assumed: every
// pair of the fixture is compared both ways. The one measured divergence is
// the sharp s against a name spelled with "ss", which the package doc names
// and no fixture row carries.
func TestPlatformAndGoAgree(t *testing.T) {
	sign := func(n int) int {
		switch {
		case n < 0:
			return -1
		case n > 0:
			return 1
		}
		return 0
	}
	for _, a := range fixture {
		for _, b := range fixture {
			if p, g := sign(platformCompare(a, b)), sign(compareGo(a, b)); p != g {
				t.Errorf("%q against %q: the platform answers %d, the Go order %d", a, b, p, g)
			}
		}
	}
}
