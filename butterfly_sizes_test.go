package gofft

import (
	"math"
	"math/cmplx"
	"testing"
)

// TestAllButterflySizes tests all implemented butterfly sizes
func TestAllButterflySizes(t *testing.T) {
	sizes := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15, 16, 17, 19, 23, 24, 27, 29, 31, 32}

	planner := NewPlanner()

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create test signal
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i), float64(i)*0.5)
			}

			// Compute expected output using naive DFT
			expected := naiveDFT(input, true)

			// Compute using our FFT
			buffer := make([]complex128, n)
			copy(buffer, input)

			fft := planner.PlanForward(n)
			fft.Process(buffer)

			// Compare results
			maxError := 0.0
			for i := range buffer {
				err := cmplx.Abs(buffer[i] - expected[i])
				if err > maxError {
					maxError = err
				}
			}

			t.Logf("Size %d: Max error = %.6e", n, maxError)

			if maxError > 1e-10 {
				t.Errorf("FFT output doesn't match expected for size %d, error %.6e", n, maxError)
				for i := range buffer {
					if cmplx.Abs(buffer[i]-expected[i]) > 1e-10 {
						t.Errorf("  [%d] got %v, want %v", i, buffer[i], expected[i])
					}
				}
			}
		})
	}
}

// TestButterflySizesRoundTrip tests forward + inverse for all butterfly sizes
func TestButterflySizesRoundTrip(t *testing.T) {
	sizes := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15, 16, 17, 19, 23, 24, 27, 29, 31, 32}

	planner := NewPlanner()

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create test signal
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(math.Sin(float64(i)), math.Cos(float64(i)))
			}
			original := make([]complex128, n)
			copy(original, input)

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

			t.Logf("Size %d: Round-trip error = %.6e", n, maxError)

			if maxError > 1e-10 {
				t.Errorf("Forward + Inverse didn't recover original signal for size %d", n)
			}
		})
	}
}
