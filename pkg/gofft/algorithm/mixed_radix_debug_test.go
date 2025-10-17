package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestMixedRadixSize6Debug - Detailed debugging for size 6 (2x3)
func TestMixedRadixSize6Debug(t *testing.T) {
	n1, n2 := 2, 3
	n := n1 * n2 // 6

	// Create simple input
	input := []complex128{
		complex(0, 0),
		complex(1, 0),
		complex(2, 0),
		complex(3, 0),
		complex(4, 0),
		complex(5, 0),
	}

	t.Logf("Input: %v", input)

	// Compute expected with DFT
	expected := make([]complex128, n)
	copy(expected, input)
	dft := NewDft(n, Forward)
	dft.ProcessWithScratch(expected, make([]complex128, n))
	t.Logf("Expected (DFT): %v", expected)

	// Compute with MixedRadix
	width := NewDft(n1, Forward)
	height := NewDft(n2, Forward)
	mr := NewMixedRadix(width, height)

	result := make([]complex128, n)
	copy(result, input)
	scratch := make([]complex128, mr.InplaceScratchLen())

	// Add detailed logging inside MixedRadix steps
	t.Logf("\n=== Starting MixedRadix ===")
	t.Logf("Width=%d, Height=%d", mr.width, mr.height)

	// Manually execute steps with logging
	selfScratch := scratch[:mr.length]

	// STEP 1: Transpose
	t.Logf("\nSTEP 1: Transpose %dx%d to %dx%d", mr.width, mr.height, mr.height, mr.width)
	transpose(mr.width, mr.height, result, selfScratch)
	t.Logf("After transpose: %v", selfScratch)

	// STEP 2: Height FFTs
	t.Logf("\nSTEP 2: Perform %d FFTs of size %d", mr.width, mr.height)
	for i := 0; i < mr.width; i++ {
		chunk := selfScratch[i*mr.height : (i+1)*mr.height]
		t.Logf("  Before FFT[%d]: %v", i, chunk)
		colScratch := make([]complex128, mr.heightFft.InplaceScratchLen())
		mr.heightFft.ProcessWithScratch(chunk, colScratch)
		t.Logf("  After FFT[%d]: %v", i, chunk)
	}
	t.Logf("After height FFTs: %v", selfScratch)

	// STEP 3: Twiddles
	t.Logf("\nSTEP 3: Apply twiddle factors")
	t.Logf("Twiddles: %v", mr.twiddles)
	for i := range selfScratch {
		old := selfScratch[i]
		selfScratch[i] = selfScratch[i] * mr.twiddles[i]
		t.Logf("  [%d]: %v * %v = %v", i, old, mr.twiddles[i], selfScratch[i])
	}
	t.Logf("After twiddles: %v", selfScratch)

	// STEP 4: Transpose back
	t.Logf("\nSTEP 4: Transpose %dx%d to %dx%d", mr.height, mr.width, mr.width, mr.height)
	transpose(mr.height, mr.width, selfScratch, result)
	t.Logf("After transpose: %v", result)

	// STEP 5: Width FFTs
	t.Logf("\nSTEP 5: Perform %d FFTs of size %d", mr.height, mr.width)
	for i := 0; i < mr.height; i++ {
		chunk := result[i*mr.width : (i+1)*mr.width]
		t.Logf("  Before FFT[%d]: %v", i, chunk)
		rowScratch := make([]complex128, mr.widthFft.InplaceScratchLen())
		mr.widthFft.ProcessWithScratch(chunk, rowScratch)
		t.Logf("  After FFT[%d]: %v", i, chunk)
	}
	t.Logf("Final result: %v", result)

	// Compare
	t.Logf("\n=== Comparison ===")
	maxErr := 0.0
	for i := range result {
		err := cmplx.Abs(result[i] - expected[i])
		if err > maxErr {
			maxErr = err
		}
		t.Logf("[%d] got=%v want=%v err=%.6e", i, result[i], expected[i], err)
	}

	t.Logf("\nMax error: %.6e", maxErr)

	if maxErr > 1e-10 {
		t.Errorf("Size 6 failed with error %.6e", maxErr)
	}
}

// TestMixedRadixManualSize6 - Manual implementation to verify algorithm
func TestMixedRadixManualSize6(t *testing.T) {
	// Size 6 = 2x3
	// Input: [0, 1, 2, 3, 4, 5]

	input := []complex128{
		complex(0, 0),
		complex(1, 0),
		complex(2, 0),
		complex(3, 0),
		complex(4, 0),
		complex(5, 0),
	}

	t.Logf("Input: %v", input)

	// Manual six-step FFT for 2x3
	// Step 1: View as 2x3 matrix (width=2, height=3)
	//   [0 1]
	//   [2 3]
	//   [4 5]

	// Step 2: Transpose to 3x2
	//   [0 2 4]
	//   [1 3 5]
	transposed := make([]complex128, 6)
	for r := 0; r < 2; r++ {
		for c := 0; c < 3; c++ {
			transposed[c*2+r] = input[r*3+c]
		}
	}
	t.Logf("After transpose (3x2): %v", transposed)

	// Step 3: FFTs on columns (each of size 3)
	// Column 0: [0, 2, 4] -> 3-point FFT
	// Column 1: [1, 3, 5] -> 3-point FFT

	col0 := []complex128{transposed[0], transposed[2], transposed[4]}
	col1 := []complex128{transposed[1], transposed[3], transposed[5]}

	t.Logf("Column 0 before FFT: %v", col0)
	t.Logf("Column 1 before FFT: %v", col1)

	// 3-point DFT on col0
	dft3 := NewDft(3, Forward)
	scratch3 := make([]complex128, 3)
	dft3.ProcessWithScratch(col0, scratch3)
	dft3.ProcessWithScratch(col1, scratch3)

	t.Logf("Column 0 after FFT: %v", col0)
	t.Logf("Column 1 after FFT: %v", col1)

	// Put back
	transposed[0], transposed[2], transposed[4] = col0[0], col0[1], col0[2]
	transposed[1], transposed[3], transposed[5] = col1[0], col1[1], col1[2]
	t.Logf("After column FFTs: %v", transposed)

	// Step 4: Apply twiddle factors
	// Twiddle[x,y] = exp(-2πi * x*y / 6)
	// In transposed layout: index = y*width + x = y*2 + x
	for y := 0; y < 3; y++ {
		for x := 0; x < 2; x++ {
			idx := y*2 + x
			angle := -2.0 * 3.14159265358979323846 * float64(x*y) / 6.0
			twiddle := complex(cosf(angle), sinf(angle))
			old := transposed[idx]
			transposed[idx] *= twiddle
			t.Logf("Twiddle[%d,%d] (idx=%d): %v * %v = %v", x, y, idx, old, twiddle, transposed[idx])
		}
	}
	t.Logf("After twiddles: %v", transposed)

	// Step 5: Transpose back to 2x3
	result := make([]complex128, 6)
	for r := 0; r < 3; r++ {
		for c := 0; c < 2; c++ {
			result[c*3+r] = transposed[r*2+c]
		}
	}
	t.Logf("After transpose (2x3): %v", result)

	// Step 6: FFTs on rows (each of size 2)
	for i := 0; i < 3; i++ {
		row := result[i*2 : (i+1)*2]
		t.Logf("Row %d before FFT: %v", i, row)
		dft2 := NewDft(2, Forward)
		scratch2 := make([]complex128, 2)
		dft2.ProcessWithScratch(row, scratch2)
		t.Logf("Row %d after FFT: %v", i, row)
	}
	t.Logf("Final result: %v", result)

	// Compare with DFT
	expected := make([]complex128, 6)
	copy(expected, input)
	dft := NewDft(6, Forward)
	dft.ProcessWithScratch(expected, make([]complex128, 6))

	t.Logf("\n=== Comparison ===")
	t.Logf("Expected: %v", expected)
	maxErr := 0.0
	for i := range result {
		err := cmplx.Abs(result[i] - expected[i])
		if err > maxErr {
			maxErr = err
		}
		t.Logf("[%d] got=%v want=%v err=%.6e", i, result[i], expected[i], err)
	}

	t.Logf("\nMax error: %.6e", maxErr)

	if maxErr > 1e-10 {
		t.Errorf("Manual size 6 failed with error %.6e", maxErr)
	}
}

func cosf(x float64) float64 {
	// Use built-in
	return real(cmplx.Exp(complex(0, x)))
}

func sinf(x float64) float64 {
	return imag(cmplx.Exp(complex(0, x)))
}
