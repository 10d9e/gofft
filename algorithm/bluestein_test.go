package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestBluesteinPrimes tests Bluestein's algorithm on prime sizes
func TestBluesteinPrimes(t *testing.T) {
	primes := []int{11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53}

	for _, n := range primes {
		t.Run("Prime"+string(rune(n+'0')), func(t *testing.T) {
			// Create input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i%7), float64(i%5)*0.3)
			}

			// Compute with Bluestein
			bluestein := NewBluestein(n, Forward)
			result := make([]complex128, n)
			copy(result, input)
			scratch := make([]complex128, bluestein.InplaceScratchLen())
			bluestein.ProcessWithScratch(result, scratch)

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

			t.Logf("Prime %d: Max error = %.6e", n, maxErr)

			if maxErr > 1e-9 {
				t.Errorf("Prime %d failed with error %.6e", n, maxErr)
			}
		})
	}
}

// TestBluesteinArbitrary tests Bluestein on arbitrary sizes
func TestBluesteinArbitrary(t *testing.T) {
	sizes := []int{100, 127, 200, 255, 300, 500, 1000}

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i%11)*0.7, float64(i%7)*0.3)
			}

			// Compute with Bluestein
			bluestein := NewBluestein(n, Forward)
			result := make([]complex128, n)
			copy(result, input)
			scratch := make([]complex128, bluestein.InplaceScratchLen())
			bluestein.ProcessWithScratch(result, scratch)

			// Compute expected with DFT (only for small sizes)
			if n <= 300 {
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

				t.Logf("Size %d: Max error = %.6e", n, maxErr)

				if maxErr > 1e-9 {
					t.Errorf("Size %d failed with error %.6e", n, maxErr)
				}
			} else {
				// For large sizes, just check round-trip
				invBluestein := NewBluestein(n, Inverse)
				invScratch := make([]complex128, invBluestein.InplaceScratchLen())
				invBluestein.ProcessWithScratch(result, invScratch)

				// Normalize
				for i := range result {
					result[i] /= complex(float64(n), 0)
				}

				// Compare
				maxErr := 0.0
				for i := range result {
					err := cmplx.Abs(result[i] - input[i])
					if err > maxErr {
						maxErr = err
					}
				}

				t.Logf("Size %d: Round-trip error = %.6e", n, maxErr)

				if maxErr > 1e-8 {
					t.Errorf("Size %d round-trip failed with error %.6e", n, maxErr)
				}
			}
		})
	}
}

// TestBluesteinRoundTrip tests forward + inverse
func TestBluesteinRoundTrip(t *testing.T) {
	sizes := []int{11, 13, 17, 23, 100, 127, 200}

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i)*0.7, float64(i)*0.3)
			}
			original := make([]complex128, n)
			copy(original, input)

			// Forward
			fwd := NewBluestein(n, Forward)
			fwdScratch := make([]complex128, fwd.InplaceScratchLen())
			fwd.ProcessWithScratch(input, fwdScratch)

			// Inverse
			inv := NewBluestein(n, Inverse)
			invScratch := make([]complex128, inv.InplaceScratchLen())
			inv.ProcessWithScratch(input, invScratch)

			// Normalize
			for i := range input {
				input[i] /= complex(float64(n), 0)
			}

			// Compare
			maxErr := 0.0
			for i := range input {
				err := cmplx.Abs(input[i] - original[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("Size %d: Round-trip error = %.6e", n, maxErr)

			if maxErr > 1e-9 {
				t.Errorf("Size %d round-trip failed with error %.6e", n, maxErr)
			}
		})
	}
}
