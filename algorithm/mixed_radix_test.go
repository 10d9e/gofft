package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestMixedRadixCorrectness tests MixedRadix against DFT
func TestMixedRadixCorrectness(t *testing.T) {
	testCases := []struct {
		n1, n2 int
	}{
		{2, 3}, // 6
		{2, 5}, // 10
		{3, 4}, // 12
		{3, 5}, // 15
		{4, 5}, // 20
		{5, 5}, // 25
		{6, 6}, // 36
		{8, 8}, // 64
	}

	for _, tc := range testCases {
		n := tc.n1 * tc.n2
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create MixedRadix
			width := NewDft(tc.n1, Forward)
			height := NewDft(tc.n2, Forward)
			mr := NewMixedRadix(width, height)

			// Test input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i), float64(i)*0.5)
			}

			// Compute with MixedRadix
			result := make([]complex128, n)
			copy(result, input)
			scratch := make([]complex128, mr.InplaceScratchLen())
			mr.ProcessWithScratch(result, scratch)

			// Compute expected with DFT
			expected := make([]complex128, n)
			copy(expected, input)
			dft := NewDft(n, Forward)
			dft.ProcessWithScratch(expected, make([]complex128, n))

			// Compare
			maxErr := 0.0
			for i := range result {
				err := cmplx.Abs(result[i] - expected[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("Size %d (%dx%d): Max error = %.6e", n, tc.n1, tc.n2, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("Size %d failed with error %.6e", n, maxErr)
				for i := range result {
					if cmplx.Abs(result[i]-expected[i]) > 1e-10 {
						t.Errorf("  [%d] got=%v want=%v", i, result[i], expected[i])
					}
				}
			}
		})
	}
}

// TestMixedRadixRoundTrip tests forward + inverse
func TestMixedRadixRoundTrip(t *testing.T) {
	testCases := []struct {
		n1, n2 int
	}{
		{2, 3}, // 6
		{2, 5}, // 10
		{3, 4}, // 12
		{3, 5}, // 15
		{4, 5}, // 20
		{5, 7}, // 35
	}

	for _, tc := range testCases {
		n := tc.n1 * tc.n2
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create forward and inverse
			widthFwd := NewDft(tc.n1, Forward)
			heightFwd := NewDft(tc.n2, Forward)
			fwd := NewMixedRadix(widthFwd, heightFwd)

			widthInv := NewDft(tc.n1, Inverse)
			heightInv := NewDft(tc.n2, Inverse)
			inv := NewMixedRadix(widthInv, heightInv)

			// Test input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i)*0.7, float64(i)*0.3)
			}
			original := make([]complex128, n)
			copy(original, input)

			// Forward
			scratchFwd := make([]complex128, fwd.InplaceScratchLen())
			fwd.ProcessWithScratch(input, scratchFwd)

			// Inverse
			scratchInv := make([]complex128, inv.InplaceScratchLen())
			inv.ProcessWithScratch(input, scratchInv)

			// Normalize
			for i := range input {
				input[i] /= complex(float64(n), 0)
			}

			// Check
			maxErr := 0.0
			for i := range input {
				err := cmplx.Abs(input[i] - original[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("Size %d (%dx%d): Round-trip error = %.6e", n, tc.n1, tc.n2, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("Size %d round-trip failed with error %.6e", n, maxErr)
			}
		})
	}
}

// TestMixedRadixWithButterflies tests MixedRadix using actual butterflies
func TestMixedRadixWithButterflies(t *testing.T) {
	testCases := []struct {
		n1, n2 int
		name   string
	}{
		{3, 4, "12=3x4"},
		{4, 3, "12=4x3"},
		{3, 5, "15=3x5"},
		{5, 3, "15=5x3"},
		{4, 5, "20=4x5"},
		{5, 4, "20=5x4"},
	}

	for _, tc := range testCases {
		n := tc.n1 * tc.n2
		t.Run(tc.name, func(t *testing.T) {
			// Create MixedRadix with actual butterflies
			var width, height FftInterface
			switch tc.n1 {
			case 2:
				width = NewButterfly2(Forward)
			case 3:
				width = NewButterfly3(Forward)
			case 4:
				width = NewButterfly4(Forward)
			case 5:
				width = NewButterfly5(Forward)
			default:
				width = NewDft(tc.n1, Forward)
			}

			switch tc.n2 {
			case 2:
				height = NewButterfly2(Forward)
			case 3:
				height = NewButterfly3(Forward)
			case 4:
				height = NewButterfly4(Forward)
			case 5:
				height = NewButterfly5(Forward)
			default:
				height = NewDft(tc.n2, Forward)
			}

			mr := NewMixedRadix(width, height)

			// Test input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i), float64((i*3)%7)*0.2)
			}

			// Compute with MixedRadix
			result := make([]complex128, n)
			copy(result, input)
			scratch := make([]complex128, mr.InplaceScratchLen())
			mr.ProcessWithScratch(result, scratch)

			// Compute expected with DFT
			expected := make([]complex128, n)
			copy(expected, input)
			dft := NewDft(n, Forward)
			dft.ProcessWithScratch(expected, make([]complex128, n))

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
			}
		})
	}
}
