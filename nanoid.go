package nanoid

import (
	"crypto/rand"
	"fmt"
	"math/bits"
)

type NanoID []rune

func New(options ...func(*Option) *Option) (NanoID, error) {
	opt := new(Option)
	for _, o := range options {
		o(opt)
	}

	if opt.length < 2 || opt.length > 255 {
		opt.length = CanonicNanoIDLenght
	}

	if opt.alphabet == nil {
		opt.alphabet = CanonicAlphabet
	}

	alphabetLen := len(opt.alphabet)
	// Runes to support unicode.
	runes := opt.alphabet

	// Because the custom alphabet is not guaranteed to have
	// 64 chars to utilise, we have to calculate a suitable mask.
	x := uint32(alphabetLen) - 1
	clz := bits.LeadingZeros32(x | 1)
	mask := (2 << (31 - clz)) - 1
	step := (opt.length / 5) * 8

	b := make([]byte, step)
	id := make([]rune, opt.length)

	j, idx := 0, 0
Outer:
	for {
		_, err := rand.Read(b)
		if err != nil {
			return nil, fmt.Errorf("nanoid: %s caused by (%w)", err.Error(), ErrInvalidByte)
		}
		for i := 0; i < step; i++ {
			idx = int(b[i]) & mask
			if idx < alphabetLen {
				id[j] = runes[idx]
				j++
				if j == opt.length {
					j = 0
					break Outer
				}
			}
		}
	}

	return NanoID(id), nil
}

func (n *NanoID) String() string {
	if n == nil {
		return "<nil>"
	}
	return string(*n)
}
