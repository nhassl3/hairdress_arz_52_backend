package reverse_enums

func ReverseMap[K, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))

	for k, v := range m {
		result[v] = k
	}

	return result
}
