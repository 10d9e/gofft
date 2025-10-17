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

