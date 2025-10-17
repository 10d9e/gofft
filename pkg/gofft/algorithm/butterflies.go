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
	output[0] = val0 + val1
	output[2] = val0 - val1
	output[1] = val2 + val3
	output[3] = val2 - val3
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
