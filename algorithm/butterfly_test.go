package algorithm

import (
	"math/cmplx"
	"testing"
)

// TestButterfly4Direct tests Butterfly4 in isolation
func TestButterfly4Direct(t *testing.T) {
	// Test with simple input
	input := []complex128{
		complex(0, 0),
		complex(1, 0),
		complex(2, 0),
		complex(3, 0),
	}

	bf4 := NewButterfly4(Forward)

	// Test in-place
	buffer := make([]complex128, 4)
	copy(buffer, input)
	bf4.performFft(buffer)

	t.Logf("Butterfly4 in-place result: %v", buffer)

	// Test out-of-place
	buffer2 := make([]complex128, 4)
	bf4.performFftOutOfPlace(input, buffer2)

	t.Logf("Butterfly4 out-of-place result: %v", buffer2)

	// They should match
	for i := range buffer {
		if cmplx.Abs(buffer[i]-buffer2[i]) > 1e-10 {
			t.Errorf("In-place and out-of-place don't match at index %d", i)
		}
	}

	// Compare with naive DFT
	dft := NewDft(4, Forward)
	expected := make([]complex128, 4)
	copy(expected, input)
	dft.ProcessWithScratch(expected, make([]complex128, 4))

	t.Logf("DFT result: %v", expected)

	for i := range buffer {
		err := cmplx.Abs(buffer[i] - expected[i])
		if err > 1e-10 {
			t.Errorf("[%d] got=%v want=%v err=%.6e", i, buffer[i], expected[i], err)
		}
	}
}

// TestButterfly4InButterfly4Stage tests if the issue is in how we call Butterfly4
func TestButterfly4InButterfly4Stage(t *testing.T) {
	// Simulate what happens in butterfly4Stage
	numColumns := 2
	data := []complex128{
		// Column 0, rows 0-3
		complex(0, 0), // row 0, col 0
		complex(2, 0), // row 1, col 0
		complex(4, 0), // row 2, col 0
		complex(6, 0), // row 3, col 0
		// Column 1, rows 0-3
		complex(1, 0), // row 0, col 1
		complex(3, 0), // row 1, col 1
		complex(5, 0), // row 2, col 1
		complex(7, 0), // row 3, col 1
	}

	t.Logf("Input data (column-major, 2 cols x 4 rows):")
	t.Logf("  Col 0: %v %v %v %v", data[0], data[2], data[4], data[6])
	t.Logf("  Col 1: %v %v %v %v", data[1], data[3], data[5], data[7])

	// All twiddles = 1 for simplicity
	twiddles := []complex128{
		complex(1, 0), complex(1, 0), complex(1, 0), // for column 0
		complex(1, 0), complex(1, 0), complex(1, 0), // for column 1
	}

	bf4 := NewButterfly4(Forward)
	butterfly4Stage(data, twiddles, numColumns, bf4)

	t.Logf("\nAfter butterfly4Stage:")
	t.Logf("  Col 0: %v %v %v %v", data[0], data[2], data[4], data[6])
	t.Logf("  Col 1: %v %v %v %v", data[1], data[3], data[5], data[7])

	// Manually compute what we expect
	// For column 0: input is [0, 2, 4, 6]
	// After 4-point FFT with twiddles all 1: should be sum, diff of evens/odds, etc
	col0Input := []complex128{complex(0, 0), complex(2, 0), complex(4, 0), complex(6, 0)}
	dft := NewDft(4, Forward)
	expected := make([]complex128, 4)
	copy(expected, col0Input)
	dft.ProcessWithScratch(expected, make([]complex128, 4))

	t.Logf("\nExpected for column 0: %v", expected)
	t.Logf("Got for column 0: [%v %v %v %v]", data[0], data[2], data[4], data[6])

	// Check column 0
	colData := []complex128{data[0], data[2], data[4], data[6]}
	for i, val := range colData {
		err := cmplx.Abs(val - expected[i])
		if err > 1e-10 {
			t.Errorf("Column 0, row %d: got=%v want=%v", i, val, expected[i])
		}
	}
}
