package gofft

import (
	"math"
	"math/cmplx"
	"testing"
)

// TestExtendedPowerOfTwoSizes tests larger power-of-two sizes
func TestExtendedPowerOfTwoSizes(t *testing.T) {
	sizes := []int{128, 256, 512, 1024, 2048, 4096}

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create test signal
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(math.Sin(float64(i)*0.1), math.Cos(float64(i)*0.1))
			}
			original := make([]complex128, n)
			copy(original, input)

			planner := NewPlanner()

			// Forward FFT
			forward := planner.PlanForward(n)
			forward.Process(input)

			// Inverse FFT
			inverse := planner.PlanInverse(n)
			inverse.Process(input)

			// Normalize
			for i := range input {
				input[i] /= complex(float64(n), 0)
			}

			// Check if we got back the original
			maxError := 0.0
			for i := range input {
				err := cmplx.Abs(input[i] - original[i])
				if err > maxError {
					maxError = err
				}
			}

			t.Logf("Size %d: Max reconstruction error: %.6e", n, maxError)

			if maxError > 1e-10 {
				t.Errorf("Forward + Inverse didn't recover original signal for size %d, error: %.6e", n, maxError)
			}
		})
	}
}

// BenchmarkFFTExtended benchmarks larger sizes
func BenchmarkFFTExtended(b *testing.B) {
	sizes := []int{1024, 4096, 16384, 65536}

	for _, n := range sizes {
		b.Run("Size"+string(rune(n+'0')), func(b *testing.B) {
			planner := NewPlanner()
			fft := planner.PlanForward(n)
			buffer := make([]complex128, n)
			scratch := make([]complex128, fft.InplaceScratchLen())

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				fft.ProcessWithScratch(buffer, scratch)
			}
		})
	}
}
