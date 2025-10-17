package simd

// Portable baseline complex helpers and butterflies.
func CmulAdd(dst, a, b []complex128) []complex128 {
	n := min(len(dst), min(len(a), len(b)))
	for i := 0; i < n; i++ {
		dst[i] += a[i] * b[i]
	}
	return dst
}

func Cmul(dst, a, b []complex128) []complex128 {
	n := min(len(dst), min(len(a), len(b)))
	for i := 0; i < n; i++ {
		dst[i] = a[i] * b[i]
	}
	return dst
}

func Scale(dst, a []complex128, s complex128) []complex128 {
	n := min(len(dst), len(a))
	for i := 0; i < n; i++ {
		dst[i] = a[i] * s
	}
	return dst
}

// Basic radix-2/4/8 butterflies WITHOUT twiddles (used only as baseline fallback).
func Butterfly2(buf []complex128, k, m int, invert bool) {
	q := m >> 1
	for j := 0; j < q; j++ {
		a := buf[k+j]
		b := buf[k+j+q]
		buf[k+j] = a + b
		buf[k+j+q] = a - b
	}
}

func Butterfly4(buf []complex128, k, m int, invert bool) {
	q := m >> 2
	var W1i complex128
	if invert { W1i = complex(0, 1) } else { W1i = complex(0, -1) }
	for j := 0; j < q; j++ {
		a := buf[k+j]
		b := buf[k+j+q]
		c := buf[k+j+2*q]
		d := buf[k+j+3*q]
		t0 := a + c
		t1 := a - c
		t2 := b + d
		t3 := (b - d) * W1i
		buf[k+j] = t0 + t2
		buf[k+j+q] = t1 + t3
		buf[k+j+2*q] = t0 - t2
		buf[k+j+3*q] = t1 - t3
	}
}

func Butterfly8(buf []complex128, k, m int, invert bool) {
	o := m >> 3
	var W1i complex128
	if invert { W1i = complex(0, 1) } else { W1i = complex(0, -1) }
	for j := 0; j < o; j++ {
		x0 := buf[k+j]
		x1 := buf[k+j+o]
		x2 := buf[k+j+2*o]
		x3 := buf[k+j+3*o]
		x4 := buf[k+j+4*o]
		x5 := buf[k+j+5*o]
		x6 := buf[k+j+6*o]
		x7 := buf[k+j+7*o]
		// layer 1
		a0 := x0 + x4
		a4 := x0 - x4
		a1 := x1 + x5
		a5 := (x1 - x5) * W1i
		a2 := x2 + x6
		a6 := x2 - x6
		a3 := x3 + x7
		a7 := (x3 - x7) * W1i
		// layer 2
		b0 := a0 + a2
		b2 := a0 - a2
		b1 := a1 + a3
		b3 := (a1 - a3) * W1i
		b4 := a4 + a6
		b6 := (a4 - a6) * W1i
		b5 := a5 + a7
		b7 := (a5 - a7) * W1i
		// outputs
		buf[k+j]       = b0 + b1
		buf[k+j+o]     = b4 + b5
		buf[k+j+2*o]   = b2 + b3
		buf[k+j+3*o]   = b6 + b7
		buf[k+j+4*o]   = b0 - b1
		buf[k+j+5*o]   = b4 - b5
		buf[k+j+6*o]   = b2 - b3
		buf[k+j+7*o]   = b6 - b7
	}
}
