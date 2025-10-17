package gofft

import (
	"math/rand"
	"testing"
)

func benchFFT(b *testing.B, n int) {
	x := make([]complex128, n)
	for i := range x { x[i] = complex(rand.Float64(), rand.Float64()) }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FFT(x)
	}
}

func BenchmarkFFT_1k(b *testing.B)   { benchFFT(b, 1024) }
func BenchmarkFFT_4k(b *testing.B)   { benchFFT(b, 4096) }
func BenchmarkFFT_16k(b *testing.B)  { benchFFT(b, 16384) }
