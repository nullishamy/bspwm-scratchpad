package main

func Remove[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func Contains[T comparable](l []T, item T) bool {
	for _, other := range l {
		if other == item {
			return true
		}
	}
	return false
}