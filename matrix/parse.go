package matrix

func (A *Matrix[T]) ToTable() [][]T {
	result := make([][]T, A.n)
	for i := range A.n {
		result[i] = make([]T, A.m)
		for j := range A.m {
			result[i][j] = A.Get(i, j)
		}
	}
	return result
}

func (A *Matrix[T]) ToRow() []T {
	return A.t
}
