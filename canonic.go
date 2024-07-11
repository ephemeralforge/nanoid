package nanoid

import (
	"crypto/rand"
	"fmt"
)

const (
	CanonicNanoIDLenght   = 21
	CanonicNanoIDAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
)

func CanonicNanoIDRandomFunc(b []byte) error {
	if _, err := rand.Read(b); err != nil {
		return fmt.Errorf("nanoid: %s caused by (%w)", err.Error(), ErrInvalidBufferRead)
	}
	return nil
}
