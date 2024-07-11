package nanoid

type Option struct {
	length     int
	alphabet   Alphabet
	randomFunc func([]byte) error
}

func WithLength(length int) func(*Option) *Option {
	return func(o *Option) *Option {
		o.length = length
		return o
	}
}

func WithAlphabet(alphabet Alphabet) func(*Option) *Option {
	return func(o *Option) *Option {
		o.alphabet = alphabet
		return o
	}
}

func WithRandomFunc(randomFunc func([]byte) error) func(*Option) *Option {
	return func(o *Option) *Option {
		o.randomFunc = randomFunc
		return o
	}
}
