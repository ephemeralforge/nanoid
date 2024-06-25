package nanoid

import "errors"

var (
	ErrInvalidBufferRead = errors.New("there is a problem reading random buffer with cripto/rand")
	ErrInvalidIDLength   = errors.New("the id length cannot be less than 2 or greater than 255")
)
