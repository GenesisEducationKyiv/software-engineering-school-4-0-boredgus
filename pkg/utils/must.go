package utils

func Must[T interface{}](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
