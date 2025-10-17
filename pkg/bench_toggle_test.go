package gofft

import (
	"math/rand"
	"testing"
	"github.com/example/gofft/pkg/gofft/simd"
)

func benchToggle(b *testing.B, n int, portable bool) {
	x := make([]complex128, n)
	for i := range x { x[i] = complex(rand.Float64(), rand.Float64()) }
	simd.ForcePortable(portable)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FFT(x)
	}
}

func BenchmarkFFT_SIMD_4k(b *testing.B)      { benchToggle(b, 4096, false) }
func BenchmarkFFT_Portable_4k(b *testing.B)  { benchToggle(b, 4096, true) }
func BenchmarkFFT_SIMD_16k(b *testing.B)     { benchToggle(b, 16384, false) }
func BenchmarkFFT_Portable_16k(b *testing.B) { benchToggle(b, 16384, true) }
