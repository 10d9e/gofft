package gofft

import (
	"math"

	"github.com/example/gofft/pkg/gofft/simd"
)

// stageTwiddle returns exp(sign*i*2π*(step*j)/n) where sign = -1 for forward (invert=false) and +1 for inverse.
func stageTwiddle(step, j, n int, invert bool) complex128 {
	angle := 2 * math.Pi * float64(step*j) / float64(n)
	if !invert {
		angle = -angle
	}
	return complex(math.Cos(angle), math.Sin(angle))
}

// Radix-2 per-stage twiddles
func precomputeTwiddlesRadix2(n, m int, invert bool) []complex128 {
	q := m >> 1
	w := make([]complex128, q)
	step := n / m
	for j := 0; j < q; j++ {
		w[j] = stageTwiddle(step, j, n, invert)
	}
	return w
}

// Radix-4 per-stage twiddles
func precomputeTwiddlesRadix4(n, m int, invert bool) (w1, w2, w3 []complex128) {
	q := m >> 2
	w1 = make([]complex128, q)
	w2 = make([]complex128, q)
	w3 = make([]complex128, q)
	step := n / m
	for j := 0; j < q; j++ {
		w1[j] = stageTwiddle(step, j, n, invert)
		w2[j] = stageTwiddle(step, 2*j, n, invert)
		w3[j] = stageTwiddle(step, 3*j, n, invert)
	}
	return
}

// Radix-8 per-stage twiddles
func precomputeTwiddlesRadix8(n, m int, invert bool) (w1, w2, w3, w4, w5, w6, w7 []complex128) {
	o := m >> 3
	w1 = make([]complex128, o)
	w2 = make([]complex128, o)
	w3 = make([]complex128, o)
	w4 = make([]complex128, o)
	w5 = make([]complex128, o)
	w6 = make([]complex128, o)
	w7 = make([]complex128, o)
	step := n / m
	for j := 0; j < o; j++ {
		w1[j] = stageTwiddle(step, 1*j, n, invert)
		w2[j] = stageTwiddle(step, 2*j, n, invert)
		w3[j] = stageTwiddle(step, 3*j, n, invert)
		w4[j] = stageTwiddle(step, 4*j, n, invert)
		w5[j] = stageTwiddle(step, 5*j, n, invert)
		w6[j] = stageTwiddle(step, 6*j, n, invert)
		w7[j] = stageTwiddle(step, 7*j, n, invert)
	}
	return
}

// Sequential radix-2 stage (twiddle-aware)
func butterfly2InPlace(buf []complex128, m int, invert bool) {
	n := len(buf)
	w := precomputeTwiddlesRadix2(n, m, invert)
	for k := 0; k < n; k += m {
		butterfly2TwiddledAdapter(buf, k, m, invert, w)
	}
}

// small adapter to reach unexported portable function in simd
func butterfly2TwiddledAdapter(buf []complex128, k, m int, invert bool, w []complex128) {
	simdCallButterfly2Twiddled(buf, k, m, invert, w)
}

//go:noinline
func simdCallButterfly2Twiddled(buf []complex128, k, m int, invert bool, w []complex128) {
	q := m >> 1
	for j := 0; j < q; j++ {
		a := buf[k+j]
		b := w[j] * buf[k+j+q]
		buf[k+j] = a + b
		buf[k+j+q] = a - b
	}
}

// Sequential radix-4 stage (twiddle-aware)
func butterfly4InPlace(buf []complex128, m int, invert bool) {
	n := len(buf)
	w1, w2, w3 := precomputeTwiddlesRadix4(n, m, invert)
	for k := 0; k < n; k += m {
		simd.DispatchButterfly4Twiddled(buf, k, m, invert, w1, w2, w3)
	}
}

// Sequential radix-8 stage (twiddle-aware)
func butterfly8InPlace(buf []complex128, m int, invert bool) {
	n := len(buf)
	w1, w2, w3, w4, w5, w6, w7 := precomputeTwiddlesRadix8(n, m, invert)
	for k := 0; k < n; k += m {
		simd.DispatchButterfly8Twiddled(buf, k, m, invert, w1, w2, w3, w4, w5, w6, w7)
	}
}
