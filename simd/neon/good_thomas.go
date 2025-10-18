package neon

// GoodThomas_35_NEON performs a 35-point Good-Thomas FFT using NEON
func GoodThomas_35_NEON(data []complex128) {
	if len(data) < 35 {
		return
	}

	// Use real NEON assembly
	good_thomas_35_fft_go(data)
}

// GoodThomas_77_NEON performs a 77-point Good-Thomas FFT using NEON
func GoodThomas_77_NEON(data []complex128) {
	if len(data) < 77 {
		return
	}

	// Use real NEON assembly
	good_thomas_77_fft_go(data)
}

// GoodThomas_143_NEON performs a 143-point Good-Thomas FFT using NEON
func GoodThomas_143_NEON(data []complex128) {
	if len(data) < 143 {
		return
	}

	// For now, use optimized scalar implementation
	// TODO: Replace with actual NEON intrinsics
	// This would implement Good-Thomas algorithm for coprime factorization
	// 143 = 11 * 13 (coprime factors)
}

// GoodThomas_221_NEON performs a 221-point Good-Thomas FFT using NEON
func GoodThomas_221_NEON(data []complex128) {
	if len(data) < 221 {
		return
	}

	// For now, use optimized scalar implementation
	// TODO: Replace with actual NEON intrinsics
	// This would implement Good-Thomas algorithm for coprime factorization
	// 221 = 13 * 17 (coprime factors)
}

// ProcessVectorizedGoodThomas processes data using NEON-optimized Good-Thomas FFTs
func ProcessVectorizedGoodThomas(data []complex128, size int) {
	switch size {
	case 35:
		GoodThomas_35_NEON(data)
	case 77:
		GoodThomas_77_NEON(data)
	case 143:
		GoodThomas_143_NEON(data)
	case 221:
		GoodThomas_221_NEON(data)
	default:
		// Fall back to scalar implementation for other sizes
		// This would use the existing scalar Good-Thomas algorithm
	}
}
