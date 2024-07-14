package matrix

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

type Matrix[T Number] struct {
	n int
	m int
	t []T
}

func New[T Number](n, m int) *Matrix[T] {
	return &Matrix[T]{n, m, make([]T, n*m)}
}

func Initialize[T Number](n, m int, f func(i, j int) T) *Matrix[T] {
	A := New[T](n, m)
	for i := range A.n {
		for j := range A.m {
			A.Set(i, j, f(i, j))
		}
	}
	return A
}

func FromVector[T Number](t []T, column bool) (A *Matrix[T]) {
	if column {
		A = New[T](len(t), 1)
	} else {
		A = New[T](1, len(t))
	}

	A.t = t
	return
}

func FromTable[T Number](t [][]T) *Matrix[T] {
	if len(t) == 0 {
		return &Matrix[T]{0, 0, []T{}}
	}

	C := New[T](len(t), len(t[0]))
	for i := range len(t) {
		for j := range len(t[0]) {
			C.Set(i, j, t[i][j])
		}
	}
	return C
}

func FromMatrix[T Number](A *Matrix[T]) *Matrix[T] {
	B := New[T](A.n, A.m)
	B.t = make([]T, A.n*A.m)
	copy(B.t, A.t)
	return B
}
