// Package gofft provides a high-performance FFT library for Go, inspired by RustFFT.
//
// gofft supports SIMD acceleration on x86_64 (SSE4.1, AVX) and ARM64 (NEON) architectures
// with automatic runtime CPU feature detection.
//
// Usage:
//
//	import "github.com/example/gofft/pkg/gofft"
//
//	// Create a planner that automatically selects the best implementation
//	planner := gofft.NewPlanner()
//	fft := planner.PlanForward(1234)
//
//	// Create a buffer of complex numbers
//	buffer := make([]complex128, 1234)
//	// ... fill buffer with data ...
//
//	// Compute FFT in-place
//	fft.Process(buffer)
//
// The library does not normalize outputs. Callers must manually normalize results
// by scaling each element by 1/sqrt(len). When doing a forward FFT followed by
// an inverse FFT, you can normalize once by scaling by 1/len.
package gofft

// Fft is the main interface for computing FFTs.
//
// All Fft implementations are thread-safe and can be used concurrently.
type Fft interface {
	// Process computes an FFT in-place on the provided buffer.
	// The buffer length must be a multiple of Len().
	// Allocates scratch space internally.
	Process(buffer []complex128)

	// ProcessWithScratch computes an FFT in-place using the provided scratch buffer.
	// The scratch buffer must have length >= InplaceScratchLen().
	// The buffer length must be a multiple of Len().
	ProcessWithScratch(buffer, scratch []complex128)

	// ProcessOutOfPlace computes an FFT from input to output using the provided scratch buffer.
	// Input and output must have the same length, which must be a multiple of Len().
	// The scratch buffer must have length >= OutOfPlaceScratchLen().
	// The contents of input are destroyed.
	ProcessOutOfPlace(input, output, scratch []complex128)

	// ProcessImmutable computes an FFT from input to output without modifying input.
	// Input and output must have the same length, which must be a multiple of Len().
	// The scratch buffer must have length >= ImmutableScratchLen().
	ProcessImmutable(input []complex128, output, scratch []complex128)

	// Len returns the FFT size that this instance processes
	Len() int

	// Direction returns whether this instance computes forward or inverse FFTs
	Direction() Direction

	// InplaceScratchLen returns the required scratch buffer size for ProcessWithScratch
	InplaceScratchLen() int

	// OutOfPlaceScratchLen returns the required scratch buffer size for ProcessOutOfPlace
	OutOfPlaceScratchLen() int

	// ImmutableScratchLen returns the required scratch buffer size for ProcessImmutable
	ImmutableScratchLen() int
}

// Fft32 is the interface for computing FFTs on complex64 values
type Fft32 interface {
	Process(buffer []complex64)
	ProcessWithScratch(buffer, scratch []complex64)
	ProcessOutOfPlace(input, output, scratch []complex64)
	ProcessImmutable(input []complex64, output, scratch []complex64)
	Len() int
	Direction() Direction
	InplaceScratchLen() int
	OutOfPlaceScratchLen() int
	ImmutableScratchLen() int
}
