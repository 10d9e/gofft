package algorithm

import "math"

// Bluestein implements the Bluestein (chirp-Z) FFT algorithm
// This algorithm can compute FFTs of ANY size in O(n log n) time by
// converting the DFT into a convolution, which is then computed using
// power-of-two FFTs.
//
// Algorithm:
// 1. Compute chirp sequence: w[k] = exp(-i*π*k²/N)
// 2. Multiply input by chirp: x[k] * w[k]
// 3. Convolve with conjugate chirp using FFT
// 4. Multiply result by chirp: result[k] * w[k]
//
// This makes ANY size O(n log n), including large primes!
type Bluestein struct {
	length         int
	direction      Direction
	fftSize        int          // Power-of-two size >= 2*length-1
	fft            *Radix4      // Power-of-two FFT
	invFft         *Radix4      // Inverse FFT
	chirp          []complex128 // Chirp sequence w[k]
	chirpConj      []complex128 // Conjugate chirp for convolution
	chirpConvolved []complex128 // Pre-convolved chirp (FFT of padded conjugate chirp)
}

// NewBluestein creates a Bluestein FFT instance for arbitrary size
func NewBluestein(length int, direction Direction) *Bluestein {
	// Find next power of two >= 2*length-1
	minSize := 2*length - 1
	fftSize := 1
	for fftSize < minSize {
		fftSize *= 2
	}

	// Create power-of-two FFTs
	fft := NewRadix4(fftSize, Forward)
	invFft := NewRadix4(fftSize, Inverse)

	// Precompute chirp sequence: w[k] = exp(-i*π*k²/N)
	chirp := make([]complex128, length)
	chirpConj := make([]complex128, fftSize)

	for k := 0; k < length; k++ {
		// angle = -π*k²/N (or +π for inverse)
		angle := -math.Pi * float64(k*k) / float64(length)
		if direction == Inverse {
			angle = -angle
		}
		chirp[k] = complex(math.Cos(angle), math.Sin(angle))
		chirpConj[k] = complexConj(chirp[k])
	}

	// Pad chirpConj with wraparound
	// For convolution, we need chirpConj[k] for k=0..length-1 and k=fftSize-(length-1)..fftSize-1
	for k := 1; k < length; k++ {
		chirpConj[fftSize-k] = chirpConj[k]
	}

	// Precompute FFT of chirpConj for convolution
	chirpConvolved := make([]complex128, fftSize)
	copy(chirpConvolved, chirpConj)
	scratch := make([]complex128, fft.InplaceScratchLen())
	fft.ProcessWithScratch(chirpConvolved, scratch)

	return &Bluestein{
		length:         length,
		direction:      direction,
		fftSize:        fftSize,
		fft:            fft,
		invFft:         invFft,
		chirp:          chirp,
		chirpConj:      chirpConj,
		chirpConvolved: chirpConvolved,
	}
}

func (b *Bluestein) Len() int                  { return b.length }
func (b *Bluestein) Direction() Direction      { return b.direction }
func (b *Bluestein) InplaceScratchLen() int    { return 2 * b.fftSize }
func (b *Bluestein) OutOfPlaceScratchLen() int { return 2 * b.fftSize }

func (b *Bluestein) ProcessWithScratch(buffer, scratch []complex128) {
	// Process each chunk of size b.length
	for i := 0; i < len(buffer); i += b.length {
		chunk := buffer[i : i+b.length]
		workScratch := scratch[:2*b.fftSize]
		b.processOne(chunk, workScratch)
	}
}

func (b *Bluestein) ProcessOutOfPlace(input, output, scratch []complex128) {
	copy(output, input)
	b.ProcessWithScratch(output, scratch)
}

func (b *Bluestein) processOne(buffer, scratch []complex128) {
	// Allocate work buffers from scratch
	x := scratch[:b.fftSize]              // Input padded to fftSize
	y := scratch[b.fftSize : 2*b.fftSize] // Temporary for convolution

	// Clear buffers
	for i := range x {
		x[i] = 0
		y[i] = 0
	}

	// Step 1: Multiply input by chirp and pad
	for k := 0; k < b.length; k++ {
		x[k] = buffer[k] * b.chirp[k]
	}

	// Step 2: FFT of x
	fftScratch := make([]complex128, b.fft.InplaceScratchLen())
	b.fft.ProcessWithScratch(x, fftScratch)

	// Step 3: Pointwise multiply with pre-convolved chirp (convolution in frequency domain)
	for k := 0; k < b.fftSize; k++ {
		x[k] = x[k] * b.chirpConvolved[k]
	}

	// Step 4: Inverse FFT
	invScratch := make([]complex128, b.invFft.InplaceScratchLen())
	b.invFft.ProcessWithScratch(x, invScratch)

	// Step 5: Normalize (inverse FFT doesn't auto-normalize)
	scale := complex(1.0/float64(b.fftSize), 0)
	for k := 0; k < b.fftSize; k++ {
		x[k] *= scale
	}

	// Step 6: Multiply by chirp and extract result
	for k := 0; k < b.length; k++ {
		buffer[k] = x[k] * b.chirp[k]
	}
}
