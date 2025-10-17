// Package algorithm contains the individual FFT algorithm implementations
package algorithm

import (
	"math"
)

// Direction represents whether an FFT is forward or inverse
type Direction int

const (
	// Forward represents a forward FFT
	Forward Direction = iota
	// Inverse represents an inverse FFT
	Inverse
)

// Dft implements a naive O(n^2) Discrete Fourier Transform
// This is primarily used for testing and as a fallback for very small sizes
type Dft struct {
	twiddles  []complex128
	direction Direction
}

// NewDft creates a new DFT instance for the given size and direction
func NewDft(length int, direction Direction) *Dft {
	twiddles := computeTwiddles(length, direction)
	return &Dft{
		twiddles:  twiddles,
		direction: direction,
	}
}

// Len returns the FFT size
func (d *Dft) Len() int {
	return len(d.twiddles)
}

// Direction returns the FFT direction
func (d *Dft) Direction() Direction {
	return d.direction
}

// computeTwiddles precomputes all twiddle factors for a given FFT size
func computeTwiddles(n int, direction Direction) []complex128 {
	twiddles := make([]complex128, n)
	for k := 0; k < n; k++ {
		angle := 2.0 * math.Pi * float64(k) / float64(n)
		if direction == Forward {
			angle = -angle
		}
		twiddles[k] = complex(math.Cos(angle), math.Sin(angle))
	}
	return twiddles
}

// InplaceScratchLen returns the scratch space needed for in-place processing
func (d *Dft) InplaceScratchLen() int {
	return d.Len()
}

// OutOfPlaceScratchLen returns the scratch space needed for out-of-place processing
func (d *Dft) OutOfPlaceScratchLen() int {
	return 0
}

// ImmutableScratchLen returns the scratch space needed for immutable processing
func (d *Dft) ImmutableScratchLen() int {
	return 0
}

// Process computes the FFT in-place
func (d *Dft) Process(buffer []complex128) {
	scratch := make([]complex128, d.InplaceScratchLen())
	d.ProcessWithScratch(buffer, scratch)
}

// ProcessWithScratch computes the FFT in-place using provided scratch space
func (d *Dft) ProcessWithScratch(buffer, scratch []complex128) {
	// For in-place operation, we use scratch as temporary output space
	for i := 0; i < len(buffer); i += d.Len() {
		chunk := buffer[i : i+d.Len()]
		selfScratch := scratch[:d.Len()]
		d.performFftOutOfPlace(chunk, selfScratch, nil)
		copy(chunk, selfScratch)
	}
}

// ProcessOutOfPlace computes the FFT from input to output
func (d *Dft) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += d.Len() {
		inChunk := input[i : i+d.Len()]
		outChunk := output[i : i+d.Len()]
		d.performFftOutOfPlace(inChunk, outChunk, scratch)
	}
}

// ProcessImmutable computes the FFT without modifying the input
func (d *Dft) ProcessImmutable(input []complex128, output, scratch []complex128) {
	for i := 0; i < len(input); i += d.Len() {
		inChunk := input[i : i+d.Len()]
		outChunk := output[i : i+d.Len()]
		d.performFftImmutable(inChunk, outChunk, scratch)
	}
}

// performFftImmutable performs the DFT computation
func (d *Dft) performFftImmutable(signal []complex128, spectrum []complex128, scratch []complex128) {
	n := len(d.twiddles)

	for k := 0; k < n; k++ {
		sum := complex(0, 0)
		twiddleIndex := 0

		for _, inputVal := range signal {
			twiddle := d.twiddles[twiddleIndex]
			sum += twiddle * inputVal

			twiddleIndex += k
			if twiddleIndex >= n {
				twiddleIndex -= n
			}
		}

		spectrum[k] = sum
	}
}

// performFftOutOfPlace is the same as performFftImmutable for DFT
func (d *Dft) performFftOutOfPlace(signal []complex128, spectrum []complex128, scratch []complex128) {
	d.performFftImmutable(signal, spectrum, scratch)
}

// computeTwiddles32 precomputes all twiddle factors for a given FFT size (complex64)
func computeTwiddles32(n int, direction Direction) []complex64 {
	twiddles := make([]complex64, n)
	for k := 0; k < n; k++ {
		angle := 2.0 * math.Pi * float64(k) / float64(n)
		if direction == Forward {
			angle = -angle
		}
		twiddles[k] = complex(float32(math.Cos(angle)), float32(math.Sin(angle)))
	}
	return twiddles
}

// Dft32 implements DFT for complex64
type Dft32 struct {
	twiddles  []complex64
	direction Direction
}

// NewDft32 creates a new DFT instance for complex64
func NewDft32(length int, direction Direction) *Dft32 {
	twiddles := computeTwiddles32(length, direction)
	return &Dft32{
		twiddles:  twiddles,
		direction: direction,
	}
}

// Len returns the FFT size
func (d *Dft32) Len() int {
	return len(d.twiddles)
}

// Direction returns the FFT direction
func (d *Dft32) Direction() Direction {
	return d.direction
}

// InplaceScratchLen returns the scratch space needed for in-place processing
func (d *Dft32) InplaceScratchLen() int {
	return d.Len()
}

// OutOfPlaceScratchLen returns the scratch space needed for out-of-place processing
func (d *Dft32) OutOfPlaceScratchLen() int {
	return 0
}

// ImmutableScratchLen returns the scratch space needed for immutable processing
func (d *Dft32) ImmutableScratchLen() int {
	return 0
}

// Process computes the FFT in-place
func (d *Dft32) Process(buffer []complex64) {
	scratch := make([]complex64, d.InplaceScratchLen())
	d.ProcessWithScratch(buffer, scratch)
}

// ProcessWithScratch computes the FFT in-place using provided scratch space
func (d *Dft32) ProcessWithScratch(buffer, scratch []complex64) {
	for i := 0; i < len(buffer); i += d.Len() {
		chunk := buffer[i : i+d.Len()]
		selfScratch := scratch[:d.Len()]
		d.performFftOutOfPlace(chunk, selfScratch, nil)
		copy(chunk, selfScratch)
	}
}

// ProcessOutOfPlace computes the FFT from input to output
func (d *Dft32) ProcessOutOfPlace(input, output, scratch []complex64) {
	for i := 0; i < len(input); i += d.Len() {
		inChunk := input[i : i+d.Len()]
		outChunk := output[i : i+d.Len()]
		d.performFftOutOfPlace(inChunk, outChunk, scratch)
	}
}

// ProcessImmutable computes the FFT without modifying the input
func (d *Dft32) ProcessImmutable(input []complex64, output, scratch []complex64) {
	for i := 0; i < len(input); i += d.Len() {
		inChunk := input[i : i+d.Len()]
		outChunk := output[i : i+d.Len()]
		d.performFftImmutable(inChunk, outChunk, scratch)
	}
}

// performFftImmutable performs the DFT computation for complex64
func (d *Dft32) performFftImmutable(signal []complex64, spectrum []complex64, scratch []complex64) {
	n := len(d.twiddles)

	for k := 0; k < n; k++ {
		sum := complex(float32(0), float32(0))
		twiddleIndex := 0

		for _, inputVal := range signal {
			twiddle := d.twiddles[twiddleIndex]
			sum += twiddle * inputVal

			twiddleIndex += k
			if twiddleIndex >= n {
				twiddleIndex -= n
			}
		}

		spectrum[k] = sum
	}
}

// performFftOutOfPlace is the same as performFftImmutable for DFT
func (d *Dft32) performFftOutOfPlace(signal []complex64, spectrum []complex64, scratch []complex64) {
	d.performFftImmutable(signal, spectrum, scratch)
}
