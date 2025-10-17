// Example program demonstrating gofft usage
package main

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/example/gofft/pkg/gofft"
)

func main() {
	fmt.Println("gofft - Go FFT Library Example")

	// Example 1: Basic FFT
	basicExample()

	// Example 2: Forward and Inverse FFT
	roundTripExample()

	// Example 3: Frequency Analysis
	frequencyAnalysisExample()
}

func basicExample() {
	fmt.Println("Example 1: Basic FFT")
	fmt.Println("--------------------")

	// Create a planner
	planner := gofft.NewPlanner()

	// Plan a forward FFT of size 16
	size := 16
	fft := planner.PlanForward(size)

	// Create input signal (simple sine wave)
	buffer := make([]complex128, size)
	for i := 0; i < size; i++ {
		// Create a sine wave at frequency = 2 cycles over the buffer
		angle := 2.0 * 2.0 * math.Pi * float64(i) / float64(size)
		buffer[i] = complex(math.Sin(angle), 0)
	}

	fmt.Printf("Input: %d-point sine wave (2 cycles)\n", size)

	// Compute FFT
	fft.Process(buffer)

	// Find peak frequency
	maxMag := 0.0
	maxIdx := 0
	for i := 0; i < size; i++ {
		mag := cmplx.Abs(buffer[i])
		if mag > maxMag {
			maxMag = mag
			maxIdx = i
		}
	}

	fmt.Printf("Peak frequency bin: %d (magnitude: %.2f)\n", maxIdx, maxMag)
	fmt.Println()
}

func roundTripExample() {
	fmt.Println("Example 2: Forward + Inverse FFT")
	fmt.Println("---------------------------------")

	planner := gofft.NewPlanner()
	size := 16

	// Create original signal
	original := make([]complex128, size)
	for i := 0; i < size; i++ {
		original[i] = complex(float64(i), float64(i)*0.5)
	}

	// Make a copy
	buffer := make([]complex128, size)
	copy(buffer, original)

	// Forward FFT
	forward := planner.PlanForward(size)
	forward.Process(buffer)

	// Inverse FFT
	inverse := planner.PlanInverse(size)
	inverse.Process(buffer)

	// Normalize (FFT doesn't normalize by default)
	for i := range buffer {
		buffer[i] /= complex(float64(size), 0)
	}

	// Check accuracy
	maxError := 0.0
	for i := 0; i < size; i++ {
		error := cmplx.Abs(buffer[i] - original[i])
		if error > maxError {
			maxError = error
		}
	}

	fmt.Printf("Size: %d\n", size)
	fmt.Printf("Max reconstruction error: %.2e\n", maxError)
	if maxError < 1e-10 {
		fmt.Println("✓ Perfect reconstruction!")
	}
	fmt.Println()
}

func frequencyAnalysisExample() {
	fmt.Println("Example 3: Frequency Analysis")
	fmt.Println("------------------------------")

	planner := gofft.NewPlanner()
	size := 64

	// Create a signal with multiple frequency components
	buffer := make([]complex128, size)
	for i := 0; i < size; i++ {
		t := float64(i) / float64(size)

		// Mix of 3 sine waves at different frequencies
		val := 0.0
		val += 1.0 * math.Sin(2.0*math.Pi*3.0*t)  // 3 Hz
		val += 0.5 * math.Sin(2.0*math.Pi*7.0*t)  // 7 Hz
		val += 0.3 * math.Sin(2.0*math.Pi*15.0*t) // 15 Hz

		buffer[i] = complex(val, 0)
	}

	// Compute FFT
	fft := planner.PlanForward(size)
	fft.Process(buffer)

	// Find the top 3 frequency components (positive frequencies only)
	fmt.Printf("Top frequency components (out of %d bins):\n", size/2)

	type peak struct {
		bin int
		mag float64
	}
	peaks := make([]peak, 0)

	for i := 0; i < size/2; i++ {
		mag := cmplx.Abs(buffer[i])
		if mag > 0.1 { // Threshold to ignore noise
			peaks = append(peaks, peak{bin: i, mag: mag})
		}
	}

	// Sort by magnitude (simple bubble sort for small arrays)
	for i := 0; i < len(peaks); i++ {
		for j := i + 1; j < len(peaks); j++ {
			if peaks[j].mag > peaks[i].mag {
				peaks[i], peaks[j] = peaks[j], peaks[i]
			}
		}
	}

	// Print top 3
	for i := 0; i < 3 && i < len(peaks); i++ {
		fmt.Printf("  Bin %2d: magnitude %.2f\n", peaks[i].bin, peaks[i].mag)
	}
	fmt.Println()
}
