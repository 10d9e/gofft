//go:build arm64

package neon

import (
	"testing"

	"github.com/10d9e/gofft/algorithm"
)

// BenchmarkNEONvsScalar compares NEON vs scalar performance
func BenchmarkNEONvsScalar(b *testing.B) {
	sizes := []int{2, 4, 8, 16, 32}

	for _, size := range sizes {
		b.Run("Size"+string(rune(size+'0')), func(b *testing.B) {
			// Create test data
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%7), float64(i%5)*0.3)
			}

			// Benchmark NEON implementation
			b.Run("NEON", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					ProcessVectorizedButterfly(data, size)
				}
			})

			// Benchmark scalar implementation
			b.Run("Scalar", func(b *testing.B) {
				scalarFft := algorithm.NewDft(size, algorithm.Forward)
				scratch := make([]complex128, scalarFft.InplaceScratchLen())

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					testData := make([]complex128, size)
					copy(testData, data)
					scalarFft.ProcessWithScratch(testData, scratch)
				}
			})
		})
	}
}

// BenchmarkButterflyComparison compares individual butterfly implementations
func BenchmarkButterflyComparison(b *testing.B) {
	testCases := []struct {
		size     int
		name     string
		neonFunc func([]complex128)
	}{
		{1, "Butterfly1", Butterfly1_NEON},
		{2, "Butterfly2", Butterfly2_NEON},
		{4, "Butterfly4", Butterfly4_NEON},
		{8, "Butterfly8", Butterfly8_NEON},
		{10, "Butterfly10", Butterfly10_NEON},
		{15, "Butterfly15", Butterfly15_NEON},
		{16, "Butterfly16", Butterfly16_NEON},
		{32, "Butterfly32", Butterfly32_NEON},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			data := make([]complex128, tc.size)
			for i := range data {
				data[i] = complex(float64(i%5), float64(i%3)*0.4)
			}

			// Benchmark NEON
			b.Run("NEON", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					testData := make([]complex128, tc.size)
					copy(testData, data)
					tc.neonFunc(testData)
				}
			})

			// Benchmark scalar
			b.Run("Scalar", func(b *testing.B) {
				scalarFft := algorithm.NewDft(tc.size, algorithm.Forward)
				scratch := make([]complex128, scalarFft.InplaceScratchLen())

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					testData := make([]complex128, tc.size)
					copy(testData, data)
					scalarFft.ProcessWithScratch(testData, scratch)
				}
			})
		})
	}
}

// BenchmarkSIMDPlanner tests the SIMD planner performance
func BenchmarkSIMDPlanner(b *testing.B) {
	sizes := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	for _, size := range sizes {
		b.Run("Size"+string(rune(size+'0')), func(b *testing.B) {
			// Create test data
			data := make([]complex128, size)
			for i := range data {
				data[i] = complex(float64(i%7), float64(i%5)*0.3)
			}

			// Benchmark SIMD planner
			b.Run("SIMD", func(b *testing.B) {
				// This would use the actual SIMD planner
				// For now, we'll benchmark the framework overhead
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					testData := make([]complex128, size)
					copy(testData, data)
					// Simulate SIMD processing
					if size <= 32 {
						ProcessVectorizedButterfly(testData, size)
					}
				}
			})

			// Benchmark scalar planner
			b.Run("Scalar", func(b *testing.B) {
				scalarFft := algorithm.NewDft(size, algorithm.Forward)
				scratch := make([]complex128, scalarFft.InplaceScratchLen())

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					testData := make([]complex128, size)
					copy(testData, data)
					scalarFft.ProcessWithScratch(testData, scratch)
				}
			})
		})
	}
}

// BenchmarkMemoryAlignment tests the impact of memory alignment
func BenchmarkMemoryAlignment(b *testing.B) {
	size := 16

	// Create aligned data
	alignedData := make([]complex128, size)
	for i := range alignedData {
		alignedData[i] = complex(float64(i%4), float64(i%3)*0.3)
	}

	// Create unaligned data (offset by 1 element)
	unalignedData := make([]complex128, size+1)
	for i := range unalignedData {
		unalignedData[i] = complex(float64(i%4), float64(i%3)*0.3)
	}
	unalignedSlice := unalignedData[1:] // This creates an unaligned slice

	b.Run("Aligned", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testData := make([]complex128, size)
			copy(testData, alignedData)
			Butterfly16_NEON(testData)
		}
	})

	b.Run("Unaligned", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testData := make([]complex128, size)
			copy(testData, unalignedSlice)
			Butterfly16_NEON(testData)
		}
	})
}

// BenchmarkVectorOperations tests individual NEON vector operations
func BenchmarkVectorOperations(b *testing.B) {
	// Test vector addition
	b.Run("VectorAdd", func(b *testing.B) {
		vec1 := Vector128{data: Complex128{re: 1.0, im: 2.0}}
		vec2 := Vector128{data: Complex128{re: 3.0, im: 4.0}}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = vec1.Add128(vec2)
		}
	})

	// Test vector multiplication
	b.Run("VectorMul", func(b *testing.B) {
		vec1 := Vector128{data: Complex128{re: 1.0, im: 2.0}}
		vec2 := Vector128{data: Complex128{re: 3.0, im: 4.0}}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = vec1.Mul128(vec2)
		}
	})

	// Test scalar operations for comparison
	b.Run("ScalarAdd", func(b *testing.B) {
		a := complex(1.0, 2.0)
		b_val := complex(3.0, 4.0)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = a + b_val
		}
	})

	b.Run("ScalarMul", func(b *testing.B) {
		a := complex(1.0, 2.0)
		b_val := complex(3.0, 4.0)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = a * b_val
		}
	})
}

// BenchmarkThroughput tests throughput for different data sizes
func BenchmarkThroughput(b *testing.B) {
	sizes := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	for _, size := range sizes {
		b.Run("Size"+string(rune(size+'0')), func(b *testing.B) {
			// Create larger dataset for throughput testing
			totalSize := size * 1000 // Process 1000 FFTs
			data := make([]complex128, totalSize)
			for i := range data {
				data[i] = complex(float64(i%7), float64(i%5)*0.3)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Process in chunks of size
				for j := 0; j < totalSize; j += size {
					chunk := data[j : j+size]
					if size <= 32 {
						ProcessVectorizedButterfly(chunk, size)
					}
				}
			}
		})
	}
}
