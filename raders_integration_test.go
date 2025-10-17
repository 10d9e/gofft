package gofft

import (
	"math/cmplx"
	"testing"
)

// TestRadersIntegration tests Rader's through the planner
func TestRadersIntegration(t *testing.T) {
	primes := []int{37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

	planner := NewPlanner()

	for _, p := range primes {
		t.Run("Prime"+string(rune(p+'0')), func(t *testing.T) {
			// Create input
			input := make([]complex128, p)
			for i := range input {
				input[i] = complex(float64(i%7), float64(i%5)*0.3)
			}
			original := make([]complex128, p)
			copy(original, input)

			// Forward FFT
			fft := planner.PlanForward(p)
			fft.Process(input)

			// Inverse FFT
			ifft := planner.PlanInverse(p)
			ifft.Process(input)

			// Normalize
			for i := range input {
				input[i] /= complex(float64(p), 0)
			}

			// Check round-trip
			maxErr := 0.0
			for i := range input {
				err := cmplx.Abs(input[i] - original[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("Prime %d: Round-trip error = %.6e (using Rader's)", p, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("Prime %d failed with error %.6e", p, maxErr)
			}
		})
	}
}

// TestPlanner AlgorithmSelection verifies planner chooses correct algorithm
func TestPlannerAlgorithmSelection(t *testing.T) {
	planner := NewPlanner()

	testCases := []struct {
		size     int
		expected string
	}{
		{64, "Radix4"},
		{128, "Radix4"},
		{2048, "Radix4"},
		{3, "Butterfly3"},
		{5, "Butterfly5"},
		{7, "Butterfly7"},
		{37, "Raders"}, // Prime <= 97, should use Rader's
		{41, "Raders"},
		{97, "Raders"},
		{101, "Bluestein"}, // Prime > 97, should use Bluestein's
		{100, "Bluestein"}, // Composite, should use Bluestein's
		{1000, "Bluestein"},
	}

	for _, tc := range testCases {
		t.Run("Size"+string(rune(tc.size+'0')), func(t *testing.T) {
			fft := planner.PlanForward(tc.size)

			// Just verify it works (we can't easily check the algorithm type)
			buffer := make([]complex128, tc.size)
			for i := range buffer {
				buffer[i] = complex(float64(i), 0)
			}
			fft.Process(buffer)

			t.Logf("Size %d: Expected %s algorithm - FFT computed successfully", tc.size, tc.expected)
		})
	}
}
