package utils

// Must is a util that help with fail fast
func Must(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	return i
}
