//go:build !arm64

package gofft

// processNEONARM64 is a stub for non-ARM64 architectures
func (a *neonButterflyAdapter) processNEONARM64(buffer []complex128) {
	// Fall back to scalar implementation on non-ARM64
	a.processNEON(buffer)
}
