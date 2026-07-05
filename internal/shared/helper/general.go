package helper

func GetDefault[T any](defaults []T) (fallback T) {
	if len(defaults) > 0 {
		fallback = defaults[0]
	}

	return
}
