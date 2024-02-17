package slices

func Map[From any, To any](slice []From, mapper func(from From) To) []To {
	ret := make([]To, 0, len(slice))
	for _, x := range slice {
		ret = append(ret, mapper(x))
	}
	return ret
}
