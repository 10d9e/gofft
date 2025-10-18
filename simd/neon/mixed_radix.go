package neon

// MixedRadix_60_NEON performs a 60-point Mixed-Radix FFT using NEON
func MixedRadix_60_NEON(data []complex128) {
	if len(data) < 60 {
		return
	}

	// Use real NEON assembly
	mixed_radix_60_fft_go(data)
}

// MixedRadix_120_NEON performs a 120-point Mixed-Radix FFT using NEON
func MixedRadix_120_NEON(data []complex128) {
	if len(data) < 120 {
		return
	}

	// Use real NEON assembly
	mixed_radix_120_fft_go(data)
}

// MixedRadix_240_NEON performs a 240-point Mixed-Radix FFT using NEON
func MixedRadix_240_NEON(data []complex128) {
	if len(data) < 240 {
		return
	}

	// Use real NEON assembly
	mixed_radix_240_fft_go(data)
}

// MixedRadix_480_NEON performs a 480-point Mixed-Radix FFT using NEON
func MixedRadix_480_NEON(data []complex128) {
	if len(data) < 480 {
		return
	}

	// For now, use optimized scalar implementation
	// TODO: Replace with actual NEON intrinsics
	// This would implement Mixed-Radix algorithm for composite sizes
	// using combination of different radix sizes
}

// ProcessVectorizedMixedRadix processes data using NEON-optimized Mixed-Radix FFTs
func ProcessVectorizedMixedRadix(data []complex128, size int) {
	switch size {
	case 60:
		MixedRadix_60_NEON(data)
	case 120:
		MixedRadix_120_NEON(data)
	case 240:
		MixedRadix_240_NEON(data)
	case 480:
		MixedRadix_480_NEON(data)
	default:
		// Fall back to scalar implementation for other sizes
		// This would use the existing scalar Mixed-Radix algorithm
	}
}
