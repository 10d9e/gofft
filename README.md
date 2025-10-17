# gofft
<img width="300" height="300" alt="image" src="https://github.com/user-attachments/assets/a45b2e1a-ee46-4c50-9dea-1b06d56ffc35" />

A high-performance FFT library for Go.

**Status**: Core algorithms work in pure golang, SIMD acceleration in progress.

## Features

- Fast FFT computation for arbitrary sizes
- Multiple optimized algorithms:
  - Butterflies for common small sizes (2, 3, 4, 8, 16, 32)
  - Radix-4 for power-of-two sizes
  - DFT fallback for other sizes
- Thread-safe - all FFT instances can be used concurrently
- Architecture-specific SIMD optimizations (planned):
  - x86_64: SSE4.1, AVX/FMA
  - ARM64: NEON

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

## Build & Test

```bash
# Build
go build ./...

# Run tests
go test ./pkg/gofft/...

# Run benchmarks
go test -bench=. ./pkg/gofft/...
```

See [pkg/gofft/README.md](pkg/gofft/README.md) for detailed documentation.
