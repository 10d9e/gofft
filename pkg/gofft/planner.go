package gofft

import (
	"sync"

	"github.com/example/gofft/pkg/gofft/algorithm"
)

// Planner creates FFT instances for arbitrary sizes
// It automatically selects the best algorithm and caches created instances
type Planner struct {
	mu          sync.Mutex
	cache       map[plannerKey]Fft
	recipeCache map[int]*recipe
}

type plannerKey struct {
	length    int
	direction Direction
}

// NewPlanner creates a new FFT planner
func NewPlanner() *Planner {
	return &Planner{
		cache:       make(map[plannerKey]Fft),
		recipeCache: make(map[int]*recipe),
	}
}

// PlanForward creates an FFT instance for computing forward FFTs of the given size
func (p *Planner) PlanForward(length int) Fft {
	return p.Plan(length, Forward)
}

// PlanInverse creates an FFT instance for computing inverse FFTs of the given size
func (p *Planner) PlanInverse(length int) Fft {
	return p.Plan(length, Inverse)
}

// Plan creates an FFT instance for the given size and direction
func (p *Planner) Plan(length int, direction Direction) Fft {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := plannerKey{length: length, direction: direction}

	// Check cache
	if fft, ok := p.cache[key]; ok {
		return fft
	}

	// Create a recipe for this FFT
	recipe, len := p.designFft(length)

	// Build the FFT from the recipe
	fft := p.buildFft(recipe, len, direction)

	// Cache it
	p.cache[key] = fft

	return fft
}

// recipe describes how to construct an FFT without actually building it
type recipe int

const (
	recipeDft recipe = iota
	recipeButterfly2
	recipeButterfly3
	recipeButterfly4
	recipeButterfly8
	recipeButterfly16
	recipeButterfly32
	recipeRadix4
)

// designFft creates a recipe for an FFT of the given length
func (p *Planner) designFft(length int) (*recipe, int) {
	if r, ok := p.recipeCache[length]; ok {
		return r, length
	}

	var r recipe

	// Choose algorithm based on length
	switch length {
	case 0, 1:
		r = recipeDft
	case 2:
		r = recipeButterfly2
	case 3:
		r = recipeButterfly3
	case 4:
		r = recipeButterfly4
	case 8:
		r = recipeButterfly8
	case 16:
		r = recipeButterfly16
	case 32:
		r = recipeButterfly32
	default:
		// Check if it's a power of two
		if isPowerOfTwo(length) && length > 32 {
			r = recipeRadix4
		} else {
			// Fall back to DFT for now
			// TODO: implement more algorithms (MixedRadix, Raders, Bluesteins, etc.)
			r = recipeDft
		}
	}

	p.recipeCache[length] = &r
	return &r, length
}

// buildFft constructs an FFT instance from a recipe
func (p *Planner) buildFft(recipe *recipe, length int, direction Direction) Fft {
	switch *recipe {
	case recipeDft:
		return &fftAdapter{inner: algorithm.NewDft(length, toAlgoDirection(direction))}
	case recipeButterfly2:
		return &fftAdapter{inner: algorithm.NewButterfly2(toAlgoDirection(direction))}
	case recipeButterfly3:
		return &fftAdapter{inner: algorithm.NewButterfly3(toAlgoDirection(direction))}
	case recipeButterfly4:
		return &fftAdapter{inner: algorithm.NewButterfly4(toAlgoDirection(direction))}
	case recipeButterfly8:
		return &fftAdapter{inner: algorithm.NewButterfly8(toAlgoDirection(direction))}
	case recipeButterfly16:
		return &fftAdapter{inner: algorithm.NewButterfly16(toAlgoDirection(direction))}
	case recipeButterfly32:
		return &fftAdapter{inner: algorithm.NewButterfly32(toAlgoDirection(direction))}
	case recipeRadix4:
		return &fftAdapter{inner: algorithm.NewRadix4(length, toAlgoDirection(direction))}
	default:
		panic("unknown recipe type")
	}
}

// toAlgoDirection converts gofft.Direction to algorithm.Direction
func toAlgoDirection(d Direction) algorithm.Direction {
	if d == Forward {
		return algorithm.Forward
	}
	return algorithm.Inverse
}

// fftAdapter adapts algorithm FFTs to the gofft.Fft interface
type fftAdapter struct {
	inner algorithm.FftInterface
}

func (f *fftAdapter) Process(buffer []complex128) {
	f.inner.ProcessWithScratch(buffer, make([]complex128, f.inner.InplaceScratchLen()))
}

func (f *fftAdapter) ProcessWithScratch(buffer, scratch []complex128) {
	f.inner.ProcessWithScratch(buffer, scratch)
}

func (f *fftAdapter) ProcessOutOfPlace(input, output, scratch []complex128) {
	// For now, do it via copy
	copy(output, input)
	f.inner.ProcessWithScratch(output, scratch)
}

func (f *fftAdapter) ProcessImmutable(input []complex128, output, scratch []complex128) {
	copy(output, input)
	f.inner.ProcessWithScratch(output, scratch)
}

func (f *fftAdapter) Len() int {
	return f.inner.Len()
}

func (f *fftAdapter) Direction() Direction {
	if f.inner.Direction() == algorithm.Forward {
		return Forward
	}
	return Inverse
}

func (f *fftAdapter) InplaceScratchLen() int {
	return f.inner.InplaceScratchLen()
}

func (f *fftAdapter) OutOfPlaceScratchLen() int {
	return 0
}

func (f *fftAdapter) ImmutableScratchLen() int {
	return f.inner.InplaceScratchLen()
}

// isPowerOfTwo checks if n is a power of two
func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// Plan32 is a planner for complex64 FFTs
type Planner32 struct {
	mu          sync.Mutex
	cache       map[plannerKey]Fft32
	recipeCache map[int]*recipe
}

// NewPlanner32 creates a new FFT planner for complex64
func NewPlanner32() *Planner32 {
	return &Planner32{
		cache:       make(map[plannerKey]Fft32),
		recipeCache: make(map[int]*recipe),
	}
}

// PlanForward creates an FFT instance for computing forward FFTs of the given size
func (p *Planner32) PlanForward(length int) Fft32 {
	return p.Plan(length, Forward)
}

// PlanInverse creates an FFT instance for computing inverse FFTs of the given size
func (p *Planner32) PlanInverse(length int) Fft32 {
	return p.Plan(length, Inverse)
}

// Plan creates an FFT instance for the given size and direction
func (p *Planner32) Plan(length int, direction Direction) Fft32 {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := plannerKey{length: length, direction: direction}

	// Check cache
	if fft, ok := p.cache[key]; ok {
		return fft
	}

	// For now, use DFT for all sizes
	// TODO: port all algorithms to complex64
	fft := &fftAdapter32{inner: algorithm.NewDft32(length, toAlgoDirection(direction))}

	// Cache it
	p.cache[key] = fft

	return fft
}

// fftAdapter32 adapts algorithm FFTs to the gofft.Fft32 interface
type fftAdapter32 struct {
	inner interface {
		Len() int
		Direction() algorithm.Direction
		InplaceScratchLen() int
		ProcessWithScratch(buffer, scratch []complex64)
	}
}

func (f *fftAdapter32) Process(buffer []complex64) {
	f.inner.ProcessWithScratch(buffer, make([]complex64, f.inner.InplaceScratchLen()))
}

func (f *fftAdapter32) ProcessWithScratch(buffer, scratch []complex64) {
	f.inner.ProcessWithScratch(buffer, scratch)
}

func (f *fftAdapter32) ProcessOutOfPlace(input, output, scratch []complex64) {
	copy(output, input)
	f.inner.ProcessWithScratch(output, scratch)
}

func (f *fftAdapter32) ProcessImmutable(input []complex64, output, scratch []complex64) {
	copy(output, input)
	f.inner.ProcessWithScratch(output, scratch)
}

func (f *fftAdapter32) Len() int {
	return f.inner.Len()
}

func (f *fftAdapter32) Direction() Direction {
	if f.inner.Direction() == algorithm.Forward {
		return Forward
	}
	return Inverse
}

func (f *fftAdapter32) InplaceScratchLen() int {
	return f.inner.InplaceScratchLen()
}

func (f *fftAdapter32) OutOfPlaceScratchLen() int {
	return 0
}

func (f *fftAdapter32) ImmutableScratchLen() int {
	return f.inner.InplaceScratchLen()
}
