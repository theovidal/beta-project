package matrix

import (
	"errors"
	"fmt"
)

func Add[T Number](A, B *Matrix[T]) (C *Matrix[T], err error) {
	if A.n != B.n || A.m != B.m {
		err = errors.New("dimensions are not equal")
		return
	}
	C = New[T](A.n, A.m)
	_ = C.Add(A)
	_ = C.Add(B)
	return
}

func Scale[T Number](A *Matrix[T], scale T) (C *Matrix[T]) {
	C = FromMatrix(A)
	C.Scale(scale)
	return
}

func Multiply[T Number](A, B *Matrix[T]) (C *Matrix[T], err error) {
	if A.m != B.n {
		err = errors.New(fmt.Sprintf("dimensions (%d, %d) and (%d, %d) don't match", A.n, A.m, B.n, B.m))
		return
	}

	C = New[T](A.n, B.m)

	for i := range A.n {
		for j := range B.m {
			var value T
			for k := range A.m {
				a := A.Get(i, k)
				b := B.Get(k, j)
				value += a * b
			}
			C.Set(i, j, value)
		}
	}
	return
}

func ApplyFunction[T Number](A *Matrix[T], f func(_ T) T) (C *Matrix[T]) {
	C = FromMatrix(A)
	C.ApplyFunction(f)
	return
}

func AreEqual[T Number](A, B *Matrix[T]) bool {
	if A.n != B.n || A.m != B.m {
		return false
	}
	for i := range A.n * A.m {
		if A.t[i] != B.t[i] {
			return false
		}
	}
	return true
}
