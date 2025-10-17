package algorithm

import (
	"testing"
)

// TestCompareAlgorithmWithRust documents what Rust does vs what we do
func TestCompareAlgorithmWithRust(t *testing.T) {
	// For size 32:
	// - Rust uses Butterfly32 as base (2^5 = 32), so exponent=5
	// - k = (5-5)/2 = 0
	// - So NO radix-4 cross-FFT stages!
	//
	// Our Go implementation:
	// - We use Butterfly8 as base (2^3 = 8), so exponent=5, base_exponent=3
	// - k = (5-3)/2 = 1
	// - So we have ONE radix-4 cross-FFT stage
	//
	// This is the fundamental difference!

	size := 32
	exponent := trailingZeros(size) // 5

	t.Logf("Size %d has exponent %d", size, exponent)

	// Rust's choice:
	rustBaseExponent := 5
	rustK := (exponent - rustBaseExponent) / 2
	t.Logf("Rust: baseExponent=%d, k=%d, uses Butterfly32 directly", rustBaseExponent, rustK)

	// Our choice:
	ourFft := NewRadix4(size, Forward)
	ourBaseExponent := trailingZeros(ourFft.baseLen)
	ourK := (exponent - ourBaseExponent) / 2
	t.Logf("Go: baseExponent=%d, k=%d, uses base size %d", ourBaseExponent, ourK, ourFft.baseLen)

	if rustK != ourK {
		t.Logf("")
		t.Logf("IMPORTANT: Rust and Go use different k values!")
		t.Logf("Rust directly uses Butterfly32 (k=0, no cross-FFTs)")
		t.Logf("Go uses Butterfly8 + one radix-4 stage (k=1)")
		t.Logf("")
		t.Logf("For perfect compatibility, we should use Butterfly32 for size 32")
	}
}

// TestWhatShouldSize32Use checks which base Rust uses for size 32
func TestWhatShouldSize32Use(t *testing.T) {
	// According to RustFFT radix4.rs lines 51-63:
	//     _ => {
	//         if exponent % 2 == 1 {
	//             (5, Arc::new(Butterfly32::new(direction)) as Arc<dyn Fft<T>>)
	//         } else {
	//             (4, Arc::new(Butterfly16::new(direction)) as Arc<dyn Fft<T>>)
	//         }
	//     }
	//
	// For size 32, exponent = 5 (odd), so it uses Butterfly32 with base_exponent=5
	// Therefore k = (5-5)/2 = 0

	// Our code should do the same! Let me check what exponent 5 gives us
	size := 32
	exponent := trailingZeros(size)

	t.Logf("Size=%d, exponent=%d", size, exponent)
	t.Logf("exponent mod 2 = %d", exponent%2)

	if exponent%2 == 1 {
		t.Logf("Exponent is ODD, so RustFFT uses Butterfly32 (base_exponent=5)")
		t.Logf("This gives k = (%d - 5) / 2 = %d", exponent, (exponent-5)/2)
		t.Logf("")
		t.Logf("Our code should use Butterfly32 for size 32, NOT Butterfly8!")
	}
}
