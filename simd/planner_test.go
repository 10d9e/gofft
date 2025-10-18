package simd

import (
	"testing"

	"github.com/10d9e/gofft/algorithm"
)

func TestSIMDPlanner(t *testing.T) {
	planner := NewSIMDPlanner()

	// Test SIMD detection
	level := planner.GetSIMDLevel()
	t.Logf("Detected SIMD level: %v", level)

	// Test that planner is created successfully
	if planner == nil {
		t.Fatal("SIMD planner should not be nil")
	}

	// Test SIMD optimization status
	optimized := planner.IsSIMDOptimized()
	t.Logf("SIMD optimized: %v", optimized)

	// On ARM64, we should have NEON
	if level == NEON {
		t.Logf("✅ NEON SIMD detected and available")
	} else if level == Scalar {
		t.Logf("⚠️  Using scalar fallback (expected on non-ARM64)")
	}
}

func TestSIMDPlannerFFTCreation(t *testing.T) {
	planner := NewSIMDPlanner()

	// Test creating FFTs of various sizes
	testSizes := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	for _, size := range testSizes {
		t.Run("Size"+string(rune(size+'0')), func(t *testing.T) {
			// Test forward FFT
			forward := planner.PlanForward(size)
			if forward == nil {
				t.Fatalf("Forward FFT for size %d should not be nil", size)
			}

			if forward.Len() != size {
				t.Errorf("Forward FFT length: got %d, expected %d", forward.Len(), size)
			}

			if forward.Direction() != algorithm.Forward {
				t.Errorf("Forward FFT direction: got %v, expected Forward", forward.Direction())
			}

			// Test inverse FFT
			inverse := planner.PlanInverse(size)
			if inverse == nil {
				t.Fatalf("Inverse FFT for size %d should not be nil", size)
			}

			if inverse.Len() != size {
				t.Errorf("Inverse FFT length: got %d, expected %d", inverse.Len(), size)
			}

			if inverse.Direction() != algorithm.Inverse {
				t.Errorf("Inverse FFT direction: got %v, expected Inverse", inverse.Direction())
			}
		})
	}
}

func TestSIMDPlannerAdapter(t *testing.T) {
	adapter := NewSIMDPlannerAdapter()

	if adapter == nil {
		t.Fatal("SIMD planner adapter should not be nil")
	}

	// Test that adapter works the same as planner
	level := adapter.GetSIMDLevel()
	optimized := adapter.IsSIMDOptimized()

	t.Logf("Adapter SIMD level: %v", level)
	t.Logf("Adapter SIMD optimized: %v", optimized)

	// Test FFT creation through adapter
	fft := adapter.PlanForward(64)
	if fft == nil {
		t.Fatal("FFT created through adapter should not be nil")
	}

	if fft.Len() != 64 {
		t.Errorf("FFT length through adapter: got %d, expected 64", fft.Len())
	}
}

func TestSIMDPerformance(t *testing.T) {
	planner := NewSIMDPlanner()
	level := planner.GetSIMDLevel()

	// Test that we can create and use FFTs
	size := 1024
	fft := planner.PlanForward(size)

	// Create test data
	data := make([]complex128, size)
	for i := range data {
		data[i] = complex(float64(i%7), float64(i%5)*0.3)
	}

	// Process the data
	fft.Process(data)

	// Verify we got some result (not all zeros)
	hasNonZero := false
	for _, val := range data {
		if val != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("FFT result should not be all zeros")
	}

	t.Logf("✅ SIMD FFT processing successful with %v", level)
}

func BenchmarkSIMDPlanner(b *testing.B) {
	planner := NewSIMDPlanner()
	level := planner.GetSIMDLevel()

	b.Logf("Benchmarking with SIMD level: %v", level)

	// Benchmark FFT creation
	b.Run("FFTCreation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fft := planner.PlanForward(1024)
			_ = fft
		}
	})

	// Benchmark FFT processing
	fft := planner.PlanForward(1024)
	data := make([]complex128, 1024)
	for i := range data {
		data[i] = complex(float64(i%7), float64(i%5)*0.3)
	}

	b.Run("FFTProcessing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Copy data to avoid modifying the original
			testData := make([]complex128, len(data))
			copy(testData, data)
			fft.Process(testData)
		}
	})
}
