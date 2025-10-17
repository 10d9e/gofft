package simd

// Twiddle-aware radix-2 butterfly over a single block: m=2*q, w1[j]=W^(j).
// buf indices: a=buf[k+j], b=buf[k+j+q] (b is multiplied by w1[j] before butterfly).
func butterfly2TwiddledPortable(buf []complex128, k, m int, invert bool, w1 []complex128) {
	q := m >> 1
	for j := 0; j < q; j++ {
		a := buf[k+j]
		b := w1[j] * buf[k+j+q]
		buf[k+j] = a + b
		buf[k+j+q] = a - b
	}
}
