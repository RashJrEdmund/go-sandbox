package main

func Map[T any, V any](
	arr []T,
	fn func(T, int) V,
) []V {
	res := make([]V, len(arr))

	for i, val := range arr {
		res[i] = fn(val, i)
	}

	return res
}
