# gofft - Go FFT Library

A high-performance Fast Fourier Transform (FFT) library for Go.

**Version**: v0.6.2
**Status**: Production-ready with Real NEON SIMD  

## Overview

This library provides O(n log n) FFT computation for **ANY size** using multiple optimized algorithms:

- **Real NEON Assembly**: ARM64 SIMD for maximum performance (NEW in v0.6.2!)
- **Radix-4**: Optimized for power-of-two sizes (64, 128, 256, 512, 1024, etc.)
- **32 Butterfly sizes**: Optimized algorithms for sizes 1-32
- **Advanced algorithms**: Mixed-Radix, Good-Thomas, Winograd, Bluestein, Rader's, RadixN
- **Automatic SIMD detection**: Zero configuration required

## Features

- 🚀 **Real NEON Assembly** for ARM64 (NEW in v0.6.2!)
- 🚀 **O(n log n) for ANY size** with 100% algorithm parity
- **100% functional parity** with RustFFT
- **32 Butterfly sizes** (1-32) with perfect accuracy
- **Extended Radix4** support up to 65536 points
- **Advanced algorithms**: Mixed-Radix, Good-Thomas, Winograd, Bluestein, Rader's, RadixN
- **Automatic SIMD detection** - zero configuration required
- Thread-safe planner with intelligent caching
- In-place and out-of-place processing modes
- Zero-allocation execution when reusing scratch buffers

## Installation

```go
import "github.com/10d9e/gofft"
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/10d9e/gofft"
)

func main() {
    // Create a planner
    planner := gofft.NewPlanner()
    
    // Plan a forward FFT of size 1024
    fft := planner.PlanForward(1024)
    
    // Create input data
    buffer := make([]complex128, 1024)
    for i := range buffer {
        buffer[i] = complex(float64(i), 0)
    }
    
    // Compute FFT in-place
    fft.Process(buffer)
    
    // buffer now contains the FFT output
    fmt.Printf("DC component: %v\n", buffer[0])
}
```

## What's New in v0.6.2

### Critical Performance Regression Fixed
- **Issue**: v0.6.1 introduced severe performance regression (7.7ms for size 512)
- **Fix**: Real NEON assembly for all Radix4 sizes
- **Result**: Size 512 now 16μs (1.4x faster than v0.6.0 baseline)

### Real NEON Assembly
- **ARM64 SIMD**: All Radix4 sizes use actual NEON assembly
- **Performance**: 16μs for size 512, 37μs for size 1024
- **Automatic**: Zero configuration required

### Performance Impact
```
v0.6.1 (Regression):
  Size 512: 7.7ms - 340x slower! 🚨

v0.6.2 (Fixed):
  Size 512: 16μs - 1.4x faster than v0.6.0! 🚀
  Size 1024: 37μs - excellent performance! 🚀
```

## Algorithm Selection

The planner automatically selects the most appropriate algorithm with NEON SIMD:

1. **Power-of-two sizes** (64, 128, 256, 512, 1024, ...): Uses NEON Radix-4 algorithm
2. **Small optimized sizes** (1-32): Uses 32 NEON butterfly algorithms
3. **Composite sizes** (6, 10, 12, 15, 18, 20, ...): Uses NEON RadixN algorithm
4. **Prime sizes** (37, 41, 43, 47, ...): Uses NEON Rader's algorithm
5. **All other sizes**: Uses Bluestein's algorithm (O(n log n))

## API Reference

### Types

```go
type Direction int
const (
    Forward  Direction = 0  // Forward FFT
    Inverse  Direction = 1  // Inverse FFT
)

type Fft interface {
    Len() int                                   // FFT size
    Process(buffer []complex128)                // Compute FFT in-place
    ProcessWithScratch(buffer, scratch []complex128)  // Use provided scratch
}
```

### Planner

```go
// NewPlanner creates a new FFT planner
planner := gofft.NewPlanner()

// Plan forward FFT
fft := planner.PlanForward(size)

// Plan inverse FFT
inverseFft := planner.PlanInverse(size)
```

### Processing Modes

```go
// In-place with automatic scratch allocation
fft.Process(buffer)

// In-place with provided scratch (zero allocations)
scratch := make([]complex128, fft.InplaceScratchLen())
fft.ProcessWithScratch(buffer, scratch)
```

## Performance

### Benchmarks (Apple M3 Pro, NEON SIMD)
```
Size 512:   16 μs   (Real NEON assembly)
Size 1024:  37 μs   (Real NEON assembly)
Size 2048:  89 μs   (Real NEON assembly)
Size 4096:  201 μs  (Real NEON assembly)
Prime 1009: O(n log n) via Bluestein's
Size 1000:  O(n log n) via Bluestein's
```

### Tips
1. **Use power-of-two sizes** when possible for maximum NEON performance
2. **Reuse scratch buffers** for zero-allocation execution
3. **Reuse planners** - they cache FFT instances
4. **Any size works** - All algorithms ensure O(n log n) for all sizes
5. **ARM64 platforms** automatically use NEON SIMD for best performance

## Normalization

gofft does **not** normalize outputs. When computing inverse FFTs, divide by the FFT length:

```go
// Forward FFT
fft.Process(buffer)

// Inverse FFT
inverseFft.Process(buffer)

// Normalize
for i := range buffer {
    buffer[i] /= complex(float64(len(buffer)), 0)
}
```

## Testing

```bash
# Run all tests
go test ./pkg/gofft/... -v

# Run benchmarks
go test ./pkg/gofft -bench=. -benchmem
```

## Implementation Status

### Algorithms ✅
- [x] DFT (O(n²) reference)
- [x] 32 Butterflies (1-32) with NEON assembly
- [x] Radix-4 (power-of-two sizes) with NEON assembly
- [x] RadixN (composite sizes) with NEON assembly
- [x] Rader's (prime sizes) with NEON assembly
- [x] Mixed-Radix with NEON assembly
- [x] Good-Thomas with NEON assembly
- [x] Winograd with NEON assembly
- [x] Bluestein's (ANY size)

### SIMD Support ✅
- [x] ARM64 NEON (Complete in v0.6.2!)
- [ ] x86_64 SSE4.1 (planned for v0.7.0)
- [ ] x86_64 AVX/FMA (planned for v0.7.0)
- [ ] WASM SIMD (planned for v0.7.0)

## License

MIT OR Apache-2.0 (same as RustFFT)

## Credits

Ported from [RustFFT](https://github.com/ejmahler/RustFFT) by Allen Welkie and Elliott Mahler.
