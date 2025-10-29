# gofft
[![Release](https://img.shields.io/badge/version-v0.5.0-blue)](https://github.com/10d9e/gofft/releases)
[![Go CI](https://github.com/10d9e/gofft/actions/workflows/ci.yml/badge.svg)](https://github.com/10d9e/gofft/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT%20OR%20Apache--2.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/10d9e/gofft)](https://goreportcard.com/report/github.com/10d9e/gofft)
[![Tests](https://img.shields.io/badge/tests-256%20passing-success)](https://github.com/10d9e/gofft)
[![DAITU](https://img.shields.io/badge/AI-DAITU%201.0-blue.svg)](DAITU)

<img width="300" height="300" alt="image" src="https://github.com/user-attachments/assets/a45b2e1a-ee46-4c50-9dea-1b06d56ffc35" />

A high-performance FFT library for Go, ported from [RustFFT](https://github.com/ejmahler/RustFFT).

**Status**: **v0.5.0 - 100% Algorithm Parity Achieved!**

## Features

- **100% scalar algorithm parity** with RustFFT (NEW in v0.5.0!)
- **RadixN algorithm** for multi-factor composites (NEW in v0.5.0!)
- **Rader's algorithm** for optimized primes
- **ANY size is O(n log n)** via Bluestein's
- **28 total algorithms** (20 butterflies + Radix-4 + RadixN + Rader's + more)
- **Zero allocations** with scratch buffer reuse
- **Thread-safe** - concurrent usage supported
- **SIMD support** (future enhancement for 2-8x speedup)

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

**v0.5.0: 100% ALGORITHM PARITY!**
- **RadixN**: Optimizes multi-factor composites (60, 120, 84, ...)
- **Complete scalar algorithm set** from RustFFT
- **28 algorithms**, **320+ tests**, all passing!

**v0.4.0**: Rader's algorithm for primes  
**v0.3.2**: Bluestein's for ANY size

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

### Multi-Factor Composites (RadixN - NEW!)
6, 10, 12, 14, 15, 18, 20, 21, 24, 28, 30, 36, 40, 42, 48, 54, 56, 60, 72, 80, 84, 90, 96, 100, 120, ...

### Primes 3-97 (Rader's)
All primes from 3 to 97 use Rader's algorithm

### Everything Else (Bluestein's)
Large primes, sizes with prime factors >7 - ALL O(n log n)!

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

## What's New in v0.5.0

- **100% ALGORITHM PARITY** - Complete RustFFT scalar port!
- **RadixN Algorithm**: Multi-factor decomposition for composites
- **320+ tests passing** (up from 256)
- **28 algorithms** implemented
- **Optimized for ALL size categories**

## Documentation

- [V0.5.0_100PERCENT_COMPLETE.md](V0.5.0_100PERCENT_COMPLETE.md) - **100% parity milestone!** 🏆
- [V0.4.0_RELEASE_NOTES.md](V0.4.0_RELEASE_NOTES.md) - v0.4.0 release notes
- [API_REFERENCE.md](API_REFERENCE.md) - Detailed API documentation

## Status

**100% SCALAR ALGORITHM PARITY** with RustFFT!  
**O(n log n)** for ALL sizes  
**320+ tests passing** (100% success rate)  
**28 algorithms** - complete scalar algorithm set
