//go:build arm64

package neon

import (
	"unsafe"
)

// NEON intrinsics wrapper functions
// These call the actual NEON assembly implementations

// solo_fft2_f64_go performs a 2-point FFT using NEON intrinsics
// Based on RustFFT's solo_fft2_f64 implementation
func solo_fft2_f64_go(left, right complex128) [2]complex128 {
	var result [2]complex128

	// Call assembly function
	solo_fft2_f64_asm(
		unsafe.Pointer(&left),
		unsafe.Pointer(&right),
		unsafe.Pointer(&result),
	)

	return result
}

// load_complex_f64_go loads a complex128 into a NEON vector
func load_complex_f64_go(ptr *complex128) complex128 {
	var result complex128

	// Call assembly function
	load_complex_f64_asm(
		unsafe.Pointer(ptr),
		unsafe.Pointer(&result),
	)

	return result
}

// store_complex_f64_go stores a complex128 from a NEON vector
func store_complex_f64_go(ptr *complex128, value complex128) {
	// Call assembly function
	store_complex_f64_asm(
		unsafe.Pointer(ptr),
		unsafe.Pointer(&value),
	)
}

// add_complex_f64_go adds two complex numbers using NEON
func add_complex_f64_go(a, b complex128) complex128 {
	// Call assembly function
	return add_complex_f64_asm(a, b)
}

// sub_complex_f64_go subtracts two complex numbers using NEON
func sub_complex_f64_go(a, b complex128) complex128 {
	// Call assembly function
	return sub_complex_f64_asm(a, b)
}

// mul_complex_f64_go multiplies two complex numbers using NEON
func mul_complex_f64_go(a, b complex128) complex128 {
	// Call assembly function
	return mul_complex_f64_asm(a, b)
}

// transpose_2x2_f64_go transposes a 2x2 matrix of complex numbers
func transpose_2x2_f64_go(a, b complex128) (complex128, complex128) {
	// Call assembly function
	return transpose_2x2_f64_asm(a, b)
}

// Assembly function declarations
// These are implemented in neon_asm.s

//go:noescape
func solo_fft2_f64_asm(left, right, result unsafe.Pointer)

//go:noescape
func load_complex_f64_asm(ptr, result unsafe.Pointer)

//go:noescape
func store_complex_f64_asm(ptr, value unsafe.Pointer)

//go:noescape
func add_complex_f64_asm(a, b complex128) complex128

//go:noescape
func sub_complex_f64_asm(a, b complex128) complex128

//go:noescape
func mul_complex_f64_asm(a, b complex128) complex128

//go:noescape
func transpose_2x2_f64_asm(a, b complex128) (complex128, complex128)

// Real NEON Butterfly implementations using actual intrinsics

// Butterfly2_NEON_Real performs a 2-point butterfly using real NEON intrinsics
func Butterfly2_NEON_Real(data []complex128) {
	if len(data) < 2 {
		return
	}

	// Use real NEON intrinsics
	result := solo_fft2_f64_go(data[0], data[1])
	data[0] = result[0]
	data[1] = result[1]
}

// Butterfly4_NEON_Real performs a 4-point butterfly using real NEON intrinsics
func Butterfly4_NEON_Real(data []complex128) {
	if len(data) < 4 {
		return
	}

	// 4-point butterfly using 2x2 mixed radix
	// Step 1: 2-point butterflies
	Butterfly2_NEON_Real(data[0:2])
	Butterfly2_NEON_Real(data[2:4])

	// Step 2: Apply twiddle factors (90-degree rotation for size-2)
	// For size-4, the twiddle factor is i (0 + 1i)
	// This means: (a + bi) * i = -b + ai
	// So we need to swap real and imaginary parts and negate the new real part

	// Apply twiddle to data[2] and data[3]
	// data[2] *= i, data[3] *= i
	for i := 2; i < 4; i++ {
		re := real(data[i])
		im := imag(data[i])
		data[i] = complex(-im, re) // (a + bi) * i = -b + ai
	}

	// Step 3: Transpose and apply 2-point butterflies
	// This is a simplified version - full implementation would use NEON transpose
	temp0 := data[0]
	temp1 := data[1]
	temp2 := data[2]
	temp3 := data[3]

	// Transpose: [0,1,2,3] -> [0,2,1,3]
	data[0] = temp0
	data[1] = temp2
	data[2] = temp1
	data[3] = temp3

	// Step 4: Final 2-point butterflies
	Butterfly2_NEON_Real(data[0:2])
	Butterfly2_NEON_Real(data[2:4])

	// Step 5: Transpose back: [0,2,1,3] -> [0,1,2,3]
	temp0 = data[0]
	temp1 = data[1]
	temp2 = data[2]
	temp3 = data[3]

	data[0] = temp0
	data[1] = temp2
	data[2] = temp1
	data[3] = temp3
}

// Butterfly8_NEON_Real performs an 8-point butterfly using real NEON intrinsics
func Butterfly8_NEON_Real(data []complex128) {
	if len(data) < 8 {
		return
	}

	// 8-point butterfly using 2x4 mixed radix
	// Step 1: 2-point butterflies
	for i := 0; i < 8; i += 2 {
		Butterfly2_NEON_Real(data[i : i+2])
	}

	// Step 2: Apply twiddle factors
	// For size-8, we have twiddle factors: 1, w, w^2, w^3 where w = e^(-2πi/8)
	// w = (1-i)/√2, w^2 = -i, w^3 = -(1+i)/√2

	// Apply twiddles to data[2], data[4], data[6]
	// data[2] *= w, data[4] *= w^2, data[6] *= w^3
	const sqrt2 = 1.4142135623730951
	const w_re = 1.0 / sqrt2
	const w_im = -1.0 / sqrt2
	const w2_re = 0.0
	const w2_im = -1.0
	const w3_re = -1.0 / sqrt2
	const w3_im = -1.0 / sqrt2

	// Apply w to data[2]
	re := real(data[2])
	im := imag(data[2])
	data[2] = complex(re*w_re-im*w_im, re*w_im+im*w_re)

	// Apply w^2 to data[4] (multiply by -i)
	re = real(data[4])
	im = imag(data[4])
	data[4] = complex(im, -re)

	// Apply w^3 to data[6]
	re = real(data[6])
	im = imag(data[6])
	data[6] = complex(re*w3_re-im*w3_im, re*w3_im+im*w3_re)

	// Step 3: Transpose and apply 4-point butterflies
	// This is a simplified version - full implementation would use NEON transpose
	temp := make([]complex128, 8)
	copy(temp, data)

	// Transpose 2x4 matrix
	data[0] = temp[0]
	data[1] = temp[2]
	data[2] = temp[4]
	data[3] = temp[6]
	data[4] = temp[1]
	data[5] = temp[3]
	data[6] = temp[5]
	data[7] = temp[7]

	// Step 4: Apply 4-point butterflies
	Butterfly4_NEON_Real(data[0:4])
	Butterfly4_NEON_Real(data[4:8])

	// Step 5: Transpose back
	copy(temp, data)
	data[0] = temp[0]
	data[1] = temp[4]
	data[2] = temp[1]
	data[3] = temp[5]
	data[4] = temp[2]
	data[5] = temp[6]
	data[6] = temp[3]
	data[7] = temp[7]
}

// ProcessVectorizedButterfly_Real processes data using real NEON-optimized butterflies
func ProcessVectorizedButterfly_Real(data []complex128, size int) {
	switch size {
	case 2:
		Butterfly2_NEON_Real(data)
	case 4:
		Butterfly4_NEON_Real(data)
	case 8:
		Butterfly8_NEON_Real(data)
	default:
		// Fall back to scalar implementation for unsupported sizes
		processScalarButterfly_Real(data, size)
	}
}

// processScalarButterfly_Real is a fallback for unsupported butterfly sizes
func processScalarButterfly_Real(data []complex128, size int) {
	// This would call the existing scalar butterfly implementation
	// For now, just a placeholder
}

// butterfly16_fft_go performs a 16-point FFT using real NEON assembly
func butterfly16_fft_go(data []complex128) {
	if len(data) < 16 {
		return
	}

	// Call assembly function
	butterfly16_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly32_fft_go performs a 32-point FFT using real NEON assembly
func butterfly32_fft_go(data []complex128) {
	if len(data) < 32 {
		return
	}

	// Call assembly function
	butterfly32_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly3_fft_go performs a 3-point FFT using real NEON assembly
func butterfly3_fft_go(data []complex128) {
	if len(data) < 3 {
		return
	}

	// Call assembly function
	butterfly3_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly5_fft_go performs a 5-point FFT using real NEON assembly
func butterfly5_fft_go(data []complex128) {
	if len(data) < 5 {
		return
	}

	// Call assembly function
	butterfly5_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly7_fft_go performs a 7-point FFT using real NEON assembly
func butterfly7_fft_go(data []complex128) {
	if len(data) < 7 {
		return
	}

	// Call assembly function
	butterfly7_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly6_fft_go performs a 6-point FFT using real NEON assembly
func butterfly6_fft_go(data []complex128) {
	if len(data) < 6 {
		return
	}

	// Call assembly function
	butterfly6_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly12_fft_go performs a 12-point FFT using real NEON assembly
func butterfly12_fft_go(data []complex128) {
	if len(data) < 12 {
		return
	}

	// Call assembly function
	butterfly12_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly24_fft_go performs a 24-point FFT using real NEON assembly
func butterfly24_fft_go(data []complex128) {
	if len(data) < 24 {
		return
	}

	// Call assembly function
	butterfly24_fft_asm(unsafe.Pointer(&data[0]))
}

// radix4_64_fft_go performs a 64-point Radix-4 FFT using real NEON assembly
func radix4_64_fft_go(data []complex128) {
	if len(data) < 64 {
		return
	}

	// Call assembly function
	radix4_64_fft_asm(unsafe.Pointer(&data[0]))
}

// radix4_128_fft_go performs a 128-point Radix-4 FFT using real NEON assembly
func radix4_128_fft_go(data []complex128) {
	if len(data) < 128 {
		return
	}

	// Call assembly function
	radix4_128_fft_asm(unsafe.Pointer(&data[0]))
}

// radix4_256_fft_go performs a 256-point Radix-4 FFT using real NEON assembly
func radix4_256_fft_go(data []complex128) {
	if len(data) < 256 {
		return
	}

	// Call assembly function
	radix4_256_fft_asm(unsafe.Pointer(&data[0]))
}

// radix4_512_fft_go performs a 512-point Radix-4 FFT using real NEON assembly
func radix4_512_fft_go(data []complex128) {
	if len(data) < 512 {
		return
	}

	// Call assembly function
	radix4_512_fft_asm(unsafe.Pointer(&data[0]))
}

// radix4_1024_fft_go performs a 1024-point Radix-4 FFT using real NEON assembly
func radix4_1024_fft_go(data []complex128) {
	if len(data) < 1024 {
		return
	}

	// Call assembly function
	radix4_1024_fft_asm(unsafe.Pointer(&data[0]))
}

// radixn_6_fft_go performs a 6-point RadixN FFT using real NEON assembly
func radixn_6_fft_go(data []complex128) {
	if len(data) < 6 {
		return
	}

	// Call assembly function
	radixn_6_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_37_fft_go performs a 37-point Rader's FFT using real NEON assembly
func raders_37_fft_go(data []complex128) {
	if len(data) < 37 {
		return
	}

	// Call assembly function
	raders_37_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_41_fft_go performs a 41-point Rader's FFT using real NEON assembly
func raders_41_fft_go(data []complex128) {
	if len(data) < 41 {
		return
	}

	// Call assembly function
	raders_41_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_43_fft_go performs a 43-point Rader's FFT using real NEON assembly
func raders_43_fft_go(data []complex128) {
	if len(data) < 43 {
		return
	}

	// Call assembly function
	raders_43_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_47_fft_go performs a 47-point Rader's FFT using real NEON assembly
func raders_47_fft_go(data []complex128) {
	if len(data) < 47 {
		return
	}

	// Call assembly function
	raders_47_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_53_fft_go performs a 53-point Rader's FFT using real NEON assembly
func raders_53_fft_go(data []complex128) {
	if len(data) < 53 {
		return
	}

	// Call assembly function
	raders_53_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_59_fft_go performs a 59-point Rader's FFT using real NEON assembly
func raders_59_fft_go(data []complex128) {
	if len(data) < 59 {
		return
	}

	// Call assembly function
	raders_59_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_61_fft_go performs a 61-point Rader's FFT using real NEON assembly
func raders_61_fft_go(data []complex128) {
	if len(data) < 61 {
		return
	}

	// Call assembly function
	raders_61_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_67_fft_go performs a 67-point Rader's FFT using real NEON assembly
func raders_67_fft_go(data []complex128) {
	if len(data) < 67 {
		return
	}

	// Call assembly function
	raders_67_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_71_fft_go performs a 71-point Rader's FFT using real NEON assembly
func raders_71_fft_go(data []complex128) {
	if len(data) < 71 {
		return
	}

	// Call assembly function
	raders_71_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_73_fft_go performs a 73-point Rader's FFT using real NEON assembly
func raders_73_fft_go(data []complex128) {
	if len(data) < 73 {
		return
	}

	// Call assembly function
	raders_73_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_79_fft_go performs a 79-point Rader's FFT using real NEON assembly
func raders_79_fft_go(data []complex128) {
	if len(data) < 79 {
		return
	}

	// Call assembly function
	raders_79_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_83_fft_go performs a 83-point Rader's FFT using real NEON assembly
func raders_83_fft_go(data []complex128) {
	if len(data) < 83 {
		return
	}

	// Call assembly function
	raders_83_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_89_fft_go performs a 89-point Rader's FFT using real NEON assembly
func raders_89_fft_go(data []complex128) {
	if len(data) < 89 {
		return
	}

	// Call assembly function
	raders_89_fft_asm(unsafe.Pointer(&data[0]))
}

// raders_97_fft_go performs a 97-point Rader's FFT using real NEON assembly
func raders_97_fft_go(data []complex128) {
	if len(data) < 97 {
		return
	}

	// Call assembly function
	raders_97_fft_asm(unsafe.Pointer(&data[0]))
}

// bluestein_128_fft_go performs a 128-point Bluestein's FFT using real NEON assembly
func bluestein_128_fft_go(data []complex128) {
	if len(data) < 128 {
		return
	}

	// Call assembly function
	bluestein_128_fft_asm(unsafe.Pointer(&data[0]))
}

// mixed_radix_60_fft_go performs a 60-point Mixed-Radix FFT using real NEON assembly
func mixed_radix_60_fft_go(data []complex128) {
	if len(data) < 60 {
		return
	}

	// Call assembly function
	mixed_radix_60_fft_asm(unsafe.Pointer(&data[0]))
}

// mixed_radix_120_fft_go performs a 120-point Mixed-Radix FFT using real NEON assembly
func mixed_radix_120_fft_go(data []complex128) {
	if len(data) < 120 {
		return
	}

	// Call assembly function
	mixed_radix_120_fft_asm(unsafe.Pointer(&data[0]))
}

// mixed_radix_240_fft_go performs a 240-point Mixed-Radix FFT using real NEON assembly
func mixed_radix_240_fft_go(data []complex128) {
	if len(data) < 240 {
		return
	}

	// Call assembly function
	mixed_radix_240_fft_asm(unsafe.Pointer(&data[0]))
}

// good_thomas_35_fft_go performs a 35-point Good-Thomas FFT using real NEON assembly
func good_thomas_35_fft_go(data []complex128) {
	if len(data) < 35 {
		return
	}

	// Call assembly function
	good_thomas_35_fft_asm(unsafe.Pointer(&data[0]))
}

// good_thomas_77_fft_go performs a 77-point Good-Thomas FFT using real NEON assembly
func good_thomas_77_fft_go(data []complex128) {
	if len(data) < 77 {
		return
	}

	// Call assembly function
	good_thomas_77_fft_asm(unsafe.Pointer(&data[0]))
}

// winograd_49_fft_go performs a 49-point Winograd FFT using real NEON assembly
func winograd_49_fft_go(data []complex128) {
	if len(data) < 49 {
		return
	}

	// Call assembly function
	winograd_49_fft_asm(unsafe.Pointer(&data[0]))
}

// winograd_121_fft_go performs a 121-point Winograd FFT using real NEON assembly
func winograd_121_fft_go(data []complex128) {
	if len(data) < 121 {
		return
	}

	// Call assembly function
	winograd_121_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly1_fft_go performs a 1-point Butterfly FFT using real NEON assembly
func butterfly1_fft_go(data []complex128) {
	if len(data) < 1 {
		return
	}

	// Call assembly function
	butterfly1_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly10_fft_go performs a 10-point Butterfly FFT using real NEON assembly
func butterfly10_fft_go(data []complex128) {
	if len(data) < 10 {
		return
	}

	// Call assembly function
	butterfly10_fft_asm(unsafe.Pointer(&data[0]))
}

// butterfly15_fft_go performs a 15-point Butterfly FFT using real NEON assembly
func butterfly15_fft_go(data []complex128) {
	if len(data) < 15 {
		return
	}

	// Call assembly function
	butterfly15_fft_asm(unsafe.Pointer(&data[0]))
}

// radixn_10_fft_go performs a 10-point RadixN FFT using real NEON assembly
func radixn_10_fft_go(data []complex128) {
	if len(data) < 10 {
		return
	}

	// Call assembly function
	radixn_10_fft_asm(unsafe.Pointer(&data[0]))
}

// radixn_12_fft_go performs a 12-point RadixN FFT using real NEON assembly
func radixn_12_fft_go(data []complex128) {
	if len(data) < 12 {
		return
	}

	// Call assembly function
	radixn_12_fft_asm(unsafe.Pointer(&data[0]))
}

// radixn_15_fft_go performs a 15-point RadixN FFT using real NEON assembly
func radixn_15_fft_go(data []complex128) {
	if len(data) < 15 {
		return
	}

	// Call assembly function
	radixn_15_fft_asm(unsafe.Pointer(&data[0]))
}

// radixn_18_fft_go performs an 18-point RadixN FFT using real NEON assembly
func radixn_18_fft_go(data []complex128) {
	if len(data) < 18 {
		return
	}

	// Call assembly function
	radixn_18_fft_asm(unsafe.Pointer(&data[0]))
}

// radixn_20_fft_go performs a 20-point RadixN FFT using real NEON assembly
func radixn_20_fft_go(data []complex128) {
	if len(data) < 20 {
		return
	}

	// Call assembly function
	radixn_20_fft_asm(unsafe.Pointer(&data[0]))
}

// Assembly function declarations for Butterfly1, Butterfly3, Butterfly5, Butterfly6, Butterfly7, Butterfly10, Butterfly12, Butterfly15, Butterfly16, Butterfly24, Butterfly32, Radix4_64, Radix4_128, Radix4_256, Radix4_512, Radix4_1024, RadixN_6, RadixN_10, RadixN_12, RadixN_15, RadixN_18, RadixN_20, Raders_37, Raders_41, Raders_43, Raders_47, Raders_53, Raders_59, Raders_61, Raders_67, Raders_71, Raders_73, Raders_79, Raders_83, Raders_89, Raders_97, Bluestein_128, MixedRadix_60, GoodThomas_35 and Winograd_49
//
//go:noescape
func butterfly1_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly3_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly5_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly6_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly7_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly10_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly12_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly15_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly16_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly24_fft_asm(data unsafe.Pointer)

//go:noescape
func butterfly32_fft_asm(data unsafe.Pointer)

//go:noescape
func radix4_64_fft_asm(data unsafe.Pointer)

//go:noescape
func radix4_128_fft_asm(data unsafe.Pointer)

//go:noescape
func radix4_256_fft_asm(data unsafe.Pointer)

//go:noescape
func radix4_512_fft_asm(data unsafe.Pointer)

//go:noescape
func radix4_1024_fft_asm(data unsafe.Pointer)

//go:noescape
func radixn_6_fft_asm(data unsafe.Pointer)

//go:noescape
func radixn_10_fft_asm(data unsafe.Pointer)

//go:noescape
func radixn_12_fft_asm(data unsafe.Pointer)

//go:noescape
func radixn_15_fft_asm(data unsafe.Pointer)

//go:noescape
func radixn_18_fft_asm(data unsafe.Pointer)

//go:noescape
func radixn_20_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_37_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_41_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_43_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_47_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_53_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_59_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_61_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_67_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_71_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_73_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_79_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_83_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_89_fft_asm(data unsafe.Pointer)

//go:noescape
func raders_97_fft_asm(data unsafe.Pointer)

//go:noescape
func bluestein_128_fft_asm(data unsafe.Pointer)

//go:noescape
func mixed_radix_60_fft_asm(data unsafe.Pointer)

//go:noescape
func mixed_radix_120_fft_asm(data unsafe.Pointer)

//go:noescape
func mixed_radix_240_fft_asm(data unsafe.Pointer)

//go:noescape
func good_thomas_35_fft_asm(data unsafe.Pointer)

//go:noescape
func good_thomas_77_fft_asm(data unsafe.Pointer)

//go:noescape
func winograd_49_fft_asm(data unsafe.Pointer)

//go:noescape
func winograd_121_fft_asm(data unsafe.Pointer)
