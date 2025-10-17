package gofft

import (
	"math/cmplx"
	"testing"
)

// TestCompositeSizes tests various composite FFT sizes
func TestCompositeSizes(t *testing.T) {
	sizes := []int{
		// Now supported with butterflies
		5, 6, 7, 9, 11, 12, 13, 24, 27,
		// Composite sizes that will need RadixN/MixedRadix
		10, 14, 15, 18, 20, 21, 24, 28, 30,
		36, 40, 48, 60, 72, 96,
	}

	planner := NewPlanner()

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i%7), float64(i%5)*0.3)
			}
			original := make([]complex128, n)
			copy(original, input)

			// Forward
			forward := planner.PlanForward(n)
			forward.Process(input)

			// Inverse
			inverse := planner.PlanInverse(n)
			inverse.Process(input)

			// Normalize
			for i := range input {
				input[i] /= complex(float64(n), 0)
			}

			// Check
			maxError := 0.0
			for i := range input {
				err := cmplx.Abs(input[i] - original[i])
				if err > maxError {
					maxError = err
				}
			}

			t.Logf("Size %d: Round-trip error = %.6e", n, maxError)

			if maxError > 1e-9 {
				t.Errorf("Size %d failed with error %.6e", n, maxError)
			}
		})
	}
}
