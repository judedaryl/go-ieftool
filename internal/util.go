package internal

func remove[T any](slice []T, s int) []T {
	if len(slice) == 1 {
		return []T{}
	}
	newArr := []T{}
	for i := range slice {
		if i != s {
			newArr = append(newArr, slice[i])
		}
	}
	return newArr
}
