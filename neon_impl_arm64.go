//go:build arm64

package gofft

import (
	"github.com/10d9e/gofft/algorithm"
	"github.com/10d9e/gofft/simd/neon"
)

// processNEONARM64 implements NEON processing for ARM64 architecture
func (a *neonButterflyAdapter) processNEONARM64(buffer []complex128) {
	// Use actual NEON implementations with proper direction support
	// Convert our Direction to algorithm.Direction
	algoDir := algorithm.Forward
	if a.direction == Inverse {
		algoDir = algorithm.Inverse
	}

	// Use actual NEON implementations based on size
	switch a.length {
	// Butterflies
	case 1:
		neon.Butterfly1_NEON(buffer, algoDir)
	case 2:
		neon.Butterfly2_NEON(buffer, algoDir)
	case 3:
		neon.Butterfly3_NEON(buffer, algoDir)
	case 4:
		neon.Butterfly4_NEON(buffer, algoDir)
	case 5:
		neon.Butterfly5_NEON(buffer, algoDir)
	case 6:
		neon.Butterfly6_NEON(buffer, algoDir)
	case 7:
		neon.Butterfly7_NEON(buffer, algoDir)
	case 8:
		neon.Butterfly8_NEON(buffer, algoDir)
	case 9:
		neon.Butterfly9_NEON(buffer, algoDir)
	case 10:
		neon.Butterfly10_NEON(buffer, algoDir)
	case 11:
		neon.Butterfly11_NEON(buffer, algoDir)
	case 12:
		neon.Butterfly12_NEON(buffer, algoDir)
	case 13:
		neon.Butterfly13_NEON(buffer, algoDir)
	case 15:
		neon.Butterfly15_NEON(buffer, algoDir)
	case 16:
		neon.Butterfly16_NEON(buffer, algoDir)
	case 17:
		neon.Butterfly17_NEON(buffer, algoDir)
	case 19:
		neon.Butterfly19_NEON(buffer, algoDir)
	case 23:
		neon.Butterfly23_NEON(buffer, algoDir)
	case 24:
		neon.Butterfly24_NEON(buffer, algoDir)
	case 27:
		neon.Butterfly27_NEON(buffer, algoDir)
	case 29:
		neon.Butterfly29_NEON(buffer, algoDir)
	case 31:
		neon.Butterfly31_NEON(buffer, algoDir)
	case 32:
		neon.Butterfly32_NEON(buffer, algoDir)
	// Radix4 sizes
	case 64:
		neon.Radix4_64_NEON(buffer, algoDir)
	case 128:
		neon.Radix4_128_NEON(buffer, algoDir)
	case 256:
		neon.Radix4_256_NEON(buffer, algoDir)
	case 512:
		neon.Radix4_512_NEON(buffer, algoDir)
	case 1024:
		neon.Radix4_1024_NEON(buffer, algoDir)
	case 2048:
		neon.Radix4_2048_NEON(buffer, algoDir)
	case 4096:
		neon.Radix4_4096_NEON(buffer, algoDir)
	case 8192:
		neon.Radix4_8192_NEON(buffer, algoDir)
	case 16384:
		neon.Radix4_16384_NEON(buffer, algoDir)
	case 32768:
		neon.Radix4_32768_NEON(buffer, algoDir)
	case 65536:
		neon.Radix4_65536_NEON(buffer, algoDir)
	default:
		// Fall back to scalar for unsupported sizes
		a.processNEON(buffer)
	}
}
