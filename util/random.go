package util

import (
	"crypto/rand"
)

// GenerateRandomBytes using a Cryptographically Secure Pseudorandom Number Generator (CSPRNG)
func GenerateRandomBytes(x int) ([]byte, error) {
	y := make([]byte, x)
	_, err := rand.Read(y[:])
	if err != nil {
		return nil, err
	}
	return y, err
}
