package gofft

import (
	"math/rand"
	"testing"
)

func benchFFT(b *testing.B, n int) {
	x := make([]complex128, n)
	for i := range x {
		x[i] = complex(rand.Float64(), rand.Float64())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FFT(x)
	}
}

func benchIFFT(b *testing.B, n int) {
	x := make([]complex128, n)
	for i := range x {
		x[i] = complex(rand.Float64(), rand.Float64())
	}
	// Pre-compute FFT for IFFT benchmark
	X, _ := FFT(x)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IFFT(X)
	}
}

func benchRoundTrip(b *testing.B, n int) {
	x := make([]complex128, n)
	for i := range x {
		x[i] = complex(rand.Float64(), rand.Float64())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		X, _ := FFT(x)
		_, _ = IFFT(X)
	}
}

// FFT Benchmarks
func BenchmarkFFT_64(b *testing.B)  { benchFFT(b, 64) }
func BenchmarkFFT_256(b *testing.B) { benchFFT(b, 256) }
func BenchmarkFFT_1k(b *testing.B)  { benchFFT(b, 1024) }
func BenchmarkFFT_4k(b *testing.B)  { benchFFT(b, 4096) }
func BenchmarkFFT_16k(b *testing.B) { benchFFT(b, 16384) }
func BenchmarkFFT_64k(b *testing.B) { benchFFT(b, 65536) }

// IFFT Benchmarks
func BenchmarkIFFT_64(b *testing.B)  { benchIFFT(b, 64) }
func BenchmarkIFFT_256(b *testing.B) { benchIFFT(b, 256) }
func BenchmarkIFFT_1k(b *testing.B)  { benchIFFT(b, 1024) }
func BenchmarkIFFT_4k(b *testing.B)  { benchIFFT(b, 4096) }
func BenchmarkIFFT_16k(b *testing.B) { benchIFFT(b, 16384) }
func BenchmarkIFFT_64k(b *testing.B) { benchIFFT(b, 65536) }

// Round Trip Benchmarks (FFT + IFFT)
func BenchmarkRoundTrip_64(b *testing.B)  { benchRoundTrip(b, 64) }
func BenchmarkRoundTrip_256(b *testing.B) { benchRoundTrip(b, 256) }
func BenchmarkRoundTrip_1k(b *testing.B)  { benchRoundTrip(b, 1024) }
func BenchmarkRoundTrip_4k(b *testing.B)  { benchRoundTrip(b, 4096) }
func BenchmarkRoundTrip_16k(b *testing.B) { benchRoundTrip(b, 16384) }
func BenchmarkRoundTrip_64k(b *testing.B) { benchRoundTrip(b, 65536) }
