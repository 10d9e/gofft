package gofft

import (
	"math/cmplx"
	"testing"
)

// TestRadixNIntegration tests RadixN through the planner
func TestRadixNIntegration(t *testing.T) {
	// Sizes that should use RadixN (only factors 2-7)
	sizes := []int{
		6,   // 2×3
		10,  // 2×5
		12,  // 2²×3
		14,  // 2×7
		15,  // 3×5
		18,  // 2×3²
		20,  // 2²×5
		21,  // 3×7
		24,  // 2³×3
		28,  // 2²×7
		30,  // 2×3×5
		36,  // 2²×3²
		40,  // 2³×5
		42,  // 2×3×7
		48,  // 2⁴×3
		54,  // 2×3³
		56,  // 2³×7
		60,  // 2²×3×5
		63,  // 3²×7
		72,  // 2³×3²
		80,  // 2⁴×5
		84,  // 2²×3×7
		90,  // 2×3²×5
		96,  // 2⁵×3
		100, // 2²×5²
		120, // 2³×3×5
	}

	planner := NewPlanner()

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i%7), float64(i%5)*0.3)
			}
			original := make([]complex128, n)
			copy(original, input)

			// Forward FFT
			fft := planner.PlanForward(n)
			fft.Process(input)

			// Inverse FFT
			ifft := planner.PlanInverse(n)
			ifft.Process(input)

			// Normalize
			for i := range input {
				input[i] /= complex(float64(n), 0)
			}

			// Check round-trip
			maxErr := 0.0
			for i := range input {
				err := cmplx.Abs(input[i] - original[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("Size %d: Round-trip error = %.6e (using RadixN)", n, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("Size %d failed with error %.6e", n, maxErr)
			}
		})
	}
}
