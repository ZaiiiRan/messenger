package utils

func RemoveFromSlice[T comparable](slice []T, id T) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if v != id {
			result = append(result, v)
		}
	}
	return result
}
