package main

import (
	"reflect"
	"testing"
)

type TestCases struct {
	pass []int
	tc   []int
	want []Color
}

func TestColors(t *testing.T) {
	testcases := []TestCases{
		{
			pass: []int{1, 7, 9, 9, 2},
			tc:   []int{1, 1, 9, 8, 1},
			want: []Color{Right, Wrong, Right, Wrong, Wrong},
		},
		{
			pass: []int{1, 1, 1, 1, 1},
			tc:   []int{1, 1, 1, 2, 1},
			want: []Color{Right, Right, Right, Wrong, Right},
		},
		{
			pass: []int{3, 2, 1, 1, 1},
			tc:   []int{1, 1, 1, 2, 1},
			want: []Color{White, Wrong, Right, White, Right},
		},
		{
			pass: []int{3, 2, 1, 1, 1},
			tc:   []int{1, 1, 3, 1, 2},
			want: []Color{White, White, White, Right, White},
		},
		{
			pass: []int{3, 2, 1, 4, 1},
			tc:   []int{1, 4, 3, 1, 2},
			want: []Color{White, White, White, White, White},
		},
		{
			pass: []int{3, 2, 1, 1, 1},
			tc:   []int{5, 5, 5, 5, 5},
			want: []Color{Wrong, Wrong, Wrong, Wrong, Wrong},
		},
		{
			pass: []int{5, 7, 5, 1, 5},
			tc:   []int{5, 7, 5, 1, 4},
			want: []Color{Right, Right, Right, Right, Wrong},
		},
		{
			pass: []int{5, 7, 5, 1, 5},
			tc:   []int{1, 1, 1, 2, 1},
			want: []Color{White, Wrong, Wrong, Wrong, Wrong},
		},
	}

	// Loop through test cases and check weird values with colors
	sz := 5
	for tn, x := range testcases {
		guess := NewGuess(sz)
		copy(guess.Num, x.tc)
		assignColors(guess, x.pass, sz)
		if !reflect.DeepEqual(guess.Color, x.want) {
			t.Errorf("Test: %d - Wanted: %v - Got: %v", tn, x.want, guess.Color)
		}
	}
}
