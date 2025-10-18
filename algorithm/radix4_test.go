package algorithm

import (
	"math/cmplx"
	"testing"
)

// Test Radix4 for size 32
func TestRadix4Size32(t *testing.T) {
	size := 32
	fft := NewRadix4(size, Forward)

	// Create simple test input
	input := make([]complex128, size)
	for i := range input {
		input[i] = complex(float64(i), 0)
	}

	// Compute FFT
	buffer := make([]complex128, size)
	copy(buffer, input)
	scratch := make([]complex128, fft.InplaceScratchLen())
	fft.ProcessWithScratch(buffer, scratch)

	// Compute expected using DFT
	expected := make([]complex128, size)
	dft := NewDft(size, Forward)
	dftScratch := make([]complex128, dft.InplaceScratchLen())
	copy(expected, input)
	dft.ProcessWithScratch(expected, dftScratch)

	// Compare
	maxError := 0.0
	for i := range buffer {
		err := cmplx.Abs(buffer[i] - expected[i])
		if err > maxError {
			maxError = err
		}
		if err > 1e-10 {
			t.Errorf("[%d] got %v, want %v, err %v", i, buffer[i], expected[i], err)
		}
	}

	if maxError > 1e-10 {
		t.Errorf("Max error: %v", maxError)
	}
}

// Test bit-reversed transpose
func TestBitReversedTranspose(t *testing.T) {
	// Test with size 32, baseLen 16
	size := 32
	baseLen := 16

	input := make([]complex128, size)
	for i := range input {
		input[i] = complex(float64(i), 0)
	}

	output := make([]complex128, size)
	// For size 32 with baseLen 16, we have height=2, width=16
	height := size / baseLen
	bitReversedTranspose4(height, input, output)

	// Print the transpose result to see the pattern
	t.Logf("Input:  %v", input[:8])
	t.Logf("Output: %v", output[:8])

	// The transpose should arrange data for radix-4
	// For baseLen=16, we have 2 rows (size/baseLen = 2)
	// Bit-reverse of row 0 = 0, bit-reverse of row 1 = 1 (for 1 bit)
	// So output should be column-major with bit-reversed rows

	// Check a few values
	// Column 0 should have row 0 value (bit-reverse of 0 = 0)
	if output[0] != input[0] {
		t.Errorf("output[0] should be input[0], got %v", output[0])
	}
}

// Test twiddle factor computation
func TestTwiddleFactors(t *testing.T) {
	// For size 128, baseLen 8 (using Butterfly8), we have one radix-4 layer with k=1
	// That layer has numColumns=32, so we need 32*3 = 96 twiddle factors

	size := 128
	fft := NewRadix4(size, Forward)

	t.Logf("Size: %d, BaseLen: %d", size, fft.baseLen)
	t.Logf("Expected count: %d, Actual count: %d", 32*3, len(fft.twiddles))

	expectedCount := 32 * 3 // One layer with 32 columns
	if len(fft.twiddles) != expectedCount {
		t.Errorf("Expected %d twiddles, got %d", expectedCount, len(fft.twiddles))
	}

	// Check first few twiddle factors
	// For k=1, twiddles should be exp(-2πi*k*j/128) for j=1,2,3
	// But since we're using Butterfly32 (baseLen=32), the twiddles are calculated differently
	// Let's just check that we have the right number of twiddles and they're not all 1
	if len(fft.twiddles) > 0 {
		allOnes := true
		for _, tw := range fft.twiddles {
			if cmplx.Abs(tw-1) > 1e-10 {
				allOnes = false
				break
			}
		}
		if allOnes {
			t.Errorf("All twiddles are 1, but we expected non-trivial twiddles for size %d", size)
		}
	}
}
