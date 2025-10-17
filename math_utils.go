package gofft

import (
	"math"
)

// PrimeFactor represents a prime factor with its count
type PrimeFactor struct {
	Value int
	Count int
}

// PrimeFactors holds the prime factorization of a number
type PrimeFactors struct {
	powerOfTwo   int
	powerOfThree int
	otherFactors []PrimeFactor
	product      int
}

// ComputePrimeFactors computes the prime factorization of n
func ComputePrimeFactors(n int) PrimeFactors {
	if n <= 1 {
		return PrimeFactors{
			powerOfTwo:   0,
			powerOfThree: 0,
			otherFactors: nil,
			product:      n,
		}
	}

	original := n
	powerOfTwo := 0
	powerOfThree := 0
	var otherFactors []PrimeFactor

	// Extract powers of 2
	for n%2 == 0 {
		powerOfTwo++
		n /= 2
	}

	// Extract powers of 3
	for n%3 == 0 {
		powerOfThree++
		n /= 3
	}

	// Extract other prime factors
	factor := 5
	for factor*factor <= n {
		count := 0
		for n%factor == 0 {
			count++
			n /= factor
		}
		if count > 0 {
			otherFactors = append(otherFactors, PrimeFactor{Value: factor, Count: count})
		}

		// Skip even numbers (we already handled 2)
		// Use 2, 4, 2, 4, ... pattern to check 5, 7, 11, 13, 17, 19, ...
		if factor == 5 {
			factor += 2
		} else {
			factor += 2
		}
	}

	// If n > 1, then it's a prime factor
	if n > 1 {
		otherFactors = append(otherFactors, PrimeFactor{Value: n, Count: 1})
	}

	return PrimeFactors{
		powerOfTwo:   powerOfTwo,
		powerOfThree: powerOfThree,
		otherFactors: otherFactors,
		product:      original,
	}
}

// IsPrime returns true if this factorization represents a prime number
func (pf PrimeFactors) IsPrime() bool {
	if pf.product <= 1 {
		return false
	}
	if pf.product == 2 || pf.product == 3 {
		return true
	}
	return pf.powerOfTwo == 0 && pf.powerOfThree == 0 && len(pf.otherFactors) == 1 && pf.otherFactors[0].Count == 1
}

// GetPowerOfTwo returns the power of 2 in the factorization
func (pf PrimeFactors) GetPowerOfTwo() int {
	return pf.powerOfTwo
}

// GetPowerOfThree returns the power of 3 in the factorization
func (pf PrimeFactors) GetPowerOfThree() int {
	return pf.powerOfThree
}

// GetOtherFactors returns the prime factors other than 2 and 3
func (pf PrimeFactors) GetOtherFactors() []PrimeFactor {
	return pf.otherFactors
}

// GetProduct returns the original number
func (pf PrimeFactors) GetProduct() int {
	return pf.product
}

// HasFactorsLeq returns true if all prime factors are <= maxFactor
func (pf PrimeFactors) HasFactorsLeq(maxFactor int) bool {
	for _, f := range pf.otherFactors {
		if f.Value > maxFactor {
			return false
		}
	}
	return true
}

// HasFactorsGt returns true if any prime factor is > threshold
func (pf PrimeFactors) HasFactorsGt(threshold int) bool {
	for _, f := range pf.otherFactors {
		if f.Value > threshold {
			return true
		}
	}
	return false
}

// ProductAbove returns the product of all prime factors > threshold
func (pf PrimeFactors) ProductAbove(threshold int) int {
	result := 1
	for _, f := range pf.otherFactors {
		if f.Value > threshold {
			for i := 0; i < f.Count; i++ {
				result *= f.Value
			}
		}
	}
	return result
}

// PartitionFactors splits the prime factors into two coprime groups
// Returns (left_factors, right_factors) where their product equals the original number
func (pf PrimeFactors) PartitionFactors() (PrimeFactors, PrimeFactors) {
	// Strategy: try to split as evenly as possible
	// Put larger factors in the left, smaller in the right
	// This is a simplified version - RustFFT has more sophisticated logic

	leftProduct := 1
	rightProduct := 1

	// Split power of 2
	halfPow2 := pf.powerOfTwo / 2
	leftPow2 := halfPow2
	rightPow2 := pf.powerOfTwo - halfPow2

	leftProduct *= 1 << leftPow2
	rightProduct *= 1 << rightPow2

	// Split power of 3
	halfPow3 := pf.powerOfThree / 2
	leftPow3 := halfPow3
	rightPow3 := pf.powerOfThree - halfPow3

	for i := 0; i < leftPow3; i++ {
		leftProduct *= 3
	}
	for i := 0; i < rightPow3; i++ {
		rightProduct *= 3
	}

	// Split other factors by alternating
	var leftOther, rightOther []PrimeFactor
	addToLeft := leftProduct < rightProduct

	for _, factor := range pf.otherFactors {
		if addToLeft {
			leftOther = append(leftOther, factor)
			for i := 0; i < factor.Count; i++ {
				leftProduct *= factor.Value
			}
		} else {
			rightOther = append(rightOther, factor)
			for i := 0; i < factor.Count; i++ {
				rightProduct *= factor.Value
			}
		}
		addToLeft = leftProduct < rightProduct
	}

	return PrimeFactors{
			powerOfTwo:   leftPow2,
			powerOfThree: leftPow3,
			otherFactors: leftOther,
			product:      leftProduct,
		}, PrimeFactors{
			powerOfTwo:   rightPow2,
			powerOfThree: rightPow3,
			otherFactors: rightOther,
			product:      rightProduct,
		}
}

// GCD computes the greatest common divisor of two numbers
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// IsPrimeSimple performs a simple primality test
func IsPrimeSimple(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	i := 5
	for i*i <= n {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
		i += 6
	}
	return true
}

// NextPowerOfTwo returns the next power of two >= n
func NextPowerOfTwo(n int) int {
	if n <= 1 {
		return 1
	}
	// Find the highest set bit
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return n
}

// TrailingZeros counts the number of trailing zero bits
func TrailingZeros(n int) int {
	if n == 0 {
		return 64 // or word size
	}
	count := 0
	for n&1 == 0 {
		count++
		n >>= 1
	}
	return count
}

// TwiddleFactor computes the twiddle factor e^(-2πik/n) for forward FFT
// or e^(2πik/n) for inverse FFT
func TwiddleFactor(k, n int, direction Direction) complex128 {
	angle := 2.0 * math.Pi * float64(k) / float64(n)
	if direction == Forward {
		angle = -angle
	}
	return complex(math.Cos(angle), math.Sin(angle))
}

// TwiddleFactor32 computes the twiddle factor for float32
func TwiddleFactor32(k, n int, direction Direction) complex64 {
	angle := 2.0 * math.Pi * float64(k) / float64(n)
	if direction == Forward {
		angle = -angle
	}
	return complex(float32(math.Cos(angle)), float32(math.Sin(angle)))
}

// ComputeTwiddles precomputes all twiddle factors for a given FFT size
func ComputeTwiddles(n int, direction Direction) []complex128 {
	twiddles := make([]complex128, n)
	for k := 0; k < n; k++ {
		twiddles[k] = TwiddleFactor(k, n, direction)
	}
	return twiddles
}

// ComputeTwiddles32 precomputes all twiddle factors for a given FFT size (float32)
func ComputeTwiddles32(n int, direction Direction) []complex64 {
	twiddles := make([]complex64, n)
	for k := 0; k < n; k++ {
		twiddles[k] = TwiddleFactor32(k, n, direction)
	}
	return twiddles
}
