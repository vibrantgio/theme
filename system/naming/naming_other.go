//go:build !darwin

package naming

// platformCompare is [compareGo] where the platform has no comparison of its
// own to ask, which is everywhere but macOS. The order is the same order:
// that the two agree is what the macOS test measures.
func platformCompare(a, b string) int { return compareGo(a, b) }
