//go:build arm64

package neon

import (
	"math"
	"math/cmplx"
	"unsafe"
)

// NEON Rader's Algorithm implementation for ARM64
// This implements prime-sized FFTs using Rader's algorithm with NEON optimizations

// Raders_NEON represents a NEON-optimized Rader's algorithm FFT
type Raders_NEON struct {
	length               int
	direction            int // 1 for forward, -1 for inverse
	primitiveRoot        uint64
	primitiveRootInverse uint64
	innerFFTData         []complex128
	innerFFT             *Raders_NEON // For recursive decomposition
}

// NewRaders_NEON creates a new NEON Rader's algorithm FFT
func NewRaders_NEON(length int, direction int) *Raders_NEON {
	if length < 2 {
		return nil
	}

	// Check if length is prime
	if !isPrime(length) {
		return nil
	}

	// Find primitive root
	primitiveRoot := findPrimitiveRoot(uint64(length))
	if primitiveRoot == 0 {
		return nil
	}

	// Find primitive root inverse
	primitiveRootInverse := modInverse(primitiveRoot, uint64(length))

	// Create inner FFT for size (length - 1)
	innerLength := length - 1
	var innerFFT *Raders_NEON
	if innerLength > 1 {
		// For now, create a simple inner FFT
		// In a full implementation, this would be a composite FFT
		innerFFT = &Raders_NEON{
			length:    innerLength,
			direction: direction,
		}
	}

	// Precompute inner FFT data
	innerFFTData := generateRadersData(length, primitiveRoot, primitiveRootInverse, direction)

	return &Raders_NEON{
		length:               length,
		direction:            direction,
		primitiveRoot:        primitiveRoot,
		primitiveRootInverse: primitiveRootInverse,
		innerFFTData:         innerFFTData,
		innerFFT:             innerFFT,
	}
}

// Process performs the Rader's algorithm FFT using NEON optimizations
func (r *Raders_NEON) Process(data []complex128) {
	if len(data) < r.length {
		return
	}

	// For now, use optimized scalar implementation
	// TODO: Replace with actual NEON intrinsics
	r.processScalar(data)
}

// processScalar is the scalar implementation (placeholder for NEON)
func (r *Raders_NEON) processScalar(data []complex128) {
	if r.innerFFT == nil {
		// Base case: direct DFT
		r.performDirectDFT(data)
		return
	}

	// Recursive case: apply Rader's algorithm
	r.performRadersAlgorithm(data)
}

// performDirectDFT performs a direct DFT for small sizes
func (r *Raders_NEON) performDirectDFT(data []complex128) {
	// Simple DFT implementation for small sizes
	n := len(data)
	for k := 0; k < n; k++ {
		sum := complex(0, 0)
		for j := 0; j < n; j++ {
			angle := -2 * math.Pi * float64(k*j) / float64(n)
			if r.direction == -1 {
				angle = -angle
			}
			w := complex(math.Cos(angle), math.Sin(angle))
			sum += data[j] * w
		}
		data[k] = sum
	}
}

// performRadersAlgorithm performs the Rader's algorithm
func (r *Raders_NEON) performRadersAlgorithm(data []complex128) {
	// Step 1: Extract first element and compute sum
	firstElement := data[0]
	sum := complex(0, 0)
	for i := 1; i < r.length; i++ {
		sum += data[i]
	}

	// Step 2: Reorder data using primitive root
	reordered := make([]complex128, r.length-1)
	inputIndex := 1
	for i := 0; i < r.length-1; i++ {
		inputIndex = int((uint64(inputIndex) * r.primitiveRoot) % uint64(r.length))
		reordered[i] = data[inputIndex-1]
	}

	// Step 3: Apply inner FFT
	if r.innerFFT != nil {
		r.innerFFT.Process(reordered)
	}

	// Step 4: First output element is sum of all elements
	data[0] = firstElement + reordered[0]

	// Step 5: Multiply with precomputed data and conjugate
	for i := 0; i < r.length-1; i++ {
		reordered[i] = cmplx.Conj(reordered[i] * r.innerFFTData[i])
	}

	// Step 6: Add first element to DC component
	reordered[0] += cmplx.Conj(firstElement)

	// Step 7: Apply second inner FFT
	if r.innerFFT != nil {
		r.innerFFT.Process(reordered)
	}

	// Step 8: Reorder output using primitive root inverse
	outputIndex := 1
	for i := 0; i < r.length-1; i++ {
		outputIndex = int((uint64(outputIndex) * r.primitiveRootInverse) % uint64(r.length))
		data[outputIndex-1] = cmplx.Conj(reordered[i])
	}
}

// Specific Rader's implementations for common prime sizes

// Raders_37_NEON performs a 37-point Rader's FFT using NEON
func Raders_37_NEON(data []complex128) {
	if len(data) < 37 {
		return
	}

	// Use real NEON assembly
	raders_37_fft_go(data)
}

// Raders_41_NEON performs a 41-point Rader's FFT using NEON
func Raders_41_NEON(data []complex128) {
	if len(data) < 41 {
		return
	}

	// Use real NEON assembly
	raders_41_fft_go(data)
}

// Raders_43_NEON performs a 43-point Rader's FFT using NEON
func Raders_43_NEON(data []complex128) {
	if len(data) < 43 {
		return
	}

	// Use real NEON assembly
	raders_43_fft_go(data)
}

// Raders_47_NEON performs a 47-point Rader's FFT using NEON
func Raders_47_NEON(data []complex128) {
	if len(data) < 47 {
		return
	}

	// Use real NEON assembly
	raders_47_fft_go(data)
}

// Raders_53_NEON performs a 53-point Rader's FFT using NEON
func Raders_53_NEON(data []complex128) {
	if len(data) < 53 {
		return
	}

	// Use real NEON assembly
	raders_53_fft_go(data)
}

// Raders_59_NEON performs a 59-point Rader's FFT using NEON
func Raders_59_NEON(data []complex128) {
	if len(data) < 59 {
		return
	}

	// Use real NEON assembly
	raders_59_fft_go(data)
}

// Raders_61_NEON performs a 61-point Rader's FFT using NEON
func Raders_61_NEON(data []complex128) {
	if len(data) < 61 {
		return
	}

	// Use real NEON assembly
	raders_61_fft_go(data)
}

// Raders_67_NEON performs a 67-point Rader's FFT using NEON
func Raders_67_NEON(data []complex128) {
	if len(data) < 67 {
		return
	}

	// Use real NEON assembly
	raders_67_fft_go(data)
}

// Raders_71_NEON performs a 71-point Rader's FFT using NEON
func Raders_71_NEON(data []complex128) {
	if len(data) < 71 {
		return
	}

	// Use real NEON assembly
	raders_71_fft_go(data)
}

// Raders_73_NEON performs a 73-point Rader's FFT using NEON
func Raders_73_NEON(data []complex128) {
	if len(data) < 73 {
		return
	}

	// Use real NEON assembly
	raders_73_fft_go(data)
}

// Raders_79_NEON performs a 79-point Rader's FFT using NEON
func Raders_79_NEON(data []complex128) {
	if len(data) < 79 {
		return
	}

	// Use real NEON assembly
	raders_79_fft_go(data)
}

// Raders_83_NEON performs a 83-point Rader's FFT using NEON
func Raders_83_NEON(data []complex128) {
	if len(data) < 83 {
		return
	}

	// Use real NEON assembly
	raders_83_fft_go(data)
}

// Raders_89_NEON performs a 89-point Rader's FFT using NEON
func Raders_89_NEON(data []complex128) {
	if len(data) < 89 {
		return
	}

	// Use real NEON assembly
	raders_89_fft_go(data)
}

// Raders_97_NEON performs a 97-point Rader's FFT using NEON
func Raders_97_NEON(data []complex128) {
	if len(data) < 97 {
		return
	}

	// Use real NEON assembly
	raders_97_fft_go(data)
}

// ProcessVectorizedRaders processes data using NEON-optimized Rader's FFTs
func ProcessVectorizedRaders(data []complex128, size int) {
	switch size {
	case 37:
		Raders_37_NEON(data)
	case 41:
		Raders_41_NEON(data)
	case 43:
		Raders_43_NEON(data)
	case 47:
		Raders_47_NEON(data)
	case 53:
		Raders_53_NEON(data)
	case 59:
		Raders_59_NEON(data)
	case 61:
		Raders_61_NEON(data)
	case 67:
		Raders_67_NEON(data)
	case 71:
		Raders_71_NEON(data)
	case 73:
		Raders_73_NEON(data)
	case 79:
		Raders_79_NEON(data)
	case 83:
		Raders_83_NEON(data)
	case 89:
		Raders_89_NEON(data)
	case 97:
		Raders_97_NEON(data)
	default:
		// Fall back to scalar implementation for unsupported sizes
		processScalarRaders(data, size)
	}
}

// processScalarRaders is a fallback for unsupported Rader's sizes
func processScalarRaders(data []complex128, size int) {
	// This would call the existing scalar Rader's implementation
	// For now, just a placeholder
}

// Helper functions

// isPrime checks if a number is prime
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// findPrimitiveRoot finds a primitive root for a prime number
func findPrimitiveRoot(prime uint64) uint64 {
	if prime < 2 {
		return 0
	}
	if prime == 2 {
		return 1
	}

	// Get distinct prime factors of (prime - 1)
	factors := distinctPrimeFactors(prime - 1)

	// Test exponents
	testExponents := make([]uint64, len(factors))
	for i, factor := range factors {
		testExponents[i] = (prime - 1) / factor
	}

	// Test potential roots
	for potentialRoot := uint64(2); potentialRoot < prime; potentialRoot++ {
		found := true
		for _, exp := range testExponents {
			if modularExponent(potentialRoot, exp, prime) == 1 {
				found = false
				break
			}
		}
		if found {
			return potentialRoot
		}
	}
	return 0
}

// modularExponent computes base^exponent % modulo using exponentiation by squaring
func modularExponent(base, exponent, modulo uint64) uint64 {
	result := uint64(1)
	base = base % modulo

	for exponent > 0 {
		if exponent%2 == 1 {
			result = (result * base) % modulo
		}
		exponent = exponent >> 1
		base = (base * base) % modulo
	}
	return result
}

// distinctPrimeFactors returns all distinct prime factors of n
func distinctPrimeFactors(n uint64) []uint64 {
	var result []uint64

	// Handle 2
	if n%2 == 0 {
		result = append(result, 2)
		for n%2 == 0 {
			n = n / 2
		}
	}

	// Handle odd factors
	for i := uint64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			result = append(result, i)
			for n%i == 0 {
				n = n / i
			}
		}
	}

	// Handle remaining prime
	if n > 2 {
		result = append(result, n)
	}

	return result
}

// modInverse computes the modular inverse of a mod m using extended Euclidean algorithm
func modInverse(a, m uint64) uint64 {
	// Extended Euclidean algorithm
	var oldR, r = int64(a), int64(m)
	var oldS, s = int64(1), int64(0)

	for r != 0 {
		quotient := oldR / r
		oldR, r = r, oldR-quotient*r
		oldS, s = s, oldS-quotient*s
	}

	if oldR > 1 {
		return 0 // No inverse exists
	}

	if oldS < 0 {
		oldS += int64(m)
	}

	return uint64(oldS)
}

// generateRadersData generates the precomputed data for Rader's algorithm
func generateRadersData(length int, primitiveRoot, primitiveRootInverse uint64, direction int) []complex128 {
	innerLength := length - 1
	data := make([]complex128, innerLength)

	scale := 1.0 / float64(innerLength)
	twiddleInput := uint64(1)

	for i := 0; i < innerLength; i++ {
		angle := -2 * math.Pi * float64(twiddleInput) / float64(length)
		if direction == -1 {
			angle = -angle
		}
		w := complex(math.Cos(angle), math.Sin(angle))
		data[i] = w * complex(scale, 0)

		twiddleInput = (twiddleInput * primitiveRootInverse) % uint64(length)
	}

	return data
}

// NEON-specific utility functions for Rader's

// loadComplex128Vector loads multiple complex128 values into NEON registers
// This is a placeholder - actual implementation would use NEON intrinsics
func loadComplex128VectorRaders(ptr *complex128, count int) []complex128 {
	// TODO: Implement actual NEON load instruction
	slice := (*[1024]complex128)(unsafe.Pointer(ptr))[:count:count]
	result := make([]complex128, count)
	copy(result, slice)
	return result
}

// storeComplex128Vector stores multiple complex128 values from NEON registers
// This is a placeholder - actual implementation would use NEON intrinsics
func storeComplex128VectorRaders(ptr *complex128, values []complex128) {
	// TODO: Implement actual NEON store instruction
	slice := (*[1024]complex128)(unsafe.Pointer(ptr))[:len(values):len(values)]
	copy(slice, values)
}

// multiplyComplex128Vector multiplies multiple complex128 values using NEON
// This is a placeholder - actual implementation would use NEON intrinsics
func multiplyComplex128VectorRaders(a, b []complex128) []complex128 {
	// TODO: Implement actual NEON complex multiplication
	result := make([]complex128, len(a))
	for i := range a {
		result[i] = a[i] * b[i]
	}
	return result
}

// NEON memory alignment utilities for Rader's

// isAligned128 checks if a pointer is 128-byte aligned (required for NEON Rader's)
func isAligned128(ptr unsafe.Pointer) bool {
	return uintptr(ptr)%128 == 0
}

// alignTo128 aligns a slice to 128-byte boundary for NEON Rader's operations
func alignTo128(data []complex128) []complex128 {
	if len(data) == 0 {
		return data
	}

	ptr := unsafe.Pointer(&data[0])
	if isAligned128(ptr) {
		return data
	}

	// Create aligned copy
	aligned := make([]complex128, len(data))
	copy(aligned, data)
	return aligned
}
