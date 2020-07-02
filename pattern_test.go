package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewPatternFromRLE(t *testing.T) {
	tests := []struct {
		rle      string
		expected *Pattern
	}{
		{
			rle: `#C This is a glider.
x = 3, y = 3
bo$2bo$3o!`,
			expected: &Pattern{
				Cells: &Cells{
					cells: [][]bool{
						{false, true, false},
						{false, false, true},
						{true, true, true},
					},
					width:  3,
					height: 3,
				},
				name: "unknown",
				comment: []string{
					"This is a glider.",
				},
				credit: "",
			},
		},
		{
			rle: `#N Beehive
#O John Conway
#C An extremely common 6-cell still life.
#C www.conwaylife.com/wiki/index.php?title=Beehive
x = 4, y = 3, rule = B3/S23
b2ob$o2bo$b2o!`,
			expected: &Pattern{
				Cells: &Cells{
					cells: [][]bool{
						{false, true, true, false},
						{true, false, false, true},
						{false, true, true, false},
					},
					width:  4,
					height: 3,
				},
				name: "Beehive",
				comment: []string{
					"An extremely common 6-cell still life.",
					"www.conwaylife.com/wiki/index.php?title=Beehive",
				},
				credit: "John Conway",
			},
		},
	}

	for _, test := range tests {
		r := strings.NewReader(test.rle)
		pattern, err := NewPatternFromRLE(r)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(pattern, test.expected) {
			t.Errorf("\nwant\t%v\ngot\t%v", test.expected, pattern)
		}
	}
}
