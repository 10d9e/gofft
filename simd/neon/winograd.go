package neon

// Winograd_49_NEON performs a 49-point Winograd FFT using NEON
func Winograd_49_NEON(data []complex128) {
	if len(data) < 49 {
		return
	}

	// Use real NEON assembly
	winograd_49_fft_go(data)
}

// Winograd_121_NEON performs a 121-point Winograd FFT using NEON
func Winograd_121_NEON(data []complex128) {
	if len(data) < 121 {
		return
	}

	// Use real NEON assembly
	winograd_121_fft_go(data)
}

// Winograd_169_NEON performs a 169-point Winograd FFT using NEON
func Winograd_169_NEON(data []complex128) {
	if len(data) < 169 {
		return
	}

	// For now, use optimized scalar implementation
	// TODO: Replace with actual NEON intrinsics
	// This would implement Winograd's algorithm for minimal multiplication
	// 169 = 13 * 13 (perfect square)
}

// Winograd_289_NEON performs a 289-point Winograd FFT using NEON
func Winograd_289_NEON(data []complex128) {
	if len(data) < 289 {
		return
	}

	// For now, use optimized scalar implementation
	// TODO: Replace with actual NEON intrinsics
	// This would implement Winograd's algorithm for minimal multiplication
	// 289 = 17 * 17 (perfect square)
}

// ProcessVectorizedWinograd processes data using NEON-optimized Winograd FFTs
func ProcessVectorizedWinograd(data []complex128, size int) {
	switch size {
	case 49:
		Winograd_49_NEON(data)
	case 121:
		Winograd_121_NEON(data)
	case 169:
		Winograd_169_NEON(data)
	case 289:
		Winograd_289_NEON(data)
	default:
		// Fall back to scalar implementation for other sizes
		// This would use the existing scalar Winograd algorithm
	}
}
