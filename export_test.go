package ensure

// Export these internal functions so they can be access by the test suite

func IsEven[T NumberType](typeStr string, i T) bool {
	return isEven[T](typeStr, i)
}
