package gofft

import (
	"github.com/example/gofft/pkg/gofft/simd"
	"github.com/example/gofft/pkg/gofft/threadpool"
)

// chunked iterates k in [start,end) with step m and runs fn(k).
func chunked(pool *threadpool.Pool, start, end, m int, fn func(int)) {
	for k := start; k < end; k += m {
		kk := k
		pool.Submit(func() { fn(kk) })
	}
	pool.Wait()
}

// parallelButterfly distributes a stage over the threadpool.
func parallelButterfly(pool *threadpool.Pool, radix int, buf []complex128, m int, invert bool) {
	n := len(buf)
	switch radix {
	case 8:
		w1, w2, w3, w4, w5, w6, w7 := precomputeTwiddlesRadix8(n, m, invert)
		chunked(pool, 0, n, m, func(k int) {
			simd.DispatchButterfly8Twiddled(buf, k, m, invert, w1, w2, w3, w4, w5, w6, w7)
		})
	case 4:
		w1, w2, w3 := precomputeTwiddlesRadix4(n, m, invert)
		chunked(pool, 0, n, m, func(k int) {
			simd.DispatchButterfly4Twiddled(buf, k, m, invert, w1, w2, w3)
		})
	case 2:
		w := precomputeTwiddlesRadix2(n, m, invert)
		chunked(pool, 0, n, m, func(k int) {
			butterfly2TwiddledAdapter(buf, k, m, invert, w)
		})
	default:
		panic("invalid radix")
	}
}
