package nanoid

func Parse[T any](_ T, _ ...func(*Option) *Option) (NanoID, error) {
	return NanoID{}, nil
}

func ParseFromString(_ string, _ ...func(*Option) *Option) (NanoID, error) {
	return NanoID{}, nil
}
