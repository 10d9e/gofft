//go:build arm64

package simd

import "golang.org/x/sys/cpu"

//go:noescape
func cmuladdNEON(dst, a, b *complex128, n int)

//go:noescape
func cmuladdFMLA(dst, a, b *complex128, n int)

//go:noescape
func butterfly4TwiddledNEON(buf, w1, w2, w3 *complex128, k, m, n, q int, invert bool)

//go:noescape
func butterfly8TwiddledNEON(buf, w1, w2, w3, w4, w5, w6, w7 *complex128, k, m, n, o int, invert bool)

func init() {
	if cpu.ARM64.HasASIMD {
		implCmulAdd = func(dst, a, b []complex128) []complex128 {
			n := min(len(dst), min(len(a), len(b)))
			if n == 0 {
				return dst
			}
			cmuladdFMLA(&dst[0], &a[0], &b[0], n)
			return dst
		}
		implButterfly4Twiddled = func(buf []complex128, k, m int, invert bool, w1, w2, w3 []complex128) {
			if len(buf) == 0 {
				return
			}
			q := m >> 2
			butterfly4TwiddledNEON(&buf[0], &w1[0], &w2[0], &w3[0], k, m, len(buf), q, invert)
		}
		implButterfly8Twiddled = func(buf []complex128, k, m int, invert bool, w1, w2, w3, w4, w5, w6, w7 []complex128) {
			// portable for now
			butterfly8TwiddledPortable(buf, k, m, invert, w1, w2, w3, w4, w5, w6, w7)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
