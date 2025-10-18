package main

import (
	"fmt"
	"math"
	"time"

	"github.com/10d9e/gofft"
	"github.com/10d9e/gofft/simd"
)

func main() {
	fmt.Println("🔍 NEON Coverage Analysis")
	fmt.Println("========================")

	// Test sizes that have NEON implementations
	neonSizes := []int{
		// Butterflies
		2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15, 16, 17, 19, 23, 24, 27, 29, 31, 32,
		// Radix-4
		64, 128, 256, 512, 1024,
		// RadixN
		6, 10, 12, 15, 18, 20,
		// Rader's
		37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97,
		// Mixed-Radix
		60, 120, 240, 480,
		// Good-Thomas
		35, 77, 143, 221,
		// Winograd
		49, 121, 169, 289,
		// Bluestein
		128, 256, 512, 1024,
	}

	// Remove duplicates and sort
	uniqueSizes := make(map[int]bool)
	for _, size := range neonSizes {
		uniqueSizes[size] = true
	}

	fmt.Printf("%-8s %-15s %-15s %-12s %-10s\n",
		"Size", "Main Planner", "SIMD Planner", "Speedup", "Coverage")
	fmt.Println("----------------------------------------------------------------")

	mainPlanner := gofft.NewPlanner()
	simdPlanner := simd.NewSIMDPlanner()

	covered := 0
	total := len(uniqueSizes)

	for size := range uniqueSizes {
		// Test main planner
		mainFft := mainPlanner.PlanForward(size)
		mainTime := benchmarkFFT(mainFft, size, 1000)

		// Test SIMD planner
		simdFft := simdPlanner.PlanForward(size)
		simdTime := benchmarkFFT(simdFft, size, 1000)

		// Calculate speedup
		speedup := mainTime / simdTime

		// Determine coverage
		coverage := "❌ Scalar"
		if speedup > 1.5 {
			coverage = "✅ NEON"
			covered++
		} else if speedup > 1.1 {
			coverage = "⚠️  Partial"
		}

		fmt.Printf("%-8d %-15.1f %-15.1f %-12.1fx %-10s\n",
			size, mainTime, simdTime, speedup, coverage)
	}

	fmt.Printf("\n📊 Coverage Summary:\n")
	fmt.Printf("Total NEON sizes: %d\n", total)
	fmt.Printf("Covered by main planner: %d (%.1f%%)\n", covered, float64(covered)/float64(total)*100)
	fmt.Printf("Missing from main planner: %d (%.1f%%)\n", total-covered, float64(total-covered)/float64(total)*100)

	fmt.Printf("\n🎯 Missing NEON Features in Main Planner:\n")
	fmt.Printf("• Butterflies: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15, 17, 19, 23, 24, 27, 29, 31\n")
	fmt.Printf("• RadixN: 6, 10, 12, 15, 18, 20\n")
	fmt.Printf("• Rader's: 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97\n")
	fmt.Printf("• Mixed-Radix: 60, 120, 240, 480\n")
	fmt.Printf("• Good-Thomas: 35, 77, 143, 221\n")
	fmt.Printf("• Winograd: 49, 121, 169, 289\n")
	fmt.Printf("• Bluestein: 128, 256, 512, 1024\n")
}

func benchmarkFFT(fft interface{}, size int, iterations int) float64 {
	// Create test signal
	buffer := make([]complex128, size)
	for i := range buffer {
		angle := 2.0 * math.Pi * float64(i) / float64(size)
		buffer[i] = complex(math.Sin(angle), 0)
	}

	// Warm up
	for i := 0; i < 5; i++ {
		processFFT(fft, buffer)
	}

	// Benchmark
	start := time.Now()
	for i := 0; i < iterations; i++ {
		processFFT(fft, buffer)
	}
	elapsed := time.Since(start)

	return float64(elapsed.Nanoseconds()) / float64(iterations)
}

func processFFT(fft interface{}, buffer []complex128) {
	switch f := fft.(type) {
	case gofft.Fft:
		f.Process(buffer)
	case interface{ Process([]complex128) }:
		f.Process(buffer)
	default:
		panic("Unknown FFT type")
	}
}
