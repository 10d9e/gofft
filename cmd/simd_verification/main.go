package main

import (
	"fmt"
	"math"
	"time"

	"github.com/10d9e/gofft"
	"github.com/10d9e/gofft/simd"
)

func main() {
	fmt.Println("🔍 SIMD Verification Test")
	fmt.Println("=========================")

	// Check SIMD detection
	simdLevel := simd.DetectSIMD()
	fmt.Printf("Detected SIMD Level: %s\n", simdLevel)

	arch, os, simdLevel := simd.PlatformInfo()
	fmt.Printf("Platform: %s/%s, SIMD: %s\n", arch, os, simdLevel)

	// Test both planners
	sizes := []int{16, 32, 64, 128}

	fmt.Println("\n📊 Performance Comparison:")
	fmt.Printf("%-8s %-15s %-15s %-12s\n",
		"Size", "Scalar (ns/op)", "SIMD (ns/op)", "Speedup")
	fmt.Println("------------------------------------------------")

	for _, size := range sizes {
		// Test Scalar planner (current main planner)
		scalarPlanner := gofft.NewPlanner()
		scalarFft := scalarPlanner.PlanForward(size)

		// Test SIMD planner
		simdPlanner := simd.NewSIMDPlanner()
		simdFft := simdPlanner.PlanForward(size)

		// Create test signal
		buffer := make([]complex128, size)
		for i := range buffer {
			angle := 2.0 * math.Pi * float64(i) / float64(size)
			buffer[i] = complex(math.Sin(angle), 0)
		}

		// Benchmark Scalar
		scalarTime := benchmarkFFT(scalarFft, buffer, 1000)

		// Benchmark SIMD
		simdTime := benchmarkFFT(simdFft, buffer, 1000)

		// Calculate speedup
		speedup := scalarTime / simdTime

		fmt.Printf("%-8d %-15.1f %-15.1f %-12.1fx\n",
			size, scalarTime, simdTime, speedup)
	}

	fmt.Println("\n🎯 Key Findings:")
	if simdLevel == simd.NEON {
		fmt.Println("✅ NEON SIMD detected on ARM64")
		fmt.Println("✅ SIMD planner should provide better performance")
	} else if simdLevel == simd.Scalar {
		fmt.Println("⚠️  Only scalar algorithms available")
		fmt.Println("⚠️  Main planner is using scalar algorithms")
	} else {
		fmt.Printf("ℹ️  SIMD level: %s\n", simdLevel)
	}

	fmt.Println("\n💡 To use SIMD:")
	fmt.Println("   Use: simd.NewSIMDPlanner() instead of gofft.NewPlanner()")
	fmt.Println("   Example:")
	fmt.Println("     planner := simd.NewSIMDPlanner()")
	fmt.Println("     fft := planner.PlanForward(1024)")
}

func benchmarkFFT(fft interface{}, buffer []complex128, iterations int) float64 {
	// Create a copy for each iteration
	testBuffer := make([]complex128, len(buffer))

	// Warm up
	for i := 0; i < 5; i++ {
		copy(testBuffer, buffer)
		processFFT(fft, testBuffer)
	}

	// Benchmark
	start := time.Now()
	for i := 0; i < iterations; i++ {
		copy(testBuffer, buffer)
		processFFT(fft, testBuffer)
	}
	elapsed := time.Since(start)

	return float64(elapsed.Nanoseconds()) / float64(iterations)
}

func processFFT(fft interface{}, buffer []complex128) {
	// Handle both gofft.Fft and simd.FftInterface
	switch f := fft.(type) {
	case gofft.Fft:
		f.Process(buffer)
	case interface{ Process([]complex128) }:
		f.Process(buffer)
	default:
		panic("Unknown FFT type")
	}
}
