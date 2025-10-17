package gofft

import (
	"math"
	"math/cmplx"
	"testing"
)

// Test helper to compare complex slices with tolerance
func complexSlicesEqual(a, b []complex128, tolerance float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if cmplx.Abs(a[i]-b[i]) > tolerance {
			return false
		}
	}
	return true
}

// naiveDFT computes a naive DFT for testing
func naiveDFT(input []complex128, forward bool) []complex128 {
	n := len(input)
	output := make([]complex128, n)

	sign := -1.0
	if !forward {
		sign = 1.0
	}

	for k := 0; k < n; k++ {
		sum := complex(0, 0)
		for j := 0; j < n; j++ {
			angle := sign * 2.0 * math.Pi * float64(k*j) / float64(n)
			twiddle := cmplx.Exp(complex(0, angle))
			sum += input[j] * twiddle
		}
		output[k] = sum
	}

	return output
}

func TestFFT_PowerOfTwo(t *testing.T) {
	sizes := []int{2, 4, 8, 16, 32, 64}

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create test signal
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(float64(i), float64(i)*0.5)
			}

			// Compute expected output using naive DFT
			expected := naiveDFT(input, true)

			// Compute using our FFT
			buffer := make([]complex128, n)
			copy(buffer, input)

			planner := NewPlanner()
			fft := planner.PlanForward(n)
			fft.Process(buffer)

			// Compare results
			if !complexSlicesEqual(buffer, expected, 1e-10) {
				t.Errorf("FFT output doesn't match expected for size %d", n)
				for i := range buffer {
					if cmplx.Abs(buffer[i]-expected[i]) > 1e-10 {
						t.Errorf("  [%d] got %v, want %v, diff %v", i, buffer[i], expected[i], buffer[i]-expected[i])
					}
				}
			}
		})
	}
}

func TestFFT_InverseProperty(t *testing.T) {
	sizes := []int{2, 4, 8, 16}

	for _, n := range sizes {
		t.Run("Size"+string(rune(n+'0')), func(t *testing.T) {
			// Create test signal
			input := make([]complex128, n)
			for i := range input {
				input[i] = complex(math.Sin(float64(i)), math.Cos(float64(i)))
			}
			original := make([]complex128, n)
			copy(original, input)

			planner := NewPlanner()

			// Forward FFT
			forward := planner.PlanForward(n)
			forward.Process(input)

			// Inverse FFT
			inverse := planner.PlanInverse(n)
			inverse.Process(input)

			// Normalize
			for i := range input {
				input[i] /= complex(float64(n), 0)
			}

			// Check if we got back the original
			if !complexSlicesEqual(input, original, 1e-10) {
				t.Errorf("Forward + Inverse didn't recover original signal for size %d", n)
				for i := range input {
					if cmplx.Abs(input[i]-original[i]) > 1e-10 {
						t.Errorf("  [%d] got %v, want %v", i, input[i], original[i])
					}
				}
			}
		})
	}
}

func TestFFT_DCValue(t *testing.T) {
	n := 16

	// Create a constant signal
	input := make([]complex128, n)
	value := complex(3.0, 2.0)
	for i := range input {
		input[i] = value
	}

	planner := NewPlanner()
	fft := planner.PlanForward(n)
	fft.Process(input)

	// DC component should be n * value
	expected := complex(float64(n), 0) * value
	if cmplx.Abs(input[0]-expected) > 1e-10 {
		t.Errorf("DC component incorrect: got %v, want %v", input[0], expected)
	}

	// All other components should be near zero
	for i := 1; i < n; i++ {
		if cmplx.Abs(input[i]) > 1e-10 {
			t.Errorf("Non-DC component [%d] should be near zero: got %v", i, input[i])
		}
	}
}

func TestFFT_ImpulseResponse(t *testing.T) {
	n := 8

	// Create an impulse signal (1 at index 0, 0 elsewhere)
	input := make([]complex128, n)
	input[0] = 1.0

	planner := NewPlanner()
	fft := planner.PlanForward(n)
	fft.Process(input)

	// FFT of impulse should be all ones
	for i := range input {
		expected := complex(1.0, 0)
		if cmplx.Abs(input[i]-expected) > 1e-10 {
			t.Errorf("Impulse response [%d] incorrect: got %v, want %v", i, input[i], expected)
		}
	}
}

func TestPrimeFactors(t *testing.T) {
	tests := []struct {
		n        int
		expected PrimeFactors
	}{
		{1, PrimeFactors{powerOfTwo: 0, powerOfThree: 0, otherFactors: nil, product: 1}},
		{2, PrimeFactors{powerOfTwo: 1, powerOfThree: 0, otherFactors: nil, product: 2}},
		{4, PrimeFactors{powerOfTwo: 2, powerOfThree: 0, otherFactors: nil, product: 4}},
		{6, PrimeFactors{powerOfTwo: 1, powerOfThree: 1, otherFactors: nil, product: 6}},
		{12, PrimeFactors{powerOfTwo: 2, powerOfThree: 1, otherFactors: nil, product: 12}},
	}

	for _, tt := range tests {
		t.Run("Factor"+string(rune(tt.n+'0')), func(t *testing.T) {
			result := ComputePrimeFactors(tt.n)
			if result.GetPowerOfTwo() != tt.expected.GetPowerOfTwo() {
				t.Errorf("Power of 2 mismatch: got %d, want %d", result.GetPowerOfTwo(), tt.expected.GetPowerOfTwo())
			}
			if result.GetPowerOfThree() != tt.expected.GetPowerOfThree() {
				t.Errorf("Power of 3 mismatch: got %d, want %d", result.GetPowerOfThree(), tt.expected.GetPowerOfThree())
			}
			if result.GetProduct() != tt.expected.GetProduct() {
				t.Errorf("Product mismatch: got %d, want %d", result.GetProduct(), tt.expected.GetProduct())
			}
		})
	}
}

func TestTwiddleFactors(t *testing.T) {
	n := 8
	twiddles := ComputeTwiddles(n, Forward)

	if len(twiddles) != n {
		t.Fatalf("Expected %d twiddle factors, got %d", n, len(twiddles))
	}

	// Twiddle[0] should be 1
	if cmplx.Abs(twiddles[0]-1.0) > 1e-10 {
		t.Errorf("Twiddle[0] should be 1, got %v", twiddles[0])
	}

	// Twiddle[n/2] should be -1 for forward FFT
	if cmplx.Abs(twiddles[n/2]-(-1.0)) > 1e-10 {
		t.Errorf("Twiddle[n/2] should be -1, got %v", twiddles[n/2])
	}

	// Twiddle[n/4] should be -i for forward FFT
	expectedQuarter := complex(0, -1)
	if cmplx.Abs(twiddles[n/4]-expectedQuarter) > 1e-10 {
		t.Errorf("Twiddle[n/4] should be -i, got %v", twiddles[n/4])
	}
}

func BenchmarkFFT(b *testing.B) {
	sizes := []int{64, 256, 1024, 4096}

	for _, n := range sizes {
		b.Run("Size"+string(rune(n+'0')), func(b *testing.B) {
			planner := NewPlanner()
			fft := planner.PlanForward(n)
			buffer := make([]complex128, n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fft.Process(buffer)
			}
		})
	}
}
