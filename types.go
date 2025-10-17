package gofft

// This file contains types that are needed by both the main package and algorithm package

// Direction represents whether an FFT is forward or inverse
type Direction int

const (
	// Forward represents a forward FFT
	Forward Direction = iota
	// Inverse represents an inverse FFT
	Inverse
)

// OppositeDirection returns the opposite direction
func (d Direction) OppositeDirection() Direction {
	if d == Forward {
		return Inverse
	}
	return Forward
}

func (d Direction) String() string {
	if d == Forward {
		return "Forward"
	}
	return "Inverse"
}
