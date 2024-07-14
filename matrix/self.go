package matrix

import (
	"errors"
)

func (A *Matrix[T]) Get(i, j int) T {
	return A.t[getIndex(A.n, A.m, i, j)]
}

func (A *Matrix[T]) GetSafe(i, j int) (result T, err error) {
	index := getIndex(A.n, A.m, i, j)
	if index < 0 || index >= A.n*A.m || i < 0 || i > A.n || j < 0 || j > A.m {
		err = errors.New("indexes out of range")
	} else {
		result = A.t[index]
	}
	return
}

func (A *Matrix[T]) Set(i, j int, value T) {
	A.t[getIndex(A.n, A.m, i, j)] = value
}

func (A *Matrix[T]) SetSafe(i, j int, value T) (err error) {
	index := getIndex(A.n, A.m, i, j)
	if index < 0 || index >= A.n*A.m {
		err = errors.New("index out of range")
	} else {
		A.t[index] = value
	}
	return
}

func (A *Matrix[T]) Add(B *Matrix[T]) (err error) {
	if A.n != B.n || A.m != B.m {
		err = errors.New("dimensions are not equal")
		return
	}
	for i := range A.t {
		A.t[i] += B.t[i]
	}
	return
}

func (A *Matrix[T]) Scale(scale T) {
	for i := range A.t {
		A.t[i] *= scale
	}
}

func (A *Matrix[T]) MultiplyLeft(B *Matrix[T]) error {
	C, err := Multiply(B, A)
	if err != nil {
		return err
	}

	A.t = C.t
	A.n = C.n
	A.m = C.m
	return nil
}

func (A *Matrix[T]) Multiply(B *Matrix[T]) error {
	C, err := Multiply(A, B)
	if err != nil {
		return err
	}

	A.t = C.t
	A.n = C.n
	A.m = C.m
	return nil
}

func (A *Matrix[T]) ApplyFunction(f func(_ T) T) {
	for i := range A.n * A.m {
		A.t[i] = f(A.t[i])
	}
}

func (A *Matrix[T]) IsEqualTo(B *Matrix[T]) bool {
	return AreEqual(A, B)
}
