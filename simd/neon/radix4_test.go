//go:build arm64

package neon

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/10d9e/gofft/algorithm"
)

func TestNEONRadix4(t *testing.T) {
	testCases := []struct {
		size int
		name string
	}{
		{64, "Radix4_64"},
		{128, "Radix4_128"},
		{256, "Radix4_256"},
		{512, "Radix4_512"},
		{1024, "Radix4_1024"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test data
			data := make([]complex128, tc.size)
			for i := range data {
				data[i] = complex(float64(i%13), float64(i%7)*0.3)
			}

			// Apply NEON Radix-4
			ProcessVectorizedRadix4(data, tc.size, algorithm.Forward)

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

func TestRadix4_64_NEON(t *testing.T) {
	// Test 64-point Radix-4 FFT
	data := make([]complex128, 64)
	for i := range data {
		data[i] = complex(float64(i%8), float64(i%5)*0.4)
	}

	original := make([]complex128, len(data))
	copy(original, data)

	Radix4_64_NEON(data, algorithm.Forward)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Radix4_64_NEON result should not be all zeros")
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
		t.Error("Radix4_64_NEON should modify the input data")
	}

	t.Logf("Radix4_64_NEON: Processed 64 elements successfully")
}

func TestRadix4_128_NEON(t *testing.T) {
	// Test 128-point Radix-4 FFT
	data := make([]complex128, 128)
	for i := range data {
		data[i] = complex(float64(i%16), float64(i%9)*0.2)
	}

	Radix4_128_NEON(data, algorithm.Forward)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Radix4_128_NEON result should not be all zeros")
	}

	t.Logf("Radix4_128_NEON: Processed 128 elements successfully")
}

func TestRadix4_256_NEON(t *testing.T) {
	// Test 256-point Radix-4 FFT
	data := make([]complex128, 256)
	for i := range data {
		data[i] = complex(float64(i%32), float64(i%11)*0.15)
	}

	Radix4_256_NEON(data, algorithm.Forward)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Radix4_256_NEON result should not be all zeros")
	}

	t.Logf("Radix4_256_NEON: Processed 256 elements successfully")
}

func TestRadix4_512_NEON(t *testing.T) {
	// Test 512-point Radix-4 FFT
	data := make([]complex128, 512)
	for i := range data {
		data[i] = complex(float64(i%64), float64(i%13)*0.1)
	}

	Radix4_512_NEON(data, algorithm.Forward)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Radix4_512_NEON result should not be all zeros")
	}

	t.Logf("Radix4_512_NEON: Processed 512 elements successfully")
}

func TestRadix4_1024_NEON(t *testing.T) {
	// Test 1024-point Radix-4 FFT
	data := make([]complex128, 1024)
	for i := range data {
		data[i] = complex(float64(i%128), float64(i%17)*0.05)
	}

	Radix4_1024_NEON(data, algorithm.Forward)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Radix4_1024_NEON result should not be all zeros")
	}

	t.Logf("Radix4_1024_NEON: Processed 1024 elements successfully")
}

func TestProcessVectorizedRadix4(t *testing.T) {
	testSizes := []int{64, 128, 256, 512, 1024}

	for _, size := range testSizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%19), float64(i%7)*0.3)
			}

			ProcessVectorizedRadix4(data, size, algorithm.Forward)

			// Verify we got a result
			hasNonZero := false
			for _, val := range data {
				if val != 0 {
					hasNonZero = true
					break
				}
			}

			if !hasNonZero {
				t.Errorf("ProcessVectorizedRadix4(size=%d) result should not be all zeros", size)
			}
		})
	}
}

func TestRadix4HelperFunctions(t *testing.T) {
	// Test isPowerOf4
	powerOf4Tests := []struct {
		n      int
		expect bool
	}{
		{1, true},
		{4, true},
		{16, true},
		{64, true},
		{256, true},
		{1024, true},
		{2, false},
		{8, false},
		{32, false},
		{128, false},
		{512, false},
		{0, false},
		{-4, false},
	}

	for _, test := range powerOf4Tests {
		result := isPowerOf4(test.n)
		if result != test.expect {
			t.Errorf("isPowerOf4(%d): got %v, expected %v", test.n, result, test.expect)
		}
	}

	// Test log4
	log4Tests := []struct {
		n      int
		expect int
	}{
		{1, 0},
		{4, 1},
		{16, 2},
		{64, 3},
		{256, 4},
		{1024, 5},
	}

	for _, test := range log4Tests {
		result := log4(test.n)
		if result != test.expect {
			t.Errorf("log4(%d): got %d, expected %d", test.n, result, test.expect)
		}
	}
}

func TestRadix4TwiddleGeneration(t *testing.T) {
	// Test twiddle factor generation
	twiddles := generateRadix4Twiddles(64, 1) // Forward direction

	// Should have 3 * 16 = 48 twiddle factors
	expectedLen := 3 * 16
	if len(twiddles) != expectedLen {
		t.Errorf("generateRadix4Twiddles(64, 1): got %d twiddles, expected %d", len(twiddles), expectedLen)
	}

	// Verify twiddle factors are not all zeros
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

	// Test inverse direction
	twiddlesInv := generateRadix4Twiddles(64, -1)
	if len(twiddlesInv) != expectedLen {
		t.Errorf("generateRadix4Twiddles(64, -1): got %d twiddles, expected %d", len(twiddlesInv), expectedLen)
	}
}

func TestNEONRadix4Alignment(t *testing.T) {
	// Test memory alignment utilities
	data := make([]complex128, 64)

	// Test alignment check
	ptr := unsafe.Pointer(&data[0])
	aligned32 := isAligned32(ptr)

	t.Logf("Data alignment: %v (32-byte aligned: %v)", ptr, aligned32)

	// Test alignment function
	alignedData := alignTo32(data)
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

func TestRadix4Constructor(t *testing.T) {
	// Test Radix4_NEON constructor
	radix4 := NewRadix4_NEON(64, 1)
	if radix4 == nil {
		t.Error("NewRadix4_NEON(64, 1) should not return nil")
	}

	if radix4.length != 64 {
		t.Errorf("Radix4 length: got %d, expected 64", radix4.length)
	}

	if radix4.direction != 1 {
		t.Errorf("Radix4 direction: got %d, expected 1", radix4.direction)
	}

	// Test invalid sizes
	invalidRadix4 := NewRadix4_NEON(32, 1) // Not a power of 4
	if invalidRadix4 != nil {
		t.Error("NewRadix4_NEON(32, 1) should return nil for non-power-of-4 size")
	}

	invalidRadix4 = NewRadix4_NEON(2, 1) // Too small
	if invalidRadix4 != nil {
		t.Error("NewRadix4_NEON(2, 1) should return nil for size < 4")
	}
}

func BenchmarkNEONRadix4(b *testing.B) {
	sizes := []int{64, 128, 256, 512, 1024}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Radix4_%d", size), func(b *testing.B) {
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%23), float64(i%11)*0.3)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ProcessVectorizedRadix4(data, size, algorithm.Forward)
			}
		})
	}
}

func BenchmarkRadix4_64_NEON(b *testing.B) {
	data := make([]complex128, 64)
	for i := range data {
		data[i] = complex(float64(i%8), float64(i%5)*0.4)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Radix4_64_NEON(data, algorithm.Forward)
	}
}

func BenchmarkRadix4_128_NEON(b *testing.B) {
	data := make([]complex128, 128)
	for i := range data {
		data[i] = complex(float64(i%16), float64(i%9)*0.2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Radix4_128_NEON(data, algorithm.Forward)
	}
}

func BenchmarkRadix4_256_NEON(b *testing.B) {
	data := make([]complex128, 256)
	for i := range data {
		data[i] = complex(float64(i%32), float64(i%11)*0.15)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Radix4_256_NEON(data, algorithm.Forward)
	}
}

func BenchmarkRadix4_512_NEON(b *testing.B) {
	data := make([]complex128, 512)
	for i := range data {
		data[i] = complex(float64(i%64), float64(i%13)*0.1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Radix4_512_NEON(data, algorithm.Forward)
	}
}

func BenchmarkRadix4_1024_NEON(b *testing.B) {
	data := make([]complex128, 1024)
	for i := range data {
		data[i] = complex(float64(i%128), float64(i%17)*0.05)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Radix4_1024_NEON(data, algorithm.Forward)
	}
}

// Helper functions
// complexEqual is defined in butterflies_test.go
