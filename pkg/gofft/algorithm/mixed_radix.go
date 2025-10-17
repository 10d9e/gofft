package algorithm

import "math"

// MixedRadix implements the Mixed-Radix FFT algorithm
// It factors a size n FFT into n1 * n2, computes several inner FFTs, then combines results
type MixedRadix struct {
	twiddles          []complex128
	widthFft          FftInterface
	width             int
	heightFft         FftInterface
	height            int
	length            int
	direction         Direction
	inplaceScratch    int
	outofplaceScratch int
}

// NewMixedRadix creates a MixedRadix FFT instance
// The FFT size will be widthFft.Len() * heightFft.Len()
func NewMixedRadix(widthFft, heightFft FftInterface) *MixedRadix {
	if widthFft.Direction() != heightFft.Direction() {
		panic("width and height FFTs must have the same direction")
	}

	direction := widthFft.Direction()
	width := widthFft.Len()
	height := heightFft.Len()
	length := width * height

	// Precompute twiddle factors
	twiddles := make([]complex128, length)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			idx := x*height + y
			angle := -2.0 * math.Pi * float64(x*y) / float64(length)
			if direction == Inverse {
				angle = -angle
			}
			twiddles[idx] = complex(math.Cos(angle), math.Sin(angle))
		}
	}

	// Calculate scratch space requirements
	heightInplace := heightFft.InplaceScratchLen()
	widthInplace := widthFft.InplaceScratchLen()
	widthOutofplace := 0 // We'll use in-place for width FFT

	maxInnerInplace := heightInplace
	if widthInplace > maxInnerInplace {
		maxInnerInplace = widthInplace
	}

	outofplaceScratch := 0
	if maxInnerInplace > length {
		outofplaceScratch = maxInnerInplace
	}

	inplaceScratch := length
	if heightInplace > length {
		inplaceScratch = length + (heightInplace - length)
	}
	if widthOutofplace > 0 && widthOutofplace > inplaceScratch-length {
		inplaceScratch = length + widthOutofplace
	}

	return &MixedRadix{
		twiddles:          twiddles,
		widthFft:          widthFft,
		width:             width,
		heightFft:         heightFft,
		height:            height,
		length:            length,
		direction:         direction,
		inplaceScratch:    inplaceScratch,
		outofplaceScratch: outofplaceScratch,
	}
}

func (m *MixedRadix) Len() int                  { return m.length }
func (m *MixedRadix) Direction() Direction      { return m.direction }
func (m *MixedRadix) InplaceScratchLen() int    { return m.inplaceScratch }
func (m *MixedRadix) OutOfPlaceScratchLen() int { return m.outofplaceScratch }
func (m *MixedRadix) ImmutableScratchLen() int  { return m.inplaceScratch }

func (m *MixedRadix) Process(buffer []complex128) {
	scratch := make([]complex128, m.InplaceScratchLen())
	m.ProcessWithScratch(buffer, scratch)
}

func (m *MixedRadix) ProcessWithScratch(buffer, scratch []complex128) {
	// Six-step FFT algorithm (based on RustFFT)
	selfScratch := scratch[:m.length]
	var innerScratch []complex128
	if len(scratch) > m.length {
		innerScratch = scratch[m.length:]
	}

	// STEP 1: Transpose input (width x height) to (height x width)
	transpose(m.width, m.height, buffer, selfScratch)

	// STEP 2: Perform height-sized FFTs
	// The heightFft will process multiple FFTs of size height
	heightScratch := buffer // Use buffer as scratch since we've copied data to selfScratch
	if len(innerScratch) >= len(buffer) {
		heightScratch = innerScratch
	}
	m.heightFft.ProcessWithScratch(selfScratch, heightScratch)

	// STEP 3: Apply twiddle factors
	for i := range selfScratch {
		selfScratch[i] = selfScratch[i] * m.twiddles[i]
	}

	// STEP 4: Transpose back to (width x height)
	transpose(m.height, m.width, selfScratch, buffer)

	// STEP 5: Perform width-sized FFTs
	// The widthFft will process multiple FFTs of size width
	m.widthFft.ProcessWithScratch(buffer, heightScratch)

	// Result is now in buffer, which is correct
}

func (m *MixedRadix) ProcessOutOfPlace(input, output, scratch []complex128) {
	copy(output, input)
	m.ProcessWithScratch(output, scratch)
}

func (m *MixedRadix) ProcessImmutable(input []complex128, output, scratch []complex128) {
	copy(output, input)
	m.ProcessWithScratch(output, scratch)
}

// transpose performs a matrix transpose
// Treats input as a rows x cols matrix and transposes to output
func transpose(rows, cols int, input, output []complex128) {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			inputIdx := r*cols + c
			outputIdx := c*rows + r
			output[outputIdx] = input[inputIdx]
		}
	}
}
