package simd

import "sync"

var (
	forcePortable bool
	flagsMu       sync.Mutex
)

// ForcePortable, when true, rebinds all SIMD-dispatched functions to their
// portable Go implementations. Intended for tests/benchmarks.
func ForcePortable(v bool) {
	flagsMu.Lock()
	defer flagsMu.Unlock()
	forcePortable = v
	if v { rebindPortable() }
}

// rebindPortable assigns all impl* to pure Go fallbacks.
func rebindPortable() {
	implButterfly2 = Butterfly2
	implButterfly4 = Butterfly4
	implButterfly8 = Butterfly8
	implCmulAdd    = CmulAdd
	implCmul       = Cmul
	implScale      = Scale
	implButterfly4Twiddled = butterfly4TwiddledPortableUnroll2j
	implButterfly8Twiddled = butterfly8TwiddledPortable
}
