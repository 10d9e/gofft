package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestRadixNSimple tests RadixN with simple factorizations
func TestRadixNSimple(t *testing.T) {
	testCases := []struct {
		factors []RadixFactor
		size    int
		name    string
	}{
		{[]RadixFactor{Factor2, Factor2}, 4, "4=2×2"},
		{[]RadixFactor{Factor2, Factor3}, 6, "6=2×3"},
		{[]RadixFactor{Factor2, Factor2, Factor2}, 8, "8=2³"},
		{[]RadixFactor{Factor3, Factor3}, 9, "9=3²"},
		{[]RadixFactor{Factor2, Factor2, Factor3}, 12, "12=2²×3"},
		{[]RadixFactor{Factor3, Factor5}, 15, "15=3×5"},
		{[]RadixFactor{Factor2, Factor2, Factor2, Factor3}, 24, "24=2³×3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create RadixN with DFT base (size 1)
			baseFft := NewDft(1, Forward)
			radixN := NewRadixN(tc.factors, baseFft)

			if radixN.Len() != tc.size {
				t.Fatalf("Expected size %d, got %d", tc.size, radixN.Len())
			}

			// Test input
			input := make([]complex128, tc.size)
			for i := range input {
				input[i] = complex(float64(i), float64(i)*0.3)
			}

			// Compute with RadixN
			result := make([]complex128, tc.size)
			copy(result, input)
			scratch := make([]complex128, radixN.InplaceScratchLen())
			radixN.ProcessWithScratch(result, scratch)

			// Compute expected with DFT
			expected := make([]complex128, tc.size)
			copy(expected, input)
			dft := NewDft(tc.size, Forward)
			dft.ProcessWithScratch(expected, make([]complex128, tc.size))

			// Compare
			maxErr := 0.0
			for i := range result {
				err := cmplx.Abs(result[i] - expected[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("%s: Max error = %.6e", tc.name, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("%s failed with error %.6e", tc.name, maxErr)
				for i := range result {
					if cmplx.Abs(result[i]-expected[i]) > 1e-10 {
						t.Errorf("  [%d] got=%v want=%v", i, result[i], expected[i])
					}
				}
			}
		})
	}
}

// TestRadixNRoundTrip tests forward + inverse
func TestRadixNRoundTrip(t *testing.T) {
	testCases := []struct {
		factors []RadixFactor
		size    int
		name    string
	}{
		{[]RadixFactor{Factor2, Factor3}, 6, "6=2×3"},
		{[]RadixFactor{Factor2, Factor2, Factor3}, 12, "12=2²×3"},
		{[]RadixFactor{Factor2, Factor2, Factor5}, 20, "20=2²×5"},
		{[]RadixFactor{Factor2, Factor3, Factor5}, 30, "30=2×3×5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create forward and inverse
			baseFwd := NewDft(1, Forward)
			fwd := NewRadixN(tc.factors, baseFwd)

			baseInv := NewDft(1, Inverse)
			inv := NewRadixN(tc.factors, baseInv)

			// Test input
			input := make([]complex128, tc.size)
			for i := range input {
				input[i] = complex(float64(i)*0.7, float64(i)*0.3)
			}
			original := make([]complex128, tc.size)
			copy(original, input)

			// Forward
			fwdScratch := make([]complex128, fwd.InplaceScratchLen())
			fwd.ProcessWithScratch(input, fwdScratch)

			// Inverse
			invScratch := make([]complex128, inv.InplaceScratchLen())
			inv.ProcessWithScratch(input, invScratch)

			// Normalize
			for i := range input {
				input[i] /= complex(float64(tc.size), 0)
			}

			// Compare
			maxErr := 0.0
			for i := range input {
				err := cmplx.Abs(input[i] - original[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("%s: Round-trip error = %.6e", tc.name, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("%s round-trip failed with error %.6e", tc.name, maxErr)
			}
		})
	}
}

// TestRemainderReversal tests the remainder reversal function
func TestRemainderReversal(t *testing.T) {
	// Test with factors [2, 2, 3] (size 12)
	// This should reverse the "digits" in base [2, 2, 3]
	factors := []TransposeFactor{
		{factor: Factor2, count: 2}, // Two factors of 2
		{factor: Factor3, count: 1}, // One factor of 3
	}

	testCases := []struct {
		input, expected int
	}{
		{0, 0},   // 0 → 0
		{1, 6},   // 001₂₂₃ → 100₃₂₂ = 6
		{2, 3},   // 010₂₂₃ → 010₃₂₂ = 3
		{3, 9},   // 011₂₂₃ → 110₃₂₂ = 9
		{4, 1},   // 100₂₂₃ → 001₃₂₂ = 1
		{5, 7},   // 101₂₂₃ → 101₃₂₂ = 7
		{6, 4},   // 110₂₂₃ → 011₃₂₂ = 4
		{7, 10},  // 111₂₂₃ → 111₃₂₂ = 10
		{8, 2},   // 200₂₂₃ → 002₃₂₂ = 2
		{9, 8},   // 201₂₂₃ → 102₃₂₂ = 8
		{10, 5},  // 210₂₂₃ → 012₃₂₂ = 5
		{11, 11}, // 211₂₂₃ → 112₃₂₂ = 11
	}

	for _, tc := range testCases {
		result := reverseRemainders(tc.input, factors)
		if result != tc.expected {
			t.Errorf("reverseRemainders(%d) = %d, expected %d", tc.input, result, tc.expected)
		} else {
			t.Logf("reverseRemainders(%d) = %d ✓", tc.input, result)
		}
	}
}
