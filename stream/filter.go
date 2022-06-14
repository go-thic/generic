package stream

func Equals[V comparable](compare ...V) func(elem V) bool {
	lookup := make(map[V]bool)
	for _, v := range compare {
		lookup[v] = true
	}

	return func(elem V) bool {
		return lookup[elem]
	}
}
