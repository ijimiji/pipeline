package ptr

func T[S any](val S) *S {
	return &val
}
