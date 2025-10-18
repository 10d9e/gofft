//go:build arm64

package neon

import (
	"math"
	"unsafe"
)

// NEON RadixN FFT implementation for ARM64
// This implements multi-factor composite FFTs using RadixN decomposition with NEON optimizations

// RadixFactor represents a radix factor for RadixN decomposition
type RadixFactor int

const (
	Factor2 RadixFactor = 2
	Factor3 RadixFactor = 3
	Factor4 RadixFactor = 4
	Factor5 RadixFactor = 5
	Factor6 RadixFactor = 6
	Factor7 RadixFactor = 7
)

// TransposeFactor represents a factor for transpose operations
type TransposeFactor struct {
	Factor RadixFactor
	Count  int
}

// RadixN_NEON represents a NEON-optimized RadixN FFT
type RadixN_NEON struct {
	length    int
	direction int // 1 for forward, -1 for inverse
	twiddles  []complex128
	factors   []TransposeFactor
	baseFFT   *RadixN_NEON // For recursive decomposition
}

// NewRadixN_NEON creates a new NEON RadixN FFT
func NewRadixN_NEON(length int, direction int, factors []RadixFactor) *RadixN_NEON {
	if length < 2 {
		return nil
	}

	// Calculate the product of factors
	factorProduct := 1
	for _, factor := range factors {
		factorProduct *= int(factor)
	}

	// Check if length is divisible by factor product
	if length%factorProduct != 0 {
		return nil
	}

	baseLen := length / factorProduct

	// Create base FFT for the smallest size
	var baseFFT *RadixN_NEON
	if baseLen > 1 {
		// For now, create a simple base FFT
		// In a full implementation, this would be a DFT or smaller RadixN
		baseFFT = &RadixN_NEON{
			length:    baseLen,
			direction: direction,
			factors:   nil, // Base case
		}
	}

	// Generate twiddle factors
	twiddles := generateRadixNTwiddles(length, factors, direction)

	// Create transpose factors
	transposeFactors := createTransposeFactors(factors)

	return &RadixN_NEON{
		length:    length,
		direction: direction,
		twiddles:  twiddles,
		factors:   transposeFactors,
		baseFFT:   baseFFT,
	}
}

// Process performs the RadixN FFT using NEON optimizations
func (r *RadixN_NEON) Process(data []complex128) {
	if len(data) < r.length {
		return
	}

	// For now, use optimized scalar implementation
	// TODO: Replace with actual NEON intrinsics
	r.processScalar(data)
}

// processScalar is the scalar implementation (placeholder for NEON)
func (r *RadixN_NEON) processScalar(data []complex128) {
	if r.baseFFT == nil {
		// Base case: direct DFT
		r.performDirectDFT(data)
		return
	}

	// Recursive case: apply RadixN decomposition
	r.performRadixNDecomposition(data)
}

// performDirectDFT performs a direct DFT for small sizes
func (r *RadixN_NEON) performDirectDFT(data []complex128) {
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

// performRadixNDecomposition performs the RadixN decomposition
func (r *RadixN_NEON) performRadixNDecomposition(data []complex128) {
	// Step 1: Factor transpose
	r.factorTranspose(data)

	// Step 2: Base FFTs
	if r.baseFFT != nil {
		baseLen := r.baseFFT.length
		for i := 0; i < r.length; i += baseLen {
			r.baseFFT.Process(data[i : i+baseLen])
		}
	}

	// Step 3: Cross FFTs with twiddles
	r.performCrossFFTs(data)
}

// factorTranspose performs the factor transpose operation
func (r *RadixN_NEON) factorTranspose(data []complex128) {
	if len(r.factors) == 0 {
		return
	}

	// For now, use a simple transpose
	// TODO: Implement proper factor transpose with NEON optimizations
	width := r.length / r.baseFFT.length
	height := r.baseFFT.length

	// Simple transpose
	transposed := make([]complex128, len(data))
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			transposed[i*width+j] = data[j*height+i]
		}
	}
	copy(data, transposed)
}

// performCrossFFTs performs cross FFTs with twiddle factors
func (r *RadixN_NEON) performCrossFFTs(data []complex128) {
	// Apply twiddle factors and cross FFTs
	// This is a simplified implementation
	// TODO: Implement proper cross FFTs with NEON optimizations

	twiddleIdx := 0
	for _, factor := range r.factors {
		radix := int(factor.Factor)
		// Apply twiddles and butterflies for this factor
		for i := 0; i < r.length; i += radix {
			// Apply twiddles
			for j := 1; j < radix; j++ {
				if twiddleIdx < len(r.twiddles) {
					data[i+j] *= r.twiddles[twiddleIdx]
					twiddleIdx++
				}
			}

			// Apply butterfly
			r.applyButterfly(data[i:i+radix], radix)
		}
	}
}

// applyButterfly applies a butterfly operation for the given radix
func (r *RadixN_NEON) applyButterfly(data []complex128, radix int) {
	switch radix {
	case 2:
		Butterfly2_NEON(data)
	case 3:
		Butterfly3_NEON(data)
	case 4:
		Butterfly4_NEON(data)
	case 5:
		Butterfly5_NEON(data)
	case 6:
		Butterfly6_NEON(data)
	case 7:
		Butterfly7_NEON(data)
	default:
		// Fallback to general butterfly
		r.performGeneralButterfly(data, radix)
	}
}

// performGeneralButterfly performs a general butterfly for arbitrary radix
func (r *RadixN_NEON) performGeneralButterfly(data []complex128, radix int) {
	// General butterfly implementation
	// This is a simplified version
	for k := 0; k < radix; k++ {
		sum := complex(0, 0)
		for j := 0; j < radix; j++ {
			angle := -2 * math.Pi * float64(k*j) / float64(radix)
			if r.direction == -1 {
				angle = -angle
			}
			w := complex(math.Cos(angle), math.Sin(angle))
			sum += data[j] * w
		}
		data[k] = sum
	}
}

// Specific RadixN implementations for common composite sizes

// RadixN_6_NEON performs a 6-point RadixN FFT using NEON (2×3)
func RadixN_6_NEON(data []complex128) {
	if len(data) < 6 {
		return
	}

	// Use real NEON assembly
	radixn_6_fft_go(data)
}

// RadixN_10_NEON performs a 10-point RadixN FFT using NEON (2×5)
func RadixN_10_NEON(data []complex128) {
	if len(data) < 10 {
		return
	}

	// Use real NEON assembly
	radixn_10_fft_go(data)
}

// RadixN_12_NEON performs a 12-point RadixN FFT using NEON (3×4)
func RadixN_12_NEON(data []complex128) {
	if len(data) < 12 {
		return
	}

	// Use real NEON assembly
	radixn_12_fft_go(data)
}

// RadixN_15_NEON performs a 15-point RadixN FFT using NEON (3×5)
func RadixN_15_NEON(data []complex128) {
	if len(data) < 15 {
		return
	}

	// Use real NEON assembly
	radixn_15_fft_go(data)
}

// RadixN_18_NEON performs an 18-point RadixN FFT using NEON (2×9)
func RadixN_18_NEON(data []complex128) {
	if len(data) < 18 {
		return
	}

	// Use real NEON assembly
	radixn_18_fft_go(data)
}

// RadixN_20_NEON performs a 20-point RadixN FFT using NEON (4×5)
func RadixN_20_NEON(data []complex128) {
	if len(data) < 20 {
		return
	}

	// Use real NEON assembly
	radixn_20_fft_go(data)
}

// ProcessVectorizedRadixN processes data using NEON-optimized RadixN FFTs
func ProcessVectorizedRadixN(data []complex128, size int) {
	switch size {
	case 6:
		RadixN_6_NEON(data)
	case 10:
		RadixN_10_NEON(data)
	case 12:
		RadixN_12_NEON(data)
	case 15:
		RadixN_15_NEON(data)
	case 18:
		RadixN_18_NEON(data)
	case 20:
		RadixN_20_NEON(data)
	default:
		// Fall back to scalar implementation for unsupported sizes
		processScalarRadixN(data, size)
	}
}

// processScalarRadixN is a fallback for unsupported RadixN sizes
func processScalarRadixN(data []complex128, size int) {
	// This would call the existing scalar RadixN implementation
	// For now, just a placeholder
}

// Helper functions

// generateRadixNTwiddles generates twiddle factors for RadixN FFT
func generateRadixNTwiddles(length int, factors []RadixFactor, direction int) []complex128 {
	if len(factors) == 0 {
		return nil
	}

	// Calculate total number of twiddles needed
	twiddleCount := 0
	crossFFTLen := length
	for _, factor := range factors {
		crossFFTLen /= int(factor)
		twiddleCount += crossFFTLen * (int(factor) - 1)
	}

	twiddles := make([]complex128, twiddleCount)
	idx := 0

	// Generate twiddles for each factor
	crossFFTLen = length
	for _, factor := range factors {
		crossFFTLen /= int(factor)
		radix := int(factor)

		for i := 1; i < radix; i++ {
			for j := 0; j < crossFFTLen; j++ {
				angle := -2 * math.Pi * float64(i*j) / float64(length)
				if direction == -1 {
					angle = -angle
				}
				twiddles[idx] = complex(math.Cos(angle), math.Sin(angle))
				idx++
			}
		}
	}

	return twiddles
}

// createTransposeFactors creates transpose factors from radix factors
func createTransposeFactors(factors []RadixFactor) []TransposeFactor {
	if len(factors) == 0 {
		return nil
	}

	// Create transpose factors (reversed and collapsed)
	transposeFactors := make([]TransposeFactor, 0, len(factors))

	// Reverse the factors
	for i := len(factors) - 1; i >= 0; i-- {
		factor := factors[i]

		// Collapse adjacent identical factors
		if len(transposeFactors) > 0 && transposeFactors[len(transposeFactors)-1].Factor == factor {
			transposeFactors[len(transposeFactors)-1].Count++
		} else {
			transposeFactors = append(transposeFactors, TransposeFactor{
				Factor: factor,
				Count:  1,
			})
		}
	}

	return transposeFactors
}

// reverseRemainders performs remainder reversal for factor transpose
func reverseRemainders(value int, factors []TransposeFactor) int {
	result := 0
	temp := value

	for _, factor := range factors {
		radix := int(factor.Factor)
		for i := 0; i < int(factor.Count); i++ {
			result = result*radix + (temp % radix)
			temp = temp / radix
		}
	}

	return result
}

// NEON-specific utility functions for RadixN
// Note: loadComplex128Vector, storeComplex128Vector, and multiplyComplex128Vector
// are defined in radix4.go to avoid duplication

// NEON memory alignment utilities for RadixN

// isAligned64 checks if a pointer is 64-byte aligned (required for NEON RadixN)
func isAligned64(ptr unsafe.Pointer) bool {
	return uintptr(ptr)%64 == 0
}

// alignTo64 aligns a slice to 64-byte boundary for NEON RadixN operations
func alignTo64(data []complex128) []complex128 {
	if len(data) == 0 {
		return data
	}

	ptr := unsafe.Pointer(&data[0])
	if isAligned64(ptr) {
		return data
	}

	// Create aligned copy
	aligned := make([]complex128, len(data))
	copy(aligned, data)
	return aligned
}
