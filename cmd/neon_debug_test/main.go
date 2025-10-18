package main

import (
	"fmt"
	"math"
	"time"

	"github.com/10d9e/gofft"
)

func main() {
	fmt.Println("🔍 NEON Debug Test")
	fmt.Println("==================")

	// Test a few specific sizes
	sizes := []int{16, 32, 64, 128}

	planner := gofft.NewPlanner()

	for _, size := range sizes {
		fmt.Printf("\n📊 Testing size %d:\n", size)

		// Create FFT
		fft := planner.PlanForward(size)

		// Create test signal
		buffer := make([]complex128, size)
		for i := range buffer {
			angle := 2.0 * math.Pi * float64(i) / float64(size)
			buffer[i] = complex(math.Sin(angle), 0)
		}

		// Benchmark
		iterations := 10000
		start := time.Now()
		for i := 0; i < iterations; i++ {
			fft.Process(buffer)
		}
		elapsed := time.Since(start)
		timePerOp := float64(elapsed.Nanoseconds()) / float64(iterations)

		fmt.Printf("  Time per operation: %.1f ns\n", timePerOp)
		fmt.Printf("  FFT type: %T\n", fft)

		// Check if it's using NEON (by checking the type name)
		fftType := fmt.Sprintf("%T", fft)
		if fftType == "*gofft.neonButterflyAdapter" {
			fmt.Printf("  ✅ Using NEON adapter\n")
		} else {
			fmt.Printf("  ❌ Using scalar implementation (%s)\n", fftType)
		}
	}
}
