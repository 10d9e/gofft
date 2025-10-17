package simd

// Portable, twiddle-aware radix-4 butterfly over a single block (k..k+m-1).
// This version is unrolled to process TWO j's per inner iteration (2j).
func butterfly4TwiddledPortableUnroll2j(buf []complex128, k, m int, invert bool, w1, w2, w3 []complex128) {
	q := m >> 2
	var W1i complex128
	if invert {
		W1i = complex(0, 1)
	} else {
		W1i = complex(0, -1)
	}

	j := 0
	for ; j+1 < q; j += 2 {
		// j
		a0 := buf[k+j]
		b0 := w1[j] * buf[k+j+q]
		c0 := w2[j] * buf[k+j+2*q]
		d0 := w3[j] * buf[k+j+3*q]

		t00 := a0 + c0
		t10 := a0 - c0
		t20 := b0 + d0
		t30 := (b0 - d0) * W1i

		buf[k+j] = t00 + t20
		buf[k+j+q] = t10 + t30
		buf[k+j+2*q] = t00 - t20
		buf[k+j+3*q] = t10 - t30

		// j+1
		j1 := j + 1
		a1 := buf[k+j1]
		b1 := w1[j1] * buf[k+j1+q]
		c1 := w2[j1] * buf[k+j1+2*q]
		d1 := w3[j1] * buf[k+j1+3*q]

		t01 := a1 + c1
		t11 := a1 - c1
		t21 := b1 + d1
		t31 := (b1 - d1) * W1i

		buf[k+j1] = t01 + t21
		buf[k+j1+q] = t11 + t31
		buf[k+j1+2*q] = t01 - t21
		buf[k+j1+3*q] = t11 - t31
	}
	if j < q {
		a := buf[k+j]
		b := w1[j] * buf[k+j+q]
		c := w2[j] * buf[k+j+2*q]
		d := w3[j] * buf[k+j+3*q]
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

// Portable, twiddle-aware radix-8 butterfly over a single block.
func butterfly8TwiddledPortable(
	buf []complex128, k, m int, invert bool,
	w1, w2, w3, w4, w5, w6, w7 []complex128,
) {
	o := m >> 3
	var W1i complex128
	if invert {
		W1i = complex(0, 1)
	} else {
		W1i = complex(0, -1)
	}

	for j := 0; j < o; j++ {
		x0 := buf[k+j]
		x1 := w1[j] * buf[k+j+o]
		x2 := w2[j] * buf[k+j+2*o]
		x3 := w3[j] * buf[k+j+3*o]
		x4 := w4[j] * buf[k+j+4*o]
		x5 := w5[j] * buf[k+j+5*o]
		x6 := w6[j] * buf[k+j+6*o]
		x7 := w7[j] * buf[k+j+7*o]

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
		buf[k+j] = b0 + b1
		buf[k+j+o] = b4 + b5
		buf[k+j+2*o] = b2 + b3
		buf[k+j+3*o] = b6 + b7
		buf[k+j+4*o] = b0 - b1
		buf[k+j+5*o] = b4 - b5
		buf[k+j+6*o] = b2 - b3
		buf[k+j+7*o] = b6 - b7
	}
}
