package toolkit

import (
	"crypto/rand"
	"math/big"
)

const randomStringSource = "abcdesfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

//Tools is the type used to instantiate this module. ANy variable of this type will have access to all the mehtods
//with the receiver *Tools

type Tools struct {
	MaxFileSize      int
	AllowedFileTypes []string
}

//This function generates a random string of length `n`.
//The randomness is derived from cryptographic source(`rand.Reader`) and simple prime number calculations.

func (t *Tools) RandomStringUsingRandPrime(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)

	for i := range s {
		// It generates a cryptographic strong prime number `p` with `len(r)` bits

		p, _ := rand.Prime(rand.Reader, len(r))

		x, y := p.Uint64(), uint64(len(r))

		s[i] = r[x%y]
	}
	return string(s)
}

func (t *Tools) RandomStringUsingRandInt(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)

	for i := range s {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(r))))
		s[i] = r[num.Int64()]
	}

	return string(s)
}
