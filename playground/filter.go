package playground

func Equals[V comparable](s ...V) func(elem VAL) bool {
	return func(elem VAL) bool {
		for _, t := range s {
			if t == elem {
				return true
			}
		}
		return false
	}
}
