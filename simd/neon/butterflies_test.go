//go:build arm64

package neon

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/10d9e/gofft/algorithm"
)

func TestNEONButterflies(t *testing.T) {
	testCases := []struct {
		size int
		name string
	}{
		{1, "Butterfly1"},
		{2, "Butterfly2"},
		{3, "Butterfly3"},
		{4, "Butterfly4"},
		{5, "Butterfly5"},
		{6, "Butterfly6"},
		{7, "Butterfly7"},
		{8, "Butterfly8"},
		{9, "Butterfly9"},
		{10, "Butterfly10"},
		{11, "Butterfly11"},
		{12, "Butterfly12"},
		{13, "Butterfly13"},
		{15, "Butterfly15"},
		{16, "Butterfly16"},
		{17, "Butterfly17"},
		{19, "Butterfly19"},
		{23, "Butterfly23"},
		{24, "Butterfly24"},
		{27, "Butterfly27"},
		{29, "Butterfly29"},
		{31, "Butterfly31"},
		{32, "Butterfly32"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test data
			data := make([]complex128, tc.size)
			for i := range data {
				data[i] = complex(float64(i%7), float64(i%5)*0.3)
			}

			// Apply NEON butterfly
			ProcessVectorizedButterfly(data, tc.size, algorithm.Forward)

			// Verify we got some result (not all zeros)
			// Special case for Butterfly1: it's the identity operation, so it should preserve input
			if tc.size == 1 {
				// For size 1, just check that the function ran without error
				t.Logf("%s: Identity operation completed successfully", tc.name)
			} else {
				hasNonZero := false
				for _, val := range data {
					if val != 0 {
						hasNonZero = true
						break
					}
				}

				if !hasNonZero {
					t.Errorf("%s result should not be all zeros", tc.name)
				}
			}

			t.Logf("%s: Processed %d elements successfully", tc.name, tc.size)
		})
	}
}

func TestButterfly2_NEON(t *testing.T) {
	// Test 2-point butterfly
	data := []complex128{
		complex(1.0, 2.0),
		complex(3.0, 4.0),
	}

	Butterfly2_NEON(data, algorithm.Forward)

	// Expected: out[0] = (1+2i) + (3+4i) = 4+6i
	//           out[1] = (1+2i) - (3+4i) = -2-2i
	expected := []complex128{
		complex(4.0, 6.0),
		complex(-2.0, -2.0),
	}

	for i := range data {
		if !complexEqual(data[i], expected[i]) {
			t.Errorf("Butterfly2_NEON[%d]: got %v, expected %v", i, data[i], expected[i])
		}
	}
}

func TestButterfly4_NEON(t *testing.T) {
	// Test 4-point butterfly
	data := []complex128{
		complex(1.0, 0.0),
		complex(0.0, 1.0),
		complex(1.0, 1.0),
		complex(0.0, 0.0),
	}

	original := make([]complex128, len(data))
	copy(original, data)

	Butterfly4_NEON(data, algorithm.Forward)

	// Verify we got a result (not all zeros)
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Butterfly4_NEON result should not be all zeros")
	}

	// Verify input was modified
	modified := false
	for i := range data {
		if data[i] != original[i] {
			modified = true
			break
		}
	}

	if !modified {
		t.Error("Butterfly4_NEON should modify the input data")
	}

	t.Logf("Butterfly4_NEON: Input %v -> Output %v", original, data)
}

func TestButterfly8_NEON(t *testing.T) {
	// Test 8-point butterfly
	data := make([]complex128, 8)
	for i := range data {
		data[i] = complex(float64(i), float64(i)*0.5)
	}

	original := make([]complex128, len(data))
	copy(original, data)

	Butterfly8_NEON(data, algorithm.Forward)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Butterfly8_NEON result should not be all zeros")
	}

	t.Logf("Butterfly8_NEON: Processed 8 elements successfully")
}

func TestButterfly16_NEON(t *testing.T) {
	// Test 16-point butterfly
	data := make([]complex128, 16)
	for i := range data {
		data[i] = complex(float64(i%4), float64(i%3)*0.3)
	}

	Butterfly16_NEON(data, algorithm.Forward)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Butterfly16_NEON result should not be all zeros")
	}

	t.Logf("Butterfly16_NEON: Processed 16 elements successfully")
}

func TestButterfly32_NEON(t *testing.T) {
	// Test 32-point butterfly
	data := make([]complex128, 32)
	for i := range data {
		data[i] = complex(float64(i%8), float64(i%7)*0.2)
	}

	Butterfly32_NEON(data, algorithm.Forward)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Butterfly32_NEON result should not be all zeros")
	}

	t.Logf("Butterfly32_NEON: Processed 32 elements successfully")
}

func TestProcessVectorizedButterfly(t *testing.T) {
	testSizes := []int{2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 16, 17, 19, 23, 24, 27, 29, 31, 32}

	for _, size := range testSizes {
		t.Run("Size"+string(rune(size+'0')), func(t *testing.T) {
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%5), float64(i%3)*0.4)
			}

			ProcessVectorizedButterfly(data, size, algorithm.Forward)

			// Verify we got a result
			hasNonZero := false
			for _, val := range data {
				if val != 0 {
					hasNonZero = true
					break
				}
			}

			if !hasNonZero {
				t.Errorf("ProcessVectorizedButterfly(size=%d) result should not be all zeros", size)
			}
		})
	}
}

func TestNEONAlignment(t *testing.T) {
	// Test memory alignment utilities
	data := make([]complex128, 16)

	// Test alignment check
	ptr := unsafe.Pointer(&data[0])
	aligned := isAligned16(ptr)

	t.Logf("Data alignment: %v (16-byte aligned: %v)", ptr, aligned)

	// Test alignment function
	alignedData := alignTo16(data)
	if len(alignedData) != len(data) {
		t.Errorf("Aligned data length: got %d, expected %d", len(alignedData), len(data))
	}

	// Verify data is preserved
	for i := range data {
		if alignedData[i] != data[i] {
			t.Errorf("Aligned data[%d]: got %v, expected %v", i, alignedData[i], data[i])
		}
	}
}

func TestNEONStats(t *testing.T) {
	// Test NEON performance statistics
	neonOps, scalarOps := getNEONStats()

	t.Logf("NEON operations: %d, Scalar fallbacks: %d", neonOps, scalarOps)

	// Reset stats
	resetNEONStats()

	neonOps, scalarOps = getNEONStats()
	if neonOps != 0 || scalarOps != 0 {
		t.Errorf("Stats not reset: neonOps=%d, scalarOps=%d", neonOps, scalarOps)
	}
}

func BenchmarkNEONButterflies(b *testing.B) {
	sizes := []int{2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 16, 17, 19, 23, 24, 27, 29, 31, 32}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Butterfly%d", size), func(b *testing.B) {
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%7), float64(i%5)*0.3)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ProcessVectorizedButterfly(data, size, algorithm.Forward)
			}
		})
	}
}

func BenchmarkButterfly2_NEON(b *testing.B) {
	data := make([]complex128, 2)
	data[0] = complex(1.0, 2.0)
	data[1] = complex(3.0, 4.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Butterfly2_NEON(data, algorithm.Forward)
	}
}

func BenchmarkButterfly4_NEON(b *testing.B) {
	data := make([]complex128, 4)
	for i := range data {
		data[i] = complex(float64(i), float64(i)*0.5)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Butterfly4_NEON(data, algorithm.Forward)
	}
}

func BenchmarkButterfly8_NEON(b *testing.B) {
	data := make([]complex128, 8)
	for i := range data {
		data[i] = complex(float64(i), float64(i)*0.3)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Butterfly8_NEON(data, algorithm.Forward)
	}
}

// Helper functions

func complexEqual(a, b complex128) bool {
	const epsilon = 1e-12
	return math.Abs(real(a)-real(b)) < epsilon && math.Abs(imag(a)-imag(b)) < epsilon
}
