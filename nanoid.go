package nanoid

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/bits"
)

type NanoID []rune

func New(options ...func(*Option) *Option) (NanoID, error) {
	opt := new(Option)
	for _, o := range options {
		o(opt)
	}
	// TODO: change defaults injection and validations
	if opt.length < 2 || opt.length > 255 {
		opt.length = CanonicNanoIDLenght
	}

	if opt.alphabet == nil {
		opt.alphabet = CanonicAlphabet
	}

	alphabetLen := len(opt.alphabet)
	// Runes to support unicode.
	runes := opt.alphabet

	// First, a bitmask is necessary to generate the ID. The bitmask makes bytes
	// values closer to the alphabet size. The bitmask calculates the closest
	// `2^31 - 1` number, which exceeds the alphabet size.
	// For example, the bitmask for the alphabet size 30 is 31 (00011111).
	alphabetUpperBound := uint32(alphabetLen) - 1
	clz := bits.LeadingZeros32(alphabetUpperBound | 1)
	mask := (2 << (31 - clz)) - 1

	// Though, the bitmask solution is not perfect since the bytes exceeding
	// the alphabet size are refused. Therefore, to reliably generate the ID,
	// the random bytes redundancy has to be satisfied.

	// Note: every hardware random generator call is performance expensive,
	// because the system call for entropy collection takes a lot of time.
	// So, to avoid additional system calls, extra bytes are requested in advance.

	// Next, a step determines how many random bytes to generate.
	// The number of random bytes gets decided upon the ID size, mask,
	// alphabet size, and magic number 1.6 (using 1.6 peaks at performance
	// according to benchmarks).
	//TODO: add readability to this step
	step := int(math.Ceil(1.6 * float64(mask*opt.length) / float64(alphabetLen)))

	randBytes := make([]byte, step)
	id := make([]rune, opt.length)

	// TODO: make this a function to avoid return with tag
	currentRuneToGeneratePos, randomAlphabetPos := 0, 0
Outer:
	for {
		_, err := rand.Read(randBytes)
		if err != nil {
			return nil, fmt.Errorf("nanoid: %s caused by (%w)", err.Error(), ErrInvalidByte)
		}

		for i := 0; i < step; i++ {
			randomAlphabetPos = int(randBytes[i]) & mask

			if randomAlphabetPos >= alphabetLen {
				continue
			}

			id[currentRuneToGeneratePos] = runes[randomAlphabetPos]

			currentRuneToGeneratePos++

			if currentRuneToGeneratePos == opt.length {
				break Outer
			}
		}
	}

	return NanoID(id), nil
}

func (n NanoID) String() string {
	if n == nil {
		return "<nil>"
	}
	return string(n)
}
