package gofft

import (
	"errors"
	"math/bits"

	"github.com/example/gofft/pkg/gofft/threadpool"
)

// FFT computes the forward FFT of x (len must be a power of two).
// It returns a new slice with the spectrum.
func FFT(x []complex128) ([]complex128, error) {
	return doFFT(x, false)
}

// IFFT computes the inverse FFT of X (len must be a power of two).
// It returns a new slice with time-domain samples (scaled by 1/N).
func IFFT(X []complex128) ([]complex128, error) {
	return doFFT(X, true)
}

func doFFT(in []complex128, invert bool) ([]complex128, error) {
	n := len(in)
	if n == 0 {
		return []complex128{}, nil
	}
	if n&(n-1) != 0 {
		return nil, errors.New("length must be power of two")
	}

	// bit-reverse copy
	buf := make([]complex128, n)
	for i := 0; i < n; i++ {
		r := bits.Reverse(uint(i)) >> (bits.UintSize - uint(bits.Len(uint(n-1))))
		buf[int(r)] = in[i]
	}

	// stages - use pure radix-2 for correctness
	// (radix-4/8 optimizations would require proper mixed-radix algorithm)
	for m := 2; m <= n; m <<= 1 {
		butterfly2InPlace(buf, m, invert)
	}

	// scale on inverse
	if invert {
		invN := 1.0 / float64(n)
		for i := range buf {
			buf[i] *= complex(invN, 0)
		}
	}
	return buf, nil
}

// FFTParallel runs FFT using a worker pool for stage-level parallelism.
func FFTParallel(x []complex128, workers int) ([]complex128, error) {
	return doFFTParallel(x, false, workers)
}

func IFFTParallel(X []complex128, workers int) ([]complex128, error) {
	return doFFTParallel(X, true, workers)
}

func doFFTParallel(in []complex128, invert bool, workers int) ([]complex128, error) {
	n := len(in)
	if n == 0 {
		return []complex128{}, nil
	}
	if n&(n-1) != 0 {
		return nil, errors.New("length must be power of two")
	}

	buf := make([]complex128, n)
	for i := 0; i < n; i++ {
		r := bits.Reverse(uint(i)) >> (bits.UintSize - uint(bits.Len(uint(n-1))))
		buf[int(r)] = in[i]
	}

	pool := threadpool.New(workers)
	defer pool.Close()

	// use pure radix-2 for correctness
	for m := 2; m <= n; m <<= 1 {
		parallelButterfly(pool, 2, buf, m, invert)
	}

	if invert {
		invN := 1.0 / float64(n)
		for i := range buf {
			buf[i] *= complex(invN, 0)
		}
	}
	return buf, nil
}
