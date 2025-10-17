package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
	"github.com/example/gofft/pkg/gofft"
)

func main() {
	sz := flag.Int("n", 16384, "FFT size (power of two)")
	w := flag.Int("workers", 0, "#workers (0 = GOMAXPROCS)")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	x := make([]complex128, *sz)
	for i := range x { x[i] = complex(rand.NormFloat64(), rand.NormFloat64()) }

	start := time.Now()
	X, err := gofft.FFTParallel(x, *w)
	if err != nil { panic(err) }
	dur := time.Since(start)

	fmt.Printf("N=%d workers=%d time=%s first bin=% .6f%+.6fi\n", *sz, *w, dur, real(X[0]), imag(X[0]))
}
