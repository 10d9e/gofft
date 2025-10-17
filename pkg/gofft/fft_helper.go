package gofft

// ProcessFunc is a function type that processes a single FFT chunk
type ProcessFunc func(chunk, scratch []complex128)

// ProcessFuncOutOfPlace is a function type for out-of-place FFT processing
type ProcessFuncOutOfPlace func(input, output, scratch []complex128)

// ProcessFuncImmutable is a function type for immutable FFT processing
type ProcessFuncImmutable func(input []complex128, output, scratch []complex128)

// fftHelperInplace processes multiple FFT chunks in-place
func fftHelperInplace(buffer, scratch []complex128, expectedLen, expectedScratch int, process ProcessFunc) {
	validateInplace(len(buffer), expectedLen, len(scratch), expectedScratch)

	// Process each chunk
	for i := 0; i < len(buffer); i += expectedLen {
		chunk := buffer[i : i+expectedLen]
		process(chunk, scratch)
	}
}

// fftHelperOutOfPlace processes multiple FFT chunks out-of-place
func fftHelperOutOfPlace(input, output, scratch []complex128, expectedLen, expectedScratch int, process ProcessFuncOutOfPlace) {
	validateOutOfPlace(len(input), len(output), expectedLen, len(scratch), expectedScratch)

	// Process each chunk
	for i := 0; i < len(input); i += expectedLen {
		inChunk := input[i : i+expectedLen]
		outChunk := output[i : i+expectedLen]
		process(inChunk, outChunk, scratch)
	}
}

// fftHelperImmutable processes multiple FFT chunks without modifying input
func fftHelperImmutable(input []complex128, output, scratch []complex128, expectedLen, expectedScratch int, process ProcessFuncImmutable) {
	validateOutOfPlace(len(input), len(output), expectedLen, len(scratch), expectedScratch)

	// Process each chunk
	for i := 0; i < len(input); i += expectedLen {
		inChunk := input[i : i+expectedLen]
		outChunk := output[i : i+expectedLen]
		process(inChunk, outChunk, scratch)
	}
}

// Complex64 versions

// ProcessFunc32 is a function type that processes a single FFT chunk (complex64)
type ProcessFunc32 func(chunk, scratch []complex64)

// ProcessFuncOutOfPlace32 is a function type for out-of-place FFT processing (complex64)
type ProcessFuncOutOfPlace32 func(input, output, scratch []complex64)

// ProcessFuncImmutable32 is a function type for immutable FFT processing (complex64)
type ProcessFuncImmutable32 func(input []complex64, output, scratch []complex64)

// fftHelperInplace32 processes multiple FFT chunks in-place (complex64)
func fftHelperInplace32(buffer, scratch []complex64, expectedLen, expectedScratch int, process ProcessFunc32) {
	validateInplace(len(buffer), expectedLen, len(scratch), expectedScratch)

	for i := 0; i < len(buffer); i += expectedLen {
		chunk := buffer[i : i+expectedLen]
		process(chunk, scratch)
	}
}

// fftHelperOutOfPlace32 processes multiple FFT chunks out-of-place (complex64)
func fftHelperOutOfPlace32(input, output, scratch []complex64, expectedLen, expectedScratch int, process ProcessFuncOutOfPlace32) {
	validateOutOfPlace(len(input), len(output), expectedLen, len(scratch), expectedScratch)

	for i := 0; i < len(input); i += expectedLen {
		inChunk := input[i : i+expectedLen]
		outChunk := output[i : i+expectedLen]
		process(inChunk, outChunk, scratch)
	}
}

// fftHelperImmutable32 processes multiple FFT chunks without modifying input (complex64)
func fftHelperImmutable32(input []complex64, output, scratch []complex64, expectedLen, expectedScratch int, process ProcessFuncImmutable32) {
	validateOutOfPlace(len(input), len(output), expectedLen, len(scratch), expectedScratch)

	for i := 0; i < len(input); i += expectedLen {
		inChunk := input[i : i+expectedLen]
		outChunk := output[i : i+expectedLen]
		process(inChunk, outChunk, scratch)
	}
}
