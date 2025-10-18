package simd

import (
	"runtime"
	"unsafe"
)

// SIMDLevel represents the available SIMD instruction set
type SIMDLevel int

const (
	Scalar SIMDLevel = iota // No SIMD support
	SSE                     // x86_64 SSE4.1
	AVX                     // x86_64 AVX/FMA
	NEON                    // ARM64 NEON
)

// String returns a human-readable representation of the SIMD level
func (s SIMDLevel) String() string {
	switch s {
	case Scalar:
		return "Scalar"
	case SSE:
		return "SSE4.1"
	case AVX:
		return "AVX/FMA"
	case NEON:
		return "NEON"
	default:
		return "Unknown"
	}
}

// DetectSIMD detects the best available SIMD instruction set
func DetectSIMD() SIMDLevel {
	switch runtime.GOARCH {
	case "amd64":
		return detectX86SIMD()
	case "arm64":
		return detectARMSIMD()
	default:
		return Scalar
	}
}

// detectX86SIMD detects x86_64 SIMD capabilities
func detectX86SIMD() SIMDLevel {
	// Check for AVX support
	if hasAVX() {
		return AVX
	}

	// Check for SSE4.1 support
	if hasSSE41() {
		return SSE
	}

	return Scalar
}

// detectARMSIMD detects ARM64 SIMD capabilities
func detectARMSIMD() SIMDLevel {
	// ARM64 always has NEON support
	return NEON
}

// hasAVX checks if the CPU supports AVX and FMA instructions
func hasAVX() bool {
	// This is a simplified check - in practice, you'd use CPUID
	// For now, we'll assume AVX is available on modern x86_64 systems
	// TODO: Implement proper CPUID check
	return false // Disabled for now, will implement proper detection
}

// hasSSE41 checks if the CPU supports SSE4.1 instructions
func hasSSE41() bool {
	// This is a simplified check - in practice, you'd use CPUID
	// For now, we'll assume SSE4.1 is available on modern x86_64 systems
	// TODO: Implement proper CPUID check
	return false // Disabled for now, will implement proper detection
}

// IsAvailable returns true if SIMD is available for the current architecture
func IsAvailable() bool {
	return DetectSIMD() != Scalar
}

// GetSIMDLevel returns the current SIMD level
func GetSIMDLevel() SIMDLevel {
	return DetectSIMD()
}

// PlatformInfo returns information about the current platform
func PlatformInfo() (arch, os string, simd SIMDLevel) {
	return runtime.GOARCH, runtime.GOOS, DetectSIMD()
}

// VectorSize returns the number of elements that can be processed in parallel
func VectorSize(simd SIMDLevel) int {
	switch simd {
	case SSE:
		return 2 // SSE can process 2 complex64 or 1 complex128
	case AVX:
		return 4 // AVX can process 4 complex64 or 2 complex128
	case NEON:
		return 2 // NEON can process 2 complex64 or 1 complex128
	default:
		return 1 // Scalar processing
	}
}

// Alignment returns the required memory alignment for SIMD operations
func Alignment(simd SIMDLevel) int {
	switch simd {
	case SSE:
		return 16 // 128-bit alignment
	case AVX:
		return 32 // 256-bit alignment
	case NEON:
		return 16 // 128-bit alignment
	default:
		return 8 // 64-bit alignment (complex128)
	}
}

// IsAligned checks if a slice is properly aligned for SIMD operations
func IsAligned(data []complex128, simd SIMDLevel) bool {
	if len(data) == 0 {
		return true
	}

	ptr := uintptr(unsafe.Pointer(&data[0]))
	alignment := Alignment(simd)
	return ptr%uintptr(alignment) == 0
}

// AlignSlice ensures a slice is properly aligned for SIMD operations
func AlignSlice(data []complex128, simd SIMDLevel) []complex128 {
	if IsAligned(data, simd) {
		return data
	}

	// Create a new aligned slice
	aligned := make([]complex128, len(data))
	copy(aligned, data)
	return aligned
}
