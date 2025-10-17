package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestPrimeButterflies tests all prime-sized butterflies
func TestPrimeButterflies(t *testing.T) {
	primes := []int{3, 5, 7, 11, 13, 17, 19, 23, 29, 31}

	for _, n := range primes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create butterfly
			var bf FftInterface
			switch n {
			case 3:
				bf = NewButterfly3(Forward)
			case 5:
				bf = NewButterfly5(Forward)
			case 7:
				bf = NewButterfly7(Forward)
			case 11:
				bf = NewButterfly11(Forward)
			case 13:
				bf = NewButterfly13(Forward)
			case 17:
				bf = NewButterfly17(Forward)
			case 19:
				bf = NewButterfly19(Forward)
			case 23:
				bf = NewButterfly23(Forward)
			case 29:
				bf = NewButterfly29(Forward)
			case 31:
				bf = NewButterfly31(Forward)
			default:
				t.Fatalf("Unknown prime size %d", n)
			}

			// Test input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i), float64(i)*0.3)
			}

			// Compute with butterfly
			buffer := make([]complex128, n)
			copy(buffer, input)
			scratch := make([]complex128, bf.InplaceScratchLen())
			bf.ProcessWithScratch(buffer, scratch)

			// Compute expected with DFT
			expected := make([]complex128, n)
			dft := NewDft(n, Forward)
			copy(expected, input)
			dft.ProcessWithScratch(expected, make([]complex128, n))

			// Compare
			maxErr := 0.0
			for i := range buffer {
				err := cmplx.Abs(buffer[i] - expected[i])
				if err > maxErr {
					maxErr = err
				}
			}

			t.Logf("Prime %d: Max error = %.6e", n, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("Prime %d failed with error %.6e", n, maxErr)
				for i := range buffer {
					if cmplx.Abs(buffer[i]-expected[i]) > 1e-10 {
						t.Errorf("  [%d] got=%v want=%v", i, buffer[i], expected[i])
					}
				}
			}
		})
	}
}

// TestPrimeButterfliesRoundTrip tests forward + inverse for all primes
func TestPrimeButterfliesRoundTrip(t *testing.T) {
	primes := []int{3, 5, 7, 11, 13, 17, 19, 23, 29, 31}

	for _, n := range primes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create butterflies
			var fwd, inv FftInterface
			switch n {
			case 3:
				fwd = NewButterfly3(Forward)
				inv = NewButterfly3(Inverse)
			case 5:
				fwd = NewButterfly5(Forward)
				inv = NewButterfly5(Inverse)
			case 7:
				fwd = NewButterfly7(Forward)
				inv = NewButterfly7(Inverse)
			case 11:
				fwd = NewButterfly11(Forward)
				inv = NewButterfly11(Inverse)
			case 13:
				fwd = NewButterfly13(Forward)
				inv = NewButterfly13(Inverse)
			case 17:
				fwd = NewButterfly17(Forward)
				inv = NewButterfly17(Inverse)
			case 19:
				fwd = NewButterfly19(Forward)
				inv = NewButterfly19(Inverse)
			case 23:
				fwd = NewButterfly23(Forward)
				inv = NewButterfly23(Inverse)
			case 29:
				fwd = NewButterfly29(Forward)
				inv = NewButterfly29(Inverse)
			case 31:
				fwd = NewButterfly31(Forward)
				inv = NewButterfly31(Inverse)
			default:
				t.Fatalf("Unknown prime %d", n)
			}

			// Test data
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

			t.Logf("Prime %d: Round-trip error = %.6e", n, maxErr)

			if maxErr > 1e-10 {
				t.Errorf("Prime %d round-trip failed with error %.6e", n, maxErr)
			}
		})
	}
}
