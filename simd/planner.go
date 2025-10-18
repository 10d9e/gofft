package simd

import (
	"github.com/10d9e/gofft"
	"github.com/10d9e/gofft/algorithm"
	"github.com/10d9e/gofft/simd/neon"
)

// SIMDPlanner creates FFT instances optimized for SIMD
type SIMDPlanner struct {
	level SIMDLevel
}

// NewSIMDPlanner creates a new SIMD planner
func NewSIMDPlanner() *SIMDPlanner {
	return &SIMDPlanner{
		level: DetectSIMD(),
	}
}

// PlanForward creates a forward FFT optimized for the detected SIMD level
func (p *SIMDPlanner) PlanForward(length int) FftInterface {
	return p.plan(length, algorithm.Forward)
}

// PlanInverse creates an inverse FFT optimized for the detected SIMD level
func (p *SIMDPlanner) PlanInverse(length int) FftInterface {
	return p.plan(length, algorithm.Inverse)
}

// plan creates an FFT instance optimized for the current SIMD level
func (p *SIMDPlanner) plan(length int, direction algorithm.Direction) FftInterface {
	switch p.level {
	case NEON:
		return p.planNEON(length, direction)
	case AVX:
		return p.planAVX(length, direction)
	case SSE:
		return p.planSSE(length, direction)
	default:
		return p.planScalar(length, direction)
	}
}

// planNEON creates a NEON-optimized FFT
func (p *SIMDPlanner) planNEON(length int, direction algorithm.Direction) FftInterface {
	// Import NEON package to avoid circular imports
	// For now, we'll create a wrapper that uses NEON when available

	// Check if we can use NEON butterflies for small sizes
	switch length {
	case 2, 3, 4, 5, 7, 8, 9, 11, 13, 16, 17, 19, 23, 27, 29, 31, 32:
		// Use NEON butterfly implementation
		return &neonButterflyAdapter{
			length:    length,
			direction: direction,
		}
	case 6, 10, 12, 15, 18, 20:
		// Use NEON RadixN implementation
		return &neonRadixNAdapter{
			length:    length,
			direction: direction,
		}
	case 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97:
		// Use NEON Rader's implementation
		return &neonRadersAdapter{
			length:    length,
			direction: direction,
		}
	case 64, 128, 256, 512, 1024:
		// Use NEON Radix-4 implementation
		return &neonRadix4Adapter{
			length:    length,
			direction: direction,
		}
	default:
		// For larger sizes, fall back to scalar for now
		// TODO: Implement more NEON algorithms
		return p.planScalar(length, direction)
	}
}

// planAVX creates an AVX-optimized FFT
func (p *SIMDPlanner) planAVX(length int, direction algorithm.Direction) FftInterface {
	// For now, fall back to scalar implementation
	// TODO: Implement actual AVX-optimized algorithms
	return p.planScalar(length, direction)
}

// planSSE creates an SSE-optimized FFT
func (p *SIMDPlanner) planSSE(length int, direction algorithm.Direction) FftInterface {
	// For now, fall back to scalar implementation
	// TODO: Implement actual SSE-optimized algorithms
	return p.planScalar(length, direction)
}

// planScalar creates a scalar FFT (fallback)
func (p *SIMDPlanner) planScalar(length int, direction algorithm.Direction) FftInterface {
	// Use the existing scalar planner
	scalarPlanner := gofft.NewPlanner()
	if direction == algorithm.Forward {
		return &fftAdapter{inner: scalarPlanner.PlanForward(length)}
	}
	return &fftAdapter{inner: scalarPlanner.PlanInverse(length)}
}

// GetSIMDLevel returns the current SIMD level
func (p *SIMDPlanner) GetSIMDLevel() SIMDLevel {
	return p.level
}

// IsSIMDOptimized returns true if the planner is using SIMD optimizations
func (p *SIMDPlanner) IsSIMDOptimized() bool {
	return p.level != Scalar
}

// FftInterface represents an FFT instance (placeholder for now)
type FftInterface interface {
	Process(buffer []complex128)
	Len() int
	Direction() algorithm.Direction
}

// fftAdapter adapts the existing FFT interface to our SIMD interface
type fftAdapter struct {
	inner gofft.Fft
}

func (a *fftAdapter) Process(buffer []complex128) {
	a.inner.Process(buffer)
}

func (a *fftAdapter) Len() int {
	return a.inner.Len()
}

func (a *fftAdapter) Direction() algorithm.Direction {
	dir := a.inner.Direction()
	if dir == gofft.Forward {
		return algorithm.Forward
	}
	return algorithm.Inverse
}

// neonButterflyAdapter adapts NEON butterflies to our SIMD interface
type neonButterflyAdapter struct {
	length    int
	direction algorithm.Direction
}

func (a *neonButterflyAdapter) Process(buffer []complex128) {
	// Use actual NEON butterfly implementations
	switch a.length {
	case 1:
		neon.Butterfly1_NEON(buffer)
	case 2:
		neon.Butterfly2_NEON(buffer)
	case 3:
		neon.Butterfly3_NEON(buffer)
	case 4:
		neon.Butterfly4_NEON(buffer)
	case 5:
		neon.Butterfly5_NEON(buffer)
	case 7:
		neon.Butterfly7_NEON(buffer)
	case 8:
		neon.Butterfly8_NEON(buffer)
	case 9:
		neon.Butterfly9_NEON(buffer)
	case 11:
		neon.Butterfly11_NEON(buffer)
	case 13:
		neon.Butterfly13_NEON(buffer)
	case 16:
		neon.Butterfly16_NEON(buffer)
	case 17:
		neon.Butterfly17_NEON(buffer)
	case 19:
		neon.Butterfly19_NEON(buffer)
	case 23:
		neon.Butterfly23_NEON(buffer)
	case 27:
		neon.Butterfly27_NEON(buffer)
	case 29:
		neon.Butterfly29_NEON(buffer)
	case 31:
		neon.Butterfly31_NEON(buffer)
	case 32:
		neon.Butterfly32_NEON(buffer)
	default:
		// Fallback to scalar for unsupported sizes
		scalarFft := algorithm.NewDft(a.length, a.direction)
		scratch := make([]complex128, scalarFft.InplaceScratchLen())
		scalarFft.ProcessWithScratch(buffer, scratch)
	}
}

func (a *neonButterflyAdapter) Len() int {
	return a.length
}

func (a *neonButterflyAdapter) Direction() algorithm.Direction {
	return a.direction
}

// neonRadix4Adapter adapts NEON Radix-4 FFTs to our SIMD interface
type neonRadix4Adapter struct {
	length    int
	direction algorithm.Direction
}

func (a *neonRadix4Adapter) Process(buffer []complex128) {
	// Use actual NEON Radix-4 implementations
	switch a.length {
	case 64:
		neon.Radix4_64_NEON(buffer)
	case 128:
		neon.Radix4_128_NEON(buffer)
	case 256:
		neon.Radix4_256_NEON(buffer)
	case 512:
		neon.Radix4_512_NEON(buffer)
	case 1024:
		neon.Radix4_1024_NEON(buffer)
	default:
		// Fallback to scalar for unsupported sizes
		scalarFft := algorithm.NewDft(a.length, a.direction)
		scratch := make([]complex128, scalarFft.InplaceScratchLen())
		scalarFft.ProcessWithScratch(buffer, scratch)
	}
}

func (a *neonRadix4Adapter) Len() int {
	return a.length
}

func (a *neonRadix4Adapter) Direction() algorithm.Direction {
	return a.direction
}

// neonRadixNAdapter adapts NEON RadixN FFTs to our SIMD interface
type neonRadixNAdapter struct {
	length    int
	direction algorithm.Direction
}

func (a *neonRadixNAdapter) Process(buffer []complex128) {
	// For now, use scalar implementation
	// TODO: Import and use actual NEON RadixN
	scalarFft := algorithm.NewDft(a.length, a.direction)
	scratch := make([]complex128, scalarFft.InplaceScratchLen())
	scalarFft.ProcessWithScratch(buffer, scratch)
}

func (a *neonRadixNAdapter) Len() int {
	return a.length
}

func (a *neonRadixNAdapter) Direction() algorithm.Direction {
	return a.direction
}

// neonRadersAdapter adapts NEON Rader's FFTs to our SIMD interface
type neonRadersAdapter struct {
	length    int
	direction algorithm.Direction
}

func (a *neonRadersAdapter) Process(buffer []complex128) {
	// For now, use scalar implementation
	// TODO: Import and use actual NEON Rader's
	scalarFft := algorithm.NewDft(a.length, a.direction)
	scratch := make([]complex128, scalarFft.InplaceScratchLen())
	scalarFft.ProcessWithScratch(buffer, scratch)
}

func (a *neonRadersAdapter) Len() int {
	return a.length
}

func (a *neonRadersAdapter) Direction() algorithm.Direction {
	return a.direction
}

// SIMDPlannerAdapter adapts the SIMD planner to work with existing code
type SIMDPlannerAdapter struct {
	*SIMDPlanner
}

// NewSIMDPlannerAdapter creates a new SIMD planner adapter
func NewSIMDPlannerAdapter() *SIMDPlannerAdapter {
	return &SIMDPlannerAdapter{
		SIMDPlanner: NewSIMDPlanner(),
	}
}

// PlanForward creates a forward FFT
func (p *SIMDPlannerAdapter) PlanForward(length int) FftInterface {
	return p.SIMDPlanner.PlanForward(length)
}

// PlanInverse creates an inverse FFT
func (p *SIMDPlannerAdapter) PlanInverse(length int) FftInterface {
	return p.SIMDPlanner.PlanInverse(length)
}
