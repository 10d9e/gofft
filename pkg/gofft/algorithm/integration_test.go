package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestAllSizesUpTo100 tests FFT correctness for all sizes up to 100
func TestAllSizesUpTo100(t *testing.T) {
	for n := 2; n <= 100; n++ {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create appropriate FFT
			var fft FftInterface
			switch n {
			case 2:
				fft = NewButterfly2(Forward)
			case 3:
				fft = NewButterfly3(Forward)
			case 4:
				fft = NewButterfly4(Forward)
			case 5:
				fft = NewButterfly5(Forward)
			case 6:
				fft = NewButterfly6(Forward)
			case 7:
				fft = NewButterfly7(Forward)
			case 8:
				fft = NewButterfly8(Forward)
			case 9:
				fft = NewButterfly9(Forward)
			case 11:
				fft = NewButterfly11(Forward)
			case 12:
				fft = NewButterfly12(Forward)
			case 13:
				fft = NewButterfly13(Forward)
			case 16:
				fft = NewButterfly16(Forward)
			case 17:
				fft = NewButterfly17(Forward)
			case 19:
				fft = NewButterfly19(Forward)
			case 23:
				fft = NewButterfly23(Forward)
			case 24:
				fft = NewButterfly24(Forward)
			case 27:
				fft = NewButterfly27(Forward)
			case 29:
				fft = NewButterfly29(Forward)
			case 31:
				fft = NewButterfly31(Forward)
			case 32:
				fft = NewButterfly32(Forward)
			case 64:
				fft = NewRadix4(64, Forward)
			default:
				// Use DFT for other sizes
				fft = NewDft(n, Forward)
			}

			// Test input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i%7), float64(i%5)*0.3)
			}

			// Compute with selected FFT
			result := make([]complex128, n)
			copy(result, input)
			scratch := make([]complex128, fft.InplaceScratchLen())
			fft.ProcessWithScratch(result, scratch)

			// Compute expected with DFT
			expected := make([]complex128, n)
			copy(expected, input)
			dft := NewDft(n, Forward)
			dftScratch := make([]complex128, n)
			dft.ProcessWithScratch(expected, dftScratch)

			// Compare
			maxErr := 0.0
			for i := range result {
				err := cmplx.Abs(result[i] - expected[i])
				if err > maxErr {
					maxErr = err
				}
			}

			if maxErr > 1e-10 {
				t.Errorf("Size %d failed with error %.6e", n, maxErr)
			}
		})
	}
}

// TestPowerOfTwoSizes tests all power-of-two sizes
func TestPowerOfTwoSizes(t *testing.T) {
	sizes := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096}

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Use Radix4 for sizes >= 64
			var fft FftInterface
			if n >= 64 {
				fft = NewRadix4(n, Forward)
			} else {
				// Use butterflies for smaller sizes
				switch n {
				case 2:
					fft = NewButterfly2(Forward)
				case 4:
					fft = NewButterfly4(Forward)
				case 8:
					fft = NewButterfly8(Forward)
				case 16:
					fft = NewButterfly16(Forward)
				case 32:
					fft = NewButterfly32(Forward)
				}
			}

			// Test input
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i%11)*0.7, float64(i%7)*0.3)
			}

			// Compute with selected FFT
			result := make([]complex128, n)
			copy(result, input)
			scratch := make([]complex128, fft.InplaceScratchLen())
			fft.ProcessWithScratch(result, scratch)

			// Compute expected with DFT (for sizes <= 256, otherwise just check round-trip)
			if n <= 256 {
				expected := make([]complex128, n)
				copy(expected, input)
				dft := NewDft(n, Forward)
				dftScratch := make([]complex128, n)
				dft.ProcessWithScratch(expected, dftScratch)

				// Compare
				maxErr := 0.0
				for i := range result {
					err := cmplx.Abs(result[i] - expected[i])
					if err > maxErr {
						maxErr = err
					}
				}

				t.Logf("Size %d: Max error = %.6e", n, maxErr)

				if maxErr > 1e-10 {
					t.Errorf("Size %d failed with error %.6e", n, maxErr)
				}
			} else {
				// For large sizes, just check round-trip
				var inv FftInterface
				if n >= 64 {
					inv = NewRadix4(n, Inverse)
				}

				invScratch := make([]complex128, inv.InplaceScratchLen())
				inv.ProcessWithScratch(result, invScratch)

				// Normalize
				for i := range result {
					result[i] /= complex(float64(n), 0)
				}

				// Check
				maxErr := 0.0
				for i := range result {
					err := cmplx.Abs(result[i] - input[i])
					if err > maxErr {
						maxErr = err
					}
				}

				t.Logf("Size %d: Round-trip error = %.6e", n, maxErr)

				if maxErr > 1e-10 {
					t.Errorf("Size %d round-trip failed with error %.6e", n, maxErr)
				}
			}
		})
	}
}
