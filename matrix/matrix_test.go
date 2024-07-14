package matrix

import "testing"

func TestMatrix_Print(t *testing.T) {
	a := FromVector[float64]([]float64{4, 5, 6, 1, 4}, false)
	b := FromTable[float64]([][]float64{
		{4, 3, 0, 1},
		{8, 7, -5, -9},
		{3, 6, -7, -9},
		{-1, 0, 0, 4},
	})
	c := FromVector[float64]([]float64{4, 5, 6, 1, 4}, true)
	d := FromTable[float64]([][]float64{})
	e := FromTable[float64]([][]float64{})
	a.Print()
	b.Print()
	c.Print()
	d.Print()
	e.Print()
}

func TestMatrix_GetSafe(t *testing.T) {
	c := FromVector[float64]([]float64{4, 5, 6, 1, 4}, true)
	c.Print()

	if _, err := c.GetSafe(1, 2); err == nil {
		t.Error("index should be out of range")
	}

	if value, err := c.GetSafe(2, 0); err != nil || value != 6 {
		t.Error("coefficient should be 6")
	}
}

func TestMatrix_SetSafe(t *testing.T) {
	c := New[float64](5, 3)

	if c.SetSafe(5, 2, 6) == nil || c.SetSafe(4, 3, 4) == nil {
		t.Error("matrix dimensions should have been violated")
	}

	err := c.SetSafe(4, 2, 8)
	if err != nil {
		t.Error(err)
	}

	var v float64
	if v, err = c.GetSafe(4, 2); v != 8 || err != nil {
		t.Error("c.Set didn't set the value")
	}
}

func TestMatrix_ToTable(t *testing.T) {
	table := [][]float64{
		{4, 3, 0, 1},
		{8, 7, -5, -9},
		{3, 6, -7, -9},
		{-1, 0, 0, 4},
	}

	a := FromTable[float64]([][]float64{
		{4, 3, 0, 1},
		{8, 7, -5, -9},
		{3, 6, -7, -9},
		{-1, 0, 0, 4},
	})
	a.Set(1, 2, -4)
	table[1][2] = -4

	res := a.ToTable()
	for i := range len(table) {
		for j := range len(table[0]) {
			if res[i][j] != table[i][j] {
				t.Errorf("coefficients (%d, %d) don't match", i, j)
			}
		}
	}
}

func TestMatrix_IsEqualTo(t *testing.T) {
	table := [][]int{
		{4, 3, 0, 1},
		{8, 7, -5, -9},
		{3, 6, -7, -9},
		{-1, 0, 0, 4},
	}
	a := FromTable(table)
	b := FromTable(table)
	c := FromVector([]int{4, 1, 5, 9}, false)

	if a.IsEqualTo(c) {
		t.Error("a and c shouldn't be equal, as dimensions don't correspond")
	}

	if !a.IsEqualTo(b) || !b.IsEqualTo(a) {
		t.Error("matrixes should be equal in both ways")
	}

	a.Set(1, 1, -5)
	if AreEqual(a, b) {
		t.Error("a should be different to b after being altered")
	}
}

func TestMatrix_Scale(t *testing.T) {
	a := FromTable([][]int{
		{4, 3, 0, 1},
		{8, 7, -5, -9},
		{3, 6, -7, -9},
		{-1, 0, 0, 4},
	})
	b := FromTable([][]int{
		{8, 6, 0, 2},
		{16, 14, -10, -18},
		{6, 12, -14, -18},
		{-2, 0, 0, 8},
	})

	c := Scale(a, 2)
	a.Scale(2)

	if !a.IsEqualTo(b) || !c.IsEqualTo(b) {
		t.Error("wrong scaling")
	}
}

func TestMatrix_Add(t *testing.T) {
	a := FromTable([][]int{
		{4, 3, 0, 1},
		{8, 7, -5, -9},
		{3, 6, -7, -9},
		{-1, 0, 0, 4},
	})
	b := FromTable([][]int{
		{1, 1, 0, 0},
		{-4, -8, 7, 1},
		{2, 0, -5, -3},
		{1, 0, 6, 5},
	})
	e := FromTable([][]int{
		{5, 4, 1, 2},
	})

	if _, err := Add(b, e); b.Add(e) == nil || err == nil {
		t.Error("matrixes should not be summed, as dimensions don't correspond")
	}

	c, err := Add(a, b)
	if err != nil {
		t.Error(err)
	}
	if a.Add(b) != nil {
		t.Error(err)
	}

	d := FromTable([][]int{
		{5, 4, 0, 1},
		{4, -1, 2, -8},
		{5, 6, -12, -12},
		{0, 0, 6, 9},
	})

	if !c.IsEqualTo(d) || !a.IsEqualTo(d) {
		t.Error("addition is wrong")
	}
}

func TestMatrix_Multiply(t *testing.T) {
	a := FromTable([][]int{
		{4, 3, 0, 1},
		{8, 7, -5, -9},
		{3, 6, -7, -9},
	})
	b := FromTable([][]int{
		{1, 1},
		{-4, -8},
		{2, 0},
		{1, 0},
	})
	e := FromTable([][]int{
		{7, 8, 9, 6},
	})

	if _, err := Multiply(a, e); a.Multiply(e) == nil || err == nil {
		t.Error("matrixes should not be multiplied, as dimensions don't correspond")
	}

	c, err := Multiply(a, b)
	if err != nil {
		t.Error(err)
	}
	if a.Multiply(b) != nil {
		t.Error(err)
	}

	if c.n != 3 || c.m != 2 || a.n != 3 || a.m != 2 {
		t.Error("wrong matrix dimension")
	}

	d := FromTable([][]int{
		{-7, -20},
		{-39, -48},
		{-44, -45},
	})

	if !c.IsEqualTo(d) || !a.IsEqualTo(d) {
		t.Error("multiplication is wrong")
	}
}

func TestMatrix_ApplyFunction(t *testing.T) {
	a := FromTable([][]int{
		{4, 3, 0, 1},
		{8, 7, -5, -9},
		{3, 6, -7, -9},
		{-1, 0, 0, 4},
	})
	b := FromTable([][]int{
		{16, 9, 0, 1},
		{64, 49, 25, 81},
		{9, 36, 49, 81},
		{1, 0, 0, 16},
	})

	square := func(x int) int {
		return x * x
	}

	c := ApplyFunction(a, square)
	a.ApplyFunction(square)

	if !a.IsEqualTo(b) || !c.IsEqualTo(b) {
		t.Error("wrong function application")
	}
}
