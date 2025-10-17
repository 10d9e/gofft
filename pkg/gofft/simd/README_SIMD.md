# SIMD backends (AVX2 / NEON)

This folder contains runtime-dispatched SIMD entrypoints and portable fallbacks.

- `cmuladd` — fully working AVX2/NEON assembly (2 complex128 per iter)
- `butterfly4Twiddled` — portable unrolled-by-2j; packed-2j SIMD hooks (stubs to keep concise here)
- `butterfly8Twiddled` — portable; SIMD hooks in place

Use `golang.org/x/sys/cpu` to select features at runtime.
