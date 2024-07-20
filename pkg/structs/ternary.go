package structs

func Or[T any](clause bool, value1 T, value2 T) T {
	if clause {
		return value1
	}

	return value2
}
