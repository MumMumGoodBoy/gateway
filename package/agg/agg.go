package agg

func SortBySlice[K comparable, T any](keys []K, data []T, keyFunc func(T) K) []T {
	m := make(map[K]T)
	for _, d := range data {
		m[keyFunc(d)] = d
	}

	res := make([]T, 0, len(keys))
	for _, k := range keys {
		if d, ok := m[k]; ok {
			res = append(res, d)
		}
	}

	return res
}
