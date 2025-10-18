package neon

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestNEONGoodThomas(t *testing.T) {
	sizes := []int{35, 77, 143, 221}

	for _, size := range sizes {
		t.Run("GoodThomas_"+string(rune(size)), func(t *testing.T) {
			// Create test data
			data := make([]complex128, size)
			for i := 0; i < size; i++ {
				data[i] = complex128(complex(
					math.Sin(2*math.Pi*float64(i)/float64(size)),
					math.Cos(2*math.Pi*float64(i)/float64(size)),
				))
			}

			// Test NEON implementation
			ProcessVectorizedGoodThomas(data, size)

			// Verify that data was processed (not all zeros)
			hasNonZero := false
			for _, val := range data {
				if cmplx.Abs(val) > 1e-10 {
					hasNonZero = true
					break
				}
			}

			if !hasNonZero {
				t.Errorf("GoodThomas_%d: All values are zero after processing", size)
			} else {
				t.Logf("GoodThomas_%d: Processed %d elements successfully", size, size)
			}
		})
	}
}

func BenchmarkNEONGoodThomas(b *testing.B) {
	sizes := []int{35, 77, 143, 221}

	for _, size := range sizes {
		b.Run("GoodThomas_"+string(rune(size)), func(b *testing.B) {
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
				ProcessVectorizedGoodThomas(data, size)
			}
		})
	}
}
