package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"time"

	"github.com/10d9e/gofft"
)

// Pure Go DFT implementation for comparison
type PureGoDFT struct {
	length    int
	direction gofft.Direction
}

func NewPureGoDFT(length int, direction gofft.Direction) *PureGoDFT {
	return &PureGoDFT{length: length, direction: direction}
}

func (d *PureGoDFT) Len() int                   { return d.length }
func (d *PureGoDFT) Direction() gofft.Direction { return d.direction }
func (d *PureGoDFT) ImmutableScratchLen() int   { return 0 }
func (d *PureGoDFT) InplaceScratchLen() int     { return 0 }
func (d *PureGoDFT) OutOfPlaceScratchLen() int  { return 0 }
func (d *PureGoDFT) Process(buffer []complex128) {
	// Pure Go DFT implementation - O(n²) complexity
	output := make([]complex128, d.length)

	for k := 0; k < d.length; k++ {
		var sum complex128
		for n := 0; n < d.length; n++ {
			angle := -2.0 * math.Pi * float64(k*n) / float64(d.length)
			if d.direction == gofft.Inverse {
				angle = -angle
			}
			sum += buffer[n] * cmplx.Exp(complex(0, angle))
		}
		output[k] = sum
	}

	copy(buffer, output)
}

func (d *PureGoDFT) ProcessImmutable(input []complex128, output []complex128, scratch []complex128) {
	copy(output, input)
	d.Process(output)
}

func (d *PureGoDFT) ProcessOutOfPlace(input []complex128, output []complex128, scratch []complex128) {
	copy(output, input)
	d.Process(output)
}

func (d *PureGoDFT) ProcessWithScratch(buffer []complex128, scratch []complex128) {
	d.Process(buffer)
}

func main() {
	fmt.Println("🚀 gofft SIMD vs Pure Go DFT Performance Comparison")
	fmt.Println("====================================================")

	// Test sizes - keeping small for reasonable comparison times
	sizes := []int{16, 32, 64, 128}

	fmt.Printf("%-8s %-15s %-15s %-12s %-10s\n",
		"Size", "SIMD (ns/op)", "Pure Go (ns/op)", "Speedup", "Algorithm")
	fmt.Println("----------------------------------------------------------------")

	for _, size := range sizes {
		// Test SIMD version (automatic - uses NEON on ARM64)
		planner := gofft.NewPlanner()
		simdFft := planner.PlanForward(size)

		// Test Pure Go version (O(n²) DFT)
		pureGoFft := NewPureGoDFT(size, gofft.Forward)

		// Create test signal
		buffer := make([]complex128, size)
		for i := range buffer {
			angle := 2.0 * math.Pi * float64(i) / float64(size)
			buffer[i] = complex(math.Sin(angle), 0)
		}

		// Benchmark SIMD version
		simdTime := benchmarkFFT(simdFft, buffer, 1000)

		// Benchmark Pure Go version (fewer iterations due to O(n²) complexity)
		iterations := 100
		if size >= 64 {
			iterations = 10 // Even fewer for larger sizes
		}
		pureGoTime := benchmarkFFT(pureGoFft, buffer, iterations)

		// Calculate speedup
		speedup := pureGoTime / simdTime

		// Determine algorithm type
		algorithm := "NEON Butterfly"
		if size > 32 {
			algorithm = "NEON Radix-4"
		}

		fmt.Printf("%-8d %-15.1f %-15.1f %-12.1fx %-10s\n",
			size, simdTime, pureGoTime, speedup, algorithm)
	}

	fmt.Println("\n🎯 Key Insights:")
	fmt.Println("• SIMD versions: O(n log n) with real ARM64 NEON assembly")
	fmt.Println("• Pure Go versions: O(n²) naive DFT implementation")
	fmt.Println("• Speedup shows the benefit of both algorithm optimization AND SIMD")
	fmt.Println("• Real-world performance gains are even higher with optimized algorithms!")
	fmt.Println("• SIMD provides both algorithmic and hardware acceleration")
}

func benchmarkFFT(fft gofft.Fft, buffer []complex128, iterations int) float64 {
	// Create a copy for each iteration
	testBuffer := make([]complex128, len(buffer))

	// Warm up
	for i := 0; i < 5; i++ {
		copy(testBuffer, buffer)
		fft.Process(testBuffer)
	}

	// Benchmark
	start := time.Now()
	for i := 0; i < iterations; i++ {
		copy(testBuffer, buffer)
		fft.Process(testBuffer)
	}
	elapsed := time.Since(start)

	return float64(elapsed.Nanoseconds()) / float64(iterations)
}
