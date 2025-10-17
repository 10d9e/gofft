# gofft (prototype SIMD-ready FFT scaffolding)

This repo contains a prototype of a RustFFT-like Go FFT library skeleton with:
- Portable radix-2/4/8 butterfly helpers
- Twiddle-aware portable kernels (radix-4 unrolled by 2j, radix-8)
- SIMD dispatch seams and **working AVX2/NEON assembly** for `cmuladd`
- Packed 2j SIMD radix-4 butterfly implementations (AVX2 and NEON)
- Runtime CPU feature detection via `golang.org/x/sys/cpu`

> Note: This is a scaffold so you can drop in highly optimized assembly for butterflies.
> The portable versions are complete and used as fallbacks. Some symbols are placeholders
> (e.g., radix-8 SIMD hooks) and can be filled later without changing public APIs.

## Build & test

```bash
# from repository root
go build ./...
# optional: add your tests under pkg/gofft
```

## Examples

Build and run the demo:

```bash
go run ./cmd/fft-demo -n 1024
go run ./cmd/fft-parallel -n 16384 -workers 0
```

## Portable vs SIMD toggle (for benchmarks)

You can force portable implementations at runtime for apples-to-apples comparisons:

```go
import "github.com/example/gofft/pkg/gofft/simd"

simd.ForcePortable(true)  // force portable
simd.ForcePortable(false) // leave current bindings (SIMD if available)
```

See `BenchmarkFFT_SIMD_*` vs `BenchmarkFFT_Portable_*`.
