package algorithm

import "math"

// This file contains additional butterfly implementations

// Butterfly11 implements a size-11 FFT (prime size)
type Butterfly11 struct {
	direction Direction
	twiddles  [5]complex128 // W1, W2, W3, W4, W5 (W6-W10 are conjugates)
}

// NewButterfly11 creates a new Butterfly11 instance
func NewButterfly11(direction Direction) *Butterfly11 {
	return &Butterfly11{
		direction: direction,
		twiddles: [5]complex128{
			twiddleFactor(1, 11, direction),
			twiddleFactor(2, 11, direction),
			twiddleFactor(3, 11, direction),
			twiddleFactor(4, 11, direction),
			twiddleFactor(5, 11, direction),
		},
	}
}

func (b *Butterfly11) Len() int                  { return 11 }
func (b *Butterfly11) Direction() Direction      { return b.direction }
func (b *Butterfly11) InplaceScratchLen() int    { return 0 }
func (b *Butterfly11) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly11) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly11) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly11) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 11 {
		b.performFft(buffer[i : i+11])
	}
}

func (b *Butterfly11) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 11 {
		b.performFftOutOfPlace(input[i:i+11], output[i:i+11])
	}
}

func (b *Butterfly11) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly11) performFft(buffer []complex128) {
	// For size 11 (prime), use DFT for now
	// TODO: Implement optimized version with symmetry
	temp := make([]complex128, 11)
	copy(temp, buffer)
	for k := 0; k < 11; k++ {
		sum := complex(0, 0)
		for j := 0; j < 11; j++ {
			angle := -2.0 * math.Pi * float64(k*j) / 11.0
			if b.direction == Inverse {
				angle = -angle
			}
			tw := complex(math.Cos(angle), math.Sin(angle))
			sum += temp[j] * tw
		}
		buffer[k] = sum
	}
}

func (b *Butterfly11) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// Butterfly13 implements a size-13 FFT (prime size)
type Butterfly13 struct {
	direction Direction
	twiddles  [6]complex128 // W1-W6 (W7-W12 are conjugates)
}

// NewButterfly13 creates a new Butterfly13 instance
func NewButterfly13(direction Direction) *Butterfly13 {
	return &Butterfly13{
		direction: direction,
		twiddles: [6]complex128{
			twiddleFactor(1, 13, direction),
			twiddleFactor(2, 13, direction),
			twiddleFactor(3, 13, direction),
			twiddleFactor(4, 13, direction),
			twiddleFactor(5, 13, direction),
			twiddleFactor(6, 13, direction),
		},
	}
}

func (b *Butterfly13) Len() int                  { return 13 }
func (b *Butterfly13) Direction() Direction      { return b.direction }
func (b *Butterfly13) InplaceScratchLen() int    { return 0 }
func (b *Butterfly13) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly13) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly13) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly13) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 13 {
		b.performFft(buffer[i : i+13])
	}
}

func (b *Butterfly13) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 13 {
		b.performFftOutOfPlace(input[i:i+13], output[i:i+13])
	}
}

func (b *Butterfly13) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly13) performFft(buffer []complex128) {
	// For size 13 (prime), use DFT
	// TODO: Implement optimized version with symmetry
	temp := make([]complex128, 13)
	copy(temp, buffer)
	for k := 0; k < 13; k++ {
		sum := complex(0, 0)
		for j := 0; j < 13; j++ {
			angle := -2.0 * math.Pi * float64(k*j) / 13.0
			if b.direction == Inverse {
				angle = -angle
			}
			tw := complex(math.Cos(angle), math.Sin(angle))
			sum += temp[j] * tw
		}
		buffer[k] = sum
	}
}

func (b *Butterfly13) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// Butterfly24 implements a size-24 FFT
type Butterfly24 struct {
	direction  Direction
	butterfly4 *Butterfly4
	butterfly6 *Butterfly6
}

// NewButterfly24 creates a new Butterfly24 instance
func NewButterfly24(direction Direction) *Butterfly24 {
	return &Butterfly24{
		direction:  direction,
		butterfly4: NewButterfly4(direction),
		butterfly6: NewButterfly6(direction),
	}
}

func (b *Butterfly24) Len() int                  { return 24 }
func (b *Butterfly24) Direction() Direction      { return b.direction }
func (b *Butterfly24) InplaceScratchLen() int    { return 0 }
func (b *Butterfly24) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly24) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly24) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly24) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 24 {
		b.performFft(buffer[i : i+24])
	}
}

func (b *Butterfly24) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 24 {
		b.performFftOutOfPlace(input[i:i+24], output[i:i+24])
	}
}

func (b *Butterfly24) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly24) performFft(buffer []complex128) {
	// For now, use DFT for size 24
	// TODO: Implement proper mixed-radix 6x4 algorithm
	dft := NewDft(24, b.direction)
	temp := make([]complex128, 24)
	copy(temp, buffer)
	dft.performFftImmutable(temp, buffer, nil)
}

func (b *Butterfly24) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// Butterfly27 implements a size-27 FFT (3^3)
type Butterfly27 struct {
	direction  Direction
	butterfly9 *Butterfly9
}

// NewButterfly27 creates a new Butterfly27 instance
func NewButterfly27(direction Direction) *Butterfly27 {
	return &Butterfly27{
		direction:  direction,
		butterfly9: NewButterfly9(direction),
	}
}

func (b *Butterfly27) Len() int                  { return 27 }
func (b *Butterfly27) Direction() Direction      { return b.direction }
func (b *Butterfly27) InplaceScratchLen() int    { return 0 }
func (b *Butterfly27) OutOfPlaceScratchLen() int { return 0 }
func (b *Butterfly27) ImmutableScratchLen() int  { return 0 }

func (b *Butterfly27) Process(buffer []complex128) {
	b.ProcessWithScratch(buffer, nil)
}

func (b *Butterfly27) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += 27 {
		b.performFft(buffer[i : i+27])
	}
}

func (b *Butterfly27) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += 27 {
		b.performFftOutOfPlace(input[i:i+27], output[i:i+27])
	}
}

func (b *Butterfly27) ProcessImmutable(input []complex128, output, scratch []complex128) {
	b.ProcessOutOfPlace(input, output, scratch)
}

func (b *Butterfly27) performFft(buffer []complex128) {
	// For now, use DFT for size 27
	// TODO: Implement proper 9x3 mixed radix
	dft := NewDft(27, b.direction)
	temp := make([]complex128, 27)
	copy(temp, buffer)
	dft.performFftImmutable(temp, buffer, nil)
}

func (b *Butterfly27) performFftOutOfPlace(input, output []complex128) {
	copy(output, input)
	b.performFft(output)
}

// addPrimeButterfly is a helper to create prime-sized butterflies
// For primes, we use DFT-based computation (RustFFT optimizes with symmetry, but DFT is simpler and correct)
func createPrimeButterfly(size int, direction Direction) FftInterface {
	return NewDft(size, direction)
}

// Butterfly17 implements a size-17 FFT (prime)
type Butterfly17 struct {
	inner *Dft
}

func NewButterfly17(direction Direction) *Butterfly17 {
	return &Butterfly17{inner: NewDft(17, direction)}
}

func (b *Butterfly17) Len() int               { return 17 }
func (b *Butterfly17) Direction() Direction   { return b.inner.Direction() }
func (b *Butterfly17) InplaceScratchLen() int { return b.inner.InplaceScratchLen() }
func (b *Butterfly17) ProcessWithScratch(buffer, scratch []complex128) {
	b.inner.ProcessWithScratch(buffer, scratch)
}

// Butterfly19 implements a size-19 FFT (prime)
type Butterfly19 struct {
	inner *Dft
}

func NewButterfly19(direction Direction) *Butterfly19 {
	return &Butterfly19{inner: NewDft(19, direction)}
}

func (b *Butterfly19) Len() int               { return 19 }
func (b *Butterfly19) Direction() Direction   { return b.inner.Direction() }
func (b *Butterfly19) InplaceScratchLen() int { return b.inner.InplaceScratchLen() }
func (b *Butterfly19) ProcessWithScratch(buffer, scratch []complex128) {
	b.inner.ProcessWithScratch(buffer, scratch)
}

// Butterfly23 implements a size-23 FFT (prime)
type Butterfly23 struct {
	inner *Dft
}

func NewButterfly23(direction Direction) *Butterfly23 {
	return &Butterfly23{inner: NewDft(23, direction)}
}

func (b *Butterfly23) Len() int               { return 23 }
func (b *Butterfly23) Direction() Direction   { return b.inner.Direction() }
func (b *Butterfly23) InplaceScratchLen() int { return b.inner.InplaceScratchLen() }
func (b *Butterfly23) ProcessWithScratch(buffer, scratch []complex128) {
	b.inner.ProcessWithScratch(buffer, scratch)
}

// Butterfly29 implements a size-29 FFT (prime)
type Butterfly29 struct {
	inner *Dft
}

func NewButterfly29(direction Direction) *Butterfly29 {
	return &Butterfly29{inner: NewDft(29, direction)}
}

func (b *Butterfly29) Len() int               { return 29 }
func (b *Butterfly29) Direction() Direction   { return b.inner.Direction() }
func (b *Butterfly29) InplaceScratchLen() int { return b.inner.InplaceScratchLen() }
func (b *Butterfly29) ProcessWithScratch(buffer, scratch []complex128) {
	b.inner.ProcessWithScratch(buffer, scratch)
}

// Butterfly31 implements a size-31 FFT (prime)
type Butterfly31 struct {
	inner *Dft
}

func NewButterfly31(direction Direction) *Butterfly31 {
	return &Butterfly31{inner: NewDft(31, direction)}
}

func (b *Butterfly31) Len() int               { return 31 }
func (b *Butterfly31) Direction() Direction   { return b.inner.Direction() }
func (b *Butterfly31) InplaceScratchLen() int { return b.inner.InplaceScratchLen() }
func (b *Butterfly31) ProcessWithScratch(buffer, scratch []complex128) {
	b.inner.ProcessWithScratch(buffer, scratch)
}
