package algorithm

import (
	"testing"
)

// TestActualLayoutAfterBaseFFTs checks the real data layout
func TestActualLayoutAfterBaseFFTs(t *testing.T) {
	size := 32
	baseLen := 8

	// Original input
	input := make([]complex128, size)
	for i := range input {
		input[i] = complex(float64(i), 0)
	}

	// Step 1: Bit-reversed transpose
	output := make([]complex128, size)
	bitReversedTranspose4(baseLen, input, output)

	t.Logf("After transpose:")
	t.Logf("output[0:8]   = %v", output[0:8])
	t.Logf("output[8:16]  = %v", output[8:16])
	t.Logf("output[16:24] = %v", output[16:24])
	t.Logf("output[24:32] = %v", output[24:32])

	// Step 2: Base FFTs (process in chunks of 8)
	baseFft := NewButterfly8(Forward)
	baseScratch := make([]complex128, baseFft.InplaceScratchLen())

	t.Logf("\nProcessing base FFTs on chunks of %d:", baseLen)
	for i := 0; i < size; i += baseLen {
		chunk := output[i : i+baseLen]
		t.Logf("  Chunk starting at %d: %v (before FFT)", i, chunk)
		baseFft.ProcessWithScratch(chunk, baseScratch)
		t.Logf("  Chunk starting at %d: %v (after FFT)", i, chunk)
	}

	t.Logf("\nAfter all base FFTs:")
	t.Logf("output[0:8]   = %v", output[0:8])
	t.Logf("output[8:16]  = %v", output[8:16])
	t.Logf("output[16:24] = %v", output[16:24])
	t.Logf("output[24:32] = %v", output[24:32])

	// Now for butterfly4Stage to work, it needs data in this layout:
	// Indices [0, 8, 16, 24] = column 0 (4 values)
	// Indices [1, 9, 17, 25] = column 1 (4 values)
	// ...
	// Indices [7, 15, 23, 31] = column 7 (4 values)

	t.Logf("\nWhat butterfly4Stage will see (numColumns=%d):", baseLen)
	for col := 0; col < baseLen; col++ {
		idx0 := col
		idx1 := col + baseLen
		idx2 := col + 2*baseLen
		idx3 := col + 3*baseLen
		t.Logf("  Column %d: [%2d,%2d,%2d,%2d] = %v",
			col, idx0, idx1, idx2, idx3,
			[]complex128{output[idx0], output[idx1], output[idx2], output[idx3]})
	}
}

// TestCompareLayouts compares our layout with what we expect
func TestCompareLayouts(t *testing.T) {
	// Simple test: what should happen for size 8 with baseLen 4?
	// This should NOT need any radix-4 cross FFTs (k=0)

	size := 8

	// Test that size 8 uses Butterfly8 directly (no Radix4 layers)
	fft := NewRadix4(size, Forward)

	t.Logf("Size=%d, baseLen=%d", size, fft.baseLen)
	t.Logf("Should k=0? Actual k=%d", (trailingZeros(size)-trailingZeros(fft.baseLen))/2)

	// For size 8, we should use baseLen=8, giving k=0
	if fft.baseLen != size {
		t.Logf("Note: baseLen != size, so there WILL be radix-4 stages")
	} else {
		t.Logf("baseLen == size, so NO radix-4 stages needed")
	}
}
