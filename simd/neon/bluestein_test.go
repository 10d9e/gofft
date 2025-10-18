package neon

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestNEONBluestein(t *testing.T) {
	sizes := []int{128, 256, 512, 1024}

	for _, size := range sizes {
		t.Run("Bluestein_"+string(rune(size)), func(t *testing.T) {
			// Create test data
			data := make([]complex128, size)
			for i := 0; i < size; i++ {
				data[i] = complex128(complex(
					math.Sin(2*math.Pi*float64(i)/float64(size)),
					math.Cos(2*math.Pi*float64(i)/float64(size)),
				))
			}

			// Test NEON implementation
			ProcessVectorizedBluestein(data, size)

			// Verify that data was processed (not all zeros)
			hasNonZero := false
			for _, val := range data {
				if cmplx.Abs(val) > 1e-10 {
					hasNonZero = true
					break
				}
			}

			if !hasNonZero {
				t.Errorf("Bluestein_%d: All values are zero after processing", size)
			} else {
				t.Logf("Bluestein_%d: Processed %d elements successfully", size, size)
			}
		})
	}
}

func BenchmarkNEONBluestein(b *testing.B) {
	sizes := []int{128, 256, 512, 1024}

	for _, size := range sizes {
		b.Run("Bluestein_"+string(rune(size)), func(b *testing.B) {
			data := make([]complex128, size)
			for i := 0; i < size; i++ {
				data[i] = complex128(complex(
					math.Sin(2*math.Pi*float64(i)/float64(size)),
					math.Cos(2*math.Pi*float64(i)/float64(size)),
				))
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				ProcessVectorizedBluestein(data, size)
			}
		})
	}
}
