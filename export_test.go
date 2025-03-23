package ensure

// Export these internal functions so they can be accessed by the test suite

func IsEven(typeStr string, i any) bool {
	return isEven(typeStr, i)
}

func IsOdd(typeStr string, i any) bool {
	return isOdd(typeStr, i)
}

func NewTypeError(msg string) TypeError {
	return TypeError{err: msg}
}
