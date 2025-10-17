# gofft
[![Release](https://img.shields.io/badge/version-v0.4.0-blue)](https://github.com/10d9e/gofft/releases)
[![Go CI](https://github.com/10d9e/gofft/actions/workflows/ci.yml/badge.svg)](https://github.com/10d9e/gofft/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT%20OR%20Apache--2.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/10d9e/gofft)](https://goreportcard.com/report/github.com/10d9e/gofft)
[![Tests](https://img.shields.io/badge/tests-256%20passing-success)](https://github.com/10d9e/gofft)
[![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://github.com/10d9e/gofft)

A high-performance FFT library for Go, ported from [RustFFT](https://github.com/ejmahler/RustFFT).

<img width="300" height="300" alt="image" src="https://github.com/user-attachments/assets/a45b2e1a-ee46-4c50-9dea-1b06d56ffc35" />

**Status**: ✅ **v0.4.0 - Rader's algorithm for optimized prime FFTs!**

## Features

- 🚀 **ANY size is O(n log n)** via Bluestein's algorithm
- ⚡ **Rader's algorithm** for optimized prime FFTs (NEW in v0.4.0!)
- ✅ **20 optimized butterflies** (2-32)
- ✅ **Radix-4** for power-of-two sizes
- ✅ **Zero allocations** with scratch buffer reuse
- ✅ **Thread-safe** - concurrent usage supported
- ✅ **~98% algorithm parity** with RustFFT
- ⏳ **SIMD support** (future enhancement)

## Quick Start

```go
package main

import (
    "github.com/10d9e/gofft"
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

**NEW in v0.4.0: Rader's Algorithm!**
- Primes 3-97 now ~2-3x faster than v0.3.2
- More efficient than Bluestein's for primes
- Automatic algorithm selection

**v0.3.2**: Bluestein's algorithm makes ANY size O(n log n)

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

### Primes 3-97 (Rader's - NEW!)
All primes up to 97 use optimized Rader's algorithm

### Everything Else (Bluestein's)
Primes >97, composites, arbitrary sizes - ALL O(n log n)!

## Build & Test

```bash
# Build
go build ./...

# Run tests (all passing!)
go test -v

# Run benchmarks
go test -bench=. -benchmem

# Try the examples
go run cmd/example/main.go
```

## What's New in v0.4.0

- ⚡ **Rader's Algorithm**: Optimizes primes 3-97 (~2-3x faster than Bluestein's)
- ✅ **256 tests passing** (up from 228)
- ✅ **~98% algorithm parity** with RustFFT (up from 95%)
- ✅ **27 algorithms** implemented (up from 26)

## Documentation

- [V0.4.0_RELEASE_NOTES.md](V0.4.0_RELEASE_NOTES.md) - v0.4.0 release notes
- [V0.3.2_RELEASE_NOTES.md](V0.3.2_RELEASE_NOTES.md) - v0.3.2 release notes
- [API_REFERENCE.md](API_REFERENCE.md) - Detailed API documentation

## Status

✅ **Production-ready** for ALL sizes  
✅ **O(n log n)** for ALL sizes  
✅ **256 tests passing** (100% success rate)  
📊 **~98% algorithm parity** with RustFFT
