//go:build arm64

package gofft

// processNEONARM64 implements NEON processing for ARM64 architecture
func (a *neonButterflyAdapter) processNEONARM64(buffer []complex128) {
	// For now, use scalar implementation with proper direction handling
	// TODO: Implement proper NEON assembly with direction support
	// The NEON functions are hardcoded to forward FFT, but we need to support both directions
	a.processNEON(buffer)
}
