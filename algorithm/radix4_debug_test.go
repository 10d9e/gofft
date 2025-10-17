package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestRadix4Deterministic creates step-by-step deterministic tests
func TestRadix4Deterministic(t *testing.T) {
	// Test with simple known input
	size := 32

	// Create input: simple sequence 0, 1, 2, ...
	input := make([]complex128, size)
	for i := range input {
		input[i] = complex(float64(i), 0)
	}

	t.Logf("Input: first 8 values: %v", input[:8])

	// Compute using Radix4
	fft := NewRadix4(size, Forward)
	buffer := make([]complex128, size)
	copy(buffer, input)
	scratch := make([]complex128, fft.InplaceScratchLen())

	t.Logf("Base FFT: %T, baseLen=%d", fft.baseFft, fft.baseLen)
	t.Logf("Number of twiddles: %d", len(fft.twiddles))
	t.Logf("Expected k value: %d", (trailingZeros(size)-trailingZeros(fft.baseLen))/2)

	fft.ProcessWithScratch(buffer, scratch)

	t.Logf("Output: first 8 values: %v", buffer[:8])

	// Compute expected using DFT
	expected := make([]complex128, size)
	copy(expected, input)
	dft := NewDft(size, Forward)
	dftScratch := make([]complex128, dft.InplaceScratchLen())
	dft.ProcessWithScratch(expected, dftScratch)

	t.Logf("Expected: first 8 values: %v", expected[:8])

	// Compare and show pattern of errors
	errors := make([]float64, size)
	for i := range buffer {
		errors[i] = cmplx.Abs(buffer[i] - expected[i])
	}

	// Find max error and its location
	maxErr := 0.0
	maxIdx := 0
	for i, err := range errors {
		if err > maxErr {
			maxErr = err
			maxIdx = i
		}
	}

	t.Logf("Max error: %.6f at index %d", maxErr, maxIdx)

	// Show error pattern
	t.Logf("Error pattern (indices with error > 0.1):")
	for i, err := range errors {
		if err > 0.1 {
			t.Logf("  [%2d] error=%.2f  got=%v  want=%v", i, err, buffer[i], expected[i])
		}
	}

	if maxErr > 1e-10 {
		t.Errorf("Max error too large: %.6e", maxErr)
	}
}

// TestBitReversedTransposeDetailed tests the transpose step in detail
func TestBitReversedTransposeDetailed(t *testing.T) {
	size := 32
	baseLen := 8

	input := make([]complex128, size)
	for i := range input {
		input[i] = complex(float64(i), 0)
	}

	output := make([]complex128, size)
	bitReversedTranspose4(baseLen, input, output)

	t.Logf("Size=%d, baseLen=%d", size, baseLen)
	t.Logf("Rows=%d, Cols=%d", size/baseLen, baseLen)

	// Show the transpose
	t.Logf("\nTranspose result:")
	for i := 0; i < size; i++ {
		t.Logf("output[%2d] = input[%2d] = %v", i, int(real(output[i])), output[i])
	}

	// For size=32, baseLen=8:
	// We have 4 rows (0-3), 8 columns (0-7)
	// Input is in row-major: [row0, row1, row2, row3] where each row has 8 elements
	// Output should be column-major with bit-reversed rows

	// Bit-reverse for 2 bits (4 rows): 0->0, 1->2, 2->1, 3->3
	bitRev := []int{0, 2, 1, 3}

	t.Logf("\nExpected mapping (col=0):")
	for row := 0; row < 4; row++ {
		revRow := bitRev[row]
		inputIdx := row*baseLen + 0
		outputIdx := 0*4 + revRow
		t.Logf("output[%2d] should be input[%2d] (row %d bit-reversed to %d)", outputIdx, inputIdx, row, revRow)
	}
}

// TestButterfly4StageIsolated tests just the butterfly stage
func TestButterfly4StageIsolated(t *testing.T) {
	// Create simple test data: 8 columns, 4 rows
	numColumns := 8
	data := make([]complex128, 32)

	// Fill with pattern: each group of 8 is [0,1,2,3,4,5,6,7], [8,9,10,...], etc
	for i := range data {
		data[i] = complex(float64(i), 0)
	}

	t.Logf("Before butterfly stage:")
	for row := 0; row < 4; row++ {
		t.Logf("Row %d:", row)
		for col := 0; col < 8; col++ {
			idx := col + row*numColumns
			t.Logf("  col=%d: data[%2d] = %v", col, idx, data[idx])
		}
	}

	// Create simple twiddles (all 1s for this test)
	twiddles := make([]complex128, numColumns*3)
	for i := range twiddles {
		twiddles[i] = complex(1, 0)
	}

	bf4 := NewButterfly4(Forward)
	butterfly4Stage(data, twiddles, numColumns, bf4)

	t.Logf("\nAfter butterfly stage:")
	for row := 0; row < 4; row++ {
		t.Logf("Row %d:", row)
		for col := 0; col < 8; col++ {
			idx := col + row*numColumns
			t.Logf("  col=%d: data[%2d] = %v", col, idx, data[idx])
		}
	}
}

// TestFullRadix4StepByStep walks through the algorithm step by step
func TestFullRadix4StepByStep(t *testing.T) {
	size := 32
	baseLen := 8

	// Simple input
	input := make([]complex128, size)
	for i := range input {
		input[i] = complex(float64(i), 0)
	}

	t.Logf("=== Step 1: Input ===")
	t.Logf("First 8: %v", input[:8])

	// Step 2: Bit-reversed transpose
	output := make([]complex128, size)
	bitReversedTranspose4(baseLen, input, output)

	t.Logf("\n=== Step 2: After bit-reversed transpose ===")
	t.Logf("First 8: %v", output[:8])
	t.Logf("Second 8: %v", output[8:16])

	// Step 3: Base FFTs (8-point)
	baseFft := NewButterfly8(Forward)
	baseScratch := make([]complex128, baseFft.InplaceScratchLen())
	for i := 0; i < size; i += baseLen {
		chunk := output[i : i+baseLen]
		baseFft.ProcessWithScratch(chunk, baseScratch)
	}

	t.Logf("\n=== Step 3: After base FFTs (8-point) ===")
	t.Logf("First 8: %v", output[:8])

	// Step 4: Cross FFTs (radix-4 stage)
	// For size=32, baseLen=8, we have one radix-4 stage with numColumns=8
	fft := NewRadix4(size, Forward)

	t.Logf("\n=== Step 4: Radix-4 cross FFT ===")
	t.Logf("numColumns should be: %d", baseLen)
	t.Logf("Twiddles available: %d (need %d)", len(fft.twiddles), baseLen*3)

	bf4 := NewButterfly4(Forward)
	butterfly4Stage(output, fft.twiddles, baseLen, bf4)

	t.Logf("First 8: %v", output[:8])
	t.Logf("Next 8: %v", output[8:16])

	// Compare with expected
	expected := make([]complex128, size)
	copy(expected, input)
	dft := NewDft(size, Forward)
	dft.ProcessWithScratch(expected, make([]complex128, dft.InplaceScratchLen()))

	t.Logf("\n=== Comparison ===")
	for i := 0; i < size; i++ {
		err := cmplx.Abs(output[i] - expected[i])
		if err > 0.1 {
			t.Logf("[%2d] got=%v want=%v err=%.3f", i, output[i], expected[i], err)
		}
	}
}

// TestCompareWithRustFFTPattern tests specific known values
func TestCompareWithRustFFTPattern(t *testing.T) {
	// Use impulse at position 0
	size := 32
	input := make([]complex128, size)
	input[0] = complex(1, 0) // Impulse

	fft := NewRadix4(size, Forward)
	buffer := make([]complex128, size)
	copy(buffer, input)
	scratch := make([]complex128, fft.InplaceScratchLen())
	fft.ProcessWithScratch(buffer, scratch)

	t.Logf("Impulse FFT output:")
	// FFT of impulse should be all 1s
	for i := 0; i < size; i++ {
		expected := complex(1, 0)
		err := cmplx.Abs(buffer[i] - expected)
		if err > 1e-10 {
			t.Errorf("[%2d] got=%v want=%v err=%.6e", i, buffer[i], expected, err)
		}
	}
}

// TestRadix4Size8 tests size 8 which should work
func TestRadix4Size8(t *testing.T) {
	size := 8
	input := make([]complex128, size)
	for i := range input {
		input[i] = complex(float64(i), 0)
	}

	fft := NewRadix4(size, Forward)
	buffer := make([]complex128, size)
	copy(buffer, input)
	scratch := make([]complex128, fft.InplaceScratchLen())
	fft.ProcessWithScratch(buffer, scratch)

	// Expected using DFT
	expected := make([]complex128, size)
	copy(expected, input)
	dft := NewDft(size, Forward)
	dft.ProcessWithScratch(expected, make([]complex128, dft.InplaceScratchLen()))

	// Compare
	maxErr := 0.0
	for i := range buffer {
		err := cmplx.Abs(buffer[i] - expected[i])
		if err > maxErr {
			maxErr = err
		}
	}

	if maxErr > 1e-10 {
		t.Errorf("Size 8 failed with max error: %.6e", maxErr)
		for i := range buffer {
			err := cmplx.Abs(buffer[i] - expected[i])
			if err > 1e-10 {
				t.Logf("[%d] got=%v want=%v", i, buffer[i], expected[i])
			}
		}
	} else {
		t.Logf("Size 8 passes! Max error: %.6e", maxErr)
	}
}
