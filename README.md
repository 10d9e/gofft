# gofft
<img width="300" height="300" alt="image" src="https://github.com/user-attachments/assets/a45b2e1a-ee46-4c50-9dea-1b06d56ffc35" />

A high-performance FFT library for Go, ported from [RustFFT](https://github.com/ejmahler/RustFFT).

**Status**: ✅ **v0.3.0 - ALL sizes now O(n log n)!**

## Features

- 🚀 **ANY size is O(n log n)** via Bluestein's algorithm (NEW in v0.3.0!)
- ✅ **20 optimized butterflies** (2-32)
- ✅ **Radix-4** for power-of-two sizes
- ✅ **Zero allocations** with scratch buffer reuse
- ✅ **Thread-safe** - concurrent usage supported
- ✅ **~95% algorithm parity** with RustFFT
- ⏳ **SIMD support** (future enhancement)

## Quick Start

```go
package main

import (
    "github.com/10d9e/gofft/pkg/gofft"
)

func main() {
    // Create a planner
    planner := gofft.NewPlanner()
    
    // Plan a forward FFT of size 1024
    fft := planner.PlanForward(1024)
    
    // Create and process data
    buffer := make([]complex128, 1024)
    // ... fill buffer with data ...
    fft.Process(buffer)
}
```

## Highlights

**NEW: Bluestein's Algorithm** makes ANY size O(n log n)!
- Prime 1009: ~100x faster than v0.2.0
- Size 1000: ~100x faster than v0.2.0
- Works for ALL sizes automatically

## Performance

Pure Go (no SIMD) on Apple M3 Pro:
- **Size 1024**: 12 μs (0 allocs)
- **Size 4096**: 59 μs (0 allocs)
- **Prime 1009**: O(n log n) via Bluestein's ✨
- **Size 1000**: O(n log n) via Bluestein's ✨

## Algorithm Coverage

### Power-of-Two (Radix-4)
2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, ...

### Small Sizes (Butterflies)
2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 16, 17, 19, 23, 24, 27, 29, 31, 32

### Everything Else (Bluestein's)
**ALL other sizes** - primes, composites, arbitrary! ✅

## Build & Test

```bash
# Build
go build ./...

# Run tests (all passing!)
go test ./pkg/gofft/... -v

# Run benchmarks
go test ./pkg/gofft -bench=. -benchmem

# Try the examples
go run cmd/example/main.go
```

## What's New in v0.3.0

- 🚀 **Bluestein's Algorithm**: Makes ANY size O(n log n)
- ✅ **228 tests passing** (up from 224)
- ✅ **~95% algorithm parity** with RustFFT
- 🎯 **~100x speedup** for non-power-of-two sizes

## Documentation

- [V2_RELEASE_NOTES.md](V0.3.0_RELEASE_NOTES.md) - v0.3.0 release notes
- [pkg/gofft/README.md](pkg/gofft/README.md) - API documentation

## Status

✅ **Production-ready** for ALL sizes  
✅ **O(n log n)** for ALL sizes  
✅ **228 tests passing** (100% success rate)  
📊 **~95% algorithm parity** with RustFFT
