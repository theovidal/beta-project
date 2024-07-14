package matrix

import "fmt"

var brackets = [][]rune{
	{'[', ']'},
	{'⎡', '⎤'},
	{'⎣', '⎦'},
	{'⎢', '⎥'},
}

func getBracket(n, i int, right bool) rune {
	var index int
	if right {
		index = 1
	}

	if n <= 1 {
		return brackets[0][index]
	} else {
		switch i {
		case 0:
			return brackets[1][index]
		case n - 1:
			return brackets[2][index]
		default:
			return brackets[3][index]
		}
	}
}

func (A *Matrix[T]) Print() {
	for i := range A.n {
		fmt.Printf("%c ", getBracket(A.n, i, false))
		for j := range A.m {
			fmt.Printf("%+v ", A.Get(i, j))
		}

		fmt.Printf("%c\n", getBracket(A.n, i, true))
	}
}
