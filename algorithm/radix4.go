package algorithm

import (
	"math"
)

// FftInterface is a minimal interface for FFT algorithms
type FftInterface interface {
	Len() int
	Direction() Direction
	InplaceScratchLen() int
	OutOfPlaceScratchLen() int
	ProcessWithScratch(buffer, scratch []complex128)
	ProcessOutOfPlace(input, output, scratch []complex128)
}

// Radix4 implements an FFT algorithm optimized for power-of-two sizes
// It uses a radix-4 decimation-in-frequency algorithm
type Radix4 struct {
	twiddles          []complex128
	baseFft           FftInterface
	baseLen           int
	length            int
	direction         Direction
	inplaceScratch    int
	outofplaceScratch int
}

// NewRadix4 creates a new Radix4 FFT instance for the given power-of-two length
func NewRadix4(length int, direction Direction) *Radix4 {
	if !isPowerOfTwo(length) {
		panic("Radix4 algorithm requires a power-of-two input size")
	}

	// Figure out which base length to use (match RustFFT's selection logic exactly)
	exponent := trailingZeros(length)
	var baseFft FftInterface
	var baseExponent int

	switch exponent {
	case 0:
		// Length 1 - trivial case
		baseExponent = 0
		baseFft = &trivialFft{direction: direction}
	case 1:
		baseExponent = 1
		baseFft = NewButterfly2(direction)
	case 2:
		baseExponent = 2
		baseFft = NewButterfly4(direction)
	case 3:
		baseExponent = 3
		baseFft = NewButterfly8(direction)
	default:
		// For larger sizes, match RustFFT's choice:
		// - Odd exponent: use Butterfly32 (base_exponent=5)
		// - Even exponent: use Butterfly16 (base_exponent=4)
		if exponent%2 == 1 {
			baseExponent = 5
			baseFft = NewButterfly32(direction)
		} else {
			baseExponent = 4
			baseFft = NewButterfly16(direction)
		}
	}

	k := (exponent - baseExponent) / 2
	return NewRadix4WithBase(k, baseFft)
}

// trailingZeros counts the number of trailing zero bits
func trailingZeros(n int) int {
	if n == 0 {
		return 64
	}
	count := 0
	for n&1 == 0 {
		count++
		n >>= 1
	}
	return count
}

// NewRadix4WithBase creates a Radix4 instance that computes FFTs of length 4^k * baseFft.Len()
func NewRadix4WithBase(k int, baseFft FftInterface) *Radix4 {
	baseLen := baseFft.Len()
	length := baseLen * (1 << (k * 2))
	direction := baseFft.Direction()

	// Precompute twiddle factors
	const rowCount = 4
	crossFftLen := baseLen
	twiddleFactors := make([]complex128, 0, length*2)

	for crossFftLen < length {
		numColumns := crossFftLen
		crossFftLen *= rowCount

		for i := 0; i < numColumns; i++ {
			for k := 1; k < rowCount; k++ {
				angle := 2.0 * math.Pi * float64(i*k) / float64(crossFftLen)
				if direction == Forward {
					angle = -angle
				}
				twiddle := complex(math.Cos(angle), math.Sin(angle))
				twiddleFactors = append(twiddleFactors, twiddle)
			}
		}
	}

	baseInplaceScratch := baseFft.InplaceScratchLen()
	inplaceScratch := crossFftLen
	if baseInplaceScratch > crossFftLen {
		inplaceScratch = crossFftLen + baseInplaceScratch
	}

	outofplaceScratch := 0
	if baseInplaceScratch > length {
		outofplaceScratch = baseInplaceScratch
	}

	return &Radix4{
		twiddles:          twiddleFactors,
		baseFft:           baseFft,
		baseLen:           baseLen,
		length:            length,
		direction:         direction,
		inplaceScratch:    inplaceScratch,
		outofplaceScratch: outofplaceScratch,
	}
}

func (r *Radix4) Len() int                  { return r.length }
func (r *Radix4) Direction() Direction      { return r.direction }
func (r *Radix4) InplaceScratchLen() int    { return r.inplaceScratch }
func (r *Radix4) OutOfPlaceScratchLen() int { return r.outofplaceScratch }
func (r *Radix4) ImmutableScratchLen() int  { return r.baseFft.InplaceScratchLen() }

func (r *Radix4) Process(buffer []complex128) {
	scratch := make([]complex128, r.InplaceScratchLen())
	r.ProcessWithScratch(buffer, scratch)
}

func (r *Radix4) ProcessWithScratch(buffer, scratch []complex128) {
	for i := 0; i < len(buffer); i += r.length {
		chunk := buffer[i : i+r.length]
		selfScratch := scratch[:r.length]
		r.performFftOutOfPlace(chunk, selfScratch, scratch[r.length:])
		copy(chunk, selfScratch)
	}
}

func (r *Radix4) ProcessOutOfPlace(input, output, scratch []complex128) {
	for i := 0; i < len(input); i += r.length {
		inChunk := input[i : i+r.length]
		outChunk := output[i : i+r.length]
		r.performFftOutOfPlace(inChunk, outChunk, scratch)
	}
}

func (r *Radix4) ProcessImmutable(input []complex128, output, scratch []complex128) {
	for i := 0; i < len(input); i += r.length {
		inChunk := input[i : i+r.length]
		outChunk := output[i : i+r.length]
		r.performFftImmutable(inChunk, outChunk, scratch)
	}
}

func (r *Radix4) performFftImmutable(input []complex128, output []complex128, scratch []complex128) {
	// Copy data with bit-reversed transpose
	if r.length == r.baseLen {
		copy(output, input)
	} else {
		bitReversedTranspose4(r.baseLen, input, output)
	}

	// Base-level FFTs
	r.baseFft.ProcessWithScratch(output, scratch)

	// Cross FFTs
	r.performCrossFfts(output)
}

func (r *Radix4) performFftOutOfPlace(input []complex128, output []complex128, scratch []complex128) {
	// Copy data with bit-reversed transpose
	if r.length == r.baseLen {
		copy(output, input)
	} else {
		bitReversedTranspose4(r.baseLen, input, output)
	}

	// Base-level FFTs
	baseScratch := scratch
	if len(scratch) == 0 {
		baseScratch = input
	}
	r.baseFft.ProcessWithScratch(output, baseScratch)

	// Cross FFTs
	r.performCrossFfts(output)
}

func (r *Radix4) performCrossFfts(output []complex128) {
	const rowCount = 4
	crossFftLen := r.baseLen
	layerTwiddles := r.twiddles
	butterfly4 := NewButterfly4(r.direction)

	for crossFftLen < len(output) {
		numColumns := crossFftLen
		crossFftLen *= rowCount

		// Process each chunk
		for offset := 0; offset < len(output); offset += crossFftLen {
			data := output[offset : offset+crossFftLen]
			butterfly4Stage(data, layerTwiddles, numColumns, butterfly4)
		}

		// Skip past twiddle factors used in this layer
		twiddleOffset := numColumns * (rowCount - 1)
		layerTwiddles = layerTwiddles[twiddleOffset:]
	}
}

// butterfly4Stage applies a radix-4 butterfly stage
func butterfly4Stage(data []complex128, twiddles []complex128, numColumns int, butterfly4 *Butterfly4) {
	// Apply twiddle factors and perform radix-4 butterflies
	for col := 0; col < numColumns; col++ {
		// Get the four values for this column
		idx0 := col
		idx1 := col + numColumns
		idx2 := col + 2*numColumns
		idx3 := col + 3*numColumns

		// Load values and apply twiddle factors (first row doesn't need twiddles)
		twIdx := col * 3
		scratch := [4]complex128{
			data[idx0],
			data[idx1] * twiddles[twIdx+0],
			data[idx2] * twiddles[twIdx+1],
			data[idx3] * twiddles[twIdx+2],
		}

		// Perform 4-point butterfly FFT on scratch array
		butterfly4.performFftOutOfPlace(scratch[:], scratch[:])

		// Store results back
		data[idx0] = scratch[0]
		data[idx1] = scratch[1]
		data[idx2] = scratch[2]
		data[idx3] = scratch[3]
	}
}

// bitReversedTranspose4 performs a bit-reversed transpose with divisor 4
// This is a port of RustFFT's bitreversed_transpose::<T, 4>
func bitReversedTranspose4(height int, input, output []complex128) {
	const D = 4 // Divisor for bit reversal
	width := len(input) / height

	if len(input)%height != 0 || len(input) != len(output) {
		panic("invalid dimensions for bitreversed_transpose")
	}

	stridedWidth := width / D

	// Compute how many "digits" we need for base-D bit reversal
	revDigits := 0
	temp := width
	for temp > 1 {
		if temp%D != 0 {
			panic("width must be a power of D")
		}
		temp /= D
		revDigits++
	}

	for x := 0; x < stridedWidth; x++ {
		// Create forward and reversed indices
		xFwd := [D]int{}
		xRev := [D]int{}

		for i := 0; i < D; i++ {
			xFwd[i] = D*x + i
			xRev[i] = reverseBitsBaseD(xFwd[i], revDigits, D)
		}

		// Transpose with bit-reversed columns
		for y := 0; y < height; y++ {
			for i := 0; i < D; i++ {
				inputIndex := xFwd[i] + y*width
				outputIndex := y + xRev[i]*height
				output[outputIndex] = input[inputIndex]
			}
		}
	}
}

// reverseBitsBaseD reverses digits in base D
// This is like bit reversal but works for any base, not just base 2
func reverseBitsBaseD(value, revDigits, D int) int {
	result := 0
	for i := 0; i < revDigits; i++ {
		result = (result * D) + (value % D)
		value = value / D
	}
	return result
}

// bitReverseIndex reverses the bits of x considering n values (n must be power of 2)
func bitReverseIndex(x, n int) int {
	bits := trailingZeros(n)
	result := 0
	for i := 0; i < bits; i++ {
		if x&(1<<i) != 0 {
			result |= 1 << (bits - 1 - i)
		}
	}
	return result
}

// isPowerOfTwo checks if n is a power of two
func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// trivialFft is a trivial FFT for length 1 (does nothing)
type trivialFft struct {
	direction Direction
}

func (t *trivialFft) Len() int                                        { return 1 }
func (t *trivialFft) Direction() Direction                            { return t.direction }
func (t *trivialFft) InplaceScratchLen() int                          { return 0 }
func (t *trivialFft) OutOfPlaceScratchLen() int                       { return 0 }
func (t *trivialFft) ImmutableScratchLen() int                        { return 0 }
func (t *trivialFft) Process(buffer []complex128)                     {}
func (t *trivialFft) ProcessWithScratch(buffer, scratch []complex128) {}
func (t *trivialFft) ProcessOutOfPlace(input, output, scratch []complex128) {
	copy(output, input)
}
func (t *trivialFft) ProcessImmutable(input []complex128, output, scratch []complex128) {
	copy(output, input)
}
