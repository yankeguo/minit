package gg

// Repeat create a slice with a value repeated
func Repeat[T any](count int, v T) []T {
	out := make([]T, count, count)
	for i := range out {
		out[i] = v
	}
	return out
}

// Ptr create a pointer to a value
func Ptr[T any](v T) *T {
	return &v
}
