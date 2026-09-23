//go:build cgo

package openpanel

// builtWithCgo is whether this test binary can reach a platform framework
// through cgo. The panel is presented through one, so it is the other half of
// what [Available] must answer, beside the platform itself.
const builtWithCgo = true
