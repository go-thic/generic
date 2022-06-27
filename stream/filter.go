package stream

func Equals[V comparable](vs ...V) func(elem V) bool {
	lookup := make(map[V]bool)
	for _, v := range vs {
		lookup[v] = true
	}

	return func(elem V) bool {
		return lookup[elem]
	}
}
