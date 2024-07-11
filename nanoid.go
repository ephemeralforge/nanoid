package nanoid

import (
	"math"
	"math/bits"
)

type NanoID []rune

// TODO: buffers optimization
// New NanoID with options or canonnic by default.
func New(options ...func(*Option) *Option) (NanoID, error) {
	opt := new(Option)

	switch len(options) {
	case 0:
		opt.length = CanonicNanoIDLenght
		opt.alphabet = AlphabetFromString(CanonicNanoIDAlphabet)
		opt.randomFunc = CanonicNanoIDRandomFunc
	default:
		for _, o := range options {
			o(opt)
		}
	}

	if opt.length < 2 || opt.length > 255 {
		return nil, ErrInvalidIDLength
	}

	if opt.alphabet == nil {
		return nil, ErrNilAlphabet
	}

	// First, a bitmask is necessary to generate the ID. The bitmask makes bytes
	// values closer to the alphabet size. The bitmask calculates the closest

	alphabetLen := len(opt.alphabet)
	alphabetUpperBound := uint32(alphabetLen) - 1
	clz := bits.LeadingZeros32(alphabetUpperBound | 1)
	mask := (2 << (31 - clz)) - 1 //nolint:mnd // `2^31 - 1` number, which exceeds the alphabet size
	// For example, the bitmask for the alphabet size 30 is 31 (00011111).

	// TODO: lurk for another approach
	// Though, the bitmask solution is not perfect since the bytes exceeding
	// the alphabet size are refused. Therefore, to reliably generate the ID,
	// the random bytes redundancy has to be satisfied.

	// Note: every hardware random generator call is performance expensive,
	// because the system call for entropy collection takes a lot of time.
	// So, to avoid additional system calls, extra bytes are requested in advance.

	// Next, a step determines how many random bytes to generate.
	// The number of random bytes gets decided upon the ID size, mask,
	// alphabet size, and magic number 1.6.

	//nolint:mnd // using 1.6 peaks at performance according to benchmarks
	stepFormula := 1.6 * float64(mask*opt.length) / float64(alphabetLen)
	step := int(math.Ceil(stepFormula))

	return generate(&generateInput{
		Option:      opt,
		alphabetLen: alphabetLen,
		mask:        mask,
		step:        step,
	})
}

func (n NanoID) String() string {
	if n == nil {
		return "<nil>"
	}
	return string(n)
}

type generateInput struct {
	*Option
	alphabetLen int
	mask        int
	step        int
}

func generate(opt *generateInput) (NanoID, error) {
	b := make([]byte, opt.step)
	id := make([]rune, opt.Option.length)

	currentRune := 0
	var currentAlphabetPosition int

	for {
		if err := opt.randomFunc(b); err != nil {
			return nil, err
		}

		for i := 0; i < opt.step; i++ {
			currentAlphabetPosition = int(b[i]) & opt.mask

			if currentAlphabetPosition >= opt.alphabetLen {
				continue
			}

			id[currentRune] = opt.alphabet[currentAlphabetPosition]
			currentRune++

			if currentRune == opt.Option.length {
				return NanoID(id), nil
			}
		}
	}
}
