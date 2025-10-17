package gofft

import (
	"fmt"
	"math"
)

// FftNum is a constraint for numeric types that can be used in FFT computations
type FftNum interface {
	~float32 | ~float64
}

// validateBuffers validates buffer sizes for in-place FFT operations
func validateInplace(bufferLen, expectedLen, scratchLen, expectedScratch int) {
	if bufferLen < expectedLen {
		panic(fmt.Sprintf("Provided FFT buffer was too small. Expected len = %d, got len = %d", expectedLen, bufferLen))
	}
	if bufferLen%expectedLen != 0 {
		panic(fmt.Sprintf("Input FFT buffer must be a multiple of FFT length. Expected multiple of %d, got len = %d", expectedLen, bufferLen))
	}
	if scratchLen < expectedScratch {
		panic(fmt.Sprintf("Not enough scratch space was provided. Expected scratch len >= %d, got scratch len = %d", expectedScratch, scratchLen))
	}
}

// validateOutOfPlace validates buffer sizes for out-of-place FFT operations
func validateOutOfPlace(inputLen, outputLen, expectedLen, scratchLen, expectedScratch int) {
	if inputLen != outputLen {
		panic(fmt.Sprintf("Provided FFT input buffer and output buffer must have the same length. Got input.len() = %d, output.len() = %d", inputLen, outputLen))
	}
	if inputLen < expectedLen {
		panic(fmt.Sprintf("Provided FFT buffer was too small. Expected len = %d, got len = %d", expectedLen, inputLen))
	}
	if inputLen%expectedLen != 0 {
		panic(fmt.Sprintf("Input FFT buffer must be a multiple of FFT length. Expected multiple of %d, got len = %d", expectedLen, inputLen))
	}
	if scratchLen < expectedScratch {
		panic(fmt.Sprintf("Not enough scratch space was provided. Expected scratch len >= %d, got scratch len = %d", expectedScratch, scratchLen))
	}
}

// Complex utility functions for complex128

// ComplexMul multiplies two complex numbers
func ComplexMul(a, b complex128) complex128 {
	return a * b
}

// ComplexMulAdd computes a * b + c
func ComplexMulAdd(a, b, c complex128) complex128 {
	return a*b + c
}

// ComplexMulSub computes a * b - c
func ComplexMulSub(a, b, c complex128) complex128 {
	return a*b - c
}

// ComplexConj returns the complex conjugate
func ComplexConj(a complex128) complex128 {
	return complex(real(a), -imag(a))
}

// ComplexScale multiplies a complex number by a real scalar
func ComplexScale(a complex128, s float64) complex128 {
	return complex(real(a)*s, imag(a)*s)
}

// ComplexFromPolar creates a complex number from polar coordinates
func ComplexFromPolar(r, theta float64) complex128 {
	return complex(r*math.Cos(theta), r*math.Sin(theta))
}

// Complex32 utility functions for complex64

// ComplexMul32 multiplies two complex64 numbers
func ComplexMul32(a, b complex64) complex64 {
	return a * b
}

// ComplexMulAdd32 computes a * b + c for complex64
func ComplexMulAdd32(a, b, c complex64) complex64 {
	return a*b + c
}

// ComplexMulSub32 computes a * b - c for complex64
func ComplexMulSub32(a, b, c complex64) complex64 {
	return a*b - c
}

// ComplexConj32 returns the complex conjugate of a complex64
func ComplexConj32(a complex64) complex64 {
	return complex(real(a), -imag(a))
}

// ComplexScale32 multiplies a complex64 number by a real scalar
func ComplexScale32(a complex64, s float32) complex64 {
	return complex(real(a)*s, imag(a)*s)
}

// ComplexFromPolar32 creates a complex64 number from polar coordinates
func ComplexFromPolar32(r, theta float32) complex64 {
	return complex(r*float32(math.Cos(float64(theta))), r*float32(math.Sin(float64(theta))))
}

// RadixFactor represents the radix of a single FFT stage
type RadixFactor uint8

const (
	Factor2 RadixFactor = 2
	Factor3 RadixFactor = 3
	Factor4 RadixFactor = 4
	Factor5 RadixFactor = 5
	Factor6 RadixFactor = 6
	Factor7 RadixFactor = 7
)

// Radix returns the numeric value of the radix factor
func (r RadixFactor) Radix() int {
	return int(r)
}
