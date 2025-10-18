//go:build arm64

package neon

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestNEONRaders(t *testing.T) {
	testCases := []struct {
		size int
		name string
	}{
		{37, "Raders_37"},
		{41, "Raders_41"},
		{43, "Raders_43"},
		{47, "Raders_47"},
		{53, "Raders_53"},
		{59, "Raders_59"},
		{61, "Raders_61"},
		{67, "Raders_67"},
		{71, "Raders_71"},
		{73, "Raders_73"},
		{79, "Raders_79"},
		{83, "Raders_83"},
		{89, "Raders_89"},
		{97, "Raders_97"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test data
			data := make([]complex128, tc.size)
			for i := range data {
				data[i] = complex(float64(i%23), float64(i%17)*0.3)
			}

			// Apply NEON Rader's
			ProcessVectorizedRaders(data, tc.size)

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

func TestRaders_37_NEON(t *testing.T) {
	// Test 37-point Rader's FFT
	data := make([]complex128, 37)
	for i := range data {
		data[i] = complex(float64(i%7), float64(i%5)*0.4)
	}

	original := make([]complex128, len(data))
	copy(original, data)

	Raders_37_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Raders_37_NEON result should not be all zeros")
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
		t.Error("Raders_37_NEON should modify the input data")
	}

	t.Logf("Raders_37_NEON: Processed 37 elements successfully")
}

func TestRaders_41_NEON(t *testing.T) {
	// Test 41-point Rader's FFT
	data := make([]complex128, 41)
	for i := range data {
		data[i] = complex(float64(i%11), float64(i%7)*0.2)
	}

	Raders_41_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Raders_41_NEON result should not be all zeros")
	}

	t.Logf("Raders_41_NEON: Processed 41 elements successfully")
}

func TestRaders_43_NEON(t *testing.T) {
	// Test 43-point Rader's FFT
	data := make([]complex128, 43)
	for i := range data {
		data[i] = complex(float64(i%13), float64(i%11)*0.3)
	}

	Raders_43_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Raders_43_NEON result should not be all zeros")
	}

	t.Logf("Raders_43_NEON: Processed 43 elements successfully")
}

func TestRaders_47_NEON(t *testing.T) {
	// Test 47-point Rader's FFT
	data := make([]complex128, 47)
	for i := range data {
		data[i] = complex(float64(i%17), float64(i%13)*0.25)
	}

	Raders_47_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Raders_47_NEON result should not be all zeros")
	}

	t.Logf("Raders_47_NEON: Processed 47 elements successfully")
}

func TestRaders_53_NEON(t *testing.T) {
	// Test 53-point Rader's FFT
	data := make([]complex128, 53)
	for i := range data {
		data[i] = complex(float64(i%19), float64(i%17)*0.15)
	}

	Raders_53_NEON(data)

	// Verify we got a result
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Raders_53_NEON result should not be all zeros")
	}

	t.Logf("Raders_53_NEON: Processed 53 elements successfully")
}

func TestProcessVectorizedRaders(t *testing.T) {
	testSizes := []int{37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

	for _, size := range testSizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%29), float64(i%19)*0.3)
			}

			ProcessVectorizedRaders(data, size)

			// Verify we got a result
			hasNonZero := false
			for _, val := range data {
				if val != 0 {
					hasNonZero = true
					break
				}
			}

			if !hasNonZero {
				t.Errorf("ProcessVectorizedRaders(size=%d) result should not be all zeros", size)
			}
		})
	}
}

func TestRadersHelperFunctions(t *testing.T) {
	// Test isPrime
	primeTests := []struct {
		n      int
		expect bool
	}{
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
		{9, false},
		{10, false},
		{11, true},
		{37, true},
		{41, true},
		{43, true},
		{47, true},
		{53, true},
		{59, true},
		{61, true},
		{67, true},
		{71, true},
		{73, true},
		{79, true},
		{83, true},
		{89, true},
		{97, true},
	}

	for _, test := range primeTests {
		result := isPrime(test.n)
		if result != test.expect {
			t.Errorf("isPrime(%d): got %v, expected %v", test.n, result, test.expect)
		}
	}

	// Test findPrimitiveRoot
	rootTests := []struct {
		prime uint64
		valid bool
	}{
		{2, true},
		{3, true},
		{5, true},
		{7, true},
		{11, true},
		{13, true},
		{17, true},
		{19, true},
		{23, true},
		{29, true},
		{31, true},
		{37, true},
		{41, true},
		{43, true},
		{47, true},
		{53, true},
		{59, true},
		{61, true},
		{67, true},
		{71, true},
		{73, true},
		{79, true},
		{83, true},
		{89, true},
		{97, true},
	}

	for _, test := range rootTests {
		root := findPrimitiveRoot(test.prime)
		if test.valid && root == 0 {
			t.Errorf("findPrimitiveRoot(%d): expected valid root, got 0", test.prime)
		}
		if !test.valid && root != 0 {
			t.Errorf("findPrimitiveRoot(%d): expected 0, got %d", test.prime, root)
		}
	}

	// Test modInverse
	inverseTests := []struct {
		a, m, expected uint64
	}{
		{3, 7, 5},    // 3 * 5 = 15 ≡ 1 (mod 7)
		{5, 11, 9},   // 5 * 9 = 45 ≡ 1 (mod 11)
		{7, 13, 2},   // 7 * 2 = 14 ≡ 1 (mod 13)
		{11, 17, 14}, // 11 * 14 = 154 ≡ 1 (mod 17)
	}

	for _, test := range inverseTests {
		result := modInverse(test.a, test.m)
		if result != test.expected {
			t.Errorf("modInverse(%d, %d): got %d, expected %d", test.a, test.m, result, test.expected)
		}
	}
}

func TestRadersConstructor(t *testing.T) {
	// Test Raders_NEON constructor
	raders := NewRaders_NEON(37, 1)
	if raders == nil {
		t.Error("NewRaders_NEON(37, 1) should not return nil")
	}

	if raders.length != 37 {
		t.Errorf("Raders length: got %d, expected 37", raders.length)
	}

	if raders.direction != 1 {
		t.Errorf("Raders direction: got %d, expected 1", raders.direction)
	}

	// Test invalid sizes
	invalidRaders := NewRaders_NEON(36, 1) // Not prime
	if invalidRaders != nil {
		t.Error("NewRaders_NEON(36, 1) should return nil for non-prime size")
	}

	invalidRaders = NewRaders_NEON(1, 1) // Too small
	if invalidRaders != nil {
		t.Error("NewRaders_NEON(1, 1) should return nil for size < 2")
	}
}

func TestNEONRadersAlignment(t *testing.T) {
	// Test memory alignment utilities
	data := make([]complex128, 97)

	// Test alignment check
	ptr := unsafe.Pointer(&data[0])
	aligned128 := isAligned128(ptr)

	t.Logf("Data alignment: %v (128-byte aligned: %v)", ptr, aligned128)

	// Test alignment function
	alignedData := alignTo128(data)
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

func TestRadersDataGeneration(t *testing.T) {
	// Test generateRadersData
	primitiveRoot := findPrimitiveRoot(37)
	if primitiveRoot == 0 {
		t.Fatal("Could not find primitive root for 37")
	}

	primitiveRootInverse := modInverse(primitiveRoot, 37)
	if primitiveRootInverse == 0 {
		t.Fatal("Could not find primitive root inverse for 37")
	}

	data := generateRadersData(37, primitiveRoot, primitiveRootInverse, 1)
	if len(data) != 36 {
		t.Errorf("generateRadersData(37): got %d elements, expected 36", len(data))
	}

	// Verify data is not all zeros
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("Generated Rader's data should not be all zeros")
	}
}

func BenchmarkNEONRaders(b *testing.B) {
	sizes := []int{37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Raders_%d", size), func(b *testing.B) {
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%31), float64(i%23)*0.3)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ProcessVectorizedRaders(data, size)
			}
		})
	}
}

func BenchmarkRaders_37_NEON(b *testing.B) {
	data := make([]complex128, 37)
	for i := range data {
		data[i] = complex(float64(i%7), float64(i%5)*0.4)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Raders_37_NEON(data)
	}
}

func BenchmarkRaders_41_NEON(b *testing.B) {
	data := make([]complex128, 41)
	for i := range data {
		data[i] = complex(float64(i%11), float64(i%7)*0.2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Raders_41_NEON(data)
	}
}

func BenchmarkRaders_43_NEON(b *testing.B) {
	data := make([]complex128, 43)
	for i := range data {
		data[i] = complex(float64(i%13), float64(i%11)*0.3)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Raders_43_NEON(data)
	}
}

func BenchmarkRaders_47_NEON(b *testing.B) {
	data := make([]complex128, 47)
	for i := range data {
		data[i] = complex(float64(i%17), float64(i%13)*0.25)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Raders_47_NEON(data)
	}
}

func BenchmarkRaders_53_NEON(b *testing.B) {
	data := make([]complex128, 53)
	for i := range data {
		data[i] = complex(float64(i%19), float64(i%17)*0.15)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Raders_53_NEON(data)
	}
}
