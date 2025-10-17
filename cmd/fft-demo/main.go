package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
	"github.com/example/gofft/pkg/gofft"
)

func main() {
	sz := flag.Int("n", 1024, "FFT size (power of two)")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	x := make([]complex128, *sz)
	for i := range x { x[i] = complex(rand.NormFloat64(), rand.NormFloat64()) }

	X, err := gofft.FFT(x)
	if err != nil { panic(err) }

	fmt.Printf("N=%d\n", *sz)
	for i := 0; i < min(8, len(X)); i++ {
		fmt.Printf("X[%d] = % .6f%+.6fi\n", i, real(X[i]), imag(X[i]))
	}
}

func min(a, b int) int { if a < b { return a }; return b }
