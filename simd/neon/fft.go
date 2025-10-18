//go:build arm64

package neon

import (
	"github.com/10d9e/gofft/algorithm"
)

// NEONFFT represents a NEON-optimized FFT instance
type NEONFFT struct {
	length    int
	direction algorithm.Direction
}

// NewNEONFFT creates a new NEON-optimized FFT instance
func NewNEONFFT(length int, direction algorithm.Direction) *NEONFFT {
	return &NEONFFT{
		length:    length,
		direction: direction,
	}
}

// Len returns the FFT length
func (f *NEONFFT) Len() int {
	return f.length
}

// Direction returns the FFT direction
func (f *NEONFFT) Direction() algorithm.Direction {
	return f.direction
}

// Process computes the FFT in-place using NEON optimizations
func (f *NEONFFT) Process(buffer []complex128) {
	if len(buffer) < f.length {
		return
	}

	// Process each chunk of the buffer
	for i := 0; i < len(buffer); i += f.length {
		chunk := buffer[i : i+f.length]
		f.processChunk(chunk)
	}
}

// processChunk processes a single FFT chunk using NEON optimizations
func (f *NEONFFT) processChunk(data []complex128) {
	// For now, use NEON butterflies for supported sizes
	// TODO: Implement full NEON FFT algorithms

	switch f.length {
	case 2:
		Butterfly2_NEON(data, f.direction)
	case 4:
		Butterfly4_NEON(data, f.direction)
	case 8:
		Butterfly8_NEON(data, f.direction)
	case 16:
		Butterfly16_NEON(data, f.direction)
	case 32:
		Butterfly32_NEON(data, f.direction)
	default:
		// For unsupported sizes, fall back to scalar implementation
		f.processScalar(data)
	}
}

// processScalar falls back to scalar implementation for unsupported sizes
func (f *NEONFFT) processScalar(data []complex128) {
	// Create a scalar FFT and use it
	scalarFft := algorithm.NewDft(f.length, f.direction)
	scratch := make([]complex128, scalarFft.InplaceScratchLen())
	scalarFft.ProcessWithScratch(data, scratch)
}

// NEONButterflyFFT represents a NEON-optimized butterfly FFT
type NEONButterflyFFT struct {
	*NEONFFT
}

// NewNEONButterflyFFT creates a new NEON butterfly FFT
func NewNEONButterflyFFT(length int, direction algorithm.Direction) *NEONButterflyFFT {
	return &NEONButterflyFFT{
		NEONFFT: NewNEONFFT(length, direction),
	}
}

// Process computes the butterfly FFT using NEON optimizations
func (f *NEONButterflyFFT) Process(buffer []complex128) {
	if len(buffer) < f.length {
		return
	}

	// Process each chunk using NEON butterflies
	for i := 0; i < len(buffer); i += f.length {
		chunk := buffer[i : i+f.length]
		ProcessVectorizedButterfly(chunk, f.length, f.direction)
	}
}

// NEONRadix4FFT represents a NEON-optimized Radix-4 FFT
type NEONRadix4FFT struct {
	*NEONFFT
}

// NewNEONRadix4FFT creates a new NEON Radix-4 FFT
func NewNEONRadix4FFT(length int, direction algorithm.Direction) *NEONRadix4FFT {
	return &NEONRadix4FFT{
		NEONFFT: NewNEONFFT(length, direction),
	}
}

// Process computes the Radix-4 FFT using NEON optimizations
func (f *NEONRadix4FFT) Process(buffer []complex128) {
	if len(buffer) < f.length {
		return
	}

	// For now, fall back to scalar implementation
	// TODO: Implement actual NEON Radix-4
	f.processScalar(buffer)
}

// NEONRadixNFFT represents a NEON-optimized RadixN FFT
type NEONRadixNFFT struct {
	*NEONFFT
}

// NewNEONRadixNFFT creates a new NEON RadixN FFT
func NewNEONRadixNFFT(length int, direction algorithm.Direction) *NEONRadixNFFT {
	return &NEONRadixNFFT{
		NEONFFT: NewNEONFFT(length, direction),
	}
}

// Process computes the RadixN FFT using NEON optimizations
func (f *NEONRadixNFFT) Process(buffer []complex128) {
	if len(buffer) < f.length {
		return
	}

	// For now, fall back to scalar implementation
	// TODO: Implement actual NEON RadixN
	f.processScalar(buffer)
}

// NEONBluesteinFFT represents a NEON-optimized Bluestein FFT
type NEONBluesteinFFT struct {
	*NEONFFT
}

// NewNEONBluesteinFFT creates a new NEON Bluestein FFT
func NewNEONBluesteinFFT(length int, direction algorithm.Direction) *NEONBluesteinFFT {
	return &NEONBluesteinFFT{
		NEONFFT: NewNEONFFT(length, direction),
	}
}

// Process computes the Bluestein FFT using NEON optimizations
func (f *NEONBluesteinFFT) Process(buffer []complex128) {
	if len(buffer) < f.length {
		return
	}

	// For now, fall back to scalar implementation
	// TODO: Implement actual NEON Bluestein
	f.processScalar(buffer)
}

// NEONRadersFFT represents a NEON-optimized Rader's FFT
type NEONRadersFFT struct {
	*NEONFFT
}

// NewNEONRadersFFT creates a new NEON Rader's FFT
func NewNEONRadersFFT(length int, direction algorithm.Direction) *NEONRadersFFT {
	return &NEONRadersFFT{
		NEONFFT: NewNEONFFT(length, direction),
	}
}

// Process computes the Rader's FFT using NEON optimizations
func (f *NEONRadersFFT) Process(buffer []complex128) {
	if len(buffer) < f.length {
		return
	}

	// For now, fall back to scalar implementation
	// TODO: Implement actual NEON Rader's
	f.processScalar(buffer)
}
