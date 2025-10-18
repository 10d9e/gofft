package gofft

import (
	"runtime"
	"sync"

	"github.com/10d9e/gofft/algorithm"
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
// Automatically uses SIMD optimizations when available
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

	// Build the FFT from the recipe (will use SIMD if available)
	fft := p.buildFft(recipe, len, direction)

	// Cache it
	p.cache[key] = fft

	return fft
}

// isSIMDAvailable checks if SIMD is available for the current architecture
func isSIMDAvailable() bool {
	switch runtime.GOARCH {
	case "arm64":
		return true // ARM64 always has NEON
	case "amd64":
		return false // TODO: Implement x86_64 SIMD detection
	default:
		return false
	}
}

// recipe describes how to construct an FFT without actually building it
type recipe int

const (
	recipeDft recipe = iota
	recipeButterfly2
	recipeButterfly3
	recipeButterfly4
	recipeButterfly5
	recipeButterfly6
	recipeButterfly7
	recipeButterfly8
	recipeButterfly9
	recipeButterfly11
	recipeButterfly12
	recipeButterfly13
	recipeButterfly16
	recipeButterfly17
	recipeButterfly19
	recipeButterfly23
	recipeButterfly24
	recipeButterfly27
	recipeButterfly29
	recipeButterfly31
	recipeButterfly32
	recipeRadix4
	recipeRadixN
	recipeRaders
	recipeBluestein
	// Additional NEON-supported algorithms
	recipeButterfly1
	recipeButterfly10
	recipeButterfly15
	recipeMixedRadix
	recipeGoodThomas
	recipeWinograd
)

// designFft creates a recipe for an FFT of the given length
func (p *Planner) designFft(length int) (*recipe, int) {
	if r, ok := p.recipeCache[length]; ok {
		return r, length
	}

	var r recipe

	// Choose algorithm based on length
	switch length {
	case 0:
		r = recipeDft
	case 1:
		r = recipeButterfly1
	case 2:
		r = recipeButterfly2
	case 3:
		r = recipeButterfly3
	case 4:
		r = recipeButterfly4
	case 5:
		r = recipeButterfly5
	case 6:
		r = recipeButterfly6
	case 7:
		r = recipeButterfly7
	case 8:
		r = recipeButterfly8
	case 9:
		r = recipeButterfly9
	case 10:
		r = recipeButterfly10
	case 11:
		r = recipeButterfly11
	case 12:
		r = recipeButterfly12
	case 13:
		r = recipeButterfly13
	case 15:
		r = recipeButterfly15
	case 16:
		r = recipeButterfly16
	case 17:
		r = recipeButterfly17
	case 19:
		r = recipeButterfly19
	case 23:
		r = recipeButterfly23
	case 24:
		r = recipeButterfly24
	case 27:
		r = recipeButterfly27
	case 29:
		r = recipeButterfly29
	case 31:
		r = recipeButterfly31
	case 32:
		r = recipeButterfly32
	// NEON-supported composite sizes
	case 35:
		r = recipeGoodThomas
	case 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97:
		r = recipeRaders
	case 49:
		r = recipeWinograd
	case 60:
		r = recipeMixedRadix
	case 77:
		r = recipeGoodThomas
	case 120:
		r = recipeMixedRadix
	case 121:
		r = recipeWinograd
	case 143:
		r = recipeGoodThomas
	case 169:
		r = recipeWinograd
	case 221:
		r = recipeGoodThomas
	case 240:
		r = recipeMixedRadix
	case 289:
		r = recipeWinograd
	case 480:
		r = recipeMixedRadix
	default:
		// Check if it's a power of two
		if isPowerOfTwo(length) && length > 32 {
			r = recipeRadix4
		} else {
			// Try RadixN for composite sizes (factors 2-7 only)
			if canUseRadixN(length) {
				r = recipeRadixN
			} else {
				// For non-factorizable sizes, check if it's prime
				factors := ComputePrimeFactors(length)
				if factors.IsPrime() && length <= 97 {
					// Use Rader's for small/medium primes (more efficient than Bluestein's)
					r = recipeRaders
				} else {
					// Use Bluestein's for large primes or sizes with factors >7
					r = recipeBluestein
				}
			}
		}
	}

	p.recipeCache[length] = &r
	return &r, length
}

// buildFft constructs an FFT instance from a recipe
func (p *Planner) buildFft(recipe *recipe, length int, direction Direction) Fft {
	dir := toAlgoDirection(direction)
	switch *recipe {
	case recipeDft:
		return &fftAdapter{inner: algorithm.NewDft(length, dir)}
	case recipeButterfly2:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 2, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly2(dir)}
	case recipeButterfly3:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 3, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly3(dir)}
	case recipeButterfly4:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 4, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly4(dir)}
	case recipeButterfly5:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 5, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly5(dir)}
	case recipeButterfly6:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 6, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly6(dir)}
	case recipeButterfly7:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 7, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly7(dir)}
	case recipeButterfly8:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 8, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly8(dir)}
	case recipeButterfly9:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 9, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly9(dir)}
	case recipeButterfly11:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 11, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly11(dir)}
	case recipeButterfly12:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 12, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly12(dir)}
	case recipeButterfly13:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 13, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly13(dir)}
	case recipeButterfly16:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 16, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly16(dir)}
	case recipeButterfly17:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 17, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly17(dir)}
	case recipeButterfly19:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 19, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly19(dir)}
	case recipeButterfly23:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 23, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly23(dir)}
	case recipeButterfly24:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 24, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly24(dir)}
	case recipeButterfly27:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 27, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly27(dir)}
	case recipeButterfly29:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 29, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly29(dir)}
	case recipeButterfly31:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 31, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly31(dir)}
	case recipeButterfly32:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 32, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewButterfly32(dir)}
	case recipeButterfly1:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 1, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewDft(1, dir)}
	case recipeButterfly10:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 10, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewDft(10, dir)}
	case recipeButterfly15:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: 15, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewDft(15, dir)}
	case recipeMixedRadix:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: length, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewDft(length, dir)}
	case recipeGoodThomas:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: length, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewDft(length, dir)}
	case recipeWinograd:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: length, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewDft(length, dir)}
	case recipeRadix4:
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: length, direction: direction}
		}
		return &fftAdapter{inner: algorithm.NewRadix4(length, dir)}
	case recipeRadixN:
		// For NEON-supported RadixN sizes, use NEON implementation
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: length, direction: direction}
		}
		// Factor the length and create RadixN for scalar implementation
		factors := factorizeForRadixN(length)
		baseFft := algorithm.NewDft(1, dir) // Base of size 1
		return &fftAdapter{inner: algorithm.NewRadixN(factors, baseFft)}
	case recipeRaders:
		// For NEON-supported Rader's sizes, use NEON implementation
		if isSIMDAvailable() {
			return &neonButterflyAdapter{length: length, direction: direction}
		}
		// Create FFT of size length-1 for Rader's algorithm
		// Recursively plan the inner FFT
		innerLength := length - 1
		innerRecipe, _ := p.designFft(innerLength)
		innerFftImpl := p.buildFft(innerRecipe, innerLength, direction)
		// Handle both fftAdapter and neonButterflyAdapter
		var innerAlgo algorithm.FftInterface
		if adapter, ok := innerFftImpl.(*fftAdapter); ok {
			innerAlgo = adapter.inner
		} else {
			// For NEON adapters, create a scalar fallback
			innerAlgo = algorithm.NewDft(innerLength, toAlgoDirection(direction))
		}
		return &fftAdapter{inner: algorithm.NewRaders(innerAlgo)}
	case recipeBluestein:
		return &fftAdapter{inner: algorithm.NewBluestein(length, dir)}
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

// neonButterflyAdapter provides NEON-optimized butterfly implementations
type neonButterflyAdapter struct {
	length    int
	direction Direction
}

func (a *neonButterflyAdapter) Len() int {
	return a.length
}

func (a *neonButterflyAdapter) Direction() Direction {
	return a.direction
}

func (a *neonButterflyAdapter) Process(buffer []complex128) {
	// Use NEON implementations when available
	// This calls the architecture-specific implementation
	a.processNEONARM64(buffer)
}

func (a *neonButterflyAdapter) processNEON(buffer []complex128) {
	// This will be implemented with build tags to call actual NEON functions
	// For now, fall back to scalar implementation
	scalarFft := algorithm.NewDft(a.length, toAlgoDirection(a.direction))
	scratch := make([]complex128, scalarFft.InplaceScratchLen())
	scalarFft.ProcessWithScratch(buffer, scratch)
}

func (a *neonButterflyAdapter) ImmutableScratchLen() int {
	return 0
}

func (a *neonButterflyAdapter) InplaceScratchLen() int {
	return 0
}

func (a *neonButterflyAdapter) OutOfPlaceScratchLen() int {
	return 0
}

func (a *neonButterflyAdapter) ProcessImmutable(input []complex128, output []complex128, scratch []complex128) {
	copy(output, input)
	a.Process(output)
}

func (a *neonButterflyAdapter) ProcessOutOfPlace(input []complex128, output []complex128, scratch []complex128) {
	copy(output, input)
	a.Process(output)
}

func (a *neonButterflyAdapter) ProcessWithScratch(buffer []complex128, scratch []complex128) {
	a.Process(buffer)
}
