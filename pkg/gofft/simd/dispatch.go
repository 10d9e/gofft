package simd

// Core function pointers used by gofft.

var (
	implButterfly2 func(buf []complex128, k, m int, invert bool)
	implButterfly4 func(buf []complex128, k, m int, invert bool)
	implButterfly8 func(buf []complex128, k, m int, invert bool)

	implCmulAdd func(dst, a, b []complex128) []complex128
	implCmul    func(dst, a, b []complex128) []complex128
	implScale   func(dst, a []complex128, s complex128) []complex128

	// Twiddle-aware seams (portable today; SIMD override when available)
	implButterfly4Twiddled func(buf []complex128, k, m int, invert bool, w1, w2, w3 []complex128)
	implButterfly8Twiddled func(buf []complex128, k, m int, invert bool, w1, w2, w3, w4, w5, w6, w7 []complex128)
)

func init() {
	implButterfly2 = Butterfly2
	implButterfly4 = Butterfly4
	implButterfly8 = Butterfly8

	implCmulAdd = CmulAdd
	implCmul    = Cmul
	implScale   = Scale

	// Default portable paths (now with 2j unroll for radix-4)
	implButterfly4Twiddled = butterfly4TwiddledPortableUnroll2j
	implButterfly8Twiddled = butterfly8TwiddledPortable
}

// Public dispatch
func DispatchButterfly2(buf []complex128, k, m int, invert bool) { implButterfly2(buf, k, m, invert) }
func DispatchButterfly4(buf []complex128, k, m int, invert bool) { implButterfly4(buf, k, m, invert) }
func DispatchButterfly8(buf []complex128, k, m int, invert bool) { implButterfly8(buf, k, m, invert) }

// Twiddle-aware public seams
func DispatchButterfly4Twiddled(buf []complex128, k, m int, invert bool, w1, w2, w3 []complex128) {
	implButterfly4Twiddled(buf, k, m, invert, w1, w2, w3)
}
func DispatchButterfly8Twiddled(
	buf []complex128, k, m int, invert bool,
	w1, w2, w3, w4, w5, w6, w7 []complex128,
) {
	implButterfly8Twiddled(buf, k, m, invert, w1, w2, w3, w4, w5, w6, w7)
}
