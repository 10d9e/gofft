package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestPrimitiveRootFinder tests the primitive root finding algorithm
func TestPrimitiveRootFinder(t *testing.T) {
	primes := []int{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}

	for _, p := range primes {
		t.Run("Prime"+string(rune(p+'0')), func(t *testing.T) {
			g := findPrimitiveRoot(p)

			t.Logf("Prime %d: primitive root = %d", p, g)

			// Verify it's actually a primitive root
			// g^k mod p should generate all values 1..(p-1)
			seen := make(map[int]bool)
			val := 1
			for k := 0; k < p-1; k++ {
				val = (val * g) % p
				seen[val] = true
			}

			// Should have seen all values 1..(p-1)
			if len(seen) != p-1 {
				t.Errorf("Primitive root %d for prime %d only generated %d unique values, expected %d",
					g, p, len(seen), p-1)
			}

			// Final value should be 1 (completing the cycle)
			if val != 1 {
				t.Errorf("Primitive root %d for prime %d: g^(p-1) mod p = %d, expected 1", g, p, val)
			}
		})
	}
}

// TestModularInverse tests modular inverse computation
func TestModularInverse(t *testing.T) {
	testCases := []struct {
		a, m, expected int
	}{
		{3, 7, 5},  // 3*5 = 15 ≡ 1 (mod 7)
		{2, 5, 3},  // 2*3 = 6 ≡ 1 (mod 5)
		{7, 11, 8}, // 7*8 = 56 ≡ 1 (mod 11)
	}

	for _, tc := range testCases {
		result := modInverse(tc.a, tc.m)

		// Verify: (a * result) mod m == 1
		product := (tc.a * result) % tc.m
		if product != 1 {
			t.Errorf("modInverse(%d, %d) = %d, but (%d * %d) mod %d = %d, expected 1",
				tc.a, tc.m, result, tc.a, result, tc.m, product)
		}

		t.Logf("modInverse(%d, %d) = %d ✓", tc.a, tc.m, result)
	}
}

// TestRadersSmallPrimes tests Rader's on small primes
func TestRadersSmallPrimes(t *testing.T) {
	primes := []int{3, 5, 7, 11, 13, 17, 19, 23, 29, 31}

	for _, p := range primes {
		t.Run("Prime"+string(rune(p+'0')), func(t *testing.T) {
			// Create Rader's with DFT inner FFT
			innerFft := NewDft(p-1, Forward)
			raders := NewRaders(innerFft)

			// Create test input
			input := make([]complex128, p)
			for i := range input {
				input[i] = complex(float64(i), float64(i)*0.3)
			}

			// Compute with Rader's
			result := make([]complex128, p)
			copy(result, input)
			scratch := make([]complex128, raders.InplaceScratchLen())
			raders.ProcessWithScratch(result, scratch)

			// Compute expected with DFT
			expected := make([]complex128, p)
			copy(expected, input)
			dft := NewDft(p, Forward)
			dft.ProcessWithScratch(expected, make([]complex128, p))

			// Compare
			maxErr := 0.0
			for i := range result {
				err := cmplx.Abs(result[i] - expected[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("Prime %d: Max error = %.6e", p, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("Prime %d failed with error %.6e", p, maxErr)
				for i := range result {
					if cmplx.Abs(result[i]-expected[i]) > 1e-10 {
						t.Errorf("  [%d] got=%v want=%v", i, result[i], expected[i])
					}
				}
			}
		})
	}
}

// TestRadersRoundTrip tests forward + inverse for Rader's
func TestRadersRoundTrip(t *testing.T) {
	primes := []int{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41}

	for _, p := range primes {
		t.Run("Prime"+string(rune(p+'0')), func(t *testing.T) {
			// Create forward and inverse Rader's
			fwdInner := NewDft(p-1, Forward)
			fwd := NewRaders(fwdInner)

			invInner := NewDft(p-1, Inverse)
			inv := NewRaders(invInner)

			// Test input
			input := make([]complex128, p)
			for i := range input {
				input[i] = complex(float64(i)*0.7, float64(i)*0.3)
			}
			original := make([]complex128, p)
			copy(original, input)

			// Forward
			fwdScratch := make([]complex128, fwd.InplaceScratchLen())
			fwd.ProcessWithScratch(input, fwdScratch)

			// Inverse
			invScratch := make([]complex128, inv.InplaceScratchLen())
			inv.ProcessWithScratch(input, invScratch)

			// Normalize
			for i := range input {
				input[i] /= complex(float64(p), 0)
			}

			// Compare
			maxErr := 0.0
			for i := range input {
				err := cmplx.Abs(input[i] - original[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("Prime %d: Round-trip error = %.6e", p, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("Prime %d round-trip failed with error %.6e", p, maxErr)
			}
		})
	}
}

// TestRadersVsBluestein compares Rader's and Bluestein's for accuracy
func TestRadersVsBluestein(t *testing.T) {
	primes := []int{11, 13, 17, 19, 23, 29, 31}

	for _, p := range primes {
		t.Run("Prime"+string(rune(p+'0')), func(t *testing.T) {
			// Test input
			input := make([]complex128, p)
			for i := range input {
				input[i] = complex(float64(i%7), float64(i%5)*0.3)
			}

			// Compute with Rader's
			raderInner := NewDft(p-1, Forward)
			raders := NewRaders(raderInner)
			raderResult := make([]complex128, p)
			copy(raderResult, input)
			raderScratch := make([]complex128, raders.InplaceScratchLen())
			raders.ProcessWithScratch(raderResult, raderScratch)

			// Compute with Bluestein's
			bluestein := NewBluestein(p, Forward)
			blueResult := make([]complex128, p)
			copy(blueResult, input)
			blueScratch := make([]complex128, bluestein.InplaceScratchLen())
			bluestein.ProcessWithScratch(blueResult, blueScratch)

			// Compare Rader's vs Bluestein's
			maxErr := 0.0
			for i := range raderResult {
				err := cmplx.Abs(raderResult[i] - blueResult[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("Prime %d: Rader's vs Bluestein's error = %.6e", p, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("Prime %d: Rader's and Bluestein's differ by %.6e", p, maxErr)
			}
		})
	}
}
