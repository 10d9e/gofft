package gofft

import "github.com/10d9e/gofft/algorithm"

// factorizeForRadixN breaks n into small radix factors (2-7)
// Returns nil if n cannot be factored into only 2-7
func factorizeForRadixN(n int) []algorithm.RadixFactor {
	if n <= 1 {
		return nil
	}

	factors := []algorithm.RadixFactor{}

	// Factor out 7s
	for n%7 == 0 {
		factors = append(factors, algorithm.Factor7)
		n /= 7
	}

	// Factor out 5s
	for n%5 == 0 {
		factors = append(factors, algorithm.Factor5)
		n /= 5
	}

	// Factor out 3s
	for n%3 == 0 {
		factors = append(factors, algorithm.Factor3)
		n /= 3
	}

	// Factor out 2s
	for n%2 == 0 {
		factors = append(factors, algorithm.Factor2)
		n /= 2
	}

	// If n != 1, there are prime factors > 7
	if n != 1 {
		return nil
	}

	return factors
}

// canUseRadixN checks if a size can be handled by RadixN
func canUseRadixN(n int) bool {
	return factorizeForRadixN(n) != nil
}
