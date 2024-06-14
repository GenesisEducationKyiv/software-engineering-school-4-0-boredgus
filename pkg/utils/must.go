package utils

// Panics if error is not nil.
func Must[T interface{}](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
