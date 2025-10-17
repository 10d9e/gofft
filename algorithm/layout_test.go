package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestDataLayout verifies we understand the data layout
func TestDataLayout(t *testing.T) {
	// After bit-reversed transpose with size=32, baseLen=8
	// We should have 4 rows of 8 columns each
	size := 32
	baseLen := 8
	numRows := size / baseLen // Should be 4

	input := make([]complex128, size)
	for i := range input {
		input[i] = complex(float64(i), 0)
	}

	output := make([]complex128, size)
	bitReversedTranspose4(baseLen, input, output)

	t.Logf("After transpose, size=%d, baseLen=%d, numRows=%d", size, baseLen, numRows)
	t.Logf("Full output: %v", output)

	// The bit-reversed transpose should give us column-major order with bit-reversed rows
	// For 4 rows (2 bits), bit reverse is: 0->0, 1->2, 2->1, 3->3

	t.Logf("\nExpected column-major layout with bit-reversed rows:")
	for col := 0; col < baseLen; col++ {
		t.Logf("Column %d:", col)
		for row := 0; row < numRows; row++ {
			idx := col*numRows + row
			t.Logf("  Row %d: output[%2d] = %v", row, idx, output[idx])
		}
	}

	// After base FFTs, we should have data organized for the butterfly stage
	// The butterfly stage operates on numColumns where each column has 4 elements
	// For size=32 after baseLen=8 FFTs, we have:
	// - Data is organized in groups of baseLen (8 elements each)
	// - Each group has been FFT'd
	// - Now we need to do radix-4 on columns

	t.Logf("\nFor butterfly4Stage with numColumns=%d:", baseLen)
	t.Logf("Each column should consist of 4 values strided by %d", baseLen)
	for col := 0; col < baseLen; col++ {
		idx0 := col
		idx1 := col + baseLen
		idx2 := col + 2*baseLen
		idx3 := col + 3*baseLen
		t.Logf("Column %d: indices [%2d, %2d, %2d, %2d] = values %v",
			col, idx0, idx1, idx2, idx3,
			[]complex128{output[idx0], output[idx1], output[idx2], output[idx3]})
	}
}

// TestButterfly4StageCorrectness checks if our butterfly4Stage is correct
func TestButterfly4StageCorrectness(t *testing.T) {
	// Create a simple test: 2 columns, 4 rows
	// Data in row-major order for clarity, then we'll reorganize
	numColumns := 2

	// Initialize as if we had 4 rows of 2 elements each
	// In column-major: col0=[row0, row1, row2, row3], col1=[row0, row1, row2, row3]
	// Strided layout: [col0row0, col1row0, col0row1, col1row1, col0row2, col1row2, col0row3, col1row3]
	// But butterfly4Stage expects: [col0row0, col0row1, col0row2, col0row3, col1row0, col1row1, col1row2, col1row3]

	data := []complex128{
		complex(1, 0), complex(2, 0), // row 0: col 0, col 1
		complex(3, 0), complex(4, 0), // row 1: col 0, col 1
		complex(5, 0), complex(6, 0), // row 2: col 0, col 1
		complex(7, 0), complex(8, 0), // row 3: col 0, col 1
	}

	t.Logf("Input (row-major, 4 rows x 2 cols):")
	for row := 0; row < 4; row++ {
		t.Logf("  Row %d: %v", row, data[row*2:(row+1)*2])
	}

	// All twiddles = 1
	twiddles := []complex128{
		complex(1, 0), complex(1, 0), complex(1, 0),
		complex(1, 0), complex(1, 0), complex(1, 0),
	}

	bf4 := NewButterfly4(Forward)
	butterfly4Stage(data, twiddles, numColumns, bf4)

	t.Logf("\nAfter butterfly4Stage:")
	for row := 0; row < 4; row++ {
		t.Logf("  Row %d: %v", row, data[row*2:(row+1)*2])
	}

	// Now compute expected:
	// Column 0 should have FFT of [1, 3, 5, 7]
	// Column 1 should have FFT of [2, 4, 6, 8]
	dft := NewDft(4, Forward)

	col0 := []complex128{complex(1, 0), complex(3, 0), complex(5, 0), complex(7, 0)}
	expected0 := make([]complex128, 4)
	copy(expected0, col0)
	dft.ProcessWithScratch(expected0, make([]complex128, 4))

	col1 := []complex128{complex(2, 0), complex(4, 0), complex(6, 0), complex(8, 0)}
	expected1 := make([]complex128, 4)
	copy(expected1, col1)
	dft.ProcessWithScratch(expected1, make([]complex128, 4))

	t.Logf("\nExpected column 0 (FFT of [1,3,5,7]): %v", expected0)
	t.Logf("Expected column 1 (FFT of [2,4,6,8]): %v", expected1)

	// Extract columns from result
	got0 := []complex128{data[0], data[2], data[4], data[6]}
	got1 := []complex128{data[1], data[3], data[5], data[7]}

	t.Logf("Got column 0: %v", got0)
	t.Logf("Got column 1: %v", got1)

	for i := 0; i < 4; i++ {
		if cmplx.Abs(got0[i]-expected0[i]) > 1e-10 {
			t.Errorf("Col 0, row %d: got=%v want=%v", i, got0[i], expected0[i])
		}
		if cmplx.Abs(got1[i]-expected1[i]) > 1e-10 {
			t.Errorf("Col 1, row %d: got=%v want=%v", i, got1[i], expected1[i])
		}
	}
}
