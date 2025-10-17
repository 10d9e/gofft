package gofft

import (
	"math"
	"math/rand"
	"testing"
)

func dft(x []complex128, invert bool) []complex128 {
	n := len(x)
	X := make([]complex128, n)
	for k := 0; k < n; k++ {
		var sum complex128
		for n0 := 0; n0 < n; n0++ {
			angle := 2 * math.Pi * float64(k*n0) / float64(n)
			if !invert { angle = -angle }
			w := complex(math.Cos(angle), math.Sin(angle))
			sum += x[n0] * w
		}
		if invert { sum /= complex(float64(n), 0) }
		X[k] = sum
	}
	return X
}

func almostEqual(a, b complex128, eps float64) bool {
	return cmplxAbs(a-b) <= eps
}

func cmplxAbs(a complex128) float64 {
	return math.Hypot(real(a), imag(a))
}

func TestRoundTrip(t *testing.T) {
	sizes := []int{2,4,8,16,32,64,128}
	for _, n := range sizes {
		x := make([]complex128, n)
		for i := range x {
			x[i] = complex(rand.NormFloat64(), rand.NormFloat64())
		}
		X, err := FFT(x)
		if err != nil { t.Fatalf("fft err: %v", err) }
		y, err := IFFT(X)
		if err != nil { t.Fatalf("ifft err: %v", err) }
		for i := range x {
			if !almostEqual(x[i], y[i], 1e-9) {
				t.Fatalf("roundtrip mismatch n=%d i=%d want=%v got=%v", n, i, x[i], y[i])
			}
		}
	}
}

func TestAgainstDFT(t *testing.T) {
	sizes := []int{8,16,32}
	for _, n := range sizes {
		x := make([]complex128, n)
		for i := range x {
			x[i] = complex(math.Sin(2*math.Pi*float64(i)/float64(n))+0.1, 0.3*math.Cos(2*math.Pi*float64(i)/float64(n)))
		}
		X1, err := FFT(x)
		if err != nil { t.Fatalf("fft err: %v", err) }
		X2 := dft(x, false)
		for k := range x {
			if !almostEqual(X1[k], X2[k], 1e-9) {
				t.Fatalf("DFT mismatch n=%d k=%d fft=%v dft=%v", n, k, X1[k], X2[k])
			}
		}
	}
}
