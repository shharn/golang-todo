package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type GetOffsetTestCase struct {
	input []int
	expected int
}

func TestGetOffset(t *testing.T) {
	tcs := []GetOffsetTestCase{
		0: GetOffsetTestCase{
			input: []int{1, 5},
			expected: 1,
		},
		1: GetOffsetTestCase{
			input: []int{2, 5},
			expected: 6,
		},
		2: GetOffsetTestCase{
			input: []int{5, 7},
			expected: 29,
		},
	}

	for _, tc := range tcs {
		actual := getOffset(tc.input[0], tc.input[1])
		assert.Equal(t, tc.expected, actual)
	}
}

type GetUpperLimitTestCase struct {
	input []int
	expected int
}

func TestGetUpperLimit(t *testing.T) {
	tcs := []GetUpperLimitTestCase{
		0: GetUpperLimitTestCase{
			input: []int{1,3,1},
			expected: 2,
		},
		1: GetUpperLimitTestCase{
			input: []int{1,3,2},
			expected: 3,
		},
		2: GetUpperLimitTestCase{
			input: []int{1,3,3},
			expected: 4,
		},
		3: GetUpperLimitTestCase{
			input: []int{1,3,4},
			expected: 4,
		},
		4: GetUpperLimitTestCase{
			input: []int{1,3,5},
			expected: 4,
		},
		5: GetUpperLimitTestCase{
			input: []int{2,1,1},
			expected: 2,
		},
		6: GetUpperLimitTestCase{
			input: []int{2,1,2},
			expected: 3,
		},
		7: GetUpperLimitTestCase{
			input: []int{2,1,3},
			expected: 3,
		},
		8: GetUpperLimitTestCase{
			input: []int{2,1,4},
			expected: 3,
		},
		9: GetUpperLimitTestCase{
			input: []int{2,1,0},
			expected: 1,
		},
	}
	for _, tc := range tcs {
		actual := getUpperLimit(tc.input[0], tc.input[1], tc.input[2])
		assert.Equal(t, tc.expected, actual)
	}
}