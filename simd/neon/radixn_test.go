//go:build arm64

package neon

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestNEONRadixN(t *testing.T) {
	testCases := []struct {
		size int
		name string
	}{
		{6, "RadixN_6"},
		{10, "RadixN_10"},
		{12, "RadixN_12"},
		{15, "RadixN_15"},
		{18, "RadixN_18"},
		{20, "RadixN_20"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test data
			data := make([]complex128, tc.size)
			for i := range data {
				data[i] = complex(float64(i%17), float64(i%11)*0.3)
			}

			// Apply NEON RadixN
			ProcessVectorizedRadixN(data, tc.size)

			// Verify we got some result (not all zeros)
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

			t.Logf("%s: Processed %d elements successfully", tc.name, tc.size)
		})
	}
}

func TestRadixN_6_NEON(t *testing.T) {
	// Test 6-point RadixN FFT (2×3)
	data := make([]complex128, 6)
	for i := range data {
		data[i] = complex(float64(i%3), float64(i%2)*0.4)
	}

	original := make([]complex128, len(data))
	copy(original, data)

	RadixN_6_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("RadixN_6_NEON result should not be all zeros")
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
		t.Error("RadixN_6_NEON should modify the input data")
	}

	t.Logf("RadixN_6_NEON: Processed 6 elements successfully")
}

func TestRadixN_10_NEON(t *testing.T) {
	// Test 10-point RadixN FFT (2×5)
	data := make([]complex128, 10)
	for i := range data {
		data[i] = complex(float64(i%5), float64(i%3)*0.2)
	}

	RadixN_10_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("RadixN_10_NEON result should not be all zeros")
	}

	t.Logf("RadixN_10_NEON: Processed 10 elements successfully")
}

func TestRadixN_12_NEON(t *testing.T) {
	// Test 12-point RadixN FFT (3×4)
	data := make([]complex128, 12)
	for i := range data {
		data[i] = complex(float64(i%4), float64(i%3)*0.3)
	}

	RadixN_12_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("RadixN_12_NEON result should not be all zeros")
	}

	t.Logf("RadixN_12_NEON: Processed 12 elements successfully")
}

func TestRadixN_15_NEON(t *testing.T) {
	// Test 15-point RadixN FFT (3×5)
	data := make([]complex128, 15)
	for i := range data {
		data[i] = complex(float64(i%5), float64(i%3)*0.25)
	}

	RadixN_15_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("RadixN_15_NEON result should not be all zeros")
	}

	t.Logf("RadixN_15_NEON: Processed 15 elements successfully")
}

func TestRadixN_18_NEON(t *testing.T) {
	// Test 18-point RadixN FFT (2×9)
	data := make([]complex128, 18)
	for i := range data {
		data[i] = complex(float64(i%9), float64(i%2)*0.15)
	}

	RadixN_18_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("RadixN_18_NEON result should not be all zeros")
	}

	t.Logf("RadixN_18_NEON: Processed 18 elements successfully")
}

func TestRadixN_20_NEON(t *testing.T) {
	// Test 20-point RadixN FFT (4×5)
	data := make([]complex128, 20)
	for i := range data {
		data[i] = complex(float64(i%5), float64(i%4)*0.1)
	}

	RadixN_20_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("RadixN_20_NEON result should not be all zeros")
	}

	t.Logf("RadixN_20_NEON: Processed 20 elements successfully")
}

func TestProcessVectorizedRadixN(t *testing.T) {
	testSizes := []int{6, 10, 12, 15, 18, 20}

	for _, size := range testSizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%19), float64(i%7)*0.3)
			}

			ProcessVectorizedRadixN(data, size)

			// Verify we got a result
			hasNonZero := false
			for _, val := range data {
				if val != 0 {
					hasNonZero = true
					break
				}
			}

			if !hasNonZero {
				t.Errorf("ProcessVectorizedRadixN(size=%d) result should not be all zeros", size)
			}
		})
	}
}

func TestRadixNHelperFunctions(t *testing.T) {
	// Test generateRadixNTwiddles
	factors := []RadixFactor{Factor2, Factor3}
	twiddles := generateRadixNTwiddles(6, factors, 1)

	// Should have some twiddles
	if len(twiddles) == 0 {
		t.Error("generateRadixNTwiddles should generate some twiddles")
	}

	// Verify twiddles are not all zeros
	hasNonZero := false
	for _, tw := range twiddles {
		if tw != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Twiddle factors should not be all zeros")
	}

	// Test createTransposeFactors
	transposeFactors := createTransposeFactors(factors)
	if len(transposeFactors) == 0 {
		t.Error("createTransposeFactors should create some factors")
	}

	// Test reverseRemainders
	result := reverseRemainders(5, transposeFactors)
	if result < 0 {
		t.Error("reverseRemainders should return non-negative result")
	}
}

func TestRadixNConstructor(t *testing.T) {
	// Test RadixN_NEON constructor
	factors := []RadixFactor{Factor2, Factor3}
	radixN := NewRadixN_NEON(6, 1, factors)
	if radixN == nil {
		t.Error("NewRadixN_NEON(6, 1, factors) should not return nil")
	}

	if radixN.length != 6 {
		t.Errorf("RadixN length: got %d, expected 6", radixN.length)
	}

	if radixN.direction != 1 {
		t.Errorf("RadixN direction: got %d, expected 1", radixN.direction)
	}

	// Test invalid sizes
	invalidRadixN := NewRadixN_NEON(7, 1, factors) // Not divisible by factor product
	if invalidRadixN != nil {
		t.Error("NewRadixN_NEON(7, 1, factors) should return nil for invalid size")
	}

	invalidRadixN = NewRadixN_NEON(1, 1, factors) // Too small
	if invalidRadixN != nil {
		t.Error("NewRadixN_NEON(1, 1, factors) should return nil for size < 2")
	}
}

func TestNEONRadixNAlignment(t *testing.T) {
	// Test memory alignment utilities
	data := make([]complex128, 20)

	// Test alignment check
	ptr := unsafe.Pointer(&data[0])
	aligned64 := isAligned64(ptr)

	t.Logf("Data alignment: %v (64-byte aligned: %v)", ptr, aligned64)

	// Test alignment function
	alignedData := alignTo64(data)
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

func TestRadixNFactorOperations(t *testing.T) {
	// Test factor operations
	factors := []RadixFactor{Factor2, Factor3, Factor4}

	// Test transpose factor creation
	transposeFactors := createTransposeFactors(factors)
	if len(transposeFactors) == 0 {
		t.Error("createTransposeFactors should create factors")
	}

	// Test remainder reversal
	for i := 0; i < 10; i++ {
		result := reverseRemainders(i, transposeFactors)
		if result < 0 {
			t.Errorf("reverseRemainders(%d) should return non-negative result", i)
		}
	}
}

func BenchmarkNEONRadixN(b *testing.B) {
	sizes := []int{6, 10, 12, 15, 18, 20}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("RadixN_%d", size), func(b *testing.B) {
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%23), float64(i%11)*0.3)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ProcessVectorizedRadixN(data, size)
			}
		})
	}
}

func BenchmarkRadixN_6_NEON(b *testing.B) {
	data := make([]complex128, 6)
	for i := range data {
		data[i] = complex(float64(i%3), float64(i%2)*0.4)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RadixN_6_NEON(data)
	}
}

func BenchmarkRadixN_10_NEON(b *testing.B) {
	data := make([]complex128, 10)
	for i := range data {
		data[i] = complex(float64(i%5), float64(i%3)*0.2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RadixN_10_NEON(data)
	}
}

func BenchmarkRadixN_12_NEON(b *testing.B) {
	data := make([]complex128, 12)
	for i := range data {
		data[i] = complex(float64(i%4), float64(i%3)*0.3)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RadixN_12_NEON(data)
	}
}

func BenchmarkRadixN_15_NEON(b *testing.B) {
	data := make([]complex128, 15)
	for i := range data {
		data[i] = complex(float64(i%5), float64(i%3)*0.25)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RadixN_15_NEON(data)
	}
}

func BenchmarkRadixN_18_NEON(b *testing.B) {
	data := make([]complex128, 18)
	for i := range data {
		data[i] = complex(float64(i%9), float64(i%2)*0.15)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RadixN_18_NEON(data)
	}
}

func BenchmarkRadixN_20_NEON(b *testing.B) {
	data := make([]complex128, 20)
	for i := range data {
		data[i] = complex(float64(i%5), float64(i%4)*0.1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RadixN_20_NEON(data)
	}
}
