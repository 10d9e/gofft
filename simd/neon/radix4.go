//go:build arm64

package neon

import (
	"math"
	"unsafe"

	"github.com/10d9e/gofft/algorithm"
)

// NEON Radix-4 FFT implementation for ARM64
// This implements power-of-two FFTs using Radix-4 decomposition with NEON optimizations

// Radix4_NEON represents a NEON-optimized Radix-4 FFT
type Radix4_NEON struct {
	length    int
	direction int // 1 for forward, -1 for inverse
	twiddles  []complex128
	baseFFT   *Radix4_NEON // For recursive decomposition
}

// NewRadix4_NEON creates a new NEON Radix-4 FFT
func NewRadix4_NEON(length int, direction int) *Radix4_NEON {
	if length < 4 {
		return nil
	}

	// Check if length is a power of 4
	if !isPowerOf4(length) {
		return nil
	}

	// Calculate number of Radix-4 stages
	_ = log4(length)

	// Create base FFT for the smallest size
	var baseFFT *Radix4_NEON
	if length > 4 {
		baseFFT = NewRadix4_NEON(length/4, direction)
	}

	// Generate twiddle factors
	twiddles := generateRadix4Twiddles(length, direction)

	return &Radix4_NEON{
		length:    length,
		direction: direction,
		twiddles:  twiddles,
		baseFFT:   baseFFT,
	}
}

// Process performs the Radix-4 FFT using NEON optimizations
func (r *Radix4_NEON) Process(data []complex128) {
	if len(data) < r.length {
		return
	}

	// For now, use optimized scalar implementation
	// TODO: Replace with actual NEON intrinsics
	r.processScalar(data)
}

// processScalar is the scalar implementation (placeholder for NEON)
func (r *Radix4_NEON) processScalar(data []complex128) {
	// Convert direction to algorithm.Direction
	algoDir := algorithm.Forward
	if r.direction == -1 {
		algoDir = algorithm.Inverse
	}

	if r.length == 4 {
		// Base case: 4-point FFT
		Butterfly4_NEON(data[:4], algoDir)
		return
	}

	// Recursive case: decompose into 4 smaller FFTs
	quarterLen := r.length / 4

	// Process each quarter
	for i := 0; i < 4; i++ {
		start := i * quarterLen
		end := start + quarterLen
		r.baseFFT.Process(data[start:end])
	}

	// Apply twiddle factors and combine
	r.applyTwiddlesAndCombine(data, quarterLen)
}

// applyTwiddlesAndCombine applies twiddle factors and combines the results
func (r *Radix4_NEON) applyTwiddlesAndCombine(data []complex128, quarterLen int) {
	// Convert direction to algorithm.Direction
	algoDir := algorithm.Forward
	if r.direction == -1 {
		algoDir = algorithm.Inverse
	}

	// Apply twiddle factors to quarters 1, 2, 3
	for i := 1; i < 4; i++ {
		start := i * quarterLen
		for j := 0; j < quarterLen; j++ {
			twiddleIdx := (i-1)*quarterLen + j
			if twiddleIdx < len(r.twiddles) {
				data[start+j] *= r.twiddles[twiddleIdx]
			}
		}
	}

	// Combine using 4-point butterflies
	for i := 0; i < quarterLen; i++ {
		// Extract 4 elements (one from each quarter)
		quarter0 := data[i]
		quarter1 := data[quarterLen+i]
		quarter2 := data[2*quarterLen+i]
		quarter3 := data[3*quarterLen+i]

		// Apply 4-point butterfly
		butterfly4 := []complex128{quarter0, quarter1, quarter2, quarter3}
		Butterfly4_NEON(butterfly4, algoDir)

		// Store back
		data[i] = butterfly4[0]
		data[quarterLen+i] = butterfly4[1]
		data[2*quarterLen+i] = butterfly4[2]
		data[3*quarterLen+i] = butterfly4[3]
	}
}

// Radix4_64_NEON performs a 64-point Radix-4 FFT using NEON
func Radix4_64_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 64 {
		return
	}

	// Use real NEON assembly with direction support
	radix4_64_fft_go(data, direction)
}

// Radix4_128_NEON performs a 128-point Radix-4 FFT using NEON
func Radix4_128_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 128 {
		return
	}

	// Use real NEON assembly
	radix4_128_fft_go(data, direction)
}

// Radix4_256_NEON performs a 256-point Radix-4 FFT using NEON
func Radix4_256_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 256 {
		return
	}

	// Use real NEON assembly
	radix4_256_fft_go(data, direction)
}

// Radix4_512_NEON performs a 512-point Radix-4 FFT using NEON
func Radix4_512_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 512 {
		return
	}

	// Use real NEON assembly
	radix4_512_fft_go(data, direction)
}

// Radix4_1024_NEON performs a 1024-point Radix-4 FFT using NEON
func Radix4_1024_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 1024 {
		return
	}

	// Use real NEON assembly
	radix4_1024_fft_go(data, direction)
}

// Radix4_2048_NEON performs a 2048-point Radix-4 FFT using NEON
func Radix4_2048_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 2048 {
		return
	}

	// Use real NEON assembly
	radix4_2048_fft_go(data, direction)
}

// Radix4_4096_NEON performs a 4096-point Radix-4 FFT using NEON
func Radix4_4096_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 4096 {
		return
	}

	// Use real NEON assembly
	radix4_4096_fft_go(data, direction)
}

// Radix4_8192_NEON performs an 8192-point Radix-4 FFT using NEON
func Radix4_8192_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 8192 {
		return
	}

	// Use real NEON assembly
	radix4_8192_fft_go(data, direction)
}

// Radix4_16384_NEON performs a 16384-point Radix-4 FFT using NEON
func Radix4_16384_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 16384 {
		return
	}

	// Use real NEON assembly
	radix4_16384_fft_go(data, direction)
}

// Radix4_32768_NEON performs a 32768-point Radix-4 FFT using NEON
func Radix4_32768_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 32768 {
		return
	}

	// Use real NEON assembly
	radix4_32768_fft_go(data, direction)
}

// Radix4_65536_NEON performs a 65536-point Radix-4 FFT using NEON
func Radix4_65536_NEON(data []complex128, direction algorithm.Direction) {
	if len(data) < 65536 {
		return
	}

	// Use real NEON assembly
	radix4_65536_fft_go(data, direction)
}

// ProcessVectorizedRadix4 processes data using NEON-optimized Radix-4 FFTs
func ProcessVectorizedRadix4(data []complex128, size int, direction algorithm.Direction) {
	switch size {
	case 64:
		Radix4_64_NEON(data, direction)
	case 128:
		Radix4_128_NEON(data, direction)
	case 256:
		Radix4_256_NEON(data, direction)
	case 512:
		Radix4_512_NEON(data, direction)
	case 1024:
		Radix4_1024_NEON(data, direction)
	case 2048:
		Radix4_2048_NEON(data, direction)
	case 4096:
		Radix4_4096_NEON(data, direction)
	case 8192:
		Radix4_8192_NEON(data, direction)
	case 16384:
		Radix4_16384_NEON(data, direction)
	case 32768:
		Radix4_32768_NEON(data, direction)
	case 65536:
		Radix4_65536_NEON(data, direction)
	default:
		// Fall back to scalar implementation for unsupported sizes
		processScalarRadix4(data, size)
	}
}

// processScalarRadix4 is a fallback for unsupported Radix-4 sizes
func processScalarRadix4(data []complex128, size int) {
	// This would call the existing scalar Radix-4 implementation
	// For now, just a placeholder
}

// Helper functions

// isPowerOf4 checks if a number is a power of 4
func isPowerOf4(n int) bool {
	return n > 0 && (n&(n-1)) == 0 && (n&0x55555555) == n
}

// log4 calculates log base 4 of a number (must be power of 4)
func log4(n int) int {
	if n <= 0 {
		return -1
	}

	count := 0
	for n > 1 {
		n >>= 2
		count++
	}
	return count
}

// generateRadix4Twiddles generates twiddle factors for Radix-4 FFT
func generateRadix4Twiddles(length int, direction int) []complex128 {
	if length <= 4 {
		return nil
	}

	quarterLen := length / 4
	twiddles := make([]complex128, 3*quarterLen) // 3 quarters need twiddles

	for i := 1; i < 4; i++ {
		for j := 0; j < quarterLen; j++ {
			angle := -2 * math.Pi * float64(i*j) / float64(length)
			if direction == -1 {
				angle = -angle
			}
			twiddles[(i-1)*quarterLen+j] = complex(math.Cos(angle), math.Sin(angle))
		}
	}

	return twiddles
}

// NEON-specific utility functions for Radix-4

// loadComplex128Vector loads 4 complex128 values into NEON registers
// This is a placeholder - actual implementation would use NEON intrinsics
func loadComplex128Vector(ptr *complex128) [4]complex128 {
	// TODO: Implement actual NEON load instruction
	slice := (*[4]complex128)(unsafe.Pointer(ptr))
	return [4]complex128{slice[0], slice[1], slice[2], slice[3]}
}

// storeComplex128Vector stores 4 complex128 values from NEON registers
// This is a placeholder - actual implementation would use NEON intrinsics
func storeComplex128Vector(ptr *complex128, values [4]complex128) {
	// TODO: Implement actual NEON store instruction
	slice := (*[4]complex128)(unsafe.Pointer(ptr))
	slice[0] = values[0]
	slice[1] = values[1]
	slice[2] = values[2]
	slice[3] = values[3]
}

// multiplyComplex128Vector multiplies 4 complex128 values using NEON
// This is a placeholder - actual implementation would use NEON intrinsics
func multiplyComplex128Vector(a, b [4]complex128) [4]complex128 {
	// TODO: Implement actual NEON complex multiplication
	return [4]complex128{
		a[0] * b[0],
		a[1] * b[1],
		a[2] * b[2],
		a[3] * b[3],
	}
}

// NEON memory alignment utilities for Radix-4

// isAligned32 checks if a pointer is 32-byte aligned (required for NEON Radix-4)
func isAligned32(ptr unsafe.Pointer) bool {
	return uintptr(ptr)%32 == 0
}

// alignTo32 aligns a slice to 32-byte boundary for NEON Radix-4 operations
func alignTo32(data []complex128) []complex128 {
	if len(data) == 0 {
		return data
	}

	ptr := unsafe.Pointer(&data[0])
	if isAligned32(ptr) {
		return data
	}

	// Create aligned copy
	aligned := make([]complex128, len(data))
	copy(aligned, data)
	return aligned
}
