# gofft

A high-performance FFT library for Go, ported from RustFFT with architecture-specific SIMD optimizations.

## Features

- **Fast FFT computation** for arbitrary sizes
- **SIMD acceleration** for x86_64 (AVX, SSE4.1) and ARM64 (NEON) *[in progress]*
- **Multiple algorithms**:
  - DFT (O(n²)) for very small sizes
  - Butterflies for common small sizes (2, 3, 4, 8, 16, 32)
  - Radix-4 for power-of-two sizes
  - Mixed-radix algorithms for composite sizes *[planned]*
  - Rader's and Bluestein's algorithms for prime sizes *[planned]*
- **Thread-safe** - all FFT instances can be used concurrently
- **Memory efficient** - planners reuse internal data across FFT instances

## Installation

```go
import "github.com/example/gofft/pkg/gofft"
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/example/gofft/pkg/gofft"
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

### Computing Inverse FFT

```go
// Plan an inverse FFT
inverseFft := planner.PlanInverse(1024)

// Apply inverse FFT
inverseFft.Process(buffer)

// Don't forget to normalize! FFT outputs are not normalized
for i := range buffer {
    buffer[i] /= complex(float64(len(buffer)), 0)
}
```

### Reusing Scratch Space

For better performance when computing multiple FFTs, reuse scratch space:

```go
fft := planner.PlanForward(1024)
scratch := make([]complex128, fft.InplaceScratchLen())

// Compute multiple FFTs
for _, buffer := range buffers {
    fft.ProcessWithScratch(buffer, scratch)
}
```

## Performance Tips

1. **Use power-of-two sizes when possible** - These use the optimized Radix-4 algorithm
2. **Reuse the same planner** for multiple FFT instances to share internal data
3. **Pre-allocate scratch space** and reuse it across FFT computations
4. **Sizes of the form 2^n * 3^m** will be fastest once mixed-radix is implemented

## Implementation Status

### Algorithms
- [x] DFT (naive O(n²) algorithm)
- [x] Butterflies (2, 3, 4, 8, 16, 32)
- [x] Radix-4 (power-of-two sizes)
- [ ] RadixN (composite sizes with small factors)
- [ ] Mixed-Radix (general composite sizes)
- [ ] Good-Thomas Algorithm (coprime factorization)
- [ ] Rader's Algorithm (prime sizes)
- [ ] Bluestein's Algorithm (large prime sizes)

### SIMD Support
- [ ] x86_64 SSE4.1
- [ ] x86_64 AVX/FMA
- [ ] ARM64 NEON

### Testing
- [ ] Comprehensive unit tests
- [ ] Accuracy tests
- [ ] Benchmarks

## Normalization

gofft does **not** normalize outputs. When computing a forward FFT followed by an inverse FFT,
you need to divide by the FFT length to recover the original input:

```go
// Forward FFT
fft.Process(buffer)

// Inverse FFT
inverseFft.Process(buffer)

// Normalize
n := float64(len(buffer))
for i := range buffer {
    buffer[i] /= complex(n, 0)
}
```

## License

Same as RustFFT: MIT OR Apache-2.0

## Credits

Ported from [RustFFT](https://github.com/ejmahler/RustFFT) by Allen Welkie and Elliott Mahler.

