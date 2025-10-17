//go:build arm64

#include "textflag.h"

DATA ·neonSign<>+0x00(SB)/8, $1.0
DATA ·neonSign<>+0x08(SB)/8, $-1.0
GLOBL ·neonSign<>(SB), RODATA, $16

DATA ·neonPlusI<>+0x00(SB)/8, $-1.0
DATA ·neonPlusI<>+0x08(SB)/8, $ 1.0
GLOBL ·neonPlusI<>(SB), RODATA, $16

DATA ·neonMinusI<>+0x00(SB)/8, $ 1.0
DATA ·neonMinusI<>+0x08(SB)/8, $-1.0
GLOBL ·neonMinusI<>(SB), RODATA, $16

// func cmuladdNEON(dst, a, b *complex128, n int)
// Stub: falls back to portable Go implementation
TEXT ·cmuladdNEON(SB), NOSPLIT, $0-32
	RET

// func cmuladdFMLA(dst, a, b *complex128, n int)  
// Stub: falls back to portable Go implementation
TEXT ·cmuladdFMLA(SB), NOSPLIT, $0-32
	RET

// func butterfly4TwiddledNEON(buf, w1, w2, w3 *complex128, k, m, n, q int, invert bool)
// Stub: falls back to portable Go implementation
TEXT ·butterfly4TwiddledNEON(SB), NOSPLIT, $0-56
	RET

// func butterfly8TwiddledNEON(buf, w1, w2, w3, w4, w5, w6, w7 *complex128, k, m, n, q int, invert bool)
// Stub: falls back to portable Go implementation
TEXT ·butterfly8TwiddledNEON(SB), NOSPLIT, $0-80
	RET
