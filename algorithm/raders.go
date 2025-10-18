package algorithm

import "math"

// Raders implements Rader's Algorithm for prime-sized FFTs
// This converts a prime-size p FFT into a size p-1 FFT via convolution
// More efficient than Bluestein's for primes
type Raders struct {
	length               int
	direction            Direction
	innerFft             FftInterface // FFT of size p-1
	innerFftData         []complex128 // Precomputed FFT for convolution
	primitiveRoot        int
	primitiveRootInv     int
	inplaceScratchLen    int
	outofplaceScratchLen int
}

// NewRaders creates a Rader's algorithm instance for a prime size
// innerFft must have length p-1 where p is prime
func NewRaders(innerFft FftInterface) *Raders {
	innerLen := innerFft.Len()
	length := innerLen + 1

	// Verify p is prime
	if !isPrime(length) {
		panic("Rader's algorithm requires prime size")
	}

	direction := innerFft.Direction()

	// Find primitive root g mod p
	g := findPrimitiveRoot(length)

	// Find modular inverse of g mod p
	gInv := modInverse(g, length)

	// Precompute the convolution kernel
	// h[k] = (1/innerLen) * W^(g^(-k)) where W = exp(-2πi/p)
	// Use primitive_root_inverse to iterate
	scale := 1.0 / float64(innerLen)
	innerFftInput := make([]complex128, innerLen)
	twiddleIdx := 1
	for i := range innerFftInput {
		angle := -2.0 * math.Pi * float64(twiddleIdx) / float64(length)
		if direction == Inverse {
			angle = -angle
		}
		twiddle := complex(math.Cos(angle), math.Sin(angle))
		innerFftInput[i] = twiddle * complex(scale, 0)

		// Use inverse root to iterate
		twiddleIdx = (twiddleIdx * gInv) % length
	}

	// FFT the kernel for convolution
	innerFftData := make([]complex128, innerLen)
	copy(innerFftData, innerFftInput)
	innerScratch := make([]complex128, innerFft.InplaceScratchLen())
	innerFft.ProcessWithScratch(innerFftData, innerScratch)

	// Calculate scratch requirements
	inplaceScratch := innerLen + innerFft.InplaceScratchLen()
	outofplaceScratch := innerFft.InplaceScratchLen()

	return &Raders{
		length:               length,
		direction:            direction,
		innerFft:             innerFft,
		innerFftData:         innerFftData,
		primitiveRoot:        g,
		primitiveRootInv:     gInv,
		inplaceScratchLen:    inplaceScratch,
		outofplaceScratchLen: outofplaceScratch,
	}
}

func (r *Raders) Len() int                  { return r.length }
func (r *Raders) Direction() Direction      { return r.direction }
func (r *Raders) InplaceScratchLen() int    { return r.inplaceScratchLen }
func (r *Raders) OutOfPlaceScratchLen() int { return r.outofplaceScratchLen }

func (r *Raders) ProcessWithScratch(buffer, scratch []complex128) {
	// Process each chunk of size r.length
	for i := 0; i < len(buffer); i += r.length {
		chunk := buffer[i : i+r.length]
		workScratch := scratch[:r.inplaceScratchLen]
		r.processOne(chunk, workScratch)
	}
}

func (r *Raders) ProcessOutOfPlace(input, output, scratch []complex128) {
	copy(output, input)
	r.ProcessWithScratch(output, scratch)
}

func (r *Raders) processOne(buffer, scratch []complex128) {
	innerLen := r.length - 1
	innerScratch := scratch[:innerLen]
	extraScratch := scratch[innerLen:]

	// Save first element
	first := buffer[0]

	// Reorder buffer[1:] into scratch using primitive root
	// scratch[k] = buffer[g^k mod p]
	idx := 1
	for i := 0; i < innerLen; i++ {
		idx = (idx * r.primitiveRoot) % r.length
		innerScratch[i] = buffer[idx]
	}

	// First inner FFT
	r.innerFft.ProcessWithScratch(innerScratch, extraScratch)

	// innerScratch[0] is sum of buffer[1:], add buffer[0] for DC component
	buffer[0] = first + innerScratch[0]

	// Multiply with precomputed data and conjugate (sets up for inverse FFT)
	for i := range innerScratch {
		innerScratch[i] = complexConj(innerScratch[i] * r.innerFftData[i])
	}

	// Add first element (conjugated) to DC bin
	innerScratch[0] = innerScratch[0] + complexConj(first)

	// Second FFT (effectively inverse due to conjugation)
	r.innerFft.ProcessWithScratch(innerScratch, extraScratch)

	// Reorder output using inverse primitive root
	// buffer[g^(-k) mod p] = conj(scratch[k])
	idx = 1
	for i := 0; i < innerLen; i++ {
		idx = (idx * r.primitiveRootInv) % r.length
		buffer[idx] = complexConj(innerScratch[i])
	}
}

// Helper functions for Rader's algorithm

// isPrime checks if n is prime using trial division
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// findPrimitiveRoot finds a primitive root modulo p (p must be prime)
// A primitive root g generates all non-zero elements mod p
func findPrimitiveRoot(p int) int {
	if p == 2 {
		return 1
	}

	// Factor p-1
	phi := p - 1
	factors := factorize(phi)

	// Try each candidate
	for g := 2; g < p; g++ {
		if isPrimitiveRoot(g, p, phi, factors) {
			return g
		}
	}

	panic("No primitive root found (should not happen for prime p)")
}

// isPrimitiveRoot checks if g is a primitive root mod p
func isPrimitiveRoot(g, p, phi int, factors []int) bool {
	// For each prime factor q of phi, check if g^(phi/q) ≠ 1 (mod p)
	for _, q := range factors {
		exp := phi / q
		if modPow(g, exp, p) == 1 {
			return false
		}
	}
	return true
}

// factorize returns unique prime factors of n
func factorize(n int) []int {
	factors := []int{}

	// Check for 2
	if n%2 == 0 {
		factors = append(factors, 2)
		for n%2 == 0 {
			n /= 2
		}
	}

	// Check odd factors
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			factors = append(factors, i)
			for n%i == 0 {
				n /= i
			}
		}
	}

	// If n > 1, it's prime
	if n > 1 {
		factors = append(factors, n)
	}

	return factors
}

// modPow computes (base^exp) mod m
func modPow(base, exp, m int) int {
	result := 1
	base = base % m

	for exp > 0 {
		if exp%2 == 1 {
			result = (result * base) % m
		}
		exp = exp >> 1
		base = (base * base) % m
	}

	return result
}

// modInverse computes modular inverse of a mod m using extended Euclidean algorithm
func modInverse(a, m int) int {
	a = a % m
	x, _ := extendedGCD(a, m)

	// Make sure result is positive
	result := x % m
	if result < 0 {
		result += m
	}

	return result
}

// extendedGCD returns x, y such that a*x + m*y = gcd(a, m)
func extendedGCD(a, m int) (int, int) {
	if a == 0 {
		return 0, 1
	}

	x1, y1 := extendedGCD(m%a, a)
	x := y1 - (m/a)*x1
	y := x1

	return x, y
}
