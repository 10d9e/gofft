package gofft

// Transpose performs an out-of-place matrix transpose
// data is treated as a rows x cols matrix stored in row-major order
func Transpose(input, output []complex128, rows, cols int) {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			output[c*rows+r] = input[r*cols+c]
		}
	}
}

// Transpose32 performs an out-of-place matrix transpose for complex64
func Transpose32(input, output []complex64, rows, cols int) {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			output[c*rows+r] = input[r*cols+c]
		}
	}
}

// TransposeInplace performs an in-place matrix transpose for square matrices
func TransposeInplace(data []complex128, n int) {
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			data[i*n+j], data[j*n+i] = data[j*n+i], data[i*n+j]
		}
	}
}

// BitReverse performs a bit-reversal permutation on the input
func BitReverse(data []complex128, logn int) {
	n := 1 << logn
	for i := 0; i < n; i++ {
		j := reverseBits(i, logn)
		if j > i {
			data[i], data[j] = data[j], data[i]
		}
	}
}

// reverseBits reverses the bottom n bits of x
func reverseBits(x, n int) int {
	result := 0
	for i := 0; i < n; i++ {
		result = (result << 1) | (x & 1)
		x >>= 1
	}
	return result
}

// BitReverse32 performs a bit-reversal permutation on complex64 input
func BitReverse32(data []complex64, logn int) {
	n := 1 << logn
	for i := 0; i < n; i++ {
		j := reverseBits(i, logn)
		if j > i {
			data[i], data[j] = data[j], data[i]
		}
	}
}

// Fill fills a slice with a constant value
func Fill(data []complex128, value complex128) {
	for i := range data {
		data[i] = value
	}
}

// Fill32 fills a complex64 slice with a constant value
func Fill32(data []complex64, value complex64) {
	for i := range data {
		data[i] = value
	}
}

// Copy copies data from src to dst
func Copy(dst, src []complex128) {
	copy(dst, src)
}

// Copy32 copies complex64 data from src to dst
func Copy32(dst, src []complex64) {
	copy(dst, src)
}
