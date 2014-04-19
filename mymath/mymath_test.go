package mymath

import "testing"

func TestAbs(t testing.T) {
	ans := Abs(-4)
	if ans != float64(4) {
		t.Errorf("Abs Error %f = %f", 4, ans)
	}
}
