package algorithm

import (
	"math"
)

// Butterfly2 implements a size-2 FFT (Cooley-Tukey butterfly)
type Butterfly2 struct {
	direction Direction
}

// NewButterfly2 creates a new Butterfly2 instance
func NewButterfly2(direction Direction) *Butterfly2 {
	return &Butterfly2{direction: direction}
}

func (b *Butterfly2) Len() int                  { return 2 }
func (b *Butterfly2) Direction() Direction      { return b.direction }
func (b *Butterfly2) InplaceScratchLen() int    { return 0 }
func (b *Butterfly2) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly2) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly2) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly2) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 2 {
		b.performFft(buffer[i : i+2])
	}
}

func (b *Butterfly2) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 2 {
		b.performFftOutOfPlace(input[i:i+2], output[i:i+2])
	}
}

func (b *Butterfly2) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly2) performFft(buffer []complex128) {
	temp := buffer[0] + buffer[1]
	buffer[1] = buffer[0] - buffer[1]
	buffer[0] = temp
}

func (b *Butterfly2) performFftOutOfPlace(input, output []complex128) {
	output[0] = input[0] + input[1]
	output[1] = input[0] - input[1]
}

// Butterfly3 implements a size-3 FFT
type Butterfly3 struct {
	twiddle   complex128
	direction Direction
}

// NewButterfly3 creates a new Butterfly3 instance
func NewButterfly3(direction Direction) *Butterfly3 {
	twiddle := twiddleFactor(1, 3, direction)
	return &Butterfly3{
		twiddle:   twiddle,
		direction: direction,
	}
}

func (b *Butterfly3) Len() int                  { return 3 }
func (b *Butterfly3) Direction() Direction      { return b.direction }
func (b *Butterfly3) InplaceScratchLen() int    { return 0 }
func (b *Butterfly3) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly3) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly3) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly3) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 3 {
		b.performFft(buffer[i : i+3])
	}
}

func (b *Butterfly3) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 3 {
		b.performFftOutOfPlace(input[i:i+3], output[i:i+3])
	}
}

func (b *Butterfly3) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly3) performFft(buffer []complex128) {
	xp := buffer[1] + buffer[2]
	xn := buffer[1] - buffer[2]
	sum := buffer[0] + xp

	tempA := buffer[0] + complex(real(b.twiddle)*real(xp), real(b.twiddle)*imag(xp))
	tempB := complex(-imag(b.twiddle)*imag(xn), imag(b.twiddle)*real(xn))

	buffer[0] = sum
	buffer[1] = tempA + tempB
	buffer[2] = tempA - tempB
}

func (b *Butterfly3) performFftOutOfPlace(input, output []complex128) {
	xp := input[1] + input[2]
	xn := input[1] - input[2]
	sum := input[0] + xp

	tempA := input[0] + complex(real(b.twiddle)*real(xp), real(b.twiddle)*imag(xp))
	tempB := complex(-imag(b.twiddle)*imag(xn), imag(b.twiddle)*real(xn))

	output[0] = sum
	output[1] = tempA + tempB
	output[2] = tempA - tempB
}

// Butterfly4 implements a size-4 FFT
type Butterfly4 struct {
	direction Direction
}

// NewButterfly4 creates a new Butterfly4 instance
func NewButterfly4(direction Direction) *Butterfly4 {
	return &Butterfly4{direction: direction}
}

func (b *Butterfly4) Len() int                  { return 4 }
func (b *Butterfly4) Direction() Direction      { return b.direction }
func (b *Butterfly4) InplaceScratchLen() int    { return 0 }
func (b *Butterfly4) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly4) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly4) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly4) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 4 {
		b.performFft(buffer[i : i+4])
	}
}

func (b *Butterfly4) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 4 {
		b.performFftOutOfPlace(input[i:i+4], output[i:i+4])
	}
}

func (b *Butterfly4) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly4) performFft(buffer []complex128) {
	// Implementation using radix-2 decomposition
	// Column FFTs
	temp0 := buffer[0] + buffer[2]
	buffer[2] = buffer[0] - buffer[2]
	buffer[0] = temp0

	temp1 := buffer[1] + buffer[3]
	buffer[3] = buffer[1] - buffer[3]
	buffer[1] = temp1

	// Apply twiddle factor (rotate by 90 degrees)
	buffer[3] = rotate90(buffer[3], b.direction)

	// Row FFTs
	temp0 = buffer[0] + buffer[1]
	buffer[1] = buffer[0] - buffer[1]
	buffer[0] = temp0

	temp2 := buffer[2] + buffer[3]
	buffer[3] = buffer[2] - buffer[3]
	buffer[2] = temp2

	// Final transpose (swap indices 1 and 2)
	buffer[1], buffer[2] = buffer[2], buffer[1]
}

func (b *Butterfly4) performFftOutOfPlace(input, output []complex128) {
	// Column FFTs
	val0 := input[0] + input[2]
	val2 := input[0] - input[2]
	val1 := input[1] + input[3]
	val3 := input[1] - input[3]

	// Apply twiddle factor
	val3 = rotate90(val3, b.direction)

	// Row FFTs
	temp0 := val0 + val1
	temp1 := val0 - val1
	temp2 := val2 + val3
	temp3 := val2 - val3

	// Final transpose by swapping indices 1 and 2
	output[0] = temp0
	output[1] = temp2 // ← This should be temp2, not temp1
	output[2] = temp1 // ← This should be temp1, not temp2
	output[3] = temp3
}

// twiddleFactor computes a single twiddle factor
func twiddleFactor(k, n int, direction Direction) complex128 {
	angle := 2.0 * math.Pi * float64(k) / float64(n)
	if direction == Forward {
		angle = -angle
	}
	return complex(math.Cos(angle), math.Sin(angle))
}

// rotate90 rotates a complex number by 90 degrees (multiply by ±i)
func rotate90(c complex128, direction Direction) complex128 {
	if direction == Forward {
		// Multiply by -i: (a + bi) * (-i) = b - ai
		return complex(imag(c), -real(c))
	}
	// Multiply by +i: (a + bi) * i = -b + ai
	return complex(-imag(c), real(c))
}

// Butterfly8 implements a size-8 FFT
type Butterfly8 struct {
	direction Direction
	root2     float64 // sqrt(0.5) for twiddle factor computation
}

// NewButterfly8 creates a new Butterfly8 instance
func NewButterfly8(direction Direction) *Butterfly8 {
	return &Butterfly8{
		direction: direction,
		root2:     math.Sqrt(0.5),
	}
}

func (b *Butterfly8) Len() int                  { return 8 }
func (b *Butterfly8) Direction() Direction      { return b.direction }
func (b *Butterfly8) InplaceScratchLen() int    { return 0 }
func (b *Butterfly8) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly8) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly8) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly8) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 8 {
		b.performFft(buffer[i : i+8])
	}
}

func (b *Butterfly8) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 8 {
		b.performFftOutOfPlace(input[i:i+8], output[i:i+8])
	}
}

func (b *Butterfly8) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly8) performFft(buffer []complex128) {
	// Mixed radix algorithm: 2x4 FFT
	bf4 := NewButterfly4(b.direction)

	// Step 1: Transpose input into scratch arrays (even and odd indices)
	scratch0 := [4]complex128{buffer[0], buffer[2], buffer[4], buffer[6]}
	scratch1 := [4]complex128{buffer[1], buffer[3], buffer[5], buffer[7]}

	// Step 2: Column FFTs (4-point FFTs)
	bf4.performFftOutOfPlace(scratch0[:], scratch0[:])
	bf4.performFftOutOfPlace(scratch1[:], scratch1[:])

	// Step 3: Apply twiddle factors
	// twiddle[1] = (rotate_90(x) + x) * sqrt(0.5)  = (x*(-i) + x) * sqrt(0.5) for forward
	// twiddle[2] = rotate_90(x) = x * (-i) for forward
	// twiddle[3] = (rotate_90(x) - x) * sqrt(0.5) = (x*(-i) - x) * sqrt(0.5) for forward

	rot1 := rotate90(scratch1[1], b.direction)
	scratch1[1] = (rot1 + scratch1[1]) * complex(b.root2, 0)

	scratch1[2] = rotate90(scratch1[2], b.direction)

	rot3 := rotate90(scratch1[3], b.direction)
	scratch1[3] = (rot3 - scratch1[3]) * complex(b.root2, 0)

	// Step 4: Transpose - skipped because we'll do non-contiguous FFTs

	// Step 5: Row FFTs (2-point FFTs between corresponding elements)
	for i := 0; i < 4; i++ {
		temp := scratch0[i] + scratch1[i]
		scratch1[i] = scratch0[i] - scratch1[i]
		scratch0[i] = temp
	}

	// Step 6: Copy data to output (no transpose needed since we skipped step 4)
	for i := 0; i < 4; i++ {
		buffer[i] = scratch0[i]
		buffer[i+4] = scratch1[i]
	}
}

func (b *Butterfly8) performFftOutOfPlace(input, output []complex128) {
	// Copy to output and do in-place
	copy(output, input)
	b.performFft(output)
}

// Butterfly16 implements a size-16 FFT
type Butterfly16 struct {
	direction Direction
	twiddles  []complex128
}

// NewButterfly16 creates a new Butterfly16 instance
func NewButterfly16(direction Direction) *Butterfly16 {
	twiddles := computeTwiddles(16, direction)
	return &Butterfly16{
		direction: direction,
		twiddles:  twiddles,
	}
}

func (b *Butterfly16) Len() int                  { return 16 }
func (b *Butterfly16) Direction() Direction      { return b.direction }
func (b *Butterfly16) InplaceScratchLen() int    { return 0 }
func (b *Butterfly16) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly16) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly16) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly16) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 16 {
		b.performFft(buffer[i : i+16])
	}
}

func (b *Butterfly16) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 16 {
		b.performFftOutOfPlace(input[i:i+16], output[i:i+16])
	}
}

func (b *Butterfly16) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly16) performFft(buffer []complex128) {
	// Use radix-4 decomposition
	bf4 := NewButterfly4(b.direction)

	// Column FFTs
	for i := 0; i < 4; i++ {
		chunk := []complex128{buffer[i], buffer[i+4], buffer[i+8], buffer[i+12]}
		bf4.performFft(chunk)
		buffer[i], buffer[i+4], buffer[i+8], buffer[i+12] = chunk[0], chunk[1], chunk[2], chunk[3]
	}

	// Apply twiddle factors
	for row := 1; row < 4; row++ {
		for col := 0; col < 4; col++ {
			idx := row*4 + col
			buffer[idx] = buffer[idx] * b.twiddles[row*col%16]
		}
	}

	// Row FFTs
	bf4.performFft(buffer[0:4])
	bf4.performFft(buffer[4:8])
	bf4.performFft(buffer[8:12])
	bf4.performFft(buffer[12:16])

	// Transpose (simplified)
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			idx1 := i*4 + j
			idx2 := j*4 + i
			buffer[idx1], buffer[idx2] = buffer[idx2], buffer[idx1]
		}
	}
}

func (b *Butterfly16) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// bitReverse performs a bit-reversal permutation on the input
func bitReverse(data []complex128, logn int) {
	n := 1 << logn
	for i := 0; i < n; i++ {
		j := reverseBits(i, logn)
		if j > i {
			data[i], data[j] = data[j], data[i]
		}
	}
}

// reverseBits reverses the bottom n bits of x
func reverseBits(x, n int) int {
	result := 0
	for i := 0; i < n; i++ {
		result = (result << 1) | (x & 1)
		x >>= 1
	}
	return result
}

// Butterfly32 implements a size-32 FFT using split-radix algorithm
type Butterfly32 struct {
	direction   Direction
	butterfly16 *Butterfly16
	butterfly8  *Butterfly8
	twiddles    [7]complex128
}

// NewButterfly32 creates a new Butterfly32 instance
func NewButterfly32(direction Direction) *Butterfly32 {
	return &Butterfly32{
		direction:   direction,
		butterfly16: NewButterfly16(direction),
		butterfly8:  NewButterfly8(direction),
		twiddles: [7]complex128{
			twiddleFactor(1, 32, direction),
			twiddleFactor(2, 32, direction),
			twiddleFactor(3, 32, direction),
			twiddleFactor(4, 32, direction),
			twiddleFactor(5, 32, direction),
			twiddleFactor(6, 32, direction),
			twiddleFactor(7, 32, direction),
		},
	}
}

func (b *Butterfly32) Len() int                  { return 32 }
func (b *Butterfly32) Direction() Direction      { return b.direction }
func (b *Butterfly32) InplaceScratchLen() int    { return 0 }
func (b *Butterfly32) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly32) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly32) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly32) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 32 {
		b.performFft(buffer[i : i+32])
	}
}

func (b *Butterfly32) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 32 {
		b.performFftOutOfPlace(input[i:i+32], output[i:i+32])
	}
}

func (b *Butterfly32) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly32) performFft(buffer []complex128) {
	// Split-radix algorithm
	// Step 1: Split into evens and odds
	scratchEvens := [16]complex128{}
	scratchOddsN1 := [8]complex128{} // Indices 1, 5, 9, 13, 17, 21, 25, 29
	scratchOddsN3 := [8]complex128{} // Indices 31, 3, 7, 11, 15, 19, 23, 27 (wrapped)

	// Copy evens (indices 0, 2, 4, 6, ..., 30)
	for i := 0; i < 16; i++ {
		scratchEvens[i] = buffer[i*2]
	}

	// Copy odds n1 (indices 1, 5, 9, 13, 17, 21, 25, 29)
	for i := 0; i < 8; i++ {
		scratchOddsN1[i] = buffer[1+i*4]
	}

	// Copy odds n3 (indices 31, 3, 7, 11, 15, 19, 23, 27)
	scratchOddsN3[0] = buffer[31]
	for i := 1; i < 8; i++ {
		scratchOddsN3[i] = buffer[3+(i-1)*4]
	}

	// Step 2: Column FFTs
	b.butterfly16.performFft(scratchEvens[:])
	b.butterfly8.performFft(scratchOddsN1[:])
	b.butterfly8.performFft(scratchOddsN3[:])

	// Step 3: Apply twiddle factors
	for i := 1; i < 8; i++ {
		scratchOddsN1[i] = scratchOddsN1[i] * b.twiddles[i-1]
		scratchOddsN3[i] = scratchOddsN3[i] * complexConj(b.twiddles[i-1])
	}

	// Step 4: Cross FFTs (2-point butterflies between odds_n1 and odds_n3)
	for i := 0; i < 8; i++ {
		temp := scratchOddsN1[i] + scratchOddsN3[i]
		scratchOddsN3[i] = scratchOddsN1[i] - scratchOddsN3[i]
		scratchOddsN1[i] = temp
	}

	// Apply 90-degree rotation to odds_n3
	for i := 0; i < 8; i++ {
		scratchOddsN3[i] = rotate90(scratchOddsN3[i], b.direction)
	}

	// Step 5: Combine results
	// Indices 0-7: evens[0:8] + odds_n1
	for i := 0; i < 8; i++ {
		buffer[i] = scratchEvens[i] + scratchOddsN1[i]
	}
	// Indices 8-15: evens[8:16] + odds_n3
	for i := 0; i < 8; i++ {
		buffer[8+i] = scratchEvens[8+i] + scratchOddsN3[i]
	}
	// Indices 16-23: evens[0:8] - odds_n1
	for i := 0; i < 8; i++ {
		buffer[16+i] = scratchEvens[i] - scratchOddsN1[i]
	}
	// Indices 24-31: evens[8:16] - odds_n3
	for i := 0; i < 8; i++ {
		buffer[24+i] = scratchEvens[8+i] - scratchOddsN3[i]
	}
}

// complexConj returns the complex conjugate
func complexConj(c complex128) complex128 {
	return complex(real(c), -imag(c))
}

func (b *Butterfly32) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// Butterfly5 implements a size-5 FFT
type Butterfly5 struct {
	direction Direction
	twiddle1  complex128
	twiddle2  complex128
}

// NewButterfly5 creates a new Butterfly5 instance
func NewButterfly5(direction Direction) *Butterfly5 {
	return &Butterfly5{
		direction: direction,
		twiddle1:  twiddleFactor(1, 5, direction),
		twiddle2:  twiddleFactor(2, 5, direction),
	}
}

func (b *Butterfly5) Len() int                  { return 5 }
func (b *Butterfly5) Direction() Direction      { return b.direction }
func (b *Butterfly5) InplaceScratchLen() int    { return 0 }
func (b *Butterfly5) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly5) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly5) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly5) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 5 {
		b.performFft(buffer[i : i+5])
	}
}

func (b *Butterfly5) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 5 {
		b.performFftOutOfPlace(input[i:i+5], output[i:i+5])
	}
}

func (b *Butterfly5) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly5) performFft(buffer []complex128) {
	// Using the formula from RustFFT with symmetry optimizations
	x14p := buffer[1] + buffer[4]
	x14n := buffer[1] - buffer[4]
	x23p := buffer[2] + buffer[3]
	x23n := buffer[2] - buffer[3]
	sum := buffer[0] + x14p + x23p

	// Compute real parts
	b14re_a := real(buffer[0]) + real(b.twiddle1)*real(x14p) + real(b.twiddle2)*real(x23p)
	b14re_b := imag(b.twiddle1)*imag(x14n) + imag(b.twiddle2)*imag(x23n)
	b23re_a := real(buffer[0]) + real(b.twiddle2)*real(x14p) + real(b.twiddle1)*real(x23p)
	b23re_b := imag(b.twiddle2)*imag(x14n) - imag(b.twiddle1)*imag(x23n)

	// Compute imaginary parts
	b14im_a := imag(buffer[0]) + real(b.twiddle1)*imag(x14p) + real(b.twiddle2)*imag(x23p)
	b14im_b := imag(b.twiddle1)*real(x14n) + imag(b.twiddle2)*real(x23n)
	b23im_a := imag(buffer[0]) + real(b.twiddle2)*imag(x14p) + real(b.twiddle1)*imag(x23p)
	b23im_b := imag(b.twiddle2)*real(x14n) - imag(b.twiddle1)*real(x23n)

	// Assemble outputs
	buffer[0] = sum
	buffer[1] = complex(b14re_a-b14re_b, b14im_a+b14im_b)
	buffer[2] = complex(b23re_a-b23re_b, b23im_a+b23im_b)
	buffer[3] = complex(b23re_a+b23re_b, b23im_a-b23im_b)
	buffer[4] = complex(b14re_a+b14re_b, b14im_a-b14im_b)
}

func (b *Butterfly5) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// Butterfly6 implements a size-6 FFT using Good-Thomas algorithm
type Butterfly6 struct {
	direction  Direction
	butterfly3 *Butterfly3
}

// NewButterfly6 creates a new Butterfly6 instance
func NewButterfly6(direction Direction) *Butterfly6 {
	return &Butterfly6{
		direction:  direction,
		butterfly3: NewButterfly3(direction),
	}
}

func (b *Butterfly6) Len() int                  { return 6 }
func (b *Butterfly6) Direction() Direction      { return b.direction }
func (b *Butterfly6) InplaceScratchLen() int    { return 0 }
func (b *Butterfly6) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly6) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly6) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly6) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 6 {
		b.performFft(buffer[i : i+6])
	}
}

func (b *Butterfly6) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 6 {
		b.performFftOutOfPlace(input[i:i+6], output[i:i+6])
	}
}

func (b *Butterfly6) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly6) performFft(buffer []complex128) {
	// Good-Thomas algorithm (GCD(2,3) = 1, so no twiddle factors needed)
	// Step 1: Reorder input
	scratchA := [3]complex128{buffer[0], buffer[2], buffer[4]}
	scratchB := [3]complex128{buffer[3], buffer[5], buffer[1]}

	// Step 2: Column FFTs (3-point)
	b.butterfly3.performFft(scratchA[:])
	b.butterfly3.performFft(scratchB[:])

	// Step 3: Twiddle factors - SKIPPED (Good-Thomas)

	// Step 4: Transpose - SKIPPED (will do non-contiguous FFTs)

	// Step 5: Row FFTs (2-point)
	for i := 0; i < 3; i++ {
		temp := scratchA[i] + scratchB[i]
		scratchB[i] = scratchA[i] - scratchB[i]
		scratchA[i] = temp
	}

	// Step 6: Reorder output (includes transpose)
	buffer[0] = scratchA[0]
	buffer[1] = scratchB[1]
	buffer[2] = scratchA[2]
	buffer[3] = scratchB[0]
	buffer[4] = scratchA[1]
	buffer[5] = scratchB[2]
}

func (b *Butterfly6) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// Butterfly7 implements a size-7 FFT
type Butterfly7 struct {
	direction Direction
	twiddle1  complex128
	twiddle2  complex128
	twiddle3  complex128
}

// NewButterfly7 creates a new Butterfly7 instance
func NewButterfly7(direction Direction) *Butterfly7 {
	return &Butterfly7{
		direction: direction,
		twiddle1:  twiddleFactor(1, 7, direction),
		twiddle2:  twiddleFactor(2, 7, direction),
		twiddle3:  twiddleFactor(3, 7, direction),
	}
}

func (b *Butterfly7) Len() int                  { return 7 }
func (b *Butterfly7) Direction() Direction      { return b.direction }
func (b *Butterfly7) InplaceScratchLen() int    { return 0 }
func (b *Butterfly7) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly7) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly7) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly7) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 7 {
		b.performFft(buffer[i : i+7])
	}
}

func (b *Butterfly7) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 7 {
		b.performFftOutOfPlace(input[i:i+7], output[i:i+7])
	}
}

func (b *Butterfly7) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly7) performFft(buffer []complex128) {
	// For size 7, use symmetry: W3=W4*, W5=W2*, W6=W1*
	x16p := buffer[1] + buffer[6]
	x16n := buffer[1] - buffer[6]
	x25p := buffer[2] + buffer[5]
	x25n := buffer[2] - buffer[5]
	x34p := buffer[3] + buffer[4]
	x34n := buffer[3] - buffer[4]

	sum := buffer[0] + x16p + x25p + x34p

	// Real parts for output 1, 6
	b16re_a := real(buffer[0]) + real(b.twiddle1)*real(x16p) + real(b.twiddle2)*real(x25p) + real(b.twiddle3)*real(x34p)
	b16re_b := imag(b.twiddle1)*imag(x16n) + imag(b.twiddle2)*imag(x25n) + imag(b.twiddle3)*imag(x34n)

	// Imaginary parts for output 1, 6
	b16im_a := imag(buffer[0]) + real(b.twiddle1)*imag(x16p) + real(b.twiddle2)*imag(x25p) + real(b.twiddle3)*imag(x34p)
	b16im_b := imag(b.twiddle1)*real(x16n) + imag(b.twiddle2)*real(x25n) + imag(b.twiddle3)*real(x34n)

	// Real parts for output 2, 5
	b25re_a := real(buffer[0]) + real(b.twiddle2)*real(x16p) + real(b.twiddle3)*real(x25p) + real(b.twiddle1)*real(x34p)
	b25re_b := imag(b.twiddle2)*imag(x16n) - imag(b.twiddle3)*imag(x25n) - imag(b.twiddle1)*imag(x34n)

	// Imaginary parts for output 2, 5
	b25im_a := imag(buffer[0]) + real(b.twiddle2)*imag(x16p) + real(b.twiddle3)*imag(x25p) + real(b.twiddle1)*imag(x34p)
	b25im_b := imag(b.twiddle2)*real(x16n) - imag(b.twiddle3)*real(x25n) - imag(b.twiddle1)*real(x34n)

	// Real parts for output 3, 4
	b34re_a := real(buffer[0]) + real(b.twiddle3)*real(x16p) + real(b.twiddle1)*real(x25p) + real(b.twiddle2)*real(x34p)
	b34re_b := imag(b.twiddle3)*imag(x16n) - imag(b.twiddle1)*imag(x25n) + imag(b.twiddle2)*imag(x34n)

	// Imaginary parts for output 3, 4
	b34im_a := imag(buffer[0]) + real(b.twiddle3)*imag(x16p) + real(b.twiddle1)*imag(x25p) + real(b.twiddle2)*imag(x34p)
	b34im_b := imag(b.twiddle3)*real(x16n) - imag(b.twiddle1)*real(x25n) + imag(b.twiddle2)*real(x34n)

	buffer[0] = sum
	buffer[1] = complex(b16re_a-b16re_b, b16im_a+b16im_b)
	buffer[2] = complex(b25re_a-b25re_b, b25im_a+b25im_b)
	buffer[3] = complex(b34re_a-b34re_b, b34im_a+b34im_b)
	buffer[4] = complex(b34re_a+b34re_b, b34im_a-b34im_b)
	buffer[5] = complex(b25re_a+b25re_b, b25im_a-b25im_b)
	buffer[6] = complex(b16re_a+b16re_b, b16im_a-b16im_b)
}

func (b *Butterfly7) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// Butterfly9 implements a size-9 FFT
type Butterfly9 struct {
	direction  Direction
	butterfly3 *Butterfly3
	twiddle1   complex128
	twiddle2   complex128
	twiddle4   complex128
}

// NewButterfly9 creates a new Butterfly9 instance
func NewButterfly9(direction Direction) *Butterfly9 {
	return &Butterfly9{
		direction:  direction,
		butterfly3: NewButterfly3(direction),
		twiddle1:   twiddleFactor(1, 9, direction),
		twiddle2:   twiddleFactor(2, 9, direction),
		twiddle4:   twiddleFactor(4, 9, direction),
	}
}

func (b *Butterfly9) Len() int                  { return 9 }
func (b *Butterfly9) Direction() Direction      { return b.direction }
func (b *Butterfly9) InplaceScratchLen() int    { return 0 }
func (b *Butterfly9) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly9) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly9) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly9) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 9 {
		b.performFft(buffer[i : i+9])
	}
}

func (b *Butterfly9) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 9 {
		b.performFftOutOfPlace(input[i:i+9], output[i:i+9])
	}
}

func (b *Butterfly9) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly9) performFft(buffer []complex128) {
	// Mixed radix algorithm: 3x3 FFT
	// Step 1: Transpose input into scratch
	scratch0 := [3]complex128{buffer[0], buffer[3], buffer[6]}
	scratch1 := [3]complex128{buffer[1], buffer[4], buffer[7]}
	scratch2 := [3]complex128{buffer[2], buffer[5], buffer[8]}

	// Step 2: Column FFTs
	b.butterfly3.performFft(scratch0[:])
	b.butterfly3.performFft(scratch1[:])
	b.butterfly3.performFft(scratch2[:])

	// Step 3: Apply twiddle factors
	scratch1[1] = scratch1[1] * b.twiddle1
	scratch1[2] = scratch1[2] * b.twiddle2
	scratch2[1] = scratch2[1] * b.twiddle2
	scratch2[2] = scratch2[2] * b.twiddle4

	// Step 4: Transpose - SKIPPED

	// Step 5: Row FFTs (3-point, strided across scratch arrays)
	performStrided3(&scratch0[0], &scratch1[0], &scratch2[0], b.butterfly3.twiddle)
	performStrided3(&scratch0[1], &scratch1[1], &scratch2[1], b.butterfly3.twiddle)
	performStrided3(&scratch0[2], &scratch1[2], &scratch2[2], b.butterfly3.twiddle)

	// Step 6: Copy to output (column-major)
	buffer[0] = scratch0[0]
	buffer[1] = scratch0[1]
	buffer[2] = scratch0[2]
	buffer[3] = scratch1[0]
	buffer[4] = scratch1[1]
	buffer[5] = scratch1[2]
	buffer[6] = scratch2[0]
	buffer[7] = scratch2[1]
	buffer[8] = scratch2[2]
}

// performStrided3 performs a 3-point FFT on values passed by pointer (strided access)
func performStrided3(val0, val1, val2 *complex128, twiddle complex128) {
	xp := *val1 + *val2
	xn := *val1 - *val2
	sum := *val0 + xp

	tempA := *val0 + complex(real(twiddle)*real(xp), real(twiddle)*imag(xp))
	tempB := complex(-imag(twiddle)*imag(xn), imag(twiddle)*real(xn))

	*val0 = sum
	*val1 = tempA + tempB
	*val2 = tempA - tempB
}

func (b *Butterfly9) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// Butterfly12 implements a size-12 FFT
type Butterfly12 struct {
	direction  Direction
	butterfly3 *Butterfly3
	butterfly4 *Butterfly4
	twiddle1   complex128
	twiddle2   complex128
}

// NewButterfly12 creates a new Butterfly12 instance
func NewButterfly12(direction Direction) *Butterfly12 {
	return &Butterfly12{
		direction:  direction,
		butterfly3: NewButterfly3(direction),
		butterfly4: NewButterfly4(direction),
		twiddle1:   twiddleFactor(1, 12, direction),
		twiddle2:   twiddleFactor(2, 12, direction),
	}
}

func (b *Butterfly12) Len() int                  { return 12 }
func (b *Butterfly12) Direction() Direction      { return b.direction }
func (b *Butterfly12) InplaceScratchLen() int    { return 0 }
func (b *Butterfly12) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly12) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly12) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly12) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 12 {
		b.performFft(buffer[i : i+12])
	}
}

func (b *Butterfly12) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 12 {
		b.performFftOutOfPlace(input[i:i+12], output[i:i+12])
	}
}

func (b *Butterfly12) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly12) performFft(buffer []complex128) {
	// Good-Thomas algorithm (GCD(4,3) = 1, so no twiddle factors needed)
	// Step 1: Reorder input with precomputed Good-Thomas indices
	scratch0 := [4]complex128{buffer[0], buffer[3], buffer[6], buffer[9]}
	scratch1 := [4]complex128{buffer[4], buffer[7], buffer[10], buffer[1]}
	scratch2 := [4]complex128{buffer[8], buffer[11], buffer[2], buffer[5]}

	// Step 2: Column FFTs (4-point)
	b.butterfly4.performFft(scratch0[:])
	b.butterfly4.performFft(scratch1[:])
	b.butterfly4.performFft(scratch2[:])

	// Step 3: Twiddle factors - SKIPPED (Good-Thomas)

	// Step 4: Transpose - SKIPPED (will do non-contiguous FFTs)

	// Step 5: Row FFTs (3-point, strided across scratch arrays)
	performStrided3(&scratch0[0], &scratch1[0], &scratch2[0], b.butterfly3.twiddle)
	performStrided3(&scratch0[1], &scratch1[1], &scratch2[1], b.butterfly3.twiddle)
	performStrided3(&scratch0[2], &scratch1[2], &scratch2[2], b.butterfly3.twiddle)
	performStrided3(&scratch0[3], &scratch1[3], &scratch2[3], b.butterfly3.twiddle)

	// Step 6: Reorder output with Good-Thomas pattern (includes transpose)
	buffer[0] = scratch0[0]
	buffer[1] = scratch1[1]
	buffer[2] = scratch2[2]
	buffer[3] = scratch0[3]
	buffer[4] = scratch1[0]
	buffer[5] = scratch2[1]
	buffer[6] = scratch0[2]
	buffer[7] = scratch1[3]
	buffer[8] = scratch2[0]
	buffer[9] = scratch0[1]
	buffer[10] = scratch1[2]
	buffer[11] = scratch2[3]
}

func (b *Butterfly12) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}
