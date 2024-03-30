package nanoid

type Alphabet []rune

func AlphabetFromString(s string) Alphabet {
	return []rune(s)
}
