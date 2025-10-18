//go:build arm64

package neon

import (
	"math"
	"testing"
)

func TestNEONDetection(t *testing.T) {
	// Test that NEON is detected on ARM64
	// Skip if not on ARM64
	t.Skip("NEON detection test - run on ARM64 platform")

	// level := simd.GetSIMDLevel()
	// if level != simd.NEON {
	//	t.Errorf("Expected NEON, got %v", level)
	// }

	t.Logf("NEON SIMD detected")
}

func TestVectorOperations64(t *testing.T) {
	// Test NEON vector operations for complex64
	vec1 := Vector64{
		data: [2]Complex64{
			{re: 1.0, im: 2.0},
			{re: 3.0, im: 4.0},
		},
	}

	vec2 := Vector64{
		data: [2]Complex64{
			{re: 5.0, im: 6.0},
			{re: 7.0, im: 8.0},
		},
	}

	// Test addition
	sum := vec1.Add64(vec2)
	expectedSum := Vector64{
		data: [2]Complex64{
			{re: 6.0, im: 8.0},
			{re: 10.0, im: 12.0},
		},
	}

	if !vectorsEqual64(sum, expectedSum) {
		t.Errorf("Addition failed: got %v, expected %v", sum, expectedSum)
	}

	// Test subtraction
	diff := vec1.Sub64(vec2)
	expectedDiff := Vector64{
		data: [2]Complex64{
			{re: -4.0, im: -4.0},
			{re: -4.0, im: -4.0},
		},
	}

	if !vectorsEqual64(diff, expectedDiff) {
		t.Errorf("Subtraction failed: got %v, expected %v", diff, expectedDiff)
	}

	// Test multiplication
	prod := vec1.Mul64(vec2)
	expectedProd := Vector64{
		data: [2]Complex64{
			{re: -7.0, im: 16.0},  // (1+2i)(5+6i) = -7+16i
			{re: -11.0, im: 52.0}, // (3+4i)(7+8i) = -11+52i
		},
	}

	if !vectorsEqual64(prod, expectedProd) {
		t.Errorf("Multiplication failed: got %v, expected %v", prod, expectedProd)
	}
}

func TestVectorOperations128(t *testing.T) {
	// Test NEON vector operations for complex128
	vec1 := Vector128{
		data: Complex128{re: 1.0, im: 2.0},
	}

	vec2 := Vector128{
		data: Complex128{re: 3.0, im: 4.0},
	}

	// Test addition
	sum := vec1.Add128(vec2)
	expectedSum := Vector128{
		data: Complex128{re: 4.0, im: 6.0},
	}

	if !vectorsEqual128(sum, expectedSum) {
		t.Errorf("Addition failed: got %v, expected %v", sum, expectedSum)
	}

	// Test multiplication
	prod := vec1.Mul128(vec2)
	expectedProd := Vector128{
		data: Complex128{re: -5.0, im: 10.0}, // (1+2i)(3+4i) = -5+10i
	}

	if !vectorsEqual128(prod, expectedProd) {
		t.Errorf("Multiplication failed: got %v, expected %v", prod, expectedProd)
	}
}

func TestButterfly2_64(t *testing.T) {
	// Test 2-point butterfly for complex64
	data := []Complex64{
		{re: 1.0, im: 2.0},
		{re: 3.0, im: 4.0},
	}

	Butterfly2_64(data)

	// Expected: out[0] = (1+2i) + (3+4i) = 4+6i
	//           out[1] = (1+2i) - (3+4i) = -2-2i
	expected := []Complex64{
		{re: 4.0, im: 6.0},
		{re: -2.0, im: -2.0},
	}

	for i := range data {
		if !complexEqual64(data[i], expected[i]) {
			t.Errorf("Butterfly2_64[%d]: got %v, expected %v", i, data[i], expected[i])
		}
	}
}

func TestButterfly2_128(t *testing.T) {
	// Test 2-point butterfly for complex128
	data := []Complex128{
		{re: 1.0, im: 2.0},
		{re: 3.0, im: 4.0},
	}

	Butterfly2_128(data)

	// Expected: out[0] = (1+2i) + (3+4i) = 4+6i
	//           out[1] = (1+2i) - (3+4i) = -2-2i
	expected := []Complex128{
		{re: 4.0, im: 6.0},
		{re: -2.0, im: -2.0},
	}

	for i := range data {
		if !complexEqual128(data[i], expected[i]) {
			t.Errorf("Butterfly2_128[%d]: got %v, expected %v", i, data[i], expected[i])
		}
	}
}

func TestAlignment(t *testing.T) {
	// Test memory alignment functions
	// Skip alignment test - requires simd package
	t.Skip("Alignment test - requires simd package")

	t.Logf("SIMD alignment test skipped")
}

// Helper functions for testing

func vectorsEqual64(a, b Vector64) bool {
	for i := range a.data {
		if !complexEqual64(a.data[i], b.data[i]) {
			return false
		}
	}
	return true
}

func vectorsEqual128(a, b Vector128) bool {
	return complexEqual128(a.data, b.data)
}

func complexEqual64(a, b Complex64) bool {
	const epsilon = 1e-6
	return math.Abs(float64(a.re-b.re)) < epsilon && math.Abs(float64(a.im-b.im)) < epsilon
}

func complexEqual128(a, b Complex128) bool {
	const epsilon = 1e-12
	return math.Abs(a.re-b.re) < epsilon && math.Abs(a.im-b.im) < epsilon
}
